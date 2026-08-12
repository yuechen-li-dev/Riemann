package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

var testReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "test"}

func authoredClaim(id semantic.ClaimID, evidence semantic.EvidenceKind) semantic.Claim {
	return semantic.Claim{ID: id, Proposition: semantic.NamedObligation{Name: string(id)}, Evidence: []semantic.Evidence{{Kind: evidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: testReference}}
}

func derivedClaim(id semantic.ClaimID, parents []semantic.ClaimID, transformation semantic.TransformationID) semantic.Claim {
	return semantic.Claim{ID: id, Proposition: semantic.NamedObligation{Name: string(id)}, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: transformation, Source: testReference}}
}

func addPair(t *testing.T, relation Relation) (*Graph, semantic.ClaimID, semantic.ClaimID) {
	t.Helper()
	g := NewGraph()
	from, to := semantic.ClaimID("a"), semantic.ClaimID("b")
	transform := semantic.TransformationID("a-to-b")
	mustAddClaim(t, g, authoredClaim(from, semantic.KnownTheoremEvidence))
	mustAddClaim(t, g, derivedClaim(to, []semantic.ClaimID{from}, transform))
	losses := []InformationLoss(nil)
	if relation == Approximation {
		losses = []InformationLoss{{Kind: ApproximationLoss, Reason: "test"}}
	}
	if err := g.AddTransformation(Transformation{ID: transform, Pass: "test", From: from, To: to, Relation: relation, Losses: losses}); err != nil {
		t.Fatal(err)
	}
	return g, from, to
}

func TestDirectionEvidenceAssumptionsAndApproximationRemainSound(t *testing.T) {
	g, a, b := addPair(t, Equivalent)
	if !g.CheckDischarge(a, b).Accepted || !g.CheckDischarge(b, a).Accepted {
		t.Fatal("equivalence was not bidirectional")
	}
	g, a, b = addPair(t, Implies)
	if !g.CheckDischarge(a, b).Accepted || !hasCode(g.CheckDischarge(b, a).Diagnostics, NoEstablishedDirection) {
		t.Fatal("implication direction was not enforced")
	}
	g, _, b = addPair(t, Approximation)
	if certified, diagnostics := g.Certify(b); certified || !hasCode(diagnostics, ApproximationBoundary) {
		t.Fatalf("approximation certified exact claim: %v %+v", certified, diagnostics)
	}
	numerical := authoredClaim("numerical", semantic.NumericalExperimentEvidence)
	mustAddClaim(t, g, numerical)
	if g.AttemptProof(numerical.ID, numerical.ID).Accepted {
		t.Fatal("numerical evidence certified an exact claim")
	}

	g = NewGraph()
	source := authoredClaim("conditional", semantic.KnownTheoremEvidence)
	source.Assumptions = []semantic.Assumption{{ID: "h", Description: "hypothesis"}}
	target := derivedClaim("unconditional", []semantic.ClaimID{source.ID}, "drop")
	mustAddClaim(t, g, source)
	mustAddClaim(t, g, target)
	if err := g.AddTransformation(Transformation{ID: "drop", Pass: "test", From: source.ID, To: target.ID, Relation: Implies}); err == nil {
		t.Fatal("assumption was silently dropped")
	}
}

func TestTypedDomainAlgebra(t *testing.T) {
	bounded := semantic.ZerosBelowHeight(semantic.RiemannZeta, 100)
	all := semantic.NontrivialZeros(semantic.RiemannZeta)
	if !semantic.IsSubset(bounded, all) || semantic.IsSubset(all, bounded) {
		t.Fatal("bounded-zero inclusion has wrong direction")
	}
	if !semantic.IsSubset(semantic.CriticalStrip(), semantic.ComplexPlane()) {
		t.Fatal("geometric inclusion missing")
	}
	if !semantic.IsSubset(semantic.CriticalLine(), semantic.CriticalStrip()) || !semantic.IsSubset(all, semantic.CriticalStrip()) {
		t.Fatal("modeled critical-strip inclusions are missing")
	}
	if semantic.IsSubset(semantic.HalfPlaneReGreaterThanOne(), semantic.CriticalStrip()) {
		t.Fatal("disjoint geometric domains treated as included")
	}
}

