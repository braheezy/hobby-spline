package bezier

import (
	"errors"
	"fmt"
)

type CurvatureVector struct {
	K   float64
	R   float64
	dk  float64
	adk float64
}

type Point struct {
	X float64
	Y float64
	Z float64
	// Parametrization variable
	t                float64
	threeDimensional bool
}

type Bezier struct {
	Points           []Point
	order            int
	threeDimensional bool
	dims             []string
	dimLength        int
	dpoints          [][]Point
}

type DerivativeFunc func(float64) Point

// NewBezier creates a new Bezier curve with the given points.
func NewBezier(use3d bool, points ...Point) (*Bezier, error) {
	count := len(points)
	if count < 3 {
		return nil, errors.New("at least 3 points are required")
	} else if count > 12 {
		return nil, errors.New("at most 12 points are allowed")
	}
	if use3d {
		if count != 8 && count != 9 && count != 12 {
			return nil, fmt.Errorf("3D Bezier curves require 8, 9, or 12 points, got %v", count)
		}
	}

	b := &Bezier{
		Points: points,
		order:  len(points) - 1,
		dims: []string{
			"x",
			"y",
		},
		dimLength: 2,
	}

	if use3d {
		b.threeDimensional = true
		b.dims = []string{
			"x",
			"y",
			"z",
		}
		b.dimLength = 3
		for i := 0; i < len(b.Points); i++ {
			b.Points[i].threeDimensional = true
		}
	}
	b.update()

	return b, nil
}

func (b *Bezier) Length() float64 {
	return length(b.derivative)
}

func (b *Bezier) Get(t float64) Point {
	return compute(t, b.Points, b.threeDimensional)
}

func (b *Bezier) Normal(t float64) Point {
	if b.threeDimensional {
		return normal3(t, b.derivative)
	} else {
		return normal2(t, b.derivative)
	}
}
func (b *Bezier) Curvature(t float64) CurvatureVector {
	return curvature(t, b.dpoints[0], b.dpoints[1], b.threeDimensional, false)
}

func (b *Bezier) update() {
	b.setDpoints()
}

func (b *Bezier) setDpoints() {
	var levels [][]Point

	points := b.Points
	for len(points) > 1 {
		var level []Point
		c := len(points) - 1
		for i := 0; i < c; i++ {
			dpt := Point{
				X: float64(c) * (points[i+1].X - points[i].X),
				Y: float64(c) * (points[i+1].Y - points[i].Y),
			}
			if b.threeDimensional {
				dpt.Z = float64(c) * (points[i+1].Z - points[i].Z)
			}
			level = append(level, dpt)
		}
		levels = append(levels, level)
		// Update points to the last calculated derivative level
		points = level
	}

	b.dpoints = levels
}

func (b *Bezier) derivative(t float64) Point {
	return compute(t, b.dpoints[0], b.threeDimensional)
}
