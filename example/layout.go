package main

import (
	"github.com/html2any/layout"
	"github.com/html2any/layout/fonts"
)

type Layout struct {
	sfnt *fonts.SFNT
}

func NewLayout(sfnt *fonts.SFNT) *Layout {
	return &Layout{sfnt: sfnt}
}
func (l *Layout) Overide(block *layout.Block) bool {
	if block.Id == "title" {
		block.Contents = []string{"Hello World!"}
		return false
	}
	return false
}

func (l *Layout) SplitLines(width float64, s string, family string, size float64, style string) (lines []string, line_widths []float64, err error) {
	return l.sfnt.SplitLines(width, s, family, size, style)
}
