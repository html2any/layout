package layout

import (
	"fmt"
	"image/color"
	"log/slog"
	"strconv"
	"strings"
)

const (
	None                float64 = 0
	AlignLeft           string  = "left"
	AlignRight          string  = "right"
	AlignCenter         string  = "center"
	VAlignTop           string  = "top"
	VAlignMiddle        string  = "middle"
	VAlignBottom        string  = "bottom"
	BorderStyleSolid    string  = "solid"
	BorderStyleDotted   string  = "dotted"
	BorderStyleNone     string  = "none"
	OverFlowHidden      string  = "hidden"
	OverFlowExtend      string  = "auto"
	FlexDirectionRow    string  = "row"
	FlexDirectionColumn string  = "column"
)

type Style struct {
	MarginTop, MarginRight, MarginBottom, MarginLeft                     float64    `json:"-"`
	PaddingTop, PaddingRight, PaddingBottom, PaddingLeft                 float64    `json:"-"`
	Align, VAlign                                                        string     `json:"-"`
	BackgroundColor                                                      color.RGBA `json:"-"`
	BackgroundImage                                                      string     `json:"-"`
	FontFamily                                                           string     `json:"-"`
	FontStyle                                                            string     `json:"-"`
	FontSize                                                             int        `json:"-"`
	Height                                                               float64    `json:"-"`
	OriginHeight                                                         float64    `json:"-"`
	LineHeight                                                           int        `json:"-"`
	FontColor                                                            color.RGBA `json:"-"`
	Width                                                                float64    `json:"-"`
	BorderTop, BorderRight, BorderBottom, BorderLeft                     float64    `json:"-"`
	BorderTopColor, BorderRightColor, BorderBottomColor, BorderLeftColor color.RGBA `json:"-"`
	BorderTopStyle, BorderRightStyle, BorderBottomStyle, BorderLeftStyle string     `json:"-"`
	Flex                                                                 int        `json:"-"`
	FlexDirection                                                        string     `json:"-"`
	OverFlow                                                             string     `json:"-"`
	Top                                                                  float64    `json:"-"`
}

func NewStyle() *Style {
	s := &Style{
		Align:             AlignLeft,
		VAlign:            VAlignTop,
		BackgroundColor:   color.RGBA{255, 255, 255, 0},
		FontFamily:        "times",
		FontStyle:         "normal",
		FontSize:          14,
		FontColor:         color.RGBA{0, 0, 0, 0},
		LineHeight:        0,
		Height:            0,
		Width:             0,
		BorderTop:         None,
		BorderRight:       None,
		BorderBottom:      None,
		BorderLeft:        None,
		BorderTopColor:    color.RGBA{0, 0, 0, 0},
		BorderRightColor:  color.RGBA{0, 0, 0, 0},
		BorderBottomColor: color.RGBA{0, 0, 0, 0},
		BorderLeftColor:   color.RGBA{0, 0, 0, 0},
		BorderTopStyle:    BorderStyleNone,
		BorderRightStyle:  BorderStyleNone,
		BorderBottomStyle: BorderStyleNone,
		BorderLeftStyle:   BorderStyleNone,
		Flex:              1,
		FlexDirection:     FlexDirectionColumn,
		OverFlow:          OverFlowExtend,
		Top:               0,
	}
	return s
}

func (s *Style) inherit(style *Style) {
	//auto inherit width when flex direction is not row
	if style.FlexDirection == FlexDirectionColumn {
		s.Width = style.Width
	}
	s.Align = style.Align
	s.VAlign = style.VAlign
	s.FontSize = style.FontSize
	s.FontFamily = style.FontFamily
	s.FontStyle = style.FontStyle
	s.FontColor = style.FontColor
	s.LineHeight = style.LineHeight
	s.Top = style.Top
	if style.FlexDirection == FlexDirectionRow {
		s.PaddingBottom = style.PaddingBottom
		s.PaddingLeft = style.PaddingLeft
		s.PaddingRight = style.PaddingRight
		s.PaddingTop = style.PaddingTop
	}
}

