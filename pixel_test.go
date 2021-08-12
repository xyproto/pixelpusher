package pixelpusher

import (
	"testing"
)

func TestColorPacking(t *testing.T) {
	cv := RGBAToColorValue(255, 127, 63, 31)
	r, g, b, a := ColorValueToRGBA(cv)
	if r != 255 {
		t.Fail()
	}
	if g != 127 {
		t.Fail()
	}
	if b != 63 {
		t.Fail()
	}
	if a != 31 {
		t.Fail()
	}
}

func TestBlend(t *testing.T) {
	red := RGBAToColorValue(255, 0, 0, 127)
	blue := RGBAToColorValue(0, 0, 255, 255)
	r, g, b, a := ColorValueToRGBA(Blend(red, blue))
	//fmt.Println(r, g, b, a)
	if r != 63 {
		t.Fail()
	}
	if g != 0 {
		t.Fail()
	}
	if b != 127 {
		t.Fail()
	}
	if a != 191 {
		t.Fail()
	}
}
