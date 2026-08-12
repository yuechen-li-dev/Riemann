package compiler

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

type ParamID string
type ParamType string

const (
	ObjectParam         ParamType = "object"
	DomainParam         ParamType = "domain"
	RepresentationParam ParamType = "representation"
	ScalarParam         ParamType = "scalar"
	QuantifierParam     ParamType = "quantifier"
	PredicateParam      ParamType = "predicate"
	AnalyticFactParam   ParamType = "analytic_fact"
)

type Parameter struct {
	ID   ParamID   `json:"id"`
	Type ParamType `json:"type"`
}

// BindingValue is a closed tagged union. Only the field selected by Type is
// semantic; keeping it comparable makes repeated-variable unification exact.
type BindingValue struct {
	Type           ParamType                   `json:"type"`
	Object         semantic.Function           `json:"object,omitempty"`
	Domain         semantic.Domain             `json:"domain,omitempty"`
	Representation semantic.RepresentationName `json:"representation,omitempty"`
	Scalar         uint64                      `json:"scalar,omitempty"`
	Quantifier     semantic.QuantifierKind     `json:"quantifier,omitempty"`
	Predicate      semantic.PredicateKind      `json:"predicate,omitempty"`
	AnalyticFact   semantic.AnalyticFactName   `json:"analytic_fact,omitempty"`
}

type Term struct {
	Parameter ParamID       `json:"parameter,omitempty"`
	Constant  *BindingValue `json:"constant,omitempty"`
}

func Var(id ParamID) Term { return Term{Parameter: id} }
func ConstObject(v semantic.Function) Term {
	return constant(BindingValue{Type: ObjectParam, Object: v})
}
func ConstDomain(v semantic.Domain) Term {
	return constant(BindingValue{Type: DomainParam, Domain: v})
}
func ConstRepresentation(v semantic.RepresentationName) Term {
	return constant(BindingValue{Type: RepresentationParam, Representation: v})
}
func ConstQuantifier(v semantic.QuantifierKind) Term {
	return constant(BindingValue{Type: QuantifierParam, Quantifier: v})
}
func ConstPredicate(v semantic.PredicateKind) Term {
	return constant(BindingValue{Type: PredicateParam, Predicate: v})
}
func ConstAnalyticFact(v semantic.AnalyticFactName) Term {
	return constant(BindingValue{Type: AnalyticFactParam, AnalyticFact: v})
}
func constant(v BindingValue) Term { return Term{Constant: &v} }

// Pattern is a typed semantic proposition template. Fields irrelevant to Kind
// must be absent; Formula and Affordances are display metadata, never match keys.
type Pattern struct {
	Kind           semantic.PropositionKind `json:"kind"`
	Object         Term                     `json:"object,omitempty"`
	Domain         Term                     `json:"domain,omitempty"`
	Representation Term                     `json:"representation,omitempty"`
	Left           Term                     `json:"left,omitempty"`
	Right          Term                     `json:"right,omitempty"`
	Quantifier     Term                     `json:"quantifier,omitempty"`
	Predicate      Term                     `json:"predicate,omitempty"`
	AnalyticFact   Term                     `json:"analytic_fact,omitempty"`
	Exactness      semantic.Exactness       `json:"exactness"`
	Formula        string                   `json:"formula,omitempty"`
	Affordances    []string                 `json:"affordances,omitempty"`
}

type TheoremTrust string

const (
	TrustedExternalTheorem TheoremTrust = "trusted_external_theorem"
	UntrustedTheorem       TheoremTrust = "untrusted"
)

type TheoremContract struct {
	ID           semantic.TheoremID `json:"id"`
	Parameters   []Parameter        `json:"parameters"`
	Premises     []Pattern          `json:"premises"`
	Conclusion   Pattern            `json:"conclusion"`
	ConclusionID semantic.ClaimID   `json:"conclusion_id,omitempty"`
	Relation     Relation           `json:"relation"`
	Evidence     semantic.Evidence  `json:"evidence"`
	Trust        TheoremTrust       `json:"trust"`
	Citation     string             `json:"citation,omitempty"`
}

type Binding struct {
	Parameter ParamID      `json:"parameter"`
	Value     BindingValue `json:"value"`
}

type PremiseMatch struct {
	Premise int              `json:"premise"`
	Claim   semantic.ClaimID `json:"claim"`
}

