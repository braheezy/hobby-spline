package bezier

import "math"

// Legendre-Gauss abscissae with n=24 (x_i values, defined at i=n as the roots of the nth order Legendre polynomial Pn(x))
var Tvalues = []float64{
	-0.0640568928626056260850430826247450385909,
	0.0640568928626056260850430826247450385909,
	-0.1911188674736163091586398207570696318404,
	0.1911188674736163091586398207570696318404,
	-0.3150426796961633743867932913198102407864,
	0.3150426796961633743867932913198102407864,
	-0.4337935076260451384870842319133497124524,
	0.4337935076260451384870842319133497124524,
	-0.5454214713888395356583756172183723700107,
	0.5454214713888395356583756172183723700107,
	-0.6480936519369755692524957869107476266696,
	0.6480936519369755692524957869107476266696,
	-0.7401241915785543642438281030999784255232,
	0.7401241915785543642438281030999784255232,
	-0.8200019859739029219539498726697452080761,
	0.8200019859739029219539498726697452080761,
	-0.8864155270044010342131543419821967550873,
	0.8864155270044010342131543419821967550873,
	-0.9382745520027327585236490017087214496548,
	0.9382745520027327585236490017087214496548,
	-0.9747285559713094981983919930081690617411,
	0.9747285559713094981983919930081690617411,
	-0.9951872199970213601799974097007368118745,
	0.9951872199970213601799974097007368118745,
}

// Legendre-Gauss weights with n=24 (w_i values, defined by a function linked to in the Bezier primer article)
var Cvalues = []float64{
	0.1279381953467521569740561652246953718517,
	0.1279381953467521569740561652246953718517,
	0.1258374563468282961213753825111836887264,
	0.1258374563468282961213753825111836887264,
	0.121670472927803391204463153476262425607,
	0.121670472927803391204463153476262425607,
	0.1155056680537256013533444839067835598622,
	0.1155056680537256013533444839067835598622,
	0.1074442701159656347825773424466062227946,
	0.1074442701159656347825773424466062227946,
	0.0976186521041138882698806644642471544279,
	0.0976186521041138882698806644642471544279,
	0.086190161531953275917185202983742667185,
	0.086190161531953275917185202983742667185,
	0.0733464814110803057340336152531165181193,
	0.0733464814110803057340336152531165181193,
	0.0592985849154367807463677585001085845412,
	0.0592985849154367807463677585001085845412,
	0.0442774388174198061686027482113382288593,
	0.0442774388174198061686027482113382288593,
	0.0285313886289336631813078159518782864491,
	0.0285313886289336631813078159518782864491,
	0.0123412297999871995468056670700372915759,
	0.0123412297999871995468056670700372915759,
}

func length(derivativeFunc DerivativeFunc) float64 {
	const z = 0.5
	l := len(Tvalues)

	sum := 0.0

	for i := 0; i < l; i++ {
		t := z*Tvalues[i] + z
		sum += Cvalues[i] * arcfn(t, derivativeFunc)
	}
	return z * sum
}

func arcfn(t float64, derivativeFunc func(float64) Point) float64 {
	d := derivativeFunc(t)
	l := d.X*d.X + d.Y*d.Y

	if d.threeDimensional {
		l += d.Z * d.Z
	}

	return math.Sqrt(l)
}

func compute(t float64, points []Point, use3d bool) Point {
	// shortcuts
	if t == 0 {
		points[0].t = 0
		return points[0]
	}
	order := len(points) - 1
	if t == 1 {
		points[order].t = 1
		return points[order]
	}

	mt := 1 - t
	p := points

	// constant?
	if order == 0 {
		points[0].t = t
		return points[0]
	}

	// linear?
	if order == 1 {
		result := Point{
			X: mt*p[0].X + t*p[1].X,
			Y: mt*p[0].Y + t*p[1].Y,
			t: t,
		}
		if use3d {
			result.Z = mt*p[0].Z + t*p[1].Z
		}
		return result
	}

	// quadratic/cubic curve?
	if order < 4 {
		mt2 := mt * mt
		t2 := t * t
		a, b, c, d := 0.0, 0.0, 0.0, 0.0
		if order == 2 {
			p = []Point{p[0], p[1], p[2], {}}
			a = mt2
			b = mt * t * 2
			c = t2
		} else if order == 3 {
			a = mt2 * mt
			b = mt2 * t * 3
			c = mt * t2 * 3
			d = t * t2
		}
		result := Point{
			X: a*p[0].X + b*p[1].X + c*p[2].X + d*p[3].X,
			Y: a*p[0].Y + b*p[1].Y + c*p[2].Y + d*p[3].Y,
			t: t,
		}
		if use3d {
			result.Z = a*p[0].Z + b*p[1].Z + c*p[2].Z + d*p[3].Z
		}
		return result
	}

	dCpts := make([]Point, len(points))
	copy(dCpts, points)

	// Perform de Casteljau's algorithm
	for len(dCpts) > 1 {
		for i := 0; i < len(dCpts)-1; i++ {
			dCpts[i].X = dCpts[i].X + (dCpts[i+1].X-dCpts[i].X)*t
			dCpts[i].Y = dCpts[i].Y + (dCpts[i+1].Y-dCpts[i].Y)*t

			if dCpts[i].threeDimensional {
				dCpts[i].Z = dCpts[i].Z + (dCpts[i+1].Z-dCpts[i].Z)*t
			}
		}
		dCpts = dCpts[:len(dCpts)-1]
	}
	dCpts[0].t = t
	return dCpts[0]
}

