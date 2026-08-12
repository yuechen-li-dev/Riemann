package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM4WeilEquivalenceAndFunctionalStructure(t *testing.T) {
	r, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	contract, ok := r.Registry.Contract(WeilCriterionTheoremID)
	if !ok || contract.Relation != Equivalent || contract.Trust != TrustedExternalTheorem || contract.Conclusion.Kind != semantic.UniversalFunctionalStatementKind {
		t.Fatalf("bad Weil theorem contract: %+v", contract)
	}
	if a := r.Graph.CheckDischarge(ZeroLocationID, WeilPositivityID); !a.Accepted {
		t.Fatalf("RH did not lower equivalently: %+v", a)
	}
	if a := r.Graph.CheckDischarge(WeilPositivityID, ZeroLocationID); !a.Accepted {
		t.Fatalf("Weil equivalence was not reversible: %+v", a)
	}
	claim, _ := r.Graph.Claim(WeilPositivityID)
	p := claim.Proposition.(semantic.UniversalFunctionalStatement)
	if p.FunctionClass.Kind != semantic.WeilNiceFunctionClass || p.Functional != semantic.WeilZetaQuadraticFunctional || p.Predicate != semantic.FunctionalNonnegative || p.TransformConvention != semantic.LagariasMellinConvention {
		t.Fatalf("positivity is not structural: %+v", p)
	}
	if certified, _ := r.Graph.Certify(WeilPositivityID); certified {
		t.Fatal("equivalent open target was spuriously certified")
	}

	def, _ := r.Graph.Claim(WeilFunctionalDefinitionID)
	q := def.Proposition.(semantic.FunctionalDefinition).Functional
	if len(q.Contributions) != 4 {
		t.Fatalf("functional decomposition missing: %+v", q.Contributions)
	}
	want := []semantic.FunctionalContributionKind{semantic.ZeroContribution, semantic.EndpointContribution, semantic.PrimePowerContribution, semantic.ArchimedeanContribution}
	for i, kind := range want {
		if q.Contributions[i].Kind != kind {
			t.Fatalf("contribution %d: got %s want %s", i, q.Contributions[i].Kind, kind)
		}
	}
	agg := q.Contributions[0].Aggregate
	if agg == nil || agg.IndexDomain != semantic.NontrivialZeros(semantic.RiemannZeta) || agg.TransformConvention != semantic.LagariasMellinConvention || !containsTheorem(agg.TheoremLineage, WeilExplicitFormulaTheoremID) || agg.Provenance.Citation == "" {
		t.Fatalf("aggregate provenance incomplete: %+v", agg)
	}
}

func TestM4FiniteRestrictionDirectionLossAndAdmissibility(t *testing.T) {
	r, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	if !r.FullToFinite.Accepted {
		t.Fatalf("valid universal restriction rejected: %+v", r.FullToFinite)
	}
	if r.FiniteToRH.Accepted || !hasCode(r.FiniteToRH.Diagnostics, InformationLost) || !hasCode(r.FiniteToRH.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("finite family discharged RH or hid loss: %+v", r.FiniteToRH)
	}
	var found bool
	for _, tr := range r.Graph.Transformations() {
		if tr.ID == RestrictWeilFamilyID {
			found = len(tr.Losses) == 1 && tr.Losses[0].Kind == FunctionSpaceRestriction && len(tr.Premises) == 2
		}
	}
	if !found {
		t.Fatal("finite restriction lacks coverage loss or admissibility premises")
	}

	g := NewGraph()
	full := semantic.WeilNiceClass()
	source := authoredMathClaim("all", functionalStatement(full), semantic.KnownTheoremEvidence, testReference)
	mustAddClaim(t, g, source)
	f := abstractWeilFunction("uncertified", 9, full)
	finite := semantic.FunctionClass{ID: "bad-family", Kind: semantic.FiniteFunctionClass, Constraints: full.Constraints, TransformConvention: full.TransformConvention, Members: []semantic.TestFunction{f}}
	if _, err := (FunctionClassRestriction{TargetID: "bad", Class: finite, TransformationID: "bad-restrict"}).Apply(g, source.ID); err == nil || !strings.Contains(err.Error(), "lacks certified admissibility") {
		t.Fatalf("missing admissibility silently accepted: %v", err)
	}
	unsound := semantic.Claim{ID: "unsound-finite", Proposition: functionalStatement(finite), Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{source.ID}, Transformation: "direct-unsound", Source: testReference}}
	mustAddClaim(t, g, unsound)
	if err := g.AddTransformation(Transformation{ID: "direct-unsound", Pass: "test", From: source.ID, To: unsound.ID, Relation: Implies, Losses: []InformationLoss{{Kind: FunctionSpaceRestriction, Reason: "test"}}, Provenance: testReference}); err == nil || !strings.Contains(err.Error(), "admissibility premise") {
		t.Fatalf("direct graph edge bypassed admissibility: %v", err)
	}
}

