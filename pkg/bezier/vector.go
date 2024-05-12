package bezier

import "math"

func (v Point) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func vAdd(a Point, b Point) Point {
	return Point{X: a.X + b.X, Y: a.Y + b.Y}
}

func vSub(a Point, b Point) Point {
	return Point{X: a.X - b.X, Y: a.Y - b.Y}
}

func Scale(v Point, s float64) Point {
	return Point{X: v.X * s, Y: v.Y * s}
}

func vDistance(from Point, to Point) float64 {
	return vSub(to, from).Length()
}

func Rotate(v Point, angle float64) Point {
	return Point{
		X: v.X*math.Cos(angle) - v.Y*math.Sin(angle),
		Y: v.X*math.Sin(angle) + v.Y*math.Cos(angle),
	}
}

func (v *Point) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}

func vAngleBetween(v Point, w Point) float64 {
	// w[1] * v[0] - w[0] * v[1], v[0] * w[0] + v[1] * w[1]
	angleX, angleY := w.Y*v.X-w.X*v.Y, v.X*w.X+v.Y*w.Y
	return math.Atan2(angleX, angleY)
}

func Normalize(v Point) Point {
	l := v.Length()
	return Point{
		X: v.X / l,
		Y: v.Y / l,
	}
}
