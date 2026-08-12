package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func authoredClaim(id semantic.ClaimID, evidence semantic.EvidenceKind) semantic.Claim {
	return semantic.Claim{
		ID: id, Proposition: semantic.NamedObligation{Name: string(id)},
		Evidence:   []semantic.Evidence{{Kind: evidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "test"}}},
		Exactness:  semantic.Exact,
		Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "test"}},
	}
}

func derivedClaim(id, parent semantic.ClaimID, transformation semantic.TransformationID) semantic.Claim {
	claim := authoredClaim(id, semantic.DerivedEvidence)
	claim.Provenance = semantic.Provenance{
		Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{parent}, Transformation: transformation,
		Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(transformation)},
	}
	return claim
}

func addPair(t *testing.T, relation Relation) (*Graph, semantic.ClaimID, semantic.ClaimID) {
	t.Helper()
	g := NewGraph()
	from, to := semantic.ClaimID("a"), semantic.ClaimID("b")
	transformation := semantic.TransformationID("a-to-b")
	if err := g.AddClaim(authoredClaim(from, semantic.KnownTheoremEvidence)); err != nil {
		t.Fatal(err)
	}
	if err := g.AddClaim(derivedClaim(to, from, transformation)); err != nil {
		t.Fatal(err)
	}
	if err := g.AddTransformation(Transformation{ID: transformation, Pass: "test", From: from, To: to, Relation: relation}); err != nil {
		t.Fatal(err)
	}
	return g, from, to
}

func TestEquivalenceIsBidirectional(t *testing.T) {
	g, a, b := addPair(t, Equivalent)
	if attempt := g.CheckDischarge(a, b); !attempt.Accepted {
		t.Fatalf("forward equivalence rejected: %+v", attempt)
	}
	if attempt := g.CheckDischarge(b, a); !attempt.Accepted {
		t.Fatalf("reverse equivalence rejected: %+v", attempt)
	}
	if certified, diagnostics := g.Certify(b); !certified {
		t.Fatalf("derived equivalent claim not certified: %+v", diagnostics)
	}
}

func TestImplicationOnlyUsesEstablishedDirection(t *testing.T) {
	g, a, b := addPair(t, Implies)
	if attempt := g.CheckDischarge(a, b); !attempt.Accepted {
		t.Fatalf("forward implication rejected: %+v", attempt)
	}
	attempt := g.CheckDischarge(b, a)
	if attempt.Accepted || !hasCode(attempt.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("reverse implication was not rejected by direction: %+v", attempt)
	}
}

func TestApproximationAndNumericalEvidenceCannotCertifyExactTheorem(t *testing.T) {
	g, _, target := addPair(t, Approximation)
	if certified, diagnostics := g.Certify(target); certified || !hasCode(diagnostics, ApproximationBoundary) {
		t.Fatalf("approximation certified exact claim: certified=%v diagnostics=%+v", certified, diagnostics)
	}

	numerical := authoredClaim("numerical", semantic.NumericalExperimentEvidence)
	if err := g.AddClaim(numerical); err != nil {
		t.Fatal(err)
	}
	if attempt := g.AttemptProof(numerical.ID, numerical.ID); attempt.Accepted || !hasCode(attempt.Diagnostics, UncertifiedEvidence) {
		t.Fatalf("numerical evidence discharged an exact target: %+v", attempt)
	}
	if certified, diagnostics := g.Certify(numerical.ID); certified || !hasCode(diagnostics, UncertifiedEvidence) {
		t.Fatalf("numerical evidence certified exact claim: certified=%v diagnostics=%+v", certified, diagnostics)
	}
}

func TestInformationLossIsDeclaredAndBlocksStrongerTarget(t *testing.T) {
	result, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	attempt := result.Attempt
	if attempt.Accepted || !hasCode(attempt.Diagnostics, NoEstablishedDirection) || !hasCode(attempt.Diagnostics, MissingCapability) {
		t.Fatalf("density-one improperly usable as RH: %+v", attempt)
	}
	transformations := result.Graph.Transformations()
	if got := transformations[1].Losses; len(got) != 1 || got[0].Property != semantic.ExceptionalSetSensitivity {
		t.Fatalf("density loss not explicit: %+v", got)
	}

	g := NewGraph()
	from := authoredClaim("sensitive", semantic.KnownTheoremEvidence)
	from.Capabilities = semantic.Properties(semantic.ExceptionalSetSensitivity)
	to := derivedClaim("insensitive", from.ID, "lossy")
	if err := g.AddClaim(from); err != nil {
		t.Fatal(err)
	}
	if err := g.AddClaim(to); err != nil {
		t.Fatal(err)
	}
	if err := g.AddTransformation(Transformation{ID: "lossy", Pass: "test", From: from.ID, To: to.ID, Relation: Relaxation}); err == nil {
		t.Fatal("undeclared capability loss was accepted")
	}
}

