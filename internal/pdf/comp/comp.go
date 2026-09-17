// Package comp implements the CV's structural components — header block,
// section heading, experience/education entry, stack line, bullet list,
// skills row, languages line — as functions over internal/pdf/draw. Each
// one is content-agnostic: templates decide ordering and which sections
// to call. Every component reads its colors/sizes from d.Theme, so the
// same components produce t1, t2, or t3 depending only on which theme the
// Doc was created with (see internal/pdf/theme).
package comp

import (
	"strings"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/schema"
)

const photoSize = 56.0

// bulletSep and pipeSep separate items on wrappable joined lines (header
// specialties/contact/links, languages). A single space on each side: any
// wider padding would be discarded anyway once the line wraps, since
// Paragraph/RunsParagraph tokenize on whitespace and rejoin with a single
// measured space — better to be consistent whether or not a given line
// happens to be long enough to wrap.
const (
	bulletSep = " • "
	pipeSep   = " | "
)

// dateStyle returns Italic or Regular depending on the theme's DateItalic
// flag — t2's editorial dates are italic, t1/t3's are upright.
func dateStyle(d *draw.Doc) draw.Style {
	if d.Theme.DateItalic {
		return draw.Italic
	}
	return draw.Regular
}

// headingText applies the theme's case convention to a section label.
func headingText(d *draw.Doc, label string) string {
	if d.Theme.HeadingUppercase {
		return strings.ToUpper(label)
	}
	return label
}

// Header draws the name/tagline/contact block, an optional photo, and the
// heavy accent rule that closes it off. Every line wraps (via Paragraph)
// rather than using the non-wrapping Line, since name/title/specialties/
// contact are all free-text that can run long enough to overflow the page
// width otherwise.
func Header(d *draw.Doc, h schema.Head) error {
	startY := d.Y()
	th := d.Theme

	if err := d.Paragraph(h.FullName, draw.TextRun{
		Size: th.NameSize, Style: draw.Bold, Color: th.Palette.Headline, Tracking: th.NameTracking,
	}, 0); err != nil {
		return err
	}

	if h.JobTitle != "" {
		if err := d.Paragraph(h.JobTitle, draw.TextRun{Size: th.SubtitleSize, Color: th.Palette.Subheadline}, 0); err != nil {
			return err
		}
	}

	if len(h.Specialties) > 0 {
		if err := d.Paragraph(strings.Join(h.Specialties, bulletSep), draw.TextRun{Size: th.ContactSize, Color: th.Palette.Body}, 0); err != nil {
			return err
		}
	}

	contact := []string{h.Address, h.Email, h.Phone}
	if err := d.Paragraph(strings.Join(nonEmpty(contact), pipeSep), draw.TextRun{Size: th.ContactSize, Color: th.Palette.Body}, 0); err != nil {
		return err
	}

	var links []string
	if h.LinkedIn != "" {
		links = append(links, h.LinkedIn)
	}
	if h.GitHub != "" {
		links = append(links, h.GitHub)
	}
	if len(links) > 0 {
		if err := d.Paragraph(strings.Join(links, pipeSep), draw.TextRun{Size: th.ContactSize, Color: th.Palette.Link}, 0); err != nil {
			return err
		}
	}

	if h.Photo != "" {
		photoX := d.MarginLeft() + d.ContentWidth() - photoSize
		if err := d.Image(h.Photo, photoX, startY, photoSize, photoSize); err != nil {
			return err
		}
	}

	d.Rule(th.HeaderRuleWeight, th.Palette.Border)
	return nil
}

// SectionHeading draws a heading (case per the theme) followed by the thin
// section rule. Callers own the vertical space before/after — headings get
// more space above than below per the design's spacing rule of thumb.
func SectionHeading(d *draw.Doc, label string) error {
	th := d.Theme
	if err := d.Line([]draw.TextRun{{
		Text: headingText(d, label), Size: th.HeadingSize, Style: draw.Bold,
		Color: th.Palette.Headline, Tracking: th.HeadingTracking,
	}}); err != nil {
		return err
	}
	d.Rule(th.SectionRuleWeight, th.Palette.Border)
	d.Space(d.Spacing.SM)
	return nil
}

// entryLineGap is the minimum breathing room required between the end of
// "Title — Meta" and the start of the right-aligned date before EntryLine
// decides they'd collide and drops meta to its own line instead.
const entryLineGap = 8.0

// entryMetaSep returns the configured meta separator for this doc, falling
// back to an em-dash if the theme has none set.
func entryMetaSep(d *draw.Doc) string {
	if d.Theme.MetaSep != "" {
		return d.Theme.MetaSep
	}
	return "  —  "
}