func (s *Style) parseWidthStyleColor(input string) (float64, string, color.RGBA) {
	colr := color.RGBA{0, 0, 0, 0}
	if input == "" {
		return 0, "", colr
	}
	width := float64(0)
	style := "solid"

	arr := strings.Split(input, " ")
	for _, item := range arr {
		if item == BorderStyleSolid || item == BorderStyleDotted {
			style = item
		} else if strings.HasPrefix(item, "#") {
			colr = MustParseHexColor(item)
		} else if strings.Contains(item, "px") {
			width, _ = strconv.ParseFloat(strings.TrimSuffix(item, "px"), 64)
		} else {
			width, _ = strconv.ParseFloat(item, 64)
		}
	}
	return width, style, colr
}
func (s *Style) borderValue(key, value string) string {
	if value != "solid" && value != "dotted" && value != "dashed" && value != "none" {
		slog.Warn("invalid value for " + key + ": " + value)
		return "solid"
	}
	return value
}

func (s *Style) directionValue(key, value string) string {
	if value != "row" && value != "column" {
		slog.Warn("invalid value for " + key + ": " + value)
		return "column"
	}
	return value
}
func (s *Style) fontweightValue(key, value string) string {
	if value != "normal" && value != "bold" && value != "italic" {
		slog.Warn("invalid value for " + key + ": " + value)
		return "normal"
	}
	return value
}

func (s *Style) overFlowValue(key, value string) string {
	if value != OverFlowHidden && value != OverFlowExtend {
		slog.Warn("invalid value for " + key + ": " + value)
		return OverFlowExtend
	}
	return value
}

// if value != "left" && value != "center" && value != "right" {
func (s *Style) textAlignValue(key, value string) string {
	if value != "left" && value != "center" && value != "right" {
		slog.Warn("invalid value for " + key + ": " + value)
		return "left"
	}
	return value
}

// if value != "top" && value != "middle" && value != "bottom" {
func (s *Style) verticalAlignValue(key, value string) string {
	if value != "top" && value != "middle" && value != "bottom" {
		slog.Warn("invalid value for " + key + ": " + value)
		return "top"
	}
	return value
}

//	if strings.HasPrefix(value, "url(") && strings.HasSuffix(value, ")") {
//		s.BackgroundImage = value[4 : len(value)-1]
func (s *Style) parseImageStyle(key, value string) string {
	if strings.HasPrefix(value, "url(") && strings.HasSuffix(value, ")") {
		return value[4 : len(value)-1]
	} else {
		slog.Warn("invalid value for " + key + ": " + value)
	}
	return value
}

