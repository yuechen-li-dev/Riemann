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

func M7HumanReport(result M7Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M7 — CERTIFIED ARCHIMEDEAN ENCLOSURE + RIGOROUS FINITE WEIL MATRIX\n\n")
	b.WriteString("BASIS\n  ordered: (f2, f3)\n  f_q(x)=x^-1/2 on [q^-2,q^2], zero outside, with midpoint endpoint values\n  admissibility: certified; unchanged from M6\n\n")
	b.WriteString("EXACT ARCHIMEDEAN DEFINITION\n  W_inf(h)=(EulerGamma+log(pi))*h(1)+integral_1^infinity [h(x)+tilde(h)(x)-2*x^-2*h(1)]*x*dx/(x^2-1)\n  u=log(x): H(u)=min(4log(min(q,r)), max(0,2log(qr)-u)) with the plateau starting at 2|log(q/r)|\n  integrand: F(u)=2*[exp(3u/2)*H(u)-h(1)]/expm1(2u)\n  cancellation-safe at zero: R=(3/4)*phi(3u/2)/phi(2u), S=(1/2)*exp(3u/2)/phi(2u), phi(z)=expm1(z)/z\n  singularity: removable at u=0; derivative kinks are isolated by certified breakpoint guards\n  tail after B=5: h(1)*log(1-exp(-2B)) (closed form, certified)\n\n")
	for _, entry := range result.Evaluation.Matrix.Entries {
		fmt.Fprintf(&b, "ENTRY G[%d,%d]\n", entry.Row+1, entry.Column+1)
		for _, contribution := range entry.Contributions {
			if contribution.SourceKind == semantic.ZeroContribution {
				continue
			}
			fmt.Fprintf(&b, "  %s (sign %+d): ", contribution.SourceKind, contribution.Sign)
			if contribution.Value.Kind == semantic.ExactValue {
				fmt.Fprintf(&b, "exact %s", contribution.Value.Exact.Real.Expression)
			} else if contribution.Value.Interval != nil {
				fmt.Fprintf(&b, "certified [%.15g, %.15g] certificate=%s", contribution.Value.Interval.RealLower, contribution.Value.Interval.RealUpper, contribution.Value.Metadata.Error.ProofObjectKind)
				if contribution.Value.Metadata.Truncation != nil {
					fmt.Fprintf(&b, " support_exhaustive=%t", contribution.Value.Metadata.Truncation.SupportExhaustive)
				}
				if contribution.Value.Metadata.Quadrature != nil {
					fmt.Fprintf(&b, " partitions=%d", contribution.Value.Metadata.Quadrature.Partitions)
				}
			}
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "  total: certified [%.15g, %.15g]\n  M6 approximate: %.15g; contained=%t\n  zero-side: unevaluated exact alternate representation\n\n", entry.Value.Interval.RealLower, entry.Value.Interval.RealUpper, result.Evaluation.ApproximateMatrix.Entries[entry.Row*2+entry.Column].Value.Approximate.Real, result.Evaluation.ApproximateContained[entry.Row*2+entry.Column])
	}
	b.WriteString("MATRIX\n  structural Hermitian: certified independently of interval overlap\n  entries: certified intervals\n\n")
	p := result.Evaluation.PSD
	fmt.Fprintf(&b, "FINITE PSD\n  a>0: %t\n  d>0: %t\n  determinant interval: [%.15g, %.15g]\n  ad-b^2>0: %t\n  finite matrix PSD: %s\n\n", p.APositive, p.DPositive, p.Determinant.Interval.RealLower, p.Determinant.Interval.RealUpper, p.DeterminantPositive, certifiedWord(p.Certified))
	fmt.Fprintf(&b, "FINITE WEIL SPAN\n  positivity on span{f2,f3}: %s\n  theorem path: certified matrix PSD + Q(sum c_i f_i)=c*Gc\n\n", certifiedWord(result.Evaluation.FiniteSpanPositivityCertified))
	b.WriteString("RH\n  unresolved\n  universal Weil positivity: unresolved\n  reason: function_space_restriction remains permanently on the universal-to-finite path; finite-span positivity cannot discharge RH\n\n")
	b.WriteString("RESEARCH BOUNDARIES\n  approximate M6 path: retained\n  Oct experiment: none in M7; the Go theorem-backed backend is authoritative\n  Octxiliary: not needed\n  when utility: not used; cancellation-safe interval Darboux decomposition was the smallest direct rigorous method\n  zero-side decomposition, inertia, rank, and orbit machinery: not begun\n")
	return b.String()
}

func certifiedWord(v bool) string {
	if v {
		return "certified"
	}
	return "unresolved"
}

