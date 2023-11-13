package layout

import (
	"errors"
	"image/color"
	"log/slog"
)

const (
	A4_Width  = 595.276
	A4_Height = 841.89
)

func SplitPages(c *Block, page_height float64) (pages []*Block, err error) {
	fromMap := make(map[*Block]int)
	pages = make([]*Block, 0)
	for {
		if ret, left, err := page(c, &fromMap, 0, page_height); err != nil {
			return nil, err
		} else {
			pages = append(pages, ret)
			if left == ALL_PAGED {
				break
			}
		}
	}
	return
}

func CacuPages(c *Block, width float64, gFamily string, ilayout ILayout) (height float64, pages []*Block, err error) {
	if _, err := CacuHeight(c, width, gFamily, ilayout); err != nil {
		return 0, nil, err
	}
	if c.OriginHeight <= 0 {
		slog.Warn("block height is ", c.OriginHeight, ",Should Set Height for Page. Now Going to flow mode")
		return c.Height, []*Block{c}, nil
	}
	pages, err = SplitPages(c, c.OriginHeight)
	return
}

func CacuHeight(c *Block, width float64, gFamily string, ilayout ILayout) (float64, error) {
	c.Width = width
	c.FontFamily = gFamily
	c.FontColor = color.RGBA{0, 0, 0, 255}
	c.ParseStyle(c.Class, nil)
	if err := cacuHeight(c, ilayout); err != nil {
		return 0, err
	} else {
		return c.Height, nil
	}
}

func cacuHeight(c *Block, ilayout ILayout) error {
	if ilayout.Overide(c) {
		return nil
	}
	for _, child := range c.Children {
		child.ParseStyle(child.Class, &c.Style)
		// child.Path = c.Path + "/" + strconv.Itoa(i)
	}
	updateChildrenWidth(c)
	if len(c.Contents) > 0 { //leaf text node
		if content_height, err := cacuContentsHeight(c, ilayout); err != nil {
			return err
		} else {
			c.Height = content_height
		}
	} else if len(c.Children) == 0 { //leaf node
		updateCurrentAndChildrenHeight(c)
	}
	//children start
	for _, child := range c.Children {
		if err := cacuHeight(child, ilayout); err != nil {
			return err
		}
	}
	updateCurrentAndChildrenHeight(c)
	return nil
}

func updateChildrenWidth(c *Block) error {
	emptyWidth := c.PaddingLeft + c.PaddingRight
	emptyWidth += c.BorderLeft + c.BorderRight
	emptyWidth += c.MarginLeft + c.MarginRight
	useable_width := c.Width - emptyWidth
	//Vertical Layout
	if c.FlexDirection == FlexDirectionColumn {
		for _, child := range c.Children {
			child.Width = useable_width
		}
		return nil
	} else if c.FlexDirection == FlexDirectionRow {
		//Horizontal Layout
		width_arr := make([]float64, len(c.Children))
		flex_width_map := make(map[int]int, 0)
		flex_width_arr := make([]int, 0)
		width_left := useable_width
		//auto set column width: (total_width-sum(column_with_width))/[flex_size]
		for i, child := range c.Children {
			if child.Width <= 0 {
				flex_width_arr = append(flex_width_arr, int(child.Flex))
				flex_width_map[i] = len(flex_width_arr) - 1
			} else {
				width_arr[i] = child.Width
				width_left -= child.Width
			}
		}
		var gcd_width []int
		var cols int
		if len(flex_width_arr) > 0 {
			gcd_width = GcdInts(flex_width_arr)
			cols = SumInt(gcd_width)
		}
		//update children width
		for i, child := range c.Children {
			if w, ok := flex_width_map[i]; ok {
				child.Width = width_left / float64(cols) * float64(gcd_width[w])
			} else {
				child.Width = width_arr[i]
			}
		}
	} else {
		return errors.New("unknown flex direction:" + c.FlexDirection)
	}
	return nil
}

func cacuContentsHeight(c *Block, ilayout ILayout) (height float64, err error) {
	contentWidth := c.Width - c.PaddingLeft - c.PaddingRight - c.BorderLeft - c.BorderRight - c.MarginLeft - c.MarginRight
	parsed_content := make([]string, 0)
	parsed_content_width := make([]float64, 0)
	for _, line := range c.Contents {
		if lines, lineWidth, err := ilayout.SplitLines(contentWidth, line, c.FontFamily, float64(c.FontSize), c.FontStyle); err != nil {
			return 0, err
		} else {
			parsed_content = append(parsed_content, lines...)
			parsed_content_width = append(parsed_content_width, lineWidth...)
		}
	}
	c.Contents = parsed_content
	c.contentsWidth = parsed_content_width
	line_count := len(parsed_content)
	if line_count > 0 {
		return float64(line_count * c.LineHeight), nil
	}
	return 0, nil
}

