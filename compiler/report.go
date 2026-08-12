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
	zeroFree, _ := result.Graph.Claim(ZeroFreeHalfPlaneID)
	euler, _ := result.Graph.Claim(EulerRepresentationID)
	fmt.Fprintf(&b, "TARGET\n  %s\n\n", zeroFree.Proposition.Describe())
	b.WriteString("REPRESENTATION\n  Riemann zeta function\n    → Euler product\n")
	if p, ok := euler.Proposition.(semantic.RepresentationProposition); ok {
		fmt.Fprintf(&b, "  formula: %s\n  domain: %s\n  relation: equivalent on stated domain\n", p.Representation.Formula, p.Representation.ValidOn.Describe())
		for _, affordance := range p.Representation.Affordances {
			fmt.Fprintf(&b, "  exposes: %s\n", affordance)
		}
	}
	b.WriteString("\nPREMISES\n")
	for _, id := range []semantic.ClaimID{EulerIdentityID, EulerConvergenceID, EulerFactorsNonzeroID, InfiniteProductTheoremID} {
		claim, _ := result.Graph.Claim(id)
		fmt.Fprintf(&b, "  %s: %s\n", id, claim.Proposition.Describe())
	}
	fmt.Fprintf(&b, "\nDERIVE\n  %s\n", zeroFree.Proposition.Describe())
	if result.ZeroFreeCertified {
		b.WriteString("\nCERTIFIED\n  relative to the listed trusted analytic premises\n")
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

func renderAttempt(b *strings.Builder, g *Graph, attempt ProofAttempt) {
	from, _ := g.Claim(attempt.From)
	target, _ := g.Claim(attempt.Target)
	b.WriteString("\nATTEMPT\n")
	fmt.Fprintf(b, "  use %q to discharge %q\n", from.Proposition.Describe(), target.Proposition.Describe())
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
}

func JSONReport(result M0Result) ([]byte, error) {
	return marshalGraph("riemann.semantic-graph.m0", result.Graph, nil, []ProofAttempt{result.Attempt})
}
func M1JSONReport(result M1Result) ([]byte, error) {
	certifications := []certificationJSON{{Claim: ZeroFreeHalfPlaneID, Certified: result.ZeroFreeCertified, Diagnostics: nonNil(result.ZeroFreeDiagnostics)}}
	return marshalGraph("riemann.semantic-graph.m1", result.Graph, certifications, []ProofAttempt{result.BoundedToRH, result.DensityToRH, result.ZeroFreeToRH})
}

func marshalGraph(schema string, g *Graph, certifications []certificationJSON, attempts []ProofAttempt) ([]byte, error) {
	report := graphJSON{Schema: schema, Certifications: nonNil(certifications), Attempts: nonNil(attempts)}
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
		report.Transformations = append(report.Transformations, transformationJSON{ID: t.ID, Pass: t.Pass, From: t.From, Premises: nonNil(t.Premises), To: t.To, Relation: t.Relation, Obligations: nonNil(t.Obligations), Losses: nonNil(t.Losses), Provenance: t.Provenance})
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
