// Package t1 renders a CV using the "ATS Classic" layout: a conservative,
// single-hue, non-italic design built for two things t2 doesn't
// prioritize — passing keyword-matching ATS software cleanly, and
// following the Moroccan/French hiring convention (reverse-chronological,
// experience-led, sober navy, no editorial tracking/uppercase flourishes
// beyond plain section headings). It shares every drawing primitive with
// t2 (internal/pdf/draw, internal/pdf/comp, internal/pdf/sections) — only
// the theme (internal/pdf/theme.T1) and the section order differ. See
// docs/pdf-rewrite-handoff.md for the layered architecture this and every
// other template is built on.
package t1

import (
	"github.com/othmaneBakkass/cv_gen/internal/pdf/comp"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/render"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/sections"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

// Render builds the t1 layout for cv under opts and returns the finished
// PDF bytes. The whole layout is built through render.FitToOnePage, which
// re-runs it against progressively tighter spacing until it fits on one
// page (or gives up at a floor and lets it spill onto a second).
func Render(cv schema.CV, opts render.Options) ([]byte, error) {
	th := opts.EffectiveMargins(opts.ApplyTheme(theme.T1))
	labels := i18n.For(opts.Lang)

	build := func(d *draw.Doc) error {
		return renderBody(d, cv, opts, labels)
	}

	d, err := render.FitToOnePage(th, opts.EffectiveSpacing(th.Spacing), build)
	if err != nil {
		return nil, err
	}
	return d.Bytes()
}

func renderBody(d *draw.Doc, cv schema.CV, opts render.Options, labels i18n.Labels) error {
	if err := comp.Header(d, cv.Head); err != nil {
		return err
	}

	// Default order puts Experience before Education (priorities 20/30,
	// the reverse of t2's 20/30 split) — the convention recruiters and
	// ATS parsers expect once a candidate has real work history, rather
	// than t2's academic-CV-style Profile-then-Education-first order.
	list := []render.Section{
		{Key: render.Profile, Priority: opts.PriorityOf(render.Profile, 10), HasData: cv.Profile != "", Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Profile); err != nil {
				return err
			}
			return comp.Profile(d, cv.Profile)
		}},
		{Key: render.Experience, Priority: opts.PriorityOf(render.Experience, 20), HasData: len(cv.Jobs) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Experience); err != nil {
				return err
			}
			return sections.Jobs(d, cv.Jobs, labels)
		}},
		{Key: render.Education, Priority: opts.PriorityOf(render.Education, 30), HasData: len(cv.Education) > 0, Body: func(d *draw.Doc) error {
			if err := comp.SectionHeading(d, labels.Education); err != nil {
				return err
			}
			return sections.Education(d, cv.Education, labels)
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
