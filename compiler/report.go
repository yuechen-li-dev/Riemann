package compiler

import (
	"bytes"
	"encoding/json"
	"fmt"
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
