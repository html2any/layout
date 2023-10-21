package fonts

import (
	"errors"

	"github.com/tdewolff/canvas/font"
)

type SFNT struct {
	fonts     *Fonts
	faceCache map[Family]*font.SFNT
}

func NewSFNT(fonts *Fonts) *SFNT {
	f := &SFNT{}
	f.faceCache = make(map[Family]*font.SFNT)
	f.fonts = fonts
	return f
}

func (f *SFNT) CacheAll() (*SFNT, error) {
	for family, data := range f.fonts.dataCache {
		if sfnt, err := font.ParseSFNT(data, 0); err != nil {
			return nil, err
		} else {
			f.faceCache[family] = sfnt
		}
	}
	return f, nil
}

func (f *SFNT) GetSFNT(family string, size float64, style FontStyle) (*font.SFNT, error) {
	key := Family{
		Name:  family,
		Style: style,
	}
	if sfnt, ok := f.faceCache[key]; ok {
		return sfnt, nil
	} else if data, ok := f.fonts.dataCache[key]; !ok {
		return nil, errors.New("font not found(" + key.String() + "). please add font before get face")
	} else if sfnt, err := font.ParseSFNT(data, 0); err != nil {
		return nil, err
	} else {
		f.faceCache[key] = sfnt
		return sfnt, nil
	}
}