type TheoremObligation struct {
	Premise     int                  `json:"premise"`
	Pattern     Pattern              `json:"pattern"`
	Proposition semantic.Proposition `json:"-"`
	Description string               `json:"description"`
}

type TheoremApplication struct {
	Theorem        semantic.TheoremID        `json:"theorem"`
	Bindings       []Binding                 `json:"bindings"`
	Matched        []PremiseMatch            `json:"matched_premises"`
	Conclusion     semantic.ClaimID          `json:"conclusion,omitempty"`
	Obligations    []TheoremObligation       `json:"obligations"`
	Complete       bool                      `json:"complete"`
	Transformation semantic.TransformationID `json:"transformation,omitempty"`
	Provenance     semantic.Reference        `json:"provenance"`
}

type ContractRegistry struct {
	contracts map[semantic.TheoremID]TheoremContract
}

func NewContractRegistry() *ContractRegistry {
	return &ContractRegistry{contracts: make(map[semantic.TheoremID]TheoremContract)}
}

func (r *ContractRegistry) Register(contract TheoremContract) error {
	if err := validateContract(contract); err != nil {
		return err
	}
	if _, exists := r.contracts[contract.ID]; exists {
		return fmt.Errorf("theorem contract %q already registered", contract.ID)
	}
	r.contracts[contract.ID] = cloneContract(contract)
	return nil
}

func (r *ContractRegistry) Contract(id semantic.TheoremID) (TheoremContract, bool) {
	c, ok := r.contracts[id]
	return cloneContract(c), ok
}

func (r *ContractRegistry) Contracts() []TheoremContract {
	ids := make([]semantic.TheoremID, 0, len(r.contracts))
	for id := range r.contracts {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	out := make([]TheoremContract, 0, len(ids))
	for _, id := range ids {
		out = append(out, cloneContract(r.contracts[id]))
	}
	return out
}

func cloneContract(c TheoremContract) TheoremContract {
	c.Parameters = append([]Parameter(nil), c.Parameters...)
	c.Premises = append([]Pattern(nil), c.Premises...)
	for i := range c.Premises {
		c.Premises[i].Affordances = append([]string(nil), c.Premises[i].Affordances...)
	}
	c.Conclusion.Affordances = append([]string(nil), c.Conclusion.Affordances...)
	return c
}

func validateContract(c TheoremContract) error {
	if c.ID == "" || !c.Relation.Valid() || (c.Trust != TrustedExternalTheorem && c.Trust != UntrustedTheorem) {
		return fmt.Errorf("invalid theorem contract metadata")
	}
	declared := make(map[ParamID]ParamType)
	for _, p := range c.Parameters {
		if p.ID == "" || !validParamType(p.Type) || declared[p.ID] != "" {
			return fmt.Errorf("theorem %s has invalid parameter %q", c.ID, p.ID)
		}
		declared[p.ID] = p.Type
	}
	for i, p := range append(append([]Pattern(nil), c.Premises...), c.Conclusion) {
		if err := validatePattern(p, declared); err != nil {
			return fmt.Errorf("theorem %s pattern %d: %w", c.ID, i, err)
		}
	}
	return nil
}

func validParamType(t ParamType) bool {
	return t == ObjectParam || t == DomainParam || t == RepresentationParam || t == ScalarParam || t == QuantifierParam || t == PredicateParam || t == AnalyticFactParam
}

func validatePattern(p Pattern, declared map[ParamID]ParamType) error {
	if p.Exactness != semantic.Exact && p.Exactness != semantic.Approximate {
		return fmt.Errorf("invalid exactness")
	}
	required := []struct {
		term Term
		kind ParamType
	}{}
	switch p.Kind {
	case semantic.RepresentationKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Representation, RepresentationParam}, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam})
	case semantic.RepresentationIdentityKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Left, RepresentationParam}, struct {
			term Term
			kind ParamType
		}{p.Right, RepresentationParam}, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam})
	case semantic.AnalyticFactKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.AnalyticFact, AnalyticFactParam}, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam})
	case semantic.QuantifiedStatementKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Quantifier, QuantifierParam}, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam}, struct {
			term Term
			kind ParamType
		}{p.Predicate, PredicateParam}, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam})
	default:
		return fmt.Errorf("unsupported proposition pattern kind %q", p.Kind)
	}
	for _, item := range required {
		if err := validateTerm(item.term, item.kind, declared); err != nil {
			return err
		}
	}
	return nil
}

