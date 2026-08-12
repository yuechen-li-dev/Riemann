package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func representationPattern(name semantic.RepresentationName) Pattern {
	return Pattern{Kind: semantic.RepresentationKind, Object: Var("F"), Representation: ConstRepresentation(name), Domain: Var("D"), Exactness: semantic.Exact, Formula: "typed formula"}
}
func factPattern(fact semantic.AnalyticFactName) Pattern {
	return Pattern{Kind: semantic.AnalyticFactKind, AnalyticFact: ConstAnalyticFact(fact), Object: Var("F"), Domain: Var("D"), Exactness: semantic.Exact}
}
func nonzeroPattern(q semantic.QuantifierKind) Pattern {
	return Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(q), Domain: Var("D"), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: Var("F"), Exactness: semantic.Exact}
}
func theoremFixture(id semantic.TheoremID, premises []Pattern, conclusion Pattern) TheoremContract {
	return TheoremContract{ID: id, Parameters: []Parameter{{ID: "F", Type: ObjectParam}, {ID: "D", Type: DomainParam}}, Premises: premises, Conclusion: conclusion, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: testReference}, Trust: TrustedExternalTheorem}
}
func exactClaim(id semantic.ClaimID, p semantic.Proposition) semantic.Claim {
	return semantic.Claim{ID: id, Proposition: p, Evidence: []semantic.Evidence{{Kind: semantic.KnownTheoremEvidence, Source: testReference}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: testReference}}
}
func representationClaim(id semantic.ClaimID, object semantic.Function, name semantic.RepresentationName, domain semantic.Domain) semantic.Claim {
	return exactClaim(id, semantic.RepresentationProposition{Representation: semantic.Representation{Object: object, Name: name, ValidOn: domain, Formula: "display only"}})
}
func factClaim(id semantic.ClaimID, object semantic.Function, fact semantic.AnalyticFactName, domain semantic.Domain) semantic.Claim {
	return exactClaim(id, semantic.AnalyticFact{Fact: fact, Object: object, Domain: domain})
}

func TestContractRegistrationIsTypedAndDeterministic(t *testing.T) {
	r := NewContractRegistry()
	b := theoremFixture("b", []Pattern{representationPattern(semantic.EulerProductRepresentation)}, factPattern(semantic.EulerProductConvergesAbsolutely))
	a := theoremFixture("a", []Pattern{representationPattern(semantic.EulerProductRepresentation)}, factPattern(semantic.EulerProductFactorsNonzero))
	if err := r.Register(b); err != nil {
		t.Fatal(err)
	}
	if err := r.Register(a); err != nil {
		t.Fatal(err)
	}
	if got := r.Contracts(); len(got) != 2 || got[0].ID != "a" || got[1].ID != "b" {
		t.Fatalf("registry order is not canonical: %+v", got)
	}
	if _, ok := r.Contract("a"); !ok {
		t.Fatal("registered theorem cannot be retrieved")
	}
	bad := theoremFixture("bad", nil, Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: Var("F"), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: Var("F"), Exactness: semantic.Exact})
	if err := r.Register(bad); err == nil {
		t.Fatal("ill-typed parameter use was registered")
	}
}

func TestRepeatedParametersUnifyAndConflictsRemainObligations(t *testing.T) {
	d := semantic.HalfPlaneReGreaterThanOne()
	g := NewGraph()
	for _, c := range []semantic.Claim{
		representationClaim("product", semantic.RiemannZeta, semantic.EulerProductRepresentation, d),
		factClaim("convergence", semantic.RiemannZeta, semantic.EulerProductConvergesAbsolutely, d),
		factClaim("factors", semantic.RiemannZeta, semantic.EulerProductFactorsNonzero, d),
	} {
		mustAddClaim(t, g, c)
	}
	r := NewContractRegistry()
	contract := theoremFixture("nonzero", []Pattern{representationPattern(semantic.EulerProductRepresentation), factPattern(semantic.EulerProductConvergesAbsolutely), factPattern(semantic.EulerProductFactorsNonzero)}, nonzeroPattern(semantic.ForAll))
	if err := r.Register(contract); err != nil {
		t.Fatal(err)
	}
	e := NewContractEngine(g, r)
	if err := e.Saturate(); err != nil {
		t.Fatal(err)
	}
	if len(e.Applications) != 1 || !e.Applications[0].Complete || len(e.Applications[0].Bindings) != 2 {
		t.Fatalf("consistent parameters did not bind: %+v", e.Applications)
	}

	g2 := NewGraph()
	mustAddClaim(t, g2, representationClaim("product", semantic.RiemannZeta, semantic.EulerProductRepresentation, d))
	mustAddClaim(t, g2, factClaim("bad-domain", semantic.RiemannZeta, semantic.EulerProductConvergesAbsolutely, semantic.CriticalStrip()))
	e2 := NewContractEngine(g2, r)
	if err := e2.Saturate(); err != nil {
		t.Fatal(err)
	}
	if len(e2.Applications) != 1 || e2.Applications[0].Complete || len(e2.Applications[0].Obligations) != 2 || !strings.Contains(e2.Applications[0].Obligations[0].Description, "Re(s) > 1") {
		t.Fatalf("domain conflict did not retain bound obligation: %+v", e2.Applications)
	}
}