func updateCurrentAndChildrenHeight(c *Block) {
	emptyHeight := c.PaddingTop + c.PaddingBottom
	emptyHeight += c.BorderTop + c.BorderBottom
	emptyHeight += c.MarginTop + c.MarginBottom
	useable_height := c.Height - emptyHeight

	height := float64(0)
	if c.FlexDirection == FlexDirectionColumn {
		for _, child := range c.Children {
			height = height + child.Height
		}
	} else if c.FlexDirection == FlexDirectionRow {
		for _, child := range c.Children {
			if child.Height > height {
				height = child.Height
			}
		}
	}

	if c.OverFlow == OverFlowExtend {
		if height < useable_height {
			height = useable_height
		}
	} else if c.OverFlow == OverFlowHidden {
		//if not set height, when set overflow:hidden, the useable_height auto set to height
		if useable_height < 0 {
			useable_height = height
		}
		if height > useable_height {
			height = useable_height
		}
	} else {
		panic("unknown overflow")
	}
	if c.FlexDirection == FlexDirectionRow {
		for _, child := range c.Children {
			child.Height = height
		}
	}
	c.Height = height + emptyHeight
}

const ALL_PAGED = 1<<31 - 1 // 2147483647

func page(c *Block, fromMap *map[*Block]int, deep int, left float64) (out *Block, rleft float64, err error) {
	from, visited := (*fromMap)[c]
	if !visited {
		if deep == 0 && c.Height < left {
			return c, ALL_PAGED, nil
		}
		//ignore step in if the height is enough
		if c.Height < left {
			return c, left - c.Height, nil
		}
		(*fromMap)[c] = 0
		from = 0
	}

	emptyTop := c.MarginTop + c.PaddingTop + c.BorderTop
	emptyBottom := c.MarginBottom + c.PaddingBottom + c.BorderBottom
	useable_left := left - emptyTop - emptyBottom
	i := from
	if len(c.Contents) > 0 {
		contents := make([]string, 0)
		height := 0.0
		maxLen := len(c.Contents)
		for ; i < maxLen; i++ {
			if useable_left-height-float64(c.LineHeight) > 0 {
				contents = append(contents, c.Contents[i])
				height += float64(c.LineHeight)
			} else {
				break
			}
		}
		(*fromMap)[c] = i
		newBlock := *c
		newBlock.Contents = make([]string, 0)
		if height > 0 {
			newBlock.Height = height + emptyBottom + emptyTop
			newBlock.Contents = append(newBlock.Contents, contents...)
			if i == maxLen { //hold all left lines
				useable_left = useable_left - height
			} else { //can not hold all left lines, but has one or more lines
				useable_left = 0
			}
		} else { //can not hold one line
			newBlock.Height = height + emptyBottom + emptyTop
			useable_left = -1
		}
		return &newBlock, useable_left, nil
	} else if len(c.Children) == 0 {
		useable_left = useable_left - c.Height
		newBlock := *c
		newBlock.Height = emptyBottom + emptyTop
		return &newBlock, useable_left, nil
	} else if c.FlexDirection == FlexDirectionRow {
		children_len := len(c.Children)
		children := make([]*Block, 0, children_len)
		min_left := useable_left
		all_none := true
		max_height := 0.0
		//find the min left, and the max height for row layout
		for i := 0; i < len(c.Children); i++ {
			child := c.Children[i]
			newchild, rleft, err := page(child, fromMap, deep+1, useable_left)
			if err != nil {
				return nil, 0, err
			}
			if rleft >= 0 {
				all_none = false
			}
			if rleft < min_left {
				min_left = rleft
			}
			if newchild.Height > max_height {
				max_height = newchild.Height
			}
			children = append(children, newchild)
		}
		//reset children height
		for i, _ := range children {
			children[i].Height = max_height
		}

		useable_left = min_left
		newBlock := *c
		newBlock.Height = max_height + emptyTop + emptyBottom
		newBlock.Children = make([]*Block, 0, len(children))
		//remove empty children
		if !all_none {
			newBlock.Children = append(newBlock.Children, children...)
		}
		return &newBlock, useable_left, nil
	} else if c.FlexDirection == FlexDirectionColumn {
		children_len := len(c.Children)
		children := make([]*Block, 0, children_len)
		total_height := 0.0
		has_break := false
		//match column layout children, and append to new block
		for ; i < children_len; i++ {
			child := c.Children[i]
			(*fromMap)[c] = i
			newchild, rleft, err := page(child, fromMap, deep+1, useable_left)
			if err != nil {
				return nil, 0, err
			}
			total_height += newchild.Height
			useable_left = rleft
			if useable_left <= 0 {
				children = append(children, newchild)
				has_break = true
				break
			} else {
				children = append(children, newchild)
			}
		}
		newBlock := *c
		newBlock.Height = total_height + emptyTop + emptyBottom
		newBlock.Children = make([]*Block, 0, len(children))
		newBlock.Children = append(newBlock.Children, children...)
		if deep == 0 && i == children_len && !has_break {
			return &newBlock, ALL_PAGED, nil
		}
		return &newBlock, useable_left, nil
	} else {
		return nil, 0, errors.New("invalid flex direction")
	}
}
