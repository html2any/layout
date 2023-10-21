package fonts

// type key struct {
// 	Name  string
// 	Style FontStyle
// 	Size  float64
// }

// type Faces struct {
// 	fonts     *Fonts
// 	faceCache map[key]truetype.IndexableFace
// 	wordCache map[truetype.IndexableFace]map[rune]Int26_6
// }

// func NewFaces(fonts *Fonts) *Faces {
// 	f := &Faces{}
// 	f.faceCache = make(map[key]truetype.IndexableFace)
// 	f.fonts = fonts
// 	return f
// }

// func (f *Faces) CacheAllFace(max int) (*Faces, error) {
// 	if max == 0 {
// 		max = 50
// 	}
// 	f.wordCache = make(map[truetype.IndexableFace]map[rune]Int26_6)
// 	for family, data := range f.fonts.dataCache {
// 		for i := 5; i < max; i++ {
// 			faceKey := key{
// 				Name:  family.Name,
// 				Style: family.Style,
// 				Size:  float64(i),
// 			}
// 			if ff, err := freetype.ParseFont(data); err != nil {
// 				return nil, err
// 			} else {
// 				face := truetype.NewFace(ff, &truetype.Options{
// 					Size:    float64(i),
// 					Hinting: font.HintingFull,
// 				})
// 				f.wordCache[face] = make(map[rune]Int26_6)
// 				font.MeasureString(face, ".")
// 				f.faceCache[faceKey] = face
// 			}
// 		}
// 	}
// 	return f, nil
// }

// func (f *Faces) CacheFace(family string, size float64, style FontStyle) (*Faces, error) {
// 	faceKey := key{
// 		Name:  family,
// 		Style: style,
// 		Size:  size,
// 	}
// 	fontKey := Family{
// 		Name:  family,
// 		Style: style,
// 	}
// 	if _, ok := f.faceCache[faceKey]; ok {
// 		return f, nil
// 	} else if data, ok := f.fonts.dataCache[fontKey]; !ok {
// 		return nil, errors.New("font not found(" + fontKey.String() + "). please add font before get face")
// 	} else if ff, err := freetype.ParseFont(data); err != nil {
// 		return nil, err
// 	} else {
// 		face := truetype.NewFace(ff, &truetype.Options{
// 			Size: size,
// 			// Hinting: font.HintingFull,
// 		})
// 		font.MeasureString(face, ".")
// 		f.faceCache[faceKey] = face
// 	}
// 	return f, nil
// }

// func (f *Faces) GetFace(family string, size float64, style FontStyle) (truetype.IndexableFace, error) {
// 	k := key{
// 		Name:  family,
// 		Style: style,
// 		Size:  size,
// 	}
// 	if face, ok := f.faceCache[k]; ok {
// 		return face, nil
// 	} else if _, err := f.CacheFace(family, size, style); err != nil {
// 		return nil, err
// 	} else {
// 		// slog.Warn("FaceWarn", "The call to freetype.NewFace takes some times, please addFace before using it.")
// 		return f.faceCache[k], nil
// 	}
// }
