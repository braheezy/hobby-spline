package bezier

import (
	"errors"
	"math"
)

// CreateHobbySpline: given a set of points, fit a Bézier spline to them.
// The chosen splines tend to have pleasing, rounded shapes.
//
// The output array will have 3n - 2 points where n is the number of input points.
// The output will contain every point in the input (these become knots in the
// Bézier spline), interspersed with pairs of new points which define positions
// of handle points on the spline. All points are in the same coordinate space.
func CreateHobbySpline(points []Point, omega float64) ([]Point, error) {
	// Solving is only possible with at least two points
	if len(points) < 2 {
		return nil, errors.New("not enough points")
	}

	// n is defined such that the points can be numbered P[0]...P[n], i.e, such that
	// there are a total of n+1 points
	n := len(points) - 1

	// chords[i] is the vector from P[i] to P[i+1].
	// d[i] is the length of the i-th chord.
	chords := make([]Point, n)
	d := make([]float64, n)
	for i := 0; i < n; i++ {
		chords[i] = vSub(points[i+1], points[i])
		d[i] = chords[i].Length()
		// no chord can be zero-length (i.e. no two successive points can be the same)
		if d[i] == 0 {
			return nil, errors.New("zero-length chord")
		}
	}

	// gamma[i] is the signed turning angle at P[i], i.e. the angle between
	// the chords from P[i-1] to P[i] and from P[i] to P[i+1].
	// gamma[0] is undefined; gamma[n] is artificially defined to be zero
	gamma := make([]float64, n+1)
	for i := 1; i < n; i++ {
		gamma[i] = vAngleBetween(chords[i-1], chords[i])
	}

	// Set up the system of linear equations (Jackowski, formula 38).
	// We're representing this system as a tridiagonal matrix, because
	// we can solve such a system in O(n) time using the Thomas algorithm.
	//
	// Here, A, B, and C are the matrix diagonals and D is the right-hand side.
	// See Wikipedia for a more detailed explanation:
	// https://en.wikipedia.org/wiki/Tridiagonal_matrix_algorithm
	A := make([]float64, n+1)
	B := make([]float64, n+1)
	C := make([]float64, n+1)
	D := make([]float64, n+1)

	B[0] = 2 + omega
	C[0] = 2*omega + 1
	D[0] = -1 * C[0] * gamma[1]

	for i := 1; i < n; i++ {
		A[i] = 1 / d[i-1]
		B[i] = (2*d[i-1] + 2*d[i]) / (d[i-1] * d[i])
		C[i] = 1 / d[i]
		D[i] = (-1 * (2*gamma[i]*d[i] + gamma[i+1]*d[i-1])) / (d[i-1] * d[i])
	}

	A[n] = 2*omega + 1
	B[n] = 2 + omega
	D[n] = 0

	// Solve the tridiagonal matrix of equations using the Thomas algorithm,
	// yielding the alpha angles for each point (these are the angles between
	// each chord[i] and the vector c0[i] - P[i], i.e. the vector from knot i
	// to the subsequent control point, which is tangent to the curve at P[i]).
	alpha := thomas(A, B, C, D)

	// Use alpha (the chord angle) and gamma (the turning angle of the chord
	// polyline) to solve for beta at each point (beta is like alpha, but for
	// the chord and handle vector arriving at P[i] rather than leaving from it).
	beta := make([]float64, n)

	for i := 0; i < n-1; i++ {
		beta[i] = -1*gamma[i+1] - alpha[i+1]
	}

	beta[n-1] = -1 * alpha[n]

	// Now that we have the angles between the handle vector and the chord
	// both arriving at and leaving from each point, we can solve for the
	// positions of the handle (control) points themselves.
	c0 := make([]Point, n)
	c1 := make([]Point, n)
	for i := 0; i < n; i++ {
		// Compute the magnitudes of the handle vectors at this point.
		// (Jackowski, formula 22)
		a := (rho(alpha[i], beta[i]) * d[i]) / 3
		b := (rho(beta[i], alpha[i]) * d[i]) / 3

		// Use the magnitudes, and the chords and turning angles, to find
		// the positions of the control points in the global coordinate space.
		c0[i] = vAdd(points[i], Scale(Normalize(Rotate(chords[i], alpha[i])), a))
		c1[i] = vSub(points[i+1], Scale(Normalize(Rotate(chords[i], -1*beta[i])), b))
	}

	// Finally, gather up and return the spline points (both knots and
	// control points) as a single ordered list of [x, y] pairs.
	var result []Point
	for i := 0; i < n; i++ {
		result = append(result, points[i], c0[i], c1[i])
	}
	result = append(result, points[n])

	return result, nil
}

func thomas(A, B, C, D []float64) []float64 {
	// A, B, and C are diagonals of the matrix. B is the main diagonal.
	// D is the vector on the right-hand-side of the equation.

	// Both B and D will have n elements. The arrays A and C will have
	// length n as well, but each has one fewer element than B (the values
	// A[0] and C[n-1] are undefined).

	// Note: n is defined here so that B[n] is valid, i.e. we are solving
	// a system of n+1 equations.
	n := len(B) - 1

	// Step 1: forward sweep to eliminate A[i] from each equation

	// allocate arrays for modified C and D coefficients
	// (p stands for prime)
	Cp := make([]float64, n+1)
	Dp := make([]float64, n+1)

	Cp[0] = C[0] / B[0]
	Dp[0] = D[0] / B[0]

	for i := 1; i <= n; i++ {
		denom := B[i] - A[i]*Cp[i-1]
		Cp[i] = C[i] / denom
		Dp[i] = (D[i] - A[i]*Dp[i-1]) / denom
	}

	// Step 2: back substitution to solve for X
	X := make([]float64, n+1)
	// start at the end, then work backwards to solve for each X[i]
	X[n] = Dp[n]
	for i := n - 1; i >= 0; i-- {
		X[i] = Dp[i] - Cp[i]*X[i+1]
	}
	return X
}

// Rho is the 'velocity function' that computes the length of the handles for
// the Bézier spline.
//
// Once the angles alpha and beta have been computed for each knot (which
// determine the direction from the knot to each of its neighboring handles),
// this function is used to compute the lengths of the vectors from the knot to
// those handles. Combining the length and angle together lets us solve for the
// handle positions.
//
// The exact choice of function is somewhat arbitrary. The aim is to return
// handle lengths that produce a Bézier curve which is a good approximation of a
// circular arc for points near the knot.
//
// Hobby and Knuth both proposed multiple candidate functions. This code uses
// the function from Jackowski formula 28, due to its simplicity. For other
// choices see Jackowski, section 5.
func rho(alpha, beta float64) float64 {
	c := 2.0 / 3.0
	return 2 / (1 + c*math.Cos(beta) + (1-c)*math.Cos(alpha))
}
