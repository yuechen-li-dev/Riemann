package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	BoundedZerosID            semantic.ClaimID = "zeta.zeros-below-height.on-critical-line"
	RHBoundedConsequenceID    semantic.ClaimID = "rh.consequence.zeros-below-height.on-critical-line"
	DirichletRepresentationID semantic.ClaimID = "zeta.representation.dirichlet-series"
	EulerRepresentationID     semantic.ClaimID = "zeta.representation.euler-product"
	EulerConvergenceID        semantic.ClaimID = "zeta.euler-product.absolute-convergence"
	EulerFactorsNonzeroID     semantic.ClaimID = "zeta.euler-product.nonzero-factors"
	ZeroFreeHalfPlaneID       semantic.ClaimID = "zeta.zero-free.re-gt-1"

	RestrictRHToBoundedID semantic.TransformationID = "restrict-domain:rh-to-bounded-height"

	DirichletTheoremID        semantic.TheoremID = "zeta.dirichlet-representation"
	EulerProductTheoremID     semantic.TheoremID = "zeta.euler-product"
	EulerConvergenceTheoremID semantic.TheoremID = "zeta.euler-product-absolute-convergence"
	EulerFactorsTheoremID     semantic.TheoremID = "zeta.euler-factors-nonzero"
	InfiniteProductTheoremID  semantic.TheoremID = "analysis.infinite-product-nonvanishing"
)

var dlmfDefinitionReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(i), equation 25.2.1 (Dirichlet series for Re(s)>1)", URI: "https://dlmf.nist.gov/25.2.E1"}
var dlmfEulerReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(iv), equation 25.2.11 (Euler product for Re(s)>1)", URI: "https://dlmf.nist.gov/25.2.E11"}
var dlmfZeroReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §1.10(ix) (infinite products) and §25.10(i) (Euler product implies ζ(s)≠0 for Re(s)>1)", URI: "https://dlmf.nist.gov/25.10.i"}
var apostolProductReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(iv), equation 25.2.11; source Apostol (1976), Theorem 11.7", URI: "https://dlmf.nist.gov/25.2.E11"}

// DomainRestriction is the generic universal restriction rule. It knows
// nothing about RH; only typed quantifier, predicate, and subset semantics.
type DomainRestriction struct {
	TargetID         semantic.ClaimID
	Domain           semantic.Domain
	TransformationID semantic.TransformationID
}

func (pass DomainRestriction) Apply(g *Graph, fromID semantic.ClaimID) (semantic.ClaimID, error) {
	from, ok := g.Claim(fromID)
	if !ok {
		return "", fmt.Errorf("domain-restriction: source %q is unknown", fromID)
	}
	statement, ok := semantic.Quantified(from.Proposition)
	if !ok || statement.Quantifier != semantic.ForAll {
		return "", fmt.Errorf("domain-restriction: %w: requires a universal statement", ErrWrongProposition)
	}
	if pass.TargetID == "" || pass.TransformationID == "" {
		return "", fmt.Errorf("domain-restriction: missing IDs")
	}
	if pass.Domain == statement.Domain || !semantic.IsSubset(pass.Domain, statement.Domain) {
		return "", fmt.Errorf("domain-restriction: %s is not a known strict subset of %s", pass.Domain.Describe(), statement.Domain.Describe())
	}
	targetStatement := statement
	targetStatement.Domain = pass.Domain
	target := semantic.Claim{
		ID: pass.TargetID, Proposition: targetStatement, Assumptions: semantic.CloneAssumptions(from.Assumptions),
		Evidence:   []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(pass.TransformationID)}, Note: "universal statement restricted to a structurally known subset"}},
		Exactness:  from.Exactness,
		Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{fromID}, Transformation: pass.TransformationID, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "ForAll domain restriction"}},
	}
	if err := g.AddClaim(target); err != nil {
		return "", err
	}
	if err := g.AddTransformation(Transformation{ID: pass.TransformationID, Pass: "restrict-domain", From: fromID, To: target.ID, Relation: Implies, Losses: []InformationLoss{{Kind: DomainScopeRestriction, Reason: "the conclusion covers only a strict subset of the source domain"}}, Provenance: target.Provenance.Source}); err != nil {
		return "", err
	}
	return target.ID, nil
}

type M1Options struct {
	TrustInfiniteProductTheorem bool
	OmitEulerFactorsTheorem     bool // focused partial-application fixture
}
type M1Result struct {
	Graph               *Graph
	Registry            *ContractRegistry
	Applications        []TheoremApplication
	ZeroFreeCertified   bool
	ZeroFreeDiagnostics []Diagnostic
	BoundedToRH         ProofAttempt
	DensityToRH         ProofAttempt
	ZeroFreeToRH        ProofAttempt
}

func CompileM1() (M1Result, error) {
	return CompileM1WithOptions(M1Options{TrustInfiniteProductTheorem: true})
}

