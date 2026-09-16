// Package serve adds the `serve` subcommand: an HTTP server that exposes
// the same CV engine the CLI uses over a small JSON API, so the web
// frontend can generate real PDFs instead of only exporting JSON. It reuses
// schema.Validate, settings.FromSettings and templates.RenderBytes verbatim
// — there is no second rendering path.
package serve

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/othmaneBakkass/cv_gen/cmd/root"
	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/common/logs"
	"github.com/othmaneBakkass/cv_gen/internal/fsc"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/templates"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
	"github.com/othmaneBakkass/cv_gen/internal/settings"
)

const maxBodyBytes = 2 << 20 // 2 MiB — a CV JSON is tiny; cap to be safe.

var command = &cobra.Command{
	Use:   "serve",
	Short: "Run the CV generator as an HTTP server.",
	Long:  "Serve the CV engine over a JSON API (POST /api/generate, GET /api/templates), optionally hosting the built web frontend from --static.",
	RunE:  handler,
}

func init() {
	command.Flags().String("addr", ":8080", "Address to listen on, e.g. :8080 or 127.0.0.1:8080.")
	command.Flags().String("static", "", "Optional directory of a built frontend (e.g. web/dist) to serve alongside the API, with SPA fallback.")
	command.Flags().String("cors", "*", "Value for Access-Control-Allow-Origin on /api routes (use a specific origin in production; empty disables CORS headers).")
	root.RootCommand.AddCommand(command)
}

func handler(cmd *cobra.Command, _ []string) error {
	addr, _ := cmd.Flags().GetString("addr")
	static, _ := cmd.Flags().GetString("static")
	cors, _ := cmd.Flags().GetString("cors")

	mux := http.NewServeMux()
	api := &apiServer{cors: cors}
	mux.HandleFunc("/api/templates", api.wrap(api.templates))
	mux.HandleFunc("/api/generate", api.wrap(api.generate))
	mux.HandleFunc("/api/health", api.wrap(api.health))

	if static != "" {
		if _, err := os.Stat(static); err != nil {
			return apperror.New("Invalid --static directory", err.Error(), apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
		}
		mux.Handle("/", spaHandler(static))
		fmt.Println(logs.InfoLog("Serving frontend from " + static))
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}
	fmt.Println(logs.SuccessLog("cv_gen listening on " + addr + " (POST /api/generate, GET /api/templates)"))
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return apperror.New("Server failed", err.Error(), apperror.ErrorCodeUnknown, apperror.ErrorSensitivityPublic)
	}
	return nil
}

type apiServer struct {
	cors string
}

// wrap applies CORS headers and OPTIONS preflight handling around an API
// handler, and logs each request.
func (s *apiServer) wrap(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.cors != "" {
			w.Header().Set("Access-Control-Allow-Origin", s.cors)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

func (s *apiServer) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *apiServer) templates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "GET only")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"templates": templates.Metadata()})
}

// generateRequest is the POST /api/generate body: a single CV (the same
// shape as one entry of the CLI's data[] array, settings included). A
// wrapper {"data":[...]} is also accepted for parity with the CLI's input
// files, in which case the first entry is used.
type generateRequest struct {
	schema.CV
	Data []schema.CV `json:"data,omitempty"`
}

func (s *apiServer) generate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "POST only")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	var req generateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON", err.Error())
		return
	}

	cv := req.CV
	if len(req.Data) > 0 {
		cv = req.Data[0]
	}

	// A photo path in the JSON is resolved relative to the input file on the
	// CLI; over HTTP there's no such file, and an unresolvable path would
	// fail the render. Drop it — photo upload is a separate future feature.
	cv.Head.Photo = ""

	if err := schema.Validate(&cv); err != nil {
		writeAppError(w, http.StatusBadRequest, err)
		return
	}

	opts, err := settings.FromSettings(cv.Settings)
	if err != nil {
		writeAppError(w, http.StatusBadRequest, err)
		return
	}

	pdf, err := templates.RenderBytes(cv, opts)
	if err != nil {
		writeAppError(w, http.StatusUnprocessableEntity, err)
		return
	}

	fileName := fsc.EnsureFileName(cv.FileName, "cv", "pdf")
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", fileName))
	w.Header().Set("X-Filename", fileName)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(pdf)
}

// spaHandler serves files from dir, falling back to index.html for any path
// that doesn't resolve to a real file so client-side routing works.
func spaHandler(dir string) http.Handler {
	fileServer := http.FileServer(http.Dir(dir))
	index := filepath.Join(dir, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clean := filepath.Clean(r.URL.Path)
		path := filepath.Join(dir, clean)
		// Guard against path traversal outside dir.
		if !strings.HasPrefix(path, filepath.Clean(dir)) {
			http.NotFound(w, r)
			return
		}
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}
		http.ServeFile(w, r, index)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorBody struct {
	Error  string   `json:"error"`
	Detail string   `json:"detail,omitempty"`
	Issues []string `json:"issues,omitempty"`
}

func writeError(w http.ResponseWriter, status int, title, detail string) {
	writeJSON(w, status, errorBody{Error: title, Detail: detail})
}

// writeAppError renders an apperror (from schema.Validate / settings / the
// renderer) as a JSON error, surfacing only public issues.
func writeAppError(w http.ResponseWriter, status int, err error) {
	var appErr apperror.AppError
	if !errors.As(err, &appErr) {
		writeError(w, status, "Request failed", err.Error())
		return
	}
	body := errorBody{Error: appErr.Title}
	if appErr.Sensitivity == apperror.ErrorSensitivityPublic {
		body.Detail = appErr.Detail
	}
	for _, issue := range appErr.Issues {
		if issue.Sensitivity == apperror.ErrorSensitivityPublic {
			body.Issues = append(body.Issues, fmt.Sprintf("%s: %s", issue.Title, issue.Detail))
		}
	}
	writeJSON(w, status, body)
}