func validateTerm(t Term, want ParamType, declared map[ParamID]ParamType) error {
	if (t.Parameter == "") == (t.Constant == nil) {
		return fmt.Errorf("term must contain exactly one parameter or constant")
	}
	if t.Parameter != "" {
		if declared[t.Parameter] != want {
			return fmt.Errorf("parameter %s has type %s, want %s", t.Parameter, declared[t.Parameter], want)
		}
		return nil
	}
	if t.Constant.Type != want {
		return fmt.Errorf("constant has type %s, want %s", t.Constant.Type, want)
	}
	return nil
}

type bindingEnv map[ParamID]BindingValue

func (e bindingEnv) clone() bindingEnv {
	out := make(bindingEnv, len(e))
	for k, v := range e {
		out[k] = v
	}
	return out
}
func (e bindingEnv) bind(t Term, value BindingValue) bool {
	if t.Constant != nil {
		return *t.Constant == value
	}
	if existing, ok := e[t.Parameter]; ok {
		return existing == value
	}
	e[t.Parameter] = value
	return true
}

func matchPattern(p Pattern, claim semantic.Claim, env bindingEnv) bool {
	if claim.Exactness != p.Exactness || claim.Proposition.Kind() != p.Kind {
		return false
	}
	switch v := claim.Proposition.(type) {
	case semantic.RepresentationProposition:
		r := v.Representation
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: r.Object}) && env.bind(p.Representation, BindingValue{Type: RepresentationParam, Representation: r.Name}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: r.ValidOn})
	case semantic.RepresentationIdentity:
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Left, BindingValue{Type: RepresentationParam, Representation: v.Left}) && env.bind(p.Right, BindingValue{Type: RepresentationParam, Representation: v.Right}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.AnalyticFact:
		return env.bind(p.AnalyticFact, BindingValue{Type: AnalyticFactParam, AnalyticFact: v.Fact}) && env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.QuantifiedStatement:
		return env.bind(p.Quantifier, BindingValue{Type: QuantifierParam, Quantifier: v.Quantifier}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain}) && env.bind(p.Predicate, BindingValue{Type: PredicateParam, Predicate: v.Predicate.Kind}) && env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Predicate.Function})
	}
	return false
}

func resolve(t Term, env bindingEnv) (BindingValue, error) {
	if t.Constant != nil {
		return *t.Constant, nil
	}
	v, ok := env[t.Parameter]
	if !ok {
		return BindingValue{}, fmt.Errorf("parameter %s is unbound", t.Parameter)
	}
	return v, nil
}

