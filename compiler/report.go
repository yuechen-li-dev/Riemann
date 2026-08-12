package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func HumanReport(result M0Result) string {
	var b strings.Builder
	root, _ := result.Graph.Claim(RHClaimID)
	fmt.Fprintf(&b, "TARGET\n  %s\n", root.Proposition.Describe())
	for _, t := range result.Graph.Transformations() {
		to, _ := result.Graph.Claim(t.To)
		fmt.Fprintf(&b, "\nLOWER\n  pass: %s\n  relation: %s\n", t.Pass, t.Relation)
		if len(t.Losses) == 0 {
			b.WriteString("  information loss: none\n")
		}
		for _, loss := range t.Losses {
			fmt.Fprintf(&b, "  information loss: %s — %s\n", loss.Kind, loss.Reason)
		}
		fmt.Fprintf(&b, "RESULT\n  %s\n", to.Proposition.Describe())
	}
	renderAttempt(&b, result.Graph, result.Attempt)
	return b.String()
}

func M1HumanReport(result M1Result) string {
	var b strings.Builder
	zeroFree, zeroFreeExists := result.Graph.Claim(ZeroFreeHalfPlaneID)
	targetDescription := semantic.QuantifiedStatement{Quantifier: semantic.ForAll, Domain: semantic.HalfPlaneReGreaterThanOne(), Predicate: semantic.Predicate{Kind: semantic.FunctionNonzeroPredicate, Function: semantic.RiemannZeta}}.Describe()
	if zeroFreeExists {
		targetDescription = zeroFree.Proposition.Describe()
	}
	fmt.Fprintf(&b, "TARGET\n  %s\n\n", targetDescription)
	for _, contract := range result.Registry.Contracts() {
		fmt.Fprintf(&b, "IMPORT THEOREM\n  schema: %s\n  relation: %s\n  trust: %s\n  source: %s\n\n", contract.ID, contract.Relation, contract.Trust, contract.Evidence.Source.Citation)
	}
	for _, application := range result.Applications {
		fmt.Fprintf(&b, "INSTANTIATE\n  theorem: %s\n", application.Theorem)
		if len(application.Bindings) == 0 {
			b.WriteString("  bindings: none\n")
		}
		for _, binding := range application.Bindings {
			fmt.Fprintf(&b, "  bind %s = %s\n", binding.Parameter, describeBinding(binding.Value))
		}
		b.WriteString("PREMISES\n")
		if len(application.Matched) == 0 {
			b.WriteString("  none (trusted import)\n")
		}
		for _, matched := range application.Matched {
			claim, _ := result.Graph.Claim(matched.Claim)
			fmt.Fprintf(&b, "  [%d] %s — satisfied by %s\n", matched.Premise, claim.Proposition.Describe(), matched.Claim)
		}
		for _, obligation := range application.Obligations {
			fmt.Fprintf(&b, "  [%d] %s — UNRESOLVED\n", obligation.Premise, obligation.Description)
		}
		if application.Complete {
			claim, _ := result.Graph.Claim(application.Conclusion)
			fmt.Fprintf(&b, "DERIVE\n  %s (%s)\n\n", claim.Proposition.Describe(), application.Conclusion)
		} else {
			b.WriteString("CANDIDATE\n  conclusion withheld until every premise is discharged\n\n")
		}
	}
	if result.ZeroFreeCertified {
		b.WriteString("CERTIFIED\n  relative to trusted imported theorem contracts\n")
	} else {
		b.WriteString("\nUNCERTIFIED\n")
		for _, d := range result.ZeroFreeDiagnostics {
			fmt.Fprintf(&b, "  [%s] %s\n", d.Code, d.Message)
		}
	}
	for _, attempt := range []ProofAttempt{result.BoundedToRH, result.DensityToRH, result.ZeroFreeToRH} {
		renderAttempt(&b, result.Graph, attempt)
	}
	return b.String()
}

