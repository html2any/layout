package fonts

import (
	"errors"
	"os"
)

type FontStyle string

var (
	Regular FontStyle = ""
	Bold    FontStyle = "bold"
	Italic  FontStyle = "italic"
)

type Family struct {
	Name  string
	Style FontStyle
}

func (f *Family) String() string {
	return "Name:" + f.Name + ",Style:" + string(f.Style)
}

type Fonts struct {
	dataCache map[Family][]byte
}

func NewFonts() *Fonts {
	f := &Fonts{}
	f.dataCache = make(map[Family][]byte)
	return f
}

func (f *Fonts) AddFont(family, path string, style FontStyle) error {
	k := Family{
		Name:  family,
		Style: style,
	}
	if _, ok := f.dataCache[k]; ok {
		return nil
	}
	if data, err := os.ReadFile(path); err != nil {
		return err
	} else {
		f.dataCache[k] = data
	}
	return nil
}

func (f *Fonts) GetFonts() map[Family][]byte {
	return f.dataCache
}

func (f *Fonts) GetFont(family string, style FontStyle) ([]byte, error) {
	k := Family{
		Name:  family,
		Style: style,
	}
	if data, ok := f.dataCache[k]; ok {
		return data, nil
	} else {
		return nil, errors.New("please add font before get data")
	}
}