func TestUniversalDomainRestrictionAndReverseRejection(t *testing.T) {
	g := NewGraph()
	source := semantic.Claim{ID: "all-zeros", Proposition: criticalLineStatement(semantic.ForAll, semantic.NontrivialZeros(semantic.RiemannZeta)), Evidence: []semantic.Evidence{{Kind: semantic.KnownTheoremEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: testReference}}
	mustAddClaim(t, g, source)
	targetID, err := (DomainRestriction{TargetID: "bounded", Domain: semantic.ZerosBelowHeight(semantic.RiemannZeta, 100), TransformationID: "restrict"}).Apply(g, source.ID)
	if err != nil {
		t.Fatal(err)
	}
	if !g.AttemptProof(source.ID, targetID).Accepted {
		t.Fatalf("forward restriction rejected: %+v", g.AttemptProof(source.ID, targetID))
	}
	reverse := g.CheckDischarge(targetID, source.ID)
	if reverse.Accepted || !hasCode(reverse.Diagnostics, DomainMismatch) || !hasCode(reverse.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("reverse strengthening not structurally rejected: %+v", reverse)
	}
	if losses := g.Transformations()[0].Losses; len(losses) != 1 || losses[0].Kind != DomainScopeRestriction {
		t.Fatalf("domain loss not derived/recorded: %+v", losses)
	}
}

func TestStructuralStrengtheningAndIntermediateDomainLossAreRejected(t *testing.T) {
	g := NewGraph()
	all := semantic.Claim{ID: "all", Proposition: criticalLineStatement(semantic.ForAll, semantic.NontrivialZeros(semantic.RiemannZeta)), Evidence: []semantic.Evidence{{Kind: semantic.KnownTheoremEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: testReference}}
	mustAddClaim(t, g, all)
	boundedID, err := (DomainRestriction{TargetID: "bounded", Domain: semantic.ZerosBelowHeight(semantic.RiemannZeta, 100), TransformationID: "restrict-loss"}).Apply(g, all.ID)
	if err != nil {
		t.Fatal(err)
	}
	unsound := semantic.Claim{ID: "unsound-expanded", Proposition: all.Proposition, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{boundedID}, Transformation: "expand", Source: testReference}}
	mustAddClaim(t, g, unsound)
	if err := g.AddTransformation(Transformation{ID: "expand", Pass: "test", From: boundedID, To: unsound.ID, Relation: Implies}); err == nil {
		t.Fatal("domain strengthening was admitted without an independent equivalence theorem")
	}

	bridge := derivedClaim("bridge", []semantic.ClaimID{boundedID}, "to-bridge")
	target := semantic.Claim{ID: "rehydrated", Proposition: all.Proposition, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{bridge.ID}, Transformation: "from-bridge", Source: testReference}}
	mustAddClaim(t, g, bridge)
	mustAddClaim(t, g, target)
	if err := g.AddTransformation(Transformation{ID: "to-bridge", Pass: "test", From: boundedID, To: bridge.ID, Relation: Implies}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddTransformation(Transformation{ID: "from-bridge", Pass: "test", From: bridge.ID, To: target.ID, Relation: Implies}); err != nil {
		t.Fatal(err)
	}
	attempt := g.CheckDischarge(all.ID, target.ID)
	if attempt.Accepted || !hasCode(attempt.Diagnostics, InformationLost) {
		t.Fatalf("intermediate domain loss was hidden by a later representation: %+v", attempt)
	}
}

func TestDensitySemanticsAreStructuralAndM0Regresses(t *testing.T) {
	result, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	density, _ := result.Graph.Claim(DensityOneID)
	q, ok := semantic.Quantified(density.Proposition)
	if !ok || q.Quantifier != semantic.DensityOne {
		t.Fatalf("density is not structural: %+v", density.Proposition)
	}
	if result.Attempt.Accepted || !hasCode(result.Attempt.Diagnostics, QuantifierMismatch) || !hasCode(result.Attempt.Diagnostics, InformationLost) {
		t.Fatalf("density-one discharged RH or lacked structural diagnostics: %+v", result.Attempt)
	}
	if loss := result.Graph.Transformations()[1].Losses; len(loss) != 1 || loss[0].Kind != QuantifierWeakening {
		t.Fatalf("quantifier loss missing: %+v", loss)
	}
}

func TestM1BoundedClaimCannotCertifyRH(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	bounded, _ := result.Graph.Claim(BoundedZerosID)
	q, _ := semantic.Quantified(bounded.Proposition)
	if q.Quantifier != semantic.ForAll || q.Domain.Kind != semantic.ZerosBelowHeightDomain || q.Domain.Bound != 1_000_000 {
		t.Fatalf("bad bounded semantics: %+v", q)
	}
	if result.BoundedToRH.Accepted || !hasCode(result.BoundedToRH.Diagnostics, DomainMismatch) {
		t.Fatalf("bounded claim certified RH: %+v", result.BoundedToRH)
	}
}

func TestEulerRepresentationsAttachObjectDomainAndAffordances(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		id   semantic.ClaimID
		name semantic.RepresentationName
	}{{DirichletRepresentationID, semantic.DirichletSeriesRepresentation}, {EulerRepresentationID, semantic.EulerProductRepresentation}} {
		claim, _ := result.Graph.Claim(tc.id)
		p, ok := claim.Proposition.(semantic.RepresentationProposition)
		if !ok || p.Representation.Object != semantic.RiemannZeta || p.Representation.Name != tc.name || p.Representation.ValidOn != semantic.HalfPlaneReGreaterThanOne() || len(p.Representation.Affordances) == 0 {
			t.Fatalf("bad representation %s: %+v", tc.id, claim.Proposition)
		}
	}
	lineage, err := result.Graph.Lineage(ZeroFreeHalfPlaneID)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []semantic.ClaimID{DirichletRepresentationID, EulerIdentityID, EulerRepresentationID, EulerConvergenceID, EulerFactorsNonzeroID, InfiniteProductTheoremID} {
		if !containsID(lineage, want) {
			t.Fatalf("zero-free provenance omits %s: %v", want, lineage)
		}
	}
}