func M3HumanReport(result M3Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M3 — FUNCTIONAL EQUATION + TYPED SYMMETRY TRANSPORT\n\n")
	for _, id := range []semantic.ClaimID{AnalyticContinuationID, CompletedXiID, XiFunctionalEquationID, ZetaConjugationConditionID, ZetaPoleID, ZetaTrivialZerosID, ZetaBoundaryZeroFreeID} {
		if claim, ok := result.Graph.Claim(id); ok {
			fmt.Fprintf(&b, "IMPORT\n  %s\n  source: %s\n\n", claim.Proposition.Describe(), claim.Provenance.Source.Citation)
		}
	}
	source, _ := result.Graph.Claim(SampleZeroID)
	fmt.Fprintf(&b, "SOURCE CLAIM\n  %s\n  classification: nontrivial zero\n\n", source.Proposition.Describe())
	for _, application := range result.Applications {
		if !application.Complete {
			continue
		}
		claim, ok := result.Graph.Claim(application.Conclusion)
		if !ok {
			continue
		}
		z, zero := claim.Proposition.(semantic.ZeroAtPoint)
		if !zero {
			continue
		}
		var sourcePoint string
		for _, binding := range application.Bindings {
			if binding.Value.Type == PointParam && binding.Parameter == "P" {
				sourcePoint = binding.Value.Point.Describe()
			}
		}
		fmt.Fprintf(&b, "TRANSPORT\n  theorem: %s\n  map: %s → %s\n", application.Theorem, sourcePoint, z.Point.Describe())
		for _, matched := range application.Matched {
			if matched.SideCondition {
				c, _ := result.Graph.Claim(matched.Claim)
				fmt.Fprintf(&b, "  side condition: %s — satisfied\n", c.Proposition.Describe())
			}
		}
		fmt.Fprintf(&b, "DERIVE\n  %s\n\n", claim.Proposition.Describe())
	}
	for _, application := range result.Applications {
		if application.Complete {
			continue
		}
		for _, obligation := range application.Obligations {
			if obligation.SideCondition && obligation.Proposition != nil {
				fmt.Fprintf(&b, "UNRESOLVED SIDE CONDITION\n  theorem: %s\n  %s\n  conclusion withheld\n\n", application.Theorem, obligation.Description)
			}
		}
	}
	b.WriteString("CLOSURE\n  four generated symmetry transforms:\n")
	for _, p := range result.Orbit.Generated {
		fmt.Fprintf(&b, "    %s\n", p.Describe())
	}
	fmt.Fprintf(&b, "  distinct zero locations after semantic deduplication: %d\n", len(result.Orbit.Distinct))
	for _, p := range result.Orbit.Distinct {
		fmt.Fprintf(&b, "    %s\n", p.Describe())
	}
	b.WriteString("\nGLOBAL GEOMETRY\n")
	if result.StripCertified {
		b.WriteString("  KNOWN: nontrivial zeros are confined to 0 < Re(s) < 1\n")
	} else {
		b.WriteString("  UNRESOLVED: critical-strip confinement\n")
	}
	if result.SymmetryCertified {
		b.WriteString("  KNOWN: the nontrivial-zero set is invariant under s → 1-conjugate(s)\n")
	} else {
		b.WriteString("  UNRESOLVED: critical-line-reflection symmetry\n")
	}
	b.WriteString("\nRH REMAINING OBLIGATION\n  TARGET: every nontrivial zero is fixed by critical-line reflection\n  equivalently: Re(ρ) = 1/2\n  unresolved structural defect: exclude off-axis symmetry orbits inside the critical strip\n\n")
	if result.StripCertified && result.SymmetryCertified {
		b.WriteString("CERTIFIED\n  global geometry is certified relative to trusted imported theorem contracts\n")
	}
	return b.String()
}

