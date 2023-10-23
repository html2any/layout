package layout

type ILayout interface {
	SplitLines(width float64, s string, family string, size float64, style string) (lines []string, line_widths []float64, err error)
	Overide(block *Block) bool
}
