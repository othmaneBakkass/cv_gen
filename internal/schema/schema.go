// Package schema defines the CV input format and its validation rules.
package schema

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	apperror "github.com/othmaneBakkass/cv_gen/internal/common/appError"
	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// InputData is the top-level structure of an input JSON file.
type InputData struct {
	Data []CV `json:"data"`
}

// CV is a single CV entry: the data for one generated document.
type CV struct {
	Template  string      `json:"template" validate:"required"`
	FileName  string      `json:"fileName" validate:"required"`
	Head      Head        `json:"head" validate:"required"`
	Profile   string      `json:"profile,omitempty"`
	Education []Education `json:"education" validate:"required,min=1,dive"`
	Jobs      []Job       `json:"jobs" validate:"required,min=1,dive"`
	Skills    []Skill     `json:"skills,omitempty" validate:"omitempty,dive"`
	Projects  []Project   `json:"projects,omitempty" validate:"omitempty,dive"`
	Languages []Language  `json:"languages" validate:"required,min=1,dive"`
	// Certifications is a flat list of certification titles (e.g. "AWS
	// Certified Solutions Architect — Amazon"), rendered one per line.
	Certifications []string `json:"certifications,omitempty" validate:"omitempty,dive,required"`
	// Settings holds this entry's own presentation options (language,
	// section visibility/order, color overrides) — see Settings. Optional;
	// an absent Settings, or an absent field within it, means "use the
	// template's default". CLI flags (see cmd/generate) override whatever
	// is set here when both are given for the same thing.
	Settings *Settings `json:"settings,omitempty" validate:"omitempty"`
}

// Settings holds per-render presentation options that are orthogonal to
// the CV's content. Every field is optional and independent, matching the
// "empty means skip, filled means include" pattern used throughout this
// schema — set only the ones you want to override.
type Settings struct {
	// Lang selects the template's own label language: "en" or "fr". Does
	// not translate your content, only headings/labels like "Stack:".
	Lang string `json:"lang,omitempty" validate:"omitempty,oneof=en fr"`
	// Density scales the template's spacing: "dense", "normal", or "airy".
	Density string `json:"density,omitempty" validate:"omitempty,oneof=dense normal airy"`
	// Sections controls per-section visibility and render order.
	Sections *SectionSettings `json:"sections,omitempty"`
	// Colors overrides the template's palette. Omit any field to keep the
	// template's default for that color.
	Colors *ColorSettings `json:"colors,omitempty"`
}

// SectionSettings holds one SectionSetting per optional/orderable section.
// Education, Jobs (Experience) and Languages have no Include here because
// the schema already requires at least one entry for those — there's
// nothing to optionally hide — but they can still be reordered.
type SectionSettings struct {
	Profile        *SectionSetting `json:"profile,omitempty"`
	Education      *SectionSetting `json:"education,omitempty"`
	Experience     *SectionSetting `json:"experience,omitempty"`
	Skills         *SectionSetting `json:"skills,omitempty"`
	Projects       *SectionSetting `json:"projects,omitempty"`
	Certifications *SectionSetting `json:"certifications,omitempty"`
	Languages      *SectionSetting `json:"languages,omitempty"`
}

// SectionSetting overrides one section's visibility, render-order priority
// (lower renders first), and/or its per-section styling. Any subset of
// fields may be set.
type SectionSetting struct {
	Include  *bool         `json:"include,omitempty"`
	Priority *int          `json:"priority,omitempty"`
	Style    *SectionStyle `json:"style,omitempty"`
}

// SectionStyle is the user-facing per-section style override — the "Section
// Overrides" layer of Global Template Defaults -> Section Overrides ->
// Rendered Section. Every field is optional; an unset field inherits the
// template's value for that section. It maps directly onto
// theme.SectionOverride at render time (see cmd/generate).
type SectionStyle struct {
	// Divider, when false, removes this section's heading rule; when true,
	// forces it on. Omit to keep the template's default.
	Divider *bool `json:"divider,omitempty"`
	// SpaceBefore overrides the vertical gap (in points) before the section.
	SpaceBefore *float64 `json:"spaceBefore,omitempty"`
	// Uppercase overrides whether this section's heading is upper-cased.
	Uppercase *bool `json:"uppercase,omitempty"`
	// Colors overrides individual palette roles for this section only. Same
	// shape (and same semantic-vs-alias rules) as the global Settings.Colors.
	Colors *ColorSettings `json:"colors,omitempty"`
}

// ColorSettings overrides the active template's semantic color palette (see
// theme.Palette). Each value is a hex color, e.g. "#1F3864" (the "#" is
// optional). Every field is independent — set only the roles you want to
// change; the rest keep the template's default.
//
// The four deprecated fields (Accent/Grey/Muted/Ink) are the original
// palette API, kept working so existing input files don't break. They map
// onto the semantic roles as: accent -> headline + border, grey ->
// metadata, muted -> profile, ink -> body. If both an alias and its
// semantic equivalent are set, the semantic value wins (see
// cmd/generate.colorOverrideFromSettings).
type ColorSettings struct {
	// Semantic role names — preferred.
	Headline    string `json:"headline,omitempty" validate:"omitempty,hexcolor"`
	Subheadline string `json:"subheadline,omitempty" validate:"omitempty,hexcolor"`
	Body        string `json:"body,omitempty" validate:"omitempty,hexcolor"`
	Metadata    string `json:"metadata,omitempty" validate:"omitempty,hexcolor"`
	Link        string `json:"link,omitempty" validate:"omitempty,hexcolor"`
	Border      string `json:"border,omitempty" validate:"omitempty,hexcolor"`
	Profile     string `json:"profile,omitempty" validate:"omitempty,hexcolor"`
	Background  string `json:"background,omitempty" validate:"omitempty,hexcolor"`

	// Deprecated aliases (original four-color API). Prefer the semantic
	// roles above; these are retained for backward compatibility.
	Accent string `json:"accent,omitempty" validate:"omitempty,hexcolor"`
	Grey   string `json:"grey,omitempty" validate:"omitempty,hexcolor"`
	Muted  string `json:"muted,omitempty" validate:"omitempty,hexcolor"`
	Ink    string `json:"ink,omitempty" validate:"omitempty,hexcolor"`
}