func M4HumanReport(result M4Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M4 — WEIL CRITERION + TEST-FUNCTION IR\n\n")
	rh, _ := result.Graph.Claim(ZeroLocationID)
	positivity, _ := result.Graph.Claim(WeilPositivityID)
	fmt.Fprintf(&b, "TARGET\n  %s\n\n", rh.Proposition.Describe())
	b.WriteString("LOWER\n  theorem: Weil positivity criterion (Lagarias, Theorem 3.2)\n  relation: equivalent\n  trust: trusted imported mathematics\n")
	b.WriteString("  normalization: M[f](s) = integral_0^infinity f(x)x^s dx/x\n")
	b.WriteString("  involution: tilde(g)(x) = x^-1 g(x^-1)\n")
	b.WriteString("  quadratic input: h = f * tilde(conjugate(f)), using multiplicative convolution\n")
	fmt.Fprintf(&b, "RESULT\n  %s\n\n", positivity.Proposition.Describe())

	definition, _ := result.Graph.Claim(WeilFunctionalDefinitionID)
	q := definition.Proposition.(semantic.FunctionalDefinition).Functional
	fmt.Fprintf(&b, "FUNCTIONAL\n  %s\n  contributions:\n", q.ID)
	for _, c := range q.Contributions {
		fmt.Fprintf(&b, "    %s [%s] (sign %+d): %s\n", c.Kind, c.RepresentationSide, c.Sign, c.Formula)
		if c.Aggregate != nil {
			fmt.Fprintf(&b, "      aggregate domain: %s\n      transform: %s\n      theorem lineage: %v\n", c.Aggregate.IndexDomain.Describe(), c.Aggregate.TransformConvention, c.Aggregate.TheoremLineage)
		}
	}
	b.WriteString("\nEXPLICIT-FORMULA TRUST BOUNDARY\n  definition: Q(f) is the zero-side quadratic functional W^(1)(f*tilde(conjugate(f)))\n  imported theorem: Wspec(h) = Warith(h) (Lagarias, Theorem 3.1)\n  derived decomposition: Q(f) = endpoint terms - prime-power terms - archimedean term\n\n")

	finite, _ := result.Graph.Claim(FiniteWeilPositivityID)
	fmt.Fprintf(&b, "RESTRICT\n  test-function domain:\n    %s\n  →\n    %s\n", positivity.Proposition.(semantic.UniversalFunctionalStatement).FunctionClass.Describe(), result.FiniteFamily.Describe())
	b.WriteString("  relation: implies\n  information loss: function_space_restriction — finite coverage is a strict subclass\n")
	fmt.Fprintf(&b, "RESULT\n  %s\n\n", finite.Proposition.Describe())

	b.WriteString("ATTEMPT\n  use finite-family positivity to certify RH\nRH CERTIFICATION\n  REJECTED\n")
	for _, d := range result.FiniteToRH.Diagnostics {
		fmt.Fprintf(&b, "  [%s] %s\n", d.Code, d.Message)
	}
	b.WriteString("\nNUMERICAL EVIDENCE BOUNDARY\n  sampled finite family {f1, f2}; hypothetical values recorded as nonnegative\nEVIDENCE\n  numerical only (approximate finite-family claim)\nRH CERTIFICATION\n  REJECTED\n")
	for _, d := range result.NumericalToUniversal.Diagnostics {
		fmt.Fprintf(&b, "  [%s] %s\n", d.Code, d.Message)
	}
	b.WriteString("\nOCT EXPERIMENT\n  path: experiments/m4_weil_toy.octest\n  setup: real reflection-pair block q(a,b)=2ab and critical fixed-point block q(a)=a^2\n  result: 3 passed; compiled 3; interpreted fallback 0; off-critical sample q(1,-1)=-2\n  scope: sign-convention toy only; no admissible function construction and no zeta-functional evaluation\n")
	b.WriteString("\nSTATUS\n  The RH ↔ universal Weil-positivity lowering is represented exactly relative to trusted imports.\n  Neither universal positivity nor RH is certified.\n")
	return b.String()
}

