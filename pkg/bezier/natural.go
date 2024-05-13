package bezier

import "fmt"

// controlPoints calculates the control points for Bezier curve segments

// controlPoints calculates control points for cubic Bezier curves.
func controlPoints(x []float64) ([]float64, []float64) {
	n := len(x) - 1
	if n < 1 {
		return nil, nil // No control points can be calculated if there are not enough points.
	}

	a := make([]float64, n)
	b := make([]float64, n)
	r := make([]float64, n)

	// Initialize the arrays for the tridiagonal matrix algorithm
	a[0] = 0
	b[0] = 2
	r[0] = x[0] + 2*x[1]
	for i := 1; i < n-1; i++ {
		a[i] = 1
		b[i] = 4
		r[i] = 4*x[i] + 2*x[i+1]
	}
	a[n-1] = 2
	b[n-1] = 7
	r[n-1] = 8*x[n-1] + x[n]

	// Forward elimination to transform the matrix into upper triangular
	for i := 1; i < n; i++ {
		m := a[i] / b[i-1]
		b[i] -= m
		r[i] -= m * r[i-1]
	}

	// Back substitution to solve for the first set of control points 'a'
	a[n-1] = r[n-1] / b[n-1]
	for i := n - 2; i >= 0; i-- {
		a[i] = (r[i] - a[i+1]) / b[i]
	}

	// Calculate the second set of control points 'b'
	b[n-1] = (x[n] + a[n-1]) / 2
	for i := 0; i < n-1; i++ {
		b[i] = 2*x[i+1] - a[i+1]
	}

	return a, b
}

// NaturalCubicSpline creates cubic Bezier curves that naturally interpolate through the given points.
func NaturalCubicSpline(points []Point) ([]Point, error) {
	if len(points) < 3 {
		return nil, fmt.Errorf("at least three points are required to form a spline")
	}

	// Separate points into x and y components
	x, y := make([]float64, len(points)), make([]float64, len(points))
	for i, p := range points {
		x[i], y[i] = p.X, p.Y
	}

	// Calculate control points
	cp1x, cp2x := controlPoints(x)
	cp1y, cp2y := controlPoints(y)

	// Generate the Bezier points (can be adjusted to create actual curve paths or visualizations)
	bezierPoints := make([]Point, 0)
	for i := 0; i < len(points)-1; i++ {
		bezierPoints = append(bezierPoints, points[i], Point{X: cp1x[i], Y: cp1y[i]}, Point{X: cp2x[i], Y: cp2y[i]})
	}
	bezierPoints = append(bezierPoints, points[len(points)-1])
	return bezierPoints, nil
}
