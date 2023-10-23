package layout

func addY(rect Rect, top float64) *Rect {
	return &Rect{X: rect.X, Y: rect.Y + top, W: rect.W, H: rect.H}
}
func RenderTo(c *Block, render Render, from_x, from_y float64) (Rect, error) {
	border, background, content := c.cacuRect(from_x, from_y)
	rbackground := addY(background, c.Top)
	if render.Overide(c, border, *rbackground, content) {
		return Rect{X: from_x, Y: from_y, W: float64(c.Width), H: c.Height}, nil
	}
	if c.parent != nil && c.BackgroundColor != c.parent.BackgroundColor && c.BackgroundColor.A != 0 {
		render.Fill(c.BackgroundColor, rbackground)
	} else if c.parent == nil {
		render.Fill(c.BackgroundColor, rbackground)
	}
	if len(c.BackgroundImage) > 0 {
		if err := render.Image(c.BackgroundImage, rbackground); err != nil {
			return Rect{}, err
		}
	}
	if len(c.Src) > 0 {
		if err := render.Image(c.Src, rbackground); err != nil {
			return Rect{}, err
		}
	}
	//Render Content or Children
	if len(c.Contents) > 0 {
		if err := renderContents(c, render, content); err != nil {
			return Rect{}, err
		}
	} else if len(c.Children) > 0 {
		if c.FlexDirection == FlexDirectionColumn {
			if err := renderRows(c, render, content); err != nil {
				return Rect{}, err
			}
		} else {
			if err := renderColumns(c, render, content); err != nil {
				return Rect{}, err
			}
		}
	}
	if c.BorderTop > 0 {
		render.Line(c.BorderTopColor, c.BorderTopStyle, c.BorderTop, border.X, border.Y+c.Top, border.X+border.W, border.Y+c.Top)
	}
	if c.BorderRight > 0 {
		render.Line(c.BorderRightColor, c.BorderRightStyle, c.BorderRight, border.X+border.W, border.Y+c.Top, border.X+border.W, border.Y+c.Top+border.H)
	}
	if c.BorderBottom > 0 {
		render.Line(c.BorderBottomColor, c.BorderBottomStyle, c.BorderBottom, border.X, border.Y+c.Top+border.H, border.X+border.W, border.Y+c.Top+border.H)
	}
	if c.BorderLeft > 0 {
		render.Line(c.BorderLeftColor, c.BorderLeftStyle, c.BorderLeft, border.X, border.Y+c.Top, border.X, border.Y+c.Top+border.H)
	}
	return Rect{X: from_x, Y: from_y, W: float64(c.Width), H: c.Height}, nil
}
func renderContents(c *Block, render Render, contentRect Rect) error {
	line_count := len(c.Contents)
	total_height := float64(line_count * c.LineHeight)
	if total_height > contentRect.H {
		total_height = contentRect.H
	}
	startY := contentRect.Y + c.Top
	if c.VAlign == VAlignMiddle {
		startY = contentRect.Y + contentRect.H/2 - total_height/2 + c.Top
	} else if c.VAlign == VAlignBottom {
		startY = contentRect.Y + contentRect.H - total_height + c.Top
	}
	for i, txt := range c.Contents {
		if float64(c.LineHeight*(i+1)) > contentRect.H {
			break
		}
		startX := contentRect.X
		if c.Align == AlignRight {
			startX = contentRect.X + contentRect.W - c.contentsWidth[i]
		} else if c.Align == AlignCenter {
			startX = contentRect.X + contentRect.W/2 - c.contentsWidth[i]/2
		}
		rect := Rect{X: startX, Y: startY + float64(c.LineHeight-c.FontSize)/2, W: contentRect.W, H: contentRect.H}
		err := render.Text(&rect, txt, c.FontFamily, c.FontStyle, c.FontSize, c.FontColor)
		if err != nil {
			return err
		}
		startY += float64(c.LineHeight)
	}
	return nil
}

func renderColumns(c *Block, render Render, content Rect) error {
	content_x, content_y := content.X, content.Y
	for _, cell := range c.Children {
		offsetY := 0.0
		if c.VAlign == VAlignMiddle {
			offsetY = (content.H - cell.Height) / 2
		} else if c.VAlign == VAlignBottom {
			offsetY = (content.H - cell.Height)
		}
		rect, err := RenderTo(cell, render, content_x, content_y+offsetY)
		if err != nil {
			return err
		}
		content_x = rect.X + rect.W
		content_y = rect.Y
	}
	return nil
}

func renderRows(c *Block, render Render, content Rect) error {
	total_height := float64(0)
	for _, row := range c.Children {
		total_height += row.Height
		if total_height > content.H {
			break
		}
		offsetY := 0.0
		if c.VAlign == VAlignMiddle {
			offsetY = (content.H - row.Height) / 2
		} else if c.VAlign == VAlignBottom {
			offsetY = (content.H - row.Height)
		}
		//do not support page on none root tag
		rect, err := RenderTo(row, render, content.X, content.Y+offsetY)
		if err != nil {
			return err
		}
		content.X = rect.X
		content.Y = rect.Y + rect.H
	}
	return nil
}
