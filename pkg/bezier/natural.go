package bezier

// controlPoints calculates the control points for Bezier curve segments
func controlPoints(points []float64) ([]float64, []float64) {
	n := len(points) - 1
	if n < 1 {
		return nil, nil
	}
	a, b, r := make([]float64, n), make([]float64, n), make([]float64, n)
	a[0], b[0], r[0] = 0, 2, points[0]+2*points[1]
	for i := 1; i < n-1; i++ {
		a[i] = 1
		b[i] = 4
		r[i] = 4*points[i] + 2*points[i+1]
	}
	a[n-1], b[n-1], r[n-1] = 2, 7, 8*points[n-1]+points[n]

	// Forward elimination
	for i := 1; i < n; i++ {
		m := a[i] / b[i-1]
		b[i] -= m
		r[i] -= m * r[i-1]
	}
	// Back substitution
	a[n-1] = r[n-1] / b[n-1]
	for i := n - 2; i >= 0; i-- {
		a[i] = (r[i] - a[i+1]) / b[i]
	}

	// Calculate cp1 and cp2 for each segment
	cp1, cp2 := make([]float64, n), make([]float64, n)
	cp1[n-1] = (points[n] + a[n-1]) / 2
	for i := 0; i < n-1; i++ {
		cp1[i] = 2*points[i+1] - a[i+1]
		cp2[i] = a[i]
	}
	return cp1, cp2
}

// NaturalCubicSpline creates cubic Bezier curves that naturally interpolate through the given points.
func NaturalCubicSpline(points []Point) ([]Point, error) {

	return nil, nil
}
