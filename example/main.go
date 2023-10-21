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
	textspliter := fonts.NewSFNT(fontcache)
	if _, err := textspliter.CacheAll(); err != nil {
		panic(err)
	}
	// init render
	render := initRender(format, fontcache)

	// load html
	root := parseHTML2Block()

	_, err := layout.CacuHeight(root, layout.A4_Width, "NotoSansSC", textspliter)
	if err != nil {
		panic(err)
	}
	if format == "img" {
		render.(*ImageRender).NewCanvas(root.Width, root.Height+root.Top)
		layout.RenderTo(root, render, 0, 0)
		out := render.GetData()
		os.WriteFile("out.png", out, 0644)
	}
}
