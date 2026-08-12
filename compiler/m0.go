package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	RHClaimID      semantic.ClaimID          = "rh"
	ZeroLocationID semantic.ClaimID          = "rh.zero-location"
	DensityOneID   semantic.ClaimID          = "rh.critical-line-density-one"
	NormalizeRHID  semantic.TransformationID = "normalize-rh:rh"
	DensityPassID  semantic.TransformationID = "critical-zero-density:rh.zero-location"
)

var rhReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "NIST Digital Library of Mathematical Functions, §25.10(i), Zeros of the Riemann zeta function",
	URI:      "https://dlmf.nist.gov/25.10.i",
}

func criticalLineStatement(q semantic.QuantifierKind, domain semantic.Domain) semantic.QuantifiedStatement {
	return semantic.QuantifiedStatement{
		Quantifier: q, Domain: domain,
		Predicate: semantic.Predicate{Kind: semantic.RealPartEqualsHalfPredicate, Function: semantic.RiemannZeta},
	}
}

// NormalizeRH retains the M0 definitional normalization, now between two
// structurally identical quantified propositions rather than capability flags.
type NormalizeRH struct{}

func (NormalizeRH) Apply(g *Graph, fromID semantic.ClaimID) (semantic.ClaimID, error) {
	from, ok := g.Claim(fromID)
	if !ok {
		return "", fmt.Errorf("normalize-rh: source claim %q is unknown", fromID)
	}
	statement, ok := semantic.Quantified(from.Proposition)
	if !ok || statement.Quantifier != semantic.ForAll || statement.Domain != semantic.NontrivialZeros(semantic.RiemannZeta) || statement.Predicate.Kind != semantic.RealPartEqualsHalfPredicate {
		return "", fmt.Errorf("normalize-rh: %w: got %s", ErrWrongProposition, from.Proposition.Kind())
	}
	target := semantic.Claim{
		ID: ZeroLocationID, Proposition: statement,
		Assumptions: semantic.CloneAssumptions(from.Assumptions),
		Evidence:    []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(NormalizeRHID)}, Note: "definitional normalization; this does not assert that RH is true"}},
		Exactness:   semantic.Exact,
		Provenance:  semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{fromID}, Transformation: NormalizeRHID, Source: rhReference},
	}
	if err := g.AddClaim(target); err != nil {
		return "", err
	}
	if err := g.AddTransformation(Transformation{ID: NormalizeRHID, Pass: "normalize-rh", From: fromID, To: target.ID, Relation: Equivalent, Provenance: rhReference}); err != nil {
		return "", err
	}
	return target.ID, nil
}

type CriticalZeroDensity struct{}

func (CriticalZeroDensity) Apply(g *Graph, fromID semantic.ClaimID) (semantic.ClaimID, error) {
	from, ok := g.Claim(fromID)
	if !ok {
		return "", fmt.Errorf("critical-zero-density: source claim %q is unknown", fromID)
	}
	statement, ok := semantic.Quantified(from.Proposition)
	if !ok || statement.Quantifier != semantic.ForAll || statement.Predicate.Kind != semantic.RealPartEqualsHalfPredicate {
		return "", fmt.Errorf("critical-zero-density: %w: got %s", ErrWrongProposition, from.Proposition.Kind())
	}
	density := statement
	density.Quantifier = semantic.DensityOne
	target := semantic.Claim{
		ID: DensityOneID, Proposition: density,
		Assumptions: semantic.CloneAssumptions(from.Assumptions),
		Evidence:    []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: string(DensityPassID)}, Note: "conditional consequence; density semantics do not exclude a density-zero exceptional set"}},
		Exactness:   semantic.Exact,
		Provenance:  semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{fromID}, Transformation: DensityPassID, Source: semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Universal membership implies density-one membership."}},
	}
	if err := g.AddClaim(target); err != nil {
		return "", err
	}
	if err := g.AddTransformation(Transformation{
		ID: DensityPassID, Pass: "critical-zero-density", From: fromID, To: target.ID, Relation: Relaxation,
		Losses:     []InformationLoss{{Kind: QuantifierWeakening, Reason: "density-one quantification permits a density-zero exceptional set"}},
		Provenance: target.Provenance.Source,
	}); err != nil {
		return "", err
	}
	return target.ID, nil
}

type M0Result struct {
	Graph   *Graph
	Attempt ProofAttempt
}

func CompileM0() (M0Result, error) {
	g := NewGraph()
	root := semantic.Claim{
		ID:          RHClaimID,
		Proposition: criticalLineStatement(semantic.ForAll, semantic.NontrivialZeros(semantic.RiemannZeta)),
		Evidence:    []semantic.Evidence{{Kind: semantic.UnverifiedConjectureEvidence, Source: rhReference, Note: "open theorem target; the reference defines the statement but does not prove it"}},
		Exactness:   semantic.Exact,
		Provenance:  semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: rhReference},
	}
	if err := g.AddClaim(root); err != nil {
		return M0Result{}, err
	}
	zeroLocation, err := (NormalizeRH{}).Apply(g, root.ID)
	if err != nil {
		return M0Result{}, err
	}
	density, err := (CriticalZeroDensity{}).Apply(g, zeroLocation)
	if err != nil {
		return M0Result{}, err
	}
	return M0Result{Graph: g, Attempt: g.AttemptProof(density, root.ID)}, nil
}
