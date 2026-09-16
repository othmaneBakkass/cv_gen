// Package draw is a thin layout layer over gopdf: cursor/page state, a
// font registry with real tracking and measurement, a greedy line wrapper,
// and the handful of primitives (rules, right-tab-stop rows, hanging-indent
// paragraphs) that the CV templates are built from. It knows nothing about
// CV content — that's internal/pdf/comp.
package draw

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/signintech/gopdf"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw/fonts"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// Style selects a registered face within the single embedded family.
type Style int

const (
	Regular Style = iota
	Bold
	Italic
	BoldItalic
)

func (s Style) gopdfStyle() string {
	switch s {
	case Bold:
		return "B"
	case Italic:
		return "I"
	case BoldItalic:
		return "BI"
	default:
		return ""
	}
}

// TextRun is one contiguous span of styled text — the unit every drawing
// primitive in this package operates on. Mixed-style lines ("**Role** —
// *Company*") are just a []TextRun.
type TextRun struct {
	Text     string
	Size     float64
	Style    Style
	Color    theme.RGB
	Tracking float64 // extra character spacing in points; 0 = none
}

// ascentRatio and leadingRatio approximate Carlito's metrics well enough
// for this layout: baseline offset from a line's top, and line-to-line
// advance, both as a fraction of font size. Tuned against rendered output
// rather than derived from the font's actual hhea table.
const (
	ascentRatio  = 0.75
	leadingRatio = 1.28
)

func ascent(size float64) float64  { return size * ascentRatio }
func leading(size float64) float64 { return size * leadingRatio }

// Doc wraps a gopdf document with an explicit vertical cursor. Every
// drawing call advances Y; callers never touch gopdf directly.
type Doc struct {
	PDF *gopdf.GoPdf

	// Theme is the active template's design tokens. comp-layer functions
	// read colors/sizes/margins from here (d.Theme.Accent, etc.) — it's
	// the only thing that differs between templates built on this layer.
	Theme theme.Theme

	// Spacing is the vertical rhythm scale in effect for this document,
	// seeded from Theme.Spacing but kept as its own field so a density
	// setting or a fit-to-page pass can scale it between renders without
	// mutating the theme itself.
	Spacing theme.Spacing

	pageW, pageH float64
	y            float64
}

// Options configures a new Doc.
type Options struct {
	// NoCompress disables content-stream compression, which makes the
	// coordinate-parsing verification trick in docs/pdf-rewrite-handoff.md
	// work directly against the saved file. Leave false for normal use.
	NoCompress bool
	// Spacing overrides the theme's base vertical rhythm scale (e.g. for a
	// denser or airier render). Zero value means th.Spacing.
	Spacing theme.Spacing
	// PageHeight overrides the page height used for the bottom-margin/
	// page-break check in EnsureSpace. Zero means real A4 (A4Height). Set
	// to something very large for a measurement-only pass (see
	// render.FitToOnePage) so content never triggers a page break, letting
	// the final cursor Y report the content's true total height.
	PageHeight float64
}

// A4Height is the point height of an A4 page — the same value draw.New
// uses by default. Exported so render.FitToOnePage can compute a real
// page's available content height without importing gopdf itself.
var A4Height = gopdf.PageSizeA4.H

// New creates an A4 document using th's margins and font family, with the
// embedded Carlito faces registered under that family name.
func New(th theme.Theme, opts Options) (*Doc, error) {
	pdf := &gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, Unit: gopdf.Unit_PT})
	if opts.NoCompress {
		pdf.SetNoCompression()
	}

	faces := []struct {
		data  []byte
		style int
	}{
		{fonts.CarlitoRegular, gopdf.Regular},
		{fonts.CarlitoBold, gopdf.Bold},
		{fonts.CarlitoItalic, gopdf.Italic},
		{fonts.CarlitoBoldItalic, gopdf.Bold | gopdf.Italic},
	}
	for _, f := range faces {
		opt := gopdf.TtfOption{Style: f.style}
		if err := pdf.AddTTFFontByReaderWithOption(th.FontFamily, bytes.NewReader(f.data), opt); err != nil {
			return nil, fmt.Errorf("register font face: %w", err)
		}
	}

	pdf.AddPage()

	spacing := opts.Spacing
	if spacing == (theme.Spacing{}) {
		spacing = th.Spacing
	}

	pageH := opts.PageHeight
	if pageH == 0 {
		pageH = A4Height
	}

	d := &Doc{
		PDF:     pdf,
		Theme:   th,
		Spacing: spacing,
		pageW:   gopdf.PageSizeA4.W,
		pageH:   pageH,
	}
	d.y = th.MarginTop
	return d, nil
}

