package compiler

import (
	"fmt"
	"sort"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	AnalyticContinuationID          semantic.ClaimID = "zeta.representation.analytic-continuation"
	AnalyticContinuationAvailableID semantic.ClaimID = "zeta.side-condition.analytic-continuation"
	CompletedXiID                   semantic.ClaimID = "zeta.representation.completed-xi"
	XiFunctionalEquationID          semantic.ClaimID = "zeta.identity.xi-functional-equation"
	ZetaConjugationConditionID      semantic.ClaimID = "zeta.side-condition.real-analytic-conjugation"
	ZetaPoleID                      semantic.ClaimID = "zeta.classification.simple-pole"
	ZetaTrivialZerosID              semantic.ClaimID = "zeta.classification.trivial-zeros"
	ZetaBoundaryZeroFreeID          semantic.ClaimID = "zeta.zero-free.re-eq-1"
	SampleZeroID                    semantic.ClaimID = "zeta.sample.nontrivial-zero"
	ConjugationInvariantID          semantic.ClaimID = "zeta.nontrivial-zeros.conjugation-invariant"
	FunctionalInvariantID           semantic.ClaimID = "zeta.nontrivial-zeros.functional-invariant"
	CriticalReflectionInvariantID   semantic.ClaimID = "zeta.nontrivial-zeros.critical-reflection-invariant"
	CriticalStripConfinementID      semantic.ClaimID = "zeta.nontrivial-zeros.critical-strip"
	RHFixedPointID                  semantic.ClaimID = "rh.critical-reflection-fixed"

	AnalyticContinuationTheoremID            semantic.TheoremID = "zeta.analytic-continuation"
	ContinuationAvailabilityTheoremID        semantic.TheoremID = "zeta.analytic-continuation-availability"
	CompletedXiTheoremID                     semantic.TheoremID = "zeta.completed-xi-definition"
	XiFunctionalEquationTheoremID            semantic.TheoremID = "zeta.xi-functional-equation"
	ZetaConjugationTheoremID                 semantic.TheoremID = "zeta.conjugation-zero-transport"
	PointInContinuationDomainTheoremID       semantic.TheoremID = "zeta.nontrivial-zero-in-continuation-domain"
	CompletionFactorTheoremID                semantic.TheoremID = "zeta.completion-factor-at-nontrivial-zero"
	ReflectedCompletionFactorTheoremID       semantic.TheoremID = "zeta.completion-factor-at-reflected-nontrivial-zero"
	LiftZeroToXiTheoremID                    semantic.TheoremID = "zeta.zero-to-xi-zero"
	FunctionalIdentityZeroTransportTheoremID semantic.TheoremID = "analysis.functional-identity-zero-transport"
	LowerXiZeroTheoremID                     semantic.TheoremID = "zeta.xi-zero-to-nontrivial-zero"
	PoleTheoremID                            semantic.TheoremID = "zeta.simple-pole-at-one"
	TrivialZerosTheoremID                    semantic.TheoremID = "zeta.trivial-zero-classification"
	BoundaryZeroFreeTheoremID                semantic.TheoremID = "zeta.boundary-zero-exclusion"
	ConjugationInvariantTheoremID            semantic.TheoremID = "zeta.conjugation-set-invariance"
	FunctionalInvariantTheoremID             semantic.TheoremID = "zeta.functional-set-invariance"
	ComposeInvarianceTheoremID               semantic.TheoremID = "sets.compose-transform-invariance"
	StripConfinementTheoremID                semantic.TheoremID = "zeta.critical-strip-confinement"
	RHFixedPointEquivalenceTheoremID         semantic.TheoremID = "geometry.critical-reflection-fixed-point-equivalence"
)

var dlmfContinuationReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.2(i): ζ is defined elsewhere by analytic continuation and is meromorphic with only a simple pole at s=1", URI: "https://dlmf.nist.gov/25.2.i"}
var dlmfFunctionalEquationReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.4, equations 25.4.3–25.4.4: ξ(s)=ξ(1-s), ξ(s)=½s(s-1)Γ(s/2)π^(-s/2)ζ(s)", URI: "https://dlmf.nist.gov/25.4.E3"}
var dlmfZerosReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §25.10(i): trivial zeros, boundary zero exclusion, critical-strip confinement, and symmetry", URI: "https://dlmf.nist.gov/25.10.i"}
var compilerTransportReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M3 typed transport composition over trusted theorem contracts"}

type M3Options struct {
	OmitCompletionFactor     bool
	OmitAnalyticContinuation bool
	SamplePoint              semantic.PointExpr
}