func (s *Style) ParseStyle(style string, last_style *Style) error {
	if last_style != nil {
		s.inherit(last_style)
	}
	style = strings.Join(strings.Split(style, "\n"), "")
	for _, item := range strings.Split(style, ";") {
		arr := strings.SplitN(item, ":", 2)
		if len(arr) != 2 {
			continue
		}
		key := strings.TrimSpace(arr[0])
		value := strings.TrimSpace(arr[1])
		switch key {
		case "a", "align", "text-align", "align-items":
			s.Align = s.textAlignValue(key, value)
		case "va", "vertical-align", "vertical-align-items":
			s.VAlign = s.verticalAlignValue(key, value)
		case "bg", "background-color":
			s.BackgroundColor = MustParseHexColor(value)
		case "bi", "background-image":
			s.BackgroundImage = s.parseImageStyle(key, value)
		case "fs", "font-size":
			s.FontSize = ToInt(value)
		case "ff", "font-family":
			s.FontFamily = value
		case "fsl", "font-style", "font-weight":
			s.FontStyle = s.fontweightValue(key, value)
		case "fc", "font-color", "color":
			s.FontColor = MustParseHexColor(value)
		case "ll", "line-height":
			s.LineHeight = ToInt(value)
		case "w", "width":
			s.Width = ToFloat64(value)
		case "h", "height":
			s.Height = ToFloat64(value)
			s.OriginHeight = s.Height
		case "mt", "margin-top":
			s.MarginTop = ToFloat64(value)
		case "mr", "margin-right":
			s.MarginRight = ToFloat64(value)
		case "mb", "margin-bottom":
			s.MarginBottom = ToFloat64(value)
		case "ml", "margin-left":
			s.MarginLeft = ToFloat64(value)
		case "pt", "padding-top":
			s.PaddingTop = ToFloat64(value)
		case "pr", "padding-right":
			s.PaddingRight = ToFloat64(value)
		case "pb", "padding-bottom":
			s.PaddingBottom = ToFloat64(value)
		case "pl", "padding-left":
			s.PaddingLeft = ToFloat64(value)

		case "bt", "border-top":
			s.BorderTop, s.BorderTopStyle, s.BorderTopColor = s.parseWidthStyleColor(value)
		case "br", "border-right":
			s.BorderRight, s.BorderRightStyle, s.BorderRightColor = s.parseWidthStyleColor(value)
		case "bb", "border-bottom":
			s.BorderBottom, s.BorderBottomStyle, s.BorderBottomColor = s.parseWidthStyleColor(value)
		case "bl", "border-left":
			s.BorderLeft, s.BorderLeftStyle, s.BorderLeftColor = s.parseWidthStyleColor(value)
		case "btc", "border-top-color":
			s.BorderTopColor = MustParseHexColor(value)
		case "brc", "border-right-color":
			s.BorderRightColor = MustParseHexColor(value)
		case "bbc", "border-bottom-color":
			s.BorderBottomColor = MustParseHexColor(value)
		case "blc", "border-left-color":
			s.BorderLeftColor = MustParseHexColor(value)
		case "bts", "border-top-style":
			s.BorderTopStyle = s.borderValue(key, value)
		case "brs", "border-right-style":
			s.BorderRightStyle = s.borderValue(key, value)
		case "bbs", "border-bottom-style":
			s.BorderBottomStyle = s.borderValue(key, value)
		case "bls", "border-left-style":
			s.BorderLeftStyle = s.borderValue(key, value)
		case "m", "margin":
			vals := strings.Split(value, " ")
			if len(vals) == 4 {
				s.MarginTop = ToFloat64(vals[0])
				s.MarginRight = ToFloat64(vals[1])
				s.MarginBottom = ToFloat64(vals[2])
				s.MarginLeft = ToFloat64(vals[3])
			} else if len(vals) == 2 {
				s.MarginTop = ToFloat64(vals[0])
				s.MarginRight = ToFloat64(vals[1])
				s.MarginBottom = ToFloat64(vals[0])
				s.MarginLeft = ToFloat64(vals[1])
			} else {
				s.MarginTop = ToFloat64(value)
				s.MarginRight = ToFloat64(value)
				s.MarginBottom = ToFloat64(value)
				s.MarginLeft = ToFloat64(value)
			}
		case "p", "padding":
			vals := strings.Split(value, " ")
			if len(vals) == 4 {
				s.PaddingTop = ToFloat64(vals[0])
				s.PaddingRight = ToFloat64(vals[1])
				s.PaddingBottom = ToFloat64(vals[2])
				s.PaddingLeft = ToFloat64(vals[3])
			} else if len(vals) == 2 {
				s.PaddingTop = ToFloat64(vals[0])
				s.PaddingRight = ToFloat64(vals[1])
				s.PaddingBottom = ToFloat64(vals[0])
				s.PaddingLeft = ToFloat64(vals[1])
			} else {
				s.PaddingTop = ToFloat64(value)
				s.PaddingRight = ToFloat64(value)
				s.PaddingBottom = ToFloat64(value)
				s.PaddingLeft = ToFloat64(value)
			}
		case "b", "border":
			width, styl, colr := s.parseWidthStyleColor(value)
			s.BorderTop, s.BorderRight, s.BorderBottom, s.BorderLeft = width, width, width, width
			s.BorderTopStyle, s.BorderRightStyle, s.BorderBottomStyle, s.BorderLeftStyle = styl, styl, styl, styl
			s.BorderTopColor, s.BorderRightColor, s.BorderBottomColor, s.BorderLeftColor = colr, colr, colr, colr
		case "bs", "border-style":
			vals := strings.Split(value, " ")
			if len(vals) == 4 {
				s.BorderTopStyle = vals[0]
				s.BorderRightStyle = vals[1]
				s.BorderBottomStyle = vals[2]
				s.BorderLeftStyle = vals[3]
			} else if len(vals) == 2 {
				s.BorderTopStyle = vals[0]
				s.BorderRightStyle = vals[1]
				s.BorderBottomStyle = vals[0]
				s.BorderLeftStyle = vals[1]
			} else {
				s.BorderTopStyle = value
				s.BorderRightStyle = value
				s.BorderBottomStyle = value
				s.BorderLeftStyle = value
			}
		case "bc", "border-color":
			vals := strings.Split(value, " ")
			if len(vals) == 4 {
				s.BorderTopColor = MustParseHexColor(vals[0])
				s.BorderRightColor = MustParseHexColor(vals[1])
				s.BorderBottomColor = MustParseHexColor(vals[2])
				s.BorderLeftColor = MustParseHexColor(vals[3])
			} else if len(vals) == 2 {
				s.BorderTopColor = MustParseHexColor(vals[0])
				s.BorderRightColor = MustParseHexColor(vals[1])
				s.BorderBottomColor = MustParseHexColor(vals[0])
				s.BorderLeftColor = MustParseHexColor(vals[1])
			} else {
				s.BorderTopColor = MustParseHexColor(value)
				s.BorderRightColor = MustParseHexColor(value)
				s.BorderBottomColor = MustParseHexColor(value)
				s.BorderLeftColor = MustParseHexColor(value)
			}
		case "f", "flex":
			s.Flex = ToInt(value)
		case "fd", "flex-direction":
			s.FlexDirection = s.directionValue(key, value)
		case "of", "overflow":
			s.OverFlow = s.overFlowValue(key, value)
		case "t", "top":
			s.Top = ToFloat64(value)
		default:
			slog.Debug("ignore css," + key + ":" + value)
		}
	}
	//auto set line height to fontsize height
	if s.LineHeight == 0 {
		s.LineHeight = s.FontSize
	}
	return nil
}

