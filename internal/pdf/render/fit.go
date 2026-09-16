package render

import (
	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

const (
	// fitFloor is the smallest spacing multiplier FitToOnePage will try.
	// Below this, compressing further would read as cramped rather than
	// "tightened to fit" — better to accept a second page than degrade
	// past this point.
	fitFloor = 0.7
	fitStep  = 0.05
	// measurePageHeight is large enough that no real CV's content would
	// ever trigger EnsureSpace's page-break check during measurement, so
	// the final cursor Y reports the content's true total height.
	measurePageHeight = 20000.0
)

// FitToOnePage renders build against progressively tighter spacing (from
// 1.0 down to fitFloor, in fitStep increments) until the result fits on
// one real A4 page, then performs the real render at whatever scale
// worked and returns that Doc. If nothing even at fitFloor fits, it gives
// up and returns the fitFloor render as-is — which may still spill onto a
// second page — rather than compressing further.
//
// Only theme.Spacing (the gaps between things) is ever scaled — font
// sizes are never touched, per pdfDesign.md's own guidance to compress
// line spacing aggressively before shrinking fonts. build must be pure
// beyond drawing into the *draw.Doc it's given: it runs once per candidate
// scale during measurement, then once more for the real render.
func FitToOnePage(th theme.Theme, baseSpacing theme.Spacing, build func(d *draw.Doc) error) (*draw.Doc, error) {
	available := draw.A4Height - th.MarginTop - th.MarginBottom

	scale := 1.0
	for {
		measureTh := th
		measureTh.Spacing = baseSpacing.Scaled(scale)
		md, err := draw.New(measureTh, draw.Options{PageHeight: measurePageHeight})
		if err != nil {
			return nil, err
		}
		if err := build(md); err != nil {
			return nil, err
		}

		used := md.Y() - th.MarginTop
		if used <= available || scale <= fitFloor {
			break
		}
		scale -= fitStep
		if scale < fitFloor {
			scale = fitFloor
		}
	}

	finalTh := th
	finalTh.Spacing = baseSpacing.Scaled(scale)
	d, err := draw.New(finalTh, draw.Options{})
	if err != nil {
		return nil, err
	}
	if err := build(d); err != nil {
		return nil, err
	}
	return d, nil
}
