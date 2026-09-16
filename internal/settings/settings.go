// Package settings converts a CV's JSON "settings" block (schema.Settings)
// into a render.Options. It is the single place that maps the user-facing
// presentation options — language, density, per-section visibility/order,
// and the semantic (and legacy-alias) color/style overrides — onto the
// engine's render types, so the CLI (cmd/generate) and the HTTP server
// (cmd/serve) produce identical options from the same input.
package settings

import (
	"fmt"

	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

const titleBadInput = "Invalid input value"

// DensityMultipliers maps the density setting's allowed values to the
// spacing/margin scale factor passed to render.Options.Density. "normal" is
// each template's own base rhythm; "dense"/"airy" match the two density
// variants in the reference mockups (Templates.claude/CV_Template.pdf vs.
// CV_Template_Airy.pdf). Exported so the CLI's --density flag reuses the
// exact same mapping.
var DensityMultipliers = map[string]float64{
	"dense":  0.85,
	"normal": 1.0,
	"airy":   1.15,
}

// FromSettings converts a CV entry's JSON "settings" block into its base
// render.Options. A nil s (no "settings" key) returns render.Default().
// Callers layer their own overrides (e.g. the CLI's flags) on top.
func FromSettings(s *schema.Settings) (render.Options, error) {
	opts := render.Default()
	if s == nil {
		return opts, nil
	}

	if s.Lang == string(i18n.FR) {
		opts.Lang = i18n.FR
	}

	if s.Density != "" {
		mult, ok := DensityMultipliers[s.Density]
		if !ok {
			return opts, badInput("settings.density must be one of: dense, normal, airy")
		}
		opts.Density = mult
	}

	if s.Sections != nil {
		for _, e := range []struct {
			key render.SectionKey
			s   *schema.SectionSetting
		}{
			{render.Profile, s.Sections.Profile},
			{render.Education, s.Sections.Education},
			{render.Experience, s.Sections.Experience},
			{render.Skills, s.Sections.Skills},
			{render.Projects, s.Sections.Projects},
			{render.Certifications, s.Sections.Certifications},
			{render.Languages, s.Sections.Languages},
		} {
			if e.s == nil {
				continue
			}
			if e.s.Include != nil {
				opts.Include[e.key] = *e.s.Include
			}
			if e.s.Priority != nil {
				opts.Priority[e.key] = *e.s.Priority
			}
			if e.s.Style != nil {
				ov, err := sectionOverride(e.s.Style)
				if err != nil {
					return opts, err
				}
				opts.SectionStyle[e.key] = ov
			}
		}
	}

	if s.Colors != nil {
		colors, err := ColorOverride(s.Colors)
		if err != nil {
			return opts, err
		}
		opts.Colors = colors
	}

	return opts, nil
}

// ColorOverride parses a ColorSettings' hex strings into a
// render.ColorOverride, erroring out with the offending field's name if any
// hex value is malformed. Both the current semantic role names and the
// deprecated four-color aliases (accent/grey/muted/ink) are accepted: the
// aliases are applied first and the semantic keys second, so if a file sets
// both, the semantic value wins. The alias mapping preserves the old
// behavior — accent drove both the headings and the rules, so it maps onto
// Headline and Border together.
func ColorOverride(c *schema.ColorSettings) (render.ColorOverride, error) {
	var out render.ColorOverride
	// Order matters: deprecated aliases first, semantic roles after (later
	// bindings overwrite earlier ones for the same destination pointer).
	bindings := []struct {
		name string
		hex  string
		dsts []**theme.RGB
	}{
		{"accent", c.Accent, []**theme.RGB{&out.Headline, &out.Border}},
		{"grey", c.Grey, []**theme.RGB{&out.Metadata}},
		{"muted", c.Muted, []**theme.RGB{&out.Profile}},
		{"ink", c.Ink, []**theme.RGB{&out.Body}},

		{"headline", c.Headline, []**theme.RGB{&out.Headline}},
		{"subheadline", c.Subheadline, []**theme.RGB{&out.Subheadline}},
		{"body", c.Body, []**theme.RGB{&out.Body}},
		{"metadata", c.Metadata, []**theme.RGB{&out.Metadata}},
		{"link", c.Link, []**theme.RGB{&out.Link}},
		{"border", c.Border, []**theme.RGB{&out.Border}},
		{"profile", c.Profile, []**theme.RGB{&out.Profile}},
		{"background", c.Background, []**theme.RGB{&out.Background}},
	}
	for _, b := range bindings {
		if b.hex == "" {
			continue
		}
		rgb, err := theme.ParseHex(b.hex)
		if err != nil {
			return out, badInput(fmt.Sprintf("settings.colors.%s: %s", b.name, err.Error()))
		}
		r := rgb
		for _, d := range b.dsts {
			*d = &r
		}
	}
	return out, nil
}

// sectionOverride converts one section's JSON style block into the
// theme.SectionOverride RunSections applies for that section. Colors reuse
// ColorOverride (and its alias handling), so per-section overrides accept
// the same keys as the global palette.
func sectionOverride(st *schema.SectionStyle) (theme.SectionOverride, error) {
	var ov theme.SectionOverride
	ov.ShowDivider = st.Divider
	ov.SpaceBefore = st.SpaceBefore
	ov.HeadingUppercase = st.Uppercase
	if st.Colors != nil {
		co, err := ColorOverride(st.Colors)
		if err != nil {
			return ov, err
		}
		ov.Headline = co.Headline
		ov.Subheadline = co.Subheadline
		ov.Body = co.Body
		ov.Metadata = co.Metadata
		ov.Link = co.Link
		ov.Border = co.Border
		ov.Profile = co.Profile
	}
	return ov, nil
}

func badInput(detail string) error {
	return apperror.New(titleBadInput, detail, apperror.ErrorCodeArgs, apperror.ErrorSensitivityPublic)
}
