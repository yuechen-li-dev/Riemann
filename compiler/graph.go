// Package compiler owns proof state and enforces transformation direction,
// evidence boundaries, capability requirements, and proof obligations.
package compiler

import (
	"errors"
	"fmt"
	"sort"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type Relation string

const (
	Equivalent    Relation = "equivalent"
	Implies       Relation = "implies"
	Relaxation    Relation = "relaxation"
	Approximation Relation = "approximation"
)

func (r Relation) Valid() bool {
	switch r {
	case Equivalent, Implies, Relaxation, Approximation:
		return true
	default:
		return false
	}
}

func (r Relation) Reversible() bool { return r == Equivalent }

type InformationLoss struct {
	Property semantic.Property
	Reason   string
}

type Transformation struct {
	ID          semantic.TransformationID
	Pass        string
	From        semantic.ClaimID
	To          semantic.ClaimID
	Relation    Relation
	Obligations []semantic.ClaimID
	Losses      []InformationLoss
	Provenance  semantic.Reference
}

type Graph struct {
	claims          map[semantic.ClaimID]semantic.Claim
	claimOrder      []semantic.ClaimID
	transformations []Transformation
}

func NewGraph() *Graph { return &Graph{claims: make(map[semantic.ClaimID]semantic.Claim)} }

func (g *Graph) AddClaim(claim semantic.Claim) error {
	if err := claim.Validate(); err != nil {
		return err
	}
	if _, exists := g.claims[claim.ID]; exists {
		return fmt.Errorf("claim %q already exists", claim.ID)
	}
	claim.Assumptions = semantic.CloneAssumptions(claim.Assumptions)
	claim.Evidence = append([]semantic.Evidence(nil), claim.Evidence...)
	claim.Provenance.Parents = append([]semantic.ClaimID(nil), claim.Provenance.Parents...)
	g.claims[claim.ID] = claim
	g.claimOrder = append(g.claimOrder, claim.ID)
	return nil
}

func (g *Graph) AddTransformation(t Transformation) error {
	if t.ID == "" || t.Pass == "" || !t.Relation.Valid() {
		return fmt.Errorf("invalid transformation metadata")
	}
	from, fromOK := g.claims[t.From]
	to, toOK := g.claims[t.To]
	if !fromOK || !toOK {
		return fmt.Errorf("transformation %q refers to an unknown endpoint", t.ID)
	}
	for _, existing := range g.transformations {
		if existing.ID == t.ID {
			return fmt.Errorf("transformation %q already exists", t.ID)
		}
	}
	for _, obligation := range t.Obligations {
		if _, ok := g.claims[obligation]; !ok {
			return fmt.Errorf("transformation %q refers to unknown obligation %q", t.ID, obligation)
		}
	}
	if !assumptionsContain(to.Assumptions, from.Assumptions) {
		return fmt.Errorf("transformation %q silently drops assumptions", t.ID)
	}
	if t.Relation == Equivalent && !assumptionsContain(from.Assumptions, to.Assumptions) {
		return fmt.Errorf("equivalence %q changes its assumption set", t.ID)
	}
	if t.Relation == Equivalent && (from.Capabilities != to.Capabilities || from.Exactness != to.Exactness) {
		return fmt.Errorf("equivalence %q changes capabilities or exactness", t.ID)
	}
	declaredLoss := semantic.PropertySet(0)
	for _, loss := range t.Losses {
		declaredLoss |= semantic.Properties(loss.Property)
	}
	actualLoss := from.Capabilities.Without(to.Capabilities)
	if !declaredLoss.Contains(actualLoss) {
		return fmt.Errorf("transformation %q fails to declare capability loss: %v", t.ID, actualLoss.Names())
	}
	gained := to.Capabilities.Without(from.Capabilities)
	if gained != 0 {
		return fmt.Errorf("transformation %q gains unsupported capabilities: %v", t.ID, gained.Names())
	}
	if to.Provenance.Kind != semantic.DerivedProvenance || to.Provenance.Transformation != t.ID || !containsID(to.Provenance.Parents, from.ID) {
		return fmt.Errorf("transformation %q target does not retain matching derived provenance", t.ID)
	}
	t.Obligations = append([]semantic.ClaimID(nil), t.Obligations...)
	t.Losses = append([]InformationLoss(nil), t.Losses...)
	g.transformations = append(g.transformations, t)
	return nil
}

func containsID(ids []semantic.ClaimID, want semantic.ClaimID) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}

func assumptionsContain(have, need []semantic.Assumption) bool {
	set := make(map[semantic.AssumptionID]bool, len(have))
	for _, assumption := range have {
		set[assumption.ID] = true
	}
	for _, assumption := range need {
		if !set[assumption.ID] {
			return false
		}
	}
	return true
}

