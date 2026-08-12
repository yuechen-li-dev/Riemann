package semantic

import "testing"

func TestPointTransformStructuralIdentityAndComposition(t *testing.T) {
	s := Point("s")
	points := []PointExpr{s, s.Apply(ConjugateTransform), s.Apply(OneMinusTransform), s.Apply(CriticalReflectionTransform)}
	seen := map[string]bool{}
	for _, p := range points {
		if seen[p.Key()] {
			t.Fatalf("unconstrained transforms collided at %s", p.Describe())
		}
		seen[p.Key()] = true
	}
	if got := s.Apply(ConjugateTransform).Apply(OneMinusTransform); got != s.Apply(CriticalReflectionTransform) {
		t.Fatalf("composition did not canonicalize: %+v", got)
	}
	if got := s.Apply(OneMinusTransform).Apply(OneMinusTransform); got != s {
		t.Fatalf("reflection not involutive: %+v", got)
	}
	if Compose(ConjugateTransform, OneMinusTransform) != CriticalReflectionTransform || Compose(OneMinusTransform, ConjugateTransform) != CriticalReflectionTransform {
		t.Fatal("commuting transform composition is wrong")
	}
}

func TestCriticalLinePointCanonicalDegeneracy(t *testing.T) {
	rho := PointOnCriticalLine("ρ")
	if rho.Apply(CriticalReflectionTransform) != rho.Canonical() {
		t.Fatal("critical-line reflection did not fix point")
	}
	if rho.Apply(OneMinusTransform) != rho.Apply(ConjugateTransform) {
		t.Fatal("critical-line equivalent transforms did not canonicalize")
	}
}

func TestCompletedXiStructurallyReferencesZeta(t *testing.T) {
	good := Representation{Object: RiemannXi, BaseObject: RiemannZeta, Name: CompletedXiRepresentation, ValidOn: ComplexPlane(), Formula: "xi formula"}
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := good
	bad.BaseObject = 0
	if err := bad.Validate(); err == nil {
		t.Fatal("completed object lost structural relationship to zeta")
	}
}
