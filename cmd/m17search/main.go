// Command m17search performs the purpose-specific reduced numerical search for
// the M17 unsaturated one-radius family. Its output is research evidence only.
package main

import (
	"flag"
	"fmt"
	"math"
)

type candidate struct{ c, r, w, contact, minimum, minimumT float64 }

func spectrum(c, r, w, t float64) float64 {
	if math.Abs(t) < 1e-8 {
		return 2 * (1 + w - c)
	}
	return 1 + 2*(math.Cos(t)-1)/(t*t) + 2*(1-c)*math.Sin(t)/t + 2*w*math.Cos(r*t)
}

// ceiling eliminates c analytically: when sin(t)>0, F(t)>=0 is an
// upper bound on c. The origin contributes c<=1+w.
func ceiling(r, w, tStep, tMax float64) (float64, float64) {
	c, contact := 1+w, 0.0
	for t := tStep; t <= tMax; t += tStep {
		s := math.Sin(t)
		if s <= 1e-10 {
			continue
		}
		rhs := 1 + 2*(math.Cos(t)-1)/(t*t) + 2*s/t + 2*w*math.Cos(r*t)
		bound := rhs * t / (2 * s)
		if bound < c {
			c, contact = bound, t
		}
	}
	return c, contact
}

func assess(r, w, tStep, tMax float64) candidate {
	c, contact := ceiling(r, w, tStep, tMax)
	q := candidate{c: c, r: r, w: w, contact: contact, minimum: spectrum(c, r, w, 0)}
	for t := tStep; t <= tMax; t += tStep {
		v := spectrum(c, r, w, t)
		if v < q.minimum {
			q.minimum, q.minimumT = v, t
		}
	}
	return q
}

func main() {
	candidateC := flag.Float64("c", 0, "assess this c instead of optimizing")
	candidateR := flag.Float64("r", 1, "candidate radius")
	candidateW := flag.Float64("w", 0, "candidate weight")
	rMin := flag.Float64("r-min", 1, "minimum radius")
	rMax := flag.Float64("r-max", 12, "maximum radius")
	rStep := flag.Float64("r-step", .02, "radius step")
	tStep := flag.Float64("t-step", .01, "frequency step")
	tMax := flag.Float64("t-max", 80, "frequency maximum")
	flag.Parse()
	if *candidateC > 0 {
		q := candidate{c: *candidateC, r: *candidateR, w: *candidateW, minimum: spectrum(*candidateC, *candidateR, *candidateW, 0)}
		for t := *tStep; t <= *tMax; t += *tStep {
			v := spectrum(q.c, q.r, q.w, t)
			if v < q.minimum {
				q.minimum, q.minimumT = v, t
			}
		}
		fmt.Printf("candidate c=%.12f r=%.12f w=%.12f sampled_min=%.12g at t=%.12f\n", q.c, q.r, q.w, q.minimum, q.minimumT)
		return
	}
	best := candidate{c: math.Inf(-1)}
	for r := *rMin; r <= *rMax+*rStep/2; r += *rStep {
		// The fixed-r envelope is concave in w. Golden-section search on the
		// necessary compact interval 0<=w<=1/2, then check every frequency.
		lo, hi := 0.0, .5
		for i := 0; i < 45; i++ {
			m1, m2 := lo+(hi-lo)*.3819660112501051, lo+(hi-lo)*.6180339887498949
			c1, _ := ceiling(r, m1, *tStep, *tMax)
			c2, _ := ceiling(r, m2, *tStep, *tMax)
			if c1 < c2 {
				lo = m1
			} else {
				hi = m2
			}
		}
		q := assess(r, (lo+hi)/2, *tStep, *tMax)
		if q.minimum >= -1e-8 && q.c > best.c {
			best = q
		}
	}
	fmt.Printf("best c=%.12f r=%.12f w=%.12f contact=%.12f sampled_min=%.12g at t=%.12f\n", best.c, best.r, best.w, best.contact, best.minimum, best.minimumT)
}
