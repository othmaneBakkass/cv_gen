// Package theme holds the plain-value design tokens each CV template is
// built from: a semantic color palette, type sizes, spacing, margins and
// rule weights. It has no dependency on any PDF library, and no dependency
// on any one template — draw.Doc carries a Theme value, and
// internal/pdf/comp reads it from there, so a new template is a new Theme
// plus a new arrangement of comp calls, not a new drawing layer.
//
// Styling is organized as three layers (see SectionOverride and
// Theme.With): Global Template Defaults -> Section Overrides -> Rendered
// Section. The template's Theme is the global default; a section may carry a
// sparse SectionOverride that is merged onto a copy of that theme for the
// duration of that one section's render, so every section is an independent,
// configurable component without any per-section styling being hardcoded
// inside the comp layer.
package theme

import (
	"fmt"
	"strconv"
	"strings"
)

// RGB is a plain 0-255 color triple.
type RGB struct {
	R, G, B uint8
}

// ParseHex parses a "#RRGGBB" or "RRGGBB" hex color into an RGB. Used to
// turn a CV's JSON color override (schema.ColorSettings) into a value
// render.Options can apply to a theme. It is also the single source of
// truth for what counts as a valid hex color: the schema's struct-tag
// validator and schemas/v1.json's regex are both kept in step with it (an
// optional "#" plus exactly six hex digits).
func ParseHex(s string) (RGB, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return RGB{}, fmt.Errorf("invalid hex color %q: want 6 hex digits, e.g. #1F3864", s)
	}
	v, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		return RGB{}, fmt.Errorf("invalid hex color %q: %w", s, err)
	}
	return RGB{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v)}, nil
}

// Palette is the CV's color vocabulary, named by the semantic role each
// color plays rather than by hue or by an abstract "accent/primary" tier.
// Templates set these; sections may override any of them (see
// SectionOverride). Roles intentionally over-separate elements that a given
// template happens to color the same — e.g. Headline and Border share a
// value in every built-in theme, but keeping them distinct lets a section
// recolor its divider without touching its heading text.
type Palette struct {
	Headline    RGB // section headings, the name, the "Stack:" label — structural emphasis
	Subheadline RGB // the tagline/role line under the name
	Body        RGB // body copy, profile text, entry titles, bullets, skills, contact line
	Link        RGB // hyperlinks (LinkedIn/GitHub)
	Border      RGB // header and section rules / dividers
	Background  RGB // page background (drawn only when not white)
}

// Spacing is the vertical rhythm scale. Values are points. A render pass
// can call Scaled to compress everything uniformly (fit-to-one-page, or a
// caller-requested "airy" vs. "dense" density).
type Spacing struct {
	XS float64 // between a bullet and the next bullet
	SM float64 // after a section heading's rule
	MD float64 // before a heading; before a new entry
	LG float64 // header block to first section
}

// Scaled returns the spacing scale multiplied by k, e.g. 0.85 for a denser
// render or 1.15 for an airier one.
func (s Spacing) Scaled(k float64) Spacing {
	return Spacing{XS: s.XS * k, SM: s.SM * k, MD: s.MD * k, LG: s.LG * k}
}