func TestExactMatchingRejectsObjectDomainQuantifierPredicateAndNumericalEvidence(t *testing.T) {
	d := semantic.HalfPlaneReGreaterThanOne()
	other := semantic.Function(2)
	cases := []struct {
		name    string
		claim   semantic.Claim
		pattern Pattern
	}{
		{"object", representationClaim("x", other, semantic.EulerProductRepresentation, d), Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.EulerProductRepresentation), Domain: ConstDomain(d), Exactness: semantic.Exact}},
		{"domain", representationClaim("x", semantic.RiemannZeta, semantic.EulerProductRepresentation, semantic.CriticalStrip()), Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.EulerProductRepresentation), Domain: ConstDomain(d), Exactness: semantic.Exact}},
		{"representation", representationClaim("x", semantic.RiemannZeta, semantic.DirichletSeriesRepresentation, d), Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.EulerProductRepresentation), Domain: ConstDomain(d), Exactness: semantic.Exact}},
		{"quantifier", exactClaim("x", semantic.QuantifiedStatement{Quantifier: semantic.DensityOne, Domain: d, Predicate: semantic.Predicate{Kind: semantic.FunctionNonzeroPredicate, Function: semantic.RiemannZeta}}), Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(d), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: semantic.Exact}},
		{"predicate", exactClaim("x", semantic.QuantifiedStatement{Quantifier: semantic.ForAll, Domain: d, Predicate: semantic.Predicate{Kind: semantic.RealPartEqualsHalfPredicate, Function: semantic.RiemannZeta}}), Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(d), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: semantic.Exact}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			env := make(bindingEnv)
			if matchPattern(tc.pattern, tc.claim, env) {
				t.Fatal("semantic mismatch matched")
			}
		})
	}
	numerical := representationClaim("numerical", semantic.RiemannZeta, semantic.EulerProductRepresentation, d)
	numerical.Evidence[0].Kind = semantic.NumericalExperimentEvidence
	g := NewGraph()
	mustAddClaim(t, g, numerical)
	r := NewContractRegistry()
	contract := theoremFixture("exact-only", []Pattern{representationPattern(semantic.EulerProductRepresentation)}, factPattern(semantic.EulerProductConvergesAbsolutely))
	if err := r.Register(contract); err != nil {
		t.Fatal(err)
	}
	e := NewContractEngine(g, r)
	if err := e.Saturate(); err != nil {
		t.Fatal(err)
	}
	if len(e.Applications) != 0 {
		t.Fatalf("numerical evidence satisfied exact premise: %+v", e.Applications)
	}
}

