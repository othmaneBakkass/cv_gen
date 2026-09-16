package schema

import (
	"testing"

	"github.com/othmaneBakkass/cv_gen/internal/pdf/theme"
)

// TestHexcolorTagMatchesParseHex is the guard for the schema-consistency
// requirement: the struct-tag "hexcolor" validator (used on
// settings.colors.*), the render-time parser theme.ParseHex, and
// schemas/v1.json's "^#?[0-9a-fA-F]{6}$" pattern must all accept and reject
// the exact same set of strings. Before the custom validator was
// registered, go-playground's built-in "hexcolor" diverged from both on
// these cases (it required a leading "#" and allowed 3/4/8-digit forms).
func TestHexcolorTagMatchesParseHex(t *testing.T) {
	cases := []struct {
		in   string
		want bool // true = should be accepted
	}{
		{"#1F3864", true},    // canonical, with hash, uppercase
		{"1F3864", true},     // no hash — v1.json's "#?" and ParseHex allow it
		{"#abcdef", true},    // lowercase
		{"#000000", true},    // all zeros
		{"#FFFFFF", true},    // all F
		{"#FFF", false},      // 3-digit shorthand — rejected by ParseHex + v1.json
		{"#1234", false},     // 4-digit shorthand
		{"#1F3864FF", false}, // 8-digit (alpha) — rejected: exactly 6 required
		{"#GGGGGG", false},   // non-hex digits
		{"1F386", false},     // 5 digits
		{"##FFFFFF", false},  // double hash
		{"not-a-color", false},
	}

	for _, c := range cases {
		// Struct-tag path: exactly the tag used on ColorSettings fields.
		tagErr := validate.Var(c.in, "hexcolor")
		tagOK := tagErr == nil

		// Render-time parser path.
		_, parseErr := theme.ParseHex(c.in)
		parseOK := parseErr == nil

		if tagOK != c.want {
			t.Errorf("hexcolor tag: %q got accepted=%v, want %v", c.in, tagOK, c.want)
		}
		if parseOK != c.want {
			t.Errorf("ParseHex: %q got accepted=%v, want %v", c.in, parseOK, c.want)
		}
		// The core invariant: the two agree on every input.
		if tagOK != parseOK {
			t.Errorf("divergence on %q: tag accepted=%v but ParseHex accepted=%v", c.in, tagOK, parseOK)
		}
	}
}

// TestHexcolorOmitemptyStillSkips confirms the "omitempty" prefix on the
// real ColorSettings tags means an empty string is skipped (not rejected),
// so an omitted color field keeps the template default rather than erroring.
func TestHexcolorOmitemptyStillSkips(t *testing.T) {
	if err := validate.Var("", "omitempty,hexcolor"); err != nil {
		t.Errorf("empty color with omitempty should pass, got %v", err)
	}
}