// Theme is everything a template's visual identity needs. Every field is a
// plain value so a new template is just a new Theme literal — no code in
// internal/pdf/draw or internal/pdf/comp changes.
type Theme struct {
	Name       string // shown nowhere; useful in logs/errors
	FontFamily string // must match a family draw.New registers faces under

	// Palette is the semantic color vocabulary (see Palette).
	Palette Palette

	// Type sizes, in points.
	NameSize        float64
	SubtitleSize    float64
	ContactSize     float64
	HeadingSize     float64
	RoleTitleSize   float64
	CompanySize     float64
	DateSize        float64
	StackLabelSize  float64
	StackValueSize  float64
	BodySize        float64
	SkillsKeySize   float64
	SkillsValueSize float64

	// Letter-spacing (tracking), in points, applied to the name and
	// section headings.
	NameTracking    float64
	HeadingTracking float64

	// HeadingUppercase controls whether section headings are upper-cased
	// (t2's editorial look) or left as typed (a plainer, more
	// conventional look for an ATS-oriented template).
	HeadingUppercase bool

	// DateItalic controls whether dates and stack values render italic
	// (t2's editorial style) or upright (plainer, more conservative).
	DateItalic bool

	// ProfileItalic controls whether the profile/summary paragraph renders
	// italic. Kept as a theme flag rather than hardcoded in comp.Profile so
	// no styling decision lives inside the component layer.
	ProfileItalic bool

	// Page geometry, in points (A4 is 595.28 x 841.89pt).
	MarginTop    float64
	MarginBottom float64
	MarginLeft   float64
	MarginRight  float64

	// Rule weights, in points.
	HeaderRuleWeight  float64
	SectionRuleWeight float64

	// Bullet layout, in points, relative to the content's left margin.
	BulletHangingIndent float64
	BulletOffset        float64
	// BulletMarker is the character drawn before each bullet item.
	BulletMarker string

	// SkillsGap is the space between a skills-row key and its value.
	SkillsGap float64

	// MetaSep is the separator placed between the title and the meta field
	// in entry lines (e.g. "  —  " or ": "). Defaults to "  —  ".
	MetaSep string

	// Spacing is this theme's base vertical rhythm before any density
	// adjustment (see Spacing.Scaled and render.Options.Density).
	Spacing Spacing
}

// SectionOverride is a sparse set of per-section style tweaks. It is the
// middle layer of Global Template Defaults -> Section Overrides -> Rendered
// Section: a nil pointer inherits the global theme's value, a set pointer
// replaces it for that one section only. Theme.With merges an override onto
// a copy of the theme; SpaceBefore is consumed separately by
// render.RunSections since it governs the gap before the section rather than
// anything inside the theme.
type SectionOverride struct {
	// Color overrides, one optional pointer per palette role.
	Headline    *RGB
	Subheadline *RGB
	Body        *RGB
	Link        *RGB
	Border      *RGB

	// ShowDivider, when set to false, suppresses this section's rule (it maps
	// to a zero SectionRuleWeight, which draw.Rule renders as nothing) — the
	// "optional visual elements must be removable" requirement. DividerWeight
	// overrides its thickness when set.
	ShowDivider   *bool
	DividerWeight *float64

	// Typography overrides.
	HeadingUppercase *bool
	HeadingSize      *float64
	BodySize         *float64

	// MetaSep overrides the separator between entry title and meta fields.
	// Empty string means "keep the theme default".
	MetaSep *string

	// SpaceBefore overrides the vertical gap render.RunSections inserts
	// before this section (read there, not applied by With).
	SpaceBefore *float64
}

// TypographyOverride holds font-size overrides for all semantic text elements.
// Used by adaptive density to scale fonts down while respecting minimums.
// A nil pointer means "keep the current value".
type TypographyOverride struct {
	NameSize        *float64
	SubtitleSize    *float64
	ContactSize     *float64
	HeadingSize     *float64
	RoleTitleSize   *float64
	CompanySize     *float64
	DateSize        *float64
	StackLabelSize  *float64
	StackValueSize  *float64
	BodySize        *float64
	SkillsKeySize   *float64
	SkillsValueSize *float64
}

