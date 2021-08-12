package pixelpusher

import (
	"fmt"
)

// Pos represents a position in 2D space
type Pos struct {
	x int32
	y int32
}

// NewPos creates a new position
func NewPos(x, y int32) *Pos {
	return &Pos{x, y}
}

func (p *Pos) String() string {
	return fmt.Sprintf("(%d, %d)", p.x, p.y)
}