func M7JSONReport(result M7Result) ([]byte, error) {
	proof, err := M6JSONReport(result.M6)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema              string                `json:"schema"`
		ProofGraph          json.RawMessage       `json:"proof_graph"`
		EvaluationContracts []TheoremContract     `json:"evaluation_theorem_contracts"`
		Basis               semantic.OrderedBasis `json:"basis"`
		Evaluation          M7Evaluation          `json:"evaluation"`
	}{Schema: "riemann.semantic-graph.m7", ProofGraph: json.RawMessage(proof), EvaluationContracts: result.Registry.Contracts(), Basis: result.Evaluation.Matrix.Basis, Evaluation: result.Evaluation}
	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(report); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func M8HumanReport(result M8Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M8 — ZERO-SIDE ORBIT DECOMPOSITION\n\n")
	b.WriteString("ZERO-SIDE LOWERING\n")
	b.WriteString("  exact M4 summand: M[f](rho) * conjugate(M[f](1-conjugate(rho)))\n")
	b.WriteString("  convention: conjugate-linear first, linear second\n")
	b.WriteString("  v(p)_i = M[f_i](p)\n")
	b.WriteString("  K(p)_ij = conjugate(v(1-conjugate(p))_i) * v(p)_j\n")
	b.WriteString("  identity: c* K(p) c = (c^T v(p))*conjugate(c^T v(1-conjugate(p)))\n\n")
	writeOrbit := func(title string, contribution semantic.OrbitMatrixContribution) {
		fmt.Fprintf(&b, "%s\n  representative: %s\n  class: %s\n  distinct geometric locations: %d\n  zero multiplicity: %d (separate from location count)\n", title, contribution.Orbit.Representative.Describe(), contribution.Orbit.Classification, contribution.Orbit.DistinctLocationCount, contribution.Orbit.ZeroMultiplicity)
		for _, point := range contribution.Orbit.TransformedPoints {
			fmt.Fprintf(&b, "    %s\n", point.Point.Describe())
		}
		fmt.Fprintf(&b, "  grouping: %d critical-reflection pair(s)\n", len(contribution.ReflectionPairs))
		for _, pair := range contribution.ReflectionPairs {
			fmt.Fprintf(&b, "    %s\n      form: %s\n      Hermitian certified: %t\n      PSD certified: %t\n      rank <= %d\n", pair.ID, pair.Formula, pair.Classification.Hermitian, pair.Classification.PositiveSemidefinite, pair.Classification.RankUpperBound)
			if pair.Classification.RankOneIfNonzero {
				b.WriteString("      rank 1 only if the evaluation vector is nonzero\n")
			}
			if pair.Classification.IndefiniteCondition != "" {
				fmt.Fprintf(&b, "      indefinite: conditional — %s\n      degenerate: %s\n", pair.Classification.IndefiniteCondition, pair.Classification.DegenerateCondition)
			}
		}
		fmt.Fprintf(&b, "  full orbit: Hermitian certified=%t, PSD certified=%t, rank<=%d\n\n", contribution.Classification.Hermitian, contribution.Classification.PositiveSemidefinite, contribution.Classification.RankUpperBound)
	}
	writeOrbit("CRITICAL-LINE ORBIT TEMPLATE", result.CriticalTemplate)
	writeOrbit("OFF-CRITICAL ORBIT TEMPLATE", result.OffCriticalTemplate)
	b.WriteString("ZERO-SIDE MATRIX\n")
	fmt.Fprintf(&b, "  %s\n  %s\n  %s\n  %s\n  symmetric limiting convention: %v\n\n", result.ZeroSide.Formula, result.ZeroSide.CriticalAggregate, result.ZeroSide.OffCriticalAggregate, result.ZeroSide.SplitIdentity, result.ZeroSide.SummationConvention)
	b.WriteString("DUAL REPRESENTATION\n")
	fmt.Fprintf(&b, "  same semantic matrix: %s\n  zero side: %s (symbolic orbit aggregate)\n  explicit-formula side: %s (%s)\n  identity theorem: %s\n  numerical identification used: %t\n\n", result.Dual.SemanticMatrixID, result.Dual.ZeroSideAggregateID, result.Dual.ExplicitFormulaMatrixID, result.Dual.ExplicitValueEvidence, result.Dual.IdentityTheorem, result.Dual.NumericalIdentification)
	b.WriteString("SYNTHETIC DIAGNOSTICS\n")
	for _, d := range result.ToyDiagnostics {
		fmt.Fprintf(&b, "  %s\n    inputs: %s\n    determinant: %.12g\n    eigenvalues: %v\n    evidence: %s\n", d.Name, d.Inputs, d.Determinant, d.Eigenvalues, d.Classification)
	}
	b.WriteString("\nSOUNDNESS\n")
	fmt.Fprintf(&b, "  finite M7 PSD => absence of off-critical zeros: accepted=%t\n", result.FinitePSDReverseProof.Accepted)
	for _, d := range result.FinitePSDReverseProof.Diagnostics {
		fmt.Fprintf(&b, "  diagnostic: %s — %s\n", d.Code, d.Message)
	}
	b.WriteString("  A PSD sum does not force every orbit summand to be PSD; compensation remains possible.\n\n")
	b.WriteString("STATUS\n  RH unresolved. M8 derives local zero-orbit matrix structure without adding aggregate inertia/rank machinery.\n")
	return b.String()
}