func instantiate(p Pattern, env bindingEnv) (semantic.Proposition, error) {
	object, err := resolve(p.Object, env)
	if err != nil {
		return nil, err
	}
	domain, err := resolve(p.Domain, env)
	if err != nil {
		return nil, err
	}
	switch p.Kind {
	case semantic.RepresentationKind:
		r, err := resolve(p.Representation, env)
		if err != nil {
			return nil, err
		}
		return semantic.RepresentationProposition{Representation: semantic.Representation{Object: object.Object, Name: r.Representation, ValidOn: domain.Domain, Formula: p.Formula, Affordances: append([]string(nil), p.Affordances...)}}, nil
	case semantic.RepresentationIdentityKind:
		left, err := resolve(p.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := resolve(p.Right, env)
		if err != nil {
			return nil, err
		}
		return semantic.RepresentationIdentity{Object: object.Object, Left: left.Representation, Right: right.Representation, Domain: domain.Domain}, nil
	case semantic.AnalyticFactKind:
		fact, err := resolve(p.AnalyticFact, env)
		if err != nil {
			return nil, err
		}
		return semantic.AnalyticFact{Fact: fact.AnalyticFact, Object: object.Object, Domain: domain.Domain}, nil
	case semantic.QuantifiedStatementKind:
		q, err := resolve(p.Quantifier, env)
		if err != nil {
			return nil, err
		}
		pred, err := resolve(p.Predicate, env)
		if err != nil {
			return nil, err
		}
		return semantic.QuantifiedStatement{Quantifier: q.Quantifier, Domain: domain.Domain, Predicate: semantic.Predicate{Kind: pred.Predicate, Function: object.Object}}, nil
	default:
		return nil, fmt.Errorf("unsupported pattern kind %s", p.Kind)
	}
}

type ContractEngine struct {
	Graph        *Graph
	Registry     *ContractRegistry
	Applications []TheoremApplication
	applied      map[string]bool
}

func NewContractEngine(g *Graph, registry *ContractRegistry) *ContractEngine {
	return &ContractEngine{Graph: g, Registry: registry, applied: make(map[string]bool)}
}

type matchState struct {
	env     bindingEnv
	matched []PremiseMatch
}

func (e *ContractEngine) Saturate() error {
	for {
		changed := false
		for _, contract := range e.Registry.Contracts() {
			for _, state := range e.fullMatches(contract) {
				key := applicationKey(contract.ID, state.env, state.matched)
				if e.applied[key] {
					continue
				}
				added, err := e.apply(contract, state, key)
				if err != nil {
					return err
				}
				e.applied[key] = true
				changed = changed || added
			}
		}
		if !changed {
			break
		}
	}
	e.addPartialApplications()
	return nil
}

func (e *ContractEngine) fullMatches(c TheoremContract) []matchState {
	states := []matchState{{env: make(bindingEnv)}}
	claims := e.sortedClaims()
	for premiseIndex, pattern := range c.Premises {
		var next []matchState
		for _, state := range states {
			for _, claim := range claims {
				env := state.env.clone()
				if matchPattern(pattern, claim, env) {
					matched := append(append([]PremiseMatch(nil), state.matched...), PremiseMatch{Premise: premiseIndex, Claim: claim.ID})
					next = append(next, matchState{env: env, matched: matched})
				}
			}
		}
		states = next
		if len(states) == 0 {
			break
		}
	}
	return states
}

func (e *ContractEngine) sortedClaims() []semantic.Claim {
	all := e.Graph.Claims()
	claims := make([]semantic.Claim, 0, len(all))
	for _, claim := range all {
		if certified, _ := e.Graph.Certify(claim.ID); certified {
			claims = append(claims, claim)
		}
	}
	sort.Slice(claims, func(i, j int) bool {
		a, b := semantic.SemanticKey(claims[i].Proposition), semantic.SemanticKey(claims[j].Proposition)
		if a == b {
			return claims[i].ID < claims[j].ID
		}
		return a < b
	})
	return claims
}

func (e *ContractEngine) apply(c TheoremContract, state matchState, key string) (bool, error) {
	proposition, err := instantiate(c.Conclusion, state.env)
	if err != nil {
		return false, err
	}
	semanticKey := semantic.SemanticKey(proposition)
	claim, exists := e.Graph.ClaimBySemanticKey(semanticKey)
	claimID := c.ConclusionID
	if claimID == "" {
		claimID = semantic.ClaimID(string(c.ID) + ":" + shortHash(semanticKey))
	}
	parents := make([]semantic.ClaimID, len(state.matched))
	for i, m := range state.matched {
		parents[i] = m.Claim
	}
	transformID := semantic.TransformationID("")
	if len(parents) > 0 {
		transformID = semantic.TransformationID("theorem:" + string(c.ID) + ":" + shortHash(key))
	}
	if exists {
		claimID = claim.ID
	} else {
		if other, ok := e.Graph.Claim(claimID); ok && semantic.SemanticKey(other.Proposition) != semanticKey {
			return false, fmt.Errorf("theorem %s conclusion ID %s collides with another semantic claim", c.ID, claimID)
		}
		evidence := []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: c.Evidence.Source, Note: "instantiated from typed theorem contract " + string(c.ID)}}
		if len(parents) == 0 {
			if c.Trust == TrustedExternalTheorem {
				evidence = []semantic.Evidence{c.Evidence}
			} else {
				evidence = []semantic.Evidence{{Kind: semantic.UnverifiedConjectureEvidence, Source: c.Evidence.Source, Note: "untrusted theorem contract " + string(c.ID)}}
			}
		}
		assumptions := e.unionAssumptions(parents)
		newClaim := semantic.Claim{ID: claimID, Proposition: proposition, Assumptions: assumptions, Evidence: evidence, Exactness: c.Conclusion.Exactness, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: append([]semantic.ClaimID(nil), parents...), Transformation: transformID, Theorem: c.ID, Source: c.Evidence.Source}}
		if err := e.Graph.AddClaim(newClaim); err != nil {
			return false, err
		}
	}
	bindings := sortedBindings(state.env)
	if len(parents) > 0 {
		t := Transformation{ID: transformID, Pass: "instantiate-theorem-contract", From: parents[0], To: claimID, Relation: c.Relation, Provenance: c.Evidence.Source, Theorem: c.ID, Bindings: bindings, Trusted: c.Trust == TrustedExternalTheorem}
		if len(parents) > 1 {
			t.Premises = append([]semantic.ClaimID(nil), parents[1:]...)
		}
		if err := e.Graph.AddTransformation(t); err != nil {
			return false, err
		}
	}
	e.Applications = append(e.Applications, TheoremApplication{Theorem: c.ID, Bindings: bindings, Matched: append([]PremiseMatch(nil), state.matched...), Conclusion: claimID, Complete: true, Transformation: transformID, Provenance: c.Evidence.Source})
	return !exists, nil
}

