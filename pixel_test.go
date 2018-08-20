package multirender

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