func M5HumanReport(result M5Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M5 — FINITE HERMITIAN LOWERING\n\n")
	b.WriteString("SOURCE\n  Weil quadratic functional Q_W\n  domain: Weil-nice test functions\n\n")
	b.WriteString("RESTRICT\n")
	fmt.Fprintf(&b, "  admissible function space\n  →\n  finite span V = span{%s, %s} over Complex\n", result.Span.Basis.Members[0].Function.Symbol, result.Span.Basis.Members[1].Function.Symbol)
	b.WriteString("  each ordered basis member carries an exact admissibility claim\n")
	b.WriteString("  information loss: function_space_restriction — finite-dimensional coverage cannot recover universal Weil coverage\n\n")
	b.WriteString("QUADRATIC PREREQUISITES\n  Q(lambda f)=|lambda|^2 Q(f)\n  parallelogram identity\n  real-valued diagonal\n  missing prerequisites would remain obligations\n\n")
	b.WriteString("POLARIZE\n")
	fmt.Fprintf(&b, "  Q_W → Hermitian form %s\n", result.Form.ID)
	b.WriteString("  convention: conjugate-linear in argument 1; linear in argument 2\n")
	fmt.Fprintf(&b, "  normalization: %s\n  formula: %s\n", result.Form.Convention.Normalization, result.Form.Convention.Formula)
	b.WriteString("  identities: B(f,g)=conjugate(B(g,f)); Q(f)=B(f,f)\n\n")
	b.WriteString("BASIS LOWERING\n")
	fmt.Fprintf(&b, "  ordered basis: %s\n  G_ij = %s(f_i,f_j)\n  matrix identity: Q_W(%s) = c* G c\n", result.Matrix.Basis.ID, result.Form.ID, result.Combination.Describe())
	fmt.Fprintf(&b, "  matrix value semantics: %s\n  Hermitian: certified from construction\n  positive semidefinite: separate open finite claim\n", result.Matrix.ValueSemantics)
	b.WriteString("  entry contributions retained:\n")
	for _, contribution := range result.Matrix.Entries[0].Contributions {
		fmt.Fprintf(&b, "    %s [%s] sign %+d\n", contribution.SourceKind, contribution.RepresentationSide, contribution.Sign)
	}
	b.WriteString("\nFINITE CLAIM\n  positivity of Q_W on V\n\nEQUIVALENT\n  for every c in Complex^2: c* G c >= 0\n  ⇔\n  G is positive semidefinite\n  theorem: finite-hermitian-psd-equivalence\n\n")
	b.WriteString("SOUNDNESS REJECTIONS\n  basis-point positivity → span positivity: REJECTED\n")
	for _, diagnostic := range result.FamilyToSpan.Diagnostics {
		fmt.Fprintf(&b, "    [%s] %s\n", diagnostic.Code, diagnostic.Message)
	}
	b.WriteString("  nonnegative diagonal → PSD: REJECTED\n")
	for _, diagnostic := range result.DiagonalToPSD.Diagnostics {
		fmt.Fprintf(&b, "    [%s] %s\n", diagnostic.Code, diagnostic.Message)
	}
	b.WriteString("  approximate numerical PSD → exact PSD: REJECTED\n")
	for _, diagnostic := range result.ApproximateToExactPSD.Diagnostics {
		fmt.Fprintf(&b, "    [%s] %s\n", diagnostic.Code, diagnostic.Message)
	}
	b.WriteString("\nRH CERTIFICATION\n  unavailable: finite function-space coverage only\n")
	for _, diagnostic := range result.MatrixPSDToRH.Diagnostics {
		fmt.Fprintf(&b, "  [%s] %s\n", diagnostic.Code, diagnostic.Message)
	}
	b.WriteString("\nEXPERIMENT\n")
	b.WriteString("  backend: Oct dev, local compiled test runner\n")
	b.WriteString("  path: experiments/m5_polarization_matrix.octest\n")
	b.WriteString("  command: oct test <path> --execution compiled\n")
	b.WriteString("  inputs: exact decimal toy coordinates for G=[[1,1+i],[1-i,3]] and G=[[0,1],[1,0]]\n")
	b.WriteString("  outputs: 6 passed; compiled 6; interpreted fallback 0; adversarial toy coefficient (1,-1) gives -2\n")
	b.WriteString("  numerical precision: Float backend; exact representable inputs here, but no interval/error certification\n")
	b.WriteString("  evidence: numerical only; convention and counterexample probe, not Weil-entry evaluation or proof\n")
	b.WriteString("  Octxiliary: not used; direct local Oct was sufficient and no transport payload entered theorem state\n")
	b.WriteString("  when utility: not used; the polarization sign probe was the obvious highest-information experiment\n\n")
	b.WriteString("STATUS\n  Finite Weil-functional reasoning now has a faithful structural Hermitian matrix lowering.\n  Universal Weil positivity and RH remain uncertified.\n")
	return b.String()
}

func describeBinding(v BindingValue) string {
	switch v.Type {
	case ObjectParam:
		return v.Object.String()
	case DomainParam:
		return v.Domain.Describe()
	case RepresentationParam:
		return string(v.Representation)
	case ScalarParam:
		return fmt.Sprint(v.Scalar)
	case QuantifierParam:
		return string(v.Quantifier)
	case PredicateParam:
		return string(v.Predicate)
	case AnalyticFactParam:
		return string(v.AnalyticFact)
	case PointParam:
		return v.Point.Describe()
	case ZeroClassParam:
		return string(v.ZeroClass)
	case SideConditionParam:
		return string(v.SideCondition)
	case TransformParam:
		return v.Transform.String()
	case ZeroSetPropertyParam:
		return string(v.ZeroSetProperty)
	case ZeroClassificationParam:
		return string(v.ZeroClassification)
	default:
		return "unknown"
	}
}

func renderAttempt(b *strings.Builder, g *Graph, attempt ProofAttempt) {
	from, fromOK := g.Claim(attempt.From)
	target, targetOK := g.Claim(attempt.Target)
	b.WriteString("\nATTEMPT\n")
	fromText, targetText := string(attempt.From), string(attempt.Target)
	if fromOK {
		fromText = from.Proposition.Describe()
	}
	if targetOK {
		targetText = target.Proposition.Describe()
	}
	fmt.Fprintf(b, "  use %q to discharge %q\n", fromText, targetText)
	if attempt.Accepted {
		b.WriteString("ACCEPTED\n")
		return
	}
	b.WriteString("REJECTED\n")
	for _, d := range attempt.Diagnostics {
		fmt.Fprintf(b, "  [%s] %s\n", d.Code, d.Message)
	}
}

