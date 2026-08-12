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
		} else {
			for _, loss := range t.Losses {
				fmt.Fprintf(&b, "  information loss: %s\n", loss.Property.String())
			}
		}
		fmt.Fprintf(&b, "RESULT\n  %s\n", to.Proposition.Describe())
	}
	density, _ := result.Graph.Claim(result.Attempt.From)
	b.WriteString("\nATTEMPT\n")
	fmt.Fprintf(&b, "  use %q to discharge %q\n", density.Proposition.Describe(), root.Proposition.Describe())
	if result.Attempt.Accepted {
		b.WriteString("ACCEPTED\n")
	} else {
		b.WriteString("REJECTED\n")
		for _, diagnostic := range result.Attempt.Diagnostics {
			fmt.Fprintf(&b, "  [%s] %s\n", diagnostic.Code, diagnostic.Message)
		}
	}
	return b.String()
}

type graphJSON struct {
	Schema          string               `json:"schema"`
	Claims          []claimJSON          `json:"claims"`
	Transformations []transformationJSON `json:"transformations"`
	Attempt         ProofAttempt         `json:"attempt"`
}

type claimJSON struct {
	ID           semantic.ClaimID      `json:"id"`
	Proposition  propositionJSON       `json:"proposition"`
	Assumptions  []semantic.Assumption `json:"assumptions"`
	Evidence     []semantic.Evidence   `json:"evidence"`
	Requirements []string              `json:"requirements"`
	Capabilities []string              `json:"capabilities"`
	Exactness    semantic.Exactness    `json:"exactness"`
	Provenance   semantic.Provenance   `json:"provenance"`
}

type propositionJSON struct {
	Kind        semantic.PropositionKind `json:"kind"`
	Description string                   `json:"description"`
}

type transformationJSON struct {
	ID          semantic.TransformationID `json:"id"`
	Pass        string                    `json:"pass"`
	From        semantic.ClaimID          `json:"from"`
	To          semantic.ClaimID          `json:"to"`
	Relation    Relation                  `json:"relation"`
	Obligations []semantic.ClaimID        `json:"obligations"`
	Losses      []lossJSON                `json:"losses"`
	Provenance  semantic.Reference        `json:"provenance"`
}

type lossJSON struct {
	Property string `json:"property"`
	Reason   string `json:"reason"`
}

func JSONReport(result M0Result) ([]byte, error) {
	report := graphJSON{Schema: "riemann.semantic-graph.m0", Attempt: result.Attempt}
	for _, claim := range result.Graph.Claims() {
		report.Claims = append(report.Claims, claimJSON{
			ID:          claim.ID,
			Proposition: propositionJSON{Kind: claim.Proposition.Kind(), Description: claim.Proposition.Describe()},
			Assumptions: nonNil(claim.Assumptions), Evidence: nonNil(claim.Evidence),
			Requirements: nonNil(claim.Requirements.Names()), Capabilities: nonNil(claim.Capabilities.Names()),
			Exactness: claim.Exactness, Provenance: claim.Provenance,
		})
	}
	for _, t := range result.Graph.Transformations() {
		item := transformationJSON{ID: t.ID, Pass: t.Pass, From: t.From, To: t.To, Relation: t.Relation,
			Obligations: nonNil(t.Obligations), Provenance: t.Provenance}
		for _, loss := range t.Losses {
			item.Losses = append(item.Losses, lossJSON{Property: loss.Property.String(), Reason: loss.Reason})
		}
		item.Losses = nonNil(item.Losses)
		report.Transformations = append(report.Transformations, item)
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
