package serve

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// minimalCV is the smallest body that passes schema.Validate.
const minimalCV = `{
  "template": "t2",
  "fileName": "jane",
  "head": {"fullName":"Jane Doe","address":"City","phone":"+1","email":"j@x.com"},
  "education": [{"school":"S","location":"L","startedAt":"2019","endedAt":"2021","degree":"BSc","description":"d"}],
  "jobs": [{"company":"C","location":"L","position":"Dev","startedAt":"2021","endedAt":"now","tools":["Go"],"highlights":["Did a thing"]}],
  "languages": [{"language":"English","level":"Native"}]
}`

func TestTemplatesEndpoint(t *testing.T) {
	s := &apiServer{cors: "*"}
	req := httptest.NewRequest(http.MethodGet, "/api/templates", nil)
	rec := httptest.NewRecorder()
	s.wrap(s.templates)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("CORS header = %q, want *", got)
	}
	var body struct {
		Templates []struct{ ID, Name string } `json:"templates"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Templates) != 3 {
		t.Errorf("got %d templates, want 3 (t1/t2/t3)", len(body.Templates))
	}
}

func TestGenerateEndpoint_OK(t *testing.T) {
	s := &apiServer{cors: "*"}
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(minimalCV))
	rec := httptest.NewRecorder()
	s.wrap(s.generate)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body=%s", rec.Code, rec.Body.String())
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/pdf" {
		t.Errorf("Content-Type = %q, want application/pdf", ct)
	}
	if !bytes.HasPrefix(rec.Body.Bytes(), []byte("%PDF")) {
		t.Errorf("body is not a PDF (first bytes: %q)", rec.Body.Bytes()[:min(8, rec.Body.Len())])
	}
	if fn := rec.Header().Get("X-Filename"); !strings.HasSuffix(fn, ".pdf") {
		t.Errorf("X-Filename = %q, want a .pdf name", fn)
	}
}

func TestGenerateEndpoint_ValidationError(t *testing.T) {
	s := &apiServer{cors: "*"}
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(`{"template":"t2","fileName":"x"}`))
	rec := httptest.NewRecorder()
	s.wrap(s.generate)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
	var body errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.Error == "" || len(body.Issues) == 0 {
		t.Errorf("expected an error with field issues, got %+v", body)
	}
}

func TestGenerateEndpoint_BadTemplate(t *testing.T) {
	s := &apiServer{cors: "*"}
	body := strings.Replace(minimalCV, `"template": "t2"`, `"template": "nope"`, 1)
	req := httptest.NewRequest(http.MethodPost, "/api/generate", strings.NewReader(body))
	rec := httptest.NewRecorder()
	s.wrap(s.generate)(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422; body=%s", rec.Code, rec.Body.String())
	}
}

func TestGenerateEndpoint_CORSPreflight(t *testing.T) {
	s := &apiServer{cors: "*"}
	req := httptest.NewRequest(http.MethodOptions, "/api/generate", nil)
	rec := httptest.NewRecorder()
	s.wrap(s.generate)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Errorf("OPTIONS preflight status = %d, want 204", rec.Code)
	}
}
