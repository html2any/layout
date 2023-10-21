package fonts

import (
	"fmt"
	"testing"
)

func TestMeasure(t *testing.T) {
	// loadFonts()
	fontcache := NewFonts()
	for _, fontfile := range []string{"../example/fonts/300.ttf", "../example/fonts/400.ttf"} {
		if err := fontcache.AddFont("NotoSansSC", fontfile, Regular); err != nil {
			t.Error(err)
		}
	}
	if sfnt, err := NewSFNT(fontcache).CacheAll(); err != nil {
		t.Error(err)
	} else if mw, err := sfnt.MeasureString("hello,你好,こんにちは,안녕하세요", "NotoSansSC", 12, ""); err != nil {
		t.Error(err)
	} else if lines, wds, err := sfnt.SplitLines(mw/2, "hello,你好,こんにちは,안녕하세요", "NotoSansSC", 12, ""); err != nil {
		t.Error(err)
	} else {
		for _, line := range lines {
			fmt.Println(line)
		}
		for _, wd := range wds {
			fmt.Println(wd)
		}
	}
}
