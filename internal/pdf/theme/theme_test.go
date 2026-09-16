package theme

import "testing"

func TestParseHex(t *testing.T) {
	cases := []struct {
		in      string
		want    RGB
		wantErr bool
	}{
		{"#1F3864", RGB{31, 56, 100}, false},
		{"1F3864", RGB{31, 56, 100}, false}, // "#" is optional
		{"#000000", RGB{0, 0, 0}, false},
		{"#FFFFFF", RGB{255, 255, 255}, false},
		{"#8b0000", RGB{139, 0, 0}, false}, // lowercase hex digits
		{"", RGB{}, true},
		{"#FFF", RGB{}, true},     // 3-digit shorthand not supported
		{"#GGGGGG", RGB{}, true},  // not hex digits
		{"#1F38644", RGB{}, true}, // too long
		{"not-a-color", RGB{}, true},
	}
	for _, c := range cases {
		got, err := ParseHex(c.in)
		if c.wantErr {
			if err == nil {
				t.Errorf("ParseHex(%q): expected an error, got %+v", c.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseHex(%q): unexpected error: %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseHex(%q) = %+v, want %+v", c.in, got, c.want)
		}
	}
}
