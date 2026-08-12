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
	Kind                   semantic.PropositionKind         `json:"kind"`
	Description            string                           `json:"description"`
	Quantifier             semantic.QuantifierKind          `json:"quantifier,omitempty"`
	Domain                 *semantic.Domain                 `json:"domain,omitempty"`
	Predicate              *semantic.Predicate              `json:"predicate,omitempty"`
	Representation         *semantic.Representation         `json:"representation,omitempty"`
	RepresentationIdentity *semantic.RepresentationIdentity `json:"representation_identity,omitempty"`
	AnalyticFact           *semantic.AnalyticFact           `json:"analytic_fact,omitempty"`
	ZeroAtPoint            *semantic.ZeroAtPoint            `json:"zero_at_point,omitempty"`
	SideCondition          *semantic.SideCondition          `json:"side_condition,omitempty"`
	FunctionalIdentity     *semantic.FunctionalIdentity     `json:"functional_identity,omitempty"`
	ZeroSetProperty        *semantic.ZeroSetProperty        `json:"zero_set_property,omitempty"`
	ZeroClassification     *semantic.ZeroClassification     `json:"zero_classification,omitempty"`
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
