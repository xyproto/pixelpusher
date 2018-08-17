package multirender

import (
	"fmt"
)

// Pos represents a position in 2D space
type Pos struct {
	x int32
	y int32
}

func (p *Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}