func M8JSONReport(result M8Result) ([]byte, error) {
	m7, err := M7JSONReport(result.M7)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema                string                            `json:"schema"`
		M7                    json.RawMessage                   `json:"m7"`
		CriticalOrbitTemplate semantic.OrbitMatrixContribution  `json:"critical_orbit_template"`
		OffCriticalTemplate   semantic.OrbitMatrixContribution  `json:"off_critical_orbit_template"`
		ZeroSideMatrix        semantic.ZeroSideMatrixAggregate  `json:"zero_side_matrix"`
		DualRepresentation    semantic.DualMatrixRepresentation `json:"dual_representation"`
		ToyDiagnostics        []ToyOrbitDiagnostic              `json:"toy_diagnostics"`
		ReverseInference      ProofAttempt                      `json:"finite_psd_reverse_inference"`
	}{"riemann.semantic-graph.m8", json.RawMessage(m7), result.CriticalTemplate, result.OffCriticalTemplate, result.ZeroSide, result.Dual, result.ToyDiagnostics, result.FinitePSDReverseProof}
	return json.MarshalIndent(report, "", "  ")
}

func M9HumanReport(result M9Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M9 — AGGREGATE SPECTRAL ACCOUNTING\n\n")
	b.WriteString("FINITE COMPRESSION\n")
	fmt.Fprintf(&b, "  matrix: %s\n  basis: %s\n  dimension: %d\n  information loss: function_space_restriction\n\n", result.Compression.MatrixID, result.Compression.BasisID, result.Compression.Dimension)
	b.WriteString("ZERO-SIDE SPECTRAL BUDGET\n\n")
	b.WriteString("G = P + Q\n\n")
	b.WriteString("P:\n")
	fmt.Fprintf(&b, "  source: critical-line orbit aggregate\n  PSD: %t\n  local rank <= %d\n  aggregate: %s\n  multiplicity: %s\n\n", result.CriticalBudget.PositiveSemidefinite, result.CriticalBudget.LocalRankUpperBound, result.CriticalBudget.RankUpperBound, result.CriticalBudget.MultiplicityEffect)
	b.WriteString("Q:\n")
	fmt.Fprintf(&b, "  source: off-critical reflection-pair aggregate\n  each block: rank <= %d, positive index <= %d, negative index <= %d\n  aggregate rank: %s\n  aggregate positive index: %s\n  aggregate negative index: %s\n  equality across blocks claimed: false\n\n", result.OffCriticalBudget.LocalRankUpperBound, result.OffCriticalBudget.LocalPositiveIndexUpperBound, result.OffCriticalBudget.LocalNegativeIndexUpperBound, result.OffCriticalBudget.AggregateRankUpperBound, result.OffCriticalBudget.AggregatePositiveIndexBound, result.OffCriticalBudget.AggregateNegativeIndexBound)
	b.WriteString("AGGREGATE THEOREM\n")
	fmt.Fprintf(&b, "  using: %v\n  derive: %s\n\n", result.DerivedTheorem.Theorems, result.DerivedTheorem.PositiveIndexInequality)
	b.WriteString("CRITICAL CONTRIBUTION BOUND\n")
	fmt.Fprintf(&b, "  %s\n  therefore: %s\n  M7 sanity case: %s\n\n", result.DerivedTheorem.CriticalRankLowerBound, result.DerivedTheorem.CriticalCountConsequence, result.DerivedTheorem.M7SanityInstance)
	b.WriteString("REPRESENTATION FUSION\n")
	fmt.Fprintf(&b, "  zero side supplies: %v\n  explicit-formula side supplies: %v\n  identity theorem: %s\n\n", result.Fusion.ZeroSideFacts, result.Fusion.ExplicitFormulaFacts, result.Fusion.IdentityTheorem)
	b.WriteString("COUNTEREXAMPLES TO STRONGER CANDIDATES\n")
	for _, c := range result.Counterexamples {
		fmt.Fprintf(&b, "  rejected: %s\n    fixture: %s\n    reason: %s\n", c.RejectedCandidate, c.ExactFixture, c.Reason)
	}
	b.WriteString("\nOCT EXPERIMENT\n")
	fmt.Fprintf(&b, "  path: %s\n  command: %s\n  setup: %s\n  trials: %d\n  result: %s\n\n", result.Experiment.Path, result.Experiment.Command, result.Experiment.Setup, result.Experiment.Trials, result.Experiment.EvidenceClassification)
	b.WriteString("NEWLY DERIVED MATHEMATICAL RESULT\n")
	fmt.Fprintf(&b, "  theorem: %s\n  assumptions: %v\n  proof route: positive-index subadditivity, n_plus(P)<=rank(P), and exact natural-number rearrangement\n  literature status: standard finite linear algebra corollary; not claimed novel\n  rank-trace inequality used: false\n  numerical/proportion consequence: none yet\n\n", result.DerivedTheorem.CriticalRankLowerBound, result.DerivedTheorem.Assumptions)
	b.WriteString("CLAUDE INERTIA-STAGE COMPARISON\n")
	fmt.Fprintf(&b, "  route: %s\n  ~1/2 reproduced: %t\n  missing inputs: %v\n  rank-trace stage begun: %t\n\n", result.ClaudeComparison.FiniteRoute, result.ClaudeComparison.Reproduced, result.ClaudeComparison.MissingInputs, result.ClaudeComparison.RankTraceUsed)
	b.WriteString("SOUNDNESS\n")
	fmt.Fprintf(&b, "  finite M7 PSD => absence of off-critical zeros: accepted=%t\n", result.FinitePSDReverse.Accepted)
	b.WriteString("  approximate Oct eigenvalues satisfy exact theorem premises: false\n\n")
	b.WriteString("REMAINING OBLIGATION\n")
	b.WriteString("  compile a height-window zero count, far-zero threshold control, and the first/second-moment lower bound for n_plus(G); these are exactly the missing bridge to the earlier 1/2 stage.\n\n")
	b.WriteString("STATUS\n  RH unresolved. M9 proves a reusable exact finite critical-contribution bound only.\n")
	return b.String()
}