func (g *Graph) Claim(id semantic.ClaimID) (semantic.Claim, bool) {
	claim, ok := g.claims[id]
	return cloneClaim(claim), ok
}

func (g *Graph) Claims() []semantic.Claim {
	claims := make([]semantic.Claim, 0, len(g.claimOrder))
	for _, id := range g.claimOrder {
		claims = append(claims, cloneClaim(g.claims[id]))
	}
	return claims
}

func (g *Graph) Transformations() []Transformation {
	out := make([]Transformation, len(g.transformations))
	for i, transformation := range g.transformations {
		out[i] = cloneTransformation(transformation)
	}
	return out
}

func cloneClaim(claim semantic.Claim) semantic.Claim {
	claim.Assumptions = semantic.CloneAssumptions(claim.Assumptions)
	claim.Evidence = append([]semantic.Evidence(nil), claim.Evidence...)
	claim.Provenance.Parents = append([]semantic.ClaimID(nil), claim.Provenance.Parents...)
	return claim
}

func cloneTransformation(transformation Transformation) Transformation {
	transformation.Obligations = append([]semantic.ClaimID(nil), transformation.Obligations...)
	transformation.Losses = append([]InformationLoss(nil), transformation.Losses...)
	return transformation
}

type DiagnosticCode string

const (
	NoEstablishedDirection DiagnosticCode = "no_established_direction"
	MissingCapability      DiagnosticCode = "missing_capability"
	ApproximationBoundary  DiagnosticCode = "approximation_cannot_certify_exact_claim"
	OpenObligation         DiagnosticCode = "open_obligation"
	UncertifiedEvidence    DiagnosticCode = "uncertified_evidence"
)

type Diagnostic struct {
	Code    DiagnosticCode `json:"code"`
	Message string         `json:"message"`
}

type ProofAttempt struct {
	From        semantic.ClaimID `json:"from"`
	Target      semantic.ClaimID `json:"target"`
	Accepted    bool             `json:"accepted"`
	Diagnostics []Diagnostic     `json:"diagnostics"`
}

// AttemptProof combines structural discharge checking with evidence
// certification. Structural errors take precedence so a failed path reports
// the mathematical mismatch rather than an incidental source-evidence error.
func (g *Graph) AttemptProof(fromID, targetID semantic.ClaimID) ProofAttempt {
	attempt := g.CheckDischarge(fromID, targetID)
	if len(attempt.Diagnostics) != 0 {
		return attempt
	}
	certified, diagnostics := g.Certify(fromID)
	if !certified {
		attempt.Diagnostics = append(attempt.Diagnostics, diagnostics...)
		attempt.Accepted = false
	}
	return attempt
}

// CheckDischarge checks whether an available claim is structurally usable for
// a target. Truth/certification of the available claim is checked separately by
// Certify, which is important for inspecting hypothetical compiler lowerings.
func (g *Graph) CheckDischarge(fromID, targetID semantic.ClaimID) ProofAttempt {
	attempt := ProofAttempt{From: fromID, Target: targetID}
	from, fromOK := g.claims[fromID]
	target, targetOK := g.claims[targetID]
	if !fromOK || !targetOK {
		attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{NoEstablishedDirection, "source or target claim is unknown"})
		return attempt
	}

	path, found := g.findPath(fromID, targetID)
	if !found {
		message := fmt.Sprintf("no established transformation permits %s → %s", fromID, targetID)
		for _, t := range g.transformations {
			if t.From == targetID && t.To == fromID && !t.Relation.Reversible() {
				message = fmt.Sprintf("%s is directional (%s → %s) and cannot be traversed backward", t.Relation, t.From, t.To)
				break
			}
		}
		if reversePath, reverseFound := g.findPath(targetID, fromID); reverseFound {
			for _, t := range reversePath {
				if !t.Relation.Reversible() {
					message = fmt.Sprintf("the established path contains %s (%s → %s), which cannot be traversed backward", t.Relation, t.From, t.To)
					break
				}
			}
		}
		attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{NoEstablishedDirection, message})
	}

	availableCapabilities := from.Capabilities
	if found {
		for _, step := range path {
			for _, loss := range step.Losses {
				availableCapabilities = availableCapabilities.Without(semantic.Properties(loss.Property))
			}
		}
	}
	missing := target.Requirements.Without(availableCapabilities)
	for _, property := range missing.Names() {
		attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{
			MissingCapability,
			fmt.Sprintf("target requires %s, which the source representation does not retain", property),
		})
	}

	if found {
		for _, step := range path {
			if step.Relation == Approximation && target.Exactness == semantic.Exact {
				attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{ApproximationBoundary, "an approximation cannot certify an exact theorem target"})
			}
			for _, obligation := range step.Obligations {
				if certified, _ := g.certify(obligation, make(map[semantic.ClaimID]bool)); !certified {
					attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{OpenObligation, fmt.Sprintf("proof obligation %s remains open", obligation)})
				}
			}
		}
	}
	attempt.Accepted = found && len(attempt.Diagnostics) == 0
	return attempt
}

