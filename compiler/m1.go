package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	BoundedZerosID            semantic.ClaimID = "zeta.zeros-below-height.on-critical-line"
	RHBoundedConsequenceID    semantic.ClaimID = "rh.consequence.zeros-below-height.on-critical-line"
	DirichletRepresentationID semantic.ClaimID = "zeta.representation.dirichlet-series"
	EulerIdentityID           semantic.ClaimID = "zeta.identity.dirichlet-to-euler"
	EulerRepresentationID     semantic.ClaimID = "zeta.representation.euler-product"
	EulerConvergenceID        semantic.ClaimID = "zeta.euler-product.absolute-convergence"
	EulerFactorsNonzeroID     semantic.ClaimID = "zeta.euler-product.nonzero-factors"
	InfiniteProductTheoremID  semantic.ClaimID = "analysis.infinite-product.nonzero-limit"
	ZeroFreeHalfPlaneID       semantic.ClaimID = "zeta.zero-free.re-gt-1"

	RestrictRHToBoundedID semantic.TransformationID = "restrict-domain:rh-to-bounded-height"
	ChangeToEulerID       semantic.TransformationID = "change-representation:dirichlet-to-euler"
	DeriveZeroFreeID      semantic.TransformationID = "derive:euler-product-zero-free-half-plane"
)

var dlmfDefinitionReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(i), equation 25.2.1 (Dirichlet series for Re(s)>1)", URI: "https://dlmf.nist.gov/25.2.E1"}
var dlmfEulerReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(iv), equation 25.2.11 (Euler product for Re(s)>1)", URI: "https://dlmf.nist.gov/25.2.E11"}
var dlmfZeroReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.10(i) (the product representation implies ζ(s)≠0 for Re(s)>1)", URI: "https://dlmf.nist.gov/25.10.i"}
var apostolProductReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "Tom M. Apostol, Introduction to Analytic Number Theory (1976), §11.5; cited by NIST DLMF §25.2(iv)", URI: "https://dlmf.nist.gov/25.2.iv"}

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

