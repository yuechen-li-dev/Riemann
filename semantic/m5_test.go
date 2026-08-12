package semantic

import (
	"math/cmplx"
	"testing"
)

func m5TestSpan(t *testing.T) FiniteSpan {
	t.Helper()
	parent := WeilNiceClass()
	f1 := TestFunction{Symbol: "f1", Kind: BasisTestFunction, DeclaredClass: parent.ID, RequiredAttributes: append([]FunctionConstraint(nil), parent.Constraints...), TransformConvention: parent.TransformConvention, Parameters: []FunctionParameter{{Name: "index", Value: "1"}}}
	f2 := CloneTestFunction(f1)
	f2.Symbol = "f2"
	f2.Parameters[0].Value = "2"
	basis := OrderedBasis{ID: "B", CoefficientField: ComplexField, ParentClass: parent, Members: []BasisMember{{Function: f1, AdmissibilityCertificate: "f1-admissible"}, {Function: f2, AdmissibilityCertificate: "f2-admissible"}}}
	span := FiniteSpan{ID: "V", Basis: basis, CoefficientField: ComplexField, ParentClass: parent}
	if err := span.Validate(); err != nil {
		t.Fatal(err)
	}
	return span
}

func TestM5OrderedBasisAndSpanHaveSeparateDeterministicIdentity(t *testing.T) {
	span := m5TestSpan(t)
	permuted := CloneFiniteSpan(span)
	permuted.Basis.Members[0], permuted.Basis.Members[1] = permuted.Basis.Members[1], permuted.Basis.Members[0]
	if span.Basis.Key() == permuted.Basis.Key() {
		t.Fatal("ordered-basis identity ignored coordinate order")
	}
	if span.Key() != permuted.Key() {
		t.Fatal("basis permutation changed underlying span identity")
	}
	if span.Key() != CloneFiniteSpan(span).Key() || span.Basis.Key() != CloneFiniteSpan(span).Basis.Key() {
		t.Fatal("span or basis identity is nondeterministic")
	}
}

func TestM5BasisRequiresAdmissibilityForEveryMember(t *testing.T) {
	span := m5TestSpan(t)
	span.Basis.Members[1].AdmissibilityCertificate = ""
	if err := span.Validate(); err == nil {
		t.Fatal("span accepted an uncertified basis member")
	}
}

func TestM5LinearCombinationRetainsSpanFieldAndCoordinates(t *testing.T) {
	span := m5TestSpan(t)
	combination := FiniteLinearCombination{ID: "x", Span: span, Coefficients: CoefficientVector{ID: "c", Field: ComplexField, Entries: []Coefficient{{Symbol: "a"}, {Symbol: "b"}}}}
	if err := combination.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := combination
	bad.Coefficients.Field = RealField
	if err := bad.Validate(); err == nil {
		t.Fatal("coefficient-field mismatch accepted")
	}
}

func TestM5ConjugateFirstPolarizationConvention(t *testing.T) {
	// G=[[1,1+i],[1-i,3]], with B conjugate-linear first.
	q := func(z [2]complex128) float64 {
		g0 := z[0] + (1+1i)*z[1]
		g1 := (1-1i)*z[0] + 3*z[1]
		return real(cmplx.Conj(z[0])*g0 + cmplx.Conj(z[1])*g1)
	}
	polarize := func(x, y [2]complex128) complex128 {
		add := func(a, b [2]complex128, scale complex128) [2]complex128 {
			return [2]complex128{a[0] + scale*b[0], a[1] + scale*b[1]}
		}
		e := q(add(x, y, 1)) - q(add(x, y, -1))
		d := q(add(x, y, 1i)) - q(add(x, y, -1i))
		return complex(e/4, -d/4)
	}
	e1, e2 := [2]complex128{1, 0}, [2]complex128{0, 1}
	if got := polarize(e1, e2); cmplx.Abs(got-(1+1i)) > 1e-12 {
		t.Fatalf("wrong polarization sign/convention: got %v", got)
	}
	if got := polarize(e2, e1); cmplx.Abs(got-cmplx.Conj(polarize(e1, e2))) > 1e-12 {
		t.Fatalf("Hermitian symmetry failed: got %v", got)
	}
	x := [2]complex128{1 - 2i, .5 + 1i}
	if got := polarize(x, x); cmplx.Abs(got-complex(q(x), 0)) > 1e-12 {
		t.Fatalf("quadratic recovery failed: B(x,x)=%v Q(x)=%v", got, q(x))
	}
	convention := ComplexConjugateFirstPolarization("polarization")
	if err := convention.Validate(); err != nil || convention.ConjugateLinearArgument != FirstArgument || convention.LinearArgument != SecondArgument {
		t.Fatalf("typed convention is wrong: %+v, %v", convention, err)
	}
}
