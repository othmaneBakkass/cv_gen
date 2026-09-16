package generate

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

// testCommand builds a fresh *cobra.Command with the same flags init()
// registers on the real one, so each test gets its own isolated
// Flags().Changed() state instead of sharing (and fighting over) the
// package-level command.
func testCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "generate"}
	cmd.Flags().Bool("french", false, "")
	cmd.Flags().String("density", "normal", "")
	cmd.Flags().String("template", "", "")
	for _, t := range sectionToggleFlags {
		cmd.Flags().Bool(t.flag, false, "")
	}
	for _, p := range sectionPriorityFlags {
		cmd.Flags().Int(p.flag, 0, "")
	}
	return cmd
}

func TestOptionsFromSettings_Nil(t *testing.T) {
	opts, err := optionsFromSettings(nil)
	if err != nil {
		t.Fatalf("optionsFromSettings(nil): %v", err)
	}
	if opts.Lang != i18n.EN {
		t.Errorf("nil settings should default to English, got %v", opts.Lang)
	}
	if opts.Density != 1.0 {
		t.Errorf("nil settings should default to density 1.0, got %v", opts.Density)
	}
}

func TestOptionsFromSettings_Full(t *testing.T) {
	include := true
	priority := 1
	s := &schema.Settings{
		Lang:    "fr",
		Density: "airy",
		Sections: &schema.SectionSettings{
			Skills:   &schema.SectionSetting{Priority: &priority},
			Projects: &schema.SectionSetting{Include: &include},
		},
		Colors: &schema.ColorSettings{Accent: "#8B0000"},
	}

	opts, err := optionsFromSettings(s)
	if err != nil {
		t.Fatalf("optionsFromSettings: %v", err)
	}
	if opts.Lang != i18n.FR {
		t.Errorf("Lang = %v, want fr", opts.Lang)
	}
	if opts.Density != 1.15 {
		t.Errorf("Density = %v, want 1.15 (airy)", opts.Density)
	}
	if got := opts.PriorityOf(render.Skills, 999); got != 1 {
		t.Errorf("Skills priority = %d, want 1", got)
	}
	if include, ok := opts.IncludeOverride(render.Projects); !ok || !include {
		t.Errorf("Projects include override = (%v, %v), want (true, true)", include, ok)
	}
	// The deprecated "accent" alias fans out to both the Headline and Border
	// semantic roles (accent historically drove headings and rules alike).
	want := theme.RGB{R: 139, G: 0, B: 0}
	if opts.Colors.Headline == nil || *opts.Colors.Headline != want {
		t.Errorf("Colors.Headline = %v, want %+v (from accent alias)", opts.Colors.Headline, want)
	}
	if opts.Colors.Border == nil || *opts.Colors.Border != want {
		t.Errorf("Colors.Border = %v, want %+v (from accent alias)", opts.Colors.Border, want)
	}
}

// TestSectionStyleFromSettings checks the per-section style block in JSON
// (settings.sections.<key>.style) maps onto opts.SectionStyle: a divider
// toggle, a spaceBefore, an uppercase flag, and a per-section color override.
func TestSectionStyleFromSettings(t *testing.T) {
	no := false
	space := 20.0
	s := &schema.Settings{Sections: &schema.SectionSettings{
		Experience: &schema.SectionSetting{Style: &schema.SectionStyle{
			Divider:     &no,
			SpaceBefore: &space,
			Colors:      &schema.ColorSettings{Headline: "#FF00FF"},
		}},
	}}

	opts, err := optionsFromSettings(s)
	if err != nil {
		t.Fatalf("optionsFromSettings: %v", err)
	}
	ov, ok := opts.SectionStyle[render.Experience]
	if !ok {
		t.Fatal("no SectionStyle recorded for experience")
	}
	if ov.ShowDivider == nil || *ov.ShowDivider != false {
		t.Errorf("ShowDivider = %v, want false", ov.ShowDivider)
	}
	if ov.SpaceBefore == nil || *ov.SpaceBefore != 20.0 {
		t.Errorf("SpaceBefore = %v, want 20", ov.SpaceBefore)
	}
	want := theme.RGB{R: 255, B: 255}
	if ov.Headline == nil || *ov.Headline != want {
		t.Errorf("Headline = %v, want %+v", ov.Headline, want)
	}
}

