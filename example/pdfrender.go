package main

import (
	"image/color"

	"github.com/html2any/layout"
	"github.com/html2any/layout/fonts"
	"github.com/signintech/gopdf"
)

type PdfRender struct {
	pdf gopdf.GoPdf
}

func NewPdfRender(font *fonts.Fonts) *PdfRender {
	p := &PdfRender{}
	if err := p.initGoPDF(font, layout.Rect{W: layout.A4_Width, H: layout.A4_Height}); err != nil {
		panic(err)
	}
	return p
}

func (p *PdfRender) initGoPDF(font *fonts.Fonts, rect layout.Rect) error {
	p.pdf = gopdf.GoPdf{}
	p.pdf.SetNoCompression()
	p.pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: rect.W, H: rect.H}}) // = A4
	for family, fdata := range font.GetFonts() {
		style := gopdf.Regular
		if family.Style == fonts.Bold {
			style = gopdf.Bold
		}
		if err := p.pdf.AddTTFFontDataWithOption(family.Name, fdata, gopdf.TtfOption{Style: style}); err != nil {
			return err
		}
	}
	return nil
}

func (p *PdfRender) GetData() []byte {
	return p.pdf.GetBytesPdf()
}

func (p *PdfRender) Overide(block *layout.Block, border, background, content layout.Rect) bool {
	return false
}

func (p *PdfRender) Text(rect *layout.Rect, text string, family string, style string, size int, rgba color.RGBA) error {
	if style == "bold" {
		style = "B"
	} else {
		style = ""
	}
	if err := p.pdf.SetFont(family, style, size); err != nil {
		return err
	}
	p.pdf.SetTextColor(rgba.R, rgba.G, rgba.B)
	p.pdf.SetX(rect.X)
	p.pdf.SetY(rect.Y)
	return p.pdf.Cell(&gopdf.Rect{W: rect.W, H: rect.H}, text)
}

func (p *PdfRender) Fill(rgba color.RGBA, rect *layout.Rect) {
	p.pdf.SetStrokeColor(rgba.R, rgba.G, rgba.B)
	p.pdf.SetFillColor(rgba.R, rgba.G, rgba.B)
	p.pdf.RectFromUpperLeftWithStyle(rect.X, rect.Y, rect.W, rect.H, "F")
}

func (p *PdfRender) Line(rgba color.RGBA, linetype string, width, x1, y1, x2, y2 float64) {
	p.pdf.SetLineWidth(width)
	p.pdf.SetLineType(linetype)
	p.pdf.SetStrokeColor(rgba.R, rgba.G, rgba.B)
	p.pdf.Line(x1, y1, x2, y2)
}

func (p *PdfRender) Image(img string, rect *layout.Rect) error {
	return p.pdf.Image(img, rect.X, rect.Y, &gopdf.Rect{W: rect.W, H: rect.H})
}
