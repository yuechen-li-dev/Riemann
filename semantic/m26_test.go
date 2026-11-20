package semantic

import (
	"math"
	"strconv"
	"testing"
)

func TestM26ContactPreservingFamilyEmbedsExactly(t *testing.T) {
	identity := VerifyContactPreservingEmbedding()
	if !identity.Verified || identity.CollectedResponse == "" {
		t.Fatalf("missing symbolic embedding certificate: %+v", identity)
	}
	for _, r := range []float64{1.25, 2, 3.5, 4, 7.25} {
		for _, x := range []float64{0, 0.4, math.Pi / 2, 1.8, 4, 9} {
			m := 0.0069
			c0 := 1 + math.Pi/4 - 2/math.Pi
			w0 := 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi)
			a, s := math.Cos(math.Pi*r/2), math.Sin(math.Pi*r/2)
			general := UnrestrictedTwoAtomSpectrum(UnrestrictedTwoAtomFamily{c0 + math.Pi*m*a/2, w0 + m*(2*a/math.Pi-r*s), m, r, UnrestrictedTwoAtomTransform}, x)
			adapted := OneRadiusSpectrum(c0, w0, x) + m*TwoAtomResponseFloat(r, x)
			if math.Abs(general-adapted) > 2e-14 {
				t.Fatalf("embedding changed at r=%g t=%g: general=%g adapted=%g", r, x, general, adapted)
			}
		}
	}
}

func TestM26ReleasedWitnessStrictlyImprovesM25AndIsWholeLine(t *testing.T) {
	w, err := CertifyM26ReleasedWitness()
	if err != nil {
		t.Fatal(err)
	}
	lo, _ := strconv.ParseFloat(w.C.Lower, 64)
	if !w.WholeLine || lo <= 1.1566433083347464 {
		t.Fatalf("M26 did not certify strict improvement: %+v", w)
	}
	fixed, _ := strconv.ParseFloat(w.FixedContactValue.Lower, 64)
	movable, _ := strconv.ParseFloat(w.OldMovableValue.Lower, 64)
	if fixed <= 0 || movable <= 0 {
		t.Fatalf("old contacts were not released: fixed=%g movable=%g", fixed, movable)
	}
}

func TestM26DoesNotFreeFirstRadiusOrAddAtoms(t *testing.T) {
	if UnrestrictedTwoAtomTransform == "" {
		t.Fatal("missing M26 representation")
	}
	f := UnrestrictedTwoAtomFamily{C: 1.1, Weight1: 0.1, Weight2: 0.01, Radius: 3.5, Transform: UnrestrictedTwoAtomTransform}
	if err := f.Validate(); err != nil {
		t.Fatal(err)
	}
}