func (g *Graph) findPath(from, target semantic.ClaimID) ([]Transformation, bool) {
	if from == target {
		return nil, true
	}
	type state struct {
		id   semantic.ClaimID
		path []Transformation
	}
	queue := []state{{id: from}}
	seen := map[semantic.ClaimID]bool{from: true}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		for _, t := range g.transformations {
			var next semantic.ClaimID
			var ok bool
			switch {
			case t.From == current.id:
				next, ok = t.To, true
			case t.Relation.Reversible() && t.To == current.id:
				next, ok = t.From, true
			}
			if !ok || seen[next] {
				continue
			}
			path := append(append([]Transformation(nil), current.path...), t)
			if next == target {
				return path, true
			}
			seen[next] = true
			queue = append(queue, state{id: next, path: path})
		}
	}
	return nil, false
}

// Certify determines whether a claim has exact certifying evidence. Numerical
// evidence and conjectures are retained but never treated as certificates.
func (g *Graph) Certify(id semantic.ClaimID) (bool, []Diagnostic) {
	return g.certify(id, make(map[semantic.ClaimID]bool))
}

func (g *Graph) certify(id semantic.ClaimID, visiting map[semantic.ClaimID]bool) (bool, []Diagnostic) {
	claim, ok := g.claims[id]
	if !ok {
		return false, []Diagnostic{{UncertifiedEvidence, fmt.Sprintf("claim %s is unknown", id)}}
	}
	if visiting[id] {
		return false, []Diagnostic{{UncertifiedEvidence, fmt.Sprintf("cyclic derivation at %s", id)}}
	}
	if missing := claim.Requirements.Without(claim.Capabilities); missing != 0 {
		return false, []Diagnostic{{MissingCapability, fmt.Sprintf("claim %s lacks required capabilities: %v", id, missing.Names())}}
	}
	for _, evidence := range claim.Evidence {
		if evidence.Kind == semantic.DefinitionEvidence || evidence.Kind == semantic.KnownTheoremEvidence {
			return true, nil
		}
	}
	visiting[id] = true
	defer delete(visiting, id)
	var diagnostics []Diagnostic
	for _, t := range g.transformations {
		parentID := semantic.ClaimID("")
		switch {
		case t.To == id:
			parentID = t.From
		case t.Relation.Reversible() && t.From == id:
			parentID = t.To
		default:
			continue
		}
		if t.Relation == Approximation && claim.Exactness == semantic.Exact {
			diagnostics = append(diagnostics, Diagnostic{ApproximationBoundary, "an approximation cannot certify an exact theorem target"})
			continue
		}
		blocked := false
		for _, obligation := range t.Obligations {
			if certified, _ := g.certify(obligation, visiting); !certified {
				diagnostics = append(diagnostics, Diagnostic{OpenObligation, fmt.Sprintf("proof obligation %s remains open", obligation)})
				blocked = true
			}
		}
		parentCertified, _ := g.certify(parentID, visiting)
		if !parentCertified {
			blocked = true
		}
		if !blocked {
			return true, nil
		}
	}
	if len(diagnostics) == 0 {
		diagnostics = append(diagnostics, Diagnostic{UncertifiedEvidence, fmt.Sprintf("claim %s has no exact certifying evidence", id)})
	}
	return false, deduplicateDiagnostics(diagnostics)
}

// Lineage returns a deterministic depth-first provenance trace ending at id.
func (g *Graph) Lineage(id semantic.ClaimID) ([]semantic.ClaimID, error) {
	if _, ok := g.claims[id]; !ok {
		return nil, fmt.Errorf("claim %q is unknown", id)
	}
	seen := make(map[semantic.ClaimID]bool)
	var out []semantic.ClaimID
	var visit func(semantic.ClaimID)
	visit = func(current semantic.ClaimID) {
		if seen[current] {
			return
		}
		seen[current] = true
		claim := g.claims[current]
		for _, parent := range sortedIDs(claim.Provenance.Parents) {
			visit(parent)
		}
		out = append(out, current)
	}
	visit(id)
	return out, nil
}

func deduplicateDiagnostics(in []Diagnostic) []Diagnostic {
	seen := make(map[string]bool)
	out := make([]Diagnostic, 0, len(in))
	for _, d := range in {
		key := string(d.Code) + "\x00" + d.Message
		if !seen[key] {
			seen[key] = true
			out = append(out, d)
		}
	}
	return out
}

var ErrWrongProposition = errors.New("pass received the wrong proposition kind")

func sortedIDs(ids []semantic.ClaimID) []semantic.ClaimID {
	out := append([]semantic.ClaimID(nil), ids...)
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
