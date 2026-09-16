package draw

import (
	"math"
	"testing"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// TestSpacingRhythmMatchesMockups pins the vertical rhythm the XS/SM/MD/LG
// spacing tokens produce, measured as baseline-to-baseline deltas, against
// what the reference mockups (Templates.claude/*.pdf) actually draw. It is a
// white-box test in package draw so it can use the same leading()/ascent()
// the renderer uses, deriving each gap from theme.T2's tokens exactly the
// way RunSections + comp do at render time (the derivations were confirmed
// to match a real t2 render to 0.01pt).
//
// Findings (see docs/pdf-rewrite-handoff.md §0.3):
//   - MD (gap before a section heading) and XS (bullet-to-bullet) land within
//     ~1pt of the mockups — the tokens are validated as correct.
//   - SM (heading -> first body) is intentionally tighter than the mockups,
//     per pdfDesign.md's rule that a heading gets more space above than below.
//   - LG (header -> first section) is intentionally the largest gap, again
//     per pdfDesign.md; the Word mockups happen to keep it equal to MD.
//
// The mockup targets below are the measured baseline deltas; the tolerances
// encode "matches the mockup" for MD/XS and "intentionally diverges by a
// known, bounded amount" for SM/LG. If someone retunes a token, this test
// makes the effect on the real rhythm explicit instead of silent.
func TestSpacingRhythmMatchesMockups(t *testing.T) {
	th := theme.T2
	body, heading := th.BodySize, th.HeadingSize

	// Gap before a section heading: previous line's descent-to-next-top
	// (leading-ascent) + MD + the heading line's ascent.
	beforeHeading := (leading(body) - ascent(body)) + th.Spacing.MD + ascent(heading)
	assertNear(t, "before-heading (MD)", beforeHeading, 21.2 /*mockup*/, 2.0)

	// Bullet-to-bullet: one body line's leading plus the XS gap between items.
	bulletToBullet := leading(body) + th.Spacing.XS
	assertNear(t, "bullet-to-bullet (XS)", bulletToBullet, 13.8 /*mockup*/, 1.5)

	// Heading -> first body: heading line's descent-to-next-top + the rule +
	// SM + the body line's ascent. Intentionally tighter than the mockup.
	headingToBody := (leading(heading) - ascent(heading)) + th.SectionRuleWeight + th.Spacing.SM + ascent(body)
	assertNear(t, "heading-to-body (SM)", headingToBody, 17.3 /*ours; mockup ~21.4*/, 1.5)

	// Header -> first section: header last line's descent-to-next-top + the
	// header rule + LG + heading ascent. Intentionally the largest gap.
	headerToSection := (leading(th.ContactSize) - ascent(th.ContactSize)) + th.HeaderRuleWeight + th.Spacing.LG + ascent(heading)
	assertNear(t, "header-to-section (LG)", headerToSection, 26.3 /*ours; mockup ~21.1*/, 1.5)
}

func assertNear(t *testing.T, name string, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("%s = %.2fpt, want %.2fpt ±%.1f", name, got, want, tol)
	}
}
