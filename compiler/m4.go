package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	WeilFunctionalDefinitionID  semantic.ClaimID = "weil.functional.definition"
	WeilExplicitFormulaID       semantic.ClaimID = "weil.explicit-formula.identity"
	WeilPositivityID            semantic.ClaimID = "rh.weil-universal-positivity"
	FiniteWeilPositivityID      semantic.ClaimID = "rh.weil-finite-family-positivity"
	NumericalFinitePositivityID semantic.ClaimID = "experiment.weil-finite-family-positivity"
	WeilF1AdmissibleID          semantic.ClaimID = "weil.function.f1.admissible"
	WeilF2AdmissibleID          semantic.ClaimID = "weil.function.f2.admissible"

	WeilCriterionTheoremID       semantic.TheoremID        = "weil-lagarias.positivity-criterion"
	WeilExplicitFormulaTheoremID semantic.TheoremID        = "weil-lagarias.explicit-formula"
	LowerRHToWeilID              semantic.TransformationID = "lower-rh-to-weil-positivity"
	RestrictWeilFamilyID         semantic.TransformationID = "restrict-weil-positivity-to-finite-family"
)

var lagariasWeilReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "J. C. Lagarias, The Riemann Hypothesis: Arithmetic and Geometry, §3, Theorems 3.1 and 3.2 (Mellin form of Weil's criterion)",
	URI:      "https://dept.math.lsa.umich.edu/~lagarias/doc/mt-holyoke-rev.pdf",
}

type M4Result struct {
	M3                   M3Result
	Graph                *Graph
	Registry             *ContractRegistry
	FiniteFamily         semantic.FunctionClass
	FullToFinite         ProofAttempt
	FiniteToRH           ProofAttempt
	NumericalToUniversal ProofAttempt
}

func CompileM4() (M4Result, error) {
	m3, err := CompileM3()
	if err != nil {
		return M4Result{}, err
	}
	g := m3.Graph
	full := semantic.WeilNiceClass()
	functional := weilFunctional()

	definition := authoredMathClaim(WeilFunctionalDefinitionID, semantic.FunctionalDefinition{Functional: functional}, semantic.DefinitionEvidence, lagariasWeilReference)
	if err := g.AddClaim(definition); err != nil {
		return M4Result{}, err
	}
	explicit := authoredMathClaim(WeilExplicitFormulaID, weilExplicitFormula(functional, full), semantic.KnownTheoremEvidence, lagariasWeilReference)
	if err := g.AddClaim(explicit); err != nil {
		return M4Result{}, err
	}

	registry, err := m4TheoremRegistry(full)
	if err != nil {
		return M4Result{}, err
	}
	if err := addWeilEquivalentTarget(g, full); err != nil {
		return M4Result{}, err
	}

	f1 := abstractWeilFunction("f1", 1, full)
	f2 := abstractWeilFunction("f2", 2, full)
	for _, item := range []struct {
		id semantic.ClaimID
		f  semantic.TestFunction
	}{{WeilF1AdmissibleID, f1}, {WeilF2AdmissibleID, f2}} {
		ref := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M4 abstract basis member defined with a Weil-nice admissibility certificate"}
		if err := g.AddClaim(authoredMathClaim(item.id, semantic.TestFunctionAdmissibility{Function: item.f, Class: full}, semantic.DefinitionEvidence, ref)); err != nil {
			return M4Result{}, err
		}
	}
	finite := semantic.FunctionClass{ID: "m4-demo-family", Kind: semantic.FiniteFunctionClass, Constraints: append([]semantic.FunctionConstraint(nil), full.Constraints...), TransformConvention: full.TransformConvention, Members: []semantic.TestFunction{f1, f2}}
	if _, err := (FunctionClassRestriction{TargetID: FiniteWeilPositivityID, Class: finite, TransformationID: RestrictWeilFamilyID}).Apply(g, WeilPositivityID); err != nil {
		return M4Result{}, err
	}

	numerical := semantic.Claim{ID: NumericalFinitePositivityID, Proposition: functionalStatement(finite), Evidence: []semantic.Evidence{{Kind: semantic.NumericalExperimentEvidence, Source: semantic.Reference{Kind: semantic.ExperimentRecord, Citation: "M4 finite-sample evidence boundary fixture"}, Note: "hypothetical sampled nonnegative values; not universal evidence"}}, Exactness: semantic.Approximate, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: semantic.Reference{Kind: semantic.ExperimentRecord, Citation: "M4 finite-sample evidence boundary fixture"}}}
	if err := g.AddClaim(numerical); err != nil {
		return M4Result{}, err
	}

	return M4Result{M3: m3, Graph: g, Registry: registry, FiniteFamily: finite, FullToFinite: g.CheckDischarge(WeilPositivityID, FiniteWeilPositivityID), FiniteToRH: g.AttemptProof(FiniteWeilPositivityID, RHClaimID), NumericalToUniversal: g.AttemptProof(NumericalFinitePositivityID, WeilPositivityID)}, nil
}

