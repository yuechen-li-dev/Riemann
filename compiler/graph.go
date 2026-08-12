// Package compiler owns proof state and enforces direction, evidence boundaries,
// structural strength, assumptions, premises, obligations, and provenance.
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
	return r == Equivalent || r == Implies || r == Relaxation || r == Approximation
}
func (r Relation) Reversible() bool { return r == Equivalent }

type LossKind string

const (
	QuantifierWeakening      LossKind = "quantifier_weakening"
	DomainScopeRestriction   LossKind = "domain_restriction"
	ApproximationLoss        LossKind = "approximation"
	FunctionSpaceRestriction LossKind = "function_space_restriction"
)

type InformationLoss struct {
	Kind   LossKind `json:"kind"`
	Reason string   `json:"reason"`
}

// Transformation has one primary source/result edge. Premises are additional
// jointly required inputs; Obligations are proof conditions of the rule itself.
// This is the smallest multi-input extension needed by M1, not a general hypergraph.
type Transformation struct {
	ID          semantic.TransformationID
	Pass        string
	From        semantic.ClaimID
	Premises    []semantic.ClaimID
	To          semantic.ClaimID
	Relation    Relation
	Obligations []semantic.ClaimID
	Losses      []InformationLoss
	Provenance  semantic.Reference
	Theorem     semantic.TheoremID
	Bindings    []Binding
	Trusted     bool
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
	claim = cloneClaim(claim)
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
	requiredParents := []semantic.ClaimID{t.From}
	seenInputs := map[semantic.ClaimID]bool{t.From: true}
	for _, id := range append(append([]semantic.ClaimID(nil), t.Premises...), t.Obligations...) {
		if _, ok := g.claims[id]; !ok {
			return fmt.Errorf("transformation %q refers to unknown premise or obligation %q", t.ID, id)
		}
		if seenInputs[id] {
			return fmt.Errorf("transformation %q repeats input %q", t.ID, id)
		}
		seenInputs[id] = true
	}
	for _, id := range t.Premises {
		requiredParents = append(requiredParents, id)
	}
	for _, id := range t.Obligations {
		requiredParents = append(requiredParents, id)
	}
	for _, id := range append(append([]semantic.ClaimID{t.From}, t.Premises...), t.Obligations...) {
		if !assumptionsContain(to.Assumptions, g.claims[id].Assumptions) {
			return fmt.Errorf("transformation %q silently drops assumptions from %q", t.ID, id)
		}
	}
	if t.Relation == Equivalent && !assumptionsContain(from.Assumptions, to.Assumptions) {
		return fmt.Errorf("equivalence %q changes its assumption set", t.ID)
	}
	if t.Relation == Equivalent && from.Exactness != to.Exactness {
		return fmt.Errorf("equivalence %q changes exactness", t.ID)
	}
	for _, required := range structuralLosses(from.Proposition, to.Proposition) {
		if !hasLoss(t.Losses, required) {
			return fmt.Errorf("transformation %q fails to declare structural loss %s", t.ID, required)
		}
		if required == FunctionSpaceRestriction {
			if err := g.validateFunctionRestrictionPremises(from, to, t.Premises); err != nil {
				return fmt.Errorf("transformation %q: %w", t.ID, err)
			}
		}
	}
	if strengthening := structuralStrengthening(from.Proposition, to.Proposition); strengthening != "" && !(t.Relation == Equivalent && len(t.Obligations) > 0) {
		return fmt.Errorf("transformation %q attempts unsupported structural strengthening: %s", t.ID, strengthening)
	}
	if t.Relation == Approximation && !hasLoss(t.Losses, ApproximationLoss) {
		return fmt.Errorf("approximation %q fails to declare approximation loss", t.ID)
	}
	primaryDerivation := to.Provenance.Kind == semantic.DerivedProvenance && to.Provenance.Transformation == t.ID
	if !primaryDerivation && t.Theorem == "" {
		return fmt.Errorf("transformation %q target does not retain matching derived provenance", t.ID)
	}
	if primaryDerivation {
		for _, id := range requiredParents {
			if !containsID(to.Provenance.Parents, id) {
				return fmt.Errorf("transformation %q target provenance omits premise %q", t.ID, id)
			}
		}
	}
	g.transformations = append(g.transformations, cloneTransformation(t))
	return nil
}

func (g *Graph) validateFunctionRestrictionPremises(from, to semantic.Claim, premises []semantic.ClaimID) error {
	source := from.Proposition.(semantic.UniversalFunctionalStatement)
	target := to.Proposition.(semantic.UniversalFunctionalStatement)
	for _, member := range target.FunctionClass.Members {
		covered := false
		for _, id := range premises {
			claim := g.claims[id]
			admissible, ok := claim.Proposition.(semantic.TestFunctionAdmissibility)
			if ok && admissible.Function.Key() == member.Key() && admissible.Class.Key() == source.FunctionClass.Key() {
				if certified, _ := g.Certify(id); certified {
					covered = true
					break
				}
			}
		}
		if !covered {
			return fmt.Errorf("function-space restriction lacks a certified admissibility premise for %s", member.Symbol)
		}
	}
	return nil
}

