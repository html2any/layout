package fonts

type Int26_6 int32

type TextSpliter interface {
	SplitLines(width float64, s string, family string, size float64, style string) (lines []string, line_widths []float64, err error)
}

// SplitLines: split string to lines with specific width
// func (f *Faces) SplitLines(width float64, s string, family string, size float64, style string) (lines []string, line_widths []float64, err error) {
// 	face, err := f.GetFace(family, size, FontStyle(style))
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	// var cache map[rune]Int26_6
// 	// if f.wordCache != nil {
// 	// 	cache = f.wordCache[face]
// 	// }

// 	prevC := rune(-1)
// 	advance := Int26_6(0)
// 	last_i := 0
// 	srune := []rune(s)
// 	rune_len := len(srune)
// 	i := 0
// 	for ; i < rune_len; i++ {
// 		c := srune[i]
// 		space := Int26_6(0)
// 		//Check space between two word
// 		if prevC >= 0 {
// 			space = Int26_6(face.Kern(prevC, c))
// 		}
// 		//Get word width
// 		var i26_6_a Int26_6
// 		// if cache != nil {
// 		// 	if ac, ok := cache[c]; ok {
// 		// 		i26_6_a = ac
// 		// 	} else {
// 		// 		a, _ := face.GlyphAdvance(c)
// 		// 		i26_6_a = Int26_6(a)
// 		// 		cache[c] = i26_6_a
// 		// 	}
// 		// } else {
// 		a, _ := face.GlyphAdvance(c)
// 		i26_6_a = Int26_6(a)
// 		// }

// 		//Check total width > limit width
// 		testw := advance + space + i26_6_a
// 		if !(float64(testw>>6) > width) { //continue
// 			advance += i26_6_a
// 		} else { //If not enough space, add line
// 			lines = append(lines, string(srune[last_i:i]))
// 			line_widths = append(line_widths, float64(advance>>6))
// 			last_i = i
// 			advance = space + i26_6_a
// 		}
// 		prevC = c
// 	}
// 	if last_i != i {
// 		lines = append(lines, string(srune[last_i:i]))
// 		line_widths = append(line_widths, float64(advance>>6))
// 	}
// 	return
// }

// func (f *Faces) MeasureString(s string, family string, size float64, style FontStyle) (float64, error) {
// 	face, err := f.GetFace(family, size, style)
// 	if err != nil {
// 		return 0, err
// 	}
// 	width := font.MeasureString(face, s)
// 	return float64(width >> 6), nil
// }

func (f *SFNT) SplitLines(width float64, s string, family string, size float64, style string) (lines []string, line_widths []float64, err error) {
	stl := Regular
	if style == "bold" {
		stl = Bold
	} else if style == "italic" {
		stl = Italic
	}

	sfnt, err := f.GetSFNT(family, size, stl)
	if err != nil {
		return nil, nil, err
	}
	nwidth := width * 1000 / size
	advance := 0
	prevC := uint16(0)
	last_i := 0
	srune := []rune(s)
	rune_len := len(srune)
	i := 0
	for ; i < rune_len; i++ {
		c := sfnt.GlyphIndex(srune[i])
		space := int16(0)
		//Check space between two word
		if prevC > 0 {
			space = sfnt.Kerning(prevC, c)
		}
		//Get word width
		a := int(int16(sfnt.GlyphAdvance(c)) + space)

		//Check total width > limit width
		testw := advance + a
		if !(float64(testw) > nwidth) { //continue
			advance += a
		} else { //If not enough space, add line
			lines = append(lines, string(srune[last_i:i]))
			line_widths = append(line_widths, float64(advance)*size/1000)
			last_i = i
			advance = a
		}
		prevC = c
	}
	if last_i != i {
		lines = append(lines, string(srune[last_i:i]))
		line_widths = append(line_widths, float64(advance)*size/1000)
	}
	return
}

func (f *SFNT) MeasureString(s string, family string, size float64, style FontStyle) (float64, error) {
	sfnt, err := f.GetSFNT(family, size, style)
	if err != nil {
		return 0, err
	}
	advance := 0
	prevC := uint16(0)
	for _, c := range s {
		glyphIndex := sfnt.GlyphIndex(c)
		if prevC > 0 {
			advance += int(sfnt.Kerning(prevC, glyphIndex))
		}
		a := sfnt.GlyphAdvance(glyphIndex)
		advance += int(a)
		prevC = glyphIndex
	}
	return float64(advance) * size / 1000, nil
}
