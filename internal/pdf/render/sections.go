package render

import (
	"fmt"
	"sort"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// Section is one entry in a template's ordered, filterable render list —
// the small "DSL" described in docs/pdf-rewrite-handoff.md §6.3, shared
// across templates since every template needs the same include/order
// logic and only the section bodies differ.
//
// Style carries the template's own per-section styling defaults (e.g. a
// template that wants Projects to drop its divider). It is the first of the
// two override layers applied on top of the global theme; the second is the
// caller's opts.SectionStyle (JSON/CLI), which wins over it. Together they
// implement Global Template Defaults -> Section Overrides -> Rendered
// Section (see theme.SectionOverride).
type Section struct {
	Key      SectionKey
	Priority int
	HasData  bool
	Style    theme.SectionOverride
	Body     func(d *draw.Doc) error
}

// RunSections sorts sections by priority, applies opts' include overrides,
// and draws each included one — the first gets d.Spacing.LG below the
// header, the rest get d.Spacing.MD between them, unless a section's
// SpaceBefore override says otherwise.
//
// Each section renders against its own effective theme: the document's base
// theme with the template's Section.Style merged on, then the caller's
// opts.SectionStyle[key] merged on top (so a user override beats a template
// default). d.Theme is swapped in for the section body and restored after,
// so the comp layer keeps reading everything from d.Theme with no knowledge
// that per-section styling exists.
func RunSections(d *draw.Doc, opts Options, sections []Section) error {
	sort.SliceStable(sections, func(i, j int) bool { return sections[i].Priority < sections[j].Priority })

	base := d.Theme
	first := true
	for _, s := range sections {
		include := s.HasData
		if v, ok := opts.IncludeOverride(s.Key); ok {
			include = v
		}
		if !include {
			continue
		}

		userStyle := opts.SectionStyle[s.Key]

		// Vertical gap before the section: an explicit SpaceBefore override
		// (user first, then template) wins; otherwise LG below the header for
		// the first section, MD between the rest.
		space := d.Spacing.LG
		if !first {
			space = d.Spacing.MD
		}
		if s.Style.SpaceBefore != nil {
			space = *s.Style.SpaceBefore
		}
		if userStyle.SpaceBefore != nil {
			space = *userStyle.SpaceBefore
		}
		d.Space(space)
		first = false

		// Resolve this section's theme: base -> template default -> user.
		d.Theme = base.With(s.Style).With(userStyle)
		err := s.Body(d)
		d.Theme = base
		if err != nil {
			return fmt.Errorf("render %s: %w", s.Key, err)
		}
	}
	return nil
}
