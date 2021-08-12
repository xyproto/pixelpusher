package fauxgl

import (
	"math"

	"github.com/fogleman/simplify"
)

type Mesh struct {
	Triangles []*Triangle
	Lines     []*Line
	box       *Box
}

func NewEmptyMesh() *Mesh {
	return &Mesh{}
}

func NewMesh(triangles []*Triangle, lines []*Line) *Mesh {
	return &Mesh{triangles, lines, nil}
}

func NewTriangleMesh(triangles []*Triangle) *Mesh {
	return &Mesh{triangles, nil, nil}
}

func NewLineMesh(lines []*Line) *Mesh {
	return &Mesh{nil, lines, nil}
}

func (m *Mesh) dirty() {
	m.box = nil
}

func (m *Mesh) Copy() *Mesh {
	triangles := make([]*Triangle, len(m.Triangles))
	lines := make([]*Line, len(m.Lines))
	for i, t := range m.Triangles {
		a := *t
		triangles[i] = &a
	}
	for i, l := range m.Lines {
		a := *l
		lines[i] = &a
	}
	return NewMesh(triangles, lines)
}

func (a *Mesh) Add(b *Mesh) {
	a.Triangles = append(a.Triangles, b.Triangles...)
	a.Lines = append(a.Lines, b.Lines...)
	a.dirty()
}

func (m *Mesh) SetColor(c Color) {
	for _, t := range m.Triangles {
		t.SetColor(c)
	}
}

func (m *Mesh) Volume() float64 {
	var v float64
	for _, t := range m.Triangles {
		p1 := t.V1.Position
		p2 := t.V2.Position
		p3 := t.V3.Position
		v += p1.X*(p2.Y*p3.Z-p3.Y*p2.Z) - p2.X*(p1.Y*p3.Z-p3.Y*p1.Z) + p3.X*(p1.Y*p2.Z-p2.Y*p1.Z)
	}
	return math.Abs(v / 6)
}

func (m *Mesh) SurfaceArea() float64 {
	var a float64
	for _, t := range m.Triangles {
		a += t.Area()
	}
	return a
}

func smoothNormalsThreshold(normal Vector, normals []Vector, threshold float64) Vector {
	result := Vector{}
	for _, x := range normals {
		if x.Dot(normal) >= threshold {
			result = result.Add(x)
		}
	}
	return result.Normalize()
}

func (m *Mesh) SmoothNormalsThreshold(radians float64) {
	threshold := math.Cos(radians)
	lookup := make(map[Vector][]Vector)
	for _, t := range m.Triangles {
		lookup[t.V1.Position] = append(lookup[t.V1.Position], t.V1.Normal)
		lookup[t.V2.Position] = append(lookup[t.V2.Position], t.V2.Normal)
		lookup[t.V3.Position] = append(lookup[t.V3.Position], t.V3.Normal)
	}
	for _, t := range m.Triangles {
		t.V1.Normal = smoothNormalsThreshold(t.V1.Normal, lookup[t.V1.Position], threshold)
		t.V2.Normal = smoothNormalsThreshold(t.V2.Normal, lookup[t.V2.Position], threshold)
		t.V3.Normal = smoothNormalsThreshold(t.V3.Normal, lookup[t.V3.Position], threshold)
	}
}

func (m *Mesh) SmoothNormals() {
	lookup := make(map[Vector]Vector)
	for _, t := range m.Triangles {
		lookup[t.V1.Position] = lookup[t.V1.Position].Add(t.V1.Normal)
		lookup[t.V2.Position] = lookup[t.V2.Position].Add(t.V2.Normal)
		lookup[t.V3.Position] = lookup[t.V3.Position].Add(t.V3.Normal)
	}
	for k, v := range lookup {
		lookup[k] = v.Normalize()
	}
	for _, t := range m.Triangles {
		t.V1.Normal = lookup[t.V1.Position]
		t.V2.Normal = lookup[t.V2.Position]
		t.V3.Normal = lookup[t.V3.Position]
	}
}

func (m *Mesh) UnitCube() Matrix {
	const r = 0.5
	return m.FitInside(Box{Vector{-r, -r, -r}, Vector{r, r, r}}, Vector{0.5, 0.5, 0.5})
}

func (m *Mesh) BiUnitCube() Matrix {
	const r = 1
	return m.FitInside(Box{Vector{-r, -r, -r}, Vector{r, r, r}}, Vector{0.5, 0.5, 0.5})
}

func (m *Mesh) MoveTo(position, anchor Vector) Matrix {
	matrix := Translate(position.Sub(m.BoundingBox().Anchor(anchor)))
	m.Transform(matrix)
	return matrix
}

func (m *Mesh) Center() Matrix {
	return m.MoveTo(Vector{}, Vector{0.5, 0.5, 0.5})
}

