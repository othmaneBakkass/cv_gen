// Package t2 renders a CV using the "Navy Rule" layout: a single-column,
// teal-accented, editorial design with tracked uppercase headings, thin
// rules under every section heading, right-aligned italic dates, and an
// optional photo. See the mockups in Templates.claude/ for the reference
// design (pdfDesign.md is a directional guide only — it disagrees with the
// mockups on several values, see docs/pdf-rewrite-handoff.md §4) and
// docs/pdf-rewrite-handoff.md for why this is built on gopdf via the
// draw/comp layers rather than on maroto.
package t2

import (
	"github.com/othmaneBakkass/cv_gen/internal/pdf/comp"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/sections"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

// Render builds the t2 layout for cv under opts and returns the finished
// PDF bytes. The whole layout is built through render.FitToOnePage, which
// re-runs it against progressively tighter spacing (and optionally font sizes)
// until it fits on one page.
func Render(cv schema.CV, opts render.Options) ([]byte, error) {
	th := opts.EffectiveMargins(opts.ApplyTheme(theme.T2))
	labels := i18n.For(opts.Lang)

	build := func(d *draw.Doc) error {
		return renderBody(d, cv, opts, labels)
	}

	d, err := render.FitToOnePage(th, opts.EffectiveSpacing(th.Spacing), opts.Typography, build)
	if err != nil {
		return nil, err
	}
	return d.Bytes()
}

func renderBody(d *draw.Doc, cv schema.CV, opts render.Options, labels i18n.Labels) error {
	if err := comp.Header(d, cv.Head); err != nil {
		return err
	}

	list := []render.Section{
		{Key: render.Profile, Priority: opts.PriorityOf(render.Profile, 10), HasData: cv.Profile != "", Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Profile); err != nil {
				return err
			}
			return comp.Profile(d, cv.Profile)
		}},
		{Key: render.Education, Priority: opts.PriorityOf(render.Education, 20), HasData: len(cv.Education) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Education); err != nil {
				return err
			}
			return sections.Education(d, cv.Education, labels)
		}},
		{Key: render.Experience, Priority: opts.PriorityOf(render.Experience, 30), HasData: len(cv.Jobs) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Experience); err != nil {
				return err
			}
			return sections.Jobs(d, cv.Jobs, labels)
		}},
		{Key: render.Skills, Priority: opts.PriorityOf(render.Skills, 40), HasData: len(cv.Skills) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Skills); err != nil {
				return err
			}
			return sections.Skills(d, cv.Skills)
		}},
		{Key: render.Projects, Priority: opts.PriorityOf(render.Projects, 50), HasData: len(cv.Projects) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Projects); err != nil {
				return err
			}
			return sections.Projects(d, cv.Projects)
		}},
		{Key: render.Certifications, Priority: opts.PriorityOf(render.Certifications, 60), HasData: len(cv.Certifications) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Certifications); err != nil {
				return err
			}
			return comp.BulletList(d, cv.Certifications)
		}},
		{Key: render.Languages, Priority: opts.PriorityOf(render.Languages, 70), HasData: len(cv.Languages) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Languages); err != nil {
				return err
			}
			return sections.Languages(d, cv.Languages)
		}},
	}

	return render.RunSections(d, opts, list)
}
