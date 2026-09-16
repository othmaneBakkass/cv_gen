// Package render holds per-generation settings that are orthogonal to the
// CV data itself: which sections to show, what order to show them in, and
// what language the template's own labels render in. See
// docs/pdf-rewrite-handoff.md §6 for the reasoning.
package render

import (
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// SectionKey identifies one of the CV's optional/orderable sections.
type SectionKey string

const (
	Profile        SectionKey = "profile"
	Education      SectionKey = "education"
	Experience     SectionKey = "experience"
	Skills         SectionKey = "skills"
	Projects       SectionKey = "projects"
	Certifications SectionKey = "certifications"
	Languages      SectionKey = "languages"
)

// Options is the full set of per-render settings.
type Options struct {
	Lang i18n.Lang

	// Include holds explicit include/exclude overrides. A section absent
	// from this map is included exactly when its data is non-empty
	// (nil = infer from data, per docs/pdf-rewrite-handoff.md §6.1).
	Include map[SectionKey]bool

	// Priority holds explicit order overrides; lower renders first. A
	// section absent from this map falls back to defaultPriority.
	Priority map[SectionKey]int

	// Density scales the base spacing tokens — e.g. 0.85 for the denser
	// reference mockup, 1.15 for the airier one. 0 means 1.0 (unscaled).
	Density float64

	// Colors overrides the active template's global palette. A nil field
	// means "keep the template's default for that color".
	Colors ColorOverride

	// SectionStyle holds per-section style overrides keyed by section. These
	// are the user-facing (JSON/CLI) layer of Global Template Defaults ->
	// Section Overrides -> Rendered Section, applied on top of the
	// template's own Section.Style in RunSections. A missing key means the
	// section keeps the template's styling.
	SectionStyle map[SectionKey]theme.SectionOverride
}

// ColorOverride holds optional replacements for the theme's semantic
// palette roles. A nil field leaves the template's own value in place. It
// mirrors theme.Palette's roles; the older accent/grey/muted/ink names are
// accepted at the JSON/CLI boundary (see cmd/generate) and mapped onto
// these before reaching here.
type ColorOverride struct {
	Headline    *theme.RGB
	Subheadline *theme.RGB
	Body        *theme.RGB
	Metadata    *theme.RGB
	Link        *theme.RGB
	Border      *theme.RGB
	Profile     *theme.RGB
	Background  *theme.RGB
}

// ApplyTheme returns th with any set global Colors overrides applied. Every
// template calls this on its own base theme (theme.T1/T2/T3) before
// passing the result to draw.New, so a JSON- or CLI-supplied color
// override works identically regardless of which template is selected. It
// reuses theme.With for the palette roles it shares, then applies
// Background (which is not part of a per-section override).
func (o Options) ApplyTheme(th theme.Theme) theme.Theme {
	th = th.With(theme.SectionOverride{
		Headline:    o.Colors.Headline,
		Subheadline: o.Colors.Subheadline,
		Body:        o.Colors.Body,
		Metadata:    o.Colors.Metadata,
		Link:        o.Colors.Link,
		Border:      o.Colors.Border,
		Profile:     o.Colors.Profile,
	})
	if o.Colors.Background != nil {
		th.Palette.Background = *o.Colors.Background
	}
	return th
}

// Default returns the settings that reproduce the template's original,
// unconfigured behavior: English labels, every section shown if its data
// is present, default order, default spacing.
func Default() Options {
	return Options{
		Lang:         i18n.EN,
		Include:      map[SectionKey]bool{},
		Priority:     map[SectionKey]int{},
		SectionStyle: map[SectionKey]theme.SectionOverride{},
		Density:      1.0,
	}
}

// EffectiveSpacing returns base scaled by Density. base is the active
// template's own theme.Spacing — there's no single global default because
// each template's base rhythm differs (see theme.T1/T2/T3).
func (o Options) EffectiveSpacing(base theme.Spacing) theme.Spacing {
	return base.Scaled(o.densityScale())
}

// EffectiveMargins returns th with its four page margins scaled by Density,
// mirroring EffectiveSpacing so a density choice moves the page's outer
// whitespace and its inner rhythm together. Airier density widens the
// margins, denser narrows them — matching the reference mockups, whose own
// left/right margins vary by density (CV_Template.pdf "dense" ≈37.5pt vs.
// CV_Template_Airy.pdf ≈46pt); before this, --density only tightened gaps
// and left the margins fixed. Font sizes are still never scaled (that stays
// FitToOnePage's job, and even it only touches spacing).
func (o Options) EffectiveMargins(th theme.Theme) theme.Theme {
	k := o.densityScale()
	th.MarginTop *= k
	th.MarginBottom *= k
	th.MarginLeft *= k
	th.MarginRight *= k
	return th
}

// densityScale is Density with the 0 == "unset" convention applied: a zero
// Density means 1.0 (unscaled), any other value is used as-is.
func (o Options) densityScale() float64 {
	if o.Density == 0 {
		return 1.0
	}
	return o.Density
}

// PriorityOf returns key's effective render-order priority: the caller's
// explicit override if one was set, else def. def is supplied by the
// template itself (each template has its own default section order — see
// e.g. t1.go vs t2.go), not by this package.
func (o Options) PriorityOf(key SectionKey, def int) int {
	if p, ok := o.Priority[key]; ok {
		return p
	}
	return def
}

// IncludeOverride reports whether the caller explicitly forced key on or
// off, and if so what. ok is false when there's no override, meaning the
// caller should fall back to "does this section have data".
func (o Options) IncludeOverride(key SectionKey) (include, ok bool) {
	v, ok := o.Include[key]
	return v, ok
}