func M9JSONReport(result M9Result) ([]byte, error) {
	m8, err := M8JSONReport(result.M8)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema            string                                     `json:"schema"`
		M8                json.RawMessage                            `json:"m8"`
		Compression       semantic.FiniteCompression                 `json:"finite_compression"`
		Invariants        []semantic.SpectralInvariantClaim          `json:"spectral_invariants"`
		Contracts         []semantic.SpectralTheoremContract         `json:"spectral_theorem_contracts"`
		CriticalBudget    semantic.CriticalAggregateBudget           `json:"critical_aggregate_budget"`
		OffCriticalBudget semantic.OffCriticalAggregateBudget        `json:"off_critical_aggregate_budget"`
		Fusion            semantic.RepresentationFusion              `json:"representation_fusion"`
		DerivedTheorem    semantic.FiniteCriticalContributionTheorem `json:"newly_derived_mathematical_result"`
		Counterexamples   []M9Counterexample                         `json:"counterexamples"`
		Experiment        M9Experiment                               `json:"oct_experiment"`
		ClaudeComparison  ClaudeInertiaComparison                    `json:"claude_inertia_comparison"`
		FinitePSDReverse  ProofAttempt                               `json:"finite_psd_reverse_inference"`
	}{"riemann.semantic-graph.m9", json.RawMessage(m8), result.Compression, result.Invariants, result.Contracts, result.CriticalBudget, result.OffCriticalBudget, result.Fusion, result.DerivedTheorem, result.Counterexamples, result.Experiment, result.ClaudeComparison, result.FinitePSDReverse}
	return json.MarshalIndent(report, "", "  ")
}