type M1Options struct{ TrustInfiniteProductTheorem bool }
type M1Result struct {
	Graph               *Graph
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

	halfPlane := semantic.HalfPlaneReGreaterThanOne()
	dirichlet := semantic.Representation{Object: semantic.RiemannZeta, Name: semantic.DirichletSeriesRepresentation, ValidOn: halfPlane, Formula: "ζ(s) = Σ_{n≥1} n^-s", Affordances: []string{"additive summation structure"}}
	euler := semantic.Representation{Object: semantic.RiemannZeta, Name: semantic.EulerProductRepresentation, ValidOn: halfPlane, Formula: "ζ(s) = ∏_p (1 - p^-s)^-1", Affordances: []string{"prime factorization / multiplicative structure", "factor-wise nonvanishing argument"}}
	dirichletClaim := authoredMathClaim(DirichletRepresentationID, semantic.RepresentationProposition{Representation: dirichlet}, semantic.DefinitionEvidence, dlmfDefinitionReference)
	identityClaim := authoredMathClaim(EulerIdentityID, semantic.RepresentationIdentity{Object: semantic.RiemannZeta, Left: semantic.DirichletSeriesRepresentation, Right: semantic.EulerProductRepresentation, Domain: halfPlane}, semantic.KnownTheoremEvidence, dlmfEulerReference)
	for _, claim := range []semantic.Claim{dirichletClaim, identityClaim} {
		if err := g.AddClaim(claim); err != nil {
			return M1Result{}, err
		}
	}
	eulerClaim := semantic.Claim{
		ID: EulerRepresentationID, Proposition: semantic.RepresentationProposition{Representation: euler},
		Evidence:   []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(ChangeToEulerID)}, Note: "representation change is valid only on the stated domain and depends on the imported Euler identity"}},
		Exactness:  semantic.Exact,
		Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{DirichletRepresentationID, EulerIdentityID}, Transformation: ChangeToEulerID, Source: dlmfEulerReference},
	}
	if err := g.AddClaim(eulerClaim); err != nil {
		return M1Result{}, err
	}
	if err := g.AddTransformation(Transformation{ID: ChangeToEulerID, Pass: "change-representation", From: DirichletRepresentationID, To: EulerRepresentationID, Relation: Equivalent, Obligations: []semantic.ClaimID{EulerIdentityID}, Provenance: dlmfEulerReference}); err != nil {
		return M1Result{}, err
	}

	convergence := authoredMathClaim(EulerConvergenceID, semantic.AnalyticFact{Fact: semantic.EulerProductConvergesAbsolutely, Object: semantic.RiemannZeta, Domain: halfPlane}, semantic.KnownTheoremEvidence, apostolProductReference)
	factors := authoredMathClaim(EulerFactorsNonzeroID, semantic.AnalyticFact{Fact: semantic.EulerProductFactorsNonzero, Object: semantic.RiemannZeta, Domain: halfPlane}, semantic.KnownTheoremEvidence, apostolProductReference)
	theoremEvidence := semantic.UnverifiedConjectureEvidence
	if options.TrustInfiniteProductTheorem {
		theoremEvidence = semantic.KnownTheoremEvidence
	}
	productTheorem := authoredMathClaim(InfiniteProductTheoremID, semantic.AnalyticFact{Fact: semantic.ConvergentProductNonzeroLimit, Object: semantic.RiemannZeta, Domain: halfPlane}, theoremEvidence, apostolProductReference)
	for _, claim := range []semantic.Claim{convergence, factors, productTheorem} {
		if err := g.AddClaim(claim); err != nil {
			return M1Result{}, err
		}
	}
	zeroFree := semantic.Claim{
		ID:          ZeroFreeHalfPlaneID,
		Proposition: semantic.QuantifiedStatement{Quantifier: semantic.ForAll, Domain: halfPlane, Predicate: semantic.Predicate{Kind: semantic.FunctionNonzeroPredicate, Function: semantic.RiemannZeta}},
		Evidence:    []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(DeriveZeroFreeID)}, Note: "derived from the Euler representation and all explicit analytic premises"}},
		Exactness:   semantic.Exact,
		Provenance:  semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{EulerRepresentationID, EulerConvergenceID, EulerFactorsNonzeroID, InfiniteProductTheoremID}, Transformation: DeriveZeroFreeID, Source: dlmfZeroReference},
	}
	if err := g.AddClaim(zeroFree); err != nil {
		return M1Result{}, err
	}
	if err := g.AddTransformation(Transformation{ID: DeriveZeroFreeID, Pass: "euler-product-nonvanishing", From: EulerRepresentationID, Premises: []semantic.ClaimID{EulerConvergenceID, EulerFactorsNonzeroID, InfiniteProductTheoremID}, To: ZeroFreeHalfPlaneID, Relation: Implies, Provenance: dlmfZeroReference}); err != nil {
		return M1Result{}, err
	}
	certified, diagnostics := g.Certify(ZeroFreeHalfPlaneID)
	return M1Result{Graph: g, ZeroFreeCertified: certified, ZeroFreeDiagnostics: diagnostics, BoundedToRH: g.AttemptProof(BoundedZerosID, RHClaimID), DensityToRH: g.AttemptProof(DensityOneID, RHClaimID), ZeroFreeToRH: g.AttemptProof(ZeroFreeHalfPlaneID, RHClaimID)}, nil
}

func authoredMathClaim(id semantic.ClaimID, proposition semantic.Proposition, evidence semantic.EvidenceKind, ref semantic.Reference) semantic.Claim {
	return semantic.Claim{ID: id, Proposition: proposition, Evidence: []semantic.Evidence{{Kind: evidence, Source: ref}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: ref}}
}