type graphJSON struct {
	Schema          string               `json:"schema"`
	Contracts       []TheoremContract    `json:"theorem_contracts"`
	Applications    []TheoremApplication `json:"theorem_applications"`
	Claims          []claimJSON          `json:"claims"`
	Transformations []transformationJSON `json:"transformations"`
	Certifications  []certificationJSON  `json:"certifications"`
	Attempts        []ProofAttempt       `json:"attempts"`
}

type certificationJSON struct {
	Claim       semantic.ClaimID `json:"claim"`
	Certified   bool             `json:"certified"`
	Diagnostics []Diagnostic     `json:"diagnostics"`
}
type claimJSON struct {
	ID          semantic.ClaimID      `json:"id"`
	Proposition propositionJSON       `json:"proposition"`
	Assumptions []semantic.Assumption `json:"assumptions"`
	Evidence    []semantic.Evidence   `json:"evidence"`
	Exactness   semantic.Exactness    `json:"exactness"`
	Provenance  semantic.Provenance   `json:"provenance"`
}
type propositionJSON struct {
	Kind                          semantic.PropositionKind                `json:"kind"`
	Description                   string                                  `json:"description"`
	Quantifier                    semantic.QuantifierKind                 `json:"quantifier,omitempty"`
	Domain                        *semantic.Domain                        `json:"domain,omitempty"`
	Predicate                     *semantic.Predicate                     `json:"predicate,omitempty"`
	Representation                *semantic.Representation                `json:"representation,omitempty"`
	RepresentationIdentity        *semantic.RepresentationIdentity        `json:"representation_identity,omitempty"`
	AnalyticFact                  *semantic.AnalyticFact                  `json:"analytic_fact,omitempty"`
	ZeroAtPoint                   *semantic.ZeroAtPoint                   `json:"zero_at_point,omitempty"`
	SideCondition                 *semantic.SideCondition                 `json:"side_condition,omitempty"`
	FunctionalIdentity            *semantic.FunctionalIdentity            `json:"functional_identity,omitempty"`
	ZeroSetProperty               *semantic.ZeroSetProperty               `json:"zero_set_property,omitempty"`
	ZeroClassification            *semantic.ZeroClassification            `json:"zero_classification,omitempty"`
	FunctionalDefinition          *semantic.FunctionalDefinition          `json:"functional_definition,omitempty"`
	UniversalFunctionalStatement  *semantic.UniversalFunctionalStatement  `json:"universal_functional_statement,omitempty"`
	TestFunctionAdmissibility     *semantic.TestFunctionAdmissibility     `json:"test_function_admissibility,omitempty"`
	ExplicitFormulaIdentity       *semantic.ExplicitFormulaIdentity       `json:"explicit_formula_identity,omitempty"`
	FiniteSpanDefinition          *semantic.FiniteSpanDefinition          `json:"finite_span_definition,omitempty"`
	QuadraticFormStructure        *semantic.QuadraticFormStructure        `json:"quadratic_form_structure,omitempty"`
	HermitianFormDefinition       *semantic.HermitianFormDefinition       `json:"hermitian_form_definition,omitempty"`
	HermitianMatrixDefinition     *semantic.HermitianMatrixDefinition     `json:"hermitian_matrix_definition,omitempty"`
	FiniteSpanFunctionalStatement *semantic.FiniteSpanFunctionalStatement `json:"finite_span_functional_statement,omitempty"`
	CoordinateQuadraticPositivity *semantic.CoordinateQuadraticPositivity `json:"coordinate_quadratic_positivity,omitempty"`
	QuadraticMatrixIdentity       *semantic.QuadraticMatrixIdentity       `json:"quadratic_matrix_identity,omitempty"`
	MatrixProperty                *semantic.MatrixProperty                `json:"matrix_property,omitempty"`
}
type transformationJSON struct {
	ID          semantic.TransformationID `json:"id"`
	Pass        string                    `json:"pass"`
	From        semantic.ClaimID          `json:"from"`
	Premises    []semantic.ClaimID        `json:"premises"`
	To          semantic.ClaimID          `json:"to"`
	Relation    Relation                  `json:"relation"`
	Obligations []semantic.ClaimID        `json:"obligations"`
	Losses      []InformationLoss         `json:"losses"`
	Provenance  semantic.Reference        `json:"provenance"`
	Theorem     semantic.TheoremID        `json:"theorem,omitempty"`
	Bindings    []Binding                 `json:"bindings"`
	Trusted     bool                      `json:"trusted_theorem,omitempty"`
}