func M10HumanReport(result M10Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M10 — HEIGHT-WINDOW COMPRESSION + THRESHOLDED INERTIA COUNTING\n\n")
	b.WriteString("HEIGHT WINDOW\n")
	fmt.Fprintf(&b, "  target: %s = (%s,%s]\n  center: %s\n  half-width: %s\n  convention: %s; %s\n", result.Compression.Window.ID, result.Compression.Window.Lower.Expression, result.Compression.Window.Upper.Expression, result.Compression.Window.Center.Expression, result.Compression.Window.HalfWidth.Expression, result.Compression.Window.Boundary, result.Compression.Window.OrdinateConvention)
	fmt.Fprintf(&b, "  localization: %s = (%s,%s]\n  boundary ownership is deterministic; conjugate negative ordinates are not counted in this positive window\n\n", result.Compression.LocalizationWindow.ID, result.Compression.LocalizationWindow.Lower.Expression, result.Compression.LocalizationWindow.Upper.Expression)
	b.WriteString("ZERO COUNTS\n")
	for _, item := range result.CountVocabulary {
		fmt.Fprintf(&b, "  %s\n", item)
	}
	b.WriteString("  reflection partners share an ordinate and one unordered pair ID; they remain two geometric locations\n\n")
	b.WriteString("COMPRESSION\n")
	fmt.Fprintf(&b, "  object: %s\n  basis: %s\n  height dependency: %s\n  dimension: %s\n  matrix: %s\n  normalization: %s\n  explicit-formula side: %s\n  M7 role: %s\n\n", result.Compression.ID, result.Compression.BasisFamily, result.Compression.BasisHeightDependency, result.Compression.DimensionExpression, result.Compression.MatrixID, result.Compression.Normalization, result.Compression.ExplicitFormula, result.M7RegressionRole)
	b.WriteString("ZERO-SIDE SPLIT\n")
	fmt.Fprintf(&b, "  %s\n  near membership: %s\n  far and off-critical are distinct classifications\n\n", result.Decomposition.Identity, result.Decomposition.MembershipRule)
	b.WriteString("FAR-ZERO CONTROL\n")
	fmt.Fprintf(&b, "  theorem: %s\n  norm: %s\n  bound: ||%s|| <= %s = %s\n  asymptotic: %s\n  uniformity: %s\n  evidence: trusted theorem contract, not a numerical estimate\n\n", result.FarBound.Theorems[0], result.FarBound.Norm, result.FarBound.MatrixID, result.FarBound.BoundSymbol, result.FarBound.BoundExpression, result.FarBound.AsymptoticStatement, result.FarBound.Uniformity)
	b.WriteString("THRESHOLDED SPECTRAL OBSERVATION\n")
	fmt.Fprintf(&b, "  definition: n_plus^theta(G)=#{lambda_i(G)>theta}\n  threshold: %s\n  scaling dependencies: %v\n  ordinary n_plus and thresholded n_plus^theta are distinct\n  exact sanity: n_plus^1(diag(2,1,0))=%d (equality at 1 is excluded)\n  approximate eigenvalues can certify this premise: false\n\n", result.Compression.Threshold.Expression, result.Compression.Threshold.Dependencies, result.ExactSanityObservation.Bound)
	b.WriteString("PERTURBATION / THRESHOLD THEOREM\n")
	fmt.Fprintf(&b, "  %s\n  premise: %s; comparison: %s\n  conclusion: n_plus^theta(G_tilde)<=n_plus(A_tilde)\n\n", result.Perturbation.Statement, result.Perturbation.ThresholdRule, result.Perturbation.Comparison)
	b.WriteString("M9 ACCOUNTING\n")
	fmt.Fprintf(&b, "  %s\n  reused theorem: %s\n\n", result.CountingTheorem.M9Accounting, M9CriticalRankBoundTheoremID)
	b.WriteString("NEWLY DERIVED MATHEMATICAL RESULT\n")
	fmt.Fprintf(&b, "  finite theorem: %s\n  assumptions: %v\n  enlarged window: %s\n  target simple-critical count: %s\n  target distinct-zero count: %s\n  count conversion: %s\n  literature status: Proposition 4.5 structurally reconstructed; not claimed novel\n  compiler contribution: representation-fused factorization and exact input boundaries\n\n", result.CountingTheorem.Name, result.CountingTheorem.Assumptions, result.CountingTheorem.EnlargedWindowBound, result.CountingTheorem.TargetWindowBound, result.CountingTheorem.DistinctZeroBound, result.CountingTheorem.CountConversion)
	b.WriteString("COMPLETE PROOF ROUTE\n")
	b.WriteString("  exact/certified n_plus^theta(G_tilde)>=L_theta\n  -> Proposition 4.2 gives ||E_tilde||_op<=theta0<=theta\n  -> Weyl gives n_plus^theta(G_tilde)<=n_plus(A_tilde)\n  -> M8/M9 pull-back accounting gives n_plus(A_tilde)<=s1+s2+p\n  -> N(I')>=s1+2s2+2p gives s1>=2L_theta-N(I')\n  -> remove the fringe I'\\I to obtain the target-window inequalities\n\n")
	b.WriteString("SOURCE PROOF ARCHITECTURE / COMPILER IR\n")
	fmt.Fprintf(&b, "  window localization: %s\n  far-zero control: %s\n  thresholded count: %s\n  orbit accounting: %s\n  asymptotic normalization: %s\n  fit: %s\n\n", result.Architecture.WindowLocalization, result.Architecture.FarZeroControl, result.Architecture.ThresholdedSpectralCount, result.Architecture.OrbitAccounting, result.Architecture.AsymptoticNormalization, result.Architecture.CompilerFit)
	b.WriteString("COUNTEREXAMPLE-FIRST FILTERING\n")
	for _, c := range result.Counterexamples {
		fmt.Fprintf(&b, "  rejected: %s\n    fixture: %s\n    reason: %s\n", c.RejectedCandidate, c.ExactFixture, c.Reason)
	}
	b.WriteString("\nOCT EXPERIMENT\n")
	fmt.Fprintf(&b, "  path: %s\n  command: %s\n  setup: %s\n  threshold: %s\n  perturbation bound: %s\n  trials: %d\n  execution: %s\n  counterexamples: %v\n  evidence: %s\n  when utility used: %t (operator-norm Weyl is dictated by Proposition 4.5)\n\n", result.Experiment.Path, result.Experiment.Command, result.Experiment.Setup, result.Experiment.Threshold, result.Experiment.PerturbationBound, result.Experiment.Trials, result.Experiment.Execution, result.Experiment.CounterexamplesFound, result.Experiment.EvidenceClassification, result.UtilitySchedulerUsed)
	b.WriteString("ASYMPTOTIC CONSEQUENCE\n")
	b.WriteString("  half-type proportion reached: false\n  no 2/3 or 67.25% result attempted\n\n")
	b.WriteString("REMAINING INPUT\n")
	fmt.Fprintf(&b, "  %s\n\n", result.CountingTheorem.RemainingInput)
	b.WriteString("RH\n  unresolved\n")
	return b.String()
}