func (m *Mesh) FitInside(box Box, anchor Vector) Matrix {
	scale := box.Size().Div(m.BoundingBox().Size()).MinComponent()
	extra := box.Size().Sub(m.BoundingBox().Size().MulScalar(scale))
	matrix := Identity()
	matrix = matrix.Translate(m.BoundingBox().Min.Negate())
	matrix = matrix.Scale(Vector{scale, scale, scale})
	matrix = matrix.Translate(box.Min.Add(extra.Mul(anchor)))
	m.Transform(matrix)
	return matrix
}

func (m *Mesh) BoundingBox() Box {
	if m.box == nil {
		box := EmptyBox
		for _, t := range m.Triangles {
			box = box.Extend(t.BoundingBox())
		}
		for _, l := range m.Lines {
			box = box.Extend(l.BoundingBox())
		}
		m.box = &box
	}
	return *m.box
}

func (m *Mesh) Transform(matrix Matrix) {
	for _, t := range m.Triangles {
		t.Transform(matrix)
	}
	for _, l := range m.Lines {
		l.Transform(matrix)
	}
	m.dirty()
}

func (m *Mesh) ReverseWinding() {
	for _, t := range m.Triangles {
		t.ReverseWinding()
	}
}

func (m *Mesh) Simplify(factor float64) {
	st := make([]*simplify.Triangle, len(m.Triangles))
	for i, t := range m.Triangles {
		v1 := simplify.Vector(t.V1.Position)
		v2 := simplify.Vector(t.V2.Position)
		v3 := simplify.Vector(t.V3.Position)
		st[i] = simplify.NewTriangle(v1, v2, v3)
	}
	sm := simplify.NewMesh(st)
	sm = sm.Simplify(factor)
	m.Triangles = make([]*Triangle, len(sm.Triangles))
	for i, t := range sm.Triangles {
		v1 := Vector(t.V1)
		v2 := Vector(t.V2)
		v3 := Vector(t.V3)
		m.Triangles[i] = NewTriangleForPoints(v1, v2, v3)
	}
	m.dirty()
}

func (m *Mesh) SaveSTL(path string) error {
	return SaveSTL(path, m)
}

func (m *Mesh) Silhouette(eye Vector, offset float64) *Mesh {
	return silhouette(m, eye, offset)
}

func (m *Mesh) SplitTriangles(maxEdgeLength float64) {
	var triangles []*Triangle

	var split func(t *Triangle)

	split = func(t *Triangle) {
		v1 := t.V1
		v2 := t.V2
		v3 := t.V3
		p1 := v1.Position
		p2 := v2.Position
		p3 := v3.Position
		d12 := p1.Distance(p2)
		d23 := p2.Distance(p3)
		d31 := p3.Distance(p1)
		max := math.Max(d12, math.Max(d23, d31))
		if max <= maxEdgeLength {
			triangles = append(triangles, t)
		} else if d12 == max {
			v := InterpolateVertexes(v1, v2, v3, VectorW{0.5, 0.5, 0, 1})
			t1 := NewTriangle(v3, v1, v)
			t2 := NewTriangle(v2, v3, v)
			split(t1)
			split(t2)
		} else if d23 == max {
			v := InterpolateVertexes(v1, v2, v3, VectorW{0, 0.5, 0.5, 1})
			t1 := NewTriangle(v1, v2, v)
			t2 := NewTriangle(v3, v1, v)
			split(t1)
			split(t2)
		} else {
			v := InterpolateVertexes(v1, v2, v3, VectorW{0.5, 0, 0.5, 1})
			t1 := NewTriangle(v2, v3, v)
			t2 := NewTriangle(v1, v2, v)
			split(t1)
			split(t2)
		}
	}

	for _, t := range m.Triangles {
		split(t)
	}

	m.Triangles = triangles
	m.dirty()
}

func (m *Mesh) SharpEdges(angleThreshold float64) *Mesh {
	type Edge struct {
		A, B Vector
	}

	makeEdge := func(a, b Vector) Edge {
		if a.Less(b) {
			return Edge{a, b}
		}
		return Edge{b, a}
	}

	var lines []*Line
	other := make(map[Edge]*Triangle)
	for _, t := range m.Triangles {
		p1 := t.V1.Position
		p2 := t.V2.Position
		p3 := t.V3.Position
		e1 := makeEdge(p1, p2)
		e2 := makeEdge(p2, p3)
		e3 := makeEdge(p3, p1)
		for _, e := range []Edge{e1, e2, e3} {
			if u, ok := other[e]; ok {
				a := math.Acos(t.Normal().Dot(u.Normal()))
				if a > angleThreshold {
					lines = append(lines, NewLineForPoints(e.A, e.B))
				}
			}
		}
		other[e1] = t
		other[e2] = t
		other[e3] = t
	}
	return NewLineMesh(lines)
}
