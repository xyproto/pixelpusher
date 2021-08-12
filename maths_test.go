package pixelpusher

import (
	"fmt"
)

func ExampleScale() {
	fmt.Println(ScaleByte(0, 0, 255, 0, 255))
	fmt.Println(ScaleByte(20, 20, 80, 0, 255))
	fmt.Println(ScaleByte(80, 20, 80, 0, 255))
	fmt.Println(ScaleByte(255, 255, 255, 0, 255))
	fmt.Println(ScaleByte(0, 0, 0, 0, 255))
	// Output:
	// 0
	// 0
	// 255
	// 255
	// 0
}