func normal2(t float64, df DerivativeFunc) Point {
	d := df(t)
	q := math.Sqrt(d.X*d.X + d.Y*d.Y)
	return Point{
		X: -d.Y / q,
		Y: d.X / q,
		t: t,
	}
}
func normal3(t float64, df DerivativeFunc) Point {
	// see http://stackoverflow.com/questions/25453159
	r1 := df(t)
	r2 := df(t + 0.01)
	q1 := math.Sqrt(r1.X*r1.X + r1.Y*r1.Y + r1.Z*r1.Z)
	q2 := math.Sqrt(r2.X*r2.X + r2.Y*r2.Y + r2.Z*r2.Z)

	r1.X /= q1
	r1.Y /= q1
	r1.Z /= q1
	r2.X /= q2
	r2.Y /= q2
	r2.Z /= q2
	// cross product
	c := Point{
		X: r2.Y*r1.Z - r2.Z*r1.Y,
		Y: r2.Z*r1.X - r2.X*r1.Z,
		Z: r2.X*r1.Y - r2.Y*r1.X,
	}
	m := math.Sqrt(c.X*c.X + c.Y*c.Y + c.Z*c.Z)
	c.X /= m
	c.Y /= m
	c.Z /= m
	// rotation matrix
	R := []float64{
		c.X * c.X,
		c.X*c.Y - c.Z,
		c.X*c.Z + c.Y,
		c.X*c.Y + c.Z,
		c.Y * c.Y,
		c.Y*c.Z - c.X,
		c.X*c.Z - c.Y,
		c.Y*c.Z + c.X,
		c.Z * c.Z,
	}
	// normal vector
	return Point{
		t: t,
		X: R[0]*r1.X + R[1]*r1.Y + R[2]*r1.Z,
		Y: R[3]*r1.X + R[4]*r1.Y + R[5]*r1.Z,
		Z: R[6]*r1.X + R[7]*r1.Y + R[8]*r1.Z,
	}
}

func curvature(t float64, d1 []Point, d2 []Point, use3d bool, kOnly bool) CurvatureVector {
	/*
		We're using the following formula for curvature:

		             x'y" - y'x"
		  k(t) = ------------------
		          (x'² + y'²)^(3/2)

		from https://en.wikipedia.org/wiki/Radius_of_curvature#Definition

		With it corresponding 3D counterpart:

		         sqrt( (y'z" - y"z')² + (z'x" - z"x')² + (x'y" - x"y')²)
		  k(t) = -------------------------------------------------------
		                    (x'² + y'² + z'²)^(3/2)
	*/
	d := compute(t, d1, use3d)
	dd := compute(t, d2, use3d)
	qdsum := d.X*d.X + d.Y*d.Y

	var num, dnm float64
	if use3d {
		num = math.Sqrt(
			math.Pow(d.Y*dd.Z-dd.Y*d.Z, 2) +
				math.Pow(d.Z*dd.X-dd.Z*d.X, 2) +
				math.Pow(d.X*dd.Y-dd.X*d.Y, 2))
		dnm = math.Pow(qdsum+d.Z*d.Z, 1.5)
	} else {
		num = d.X*dd.Y - d.Y*dd.X
		dnm = math.Pow(qdsum, 1.5)
	}

	if num == 0 || dnm == 0 {
		return CurvatureVector{}
	}

	k := num / dnm
	r := dnm / num

	// We're also computing the derivative of kappa, because
	// there is value in knowing the rate of change for the
	// curvature along the curve. And we're just going to
	// ballpark it based on an epsilon.
	var dk, adk float64
	if !kOnly {
		// compute k'(t) based on the interval before, and after it,
		// to at least try to not introduce forward/backward pass bias.
		pk := curvature(t-.001, d1, d2, use3d, true).K
		nk := curvature(t+.001, d1, d2, use3d, true).K
		dk = (nk - k + (k - pk)) / 2
		adk = (math.Abs(nk-k) + math.Abs(k-pk)) / 2
	}

	return CurvatureVector{
		K:   k,
		R:   r,
		dk:  dk,
		adk: adk,
	}

}
