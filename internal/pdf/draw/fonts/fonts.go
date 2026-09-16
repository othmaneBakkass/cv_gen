// Package fonts embeds the Carlito TTF faces used by the draw layer.
// Carlito is metric-compatible with Calibri and is licensed under the SIL
// Open Font License (see carlito/OFL.txt).
package fonts

import _ "embed"

//go:embed carlito/Carlito-Regular.ttf
var CarlitoRegular []byte

//go:embed carlito/Carlito-Bold.ttf
var CarlitoBold []byte

//go:embed carlito/Carlito-Italic.ttf
var CarlitoItalic []byte

//go:embed carlito/Carlito-BoldItalic.ttf
var CarlitoBoldItalic []byte
