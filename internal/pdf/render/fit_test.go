package render

import (
	"testing"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/draw"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// buildNLines returns a build func that draws n short lines, each preceded
// by a Spacing.MD gap — a synthetic stand-in for "n CV entries" with a
// controllable, predictable total height: n * (MD_at_current_scale + one
// line's leading). Real font leading doesn't scale with density (only the
// MD gap does), matching how actual templates behave.
func buildNLines(n int) func(d *draw.Doc) error {
	return func(d *draw.Doc) error {
		for i := 0; i < n; i++ {
			d.Space(d.Spacing.MD)
			if err := d.Line([]draw.TextRun{{Text: "x", Size: 10, Color: theme.T2.Palette.Body}}); err != nil {
				return err
			}
		}
		return nil
	}
}

func availableHeight(th theme.Theme) float64 {
	return draw.A4Height - th.MarginTop - th.MarginBottom
}

func TestFitToOnePage_FitsWithoutShrinking(t *testing.T) {
	th := theme.T2
	d, err := FitToOnePage(th, th.Spacing, TypographySettings{}, buildNLines(10))
	if err != nil {
		t.Fatalf("FitToOnePage: %v", err)
	}
	if d.Spacing.MD != th.Spacing.MD {
		t.Errorf("expected no shrink for content that already fits, got MD %v want %v", d.Spacing.MD, th.Spacing.MD)
	}
	if used := d.Y() - th.MarginTop; used > availableHeight(th) {
		t.Errorf("content unexpectedly overflows one page: used %v, available %v", used, availableHeight(th))
	}
}

func TestFitToOnePage_ShrinksToFit(t *testing.T) {
	th := theme.T2
	// Tuned so it overflows one page at scale 1.0 but fits by the time
	// scale drops to fitFloor (see the derivation in the handoff/PR: at
	// n=42, 1.0 -> 831.6pt (overflow vs. 794pt available), 0.7 -> 743.4pt
	// (fits)) — this must actually need shrinking, not just barely fit.
	n := 42
	build := buildNLines(n)

	// Sanity-check the premise: unscaled would overflow.
	unscaled, err := draw.New(th, draw.Options{PageHeight: 20000})
	if err != nil {
		t.Fatalf("draw.New: %v", err)
	}
	if err := build(unscaled); err != nil {
		t.Fatalf("build: %v", err)
	}
	if used := unscaled.Y() - th.MarginTop; used <= availableHeight(th) {
		t.Fatalf("test premise broken: n=%d already fits unscaled (used %v, available %v) — increase n", n, used, availableHeight(th))
	}

	d, err := FitToOnePage(th, th.Spacing, TypographySettings{}, build)
	if err != nil {
		t.Fatalf("FitToOnePage: %v", err)
	}
	if d.Spacing.MD >= th.Spacing.MD {
		t.Errorf("expected spacing to shrink below base, got MD %v want < %v", d.Spacing.MD, th.Spacing.MD)
	}
	if d.Spacing.MD < th.Spacing.MD*fitFloor {
		t.Errorf("shrunk past the floor: got MD %v, floor is %v", d.Spacing.MD, th.Spacing.MD*fitFloor)
	}
	if used := d.Y() - th.MarginTop; used > availableHeight(th) {
		t.Errorf("still overflows after shrinking: used %v, available %v", used, availableHeight(th))
	}
}

func TestFitToOnePage_GivesUpGracefullyAtFloor(t *testing.T) {
	th := theme.T2
	// Way more content than fitFloor could ever reclaim.
	d, err := FitToOnePage(th, th.Spacing, TypographySettings{}, buildNLines(100))
	if err != nil {
		t.Fatalf("FitToOnePage should not error even when content can't fit: %v", err)
	}
	wantFloorMD := th.Spacing.MD * fitFloor
	if d.Spacing.MD != wantFloorMD {
		t.Errorf("expected spacing pinned at the floor, got MD %v want %v", d.Spacing.MD, wantFloorMD)
	}
	// The real render legitimately paginates once content doesn't fit —
	// that's the whole point of "give up gracefully" — so the assertion
	// has to be page count, not leftover Y on whatever the last page is.
	if pages := d.PDF.GetNumberOfPages(); pages <= 1 {
		t.Errorf("expected this content to spill onto more than one real page even at the floor, got %d page(s) — otherwise the test isn't exercising the give-up path", pages)
	}
}