func M10JSONReport(result M10Result) ([]byte, error) {
	m9, err := M9JSONReport(result.M9)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema                 string                                 `json:"schema"`
		M9                     json.RawMessage                        `json:"m9"`
		Compression            semantic.WindowCompression             `json:"window_compression"`
		CountVocabulary        []string                               `json:"zero_count_vocabulary"`
		Decomposition          semantic.NearFarZeroDecomposition      `json:"near_far_decomposition"`
		FarBound               semantic.FarZeroContributionBound      `json:"far_zero_bound"`
		Perturbation           semantic.ThresholdPerturbationContract `json:"threshold_perturbation_theorem"`
		CountingTheorem        semantic.FiniteWindowCountingTheorem   `json:"newly_derived_mathematical_result"`
		ExactSanityObservation semantic.ThresholdedPositiveIndexClaim `json:"exact_threshold_sanity_observation"`
		Architecture           ProofArchitectureMap                   `json:"proof_architecture"`
		Counterexamples        []M10Counterexample                    `json:"counterexamples"`
		Experiment             M10Experiment                          `json:"oct_experiment"`
		UtilitySchedulerUsed   bool                                   `json:"when_utility_used"`
		M7RegressionRole       string                                 `json:"m7_regression_role"`
		Sources                []semantic.Reference                   `json:"sources"`
	}{"riemann.semantic-graph.m10", json.RawMessage(m9), result.Compression, result.CountVocabulary, result.Decomposition, result.FarBound, result.Perturbation, result.CountingTheorem, result.ExactSanityObservation, result.Architecture, result.Counterexamples, result.Experiment, result.UtilitySchedulerUsed, result.M7RegressionRole, result.Sources}
	return json.MarshalIndent(report, "", "  ")
}

func M11HumanReport(result M11Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M11 — FIRST/SECOND MOMENTS TO THRESHOLDED POSITIVE INDEX\n\n")
	b.WriteString("MOMENT INPUT\n")
	for _, m := range result.Moments {
		fmt.Fprintf(&b, "  %s\n    matrix: %s\n    statement: %s\n    evidence: %s\n", m.Kind, m.MatrixID, m.Expression, m.Evidence)
	}
	b.WriteString("  normalization: G_tilde_T=G_T/L, L=lambda*log(T/2pi)\n\n")
	b.WriteString("FINITE SPECTRAL THEOREM\n")
	fmt.Fprintf(&b, "  %s\n  assumptions: %v\n  partition: %s\n  trace residual: %s\n  Cauchy-Schwarz: %s\n  real bound: %s\n  integer bound: %s\n\n", result.FiniteTheorem.Name, result.FiniteTheorem.Assumptions, result.FiniteTheorem.Partition, result.FiniteTheorem.TraceResidual, result.FiniteTheorem.CauchySchwarz, result.FiniteTheorem.RealConclusion, result.FiniteTheorem.IntegerConclusion)
	b.WriteString("ASYMPTOTIC DISCHARGE\n")
	for _, m := range result.AsymptoticMoments {
		fmt.Fprintf(&b, "  %s\n    main: %s\n    remainder: %s\n    scale: %s\n", m.Moment.Kind, m.MainTerm, m.Remainder, m.Scale)
	}
	fmt.Fprintf(&b, "  moment error: %s\n  dimension: %s\n  threshold: %s\n  penalty: %s\n  relative penalty: %s\n  conclusion: %s\n\n", result.ThresholdScaling.MomentError, result.ThresholdScaling.DimensionAsymptotic, result.ThresholdScaling.Threshold, result.ThresholdScaling.Penalty, result.ThresholdScaling.RelativePenalty, result.ThresholdScaling.Conclusion)
	b.WriteString("THRESHOLDED POSITIVE INDEX\n")
	fmt.Fprintf(&b, "  eventually finite: %s\n  asymptotic: %s\n  F(lambda): %s\n  endpoint: %s\n\n", result.AsymptoticCount.FiniteEpsilonBound, result.AsymptoticCount.NormalizedLowerBound, result.AsymptoticCount.BandwidthFunction, result.AsymptoticCount.EndpointIndexBound)
	b.WriteString("M10 COMPOSITION\n")
	fmt.Fprintf(&b, "  simple critical zeros: %s\n  distinct zeros: %s\n  fringe: %s\n\n", result.AsymptoticCount.M10SimpleComposition, result.AsymptoticCount.M10DistinctComposition, result.AsymptoticCount.Fringe)
	b.WriteString("ASYMPTOTIC CONSEQUENCE\n")
	fmt.Fprintf(&b, "  %s\n  exact constant: %s\n  %s\n  half-type reproduced: %t\n  no stronger rank/trace optimization attempted\n\n", result.AsymptoticCount.SimpleCriticalLiminf, result.AsymptoticCount.ExactSimpleConstant, result.AsymptoticCount.DistinctLiminf, result.AsymptoticCount.HalfTypeReproduced)
	b.WriteString("COUNTEREXAMPLES TO REJECTED VARIANTS\n")
	for _, c := range result.Counterexamples {
		fmt.Fprintf(&b, "  rejected: %s\n    spectrum: %s\n    reason: %s\n", c.RejectedCandidate, c.ExactSpectrum, c.Reason)
	}
	b.WriteString("\nOCT EXPERIMENT\n")
	fmt.Fprintf(&b, "  path: %s\n  check: %s\n  run: %s\n  setup: %s\n  trials: %d\n  execution: %s\n  compiler: %s\n  limits: %s\n  findings: %v\n  evidence: %s\n  when utility used: %t (the paper's Lemma 3.3 fixes the simple branch; no live route choice remained)\n\n", result.Experiment.Path, result.Experiment.CheckCommand, result.Experiment.RunCommand, result.Experiment.Setup, result.Experiment.Trials, result.Experiment.Execution, result.Experiment.CompilerIdentity, result.Experiment.TimingAndLimits, result.Experiment.Findings, result.Experiment.EvidenceClassification, result.UtilitySchedulerUsed)
	b.WriteString("IMPORTED FROM LITERATURE\n")
	for _, x := range result.ImportedMathematics {
		fmt.Fprintf(&b, "  %s\n", x)
	}
	b.WriteString("\nNEWLY DERIVED BY COMPILER / RESEARCH LOOP\n")
	for _, x := range result.DerivedMathematics {
		fmt.Fprintf(&b, "  %s\n", x)
	}
	b.WriteString("\nREPRESENTATION FUSION\n")
	for _, x := range result.Fusion {
		fmt.Fprintf(&b, "  %s\n", x)
	}
	b.WriteString("\nARCHITECTURAL AWKWARDNESS\n  asymptotic moment claims and exact finite bounds require a deliberate Eventually adapter; encoding them as one scalar claim would launder o-terms. The current string-valued asymptotic algebra is minimal but should not grow into an untyped symbolic engine.\n")
	b.WriteString("\nCOMPILER THEORY\n  theorem compilation must preserve statistic kind, evidence grade, scale, and integer codomain across representations. The decisive operation is staged representation fusion, not formula substitution.\n")
	b.WriteString("\nONE NEXT MILESTONE\n  M12: encode and verify the paper's finite rank-trace inequality as a generic Hermitian theorem, then compare its certified index/count consequence with M11's sharp first-two-moment ceiling.\n")
	b.WriteString("\nSTATUS\n  known half-stage result reproduced: yes\n  simple critical liminf: 1/2\n  distinct-zero liminf from this route: 3/4\n  RH\n  unresolved\n")
	return b.String()
}

