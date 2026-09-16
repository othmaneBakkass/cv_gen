// Package sections renders each CV section's content (education entries,
// jobs, skills, languages, projects) against internal/pdf/comp. It's the
// layer between raw schema data and the comp primitives — shared by every
// template, since the content formatting (e.g. "Degree — School, City" /
// right-aligned date range) doesn't change between t1, t2, and t3; only
// the theme and section order do.
package sections

import (
	"fmt"
	"strings"

	"github.com/othmaneBakkass/cv_gen/internal/common/stringc"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/comp"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/i18n"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

// Education draws each entry as a title/meta/date EntryLine followed by an
// optional description paragraph.
func Education(d *draw.Doc, entries []schema.Education, labels i18n.Labels) error {
	for i, e := range entries {
		if i > 0 {
			d.Space(d.Spacing.MD)
		}
		title := stringc.ToCapital(e.Degree)
		meta := fmt.Sprintf("%s, %s", e.School, e.Location)
		dateRange := e.StartedAt + labels.DateSep + e.EndedAt
		if err := comp.EntryLine(d, title, meta, dateRange); err != nil {
			return err
		}
		if e.Description != "" {
			d.Space(d.Spacing.XS)
			if err := d.Paragraph(e.Description, draw.TextRun{Size: d.Theme.BodySize, Color: d.Theme.Palette.Body}, 0); err != nil {
				return err
			}
		}
	}
	return nil
}

// Jobs draws each entry as a title/meta/date EntryLine, an optional Stack
// line, then a bullet list of highlights.
func Jobs(d *draw.Doc, entries []schema.Job, labels i18n.Labels) error {
	for i, j := range entries {
		if i > 0 {
			d.Space(d.Spacing.MD)
		}
		title := stringc.ToCapital(j.Position)
		meta := fmt.Sprintf("%s, %s", stringc.ToCapital(j.Company), j.Location)
		dateRange := j.StartedAt + labels.DateSep + j.EndedAt
		if err := comp.EntryLine(d, title, meta, dateRange); err != nil {
			return err
		}
		if len(j.Tools) > 0 {
			d.Space(d.Spacing.XS)
			if err := comp.StackLine(d, labels.Stack, strings.Join(j.Tools, ", ")); err != nil {
				return err
			}
		}
		if len(j.Highlights) > 0 {
			d.Space(d.Spacing.SM)
			if err := comp.BulletList(d, j.Highlights); err != nil {
				return err
			}
		}
	}
	return nil
}

// Skills draws each category as a "**Category:** item, item" row.
func Skills(d *draw.Doc, entries []schema.Skill) error {
	rows := make([][2]string, len(entries))
	for i, s := range entries {
		rows[i] = [2]string{s.Category + ":", strings.Join(s.Items, ", ")}
	}
	return comp.SkillsList(d, rows)
}

// Languages draws every language on one inline, bullet-separated line.
func Languages(d *draw.Doc, entries []schema.Language) error {
	parts := make([]string, len(entries))
	for i, l := range entries {
		parts[i] = fmt.Sprintf("%s: %s", stringc.ToCapital(l.Language), stringc.ToCapital(l.Level))
	}
	return comp.Languages(d, parts)
}

// Projects draws each project as a "**Name** (Tech): Description" bullet —
// no em dash, tech folded into a parenthetical right after the name.
func Projects(d *draw.Doc, entries []schema.Project) error {
	items := make([][2]string, len(entries))
	for i, p := range entries {
		items[i] = [2]string{p.Name, fmt.Sprintf("(%s): %s", p.Tech, p.Description)}
	}
	return comp.BulletRunsList(d, items)
}