func TestEulerMultiPremiseDerivationAndMissingPremise(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	if !result.ZeroFreeCertified {
		t.Fatalf("zero-free result not certified: %+v", result.ZeroFreeDiagnostics)
	}
	var derivation Transformation
	for _, tr := range result.Graph.Transformations() {
		if tr.ID == DeriveZeroFreeID {
			derivation = tr
		}
	}
	if len(derivation.Premises) != 3 {
		t.Fatalf("analytic premises not explicit: %+v", derivation)
	}
	missing, err := CompileM1WithOptions(M1Options{TrustInfiniteProductTheorem: false})
	if err != nil {
		t.Fatal(err)
	}
	if missing.ZeroFreeCertified || !hasCode(missing.ZeroFreeDiagnostics, OpenObligation) || !strings.Contains(diagnosticsText(missing.ZeroFreeDiagnostics), string(InfiniteProductTheoremID)) {
		t.Fatalf("missing premise not visible: %+v", missing.ZeroFreeDiagnostics)
	}
}

func TestZeroFreeHalfPlaneCannotDischargeRH(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	if result.ZeroFreeToRH.Accepted || !hasCode(result.ZeroFreeToRH.Diagnostics, DomainMismatch) || !hasCode(result.ZeroFreeToRH.Diagnostics, PredicateMismatch) {
		t.Fatalf("zero-free half-plane discharged RH: %+v", result.ZeroFreeToRH)
	}
}

func TestGraphOwnershipAndOpenObligations(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	claim, _ := result.Graph.Claim(EulerRepresentationID)
	p := claim.Proposition.(semantic.RepresentationProposition)
	p.Representation.Affordances[0] = "mutated"
	again, _ := result.Graph.Claim(EulerRepresentationID)
	if again.Proposition.(semantic.RepresentationProposition).Representation.Affordances[0] == "mutated" {
		t.Fatal("claim accessor exposed owned slice")
	}
	transforms := result.Graph.Transformations()
	transforms[len(transforms)-1].Premises[0] = "mutated"
	if result.Graph.Transformations()[len(transforms)-1].Premises[0] == "mutated" {
		t.Fatal("transformation accessor exposed owned slice")
	}

	g := NewGraph()
	source := authoredClaim("source", semantic.KnownTheoremEvidence)
	obligation := authoredClaim("open", semantic.UnverifiedConjectureEvidence)
	target := derivedClaim("target", []semantic.ClaimID{source.ID, obligation.ID}, "conditional")
	for _, c := range []semantic.Claim{source, obligation, target} {
		mustAddClaim(t, g, c)
	}
	if err := g.AddTransformation(Transformation{ID: "conditional", Pass: "test", From: source.ID, To: target.ID, Relation: Implies, Obligations: []semantic.ClaimID{obligation.ID}}); err != nil {
		t.Fatal(err)
	}
	if certified, diagnostics := g.Certify(target.ID); certified || !hasCode(diagnostics, OpenObligation) {
		t.Fatalf("open obligation disappeared: %v %+v", certified, diagnostics)
	}
}

func TestReportsAreDeterministic(t *testing.T) {
	m0, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	m0JSON, err := JSONReport(m0)
	if err != nil {
		t.Fatal(err)
	}
	m0Human := HumanReport(m0)
	m1, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := M1JSONReport(m1)
	if err != nil {
		t.Fatal(err)
	}
	wantHuman := M1HumanReport(m1)
	for i := 0; i < 5; i++ {
		got0, _ := CompileM0()
		j0, _ := JSONReport(got0)
		if !bytes.Equal(j0, m0JSON) || HumanReport(got0) != m0Human {
			t.Fatal("M0 report changed")
		}
		got1, _ := CompileM1()
		j1, _ := M1JSONReport(got1)
		if !bytes.Equal(j1, wantJSON) || M1HumanReport(got1) != wantHuman {
			t.Fatal("M1 report changed")
		}
	}
	if !strings.Contains(string(wantJSON), `"schema": "riemann.semantic-graph.m1"`) || !strings.Contains(string(wantJSON), `"quantifier": "for_all"`) {
		t.Fatal("M1 structural JSON missing")
	}
	if !strings.Contains(string(wantJSON), `"certified": true`) {
		t.Fatal("M1 certification status missing")
	}
}

func mustAddClaim(t *testing.T, g *Graph, c semantic.Claim) {
	t.Helper()
	if err := g.AddClaim(c); err != nil {
		t.Fatal(err)
	}
}
func hasCode(ds []Diagnostic, code DiagnosticCode) bool {
	for _, d := range ds {
		if d.Code == code {
			return true
		}
	}
	return false
}
func diagnosticsText(ds []Diagnostic) string {
	var b strings.Builder
	for _, d := range ds {
		b.WriteString(d.Message)
	}
	return b.String()
}