type SymmetryOrbit struct {
	Generated []semantic.PointExpr `json:"generated"`
	Distinct  []semantic.PointExpr `json:"distinct"`
}

type M3Result struct {
	M1                  M1Result
	Graph               *Graph
	Registry            *ContractRegistry
	Applications        []TheoremApplication
	Orbit               SymmetryOrbit
	StripCertified      bool
	StripDiagnostics    []Diagnostic
	SymmetryCertified   bool
	SymmetryDiagnostics []Diagnostic
}

func CompileM3() (M3Result, error) { return CompileM3WithOptions(M3Options{}) }

func CompileM3WithOptions(options M3Options) (M3Result, error) {
	m1, err := CompileM1()
	if err != nil {
		return M3Result{}, err
	}
	if options.SamplePoint.Symbol == "" {
		options.SamplePoint = semantic.Point("ρ")
	}
	sampleRef := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M3 conditional source claim: let ρ be a classified nontrivial zero"}
	sample := semantic.Claim{ID: SampleZeroID, Proposition: semantic.ZeroAtPoint{Object: semantic.RiemannZeta, Point: options.SamplePoint.Canonical(), Classification: semantic.NontrivialZero}, Assumptions: []semantic.Assumption{{ID: "rho-is-nontrivial-zero", Description: "ρ is a nontrivial zero of ζ"}}, Evidence: []semantic.Evidence{{Kind: semantic.DefinitionEvidence, Source: sampleRef, Note: "conditional input, not a newly asserted zero"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: sampleRef}}
	if err := m1.Graph.AddClaim(sample); err != nil {
		return M3Result{}, err
	}
	registry, err := m3TheoremRegistry(options)
	if err != nil {
		return M3Result{}, err
	}
	engine := NewContractEngine(m1.Graph, registry)
	if err := engine.Saturate(); err != nil {
		return M3Result{}, err
	}
	if err := addRHFixedPointTarget(m1.Graph, registry); err != nil {
		return M3Result{}, err
	}
	orbit := collectOrbit(m1.Graph, options.SamplePoint)
	strip, stripDiag := m1.Graph.Certify(CriticalStripConfinementID)
	sym, symDiag := m1.Graph.Certify(CriticalReflectionInvariantID)
	return M3Result{M1: m1, Graph: m1.Graph, Registry: registry, Applications: append([]TheoremApplication(nil), engine.Applications...), Orbit: orbit, StripCertified: strip, StripDiagnostics: stripDiag, SymmetryCertified: sym, SymmetryDiagnostics: symDiag}, nil
}

func m3TheoremRegistry(options M3Options) (*ContractRegistry, error) {
	r := NewContractRegistry()
	exact := semantic.Exact
	known := func(ref semantic.Reference) semantic.Evidence {
		return semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: ref}
	}
	derived := semantic.Evidence{Kind: semantic.DerivedEvidence, Source: compilerTransportReference}
	complex := semantic.ComplexPlane()
	punctured := semantic.ComplexPlaneExceptOne()
	nz := semantic.NontrivialZeros(semantic.RiemannZeta)
	contracts := []TheoremContract{}
	if !options.OmitAnalyticContinuation {
		contracts = append(contracts,
			TheoremContract{ID: AnalyticContinuationTheoremID, ConclusionID: AnalyticContinuationID, Conclusion: Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.AnalyticContinuationRepresentation), Domain: ConstDomain(punctured), Exactness: exact, Formula: "meromorphic continuation of ζ(s), finite on ℂ\\{1}", Affordances: []string{"global evaluation away from the pole"}}, Relation: Equivalent, Evidence: known(dlmfContinuationReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.2(i)"},
			TheoremContract{ID: ContinuationAvailabilityTheoremID, Premises: []Pattern{{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.AnalyticContinuationRepresentation), Domain: ConstDomain(punctured), Exactness: exact}}, ConclusionID: AnalyticContinuationAvailableID, Conclusion: Pattern{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(semantic.AnalyticContinuationAvailable), Object: ConstObject(semantic.RiemannZeta), Domain: ConstDomain(punctured), Exactness: exact}, Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule})
	}
	contracts = append(contracts,
		TheoremContract{ID: CompletedXiTheoremID, ConclusionID: CompletedXiID, Conclusion: Pattern{Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannXi), BaseObject: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.CompletedXiRepresentation), Domain: ConstDomain(complex), Exactness: exact, Formula: "ξ(s) = ½ s(s-1) Γ(s/2) π^(-s/2) ζ(s)", Affordances: []string{"entire completion", "zero correspondence away from exceptional factors"}}, Relation: Equivalent, Evidence: known(dlmfFunctionalEquationReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.4.4"},
		TheoremContract{ID: XiFunctionalEquationTheoremID, ConclusionID: XiFunctionalEquationID, Conclusion: Pattern{Kind: semantic.FunctionalIdentityKind, Object: ConstObject(semantic.RiemannXi), Left: ConstTransform(semantic.IdentityTransform), Right: ConstTransform(semantic.OneMinusTransform), Domain: ConstDomain(complex), Exactness: exact, Formula: "ξ(s) = ξ(1-s)"}, Relation: Equivalent, Evidence: known(dlmfFunctionalEquationReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.4.3"},
		TheoremContract{ID: PoleTheoremID, ConclusionID: ZetaPoleID, Conclusion: Pattern{Kind: semantic.ZeroClassificationKind, Object: ConstObject(semantic.RiemannZeta), Classification: ConstZeroClassification(semantic.SimplePoleAtOne), Exactness: exact}, Relation: Implies, Evidence: known(dlmfContinuationReference), Trust: TrustedExternalTheorem},
		TheoremContract{ID: TrivialZerosTheoremID, ConclusionID: ZetaTrivialZerosID, Conclusion: Pattern{Kind: semantic.ZeroClassificationKind, Object: ConstObject(semantic.RiemannZeta), Classification: ConstZeroClassification(semantic.TrivialZerosExactlyNegativeEven), Exactness: exact}, Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem},
		TheoremContract{ID: BoundaryZeroFreeTheoremID, ConclusionID: ZetaBoundaryZeroFreeID, Conclusion: Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(semantic.LineReEqualsOne()), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: exact}, Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.10(i): ζ(s)≠0 on Re(s)=1"},
		TheoremContract{ID: "zeta.real-analytic-conjugation", ConclusionID: ZetaConjugationConditionID, Conclusion: Pattern{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(semantic.RealAnalyticConjugationAvailable), Object: ConstObject(semantic.RiemannZeta), Domain: ConstDomain(punctured), Exactness: exact}, Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "NIST DLMF 25.10(i)"})

	pointParams := []Parameter{{ID: "P", Type: PointParam}}
	if !options.OmitCompletionFactor {
		contracts = append(contracts,
			TheoremContract{ID: CompletionFactorTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero), trivialClassificationPattern(), poleClassificationPattern()}, Conclusion: sidePointPattern(semantic.CompletionFactorRegularNonzero, semantic.RiemannXi, Var("P")), Relation: Implies, Evidence: known(dlmfFunctionalEquationReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.4.4 with the pole and nontrivial-zero classifications in 25.2(i), 25.10(i)"},
			TheoremContract{ID: ReflectedCompletionFactorTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero), zeroFreeRightPattern(), trivialClassificationPattern(), poleClassificationPattern()}, Conclusion: sidePointPattern(semantic.CompletionFactorRegularNonzero, semantic.RiemannXi, TransformedPoint("P", semantic.OneMinusTransform)), Relation: Implies, Evidence: known(dlmfFunctionalEquationReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.4.3–25.4.4; M1 zero exclusion; DLMF 25.2(i), 25.10(i) classifications"})
	}
	contracts = append(contracts,
		TheoremContract{ID: PointInContinuationDomainTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero), poleClassificationPattern()}, Conclusion: sidePointDomainPattern(semantic.PointInValidityDomain, semantic.RiemannZeta, Var("P"), punctured), Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule},
		TheoremContract{ID: ZetaConjugationTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero)}, SideConditions: []Pattern{sidePointDomainPattern(semantic.PointInValidityDomain, semantic.RiemannZeta, Var("P"), punctured), sideDomainPattern(semantic.AnalyticContinuationAvailable, semantic.RiemannZeta, punctured), sideDomainPattern(semantic.RealAnalyticConjugationAvailable, semantic.RiemannZeta, punctured)}, Conclusion: zeroPattern(semantic.RiemannZeta, TransformedPoint("P", semantic.ConjugateTransform), semantic.NontrivialZero), Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "ζ(conjugate(s))=conjugate(ζ(s)); DLMF 25.10(i)"},
		TheoremContract{ID: LiftZeroToXiTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero)}, SideConditions: []Pattern{sidePointPattern(semantic.CompletionFactorRegularNonzero, semantic.RiemannXi, Var("P")), {Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannXi), BaseObject: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.CompletedXiRepresentation), Domain: ConstDomain(complex), Exactness: exact}}, Conclusion: zeroPattern(semantic.RiemannXi, Var("P"), semantic.NontrivialZero), Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule},
		TheoremContract{ID: FunctionalIdentityZeroTransportTheoremID, Parameters: []Parameter{{ID: "F", Type: ObjectParam}, {ID: "P", Type: PointParam}, {ID: "Z", Type: ZeroClassParam}, {ID: "D", Type: DomainParam}, {ID: "T", Type: TransformParam}}, Premises: []Pattern{zeroPatternTerms(Var("F"), Var("P"), Var("Z")), {Kind: semantic.FunctionalIdentityKind, Object: Var("F"), Left: ConstTransform(semantic.IdentityTransform), Right: Var("T"), Domain: Var("D"), Exactness: exact}}, Conclusion: zeroPatternTerms(Var("F"), PointTransformedBy("P", "T"), Var("Z")), Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule},
		TheoremContract{ID: LowerXiZeroTheoremID, Parameters: pointParams, Premises: []Pattern{zeroPattern(semantic.RiemannXi, Var("P"), semantic.NontrivialZero)}, SideConditions: []Pattern{sidePointPattern(semantic.CompletionFactorRegularNonzero, semantic.RiemannXi, Var("P")), {Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannXi), BaseObject: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.CompletedXiRepresentation), Domain: ConstDomain(complex), Exactness: exact}}, Conclusion: zeroPattern(semantic.RiemannZeta, Var("P"), semantic.NontrivialZero), Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule})

	contracts = append(contracts,
		TheoremContract{ID: ConjugationInvariantTheoremID, Premises: []Pattern{{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(semantic.RealAnalyticConjugationAvailable), Object: ConstObject(semantic.RiemannZeta), Domain: ConstDomain(punctured), Exactness: exact}, {Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.AnalyticContinuationRepresentation), Domain: ConstDomain(punctured), Exactness: exact}}, ConclusionID: ConjugationInvariantID, Conclusion: setInvariantPattern(nz, semantic.ConjugateTransform), Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.10(i): symmetry about the real axis"},
		TheoremContract{ID: FunctionalInvariantTheoremID, Premises: []Pattern{{Kind: semantic.FunctionalIdentityKind, Object: ConstObject(semantic.RiemannXi), Left: ConstTransform(semantic.IdentityTransform), Right: ConstTransform(semantic.OneMinusTransform), Domain: ConstDomain(complex), Exactness: exact}, zeroFreeRightPattern(), trivialClassificationPattern(), poleClassificationPattern(), {Kind: semantic.RepresentationKind, Object: ConstObject(semantic.RiemannXi), BaseObject: ConstObject(semantic.RiemannZeta), Representation: ConstRepresentation(semantic.CompletedXiRepresentation), Domain: ConstDomain(complex), Exactness: exact}}, ConclusionID: FunctionalInvariantID, Conclusion: setInvariantPattern(nz, semantic.OneMinusTransform), Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.4.3–25.4.4 and 25.10(i)"},
		TheoremContract{ID: ComposeInvarianceTheoremID, Premises: []Pattern{setInvariantPattern(nz, semantic.ConjugateTransform), setInvariantPattern(nz, semantic.OneMinusTransform)}, ConclusionID: CriticalReflectionInvariantID, Conclusion: setInvariantPattern(nz, semantic.CriticalReflectionTransform), Relation: Implies, Evidence: derived, Trust: CompilerVerifiedRule},
		TheoremContract{ID: StripConfinementTheoremID, Premises: []Pattern{zeroFreeRightPattern(), setInvariantPattern(nz, semantic.OneMinusTransform), trivialClassificationPattern(), {Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(semantic.LineReEqualsOne()), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: exact}}, ConclusionID: CriticalStripConfinementID, Conclusion: Pattern{Kind: semantic.ZeroSetPropertyKind, Domain: ConstDomain(nz), Property: ConstZeroSetProperty(semantic.ConfinedToRegion), Region: ConstDomain(semantic.CriticalStrip()), Exactness: exact}, Relation: Implies, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "DLMF 25.10(i), instantiated with explicit zero-free, reflection, trivial-zero, and boundary premises"},
		TheoremContract{ID: RHFixedPointEquivalenceTheoremID, Premises: []Pattern{{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(nz), Predicate: ConstPredicate(semantic.RealPartEqualsHalfPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: exact}}, Conclusion: Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(nz), Predicate: ConstPredicate(semantic.CriticalReflectionFixedPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: exact}, Relation: Equivalent, Evidence: known(dlmfZerosReference), Trust: TrustedExternalTheorem, Citation: "fixed points of s↦1-conjugate(s) are exactly Re(s)=1/2"})
	for _, c := range contracts {
		if err := r.Register(c); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func zeroPattern(f semantic.Function, p Term, z semantic.ZeroClass) Pattern {
	return zeroPatternTerms(ConstObject(f), p, ConstZeroClass(z))
}
func zeroPatternTerms(f, p, z Term) Pattern {
	return Pattern{Kind: semantic.ZeroAtPointKind, Object: f, Point: p, Classification: z, Exactness: semantic.Exact}
}
func sidePointPattern(name semantic.SideConditionName, f semantic.Function, p Term) Pattern {
	return Pattern{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(name), Object: ConstObject(f), Point: p, Exactness: semantic.Exact}
}
func sidePointDomainPattern(name semantic.SideConditionName, f semantic.Function, p Term, d semantic.Domain) Pattern {
	return Pattern{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(name), Object: ConstObject(f), Point: p, Domain: ConstDomain(d), Exactness: semantic.Exact}
}
func sideDomainPattern(name semantic.SideConditionName, f semantic.Function, d semantic.Domain) Pattern {
	return Pattern{Kind: semantic.SideConditionKind, SideCondition: ConstSideCondition(name), Object: ConstObject(f), Domain: ConstDomain(d), Exactness: semantic.Exact}
}
func setInvariantPattern(d semantic.Domain, t semantic.PointTransform) Pattern {
	return Pattern{Kind: semantic.ZeroSetPropertyKind, Domain: ConstDomain(d), Property: ConstZeroSetProperty(semantic.InvariantUnderTransform), Transform: ConstTransform(t), Exactness: semantic.Exact}
}
func zeroFreeRightPattern() Pattern {
	return Pattern{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(semantic.HalfPlaneReGreaterThanOne()), Predicate: ConstPredicate(semantic.FunctionNonzeroPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: semantic.Exact}
}
func trivialClassificationPattern() Pattern {
	return Pattern{Kind: semantic.ZeroClassificationKind, Object: ConstObject(semantic.RiemannZeta), Classification: ConstZeroClassification(semantic.TrivialZerosExactlyNegativeEven), Exactness: semantic.Exact}
}
func poleClassificationPattern() Pattern {
	return Pattern{Kind: semantic.ZeroClassificationKind, Object: ConstObject(semantic.RiemannZeta), Classification: ConstZeroClassification(semantic.SimplePoleAtOne), Exactness: semantic.Exact}
}

func addRHFixedPointTarget(g *Graph, r *ContractRegistry) error {
	if _, ok := r.Contract(RHFixedPointEquivalenceTheoremID); !ok {
		return fmt.Errorf("RH fixed-point equivalence contract missing")
	}
	source, ok := g.Claim(ZeroLocationID)
	if !ok {
		return fmt.Errorf("RH normalized target missing")
	}
	p := source.Proposition.(semantic.QuantifiedStatement)
	p.Predicate.Kind = semantic.CriticalReflectionFixedPredicate
	tid := semantic.TransformationID("theorem:" + string(RHFixedPointEquivalenceTheoremID) + ":goal")
	target := semantic.Claim{ID: RHFixedPointID, Proposition: p, Assumptions: semantic.CloneAssumptions(source.Assumptions), Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: dlmfZerosReference, Note: "equivalent target normalization; does not prove RH"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{source.ID}, Transformation: tid, Theorem: RHFixedPointEquivalenceTheoremID, Source: dlmfZerosReference}}
	if err := g.AddClaim(target); err != nil {
		return err
	}
	return g.AddTransformation(Transformation{ID: tid, Pass: "normalize-goal-by-theorem-contract", From: source.ID, To: target.ID, Relation: Equivalent, Provenance: dlmfZerosReference, Theorem: RHFixedPointEquivalenceTheoremID, Trusted: true})
}

func collectOrbit(g *Graph, base semantic.PointExpr) SymmetryOrbit {
	symbol := base.Symbol
	generated := []semantic.PointExpr{base.Canonical(), base.Apply(semantic.ConjugateTransform), base.Apply(semantic.OneMinusTransform), base.Apply(semantic.CriticalReflectionTransform)}
	byKey := map[string]semantic.PointExpr{}
	for _, claim := range g.Claims() {
		if z, ok := claim.Proposition.(semantic.ZeroAtPoint); ok && z.Object == semantic.RiemannZeta && z.Classification == semantic.NontrivialZero && z.Point.Symbol == symbol {
			byKey[z.Point.Key()] = z.Point.Canonical()
		}
	}
	keys := make([]string, 0, len(byKey))
	for k := range byKey {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	distinct := make([]semantic.PointExpr, 0, len(keys))
	for _, k := range keys {
		distinct = append(distinct, byKey[k])
	}
	return SymmetryOrbit{Generated: generated, Distinct: distinct}
}