// ContentWidth is the usable width between the left and right margins.
func (d *Doc) ContentWidth() float64 { return d.pageW - d.Theme.MarginLeft - d.Theme.MarginRight }

// MarginLeft is the left edge of the content area.
func (d *Doc) MarginLeft() float64 { return d.Theme.MarginLeft }

// Y is the current vertical cursor: the top of the next block to draw.
func (d *Doc) Y() float64 { return d.y }

// SetY overrides the vertical cursor directly. Used by fit-to-page passes.
func (d *Doc) SetY(y float64) { d.y = y }

// Bytes returns the finished PDF.
func (d *Doc) Bytes() ([]byte, error) { return d.PDF.GetBytesPdfReturnErr() }

// setRunFont selects r's face/size/color/tracking as the active gopdf font.
func (d *Doc) setRunFont(r TextRun) error {
	if err := d.PDF.SetFont(d.Theme.FontFamily, r.Style.gopdfStyle(), r.Size); err != nil {
		return fmt.Errorf("set font: %w", err)
	}
	d.PDF.SetTextColor(r.Color.R, r.Color.G, r.Color.B)
	if err := d.PDF.SetCharSpacing(r.Tracking); err != nil {
		return fmt.Errorf("set tracking: %w", err)
	}
	return nil
}

// MeasureRun returns r's rendered width under its own font/size/tracking.
func (d *Doc) MeasureRun(r TextRun) (float64, error) {
	if err := d.setRunFont(r); err != nil {
		return 0, err
	}
	return d.PDF.MeasureTextWidth(r.Text)
}

// EnsureSpace starts a new page if h more points would overflow the
// bottom margin, resetting the cursor to the top margin.
func (d *Doc) EnsureSpace(h float64) {
	if d.y+h > d.pageH-d.Theme.MarginBottom {
		d.PDF.AddPage()
		d.y = d.Theme.MarginTop
	}
}

// Space advances the cursor by pt without drawing anything.
func (d *Doc) Space(pt float64) { d.y += pt }

// Rule draws a horizontal line spanning the full content width at the
// current cursor position, then advances past it. A theme can opt out of
// the rule entirely (e.g. a heading style that relies on weight/color
// alone) by passing thickness <= 0, in which case this is a no-op.
func (d *Doc) Rule(thickness float64, color theme.RGB) {
	if thickness <= 0 {
		return
	}
	d.EnsureSpace(thickness)
	d.PDF.SetLineWidth(thickness)
	d.PDF.SetStrokeColor(color.R, color.G, color.B)
	y := d.y + thickness/2
	d.PDF.Line(d.Theme.MarginLeft, y, d.Theme.MarginLeft+d.ContentWidth(), y)
	d.y += thickness
}

// drawAt places a single run with its baseline at (x, baseline).
func (d *Doc) drawAt(r TextRun, x, baseline float64) (float64, error) {
	if err := d.setRunFont(r); err != nil {
		return 0, err
	}
	d.PDF.SetXY(x, baseline)
	if err := d.PDF.Text(r.Text); err != nil {
		return 0, fmt.Errorf("draw text: %w", err)
	}
	return d.PDF.MeasureTextWidth(r.Text)
}

// maxSize returns the largest font size among runs, falling back to 0 for
// an empty slice.
func maxSize(runs []TextRun) float64 {
	var m float64
	for _, r := range runs {
		if r.Size > m {
			m = r.Size
		}
	}
	return m
}

// Line draws runs left-to-right starting at the left margin, all sharing
// one baseline (computed from the largest run, so mixed sizes never drift
// relative to each other), then advances the cursor past the line.
func (d *Doc) Line(runs []TextRun) error {
	return d.LineAt(runs, d.Theme.MarginLeft)
}

// LineAt is Line with an explicit start X (for indented lines).
func (d *Doc) LineAt(runs []TextRun, x float64) error {
	size := maxSize(runs)
	d.EnsureSpace(leading(size))
	baseline := d.y + ascent(size)
	for _, r := range runs {
		w, err := d.drawAt(r, x, baseline)
		if err != nil {
			return err
		}
		x += w
	}
	d.y += leading(size)
	return nil
}

