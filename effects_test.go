package multirender

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/xyproto/pf"
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

func TestStretchContrast2(t *testing.T) {
	// Define an image, with red pixels that are not fully bright and not fully dark
	pixels := make([]uint32, 320*200) // ARGB

	pitch := int32(320) // 320 uint32s wide (1x3)
	cores := runtime.NumCPU()

	// Set a few bits to gray
	pixels[1000] = 0xffc0c0c0
	pixels[2000] = 0xffc0c0c0
	pixels[42] = 0xffc0c0c0
	pixels[256] = 0xffc0c0c0

	// Turn on all the green bits
	green := pf.SetGreenBits
	pf.Map(cores, green, pixels)

	pixels1 := make([]uint32, len(pixels))
	pixels2 := make([]uint32, len(pixels))

	copy(pixels1, pixels)
	copy(pixels2, pixels)

	// Stretch the contrast
	StretchContrast(cores, pixels1, pitch, 0.9)

	// Stretch the contrast, concurrently
	StretchContrast2(cores, pixels2, pitch, 0.9)

	// Check that they are equal
	for i := range pixels {
		if pixels1[i] != pixels2[i] {
			fmt.Println("pixels1 and pixels2 differ at position", i)
			fmt.Printf("pixels1[i]: %x, pixels2[i]: %x\n", pixels1[i], pixels2[i])
			t.Fail()
			return
		}
	}

}