// EntryLine draws "Title <sep> Meta" on the left and a right-aligned date
// range, all on one shared baseline. Used for both experience and
// education entries. If title+meta is long enough that it would collide
// with the date, meta drops to a plain line below instead of overlapping
// it — title and date always keep their own baseline.
func EntryLine(d *draw.Doc, title, meta, dateRange string) error {
	th := d.Theme
	sep := entryMetaSep(d)
	titleRun := draw.TextRun{Text: title, Size: th.RoleTitleSize, Style: draw.Bold, Color: th.Palette.Body}
	dateRun := draw.TextRun{Text: dateRange, Size: th.DateSize, Style: dateStyle(d), Color: th.Palette.Body}

	if meta == "" {
		return d.LineWithTail([]draw.TextRun{titleRun}, []draw.TextRun{dateRun})
	}

	metaRun := draw.TextRun{Text: sep + meta, Size: th.CompanySize, Style: draw.Bold, Color: th.Palette.Body}
	titleW, err := d.MeasureRun(titleRun)
	if err != nil {
		return err
	}
	metaW, err := d.MeasureRun(metaRun)
	if err != nil {
		return err
	}
	dateW, err := d.MeasureRun(dateRun)
	if err != nil {
		return err
	}

	if titleW+metaW+entryLineGap+dateW <= d.ContentWidth() {
		return d.LineWithTail([]draw.TextRun{titleRun, metaRun}, []draw.TextRun{dateRun})
	}

	if err := d.LineWithTail([]draw.TextRun{titleRun}, []draw.TextRun{dateRun}); err != nil {
		return err
	}
	return d.Line([]draw.TextRun{{Text: meta, Size: th.CompanySize, Style: draw.Bold, Color: th.Palette.Body}})
}

// StackLine draws the bold, accent "Stack:" label followed by a value that
// wraps onto continuation lines when it is too long to fit on one line.
// Wrapped lines are indented to the start of the value (after the label).
func StackLine(d *draw.Doc, label, value string) error {
	th := d.Theme
	labelRun := draw.TextRun{Text: label, Size: th.StackLabelSize, Style: draw.Bold, Color: th.Palette.Headline}
	valueRun := draw.TextRun{Text: value, Size: th.StackValueSize, Style: dateStyle(d), Color: th.Palette.Body}
	labelW, err := d.MeasureRun(labelRun)
	if err != nil {
		return err
	}
	// RunsParagraph handles wrapping; the label is drawn as a hanging
	// "marker" so wrapped lines start under the value, not the label.
	return d.RunsParagraph([]draw.TextRun{valueRun}, labelW, label, labelRun, 0)
}

// BulletList draws each item as a single-style bulleted, hanging-indent
// paragraph, with space-xs between items.
func BulletList(d *draw.Doc, items []string) error {
	th := d.Theme
	run := draw.TextRun{Size: th.BodySize, Color: th.Palette.Body}
	for i, item := range items {
		if i > 0 {
			d.Space(d.Spacing.XS)
		}
		if err := d.Bullet(th.BulletMarker, item, run, th.BulletHangingIndent, th.BulletOffset); err != nil {
			return err
		}
	}
	return nil
}

// BulletRunsList is BulletList for items that need one bold lead-in
// segment followed by regular text (e.g. "**Project** (Tech): description"),
// flowing as one wrapped paragraph per bullet. No em dash — plain space,
// same understated punctuation as the Stack line.
func BulletRunsList(d *draw.Doc, items [][2]string) error {
	th := d.Theme
	markerRun := draw.TextRun{Text: th.BulletMarker, Size: th.BodySize, Color: th.Palette.Body}
	for i, item := range items {
		if i > 0 {
			d.Space(d.Spacing.XS)
		}
		lead, rest := item[0], item[1]
		runs := []draw.TextRun{
			{Text: lead, Size: th.BodySize, Style: draw.Bold, Color: th.Palette.Body},
		}
		if rest != "" {
			runs = append(runs, draw.TextRun{Text: " " + rest, Size: th.BodySize, Color: th.Palette.Body})
		}
		if err := d.RunsParagraph(runs, th.BulletHangingIndent, th.BulletMarker, markerRun, th.BulletOffset); err != nil {
			return err
		}
	}
	return nil
}

// SkillsList draws each row as "**Category:** item, item, item" — a bold
// key inline with its value. Not a fixed-width table: the mockups don't
// align these into columns, so each row's value starts right after its
// own key plus a small gap, and wrapped continuation lines hang under the
// value's start rather than back at the margin.
func SkillsList(d *draw.Doc, rows [][2]string) error {
	th := d.Theme
	for _, r := range rows {
		key := draw.TextRun{Text: r[0], Size: th.SkillsKeySize, Style: draw.Bold, Color: th.Palette.Body}
		val := draw.TextRun{Text: r[1], Size: th.SkillsValueSize, Color: th.Palette.Body}
		keyW, err := d.MeasureRun(key)
		if err != nil {
			return err
		}
		indent := keyW + th.SkillsGap
		if err := d.RunsParagraph([]draw.TextRun{val}, indent, key.Text, key, th.SkillsGap); err != nil {
			return err
		}
	}
	return nil
}

// Languages draws every language on one inline, bullet-separated line —
// the mockups' actual layout, not a stacked table.
func Languages(d *draw.Doc, parts []string) error {
	return d.Paragraph(strings.Join(parts, bulletSep), draw.TextRun{Size: d.Theme.BodySize, Color: d.Theme.Palette.Body}, 0)
}

// Profile draws the summary paragraph in the theme's dedicated Profile tone
// — typically paler than the metadata grey used elsewhere — and in italic
// or upright per the theme's ProfileItalic flag (a summary reads as italic
// by a different convention than dates, so it has its own flag rather than
// borrowing DateItalic or hardcoding the style here).
func Profile(d *draw.Doc, summary string) error {
	style := draw.Regular
	if d.Theme.ProfileItalic {
		style = draw.Italic
	}
	return d.Paragraph(summary, draw.TextRun{Size: d.Theme.BodySize, Style: style, Color: d.Theme.Palette.Body}, 0)
}

func nonEmpty(ss []string) []string {
	out := make([]string, 0, len(ss))
	for _, s := range ss {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
