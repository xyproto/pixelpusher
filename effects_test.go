package multirender

import (
	"fmt"
	"runtime"
	"testing"
)

func TestStretchContrast(t *testing.T) {
	// Define an image, with red pixels that are not fully bright and not fully dark
	pixels := []uint32{0xff800000, 0xff200000, 0xff500000} // ARGB, ARGB, ARGB

	pitch := int32(1) // 1 uint32 wide (1x3)
	cores := runtime.NumCPU()

	// Stretch the contrast
	StretchContrast(cores, pixels, pitch, 0.9)

	tmp := []uint32{0xffaca8a8, 0xffafa8a8, 0xff2da8a8}
	if pixels[0] != tmp[0] {
		fmt.Println(pixels)
		t.Fail()
	}
	if pixels[1] != tmp[1] {
		fmt.Println(pixels)
		t.Fail()
	}
	if pixels[2] != tmp[2] {
		fmt.Println(pixels)
		t.Fail()
	}
}
