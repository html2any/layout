package main

import (
	"bytes"
	"fmt"
	"image/color"
	"strconv"

	"github.com/html2any/layout"
)

type TextRender struct {
	buf bytes.Buffer
}

func NewTextRender() *TextRender {
	p := &TextRender{}
	return p
}

func (p *TextRender) GetData() []byte {
	return p.buf.Bytes()
}

func (p *TextRender) Overide(block *layout.Block, border, background, content layout.Rect) bool {
	return false
}

func (p *TextRender) Text(rect *layout.Rect, text string, family string, style string, size int, rgba color.RGBA) error {
	p.buf.WriteString("DrawText=>" + rect.String() + ",Family:" + family + ",Style:" + style + ",Size:" + strconv.Itoa(size) + ",Color" + layout.RGBA2Hex(rgba) + ",Text:" + text + "\n")
	return nil
}

func (p *TextRender) Fill(rgba color.RGBA, rect *layout.Rect) {
	p.buf.WriteString("FillBgrd=>" + rect.String() + ",Fill:" + layout.RGBA2Hex(rgba) + "\n")
}

func (p *TextRender) Line(rgba color.RGBA, linetype string, width, x1, y1, x2, y2 float64) {
	p.buf.WriteString(fmt.Sprintf("DrawLine=>From:{X:%.1f,Y:%.1f},To:{X:%.1f,Y:%.1f},Color:%s,Type:%s,Width:%.1f\n", x1, y1, x2, y2, layout.RGBA2Hex(rgba), linetype, width))
}

func (p *TextRender) Image(img string, rect *layout.Rect) error {
	p.buf.WriteString("DrawImage=>" + rect.String() + ",Image:" + img + "\n")
	return nil
}