// LineWithTail draws left starting at the left margin and right
// right-aligned to the content's right edge, both on one shared baseline —
// this is the entry-line component (title/company … date) and the only
// place this template needs a right tab stop.
func (d *Doc) LineWithTail(left, right []TextRun) error {
	size := maxSize(append(append([]TextRun{}, left...), right...))
	d.EnsureSpace(leading(size))
	baseline := d.y + ascent(size)

	x := d.Theme.MarginLeft
	for _, r := range left {
		w, err := d.drawAt(r, x, baseline)
		if err != nil {
			return err
		}
		x += w
	}

	var totalRight float64
	for _, r := range right {
		w, err := d.MeasureRun(r)
		if err != nil {
			return err
		}
		totalRight += w
	}
	rx := d.Theme.MarginLeft + d.ContentWidth() - totalRight
	for _, r := range right {
		w, err := d.drawAt(r, rx, baseline)
		if err != nil {
			return err
		}
		rx += w
	}

	d.y += leading(size)
	return nil
}

// WrapText greedily wraps text into lines no wider than width, under run's
// font. Words longer than width are placed on their own line rather than
// broken.
func (d *Doc) WrapText(text string, width float64, run TextRun) ([]string, error) {
	if err := d.setRunFont(run); err != nil {
		return nil, err
	}
	return d.PDF.SplitTextWithWordWrap(text, width)
}

// Paragraph wraps and draws text as a single-style block, every line
// (including the first) starting at MarginLeft()+indent. That's a true
// hanging indent for callers that draw a bullet marker separately before
// calling this.
func (d *Doc) Paragraph(text string, run TextRun, indent float64) error {
	run.Text = text
	return d.RunsParagraph([]TextRun{run}, indent, "", TextRun{}, 0)
}

// Bullet draws text as a bulleted, hanging-indent paragraph: marker sits
// at indent-markerOffset, wrapped lines all start at indent.
func (d *Doc) Bullet(marker, text string, run TextRun, indent, markerOffset float64) error {
	run.Text = text
	markerRun := run
	markerRun.Text = marker
	return d.RunsParagraph([]TextRun{run}, indent, marker, markerRun, markerOffset)
}

// RunsParagraph greedily word-wraps a sequence of mixed-style runs as one
// flowing paragraph — the inline mixed-style text ("**Role** — *Company*")
// that a row/cell-based renderer can't measure. Every line starts at
// MarginLeft()+indent. If marker is non-empty, it's drawn at
// indent-markerOffset on the first line's baseline only (a bullet that
// doesn't participate in wrapping).
func (d *Doc) RunsParagraph(runs []TextRun, indent float64, marker string, markerRun TextRun, markerOffset float64) error {
	type word struct {
		text  string
		run   TextRun
		width float64
	}

	var words []word
	for _, r := range runs {
		for _, part := range strings.Fields(r.Text) {
			wr := r
			wr.Text = part
			w, err := d.MeasureRun(wr)
			if err != nil {
				return err
			}
			words = append(words, word{part, wr, w})
		}
	}
	if len(words) == 0 {
		return nil
	}

	width := d.ContentWidth() - indent
	x0 := d.Theme.MarginLeft + indent
	first := true

	spaceWidth := func(r TextRun) (float64, error) {
		sr := r
		sr.Text = " "
		return d.MeasureRun(sr)
	}

	var line []word
	var lineW float64

	drawLine := func() error {
		if len(line) == 0 {
			return nil
		}
		size := 0.0
		for _, w := range line {
			if w.run.Size > size {
				size = w.run.Size
			}
		}
		if marker != "" && markerRun.Size > size {
			size = markerRun.Size
		}
		d.EnsureSpace(leading(size))
		baseline := d.y + ascent(size)
		if first && marker != "" {
			markerW, err := d.MeasureRun(markerRun)
			if err != nil {
				return err
			}
			if _, err := d.drawAt(markerRun, x0-markerOffset-markerW, baseline); err != nil {
				return err
			}
		}
		x := x0
		for i, w := range line {
			if i > 0 {
				sw, err := spaceWidth(w.run)
				if err != nil {
					return err
				}
				x += sw
			}
			if _, err := d.drawAt(w.run, x, baseline); err != nil {
				return err
			}
			x += w.width
		}
		d.y += leading(size)
		first = false
		return nil
	}

	for _, w := range words {
		add := w.width
		if len(line) > 0 {
			sw, err := spaceWidth(line[len(line)-1].run)
			if err != nil {
				return err
			}
			add += sw
		}
		if len(line) > 0 && lineW+add > width {
			if err := drawLine(); err != nil {
				return err
			}
			line = nil
			lineW = 0
			add = w.width
		}
		line = append(line, w)
		lineW += add
	}
	return drawLine()
}

// Image places a raster image with its top-left corner at the current
// cursor and given width/height, without advancing the cursor (callers
// with a photo alongside text manage the row height themselves).
func (d *Doc) Image(path string, x, y, w, h float64) error {
	return d.PDF.Image(path, x, y, &gopdf.Rect{W: w, H: h})
}