// With returns a copy of t with o's set fields applied. Colors land on the
// copy's Palette; ShowDivider/DividerWeight fold into SectionRuleWeight so
// the existing "weight <= 0 means no rule" path in draw.Rule does the work.
// SpaceBefore is intentionally ignored here — it is not a theme property.
func (t Theme) With(o SectionOverride) Theme {
	applyRGB(&t.Palette.Headline, o.Headline)
	applyRGB(&t.Palette.Subheadline, o.Subheadline)
	applyRGB(&t.Palette.Body, o.Body)
	applyRGB(&t.Palette.Link, o.Link)
	applyRGB(&t.Palette.Border, o.Border)

	if o.ShowDivider != nil && !*o.ShowDivider {
		t.SectionRuleWeight = 0
	}
	if o.DividerWeight != nil {
		t.SectionRuleWeight = *o.DividerWeight
	}
	if o.HeadingUppercase != nil {
		t.HeadingUppercase = *o.HeadingUppercase
	}
	if o.HeadingSize != nil {
		t.HeadingSize = *o.HeadingSize
	}
	if o.BodySize != nil {
		t.BodySize = *o.BodySize
	}
	if o.MetaSep != nil {
		t.MetaSep = *o.MetaSep
	}
	return t
}

func applyRGB(dst *RGB, src *RGB) {
	if src != nil {
		*dst = *src
	}
}

// ApplyTypography returns a copy of t with ty's set fields applied.
// Used by adaptive density to progressively scale font sizes.
func (t Theme) ApplyTypography(ty TypographyOverride) Theme {
	if ty.NameSize != nil {
		t.NameSize = *ty.NameSize
	}
	if ty.SubtitleSize != nil {
		t.SubtitleSize = *ty.SubtitleSize
	}
	if ty.ContactSize != nil {
		t.ContactSize = *ty.ContactSize
	}
	if ty.HeadingSize != nil {
		t.HeadingSize = *ty.HeadingSize
	}
	if ty.RoleTitleSize != nil {
		t.RoleTitleSize = *ty.RoleTitleSize
	}
	if ty.CompanySize != nil {
		t.CompanySize = *ty.CompanySize
	}
	if ty.DateSize != nil {
		t.DateSize = *ty.DateSize
	}
	if ty.StackLabelSize != nil {
		t.StackLabelSize = *ty.StackLabelSize
	}
	if ty.StackValueSize != nil {
		t.StackValueSize = *ty.StackValueSize
	}
	if ty.BodySize != nil {
		t.BodySize = *ty.BodySize
	}
	if ty.SkillsKeySize != nil {
		t.SkillsKeySize = *ty.SkillsKeySize
	}
	if ty.SkillsValueSize != nil {
		t.SkillsValueSize = *ty.SkillsValueSize
	}
	return t
}

// white is the default page background — pure white, never off-white/cream.
var white = RGB{255, 255, 255}

// T2 is the "Navy Rule" theme: editorial, teal-accented, tracked uppercase
// headings, italic dates.
var T2 = Theme{
	Name:       "t2-navy-rule",
	FontFamily: "Carlito",

	Palette: Palette{
		Headline:    RGB{31, 86, 115},
		Subheadline: RGB{26, 26, 26},
		Body:        RGB{26, 26, 26},
		Link:        RGB{26, 26, 26},
		Border:      RGB{31, 86, 115},
	},

	NameSize:        20.0,
	SubtitleSize:    11.5,
	ContactSize:     9.5,
	HeadingSize:     11.0,
	RoleTitleSize:   10.5,
	CompanySize:     10.5,
	DateSize:        9.0,
	StackLabelSize:  9.0,
	StackValueSize:  9.0,
	BodySize:        10.0,
	SkillsKeySize:   10.0,
	SkillsValueSize: 10.0,

	NameTracking:    0.5,
	HeadingTracking: 0.5,

	HeadingUppercase: true,
	DateItalic:       true,
	ProfileItalic:    true,

	MarginTop:    24.0,
	MarginBottom: 24.0,
	MarginLeft:   42.5,
	MarginRight:  42.5,

	HeaderRuleWeight:  1.0,
	SectionRuleWeight: 1.0,

	BulletHangingIndent: 13.0,
	BulletOffset:        3.0,
	BulletMarker:        "•",
	SkillsGap:           4.0,
	MetaSep:             "  —  ",

	Spacing: Spacing{XS: 1.5, SM: 3.0, MD: 7.0, LG: 12.0},
}