func structuralLosses(from, to semantic.Proposition) []LossKind {
	if a, aok := from.(semantic.UniversalFunctionalStatement); aok {
		if b, bok := to.(semantic.UniversalFunctionalStatement); bok && sameFunctionalPredicate(a, b) && a.FunctionClass.Key() != b.FunctionClass.Key() && functionClassShapeSubset(b.FunctionClass, a.FunctionClass) {
			return []LossKind{FunctionSpaceRestriction}
		}
	}
	a, aok := semantic.Quantified(from)
	b, bok := semantic.Quantified(to)
	if !aok || !bok || a.Predicate != b.Predicate {
		return nil
	}
	var losses []LossKind
	if a.Quantifier == semantic.ForAll && b.Quantifier == semantic.DensityOne {
		losses = append(losses, QuantifierWeakening)
	}
	if a.Quantifier == b.Quantifier && a.Domain != b.Domain && semantic.IsSubset(b.Domain, a.Domain) {
		losses = append(losses, DomainScopeRestriction)
	}
	return losses
}

func structuralStrengthening(from, to semantic.Proposition) string {
	if a, aok := from.(semantic.UniversalFunctionalStatement); aok {
		if b, bok := to.(semantic.UniversalFunctionalStatement); bok && sameFunctionalPredicate(a, b) && a.FunctionClass.Key() != b.FunctionClass.Key() && !functionClassShapeSubset(b.FunctionClass, a.FunctionClass) {
			return "conclusion function class is not covered by source function class"
		}
	}
	a, aok := semantic.Quantified(from)
	b, bok := semantic.Quantified(to)
	if !aok || !bok || a.Predicate != b.Predicate {
		return ""
	}
	if a.Quantifier == semantic.DensityOne && b.Quantifier == semantic.ForAll {
		return "density-one to universal quantification"
	}
	if a.Quantifier == b.Quantifier && !semantic.IsSubset(b.Domain, a.Domain) {
		return "conclusion domain is not covered by source domain"
	}
	return ""
}

func hasLoss(losses []InformationLoss, kind LossKind) bool {
	for _, loss := range losses {
		if loss.Kind == kind {
			return true
		}
	}
	return false
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
	out := make([]semantic.Claim, 0, len(g.claimOrder))
	for _, id := range g.claimOrder {
		out = append(out, cloneClaim(g.claims[id]))
	}
	return out
}

func (g *Graph) ClaimBySemanticKey(key string) (semantic.Claim, bool) {
	for _, id := range g.claimOrder {
		if semantic.SemanticKey(g.claims[id].Proposition) == key {
			return cloneClaim(g.claims[id]), true
		}
	}
	return semantic.Claim{}, false
}
func (g *Graph) Transformations() []Transformation {
	out := make([]Transformation, len(g.transformations))
	for i, t := range g.transformations {
		out[i] = cloneTransformation(t)
	}
	return out
}

func cloneClaim(claim semantic.Claim) semantic.Claim {
	claim.Assumptions = semantic.CloneAssumptions(claim.Assumptions)
	claim.Evidence = append([]semantic.Evidence(nil), claim.Evidence...)
	claim.Provenance.Parents = append([]semantic.ClaimID(nil), claim.Provenance.Parents...)
	if p, ok := claim.Proposition.(semantic.RepresentationProposition); ok {
		p.Representation = semantic.CloneRepresentation(p.Representation)
		claim.Proposition = p
	}
	switch p := claim.Proposition.(type) {
	case semantic.FunctionalDefinition:
		p.Functional = semantic.CloneQuadraticFunctional(p.Functional)
		claim.Proposition = p
	case semantic.UniversalFunctionalStatement:
		p.FunctionClass = semantic.CloneFunctionClass(p.FunctionClass)
		claim.Proposition = p
	case semantic.TestFunctionAdmissibility:
		p.Function = semantic.CloneTestFunction(p.Function)
		p.Class = semantic.CloneFunctionClass(p.Class)
		claim.Proposition = p
	case semantic.ExplicitFormulaIdentity:
		p.FunctionClass = semantic.CloneFunctionClass(p.FunctionClass)
		p.ZeroSide = semantic.CloneAggregate(p.ZeroSide)
		p.ArithmeticSide = append([]semantic.FunctionalContribution(nil), p.ArithmeticSide...)
		claim.Proposition = p
	}
	return claim
}
func cloneTransformation(t Transformation) Transformation {
	t.Premises = append([]semantic.ClaimID(nil), t.Premises...)
	t.Obligations = append([]semantic.ClaimID(nil), t.Obligations...)
	t.Losses = append([]InformationLoss(nil), t.Losses...)
	t.Bindings = append([]Binding(nil), t.Bindings...)
	return t
}