func (e *ContractEngine) unionAssumptions(ids []semantic.ClaimID) []semantic.Assumption {
	byID := make(map[semantic.AssumptionID]semantic.Assumption)
	for _, id := range ids {
		c, _ := e.Graph.Claim(id)
		for _, a := range c.Assumptions {
			byID[a.ID] = a
		}
	}
	keys := make([]semantic.AssumptionID, 0, len(byID))
	for id := range byID {
		keys = append(keys, id)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	out := make([]semantic.Assumption, 0, len(keys))
	for _, id := range keys {
		out = append(out, byID[id])
	}
	return out
}

func (e *ContractEngine) addPartialApplications() {
	for _, c := range e.Registry.Contracts() {
		if len(c.Premises) == 0 || len(e.fullMatches(c)) > 0 {
			continue
		}
		state := matchState{env: make(bindingEnv)}
		used := make(map[semantic.ClaimID]bool)
		for i, p := range c.Premises {
			for _, claim := range e.sortedClaims() {
				if used[claim.ID] {
					continue
				}
				env := state.env.clone()
				if matchPattern(p, claim, env) {
					state.env = env
					state.matched = append(state.matched, PremiseMatch{Premise: i, Claim: claim.ID})
					used[claim.ID] = true
					break
				}
			}
		}
		if len(state.matched) == 0 {
			continue
		}
		matchedIndex := make(map[int]bool)
		for _, m := range state.matched {
			matchedIndex[m.Premise] = true
		}
		var obligations []TheoremObligation
		for i, p := range c.Premises {
			if matchedIndex[i] {
				continue
			}
			o := TheoremObligation{Premise: i, Pattern: p}
			if prop, err := instantiate(p, state.env); err == nil {
				o.Proposition = prop
				o.Description = prop.Describe()
			} else {
				o.Description = describePattern(p, state.env)
			}
			obligations = append(obligations, o)
		}
		e.Applications = append(e.Applications, TheoremApplication{Theorem: c.ID, Bindings: sortedBindings(state.env), Matched: state.matched, Obligations: obligations, Complete: false, Provenance: c.Evidence.Source})
	}
}

func sortedBindings(env bindingEnv) []Binding {
	ids := make([]ParamID, 0, len(env))
	for id := range env {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	out := make([]Binding, 0, len(ids))
	for _, id := range ids {
		out = append(out, Binding{Parameter: id, Value: env[id]})
	}
	return out
}
func applicationKey(id semantic.TheoremID, env bindingEnv, matches []PremiseMatch) string {
	var b strings.Builder
	b.WriteString(string(id))
	for _, x := range sortedBindings(env) {
		fmt.Fprintf(&b, "|%s=%v", x.Parameter, x.Value)
	}
	for _, m := range matches {
		fmt.Fprintf(&b, "|%d=%s", m.Premise, m.Claim)
	}
	return b.String()
}
func shortHash(s string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte(s))
	return fmt.Sprintf("%016x", h.Sum64())
}
func describePattern(p Pattern, env bindingEnv) string {
	return fmt.Sprintf("%s pattern with retained bindings %v", p.Kind, sortedBindings(env))
}
