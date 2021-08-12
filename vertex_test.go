package pixelpusher

import (
	"fmt"
	"testing"
)

func TestColor(t *testing.T) {
	x := float32(0.0)
	y := float32(1.0)
	z := float32(2.0)
	r := uint8(3)
	g := uint8(4)
	b := uint8(5)
	a := uint8(7)
	v := NewVertex(x, y, z, r, g, b, a)
	if v.X() != x {
		t.Fail()
	}
	if v.Y() != y {
		t.Error(v.Y(), "!=", y)
	}
	if v.Z() != z {
		t.Error(v.Z(), "!=", z)
	}
	if v.R() != r {
		t.Error(v.R(), "!=", r)
	}
	if v.G() != g {
		t.Error(v.G(), "!=", g)
	}
	if v.B() != b {
		t.Error(v.B(), "!=", b)
	}
	if v.A() != a {
		t.Error(v.A(), "!=", a)
	}
}

func ExampleNewVertex() {
	v := NewVertex(0, 1, 2, 3, 4, 5, 6)
	fmt.Println(v)
	// Output:
	// v(0, 1, 2) color(3, 4, 5, 6)
}
