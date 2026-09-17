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
	fitFloor = 0.5
	fitStep  = 0.05
	// fontScaleStep is the default amount to reduce font sizes per iteration.
	fontScaleStep = 0.25
	// minFontSize is the absolute minimum font size (points) we'll accept.
	minFontSize = 6.0
	// measurePageHeight is large enough that no real CV's content would
	// ever trigger EnsureSpace's page-break check during measurement, so
	// the final cursor Y reports the content's true total height.
	measurePageHeight = 20000.0
)

// FitToOnePage renders build against progressively tighter spacing (from
// 1.0 down to fitFloor, in fitStep increments) until the result fits on
// one real A4 page, then performs the real render at whatever scale
// worked and returns that Doc. If spacing alone isn't enough, it then
// progressively reduces font sizes (respecting minSize from TypographySettings)
// until everything fits.
//
// Only theme.Spacing (the gaps between things) and font sizes are scaled.
// build must be pure beyond drawing into the *draw.Doc it's given: it runs
// once per candidate scale during measurement, then once more for the real render.
func FitToOnePage(th theme.Theme, baseSpacing theme.Spacing, typo TypographySettings, build func(d *draw.Doc) error) (*draw.Doc, error) {
	available := draw.A4Height - th.MarginTop - th.MarginBottom

	// First pass: compress spacing only
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

	// Apply the spacing scale we found
	finalTh := th
	finalTh.Spacing = baseSpacing.Scaled(scale)

	// Check if we still need font scaling
	md, err := draw.New(finalTh, draw.Options{PageHeight: measurePageHeight})
	if err != nil {
		return nil, err
	}
	if err := build(md); err != nil {
		return nil, err
	}

	used := md.Y() - th.MarginTop
	if used <= available {
		// Spacing alone was enough
		d, err := draw.New(finalTh, draw.Options{})
		if err != nil {
			return nil, err
		}
		if err := build(d); err != nil {
			return nil, err
		}
		return d, nil
	}

	// Second pass: also scale fonts
	fontScale := 0.0
	for {
		ty := buildTypographyOverride(th, typo, fontScale)
		measureTh := finalTh.ApplyTypography(ty)
		md, err := draw.New(measureTh, draw.Options{PageHeight: measurePageHeight})
		if err != nil {
			return nil, err
		}
		if err := build(md); err != nil {
			return nil, err
		}

		used := md.Y() - th.MarginTop
		if used <= available {
			break
		}

		// Check if we've hit minimums for all fonts
		if atMinimum(th, typo, fontScale) {
			break
		}

		fontScale += getMinScaleStep(typo)
	}

	// Final render with both spacing and font scaling
	ty := buildTypographyOverride(th, typo, fontScale)
	finalTh = finalTh.ApplyTypography(ty)
	d, err := draw.New(finalTh, draw.Options{})
	if err != nil {
		return nil, err
	}
	if err := build(d); err != nil {
		return nil, err
	}
	return d, nil
}

// buildTypographyOverride creates a TypographyOverride by reducing each font
// size by fontScale points, respecting minimums.
func buildTypographyOverride(th theme.Theme, typo TypographySettings, fontScale float64) theme.TypographyOverride {
	var ty theme.TypographyOverride

	ty.NameSize = scaleFont(th.NameSize, typo.Name, fontScale)
	ty.SubtitleSize = scaleFont(th.SubtitleSize, typo.Subtitle, fontScale)
	ty.ContactSize = scaleFont(th.ContactSize, typo.Contact, fontScale)
	ty.HeadingSize = scaleFont(th.HeadingSize, typo.Heading, fontScale)
	ty.RoleTitleSize = scaleFont(th.RoleTitleSize, typo.RoleTitle, fontScale)
	ty.CompanySize = scaleFont(th.CompanySize, typo.Company, fontScale)
	ty.DateSize = scaleFont(th.DateSize, typo.Date, fontScale)
	ty.StackLabelSize = scaleFont(th.StackLabelSize, typo.StackLabel, fontScale)
	ty.StackValueSize = scaleFont(th.StackValueSize, typo.StackValue, fontScale)
	ty.BodySize = scaleFont(th.BodySize, typo.Body, fontScale)
	ty.SkillsKeySize = scaleFont(th.SkillsKeySize, typo.SkillsKey, fontScale)
	ty.SkillsValueSize = scaleFont(th.SkillsValueSize, typo.SkillsValue, fontScale)

	return ty
}

// scaleFont returns a pointer to the scaled font size, respecting minimums.
func scaleFont(current float64, elem *TypographyElement, scale float64) *float64 {
	minSize := minFontSize
	if elem != nil && elem.MinSize != nil && *elem.MinSize > 0 {
		minSize = *elem.MinSize
	}

	newSize := current - scale
	if newSize < minSize {
		newSize = minSize
	}
	return &newSize
}

// atMinimum checks if all fonts have reached their minimum sizes.
func atMinimum(th theme.Theme, typo TypographySettings, fontScale float64) bool {
	elements := []struct {
		current float64
		elem    *TypographyElement
	}{
		{th.NameSize, typo.Name},
		{th.SubtitleSize, typo.Subtitle},
		{th.ContactSize, typo.Contact},
		{th.HeadingSize, typo.Heading},
		{th.RoleTitleSize, typo.RoleTitle},
		{th.CompanySize, typo.Company},
		{th.DateSize, typo.Date},
		{th.StackLabelSize, typo.StackLabel},
		{th.StackValueSize, typo.StackValue},
		{th.BodySize, typo.Body},
		{th.SkillsKeySize, typo.SkillsKey},
		{th.SkillsValueSize, typo.SkillsValue},
	}

	for _, e := range elements {
		minSize := minFontSize
		if e.elem != nil && e.elem.MinSize != nil && *e.elem.MinSize > 0 {
			minSize = *e.elem.MinSize
		}
		if e.current-fontScale > minSize {
			return false
		}
	}
	return true
}

// getMinScaleStep returns the smallest scale step from all configured elements.
func getMinScaleStep(typo TypographySettings) float64 {
	step := fontScaleStep
	elements := []*TypographyElement{
		typo.Name, typo.Subtitle, typo.Contact, typo.Heading,
		typo.RoleTitle, typo.Company, typo.Date,
		typo.StackLabel, typo.StackValue, typo.Body,
		typo.SkillsKey, typo.SkillsValue,
	}
	for _, elem := range elements {
		if elem != nil && elem.ScaleStep != nil && *elem.ScaleStep > 0 && *elem.ScaleStep < step {
			step = *elem.ScaleStep
		}
	}
	return step
}