func M11JSONReport(result M11Result) ([]byte, error) {
	if err := validateM11Result(result); err != nil {
		return nil, err
	}
	m10, err := M10JSONReport(result.M10)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema               string                               `json:"schema"`
		M10                  json.RawMessage                      `json:"m10"`
		Moments              []semantic.SpectralMomentClaim       `json:"moment_inputs"`
		AsymptoticMoments    []semantic.AsymptoticMomentStatement `json:"asymptotic_moments"`
		EventuallyBounds     []semantic.EventuallyBound           `json:"eventually_finite_bounds"`
		FiniteTheorem        semantic.FiniteMomentCountTheorem    `json:"finite_spectral_theorem"`
		ExactSanity          semantic.MomentCountResult           `json:"exact_sanity"`
		ThresholdScaling     M11ThresholdScaling                  `json:"threshold_scaling"`
		AsymptoticCount      M11AsymptoticCount                   `json:"asymptotic_count"`
		M10ReuseSanity       WindowCountBounds                    `json:"m10_reuse_sanity"`
		Counterexamples      []M11Counterexample                  `json:"counterexamples"`
		Experiment           M11Experiment                        `json:"oct_experiment"`
		UtilitySchedulerUsed bool                                 `json:"when_utility_used"`
		Fusion               []string                             `json:"representation_fusion"`
		Imported             []string                             `json:"imported_from_literature"`
		Derived              []string                             `json:"newly_derived_by_compiler_research_loop"`
		Sources              []semantic.Reference                 `json:"sources"`
	}{"riemann.semantic-graph.m11", json.RawMessage(m10), result.Moments, result.AsymptoticMoments, result.EventuallyBounds, result.FiniteTheorem, result.ExactSanity, result.ThresholdScaling, result.AsymptoticCount, result.M10ReuseSanity, result.Counterexamples, result.Experiment, result.UtilitySchedulerUsed, result.Fusion, result.ImportedMathematics, result.DerivedMathematics, result.Sources}
	return json.MarshalIndent(report, "", "  ")
}

