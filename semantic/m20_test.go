package semantic

import (
	"math"
	"testing"
)

func TestM20ZeroPerturbationRegressesExactlyToM18FloatTransform(t *testing.T) {
	v, err := ContactAdaptedAtom(0, 4)
	if err != nil {
		t.Fatal(err)
	}
	if v.CIncrease != 0 || v.BoundaryWeightChange != 0 {
		t.Fatalf("zero perturbation changed parameters: %+v", v)
	}
	for _, x := range []float64{0, .2, M20Contact, 4, 17} {
		if got, want := AtomExtendedSpectrum(v, 4, x), M20BaselineSpectrum(x); got != want {
			t.Fatalf("M18 regression failed at %g: %.17g != %.17g", x, got, want)
		}
		interval, err := ContactAdaptedInterval(0, 4, .25)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := IntervalExtendedSpectrum(interval, 4, .25, x), M20BaselineSpectrum(x); got != want {
			t.Fatalf("density M18 regression failed at %g", x)
		}
		packet, err := ContactAdaptedLatticePacket(0)
		if err != nil {
			t.Fatal(err)
		}
		if got, want := LatticePacketExtendedSpectrum(packet, x), M20BaselineSpectrum(x); got != want {
			t.Fatalf("lattice M18 regression failed at %g", x)
		}
	}
	if DeriveOneRadiusTangencyCandidate().ExactCExpression.Expression != "1+pi/4-2/pi" {
		t.Fatal("M18 exact optimum drifted")
	}
}

func TestM20FirstVariationMatchesFiniteDifferenceAndPreservesContact(t *testing.T) {
	epsilon := 1e-7
	v, err := ContactAdaptedAtom(epsilon, 4)
	if err != nil {
		t.Fatal(err)
	}
	if !v.PreservesValue || !v.PreservesDerivative || v.CIncrease <= 0 {
		t.Fatalf("bad contact variation: %+v", v)
	}
	if math.Abs(v.CIncrease/epsilon-M20Contact) > 1e-12 || math.Abs(v.BoundaryWeightChange/epsilon-2/math.Pi) > 1e-12 {
		t.Fatalf("wrong r=4 direction: %+v", v)
	}
	if math.Abs(AtomExtendedSpectrum(v, 4, M20Contact)) > 2e-14 {
		t.Fatal("finite perturbation lost exact linear contact")
	}
	h := 1e-5
	d := (AtomExtendedSpectrum(v, 4, M20Contact+h) - AtomExtendedSpectrum(v, 4, M20Contact-h)) / (2 * h)
	if math.Abs(d) > 2e-10 {
		t.Fatalf("retuned derivative is %g", d)
	}
}

func TestM20ContactShiftPredictionMatchesUnretunedNumerics(t *testing.T) {
	mass := 1e-6
	deltaPrime := -2 * mass * 3 * math.Sin(3*M20Contact)
	predicted := PredictContactShift(deltaPrime)
	// Newton's first step for the perturbed derivative is the linearized shift.
	numerical := -deltaPrime / M20BaselineCurvature()
	if math.Abs(predicted-numerical) > 1e-15 || predicted >= 0 {
		t.Fatalf("wrong contact shift: %g", predicted)
	}
}

func TestM20CandidateLegalityPrimalGuardAndEvidenceBoundary(t *testing.T) {
	if _, err := ContactAdaptedInterval(.001, 1.1, .2); err == nil {
		t.Fatal("density crossing support boundary accepted")
	}
	if _, err := ContactAdaptedInterval(.001, 4, .25); err != nil {
		t.Fatal(err)
	}
	if _, err := ContactAdaptedLatticePacket(-.001); err == nil {
		t.Fatal("negative lattice mass accepted")
	}
	if err := (M20Probe{Family: "probe", CandidateC: 1.2, Evidence: NumericalEvidence, TheoremClaim: true}).Validate(); err == nil {
		t.Fatal("numerical theorem claim accepted")
	}
	if err := (M20Probe{Family: "probe", CandidateC: CompatiblePrimalC, Evidence: NumericalEvidence}).Validate(); err == nil {
		t.Fatal("primal guard failed")
	}
	if err := (M20Probe{Family: "probe", CandidateC: 1.2, Evidence: NumericalEvidence}).Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestM20UtilityIsDeterministicResearchMetadata(t *testing.T) {
	u := RepresentationUtility{Family: "contact-adapted", Improvement: 5, Tractability: 5, WholeLineCost: 4, ExactPotential: 5, Reuse: 4, Implementation: 5, Clarity: 5}
	if u.Score() != 52 || CertifiedEvidence == EvidenceClass(u.Family) {
		t.Fatalf("utility leaked into evidence semantics: %+v score=%d", u, u.Score())
	}
}
