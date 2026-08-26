// Command m22search performs the deliberately bounded numerical search used by
// M22.  Its floating-point output is research evidence, never theorem input.
package main

import (
	"fmt"
	"math"
)

const (
	pi = math.Pi
	a  = pi / 2
)

func f0(t float64) float64 {
	c0 := 1 + pi/4 - 2/pi
	w0 := 1/pi - 4/(pi*pi) + 8/(pi*pi*pi)
	if math.Abs(t) < 1e-9 {
		return 2 * (1 + w0 - c0)
	}
	return 1 + 2*(math.Cos(t)-1)/(t*t) + 2*(1-c0)*math.Sin(t)/t + 2*w0*math.Cos(t)
}

func h(r, t float64) float64 {
	c, s := math.Cos(a*r), math.Sin(a*r)
	sinc := 1.0
	if math.Abs(t) >= 1e-9 {
		sinc = math.Sin(t) / t
	}
	return -pi*c*sinc + 2*(2*c/pi-r*s)*math.Cos(t) + 2*math.Cos(r*t)
}

func hPrime(r, t float64) float64 {
	c, s := math.Cos(a*r), math.Sin(a*r)
	return -pi*c*(t*math.Cos(t)-math.Sin(t))/(t*t) - 2*(2*c/pi-r*s)*math.Sin(t) - 2*r*math.Sin(r*t)
}

func hRadius(r, t float64) float64 {
	c, s := math.Cos(a*r), math.Sin(a*r)
	cp := -a * s
	bp := -2*s - a*r*c
	return -pi*cp*math.Sin(t)/t + 2*bp*math.Cos(t) - 2*t*math.Sin(r*t)
}

func hPrimeRadius(r, t float64) float64 {
	c, s := math.Cos(a*r), math.Sin(a*r)
	cp := -a * s
	bp := -2*s - a*r*c
	return -pi*cp*(t*math.Cos(t)-math.Sin(t))/(t*t) - 2*bp*math.Sin(t) - 2*math.Sin(r*t) - 2*r*t*math.Cos(r*t)
}

func envelope(r, t float64) float64 {
	hv := h(r, t)
	if hv >= -1e-14 {
		return math.Inf(1)
	}
	return f0(t) / -hv
}

func golden(r, left, right float64) (float64, float64) {
	const q = 0.6180339887498948482
	x1, x2 := right-q*(right-left), left+q*(right-left)
	y1, y2 := envelope(r, x1), envelope(r, x2)
	for i := 0; i < 100; i++ {
		if y1 < y2 {
			right, x2, y2 = x2, x1, y1
			x1 = right - q*(right-left)
			y1 = envelope(r, x1)
		} else {
			left, x1, y1 = x1, x2, y2
			x2 = left + q*(right-left)
			y2 = envelope(r, x2)
		}
	}
	t := (left + right) / 2
	return t, envelope(r, t)
}

func bestForRadius(r float64) (float64, float64, float64) {
	bestT, bestM := 0.0, math.Inf(1)
	const limit, n = 80.0, 160000
	step := limit / n
	prev := envelope(r, 1e-6)
	cur := envelope(r, step)
	for i := 2; i <= n; i++ {
		next := envelope(r, float64(i)*step)
		if cur <= prev && cur <= next && cur < bestM {
			t, m := golden(r, float64(i-2)*step, float64(i)*step)
			if m < bestM {
				bestT, bestM = t, m
			}
		}
		prev, cur = cur, next
	}
	c := 1 + pi/4 - 2/pi + a*math.Cos(a*r)*bestM
	return bestT, bestM, c
}

func main() {
	bestR, bestT, bestM, bestC := 0.0, 0.0, 0.0, -math.MaxFloat64
	for i := 0; i <= 2000; i++ {
		r := 3.0 + 2.0*float64(i)/2000
		t, m, c := bestForRadius(r)
		if c > bestC {
			bestR, bestT, bestM, bestC = r, t, m, c
		}
	}
	fmt.Printf("coarse r=%.15g t=%.15g m=%.15g c=%.15g\n", bestR, bestT, bestM, bestC)
	left, right := bestR-.002, bestR+.002
	const q = 0.6180339887498948482
	for i := 0; i < 60; i++ {
		r1, r2 := right-q*(right-left), left+q*(right-left)
		_, _, c1 := bestForRadius(r1)
		_, _, c2 := bestForRadius(r2)
		if c1 > c2 {
			right = r2
		} else {
			left = r1
		}
	}
	bestR = (left + right) / 2
	bestT, bestM, bestC = bestForRadius(bestR)
	c, s := math.Cos(a*bestR), math.Sin(a*bestR)
	w0 := 1/pi - 4/(pi*pi) + 8/(pi*pi*pi)
	w1 := w0 + bestM*(2*c/pi-bestR*s)
	fmt.Printf("refined r=%.17g t=%.17g m=%.17g w1=%.17g c=%.17g curvature=%.17g\n", bestR, bestT, bestM, w1, bestC, func() float64 {
		e := 1e-5
		return (f0(bestT+e)+bestM*h(bestR, bestT+e)-2*(f0(bestT)+bestM*h(bestR, bestT))+f0(bestT-e)+bestM*h(bestR, bestT-e))/(e*e)
	}())
	prevV := f0(0) + bestM*h(bestR, 0)
	curV := f0(.001) + bestM*h(bestR, .001)
	fmt.Print("secondary minima:")
	printed := 0
	for i := 2; i <= 80000 && printed < 8; i++ {
		t := float64(i) / 1000
		nextV := f0(t) + bestM*h(bestR, t)
		if curV <= prevV && curV <= nextV && t-.001 > 2.0 {
			fmt.Printf(" (t=%.3f,F=%.9g)", t-.001, curV)
			printed++
		}
		prevV, curV = curV, nextV
	}
	fmt.Println()
	m21r, m21t, m21m := 4.0, 1.36474957003247, 0.00387226734364153
	dmDr := -m21m * hRadius(m21r, m21t) / h(m21r, m21t)
	curvatureStep := 1e-5
	familySecond := func(r, m, t float64) float64 {
		return (f0(t+curvatureStep) + m*h(r, t+curvatureStep) - 2*(f0(t)+m*h(r, t)) + f0(t-curvatureStep) + m*h(r, t-curvatureStep)) / (curvatureStep * curvatureStep)
	}
	dtDr := -(hPrime(m21r, m21t)*dmDr + m21m*hPrimeRadius(m21r, m21t)) / familySecond(m21r, m21m, m21t)
	fmt.Printf("M21 sensitivity dc/dr=%.17g dw1/dr=%.17g dF/dr=%.17g dFprime/dr=%.17g dm*/dr=%.17g dt1/dr=%.17g\n",
		0.0, -2*pi*m21m, m21m*hRadius(m21r, m21t), m21m*hPrimeRadius(m21r, m21t), dmDr, dtDr)
	for _, r := range []float64{2, 3, 3.5, 3.75, 3.875, 4, 4.125, 4.25, 5} {
		t, m, c := bestForRadius(r)
		fmt.Printf("r=%g t=%.12g m=%.12g c=%.15g\n", r, t, m, c)
	}
}