func JSONReport(result M0Result) ([]byte, error) {
	return marshalGraph("riemann.semantic-graph.m0", result.Graph, nil, []ProofAttempt{result.Attempt})
}
func M1JSONReport(result M1Result) ([]byte, error) {
	certifications := []certificationJSON{{Claim: ZeroFreeHalfPlaneID, Certified: result.ZeroFreeCertified, Diagnostics: nonNil(result.ZeroFreeDiagnostics)}}
	return marshalGraphWithContracts("riemann.semantic-graph.m2", result.Graph, result.Registry.Contracts(), result.Applications, certifications, []ProofAttempt{result.BoundedToRH, result.DensityToRH, result.ZeroFreeToRH})
}

func M3JSONReport(result M3Result) ([]byte, error) {
	contracts := append(result.M1.Registry.Contracts(), result.Registry.Contracts()...)
	sort.Slice(contracts, func(i, j int) bool { return contracts[i].ID < contracts[j].ID })
	applications := append(append([]TheoremApplication(nil), result.M1.Applications...), result.Applications...)
	certifications := []certificationJSON{{Claim: CriticalReflectionInvariantID, Certified: result.SymmetryCertified, Diagnostics: nonNil(result.SymmetryDiagnostics)}, {Claim: CriticalStripConfinementID, Certified: result.StripCertified, Diagnostics: nonNil(result.StripDiagnostics)}}
	return marshalGraphWithContracts("riemann.semantic-graph.m3", result.Graph, contracts, applications, certifications, []ProofAttempt{result.M1.BoundedToRH, result.M1.DensityToRH, result.M1.ZeroFreeToRH})
}

func M4JSONReport(result M4Result) ([]byte, error) {
	contracts := append(result.M3.M1.Registry.Contracts(), result.M3.Registry.Contracts()...)
	contracts = append(contracts, result.Registry.Contracts()...)
	sort.Slice(contracts, func(i, j int) bool { return contracts[i].ID < contracts[j].ID })
	applications := append(append([]TheoremApplication(nil), result.M3.M1.Applications...), result.M3.Applications...)
	positivityCertified, positivityDiagnostics := result.Graph.Certify(WeilPositivityID)
	finiteCertified, finiteDiagnostics := result.Graph.Certify(FiniteWeilPositivityID)
	certifications := []certificationJSON{
		{Claim: CriticalReflectionInvariantID, Certified: result.M3.SymmetryCertified, Diagnostics: nonNil(result.M3.SymmetryDiagnostics)},
		{Claim: CriticalStripConfinementID, Certified: result.M3.StripCertified, Diagnostics: nonNil(result.M3.StripDiagnostics)},
		{Claim: WeilPositivityID, Certified: positivityCertified, Diagnostics: nonNil(positivityDiagnostics)},
		{Claim: FiniteWeilPositivityID, Certified: finiteCertified, Diagnostics: nonNil(finiteDiagnostics)},
	}
	return marshalGraphWithContracts("riemann.semantic-graph.m4", result.Graph, contracts, applications, certifications, []ProofAttempt{result.FullToFinite, result.FiniteToRH, result.NumericalToUniversal})
}

func M5JSONReport(result M5Result) ([]byte, error) {
	contracts := append(result.M4.M3.M1.Registry.Contracts(), result.M4.M3.Registry.Contracts()...)
	contracts = append(contracts, result.M4.Registry.Contracts()...)
	contracts = append(contracts, result.Registry.Contracts()...)
	sort.Slice(contracts, func(i, j int) bool { return contracts[i].ID < contracts[j].ID })
	applications := append(append([]TheoremApplication(nil), result.M4.M3.M1.Applications...), result.M4.M3.Applications...)
	spanCertified, spanDiagnostics := result.Graph.Certify(M5SpanPositivityID)
	psdCertified, psdDiagnostics := result.Graph.Certify(M5MatrixPSDID)
	certifications := []certificationJSON{
		{Claim: M5HermitianPropertyID, Certified: result.HermitianCertified, Diagnostics: nonNil(result.HermitianDiagnostics)},
		{Claim: M5SpanPositivityID, Certified: spanCertified, Diagnostics: nonNil(spanDiagnostics)},
		{Claim: M5MatrixPSDID, Certified: psdCertified, Diagnostics: nonNil(psdDiagnostics)},
	}
	attempts := []ProofAttempt{result.FullToSpan, result.FamilyToSpan, result.DiagonalToPSD, result.ApproximateToExactPSD, result.MatrixPSDToRH}
	return marshalGraphWithContracts("riemann.semantic-graph.m5", result.Graph, contracts, applications, certifications, attempts)
}