func TestM4PredicateFunctionalTransformAndNumericalBoundaries(t *testing.T) {
	r, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	full := semantic.WeilNiceClass()
	base := exactClaim("base", functionalStatement(full))
	pattern := functionalPattern(full)
	if !matchPattern(pattern, base, make(bindingEnv)) {
		t.Fatal("correct functional premise did not match")
	}
	wrongFunctional := base
	p := wrongFunctional.Proposition.(semantic.UniversalFunctionalStatement)
	p.Functional = "unrelated"
	wrongFunctional.Proposition = p
	if matchPattern(pattern, wrongFunctional, make(bindingEnv)) {
		t.Fatal("wrong functional matched Weil premise")
	}
	wrongPredicate := base
	p = wrongPredicate.Proposition.(semantic.UniversalFunctionalStatement)
	p.Predicate = "strictly_positive"
	wrongPredicate.Proposition = p
	if matchPattern(pattern, wrongPredicate, make(bindingEnv)) {
		t.Fatal("unrelated predicate matched")
	}
	wrongTransform := base
	p = wrongTransform.Proposition.(semantic.UniversalFunctionalStatement)
	p.FunctionClass.TransformConvention = "other"
	p.TransformConvention = "other"
	wrongTransform.Proposition = p
	if matchPattern(pattern, wrongTransform, make(bindingEnv)) {
		t.Fatal("incompatible transform convention matched")
	}
	if r.NumericalToUniversal.Accepted || !hasCode(r.NumericalToUniversal.Diagnostics, NoEstablishedDirection) || !hasCode(r.NumericalToUniversal.Diagnostics, ApproximationBoundary) || !hasCode(r.NumericalToUniversal.Diagnostics, DomainMismatch) {
		t.Fatalf("numerical finite evidence crossed theorem boundary: %+v", r.NumericalToUniversal)
	}
}

func TestM4ReportsDeterministicAndM0ThroughM3Regress(t *testing.T) {
	r, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := M4JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	wantHuman := M4HumanReport(r)
	for i := 0; i < 3; i++ {
		got, err := CompileM4()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := M4JSONReport(got)
		if !bytes.Equal(data, wantJSON) || M4HumanReport(got) != wantHuman {
			t.Fatal("M4 output is nondeterministic")
		}
	}
	for _, text := range []string{"RIEMANN-M4", "Weil", "function_space_restriction", "RH CERTIFICATION\n  REJECTED", "numerical only"} {
		if !strings.Contains(wantHuman, text) {
			t.Fatalf("human report omits %q", text)
		}
	}
	for _, text := range []string{`"schema": "riemann.semantic-graph.m4"`, `"universal_functional_statement"`, `"prime_power"`, `"transform_convention"`} {
		if !strings.Contains(string(wantJSON), text) {
			t.Fatalf("JSON omits %q", text)
		}
	}
	if !r.M3.M1.ZeroFreeCertified || !r.M3.StripCertified || !r.M3.SymmetryCertified {
		t.Fatal("M0-M3 soundness regression")
	}
}

func containsTheorem(in []semantic.TheoremID, want semantic.TheoremID) bool {
	for _, id := range in {
		if id == want {
			return true
		}
	}
	return false
}