func M12HumanReport(result M12Result) string {
	var b strings.Builder
	b.WriteString("RIEMANN-M12 — RANK-TRACE INEQUALITY BEYOND THE FIRST-TWO-MOMENT CEILING\n\n")
	b.WriteString("RANK-TRACE INPUT\n")
	fmt.Fprintf(&b, "  %s\n  P: Hermitian, PSD, rank(P)<=r\n  Q: Hermitian, n_plus(Q)<=b; trace(Q) has no sign premise\n  total: ||G||_F^2=tr(G^2)\n  parameter: %s, domain %s\n\n", result.Decomposition.Identity, result.Parameter.Symbol, result.Parameter.Domain)
	b.WriteString("GENERIC FINITE THEOREM\n")
	fmt.Fprintf(&b, "  assumptions: %v\n  expansion: %s\n  von Neumann step: %s\n  scalar steps: %v\n  conclusion: %s\n  equality: %s\n\n", result.FiniteTheorem.Assumptions, result.FiniteTheorem.Expansion, result.FiniteTheorem.VonNeumannStep, result.FiniteTheorem.ScalarSteps, result.FiniteTheorem.Conclusion, result.FiniteTheorem.EqualityCase)
	b.WriteString("RANK CONSEQUENCE\n")
	fmt.Fprintf(&b, "  %s\n  %s\n\n", result.FiniteTheorem.Specialization, result.FiniteWindow.AllCriticalRankBound)
	b.WriteString("ZERO COUNT CONSEQUENCE\n")
	fmt.Fprintf(&b, "  regrouping: %s\n  simple critical: %s\n  critical distinct: %s\n  all distinct: %s\n  tail/fringe: %s\n\n", result.FiniteWindow.SimpleRegrouping, result.FiniteWindow.SimpleCountBound, result.FiniteWindow.CriticalCountBound, result.FiniteWindow.DistinctCountBound, result.FiniteWindow.TailTransfer)
	b.WriteString("ASYMPTOTIC RESULT\n")
	fmt.Fprintf(&b, "  trace: %s\n  Frobenius: %s\n  algebra: %s\n  %s\n  %s\n  %s\n\n", result.AsymptoticCount.NormalizedTrace, result.AsymptoticCount.NormalizedFrobenius, result.AsymptoticCount.BandwidthFunction, result.AsymptoticCount.SimpleCriticalLiminf, result.AsymptoticCount.CriticalLiminf, result.AsymptoticCount.DistinctLiminf)
	b.WriteString("COMPARISON\n  M11 sees only total first/second moments: 1/2 simple-critical stage (sharp for that IR).\n  M12 also consumes the exact component identity, P PSD/rank semantics, Q positive-index budget, and component traces: 2/3 simple-critical stage.\n  Same G=I_2 admits P=I_2,Q=0 and P=0,Q=I_2, so total moments alone cannot recover rank(P).\n\n")
	b.WriteString("COUNTEREXAMPLES AND EDGE FIXTURES\n")
	for _, c := range result.Counterexamples {
		fmt.Fprintf(&b, "  rejected: %s\n    fixture: %s\n    result: %s\n", c.RejectedCandidate, c.ExactFixture, c.Failure)
	}
	b.WriteString("\nOCT / OCTGO\n")
	fmt.Fprintf(&b, "  path: %s\n  command: %s\n  trials: %d\n  execution: %s\n  evidence: %s\n  OctGo used: %t\n  when utility used: %t (the source and algebra selected the von Neumann route directly)\n\n", result.Experiment.Path, result.Experiment.Command, result.Experiment.Trials, result.Experiment.Execution, result.Experiment.EvidenceClassification, result.OctGoUsed, result.UtilitySchedulerUsed)
	b.WriteString("IMPORTED MATHEMATICS\n")
	for _, x := range result.ImportedMathematics {
		fmt.Fprintf(&b, "  %s\n", x)
	}
	b.WriteString("\nNEWLY DERIVED / RECONSTRUCTED MATHEMATICAL RESULT\n\n")
	b.WriteString("GENERIC FINITE THEOREM\n  The exact parameterized Lemma 3.2 consequence above is reconstructed from finite algebra. It is standard/known here, not claimed novel. Imported pieces are the Hermitian spectral decomposition and von Neumann trace inequality.\n\n")
	b.WriteString("ZETA CONSEQUENCE\n  The paper's bandwidth-one normalization and existing M8-M11 inputs yield 2/3 for simple and distinct critical-line zeros and 5/6 for all distinct zeros. This matches Theorems A-C; no 67.25% optimization was begun.\n\n")
	b.WriteString("REPRESENTATION FUSION\n")
	for _, x := range result.Fusion {
		fmt.Fprintf(&b, "  %s\n", x)
	}
	b.WriteString("\nARCHITECTURAL AWKWARDNESS\n  M11 stores Theorem 5.8 in G_tilde units while rank-trace is scale-sensitive and must use G_hat=G_tilde/(aL). M12 records that normalization adapter explicitly; a future typed scale-change IR would be safer than extending string algebra.\n")
	b.WriteString("\nCOMPILER THEORY\n  Yes: preserving G=P+Q lets the compiler prove strictly more than a sharp total-moment theorem. The gain comes from component identity, PSD/rank information, positive rather than total index of Q, and trace linearity.\n")
	b.WriteString("\nONE NEXT MILESTONE\n  M13: type and verify the one-variable test-window functional, then reproduce the Montgomery-Taylor 67.25% optimization without changing the M12 finite theorem.\n")
	b.WriteString("\nSTATUS\n  M12 two-thirds stage reproduced: yes\n  simple critical liminf: 2/3\n  critical distinct liminf: 2/3\n  all distinct liminf: 5/6\n  RH\n  unresolved\n")
	return b.String()
}

func M12JSONReport(result M12Result) ([]byte, error) {
	if err := validateM12Result(result); err != nil {
		return nil, err
	}
	m11, err := M11JSONReport(result.M11)
	if err != nil {
		return nil, err
	}
	report := struct {
		Schema string          `json:"schema"`
		M11    json.RawMessage `json:"m11"`
		Result M12Result       `json:"m12"`
	}{"riemann.semantic-graph.m12", json.RawMessage(m11), result}
	return json.MarshalIndent(report, "", "  ")
}