func marshalGraph(schema string, g *Graph, certifications []certificationJSON, attempts []ProofAttempt) ([]byte, error) {
	return marshalGraphWithContracts(schema, g, nil, nil, certifications, attempts)
}

func marshalGraphWithContracts(schema string, g *Graph, contracts []TheoremContract, applications []TheoremApplication, certifications []certificationJSON, attempts []ProofAttempt) ([]byte, error) {
	report := graphJSON{Schema: schema, Contracts: nonNil(contracts), Applications: nonNil(applications), Certifications: nonNil(certifications), Attempts: nonNil(attempts)}
	for _, claim := range g.Claims() {
		item := propositionJSON{Kind: claim.Proposition.Kind(), Description: claim.Proposition.Describe()}
		switch p := claim.Proposition.(type) {
		case semantic.QuantifiedStatement:
			item.Quantifier = p.Quantifier
			item.Domain = &p.Domain
			item.Predicate = &p.Predicate
		case semantic.RepresentationProposition:
			r := semantic.CloneRepresentation(p.Representation)
			item.Representation = &r
		case semantic.RepresentationIdentity:
			item.RepresentationIdentity = &p
		case semantic.AnalyticFact:
			item.AnalyticFact = &p
		case semantic.ZeroAtPoint:
			item.ZeroAtPoint = &p
		case semantic.SideCondition:
			item.SideCondition = &p
		case semantic.FunctionalIdentity:
			item.FunctionalIdentity = &p
		case semantic.ZeroSetProperty:
			item.ZeroSetProperty = &p
		case semantic.ZeroClassification:
			item.ZeroClassification = &p
		case semantic.FunctionalDefinition:
			item.FunctionalDefinition = &p
		case semantic.UniversalFunctionalStatement:
			item.UniversalFunctionalStatement = &p
		case semantic.TestFunctionAdmissibility:
			item.TestFunctionAdmissibility = &p
		case semantic.ExplicitFormulaIdentity:
			item.ExplicitFormulaIdentity = &p
		case semantic.FiniteSpanDefinition:
			item.FiniteSpanDefinition = &p
		case semantic.QuadraticFormStructure:
			item.QuadraticFormStructure = &p
		case semantic.HermitianFormDefinition:
			item.HermitianFormDefinition = &p
		case semantic.HermitianMatrixDefinition:
			item.HermitianMatrixDefinition = &p
		case semantic.FiniteSpanFunctionalStatement:
			item.FiniteSpanFunctionalStatement = &p
		case semantic.CoordinateQuadraticPositivity:
			item.CoordinateQuadraticPositivity = &p
		case semantic.QuadraticMatrixIdentity:
			item.QuadraticMatrixIdentity = &p
		case semantic.MatrixProperty:
			item.MatrixProperty = &p
		}
		report.Claims = append(report.Claims, claimJSON{ID: claim.ID, Proposition: item, Assumptions: nonNil(claim.Assumptions), Evidence: nonNil(claim.Evidence), Exactness: claim.Exactness, Provenance: claim.Provenance})
	}
	for _, t := range g.Transformations() {
		report.Transformations = append(report.Transformations, transformationJSON{ID: t.ID, Pass: t.Pass, From: t.From, Premises: nonNil(t.Premises), To: t.To, Relation: t.Relation, Obligations: nonNil(t.Obligations), Losses: nonNil(t.Losses), Provenance: t.Provenance, Theorem: t.Theorem, Bindings: nonNil(t.Bindings), Trusted: t.Trusted})
	}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(report); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func nonNil[T any](values []T) []T {
	if values == nil {
		return []T{}
	}
	return values
}

func M6HumanReport(result M6Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M6 — TYPED WEIL MATRIX ENTRY EVALUATION\n\nBASIS\n")
	b.WriteString("  V = span{f2, f3} over Complex\n")
	b.WriteString("  f_q(x) = x^-1/2 on [q^-2,q^2], zero outside; midpoint values at endpoints\n")
	b.WriteString("  admissibility: certified from Lagarias §3 nice-function definition by direct constraint check\n")
	b.WriteString("  transform: M[f](s)=integral_0^infinity f(x)x^s dx/x\n")
	b.WriteString("  backend: deterministic Go float64 evaluator (53-bit precision); Oct is an independent experimental probe\n")
	for _, entry := range result.Evaluation.Matrix.Entries {
		fmt.Fprintf(&b, "\nENTRY G[%d,%d]\n  definition: %s (exact, independently of value evaluation)\n", entry.Row+1, entry.Column+1, entry.Definition)
		for _, contribution := range entry.Contributions {
			fmt.Fprintf(&b, "  %s (sign %+d): %s", contribution.SourceKind, contribution.Sign, contribution.Value.Kind)
			switch contribution.Value.Kind {
			case semantic.ExactValue:
				fmt.Fprintf(&b, " value=%s", contribution.Value.Exact.Real.Expression)
			case semantic.ApproximateValue:
				fmt.Fprintf(&b, " value=%.12g", contribution.Value.Approximate.Real)
				if contribution.Value.Metadata.Truncation != nil {
					fmt.Fprintf(&b, " truncation={%s; remainder=%s}", contribution.Value.Metadata.Truncation.Bound, contribution.Value.Metadata.Truncation.RemainderStatus)
				}
				if contribution.Value.Metadata.Quadrature != nil {
					fmt.Fprintf(&b, " quadrature={%s; tolerance=%g; rigorous=%t}", contribution.Value.Metadata.Quadrature.Method, contribution.Value.Metadata.Quadrature.Tolerance, contribution.Value.Metadata.Quadrature.ErrorRigorous)
				}
			}
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "  total: %s value=%.12g\n", entry.Value.Kind, entry.Value.Approximate.Real)
		b.WriteString("  zero-side: retained as unevaluated theorem-linked alternate representation\n")
	}
	b.WriteString("\nMATRIX\n  semantic structure: exact Hermitian form (certified)\n  numerical realization: approximate\n  values:\n")
	m := result.Evaluation.Matrix
	for i := 0; i < m.Rows; i++ {
		b.WriteString("    [")
		for j := 0; j < m.Columns; j++ {
			if j > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%.12g", m.Entries[i*m.Columns+j].Value.Approximate.Real)
		}
		b.WriteString("]\n")
	}
	maxSymmetry := 0.0
	for _, d := range result.Evaluation.HermitianDiagnostics {
		if d.Discrepancy > maxSymmetry {
			maxSymmetry = d.Discrepancy
		}
	}
	fmt.Fprintf(&b, "  Hermitian consistency max |G_ij-conj(G_ji)|: %.3g (diagnostic only)\n", maxSymmetry)
	for i, check := range result.Evaluation.DirectMatrixChecks {
		fmt.Fprintf(&b, "  direct-vs-matrix probe %d discrepancy: %.3g (tolerance %.3g; numerical only)\n", i+1, check.Discrepancy, check.Tolerance)
	}
	fmt.Fprintf(&b, "  approximate eigenvalues: [%.12g, %.12g]; condition number: %.12g\n", result.Evaluation.EigenDiagnostic.Eigenvalues[0], result.Evaluation.EigenDiagnostic.Eigenvalues[1], result.Evaluation.EigenDiagnostic.Condition)
	b.WriteString("  PSD: exact claim remains open; approximate eigenvalues do not certify it\n")
	b.WriteString("\nRH CERTIFICATION\n  unchanged: unavailable\n  finite function_space_restriction remains on the proof path\n")
	b.WriteString("\nRESEARCH BOUNDARIES\n  Octxiliary: not used; direct local Oct invocation was smaller\n  when utility: not used; the centered log-box route was clearly preferable after source inspection\n  zero-side numerical cross-check: not performed\n")
	return b.String()
}

func M6JSONReport(result M6Result) ([]byte, error) {
	proof, err := M5JSONReport(result.M5)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema              string                `json:"schema"`
		ProofGraph          json.RawMessage       `json:"proof_graph"`
		EvaluationContracts []TheoremContract     `json:"evaluation_theorem_contracts"`
		Basis               semantic.OrderedBasis `json:"basis"`
		Evaluation          MatrixEvaluation      `json:"evaluation"`
	}{Schema: "riemann.semantic-graph.m6", ProofGraph: json.RawMessage(proof), EvaluationContracts: result.Registry.Contracts(), Basis: result.Evaluation.Matrix.Basis, Evaluation: result.Evaluation}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(report); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
