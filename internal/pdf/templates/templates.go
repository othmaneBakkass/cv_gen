// Package templates maps a template name to its renderer and turns a
// validated CV into a PDF file on disk.
package templates

import (
	"fmt"
	"os"
	"sort"

	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/templates/t1"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/templates/t2"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/templates/t3"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

// registry maps a template name to the function that renders it to bytes.
// Add a new entry here to support another layout. All three are built on
// the shared internal/pdf/draw+comp+sections layers and differ only by
// theme (internal/pdf/theme) and section order — see docs/pdf-rewrite-handoff.md.
var registry = map[string]func(schema.CV, render.Options) ([]byte, error){
	"t1": t1.Render, // ATS Classic — conservative, Morocco/French hiring convention
	"t2": t2.Render, // Navy Rule — editorial, teal-accented
	"t3": t3.Render, // Modern Minimal — compact, indigo-accented, no heading rules
}

// Meta describes a template for UIs that list the available layouts (e.g.
// the web frontend's /api/templates endpoint). Kept next to the registry so
// the two never drift.
type Meta struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var meta = map[string]Meta{
	"t1": {ID: "t1", Name: "ATS Classic", Description: "Conservative navy, single-column, upright dates. Built to parse cleanly through applicant-tracking software and read as a sober professional CV."},
	"t2": {ID: "t2", Name: "Navy Rule", Description: "Editorial and teal-accented, with tracked uppercase headings, thin rules under every section, and italic right-aligned dates."},
	"t3": {ID: "t3", Name: "Modern Minimal", Description: "Compact and indigo-accented, with title-case headings, no rule lines, and an en-dash bullet — a contemporary, quieter look."},
}

// Metadata returns descriptive metadata for every template, sorted by ID.
func Metadata() []Meta {
	out := make([]Meta, 0, len(meta))
	for _, name := range Names() {
		if m, ok := meta[name]; ok {
			out = append(out, m)
		}
	}
	return out
}

// Names returns the supported template names, sorted.
func Names() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// RenderBytes renders cv with its template under opts and returns the PDF
// bytes. Used by the HTTP server (cmd/serve), which streams the result back
// rather than writing it to disk; Generate builds on it for the CLI.
func RenderBytes(cv schema.CV, opts render.Options) ([]byte, error) {
	renderFn, ok := registry[cv.Template]
	if !ok {
		return nil, apperror.New(
			"Unsupported template",
			fmt.Sprintf("template %q is not supported (available: %v)", cv.Template, Names()),
			apperror.ErrorCodeArgs,
			apperror.ErrorSensitivityPublic,
		)
	}

	pdfBytes, err := renderFn(cv, opts)
	if err != nil {
		return nil, apperror.New(
			"PDF generation failed",
			err.Error(),
			apperror.ErrorCodeUnknown,
			apperror.ErrorSensitivityPublic,
		)
	}
	return pdfBytes, nil
}

// Generate renders cv with its template under opts and writes the PDF to
// outPath.
func Generate(outPath string, cv schema.CV, opts render.Options) error {
	pdfBytes, err := RenderBytes(cv, opts)
	if err != nil {
		return err
	}

	if err := os.WriteFile(outPath, pdfBytes, 0644); err != nil {
		return apperror.New(
			"PDF save failed",
			err.Error(),
			apperror.ErrorCodeUnknown,
			apperror.ErrorSensitivityPublic,
		)
	}

	return nil
}
