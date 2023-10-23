package layout

import (
	"strings"

	"github.com/html2any/parser"
)

type block_global struct {
	cssMap *map[string]string
	idMap  map[string]*Block
	// render Render
}

type Block struct {
	Style
	Tag      string   `json:"tag"`
	Class    string   `json:"style"`
	Src      string   `json:"src"`
	Children []*Block `json:"children"`
	Contents []string `json:"contents"`

	// Path          string        `json:"-"`
	g             *block_global `json:"-"`
	parent        *Block        `json:"-"`
	contentsWidth []float64     `json:"-"`
	Id            string        `json:"-"`
}

func NewBlock(cssMap *map[string]string) *Block {
	root := &Block{}
	g := &block_global{
		cssMap: cssMap,
		idMap:  make(map[string]*Block),
	}
	root.g = g
	root.Style = *NewStyle()
	return root
}

func (t *Block) SetTagName(tagname string) parser.IHtmlTag {
	t.Tag = tagname
	t.Children = make([]*Block, 0)

	if t.parent != nil && len(t.parent.Id) > 0 {
		t.Id = t.parent.Id
	}

	if t.g.cssMap != nil {
		if style, ok := (*t.g.cssMap)[tagname]; ok {
			t.Class = t.Class + ";" + style
		}
		if len(t.Id) > 0 { //support #id>td style find
			if style, ok := (*t.g.cssMap)["#"+t.Id+">"+tagname]; ok {
				t.Class = t.Class + ";" + style
			}
		}
	}
	return t
}

// GetTagName() string
func (t *Block) GetTagName() string {
	return t.Tag
}

// SetContent(content string)
func (t *Block) SetContent(content string) parser.IHtmlTag {
	t.Contents = strings.Split(content, "\n")
	return t
}

// AddChild(child *IHtmlTree)
func (t *Block) AddChild(child parser.IHtmlTag) parser.IHtmlTag {
	t.Children = append(t.Children, child.(*Block))
	return t
}

func (t *Block) SetAttr(attr string, value string) parser.IHtmlTag {
	if attr == "style" {
		t.Class = t.Class + ";" + strings.TrimSpace(value)
	} else if attr == "src" {
		t.Src = strings.TrimSpace(value)
	} else if attr == "id" {
		t.g.idMap[value] = t
		t.Id = value
		if t.g.cssMap != nil {
			if style, ok := (*t.g.cssMap)["#"+value]; ok {
				t.Class = t.Class + ";" + style
			}
		}
	} else if attr == "class" {
		if t.g.cssMap != nil {
			for _, class := range strings.Split(value, " ") {
				if style, ok := (*t.g.cssMap)["."+class]; ok {
					t.Class = t.Class + ";" + style
				}
			}
		}
	}
	return t
}

func (t *Block) NewTag() parser.IHtmlTag {
	ntag := new(Block)
	ntag.g = t.g
	ntag.parent = t
	ntag.Style = *NewStyle()
	return ntag
}
