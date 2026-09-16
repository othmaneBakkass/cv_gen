package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

func minimalCV(template string) schema.CV {
	return schema.CV{
		Template: template,
		FileName: "test",
		Head:     schema.Head{FullName: "Test User", Address: "City", Phone: "000", Email: "a@b.com"},
		Education: []schema.Education{
			{School: "S", Location: "L", StartedAt: "2020", EndedAt: "2021", Degree: "D", Description: "Desc"},
		},
		Jobs: []schema.Job{
			{Company: "C", Location: "L", Position: "P", StartedAt: "2021", EndedAt: "2022", Tools: []string{"Go"}, Highlights: []string{"Did a thing."}},
		},
		Languages: []schema.Language{{Language: "English", Level: "Fluent"}},
	}
}

func TestNames(t *testing.T) {
	got := Names()
	want := []string{"t1", "t2", "t3"}
	if len(got) != len(want) {
		t.Fatalf("Names() = %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Names()[%d] = %q, want %q (must stay sorted)", i, got[i], want[i])
		}
	}
}

// TestGenerate_AllTemplates renders every registered template against the
// same minimal CV and checks the output is a real, non-trivial PDF. It's
// intentionally not pinned to exact byte counts or coordinates — the
// per-theme numeric checks live in theme_test.go/options_test.go/
// fit_test.go; this just guards against a template silently returning
// empty or garbage bytes.
func TestGenerate_AllTemplates(t *testing.T) {
	for _, name := range Names() {
		t.Run(name, func(t *testing.T) {
			cv := minimalCV(name)
			outPath := filepath.Join(t.TempDir(), "out.pdf")
			if err := Generate(outPath, cv, render.Default()); err != nil {
				t.Fatalf("Generate(%s): %v", name, err)
			}
			data, err := os.ReadFile(outPath)
			if err != nil {
				t.Fatalf("read output: %v", err)
			}
			if !strings.HasPrefix(string(data), "%PDF-") {
				t.Errorf("output doesn't start with a PDF header: %q", data[:min(20, len(data))])
			}
			if len(data) < 500 {
				t.Errorf("output suspiciously small (%d bytes) for a rendered CV", len(data))
			}
		})
	}
}

func TestGenerate_UnsupportedTemplate(t *testing.T) {
	cv := minimalCV("bogus")
	err := Generate(filepath.Join(t.TempDir(), "out.pdf"), cv, render.Default())
	if err == nil {
		t.Error("expected an error for an unsupported template name")
	}
}