// Head holds the contact block shown at the top of the CV. FullName,
// Address, Phone and Email are the only fields every template can rely
// on; the rest are optional extras that richer templates (e.g. t2) may
// render and simpler ones may ignore.
type Head struct {
	FullName string `json:"fullName" validate:"required"`
	Address  string `json:"address" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Email    string `json:"email" validate:"required,email"`

	JobTitle    string   `json:"jobTitle,omitempty"`
	Specialties []string `json:"specialties,omitempty"`
	LinkedIn    string   `json:"linkedin,omitempty"`
	GitHub      string   `json:"github,omitempty"`
	// Photo is a path to a local image file, resolved relative to the
	// input JSON's directory. Templates that support a photo embed it
	// when set and silently skip it when empty.
	Photo string `json:"photo,omitempty"`
}

// Skill is one category in the technical skills section (e.g.
// "Languages" -> ["Go", "TypeScript"]).
type Skill struct {
	Category string   `json:"category" validate:"required"`
	Items    []string `json:"items" validate:"required,min=1,dive,required"`
}

// Project is a single entry in the projects section.
type Project struct {
	Name        string `json:"name" validate:"required"`
	Tech        string `json:"tech" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Education is a single entry in the education section.
type Education struct {
	School      string `json:"school" validate:"required"`
	Location    string `json:"location" validate:"required"`
	StartedAt   string `json:"startedAt" validate:"required"`
	EndedAt     string `json:"endedAt" validate:"required"`
	Degree      string `json:"degree" validate:"required"`
	Description string `json:"description" validate:"required"`
}

// Job is a single entry in the experience section.
type Job struct {
	Company    string   `json:"company" validate:"required"`
	Location   string   `json:"location" validate:"required"`
	Position   string   `json:"position" validate:"required"`
	StartedAt  string   `json:"startedAt" validate:"required"`
	EndedAt    string   `json:"endedAt" validate:"required"`
	Tools      []string `json:"tools" validate:"required,dive,required"`
	Highlights []string `json:"highlights" validate:"required,dive,required"`
}

// Language is a single entry in the languages section.
type Language struct {
	Language string `json:"language" validate:"required"`
	Level    string `json:"level" validate:"required"`
}

// validate is reused across calls; it caches struct metadata and is
// safe for concurrent use. It reports field names using their `json`
// tag so messages match what the user wrote.
var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	// Override the built-in "hexcolor" so struct-tag validation accepts
	// exactly what theme.ParseHex (the real render-time parser, in
	// cmd/generate's colorOverrideFromSettings) and schemas/v1.json's
	// "^#?[0-9a-fA-F]{6}$" accept: an optional "#" plus exactly six hex
	// digits. go-playground's built-in requires a leading "#" and also
	// allows 3/4/8-digit forms, so it would reject valid values like
	// "1F3864" and accept "#FFF"/"#1F3864FF" that ParseHex later rejects at
	// render time — a divergence between what validates and what renders.
	// Delegating to ParseHex keeps all three layers identical by
	// construction.
	if err := v.RegisterValidation("hexcolor", func(fl validator.FieldLevel) bool {
		s := fl.Field().String()
		if s == "" {
			return true // empty is handled by the omitempty prefix on each tag
		}
		_, err := theme.ParseHex(s)
		return err == nil
	}); err != nil {
		panic(fmt.Sprintf("registering hexcolor validator: %v", err))
	}
	return v
}

// Validate checks cv against the schema rules. On failure it returns an
// apperror.AppError whose issues describe every field that is missing or
// invalid.
func Validate(cv *CV) error {
	err := validate.Struct(cv)
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return apperror.New(
			"Internal validation failure",
			"An unexpected error occurred while validating the JSON file.",
			apperror.ErrorCodeUnknown,
			apperror.ErrorSensitivityPublic,
		)
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return apperror.New(
			"Validation failed",
			err.Error(),
			apperror.ErrorCodeUnknown,
			apperror.ErrorSensitivityPublic,
		)
	}

	issues := make([]apperror.AppErrorIssue, 0, len(validationErrors))
	for _, ve := range validationErrors {
		// Namespace looks like "CV.education[0].school"; drop the leading
		// struct name so it reads like the user's JSON.
		field := strings.TrimPrefix(ve.Namespace(), "CV.")
		issues = append(issues, apperror.AppErrorIssue{
			Title:       fmt.Sprintf("Invalid field: %s", field),
			Detail:      fmt.Sprintf("%q is missing or invalid (rule: %s)", field, ve.Tag()),
			Sensitivity: apperror.ErrorSensitivityPublic,
		})
	}

	return apperror.New(
		"Invalid JSON format",
		"The JSON file does not match the expected format. See schemas/v1.json.",
		apperror.ErrorCodeArgs,
		apperror.ErrorSensitivityPublic,
	).WithIssues(issues...)
}
