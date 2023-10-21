package layout

import (
	"fmt"
	"image/color"
)

type Rect struct {
	X float64
	Y float64
	W float64
	H float64
}

func (r Rect) String() string {
	return fmt.Sprintf("Rect:{X:%.1f, Y:%.1f, W:%.1f, H:%.1f}", r.X, r.Y, r.W, r.H)
}

type Render interface {
	Text(rect *Rect, text string, family string, style string, size int, rgba color.RGBA) error
	Fill(rgba color.RGBA, rect *Rect)
	Image(img string, rect *Rect) error
	Line(clr color.RGBA, linetype string, width, x1, y1, x2, y2 float64)
	GetData() []byte
}
