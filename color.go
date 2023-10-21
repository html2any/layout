package layout

import (
	"errors"
	"fmt"
	"image/color"
)

var (
	errInvalidFormat = errors.New("invalid color format")
)

// ParseHexColor: parses a hex color to a color.RGBA
// s must be a hex color. #RRGGBB
// returns color.RGBA, error
func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

// MustParseHexColor: an alias of ParseHexColor, if any error occurs, it will return color.RGBA(0, 0, 0, 255)
func MustParseHexColor(s string) color.RGBA {
	c, err := ParseHexColor(s)
	if err != nil {
		return color.RGBA{0, 0, 0, 255}
	}
	return c
}

// ColorHex: converts a color.RGBA to a hex string
func RGBA2Hex(p color.RGBA) string {
	return fmt.Sprintf("#%.2x%.2x%.2x", p.R, p.G, p.B)
}
