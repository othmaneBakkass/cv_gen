package render

import (
	"math"
	"testing"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

func TestPriorityOf(t *testing.T) {
	opts := Default()
	opts.Priority[Skills] = 5

	if got := opts.PriorityOf(Skills, 40); got != 5 {
		t.Errorf("PriorityOf(Skills, 40) = %d, want 5 (explicit override should win)", got)
	}
	if got := opts.PriorityOf(Languages, 70); got != 70 {
		t.Errorf("PriorityOf(Languages, 70) = %d, want 70 (no override, should fall back to def)", got)
	}
}

func TestIncludeOverride(t *testing.T) {
	opts := Default()
	opts.Include[Profile] = false

	if include, ok := opts.IncludeOverride(Profile); !ok || include {
		t.Errorf("IncludeOverride(Profile) = (%v, %v), want (false, true)", include, ok)
	}
	if _, ok := opts.IncludeOverride(Projects); ok {
		t.Errorf("IncludeOverride(Projects) reported an override that was never set")
	}
}

func TestEffectiveSpacing(t *testing.T) {
	base := theme.Spacing{XS: 1, SM: 2, MD: 3, LG: 4}

	opts := Default() // Density defaults to 1.0
	if got := opts.EffectiveSpacing(base); got != base {
		t.Errorf("EffectiveSpacing with default density = %+v, want unscaled %+v", got, base)
	}

	opts.Density = 0.5
	want := theme.Spacing{XS: 0.5, SM: 1, MD: 1.5, LG: 2}
	if got := opts.EffectiveSpacing(base); got != want {
		t.Errorf("EffectiveSpacing at 0.5 density = %+v, want %+v", got, want)
	}

	var zero Options
	if got := zero.EffectiveSpacing(base); got != base {
		t.Errorf("EffectiveSpacing with zero-value Options (Density=0) = %+v, want unscaled %+v (0 means 1.0)", got, base)
	}
}

func TestEffectiveMargins(t *testing.T) {
	base := theme.Theme{MarginTop: 24, MarginBottom: 24, MarginLeft: 42.5, MarginRight: 42.5, HeadingSize: 11}

	close := func(got, want float64) bool { return math.Abs(got-want) < 1e-9 }
	check := func(name string, th theme.Theme, k float64) {
		if !close(th.MarginTop, 24*k) || !close(th.MarginBottom, 24*k) ||
			!close(th.MarginLeft, 42.5*k) || !close(th.MarginRight, 42.5*k) {
			t.Errorf("%s: margins = top %v/bottom %v/left %v/right %v, want scale %v of 24/24/42.5/42.5",
				name, th.MarginTop, th.MarginBottom, th.MarginLeft, th.MarginRight, k)
		}
	}

	// Default density (1.0) leaves margins untouched; airy widens, dense narrows.
	check("default", Default().EffectiveMargins(base), 1.0)
	check("airy", Options{Density: 1.2}.EffectiveMargins(base), 1.2)
	check("dense", Options{Density: 0.85}.EffectiveMargins(base), 0.85)
	// Zero-value Density (0) means 1.0, not "scale to zero margins".
	check("zero-density", Options{}.EffectiveMargins(base), 1.0)

	// Non-margin fields must pass through untouched.
	if got := (Options{}).EffectiveMargins(base); got.HeadingSize != 11 {
		t.Errorf("EffectiveMargins altered a non-margin field: HeadingSize = %v, want 11", got.HeadingSize)
	}
}

func TestApplyTheme(t *testing.T) {
	base := theme.T2
	red := theme.RGB{R: 255}

	var opts Options
	opts.Colors.Headline = &red
	got := opts.ApplyTheme(base)

	if got.Palette.Headline != red {
		t.Errorf("ApplyTheme: Palette.Headline = %+v, want override %+v", got.Palette.Headline, red)
	}
	// Every other palette role must be untouched.
	if got.Palette.Body != base.Palette.Body || got.Palette.Metadata != base.Palette.Metadata ||
		got.Palette.Profile != base.Palette.Profile || got.Palette.Border != base.Palette.Border ||
		got.Palette.Subheadline != base.Palette.Subheadline || got.Palette.Link != base.Palette.Link {
		t.Errorf("ApplyTheme changed a color that had no override: got %+v, base %+v", got.Palette, base.Palette)
	}
	// Non-color fields must pass through untouched.
	if got.HeadingSize != base.HeadingSize || got.FontFamily != base.FontFamily {
		t.Errorf("ApplyTheme altered non-color fields: got %+v, base %+v", got, base)
	}
}

// TestSectionOverride_ThemeWith checks the middle layer of the styling
// hierarchy: a SectionOverride merged onto a theme changes only what it
// names, including folding ShowDivider=false into a zero rule weight.
func TestSectionOverride_ThemeWith(t *testing.T) {
	base := theme.T2
	green := theme.RGB{G: 128}
	no := false

	got := base.With(theme.SectionOverride{Body: &green, ShowDivider: &no})

	if got.Palette.Body != green {
		t.Errorf("With: Palette.Body = %+v, want %+v", got.Palette.Body, green)
	}
	if got.SectionRuleWeight != 0 {
		t.Errorf("With: ShowDivider=false should zero SectionRuleWeight, got %v", got.SectionRuleWeight)
	}
	// Untouched fields survive, and the base theme is not mutated.
	if got.Palette.Headline != base.Palette.Headline {
		t.Errorf("With mutated an unset role: Headline changed to %+v", got.Palette.Headline)
	}
	if base.SectionRuleWeight == 0 {
		t.Error("With mutated the receiver theme's SectionRuleWeight (must return a copy)")
	}
}
