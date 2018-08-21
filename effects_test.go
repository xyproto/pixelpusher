package multirender

import (
	"fmt"
	"runtime"
	"testing"
)

func wipTestStretchContrast(t *testing.T) {
	// Define an image, with red pixels that are not fully bright and not fully dark
	pixels := []uint32{0xff800000, 0xff200000, 0xff500000} // ARGB, ARGB, ARGB

	//fmt.Println("RED0", Red(pixels[0]))
	//fmt.Println("RED1", Red(pixels[1]))
	//fmt.Println("GREEN0", Green(pixels[0]))
	//fmt.Println("GREEN1", Green(pixels[1]))
	//fmt.Println("BLUE0", Blue(pixels[0]))
	//fmt.Println("BLUE1", Blue(pixels[1]))
	//fmt.Println("A0", Alpha(pixels[0]))
	//fmt.Println("A1", Alpha(pixels[1]))

	pitch := int32(1) // 1 uint32 wide (1x3)
	cores := runtime.NumCPU()

	// Stretch the contrast
	wipStretchContrast(cores, pixels, pitch, 0, 255)

	//fmt.Println("STRETCHED:")
	//fmt.Println("RED0", Red(pixels[0]))
	//fmt.Println("RED1", Red(pixels[1]))
	//fmt.Println("GREEN0", Green(pixels[0]))
	//fmt.Println("GREEN1", Green(pixels[1]))
	//fmt.Println("BLUE0", Blue(pixels[0]))
	//fmt.Println("BLUE1", Blue(pixels[1]))
	//fmt.Println("A0", Alpha(pixels[0]))
	//fmt.Println("A1", Alpha(pixels[1]))
	//fmt.Printf("1: red %d hex %x\n", Red(pixels[0]), pixels[0])
	//fmt.Printf("2: red %d hex %x\n", Red(pixels[1]), pixels[1])

	if Red(pixels[0]) != 0xff { // != 0xffff0000
		fmt.Printf("0: red %d hex %x\n", Red(pixels[0]), pixels[0])
		t.Fail()
	}

	// Check if the second pixel was stretched to the bottom
	if Red(pixels[1]) != 00 { // != 0xff000000 {
		fmt.Printf("1: red %d hex %x\n", Red(pixels[1]), pixels[1])
		t.Fail()
	}

	// Check if the third pixel was stretched to the middle
	if Red(pixels[2]) != 0x7f { // != 0xff000000 {
		fmt.Printf("2: red %d hex %x\n", Red(pixels[2]), pixels[2])
		t.Fail()
	}

}