func TestDeterministicCompositionCyclesAndAlternateDerivations(t *testing.T) {
	build := func(reverse bool) ([]byte, *ContractEngine) {
		g := NewGraph()
		d := semantic.HalfPlaneReGreaterThanOne()
		mustAddClaim(t, g, representationClaim("product", semantic.RiemannZeta, semantic.EulerProductRepresentation, d))
		mustAddClaim(t, g, factClaim("factors", semantic.RiemannZeta, semantic.EulerProductFactorsNonzero, d))
		r := NewContractRegistry()
		a := theoremFixture("a.product-to-convergence", []Pattern{representationPattern(semantic.EulerProductRepresentation)}, factPattern(semantic.EulerProductConvergesAbsolutely))
		b := theoremFixture("b.to-nonzero", []Pattern{factPattern(semantic.EulerProductConvergesAbsolutely), factPattern(semantic.EulerProductFactorsNonzero)}, nonzeroPattern(semantic.ForAll))
		alt := theoremFixture("c.alternate", []Pattern{representationPattern(semantic.EulerProductRepresentation), factPattern(semantic.EulerProductFactorsNonzero)}, nonzeroPattern(semantic.ForAll))
		cycle := theoremFixture("d.cycle", []Pattern{factPattern(semantic.EulerProductConvergesAbsolutely)}, representationPattern(semantic.EulerProductRepresentation))
		contracts := []TheoremContract{a, b, alt, cycle}
		if reverse {
			for i, j := 0, len(contracts)-1; i < j; i, j = i+1, j-1 {
				contracts[i], contracts[j] = contracts[j], contracts[i]
			}
		}
		for _, c := range contracts {
			if err := r.Register(c); err != nil {
				t.Fatal(err)
			}
		}
		e := NewContractEngine(g, r)
		if err := e.Saturate(); err != nil {
			t.Fatal(err)
		}
		result := M1Result{Graph: g, Registry: r, Applications: e.Applications}
		data, err := marshalGraphWithContracts("test", g, r.Contracts(), e.Applications, nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		_ = result
		return data, e
	}
	one, e1 := build(false)
	two, _ := build(true)
	if !bytes.Equal(one, two) {
		t.Fatal("composition depends on registration order")
	}
	var conclusion semantic.ClaimID
	derivations := 0
	for _, a := range e1.Applications {
		if a.Complete && (a.Theorem == "b.to-nonzero" || a.Theorem == "c.alternate") {
			conclusion = a.Conclusion
			derivations++
		}
	}
	if derivations != 2 {
		t.Fatalf("alternate theorem applications lost: %+v", e1.Applications)
	}
	edges := 0
	for _, tr := range e1.Graph.Transformations() {
		if tr.To == conclusion {
			edges++
		}
	}
	if edges != 2 {
		t.Fatalf("alternate proof edges overwritten: %d", edges)
	}
	if len(e1.Graph.Claims()) > 5 {
		t.Fatalf("cycle produced duplicate claims: %d", len(e1.Graph.Claims()))
	}
}

func TestM2EulerRegressionUsesOnlyContracts(t *testing.T) {
	result, err := CompileM1()
	if err != nil {
		t.Fatal(err)
	}
	if !result.ZeroFreeCertified {
		t.Fatalf("contract chain did not certify zero-free half-plane: %+v", result.ZeroFreeDiagnostics)
	}
	want := []semantic.TheoremID{DirichletTheoremID, EulerProductTheoremID, EulerConvergenceTheoremID, EulerFactorsTheoremID, InfiniteProductTheoremID}
	for _, id := range want {
		if !hasApplication(result.Applications, id) {
			t.Fatalf("missing theorem instance %s", id)
		}
	}
	for _, tr := range result.Graph.Transformations() {
		if tr.To == EulerRepresentationID || tr.To == ZeroFreeHalfPlaneID {
			if tr.Pass != "instantiate-theorem-contract" || tr.Theorem == "" {
				t.Fatalf("Euler result still uses bespoke pass: %+v", tr)
			}
		}
	}
	if result.ZeroFreeToRH.Accepted {
		t.Fatal("zero-free half-plane discharged RH")
	}
	human := M1HumanReport(result)
	for _, text := range []string{"IMPORT THEOREM", "INSTANTIATE", "bindings:", "CERTIFIED\n  relative to trusted imported theorem contracts"} {
		if !strings.Contains(human, text) {
			t.Fatalf("human report omits %q", text)
		}
	}
}

func TestEulerMissingPremiseProducesBoundObligationWithoutConclusion(t *testing.T) {
	result, err := CompileM1WithOptions(M1Options{TrustInfiniteProductTheorem: true, OmitEulerFactorsTheorem: true})
	if err != nil {
		t.Fatal(err)
	}
	if _, exists := result.Graph.Claim(ZeroFreeHalfPlaneID); exists || result.ZeroFreeCertified {
		t.Fatal("partial theorem application emitted a certifiable conclusion")
	}
	for _, application := range result.Applications {
		if application.Theorem != InfiniteProductTheoremID {
			continue
		}
		if application.Complete || len(application.Obligations) != 1 {
			t.Fatalf("missing premise was not explicit: %+v", application)
		}
		if !strings.Contains(application.Obligations[0].Description, "nonzero") || !strings.Contains(application.Obligations[0].Description, "Re(s) > 1") {
			t.Fatalf("obligation lost inferred bindings: %+v", application.Obligations[0])
		}
		report := M1HumanReport(result)
		if !strings.Contains(report, "UNRESOLVED") || !strings.Contains(report, "conclusion withheld") {
			t.Fatalf("human report hides partial application:\n%s", report)
		}
		return
	}
	t.Fatal("missing-premise candidate application was not recorded")
}