func CompileM1WithOptions(options M1Options) (M1Result, error) {
	m0, err := CompileM0()
	if err != nil {
		return M1Result{}, err
	}
	g := m0.Graph
	_, err = (DomainRestriction{TargetID: RHBoundedConsequenceID, Domain: semantic.ZerosBelowHeight(semantic.RiemannZeta, 1_000_000), TransformationID: RestrictRHToBoundedID}).Apply(g, ZeroLocationID)
	if err != nil {
		return M1Result{}, err
	}
	boundedFixtureReference := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M1 hypothetical trusted bounded-verification fixture (no actual zero computation asserted)"}
	boundedClaim := authoredMathClaim(BoundedZerosID, criticalLineStatement(semantic.ForAll, semantic.ZerosBelowHeight(semantic.RiemannZeta, 1_000_000)), semantic.KnownTheoremEvidence, boundedFixtureReference)
	if err := g.AddClaim(boundedClaim); err != nil {
		return M1Result{}, err
	}

	registry, err := m2TheoremRegistry(options)
	if err != nil {
		return M1Result{}, err
	}
	engine := NewContractEngine(g, registry)
	if err := engine.Saturate(); err != nil {
		return M1Result{}, err
	}
	certified, diagnostics := g.Certify(ZeroFreeHalfPlaneID)
	return M1Result{Graph: g, Registry: registry, Applications: append([]TheoremApplication(nil), engine.Applications...), ZeroFreeCertified: certified, ZeroFreeDiagnostics: diagnostics, BoundedToRH: g.AttemptProof(BoundedZerosID, RHClaimID), DensityToRH: g.AttemptProof(DensityOneID, RHClaimID), ZeroFreeToRH: g.AttemptProof(ZeroFreeHalfPlaneID, RHClaimID)}, nil
}

func m2TheoremRegistry(options M1Options) (*ContractRegistry, error) {
	halfPlane := semantic.HalfPlaneReGreaterThanOne()
	exact := semantic.Exact
	registry := NewContractRegistry()
	known := func(ref semantic.Reference) semantic.Evidence {
		return semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: ref}
	}
	contracts := []TheoremContract{
		{ID: DirichletTheoremID, ConclusionID: DirichletRepresentationID, Conclusion: Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.DirichletSeriesRepresentation), Domain: ConstDomain(halfPlane), Exactness: exact, Formula: "ζ(s) = Σ_{n≥1} n^-s", Affordances: []string{"additive summation structure"}}, Relation: Equivalent, Evidence: semantic.Evidence{Kind: semantic.DefinitionEvidence, Source: dlmfDefinitionReference}, Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.2.1"},
		{ID: EulerProductTheoremID, Premises: []Pattern{{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.DirichletSeriesRepresentation), Domain: ConstDomain(halfPlane), Exactness: exact}}, ConclusionID: EulerRepresentationID, Conclusion: Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.EulerProductRepresentation), Domain: ConstDomain(halfPlane), Exactness: exact, Formula: "ζ(s) = ∏_p (1 - p^-s)^-1", Affordances: []string{"prime factorization / multiplicative structure", "factor-wise nonvanishing argument"}}, Relation: Equivalent, Evidence: known(dlmfEulerReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.2.11"},
		{ID: EulerConvergenceTheoremID, ConclusionID: EulerConvergenceID, Conclusion: Pattern{Kind: semantic.AnalyticFactKind, AnalyticFact: ConstAnalyticFact(semantic.EulerProductConvergesAbsolutely), Object: ConstObject(semantic.RiemannZeta), Domain: ConstDomain(halfPlane), Exactness: exact}, Relation: Implies, Evidence: known(apostolProductReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.2(iv); Apostol Theorem 11.7"},
	}
	if !options.OmitEulerFactorsTheorem {
		contracts = append(contracts, TheoremContract{ID: EulerFactorsTheoremID, ConclusionID: EulerFactorsNonzeroID, Conclusion: Pattern{Kind: semantic.AnalyticFactKind, AnalyticFact: ConstAnalyticFact(semantic.EulerProductFactorsNonzero), Object: ConstObject(semantic.RiemannZeta), Domain: ConstDomain(halfPlane), Exactness: exact}, Relation: Implies, Evidence: known(dlmfEulerReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.2.11"})
	}
	trust := TrustedExternalTheorem
	if !options.TrustInfiniteProductTheorem {
		trust = UntrustedTheorem
	}
	contracts = append(contracts, TheoremContract{ID: InfiniteProductTheoremID, Parameters: []Parameter{{ID: "F", Type: ObjectParam}, {ID: "D", Type: DomainParam}}, Premises: []Pattern{
		{Kind: semantic.RepresentationKind, Object: Var("F"), Representation: ConstRepresentation(semantic.EulerProductRepresentation), Domain: Var("D"), Exactness: exact},
		{Kind: semantic.AnalyticFactKind, AnalyticFact: ConstAnalyticFact(semantic.EulerProductConvergesAbsolutely), Object: Var("F"), Domain: Var("D"), Exactness: exact},
		{Kind: semantic.AnalyticFactKind, AnalyticFact: ConstAnalyticFact(semantic.EulerProductFactorsNonzero), Object: Var("F"), Domain: Var("D"), Exactness: exact},
	}, ConclusionID: ZeroFreeHalfPlaneID, Conclusion: Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: Var("D"), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: Var("F"), Exactness: exact}, Relation: Implies, Evidence: known(dlmfZeroReference), Trust: trust, Citation: "NIST DLMF 1.10(ix), 25.10(i)"})
	for _, contract := range contracts {
		if err := registry.Register(contract); err != nil {
			return nil, err
		}
	}
	return registry, nil
}

func authoredMathClaim(id semantic.ClaimID, proposition semantic.Proposition, evidence semantic.EvidenceKind, ref semantic.Reference) semantic.Claim {
	return semantic.Claim{ID: id, Proposition: proposition, Evidence: []semantic.Evidence{{Kind: evidence, Source: ref}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: ref}}
}
