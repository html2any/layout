package main

import (
	"os"

	"github.com/html2any/layout"
	"github.com/html2any/layout/fonts"
	"github.com/html2any/parser"
)

func initRender(format string, fontcache *fonts.Fonts) (render layout.Render) {
	if format == "pdf" {
		render = NewPdfRender(fontcache)
	} else if format == "img" {
		render, _ = NewImageRender(fontcache)
	} else if format == "txt" {
		render = NewTextRender()
	} else {
		panic("The implement Render: pdf, img, txt")
	}
	return
}

func loadFonts() (fontcache *fonts.Fonts) {
	fontcache = fonts.NewFonts()
	for style, fontfile := range map[string]string{
		string(fonts.Regular): "fonts/300.ttf",
		string(fonts.Bold):    "fonts/400.ttf"} {
		if err := fontcache.AddFont("NotoSansSC", fontfile, fonts.FontStyle(style)); err != nil {
			panic(err)
		}
	}
	return
}

func parseHTML2Block() (root *layout.Block) {
	file := "complicated.html"
	if len(os.Args) > 2 {
		file = os.Args[2]
	}
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	root = layout.NewBlock(&cssMap)
	if err := parser.ParseHTML(data, root); err != nil {
		panic(err)
	}
	return
}

func main() {
	format := os.Args[1]
	// loadFonts and CacheAll
	fontcache := loadFonts()
	sfnt := fonts.NewSFNT(fontcache)
	if _, err := sfnt.CacheAll(); err != nil {
		panic(err)
	}
	// init render
	render := initRender(format, fontcache)
	// load html
	root := parseHTML2Block()
	if format == "img" {
		_, err := layout.CacuHeight(root, layout.A4_Width, "NotoSansSC", NewLayout(sfnt))
		if err != nil {
			panic(err)
		}
		render.(*ImageRender).NewCanvas(root.Width, root.Height+root.Top)
		layout.RenderTo(root, render, 0, 0)
		out := render.GetData()
		os.WriteFile("out.png", out, 0644)
	} else if format == "txt" {
		_, err := layout.CacuHeight(root, layout.A4_Width, "NotoSansSC", NewLayout(sfnt))
		if err != nil {
			panic(err)
		}
		layout.RenderTo(root, render, 0, 0)
		out := render.GetData()
		os.WriteFile("out.txt", out, 0644)
	} else if format == "pdf" {
		if _, pages, err := layout.CacuPages(root, layout.A4_Width, "NotoSansSC", NewLayout(sfnt)); err != nil {
			panic(err)
		} else {
			pdf_render := render.(*PdfRender)
			for _, page := range pages {
				pdf_render.pdf.AddPage()
				layout.RenderTo(page, pdf_render, 0, 0)
			}
			out := render.GetData()
			os.WriteFile("out.pdf", out, 0644)
		}
	}
}