func m4TheoremRegistry(full semantic.FunctionClass) (*ContractRegistry, error) {
	r := NewContractRegistry()
	contract := TheoremContract{
		ID:         WeilCriterionTheoremID,
		Premises:   []Pattern{{Kind: semantic.QuantifiedStatementKind, Quantifier: ConstQuantifier(semantic.ForAll), Domain: ConstDomain(semantic.NontrivialZeros(semantic.RiemannZeta)), Predicate: ConstPredicate(semantic.RealPartEqualsHalfPredicate), Object: ConstObject(semantic.RiemannZeta), Exactness: semantic.Exact}},
		Conclusion: functionalPattern(full), Relation: Equivalent,
		Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: lagariasWeilReference, Note: "trusted theorem import; it does not assert either equivalent open statement"},
		Trust:    TrustedExternalTheorem, Citation: "Lagarias §3, Theorem 3.2; after Theorem 3.1's explicit formula",
	}
	if err := r.Register(contract); err != nil {
		return nil, err
	}
	return r, nil
}

func functionalPattern(class semantic.FunctionClass) Pattern {
	c := semantic.CloneFunctionClass(class)
	return Pattern{Kind: semantic.UniversalFunctionalStatementKind, Exactness: semantic.Exact, FunctionClass: &c, Functional: semantic.WeilZetaQuadraticFunctional, FunctionalPredicate: semantic.FunctionalNonnegative, TransformConvention: semantic.LagariasMellinConvention}
}
func functionalStatement(class semantic.FunctionClass) semantic.UniversalFunctionalStatement {
	return semantic.UniversalFunctionalStatement{Quantifier: semantic.ForAll, Variable: "f", FunctionClass: semantic.CloneFunctionClass(class), Functional: semantic.WeilZetaQuadraticFunctional, Predicate: semantic.FunctionalNonnegative, TransformConvention: semantic.LagariasMellinConvention}
}

func addWeilEquivalentTarget(g *Graph, full semantic.FunctionClass) error {
	source, ok := g.Claim(ZeroLocationID)
	if !ok {
		return fmt.Errorf("RH zero-location target missing")
	}
	target := semantic.Claim{ID: WeilPositivityID, Proposition: functionalStatement(full), Assumptions: semantic.CloneAssumptions(source.Assumptions), Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: lagariasWeilReference, Note: "equivalent open target, not a proof of positivity or RH"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{source.ID, WeilFunctionalDefinitionID, WeilExplicitFormulaID}, Transformation: LowerRHToWeilID, Theorem: WeilCriterionTheoremID, Source: lagariasWeilReference}}
	if err := g.AddClaim(target); err != nil {
		return err
	}
	return g.AddTransformation(Transformation{ID: LowerRHToWeilID, Pass: "lower-by-trusted-weil-criterion", From: source.ID, Premises: []semantic.ClaimID{WeilFunctionalDefinitionID, WeilExplicitFormulaID}, To: target.ID, Relation: Equivalent, Provenance: lagariasWeilReference, Theorem: WeilCriterionTheoremID, Trusted: true})
}

// FunctionClassRestriction is generic: it knows neither RH nor Weil. Every
// finite member must have an exact, certified admissibility proposition.
type FunctionClassRestriction struct {
	TargetID         semantic.ClaimID
	Class            semantic.FunctionClass
	TransformationID semantic.TransformationID
}

