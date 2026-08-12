package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM3ConjugationAndFunctionalReflectionTransport(t *testing.T) {
	r := mustM3(t, M3Options{})
	for _, point := range []semantic.PointExpr{semantic.Point("ρ").Apply(semantic.ConjugateTransform), semantic.Point("ρ").Apply(semantic.OneMinusTransform), semantic.Point("ρ").Apply(semantic.CriticalReflectionTransform)} {
		claim, ok := zeroClaimAt(r.Graph, semantic.RiemannZeta, point)
		if !ok {
			t.Fatalf("missing transported zero at %s", point.Describe())
		}
		if claim.Provenance.Theorem == "" || len(claim.Provenance.Parents) == 0 {
			t.Fatalf("transport lost theorem provenance: %+v", claim.Provenance)
		}
	}
	if !hasCompleteApplication(r.Applications, ZetaConjugationTheoremID) || !hasCompleteApplication(r.Applications, FunctionalIdentityZeroTransportTheoremID) {
		t.Fatal("zero claims did not arise through theorem applications")
	}
}

func TestM3SideConditionObligationAndInvalidTransport(t *testing.T) {
	missing := mustM3(t, M3Options{OmitCompletionFactor: true})
	if _, ok := zeroClaimAt(missing.Graph, semantic.RiemannZeta, semantic.Point("ρ").Apply(semantic.OneMinusTransform)); ok {
		t.Fatal("functional reflection crossed a missing completion-factor condition")
	}
	if !hasSideObligation(missing.Applications, LiftZeroToXiTheoremID, semantic.CompletionFactorRegularNonzero) {
		t.Fatal("missing completion factor was not an explicit typed obligation")
	}

	boundary := mustM3(t, M3Options{OmitAnalyticContinuation: true})
	if _, ok := zeroClaimAt(boundary.Graph, semantic.RiemannZeta, semantic.Point("ρ").Apply(semantic.ConjugateTransform)); ok {
		t.Fatal("Dirichlet-series facts enabled global conjugation transport")
	}
	if !hasSideObligation(boundary.Applications, ZetaConjugationTheoremID, semantic.AnalyticContinuationAvailable) {
		t.Fatal("analytic continuation boundary did not remain explicit")
	}
	dirichlet, _ := boundary.Graph.Claim(DirichletRepresentationID)
	rep := dirichlet.Proposition.(semantic.RepresentationProposition).Representation
	if rep.ValidOn != semantic.HalfPlaneReGreaterThanOne() {
		t.Fatalf("Dirichlet representation escaped its domain: %+v", rep.ValidOn)
	}
}

func TestM3SymmetryClosureAndCriticalLineDegeneracy(t *testing.T) {
	r := mustM3(t, M3Options{})
	if len(r.Orbit.Generated) != 4 || len(r.Orbit.Distinct) != 4 {
		t.Fatalf("generic orbit mismatch: %+v", r.Orbit)
	}
	line := mustM3(t, M3Options{SamplePoint: semantic.PointOnCriticalLine("ρ")})
	if len(line.Orbit.Generated) != 4 || len(line.Orbit.Distinct) != 2 {
		t.Fatalf("critical-line orbit reported a false quartet: %+v", line.Orbit)
	}
}

func TestM3GlobalGeometryDerivationAndRHFixedPointEquivalence(t *testing.T) {
	r := mustM3(t, M3Options{})
	if !r.StripCertified || !r.SymmetryCertified {
		t.Fatalf("global geometry not certified: strip=%v %+v symmetry=%v %+v", r.StripCertified, r.StripDiagnostics, r.SymmetryCertified, r.SymmetryDiagnostics)
	}
	lineage, err := r.Graph.Lineage(CriticalStripConfinementID)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []semantic.ClaimID{ZeroFreeHalfPlaneID, FunctionalInvariantID, ZetaTrivialZerosID, ZetaBoundaryZeroFreeID} {
		if !containsID(lineage, id) {
			t.Fatalf("strip derivation omits %s: %v", id, lineage)
		}
	}
	fixed, _ := r.Graph.Claim(RHFixedPointID)
	q, ok := semantic.Quantified(fixed.Proposition)
	if !ok || q.Predicate.Kind != semantic.CriticalReflectionFixedPredicate {
		t.Fatalf("RH fixed-point target is not structural: %+v", fixed.Proposition)
	}
	attempt := r.Graph.CheckDischarge(ZeroLocationID, RHFixedPointID)
	if !attempt.Accepted {
		t.Fatalf("registered exact equivalence unavailable: %+v", attempt)
	}
	if certified, _ := r.Graph.Certify(RHFixedPointID); certified {
		t.Fatal("equivalent normalization accidentally proved RH")
	}
}

func TestM3TrustProvenanceAndReportsAreDeterministic(t *testing.T) {
	r := mustM3(t, M3Options{})
	claim, _ := zeroClaimAt(r.Graph, semantic.RiemannZeta, semantic.Point("ρ").Apply(semantic.CriticalReflectionTransform))
	if !assumptionsContain(claim.Assumptions, []semantic.Assumption{{ID: "rho-is-nontrivial-zero"}}) {
		t.Fatalf("transport dropped source hypothesis: %+v", claim.Assumptions)
	}
	lineage, err := r.Graph.Lineage(claim.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []semantic.ClaimID{SampleZeroID, CompletedXiID, XiFunctionalEquationID} {
		if !containsID(lineage, id) {
			t.Fatalf("symmetry provenance omits %s: %v", id, lineage)
		}
	}
	wantJSON, err := M3JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	wantHuman := M3HumanReport(r)
	for i := 0; i < 3; i++ {
		again := mustM3(t, M3Options{})
		got, _ := M3JSONReport(again)
		if !bytes.Equal(got, wantJSON) || M3HumanReport(again) != wantHuman {
			t.Fatal("M3 reports are nondeterministic")
		}
	}
	for _, text := range []string{"riemann.semantic-graph.m3", "off-axis symmetry orbits", "four generated symmetry transforms", "distinct zero locations"} {
		if !strings.Contains(string(wantJSON)+wantHuman, text) {
			t.Fatalf("report omits %q", text)
		}
	}
}

func mustM3(t *testing.T, o M3Options) M3Result {
	t.Helper()
	r, err := CompileM3WithOptions(o)
	if err != nil {
		t.Fatal(err)
	}
	return r
}
func zeroClaimAt(g *Graph, f semantic.Function, p semantic.PointExpr) (semantic.Claim, bool) {
	return g.ClaimBySemanticKey(semantic.SemanticKey(semantic.ZeroAtPoint{Object: f, Point: p.Canonical(), Classification: semantic.NontrivialZero}))
}
func hasCompleteApplication(apps []TheoremApplication, id semantic.TheoremID) bool {
	for _, a := range apps {
		if a.Theorem == id && a.Complete {
			return true
		}
	}
	return false
}
func hasSideObligation(apps []TheoremApplication, id semantic.TheoremID, name semantic.SideConditionName) bool {
	for _, a := range apps {
		if a.Theorem != id || a.Complete {
			continue
		}
		for _, o := range a.Obligations {
			if !o.SideCondition {
				continue
			}
			if s, ok := o.Proposition.(semantic.SideCondition); ok && s.Condition == name {
				return true
			}
		}
	}
	return false
}