func (s *Style) cacuRect(from_x, from_y float64) (border_rect Rect, back_rect Rect, cont_rect Rect) {
	border_x := from_x + s.MarginLeft
	border_y := from_y + s.MarginTop
	border_w := float64(s.Width) - s.MarginLeft - s.MarginRight
	border_h := float64(s.Height) - s.MarginTop - s.MarginBottom
	border_rect = Rect{X: border_x, Y: border_y, W: border_w, H: border_h}

	back_x := border_x + s.BorderLeft
	back_y := border_y + s.BorderTop
	back_w := border_w - s.BorderLeft - s.BorderRight
	back_h := border_h - s.BorderTop - s.BorderBottom
	back_rect = Rect{X: back_x, Y: back_y, W: back_w, H: back_h}

	cont_x := back_x + s.PaddingLeft
	cont_y := back_y + s.PaddingTop
	cont_w := back_w - s.PaddingLeft - s.PaddingRight
	cont_h := back_h - s.PaddingTop - s.PaddingBottom
	cont_rect = Rect{X: cont_x, Y: cont_y, W: cont_w, H: cont_h}

	return
}

func (s Style) String() string {
	return fmt.Sprintf(
		    "top:%.1f;"+
			"margin-top:%.1f;"+
			"margin-right:%.1f;"+
			"margin-bottom:%.1f;"+
			"margin-left:%.1f;"+
			"padding-top:%.1f;"+
			"padding-right:%.1f;"+
			"padding-bottom:%.1f;"+
			"padding-left:%.1f;"+
			"align:%s;"+
			"v-align:%s;"+
			"background-color:%s;"+
			"font-family:%s;"+
			"font-style:%s;"+
			"font-size:%d;"+
			"height:%.1f;"+
			"line-height:%d;"+
			"font-color:%s;"+
			"width:%.1f;"+
			"border-top:%.1f;"+
			"border-right:%.1f;"+
			"border-bottom:%.1f;"+
			"border-left:%.1f;"+
			"border-top-color:%s;"+
			"border-right-color:%s;"+
			"border-bottom-color:%s;"+
			"border-left-color:%s;"+
			"border-top-style:%s;"+
			"border-right-style:%s;"+
			"border-bottom-style:%s;"+
			"border-left-style:%s;"+
			"flex:%d;"+
			"flex-direction:%s;"+
			"overflow:%s;",
		s.Top,
		s.MarginTop,
		s.MarginRight,
		s.MarginBottom,
		s.MarginLeft,
		s.PaddingTop,
		s.PaddingRight,
		s.PaddingBottom,
		s.PaddingLeft,
		s.Align,
		s.VAlign,
		RGBA2Hex(s.BackgroundColor),
		s.FontFamily,
		s.FontStyle,
		s.FontSize,
		s.Height,
		s.LineHeight,
		RGBA2Hex(s.FontColor),
		s.Width,
		s.BorderTop,
		s.BorderRight,
		s.BorderBottom,
		s.BorderLeft,
		RGBA2Hex(s.BorderTopColor),
		RGBA2Hex(s.BorderRightColor),
		RGBA2Hex(s.BorderBottomColor),
		RGBA2Hex(s.BorderLeftColor),
		s.BorderTopStyle,
		s.BorderRightStyle,
		s.BorderBottomStyle,
		s.BorderLeftStyle,
		s.Flex,
		s.FlexDirection,
		s.OverFlow)
}