type DiagnosticCode string

const (
	NoEstablishedDirection DiagnosticCode = "no_established_direction"
	QuantifierMismatch     DiagnosticCode = "quantifier_mismatch"
	DomainMismatch         DiagnosticCode = "domain_mismatch"
	PredicateMismatch      DiagnosticCode = "predicate_mismatch"
	InformationLost        DiagnosticCode = "information_lost_on_path"
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

func (g *Graph) AttemptProof(fromID, targetID semantic.ClaimID) ProofAttempt {
	attempt := g.CheckDischarge(fromID, targetID)
	if len(attempt.Diagnostics) != 0 {
		return attempt
	}
	if certified, diagnostics := g.Certify(fromID); !certified {
		attempt.Diagnostics = append(attempt.Diagnostics, diagnostics...)
		attempt.Accepted = false
	}
	return attempt
}

func (g *Graph) CheckDischarge(fromID, targetID semantic.ClaimID) ProofAttempt {
	attempt := ProofAttempt{From: fromID, Target: targetID}
	from, fromOK := g.claims[fromID]
	target, targetOK := g.claims[targetID]
	if !fromOK || !targetOK {
		attempt.Diagnostics = []Diagnostic{{NoEstablishedDirection, "source or target claim is unknown"}}
		return attempt
	}
	if from.Exactness == semantic.Approximate && target.Exactness == semantic.Exact {
		attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{ApproximationBoundary, "an approximate or numerical claim cannot certify an exact theorem target"})
	}
	path, found := g.findPath(fromID, targetID)
	if !found {
		message := fmt.Sprintf("no established transformation permits %s → %s", fromID, targetID)
		if reversePath, ok := g.findPath(targetID, fromID); ok {
			for _, step := range reversePath {
				if !step.Relation.Reversible() {
					message = fmt.Sprintf("%s is directional (%s → %s) and cannot be traversed backward", step.Relation, step.From, step.To)
				}
				for _, loss := range step.Losses {
					if g.lossBlocksTarget(step, loss.Kind, target.Proposition) {
						attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{InformationLost, fmt.Sprintf("the forward path loses %s: %s", loss.Kind, loss.Reason)})
					}
				}
			}
		}
		attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{NoEstablishedDirection, message})
		attempt.Diagnostics = append(attempt.Diagnostics, compareStructuralStrength(from.Proposition, target.Proposition)...)
	}
	if found {
		for _, step := range path {
			for _, loss := range step.Losses {
				if g.lossBlocksTarget(step, loss.Kind, target.Proposition) {
					attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{InformationLost, fmt.Sprintf("path loses %s: %s", loss.Kind, loss.Reason)})
				}
			}
			if step.Relation == Approximation && target.Exactness == semantic.Exact {
				attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{ApproximationBoundary, "an approximation cannot certify an exact theorem target"})
			}
			for _, id := range append(append([]semantic.ClaimID(nil), step.Premises...), step.Obligations...) {
				if certified, _ := g.certify(id, make(map[semantic.ClaimID]bool)); !certified {
					attempt.Diagnostics = append(attempt.Diagnostics, Diagnostic{OpenObligation, fmt.Sprintf("required premise or proof obligation %s remains open", id)})
				}
			}
		}
	}
	attempt.Diagnostics = deduplicateDiagnostics(attempt.Diagnostics)
	attempt.Accepted = found && len(attempt.Diagnostics) == 0
	return attempt
}

func compareStructuralStrength(from, target semantic.Proposition) []Diagnostic {
	if a, aok := from.(semantic.UniversalFunctionalStatement); aok {
		if b, bok := target.(semantic.UniversalFunctionalStatement); bok {
			var out []Diagnostic
			if a.Functional != b.Functional || a.Predicate != b.Predicate {
				out = append(out, Diagnostic{PredicateMismatch, "source and target use different functional predicates"})
			}
			if a.TransformConvention != b.TransformConvention {
				out = append(out, Diagnostic{PredicateMismatch, "source and target use incompatible transform conventions"})
			}
			if !functionClassShapeSubset(b.FunctionClass, a.FunctionClass) {
				out = append(out, Diagnostic{DomainMismatch, fmt.Sprintf("source function class %s does not cover target function class %s", a.FunctionClass.Describe(), b.FunctionClass.Describe())})
			}
			return out
		}
	}
	a, aok := semantic.Quantified(from)
	b, bok := semantic.Quantified(target)
	if !aok || !bok {
		return nil
	}
	var out []Diagnostic
	if a.Predicate != b.Predicate {
		out = append(out, Diagnostic{PredicateMismatch, "source and target assert different predicates"})
	}
	if a.Quantifier != b.Quantifier {
		message := fmt.Sprintf("%s quantification cannot discharge %s quantification", a.Quantifier, b.Quantifier)
		if a.Quantifier == semantic.DensityOne && b.Quantifier == semantic.ForAll {
			message = "asymptotic/density-one quantification cannot discharge a universal target"
		}
		out = append(out, Diagnostic{QuantifierMismatch, message})
	}
	if !semantic.IsSubset(b.Domain, a.Domain) {
		out = append(out, Diagnostic{DomainMismatch, fmt.Sprintf("source domain %s does not cover target domain %s", a.Domain.Describe(), b.Domain.Describe())})
	}
	return out
}