// TestColorSettings_SemanticWinsOverAlias proves that when both a deprecated
// alias and its semantic equivalent are supplied, the semantic value wins.
func TestColorSettings_SemanticWinsOverAlias(t *testing.T) {
	s := &schema.Settings{Colors: &schema.ColorSettings{
		Ink:  "#111111", // alias for body
		Body: "#222222", // semantic — should win
	}}
	opts, err := optionsFromSettings(s)
	if err != nil {
		t.Fatalf("optionsFromSettings: %v", err)
	}
	want := theme.RGB{R: 0x22, G: 0x22, B: 0x22}
	if opts.Colors.Body == nil || *opts.Colors.Body != want {
		t.Errorf("Colors.Body = %v, want %+v (semantic key should beat the ink alias)", opts.Colors.Body, want)
	}
}

func TestOptionsFromSettings_InvalidDensity(t *testing.T) {
	s := &schema.Settings{Density: "extremely-airy"}
	if _, err := optionsFromSettings(s); err == nil {
		t.Error("expected an error for an invalid settings.density value")
	}
}

func TestOptionsFromSettings_InvalidColor(t *testing.T) {
	s := &schema.Settings{Colors: &schema.ColorSettings{Accent: "not-a-color"}}
	if _, err := optionsFromSettings(s); err == nil {
		t.Error("expected an error for an invalid settings.colors.accent value")
	}
}

// TestCLIOverridesJSON is the precedence contract: a JSON settings block
// provides the base, and only a flag the caller actually typed
// (Flags().Changed) may override it.
func TestCLIOverridesJSON(t *testing.T) {
	s := &schema.Settings{Lang: "fr"}
	base, err := optionsFromSettings(s)
	if err != nil {
		t.Fatalf("optionsFromSettings: %v", err)
	}

	cmd := testCommand()
	// french flag left untouched: JSON's "fr" must survive.
	merged, err := applyCLIOverrides(cmd, base)
	if err != nil {
		t.Fatalf("applyCLIOverrides: %v", err)
	}
	if merged.Lang != i18n.FR {
		t.Errorf("an unset --french flag should not clobber JSON's lang=fr, got %v", merged.Lang)
	}

	// Now explicitly pass --french=false: the flag was typed, so it wins.
	cmd2 := testCommand()
	if err := cmd2.Flags().Set("french", "false"); err != nil {
		t.Fatal(err)
	}
	merged2, err := applyCLIOverrides(cmd2, base)
	if err != nil {
		t.Fatalf("applyCLIOverrides: %v", err)
	}
	if merged2.Lang != i18n.EN {
		t.Errorf("an explicitly-set --french=false should override JSON's lang=fr, got %v", merged2.Lang)
	}
}

func TestCLIOverrides_SectionToggleAndPriority(t *testing.T) {
	base := render.Default()
	cmd := testCommand()
	if err := cmd.Flags().Set("no-skills", "true"); err != nil {
		t.Fatal(err)
	}
	if err := cmd.Flags().Set("experience", "1"); err != nil {
		t.Fatal(err)
	}

	merged, err := applyCLIOverrides(cmd, base)
	if err != nil {
		t.Fatalf("applyCLIOverrides: %v", err)
	}
	if include, ok := merged.IncludeOverride(render.Skills); !ok || include {
		t.Errorf("--no-skills should force Skills off, got (%v, %v)", include, ok)
	}
	if got := merged.PriorityOf(render.Experience, 30); got != 1 {
		t.Errorf("--experience=1 should set priority to 1, got %d", got)
	}
	// Untouched section flags must not appear as overrides at all.
	if _, ok := merged.IncludeOverride(render.Projects); ok {
		t.Error("an untouched --no-projects flag should not create an Include entry")
	}
}