func TestIntermediateLossCannotBeHiddenByALaterPass(t *testing.T) {
	g := NewGraph()
	a := authoredClaim("a", semantic.KnownTheoremEvidence)
	a.Capabilities = semantic.Properties(semantic.ExceptionalSetSensitivity)
	b := derivedClaim("b", a.ID, "a-to-b")
	c := derivedClaim("c", b.ID, "b-to-c")
	c.Requirements = semantic.Properties(semantic.ExceptionalSetSensitivity)
	obligation := authoredClaim("restore-detail", semantic.KnownTheoremEvidence)
	for _, claim := range []semantic.Claim{a, b, c, obligation} {
		if err := g.AddClaim(claim); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.AddTransformation(Transformation{
		ID: "a-to-b", Pass: "lossy", From: a.ID, To: b.ID, Relation: Relaxation,
		Losses: []InformationLoss{{Property: semantic.ExceptionalSetSensitivity, Reason: "test loss"}},
	}); err != nil {
		t.Fatal(err)
	}
	if err := g.AddTransformation(Transformation{
		ID: "b-to-c", Pass: "restore", From: b.ID, To: c.ID, Relation: Implies,
		Obligations: []semantic.ClaimID{obligation.ID},
	}); err != nil {
		t.Fatal(err)
	}
	if attempt := g.CheckDischarge(a.ID, c.ID); attempt.Accepted || !hasCode(attempt.Diagnostics, MissingCapability) {
		t.Fatalf("intermediate information loss was hidden: %+v", attempt)
	}
}

func TestGraphAccessorsDoNotExposeMutableSlices(t *testing.T) {
	result, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	transformations := result.Graph.Transformations()
	transformations[1].Losses[0].Reason = "mutated"
	if got := result.Graph.Transformations()[1].Losses[0].Reason; got == "mutated" {
		t.Fatal("transformation accessor exposed graph-owned loss slice")
	}
	claim, _ := result.Graph.Claim(DensityOneID)
	claim.Provenance.Parents[0] = "mutated"
	lineage, err := result.Graph.Lineage(DensityOneID)
	if err != nil {
		t.Fatal(err)
	}
	if lineage[0] != RHClaimID {
		t.Fatalf("claim accessor mutated provenance: %v", lineage)
	}
}

func TestOpenObligationRemainsVisibleAndPreventsCertification(t *testing.T) {
	g := NewGraph()
	source := authoredClaim("premise", semantic.KnownTheoremEvidence)
	obligation := authoredClaim("open-lemma", semantic.UnverifiedConjectureEvidence)
	target := derivedClaim("conclusion", source.ID, "conditional")
	for _, claim := range []semantic.Claim{source, obligation, target} {
		if err := g.AddClaim(claim); err != nil {
			t.Fatal(err)
		}
	}
	if err := g.AddTransformation(Transformation{
		ID: "conditional", Pass: "test", From: source.ID, To: target.ID,
		Relation: Implies, Obligations: []semantic.ClaimID{obligation.ID},
	}); err != nil {
		t.Fatal(err)
	}
	if attempt := g.CheckDischarge(source.ID, target.ID); attempt.Accepted || !hasCode(attempt.Diagnostics, OpenObligation) {
		t.Fatalf("open obligation disappeared from proof use: %+v", attempt)
	}
	if certified, diagnostics := g.Certify(target.ID); certified || !hasCode(diagnostics, OpenObligation) {
		t.Fatalf("open obligation did not prevent certification: %v %+v", certified, diagnostics)
	}
}

func TestTransformationCannotSilentlyDropAssumptions(t *testing.T) {
	g := NewGraph()
	source := authoredClaim("conditional-premise", semantic.KnownTheoremEvidence)
	source.Assumptions = []semantic.Assumption{{ID: "hypothesis", Description: "a declared hypothesis"}}
	target := derivedClaim("unconditional-conclusion", source.ID, "drop-assumption")
	if err := g.AddClaim(source); err != nil {
		t.Fatal(err)
	}
	if err := g.AddClaim(target); err != nil {
		t.Fatal(err)
	}
	if err := g.AddTransformation(Transformation{
		ID: "drop-assumption", Pass: "test", From: source.ID, To: target.ID, Relation: Implies,
	}); err == nil {
		t.Fatal("transformation silently dropped an assumption")
	}
}

func TestProvenanceTraceReachesRoot(t *testing.T) {
	result, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	lineage, err := result.Graph.Lineage(DensityOneID)
	if err != nil {
		t.Fatal(err)
	}
	want := []semantic.ClaimID{RHClaimID, ZeroLocationID, DensityOneID}
	if len(lineage) != len(want) {
		t.Fatalf("lineage = %v, want %v", lineage, want)
	}
	for i := range want {
		if lineage[i] != want[i] {
			t.Fatalf("lineage = %v, want %v", lineage, want)
		}
	}
}

func TestReportsAreDeterministic(t *testing.T) {
	first, err := CompileM0()
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := JSONReport(first)
	if err != nil {
		t.Fatal(err)
	}
	wantHuman := HumanReport(first)
	for i := 0; i < 5; i++ {
		result, err := CompileM0()
		if err != nil {
			t.Fatal(err)
		}
		gotJSON, err := JSONReport(result)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(gotJSON, wantJSON) {
			t.Fatal("JSON report order changed")
		}
		if got := HumanReport(result); got != wantHuman {
			t.Fatal("human report order changed")
		}
	}
	if !strings.Contains(string(wantJSON), `"schema": "riemann.semantic-graph.m0"`) {
		t.Fatal("missing graph schema")
	}
}

func hasCode(diagnostics []Diagnostic, code DiagnosticCode) bool {
	for _, diagnostic := range diagnostics {
		if diagnostic.Code == code {
			return true
		}
	}
	return false
}