func (g *Graph) lossBlocksTarget(step Transformation, kind LossKind, target semantic.Proposition) bool {
	q, ok := semantic.Quantified(target)
	if !ok {
		return false
	}
	if kind == QuantifierWeakening && q.Quantifier == semantic.ForAll {
		return true
	}
	if kind == DomainScopeRestriction {
		narrowed, ok := semantic.Quantified(g.claims[step.To].Proposition)
		return ok && !semantic.IsSubset(q.Domain, narrowed.Domain)
	}
	if kind == FunctionSpaceRestriction {
		targetFunctional, ok := target.(semantic.UniversalFunctionalStatement)
		narrowed, narrowedOK := g.claims[step.To].Proposition.(semantic.UniversalFunctionalStatement)
		// A loss of universal function-space coverage cannot disappear merely
		// because a later equivalent representation (such as RH) changes IR family.
		return !ok || (narrowedOK && !functionClassShapeSubset(targetFunctional.FunctionClass, narrowed.FunctionClass))
	}
	return kind == ApproximationLoss
}

func sameFunctionalPredicate(a, b semantic.UniversalFunctionalStatement) bool {
	return a.Quantifier == b.Quantifier && a.Functional == b.Functional && a.Predicate == b.Predicate && a.TransformConvention == b.TransformConvention
}

// Shape inclusion is intentionally small. Membership evidence is checked by
// FunctionClassRestriction before a graph edge is admitted.
func functionClassShapeSubset(sub, sup semantic.FunctionClass) bool {
	if sub.Key() == sup.Key() {
		return true
	}
	return sub.Kind == semantic.FiniteFunctionClass && sup.Kind == semantic.WeilNiceFunctionClass && sub.TransformConvention == sup.TransformConvention
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
			if t.From == current.id {
				next = t.To
			} else if t.Relation.Reversible() && t.To == current.id {
				next = t.From
			} else {
				continue
			}
			if seen[next] {
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
	for _, evidence := range claim.Evidence {
		if evidence.Kind == semantic.DefinitionEvidence || evidence.Kind == semantic.KnownTheoremEvidence {
			return true, nil
		}
	}
	visiting[id] = true
	defer delete(visiting, id)
	var diagnostics []Diagnostic
	for _, t := range g.transformations {
		var primary semantic.ClaimID
		switch {
		case t.To == id:
			primary = t.From
		case t.Relation.Reversible() && t.From == id:
			primary = t.To
		default:
			continue
		}
		if t.Relation == Approximation && claim.Exactness == semantic.Exact {
			diagnostics = append(diagnostics, Diagnostic{ApproximationBoundary, "an approximation cannot certify an exact theorem target"})
			continue
		}
		if t.Theorem != "" && !t.Trusted {
			diagnostics = append(diagnostics, Diagnostic{UncertifiedEvidence, fmt.Sprintf("theorem contract %s is not trusted", t.Theorem)})
			continue
		}
		blocked := false
		for _, required := range append(append([]semantic.ClaimID{primary}, t.Premises...), t.Obligations...) {
			if certified, _ := g.certify(required, visiting); !certified {
				diagnostics = append(diagnostics, Diagnostic{OpenObligation, fmt.Sprintf("required premise or proof obligation %s remains open", required)})
				blocked = true
			}
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
		parents := append([]semantic.ClaimID(nil), g.claims[current].Provenance.Parents...)
		for _, transformation := range g.transformations {
			if transformation.To == current {
				parents = append(parents, transformation.From)
				parents = append(parents, transformation.Premises...)
				parents = append(parents, transformation.Obligations...)
			}
		}
		for _, parent := range sortedIDs(uniqueIDs(parents)) {
			visit(parent)
		}
		out = append(out, current)
	}
	visit(id)
	return out, nil
}

func uniqueIDs(ids []semantic.ClaimID) []semantic.ClaimID {
	seen := make(map[semantic.ClaimID]bool, len(ids))
	out := make([]semantic.ClaimID, 0, len(ids))
	for _, id := range ids {
		if id != "" && !seen[id] {
			seen[id] = true
			out = append(out, id)
		}
	}
	return out
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
