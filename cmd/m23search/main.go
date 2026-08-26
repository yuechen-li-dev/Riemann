// Command m23search independently refines the three active M23 equations.
// Its binary64 output is research evidence only; semantic/m23.go verifies the
// resulting rational box with directed rational interval arithmetic.
package main

import (
	"fmt"
	"math"
)

const a = math.Pi / 2

func f0(t float64) float64 {
	c0 := 1 + math.Pi/4 - 2/math.Pi
	w0 := 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi)
	return 1 + 2*(math.Cos(t)-1)/(t*t) + 2*(1-c0)*math.Sin(t)/t + 2*w0*math.Cos(t)
}

func h(r, t float64) float64 {
	c, s := math.Cos(a*r), math.Sin(a*r)
	return -math.Pi*c*math.Sin(t)/t + 2*(2*c/math.Pi-r*s)*math.Cos(t) + 2*math.Cos(r*t)
}

func q(x [3]float64) [3]float64 {
	m, r, t := x[0], x[1], x[2]
	c, s := math.Cos(a*r), math.Sin(a*r)
	cp := -a * s
	bp := -2*s - a*r*c
	hr := -math.Pi*cp*math.Sin(t)/t + 2*bp*math.Cos(t) - 2*t*math.Sin(r*t)
	ht := -math.Pi*c*(t*math.Cos(t)-math.Sin(t))/(t*t) - 2*(2*c/math.Pi-r*s)*math.Sin(t) - 2*r*math.Sin(r*t)
	// A central difference is sufficient here because this command is explicitly
	// non-certifying and only supplies decimal box centers.
	e := 1e-6
	f0t := (f0(t+e) - f0(t-e)) / (2 * e)
	return [3]float64{f0(t) + m*h(r, t), f0t + m*ht, cp*h(r, t) - c*hr}
}

func solve3(a [3][3]float64, b [3]float64) [3]float64 {
	for i := 0; i < 3; i++ {
		p := i
		for k := i + 1; k < 3; k++ {
			if math.Abs(a[k][i]) > math.Abs(a[p][i]) {
				p = k
			}
		}
		a[i], a[p], b[i], b[p] = a[p], a[i], b[p], b[i]
		for k := i + 1; k < 3; k++ {
			f := a[k][i] / a[i][i]
			for j := i; j < 3; j++ {
				a[k][j] -= f * a[i][j]
			}
			b[k] -= f * b[i]
		}
	}
	var x [3]float64
	for i := 2; i >= 0; i-- {
		s := b[i]
		for j := i + 1; j < 3; j++ {
			s -= a[i][j] * x[j]
		}
		x[i] = s / a[i][i]
	}
	return x
}

func main() {
	x := [3]float64{0.006995089360323918, 3.507853493599018, 1.7953133218061055}
	for n := 0; n < 12; n++ {
		v := q(x)
		var j [3][3]float64
		for k := 0; k < 3; k++ {
			e := 1e-5
			if k == 0 {
				e = 1e-8
			}
			y, z := x, x
			y[k], z[k] = y[k]+e, z[k]-e
			qy, qz := q(y), q(z)
			for i := 0; i < 3; i++ {
				j[i][k] = (qy[i] - qz[i]) / (2 * e)
			}
		}
		d := solve3(j, [3]float64{-v[0], -v[1], -v[2]})
		for k := 0; k < 3; k++ {
			x[k] += d[k]
		}
	}
	v := q(x)
	c := 1 + math.Pi/4 - 2/math.Pi + a*x[0]*math.Cos(a*x[1])
	fmt.Printf("m=%.17g r=%.17g t=%.17g c=%.17g Q=(%.3g,%.3g,%.3g)\n", x[0], x[1], x[2], c, v[0], v[1], v[2])
	var j [3][3]float64
	for k := 0; k < 3; k++ {
		e := 1e-6
		if k == 0 {
			e = 1e-9
		}
		y, z := x, x
		y[k] += e
		z[k] -= e
		qy, qz := q(y), q(z)
		for i := 0; i < 3; i++ {
			j[i][k] = (qy[i] - qz[i]) / (2 * e)
		}
	}
	fmt.Printf("J=%v\n", j)
	var inv [3][3]float64
	for k := 0; k < 3; k++ {
		var e [3]float64
		e[k] = 1
		col := solve3(j, e)
		for i := 0; i < 3; i++ {
			inv[i][k] = col[i]
		}
	}
	fmt.Printf("R=%v\n", inv)
}
