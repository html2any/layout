package main

import (
	"bufio"
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"strconv"

	"github.com/html2any/layout"
	"github.com/html2any/layout/fonts"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/rasterizer"
)

type ImageRender struct {
	c       *canvas.Canvas
	cctx    *canvas.Context
	font    *fonts.Fonts
	fontMap map[string]*canvas.FontFamily
}

func NewImageRender(font *fonts.Fonts) (*ImageRender, error) {
	p := &ImageRender{font: font}
	if err := p.initCanvas(); err != nil {
		return p, err
	}
	return p, nil
}

func (p *ImageRender) NewCanvas(width, height float64) {
	cW := math.Ceil(width * scale)
	cH := height * scale
	p.c = canvas.New(cW, cH)
	p.cctx = canvas.NewContext(p.c)
	p.cctx.SetCoordSystem(canvas.CartesianIV)
	p.cctx.SetFillColor(color.NRGBA{R: 255, G: 254, B: 254, A: 255})
	p.cctx.DrawPath(0, 0, canvas.Rectangle(cW, cH))
}

func (p *ImageRender) initCanvas() error {
	p.fontMap = make(map[string]*canvas.FontFamily)
	for family, fdata := range p.font.GetFonts() {
		fontf := canvas.NewFontFamily(family.Name)
		style := canvas.FontRegular
		if family.Style == fonts.Bold {
			style = canvas.FontBold
		}
		if err := fontf.LoadFont(fdata, 0, style); err != nil {
			return err
		}
		key := family.Name + "-" + strconv.Itoa(int(style))
		p.fontMap[key] = fontf
	}
	return nil
}

func (p *ImageRender) GetJPEG() []byte {
	resolution := canvas.DPMM(1.0)
	colorSpace := canvas.DefaultColorSpace
	img := rasterizer.Draw(p.c, resolution, colorSpace)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := jpeg.Encode(w, img, &jpeg.Options{Quality: 90}); err != nil {
		return nil
	} else {
		return b.Bytes()
	}
}
func (p *ImageRender) GetData() []byte {
	return p.GetPNG()
}
func (p *ImageRender) GetPNG() []byte {
	resolution := canvas.DPMM(1.0)
	colorSpace := canvas.DefaultColorSpace
	img := rasterizer.Draw(p.c, resolution, colorSpace)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	e := &png.Encoder{CompressionLevel: png.BestSpeed}
	if err := e.Encode(w, img); err != nil {
		return nil
	} else {
		return b.Bytes()
	}
}

func (p *ImageRender) GetGIF() []byte {
	resolution := canvas.DPMM(1.0)
	colorSpace := canvas.DefaultColorSpace
	img := rasterizer.Draw(p.c, resolution, colorSpace)
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	if err := gif.Encode(w, img, &gif.Options{NumColors: 256}); err != nil {
		return nil
	} else {
		return b.Bytes()
	}
}

const scale = 2

func (p *ImageRender) transformRect(rect *layout.Rect) {
	rect.X = rect.X * scale
	rect.Y = rect.Y * scale
	rect.W = rect.W * scale
	rect.H = rect.H * scale
}
func (p *ImageRender) Text(rect *layout.Rect, text string, family string, style string, size int, rgba color.RGBA) error {
	p.transformRect(rect)
	fstyle := canvas.FontRegular
	if style == "bold" {
		fstyle = canvas.FontBold
	} else if style == "italic" {
		fstyle = canvas.FontItalic
	}
	key := family + "-" + strconv.Itoa(int(fstyle))
	if fontf, ok := p.fontMap[key]; !ok {
		return errors.New("font not found")
	} else {
		scale_size := float64(size) * scale
		baseline := scale_size * 0.1
		txt := canvas.NewTextLine(fontf.Face(float64(size)*72/25.4*scale, fstyle, rgba), text, canvas.Left)
		p.cctx.DrawText(rect.X, rect.Y+scale_size-baseline, txt)
	}
	return nil
}

func (p *ImageRender) Fill(rgba color.RGBA, rect *layout.Rect) {
	p.transformRect(rect)
	p.cctx.SetStrokeColor(canvas.Transparent)
	p.cctx.SetFillColor(rgba)
	p.cctx.DrawPath(rect.X, rect.Y, canvas.Rectangle(rect.W, rect.H))
}

func (p *ImageRender) Line(rgba color.RGBA, linetype string, width, x1, y1, x2, y2 float64) {
	// fmt.Printf("Line:%s,Type:%s,Width:%.1f,From(X:%.1f,Y:%.1f) To(X:%.1f,Y:%.1f)\n", utils.Color(rgba), linetype, width, x1, y1, x2, y2)
	x1 = x1 * scale
	y1 = y1 * scale
	x2 = x2 * scale
	y2 = y2 * scale

	p.cctx.SetFillColor(canvas.Transparent)
	p.cctx.SetStrokeColor(rgba)
	p.cctx.SetStrokeWidth(width * scale)
	if linetype == "solid" {
		p.cctx.SetDashes(0.0)
	} else {
		p.cctx.SetDashes(0.0, 1.0)
	}
	l := &canvas.Path{}
	l.MoveTo(x1, y1)
	l.LineTo(x2, y2)
	p.cctx.DrawPath(0, 0, l)
}

func (p *ImageRender) Image(img string, rect *layout.Rect) error {
	p.transformRect(rect)
	var imgData *[]byte
	if data, err := os.ReadFile(img); err != nil {
		return err
	} else {
		imgData = &data
	}
	if data, _, err := image.Decode(bytes.NewReader(*imgData)); err != nil {
		return err
	} else {
		p.cctx.FitImage(data, canvas.Rect{W: rect.W, H: rect.H, X: rect.X, Y: rect.Y}, canvas.ImageFill)
	}
	return nil
}