// T1 is the "ATS Classic" theme: a conservative, single-hue, non-italic,
// slightly larger-type design aimed at the two things t2 doesn't
// prioritize — passing keyword-matching ATS software cleanly, and reading
// as a conventional professional CV in the Moroccan/French hiring
// convention (reverse-chronological, sober navy, no tracked/editorial
// flourishes). Structural ATS-safety (no real tables, no multi-column
// body, text drawn in reading order) comes from the shared draw/comp
// layer, not from this theme.
var T1 = Theme{
	Name:       "t1-ats-classic",
	FontFamily: "Carlito",

	Palette: Palette{
		Headline:    RGB{31, 56, 100},
		Subheadline: RGB{26, 26, 26},
		Body:        RGB{26, 26, 26},
		Link:        RGB{26, 26, 26},
		Border:      RGB{31, 56, 100},
	},

	NameSize:        18.0,
	SubtitleSize:    11.5,
	ContactSize:     10.0,
	HeadingSize:     11.0,
	RoleTitleSize:   10.5,
	CompanySize:     10.5,
	DateSize:        10.0,
	StackLabelSize:  10.0,
	StackValueSize:  10.0,
	BodySize:        10.5,
	SkillsKeySize:   10.0,
	SkillsValueSize: 10.0,

	NameTracking:    0, // no tracking — plain, conventional
	HeadingTracking: 0,

	HeadingUppercase: true, // still expected by ATS parsers and recruiters
	DateItalic:       false,
	ProfileItalic:    false,

	MarginTop:    28.0,
	MarginBottom: 28.0,
	MarginLeft:   46.0,
	MarginRight:  46.0,

	HeaderRuleWeight:  1.0,
	SectionRuleWeight: 0.5,

	BulletHangingIndent: 14.0,
	BulletOffset:        3.0,
	BulletMarker:        "•",
	SkillsGap:           4.0,
	MetaSep:             "  —  ",

	Spacing: Spacing{XS: 2.0, SM: 3.5, MD: 8.0, LG: 13.0},
}

// T3 is the "Modern Minimal" theme: a compact, indigo-accented design with
// title-case (not tracked-uppercase) headings that carry no rule line
// (hierarchy comes from weight/color/spacing alone — SectionRuleWeight is
// 0), an en-dash bullet marker instead of a round one, and tighter overall
// spacing, for an owner who wants a visibly different, contemporary look
// rather than t2's editorial style or t1's conservative one.
var T3 = Theme{
	Name:       "t3-modern-minimal",
	FontFamily: "Carlito",

	Palette: Palette{
		Headline:    RGB{55, 48, 163},
		Subheadline: RGB{23, 23, 23},
		Body:        RGB{23, 23, 23},
		Link:        RGB{55, 48, 163},
		Border:      RGB{55, 48, 163},
	},

	NameSize:        19.0,
	SubtitleSize:    11.0,
	ContactSize:     9.5,
	HeadingSize:     11.0,
	RoleTitleSize:   10.5,
	CompanySize:     10.5,
	DateSize:        9.5,
	StackLabelSize:  9.5,
	StackValueSize:  9.5,
	BodySize:        10.0,
	SkillsKeySize:   9.5,
	SkillsValueSize: 9.5,

	NameTracking:    0,
	HeadingTracking: 0,

	HeadingUppercase: false, // title case, not tracked caps — plainer hierarchy signal
	DateItalic:       false,
	ProfileItalic:    true,

	MarginTop:    20.0,
	MarginBottom: 20.0,
	MarginLeft:   36.0,
	MarginRight:  36.0,

	HeaderRuleWeight:  1.0,
	SectionRuleWeight: 0, // no rule under headings — weight/color/spacing carry hierarchy

	BulletHangingIndent: 12.0,
	BulletOffset:        3.0,
	BulletMarker:        "–", // en dash, not a round bullet
	SkillsGap:           4.0,
	MetaSep:             "  —  ",

	Spacing: Spacing{XS: 1.2, SM: 2.5, MD: 6.0, LG: 10.0},
}