func (pass FunctionClassRestriction) Apply(g *Graph, fromID semantic.ClaimID) (semantic.ClaimID, error) {
	from, ok := g.Claim(fromID)
	if !ok {
		return "", fmt.Errorf("function-class restriction: source %q is unknown", fromID)
	}
	statement, ok := from.Proposition.(semantic.UniversalFunctionalStatement)
	if !ok || statement.Quantifier != semantic.ForAll {
		return "", fmt.Errorf("function-class restriction: %w", ErrWrongProposition)
	}
	if pass.TargetID == "" || pass.TransformationID == "" || pass.Class.Kind != semantic.FiniteFunctionClass {
		return "", fmt.Errorf("function-class restriction: invalid target")
	}
	if err := pass.Class.Validate(); err != nil {
		return "", err
	}
	if !functionClassShapeSubset(pass.Class, statement.FunctionClass) {
		return "", fmt.Errorf("function-class restriction: incompatible class or transform convention")
	}
	var certificates []semantic.ClaimID
	for _, member := range pass.Class.Members {
		found := false
		for _, claim := range g.Claims() {
			admissible, ok := claim.Proposition.(semantic.TestFunctionAdmissibility)
			if ok && admissible.Function.Key() == member.Key() && admissible.Class.Key() == statement.FunctionClass.Key() {
				if certified, _ := g.Certify(claim.ID); certified {
					certificates = append(certificates, claim.ID)
					found = true
					break
				}
			}
		}
		if !found {
			return "", fmt.Errorf("function-class restriction: member %s lacks certified admissibility in %s", member.Symbol, statement.FunctionClass.ID)
		}
	}
	targetStatement := statement
	targetStatement.FunctionClass = semantic.CloneFunctionClass(pass.Class)
	parents := append([]semantic.ClaimID{fromID}, certificates...)
	ref := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "ForAll restriction to a certified function subclass"}
	target := semantic.Claim{ID: pass.TargetID, Proposition: targetStatement, Assumptions: semantic.CloneAssumptions(from.Assumptions), Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: ref, Note: "function-space coverage is strictly restricted"}}, Exactness: from.Exactness, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: pass.TransformationID, Source: ref}}
	if err := g.AddClaim(target); err != nil {
		return "", err
	}
	return target.ID, g.AddTransformation(Transformation{ID: pass.TransformationID, Pass: "restrict-function-class", From: fromID, Premises: certificates, To: target.ID, Relation: Implies, Losses: []InformationLoss{{Kind: FunctionSpaceRestriction, Reason: "source covers every admissible test function; conclusion covers only a strict finite family"}}, Provenance: ref})
}

func abstractWeilFunction(symbol string, index int, class semantic.FunctionClass) semantic.TestFunction {
	return semantic.TestFunction{Symbol: symbol, Kind: semantic.BasisTestFunction, DeclaredClass: class.ID, RequiredAttributes: append([]semantic.FunctionConstraint(nil), class.Constraints...), TransformConvention: class.TransformConvention, Parameters: []semantic.FunctionParameter{{Name: "index", Value: fmt.Sprint(index)}, {Name: "family", Value: "m4-abstract-weil-basis"}}}
}

func weilFunctional() semantic.QuadraticFunctional {
	aggregate := semantic.Aggregate{Kind: semantic.SumOverDomain, IndexDomain: semantic.NontrivialZeros(semantic.RiemannZeta), Summand: "M[f](rho) * conjugate(M[f](1-conjugate(rho)))", Convergence: []string{"zeros counted with multiplicity", "symmetric limiting order inherited from the explicit formula", "M[f] holomorphic on a strip containing 0 <= Re(s) <= 1 and rapidly decreasing vertically"}, TransformConvention: semantic.LagariasMellinConvention, TheoremLineage: []semantic.TheoremID{WeilExplicitFormulaTheoremID, WeilCriterionTheoremID}, Provenance: lagariasWeilReference}
	return semantic.QuadraticFunctional{ID: semantic.WeilZetaQuadraticFunctional, Object: semantic.RiemannZeta, InputConstruction: "h=f multiplicative-convolution tilde(conjugate(f)); tilde(g)(x)=x^-1 g(x^-1)", TransformConvention: semantic.LagariasMellinConvention, Contributions: []semantic.FunctionalContribution{
		{Kind: semantic.ZeroContribution, RepresentationSide: "zero_side", Sign: 1, Formula: "W^(1)(h)=sum_rho M[h](rho)", Aggregate: &aggregate},
		{Kind: semantic.EndpointContribution, RepresentationSide: "explicit_formula_side", Sign: 1, Formula: "M[h](1)+M[h](0)"},
		{Kind: semantic.PrimePowerContribution, RepresentationSide: "explicit_formula_side", Sign: -1, Formula: "sum_p log(p) sum_{n>=1} (h(p^n)+tilde(h)(p^n))", Index: "p prime and p^n prime power", Weight: "log(p), equivalently von Mangoldt Lambda(p^n)"},
		{Kind: semantic.ArchimedeanContribution, RepresentationSide: "explicit_formula_side", Sign: -1, Formula: "W_infinity(h) from the Gamma factor at the real place", Index: "infinite prime / real place"},
	}}
}

func weilExplicitFormula(q semantic.QuadraticFunctional, class semantic.FunctionClass) semantic.ExplicitFormulaIdentity {
	zero := *semantic.CloneQuadraticFunctional(q).Contributions[0].Aggregate
	return semantic.ExplicitFormulaIdentity{Functional: q.ID, FunctionClass: semantic.CloneFunctionClass(class), TransformConvention: semantic.LagariasMellinTransform(), ZeroSide: zero, ArithmeticSide: append([]semantic.FunctionalContribution(nil), q.Contributions[1:]...), Theorem: WeilExplicitFormulaTheoremID}
}
