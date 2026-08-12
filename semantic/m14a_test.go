package semantic

import "testing"

func testSupportOneClass() SupportOneExtremalClass {
	return SupportOneExtremalClass{ID: "ep3.1-zeta", DataRadius: ExactRational{Numerator: 1, Denominator: 1}, TailSemantics: NonpositiveTail, TransformConvention: "hat g(alpha)=integral g(x)e^{-2pi i x alpha}dx", Even: true, Continuous: true, FunctionL1: true, TransformL1: true, FunctionNonnegative: true, Normalization: "g(0)=1", ObjectiveID: "zeta-multiplicity", Provenance: []Reference{{Kind: StandardReference, Citation: "CGdL; Ramos EP3.1"}}}
}

func TestM14AClassDistinguishesSignCutoffFromCompactSupport(t *testing.T) {
	c := testSupportOneClass()
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	if c.TailSemantics == ExactFourierSupport {
		t.Fatal("broader LP collapsed to the M13 bandlimited subclass")
	}
}

func TestM14AMembershipRejectsLateTailAndWrongNormalization(t *testing.T) {
	c := testSupportOneClass()
	m := ExtremalCandidateMembership{CandidateID: "candidate", Even: true, Continuous: true, FunctionL1: true, TransformL1: true, FunctionNonnegative: true, TailSemantics: NonpositiveTail, TailStartsAt: ExactRational{Numerator: 101, Denominator: 100}, ValueAtZero: ExactRational{Numerator: 1, Denominator: 1}, Normalization: "g(0)=1", Proof: "fixture"}
	if err := m.ValidateFor(c); err == nil {
		t.Fatal("support beyond one entered the class")
	}
	m.TailStartsAt = ExactRational{Numerator: 1, Denominator: 1}
	m.ValueAtZero = ExactRational{Numerator: 2, Denominator: 1}
	if err := m.ValidateFor(c); err == nil {
		t.Fatal("wrong normalization entered the class")
	}
}

func TestM14AGridAndMissingTailCannotCertify(t *testing.T) {
	c := testSupportOneClass()
	w := DualCompletionWitness{ID: "grid", MultiplicityLower: ExactRational{Numerator: 9, Denominator: 8}, OutsideMeasure: "sampled", CompletionDistribution: "sampled", FourierImage: "sampled", PositivityEvidence: GridPositivityOnly, ExactSupportScope: true, WholeLineControl: false, TailControl: false, Provenance: Reference{Kind: CompilerRecord, Citation: "fixture"}}
	if _, err := ApplyDualCompletion(c, w); err == nil {
		t.Fatal("grid positivity became theorem evidence")
	}
	w.PositivityEvidence = GlobalAnalyticPD
	w.WholeLineControl = true
	if _, err := ApplyDualCompletion(c, w); err == nil {
		t.Fatal("missing tail control became theorem evidence")
	}
}

func TestM14ABaselineCompletionUsesWeakDuality(t *testing.T) {
	c := testSupportOneClass()
	w := DualCompletionWitness{ID: "sine-baseline", MultiplicityLower: ExactRational{Numerator: 1, Denominator: 1}, OutsideMeasure: "1_{|x|>1} dx", CompletionDistribution: "delta_0-(1-|x|)_+ dx", FourierImage: "1-(sin(pi xi)/(pi xi))^2 >= 0", PositivityEvidence: GlobalAnalyticPD, ExactSupportScope: true, WholeLineControl: true, TailControl: true, Provenance: Reference{Kind: StandardReference, Citation: "Fourier transform of the triangle function"}}
	b, err := ApplyDualCompletion(c, w)
	if err != nil {
		t.Fatal(err)
	}
	if b.SimpleUpper != (ExactRational{Numerator: 1, Denominator: 1}) {
		t.Fatalf("wrong weak-dual conversion: %+v", b)
	}
}
