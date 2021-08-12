package main

import (
	"errors"
	"fmt"
	"github.com/fogleman/fauxgl"
	"github.com/nfnt/resize"
	"github.com/xyproto/pixelpusher"
	"image"
	"math"
)

func LoadMeshOBJ(filename string) (*fauxgl.Mesh, error) {
	fmt.Printf("Loading %s... ", filename)
	// load the mesh
	mesh, err := fauxgl.LoadOBJ(filename)
	if err != nil {
		return nil, err
	}

	// fit mesh in a bi-unit cube centered at the origin
	mesh.BiUnitCube()

	// smooth the normals
	mesh.SmoothNormalsThreshold(fauxgl.Radians(30))

	fmt.Println("ok")

	// Return the processed mesh
	return mesh, nil
}

func DrawMesh(pixels []uint32, pitch int32, mesh *fauxgl.Mesh, cameraAngle float32, hexColor string) error {
	const (
		scale = 2  // optional supersampling
		fovy  = 20 // vertical field of view in degrees
		near  = 1  // near clipping plane
		far   = 10 // far clipping plane
	)
	var (
		up    = fauxgl.V(0, 1, 0)                    // up vector
		light = fauxgl.V(-0.75, 1, 0.25).Normalize() // light direction
		color = fauxgl.HexColor(hexColor)            // object color
	)

	// Camera position, calculated from cameraAngle
	cameraX := math.Cos(float64(cameraAngle/10.0)) * 4.0
	cameraY := math.Sin(float64(cameraAngle/10.0)) * 4.0
	cameraZ := math.Cos(float64(cameraAngle/10.0)) * 4.0
	camera := fauxgl.V(cameraX, cameraY, cameraZ)

	// View center position
	centerX := math.Cos(float64(cameraAngle/10.0)) * -1.0
	centerY := math.Sin(float64(cameraAngle/10.0)) * -1.0
	center := fauxgl.V(centerX, centerY, 0)

	// create a rendering context
	context := fauxgl.NewContext(width*scale, height*scale)

	// white transparent background
	context.ClearColorBufferWith(fauxgl.HexColor("#fffffffff")) // #FFF8E3

	// create transformation matrix and light direction
	aspect := float64(width) / float64(height)
	matrix := fauxgl.LookAt(camera, center, up).Perspective(fovy, aspect, near, far)

	// use builtin phong shader
	shader := fauxgl.NewPhongShader(matrix, light, camera)
	shader.ObjectColor = color
	context.Shader = shader

	// render
	context.DrawMesh(mesh)

	// downsample image for antialiasing
	resized := resize.Resize(width, height, context.Image(), resize.Bilinear)

	img, ok := resized.(*image.RGBA)
	if !ok {
		return errors.New("Could not convert image to *image.RGBA")
	}

	return pixelpusher.BlitImage(pixels, pitch, img)
}
