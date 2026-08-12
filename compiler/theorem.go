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
	ObjectParam             ParamType = "object"
	DomainParam             ParamType = "domain"
	RepresentationParam     ParamType = "representation"
	ScalarParam             ParamType = "scalar"
	QuantifierParam         ParamType = "quantifier"
	PredicateParam          ParamType = "predicate"
	AnalyticFactParam       ParamType = "analytic_fact"
	PointParam              ParamType = "point"
	ZeroClassParam          ParamType = "zero_class"
	SideConditionParam      ParamType = "side_condition"
	TransformParam          ParamType = "point_transform"
	ZeroSetPropertyParam    ParamType = "zero_set_property"
	ZeroClassificationParam ParamType = "zero_classification"
)

type Parameter struct {
	ID   ParamID   `json:"id"`
	Type ParamType `json:"type"`
}

// BindingValue is a closed tagged union. Only the field selected by Type is
// semantic; keeping it comparable makes repeated-variable unification exact.
type BindingValue struct {
	Type               ParamType                       `json:"type"`
	Object             semantic.Function               `json:"object,omitempty"`
	Domain             semantic.Domain                 `json:"domain,omitempty"`
	Representation     semantic.RepresentationName     `json:"representation,omitempty"`
	Scalar             uint64                          `json:"scalar,omitempty"`
	Quantifier         semantic.QuantifierKind         `json:"quantifier,omitempty"`
	Predicate          semantic.PredicateKind          `json:"predicate,omitempty"`
	AnalyticFact       semantic.AnalyticFactName       `json:"analytic_fact,omitempty"`
	Point              semantic.PointExpr              `json:"point,omitempty"`
	ZeroClass          semantic.ZeroClass              `json:"zero_class,omitempty"`
	SideCondition      semantic.SideConditionName      `json:"side_condition,omitempty"`
	Transform          semantic.PointTransform         `json:"point_transform,omitempty"`
	ZeroSetProperty    semantic.ZeroSetPropertyName    `json:"zero_set_property,omitempty"`
	ZeroClassification semantic.ZeroClassificationName `json:"zero_classification,omitempty"`
}

type Term struct {
	Parameter ParamID       `json:"parameter,omitempty"`
	Constant  *BindingValue `json:"constant,omitempty"`
	// Transform is meaningful only for point terms. It records transport in
	// the contract itself and is resolved through semantic transform composition.
	Transform          semantic.PointTransform `json:"point_transform,omitempty"`
	TransformParameter ParamID                 `json:"point_transform_parameter,omitempty"`
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
func ConstPoint(v semantic.PointExpr) Term {
	return constant(BindingValue{Type: PointParam, Point: v.Canonical()})
}
func TransformedPoint(id ParamID, transform semantic.PointTransform) Term {
	return Term{Parameter: id, Transform: transform}
}
func PointTransformedBy(id, transform ParamID) Term {
	return Term{Parameter: id, TransformParameter: transform}
}
func ConstZeroClass(v semantic.ZeroClass) Term {
	return constant(BindingValue{Type: ZeroClassParam, ZeroClass: v})
}
func ConstSideCondition(v semantic.SideConditionName) Term {
	return constant(BindingValue{Type: SideConditionParam, SideCondition: v})
}
func ConstTransform(v semantic.PointTransform) Term {
	return constant(BindingValue{Type: TransformParam, Transform: v})
}
func ConstZeroSetProperty(v semantic.ZeroSetPropertyName) Term {
	return constant(BindingValue{Type: ZeroSetPropertyParam, ZeroSetProperty: v})
}
func ConstZeroClassification(v semantic.ZeroClassificationName) Term {
	return constant(BindingValue{Type: ZeroClassificationParam, ZeroClassification: v})
}
func constant(v BindingValue) Term { return Term{Constant: &v} }

// Pattern is a typed semantic proposition template. Fields irrelevant to Kind
// must be absent; Formula and Affordances are display metadata, never match keys.
type Pattern struct {
	Kind           semantic.PropositionKind `json:"kind"`
	Object         Term                     `json:"object,omitempty"`
	BaseObject     Term                     `json:"base_object,omitempty"`
	Domain         Term                     `json:"domain,omitempty"`
	Representation Term                     `json:"representation,omitempty"`
	Left           Term                     `json:"left,omitempty"`
	Right          Term                     `json:"right,omitempty"`
	Quantifier     Term                     `json:"quantifier,omitempty"`
	Predicate      Term                     `json:"predicate,omitempty"`
	AnalyticFact   Term                     `json:"analytic_fact,omitempty"`
	Point          Term                     `json:"point,omitempty"`
	Classification Term                     `json:"classification,omitempty"`
	SideCondition  Term                     `json:"side_condition,omitempty"`
	Transform      Term                     `json:"point_transform,omitempty"`
	Property       Term                     `json:"property,omitempty"`
	Region         Term                     `json:"region,omitempty"`
	Exactness      semantic.Exactness       `json:"exactness"`
	Formula        string                   `json:"formula,omitempty"`
	Affordances    []string                 `json:"affordances,omitempty"`
	// M4's first functional theorem is concrete. Keeping these fields typed and
	// non-parametric avoids pretending the compiler has a general binder calculus.
	FunctionClass       *semantic.FunctionClass          `json:"function_class,omitempty"`
	Functional          semantic.FunctionalID            `json:"functional,omitempty"`
	FunctionalPredicate semantic.FunctionalPredicateKind `json:"functional_predicate,omitempty"`
	TransformConvention semantic.TransformConventionID   `json:"transform_convention,omitempty"`
}

type TheoremTrust string

const (
	TrustedExternalTheorem TheoremTrust = "trusted_external_theorem"
	CompilerVerifiedRule   TheoremTrust = "compiler_verified_rule"
	UntrustedTheorem       TheoremTrust = "untrusted"
)

type TheoremContract struct {
	ID             semantic.TheoremID `json:"id"`
	Parameters     []Parameter        `json:"parameters"`
	Premises       []Pattern          `json:"premises"`
	SideConditions []Pattern          `json:"side_conditions"`
	Conclusion     Pattern            `json:"conclusion"`
	ConclusionID   semantic.ClaimID   `json:"conclusion_id,omitempty"`
	Relation       Relation           `json:"relation"`
	Evidence       semantic.Evidence  `json:"evidence"`
	Trust          TheoremTrust       `json:"trust"`
	Citation       string             `json:"citation,omitempty"`
}

type Binding struct {
	Parameter ParamID      `json:"parameter"`
	Value     BindingValue `json:"value"`
}

type PremiseMatch struct {
	Premise       int              `json:"premise"`
	Claim         semantic.ClaimID `json:"claim"`
	SideCondition bool             `json:"side_condition,omitempty"`
}

type TheoremObligation struct {
	Premise       int                  `json:"premise"`
	Pattern       Pattern              `json:"pattern"`
	Proposition   semantic.Proposition `json:"-"`
	Description   string               `json:"description"`
	SideCondition bool                 `json:"side_condition,omitempty"`
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
	c.SideConditions = append([]Pattern(nil), c.SideConditions...)
	for i := range c.Premises {
		c.Premises[i].Affordances = append([]string(nil), c.Premises[i].Affordances...)
		if c.Premises[i].FunctionClass != nil {
			x := semantic.CloneFunctionClass(*c.Premises[i].FunctionClass)
			c.Premises[i].FunctionClass = &x
		}
	}
	for i := range c.SideConditions {
		if c.SideConditions[i].FunctionClass != nil {
			x := semantic.CloneFunctionClass(*c.SideConditions[i].FunctionClass)
			c.SideConditions[i].FunctionClass = &x
		}
	}
	c.Conclusion.Affordances = append([]string(nil), c.Conclusion.Affordances...)
	if c.Conclusion.FunctionClass != nil {
		x := semantic.CloneFunctionClass(*c.Conclusion.FunctionClass)
		c.Conclusion.FunctionClass = &x
	}
	return c
}

func validateContract(c TheoremContract) error {
	if c.ID == "" || !c.Relation.Valid() || (c.Trust != TrustedExternalTheorem && c.Trust != CompilerVerifiedRule && c.Trust != UntrustedTheorem) {
		return fmt.Errorf("invalid theorem contract metadata")
	}
	declared := make(map[ParamID]ParamType)
	for _, p := range c.Parameters {
		if p.ID == "" || !validParamType(p.Type) || declared[p.ID] != "" {
			return fmt.Errorf("theorem %s has invalid parameter %q", c.ID, p.ID)
		}
		declared[p.ID] = p.Type
	}
	patterns := append(append(append([]Pattern(nil), c.Premises...), c.SideConditions...), c.Conclusion)
	for i, p := range patterns {
		if err := validatePattern(p, declared); err != nil {
			return fmt.Errorf("theorem %s pattern %d: %w", c.ID, i, err)
		}
	}
	return nil
}

func validParamType(t ParamType) bool {
	return t == ObjectParam || t == DomainParam || t == RepresentationParam || t == ScalarParam || t == QuantifierParam || t == PredicateParam || t == AnalyticFactParam || t == PointParam || t == ZeroClassParam || t == SideConditionParam || t == TransformParam || t == ZeroSetPropertyParam || t == ZeroClassificationParam
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
		if p.Representation.Constant != nil && p.Representation.Constant.Representation == semantic.CompletedXiRepresentation {
			required = append(required, struct {
				term Term
				kind ParamType
			}{p.BaseObject, ObjectParam})
		}
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
	case semantic.ZeroAtPointKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Point, PointParam}, struct {
			term Term
			kind ParamType
		}{p.Classification, ZeroClassParam})
	case semantic.SideConditionKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.SideCondition, SideConditionParam}, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam})
		if p.SideCondition.Constant != nil && (p.SideCondition.Constant.SideCondition == semantic.PointInValidityDomain || p.SideCondition.Constant.SideCondition == semantic.CompletionFactorRegularNonzero) {
			required = append(required, struct {
				term Term
				kind ParamType
			}{p.Point, PointParam})
			if p.SideCondition.Constant.SideCondition == semantic.PointInValidityDomain {
				required = append(required, struct {
					term Term
					kind ParamType
				}{p.Domain, DomainParam})
			}
		} else {
			required = append(required, struct {
				term Term
				kind ParamType
			}{p.Domain, DomainParam})
		}
	case semantic.FunctionalIdentityKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Left, TransformParam}, struct {
			term Term
			kind ParamType
		}{p.Right, TransformParam}, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam})
	case semantic.ZeroSetPropertyKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Domain, DomainParam}, struct {
			term Term
			kind ParamType
		}{p.Property, ZeroSetPropertyParam})
		if p.Property.Constant != nil && p.Property.Constant.ZeroSetProperty == semantic.InvariantUnderTransform {
			required = append(required, struct {
				term Term
				kind ParamType
			}{p.Transform, TransformParam})
		} else {
			required = append(required, struct {
				term Term
				kind ParamType
			}{p.Region, DomainParam})
		}
	case semantic.ZeroClassificationKind:
		required = append(required, struct {
			term Term
			kind ParamType
		}{p.Object, ObjectParam}, struct {
			term Term
			kind ParamType
		}{p.Classification, ZeroClassificationParam})
	case semantic.UniversalFunctionalStatementKind:
		if p.FunctionClass == nil || p.Functional == "" || p.FunctionalPredicate != semantic.FunctionalNonnegative || p.TransformConvention == "" {
			return fmt.Errorf("incomplete functional statement pattern")
		}
		if err := p.FunctionClass.Validate(); err != nil {
			return err
		}
		if p.FunctionClass.TransformConvention != p.TransformConvention {
			return fmt.Errorf("functional pattern transform mismatch")
		}
	case semantic.FunctionalDefinitionKind,
		semantic.ExplicitFormulaIdentityKind,
		semantic.FiniteSpanDefinitionKind,
		semantic.QuadraticFormStructureKind,
		semantic.HermitianFormDefinitionKind,
		semantic.HermitianMatrixDefinitionKind,
		semantic.FiniteSpanFunctionalStatementKind,
		semantic.CoordinatePositivityKind,
		semantic.QuadraticMatrixIdentityKind,
		semantic.MatrixPropertyStatementKind,
		semantic.TwoByTwoMinorCertificateKind:
		// M5-M7 contracts are ground theorem schemas whose full typed payload lives
		// in the graph propositions. They are validated by those proposition
		// types and applied by the dedicated lowering passes, not instantiated
		// through the older M1-M4 parameter matcher.
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
		if t.Transform != semantic.IdentityTransform && want != PointParam {
			return fmt.Errorf("only point terms may be transformed")
		}
		if t.TransformParameter != "" && (want != PointParam || declared[t.TransformParameter] != TransformParam) {
			return fmt.Errorf("point transform parameter %s is not declared with transform type", t.TransformParameter)
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
	if t.TransformParameter != "" {
		tr, ok := e[t.TransformParameter]
		if !ok || tr.Type != TransformParam || value.Type != PointParam {
			return false
		}
		value.Point = value.Point.Apply(tr.Transform)
	}
	if t.Transform != semantic.IdentityTransform {
		if value.Type != PointParam {
			return false
		}
		// Every member of this Klein-four transform group is self-inverse.
		value.Point = value.Point.Apply(t.Transform)
	}
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
		ok := env.bind(p.Object, BindingValue{Type: ObjectParam, Object: r.Object}) && env.bind(p.Representation, BindingValue{Type: RepresentationParam, Representation: r.Name}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: r.ValidOn})
		if !ok {
			return false
		}
		if p.BaseObject.Parameter != "" || p.BaseObject.Constant != nil {
			return env.bind(p.BaseObject, BindingValue{Type: ObjectParam, Object: r.BaseObject})
		}
		return r.BaseObject == 0
	case semantic.RepresentationIdentity:
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Left, BindingValue{Type: RepresentationParam, Representation: v.Left}) && env.bind(p.Right, BindingValue{Type: RepresentationParam, Representation: v.Right}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.AnalyticFact:
		return env.bind(p.AnalyticFact, BindingValue{Type: AnalyticFactParam, AnalyticFact: v.Fact}) && env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.QuantifiedStatement:
		return env.bind(p.Quantifier, BindingValue{Type: QuantifierParam, Quantifier: v.Quantifier}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain}) && env.bind(p.Predicate, BindingValue{Type: PredicateParam, Predicate: v.Predicate.Kind}) && env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Predicate.Function})
	case semantic.ZeroAtPoint:
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Point, BindingValue{Type: PointParam, Point: v.Point.Canonical()}) && env.bind(p.Classification, BindingValue{Type: ZeroClassParam, ZeroClass: v.Classification})
	case semantic.SideCondition:
		if !env.bind(p.SideCondition, BindingValue{Type: SideConditionParam, SideCondition: v.Condition}) || !env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) {
			return false
		}
		if v.Condition == semantic.PointInValidityDomain || v.Condition == semantic.CompletionFactorRegularNonzero {
			ok := env.bind(p.Point, BindingValue{Type: PointParam, Point: v.Point.Canonical()})
			if v.Condition == semantic.PointInValidityDomain {
				ok = ok && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
			}
			return ok
		}
		return env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.FunctionalIdentity:
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Left, BindingValue{Type: TransformParam, Transform: v.Left}) && env.bind(p.Right, BindingValue{Type: TransformParam, Transform: v.Right}) && env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Domain})
	case semantic.ZeroSetProperty:
		if !env.bind(p.Domain, BindingValue{Type: DomainParam, Domain: v.Set}) || !env.bind(p.Property, BindingValue{Type: ZeroSetPropertyParam, ZeroSetProperty: v.Property}) {
			return false
		}
		if v.Property == semantic.InvariantUnderTransform {
			return env.bind(p.Transform, BindingValue{Type: TransformParam, Transform: v.Transform})
		}
		return env.bind(p.Region, BindingValue{Type: DomainParam, Domain: v.Region})
	case semantic.ZeroClassification:
		return env.bind(p.Object, BindingValue{Type: ObjectParam, Object: v.Object}) && env.bind(p.Classification, BindingValue{Type: ZeroClassificationParam, ZeroClassification: v.Classification})
	case semantic.UniversalFunctionalStatement:
		return p.FunctionClass != nil && v.Quantifier == semantic.ForAll && v.FunctionClass.Key() == p.FunctionClass.Key() && v.Functional == p.Functional && v.Predicate == p.FunctionalPredicate && v.TransformConvention == p.TransformConvention
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
	if t.Transform != semantic.IdentityTransform {
		if v.Type != PointParam {
			return BindingValue{}, fmt.Errorf("transform applied to non-point")
		}
		v.Point = v.Point.Apply(t.Transform)
	}
	if t.TransformParameter != "" {
		tr, ok := env[t.TransformParameter]
		if !ok {
			return BindingValue{}, fmt.Errorf("transform parameter %s is unbound", t.TransformParameter)
		}
		if v.Type != PointParam || tr.Type != TransformParam {
			return BindingValue{}, fmt.Errorf("dynamic transform applied to non-point")
		}
		v.Point = v.Point.Apply(tr.Transform)
	}
	return v, nil
}

func instantiate(p Pattern, env bindingEnv) (semantic.Proposition, error) {
	switch p.Kind {
	case semantic.RepresentationKind:
		object, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		domain, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
		r, err := resolve(p.Representation, env)
		if err != nil {
			return nil, err
		}
		base := semantic.Function(0)
		if p.BaseObject.Parameter != "" || p.BaseObject.Constant != nil {
			b, err := resolve(p.BaseObject, env)
			if err != nil {
				return nil, err
			}
			base = b.Object
		}
		return semantic.RepresentationProposition{Representation: semantic.Representation{Object: object.Object, BaseObject: base, Name: r.Representation, ValidOn: domain.Domain, Formula: p.Formula, Affordances: append([]string(nil), p.Affordances...)}}, nil
	case semantic.RepresentationIdentityKind:
		object, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		domain, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
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
		object, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		domain, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
		fact, err := resolve(p.AnalyticFact, env)
		if err != nil {
			return nil, err
		}
		return semantic.AnalyticFact{Fact: fact.AnalyticFact, Object: object.Object, Domain: domain.Domain}, nil
	case semantic.QuantifiedStatementKind:
		object, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		domain, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
		q, err := resolve(p.Quantifier, env)
		if err != nil {
			return nil, err
		}
		pred, err := resolve(p.Predicate, env)
		if err != nil {
			return nil, err
		}
		return semantic.QuantifiedStatement{Quantifier: q.Quantifier, Domain: domain.Domain, Predicate: semantic.Predicate{Kind: pred.Predicate, Function: object.Object}}, nil
	case semantic.ZeroAtPointKind:
		o, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		pt, err := resolve(p.Point, env)
		if err != nil {
			return nil, err
		}
		z, err := resolve(p.Classification, env)
		if err != nil {
			return nil, err
		}
		return semantic.ZeroAtPoint{Object: o.Object, Point: pt.Point.Canonical(), Classification: z.ZeroClass}, nil
	case semantic.SideConditionKind:
		o, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		s, err := resolve(p.SideCondition, env)
		if err != nil {
			return nil, err
		}
		out := semantic.SideCondition{Condition: s.SideCondition, Object: o.Object}
		if s.SideCondition == semantic.PointInValidityDomain || s.SideCondition == semantic.CompletionFactorRegularNonzero {
			pt, err := resolve(p.Point, env)
			if err != nil {
				return nil, err
			}
			out.Point = pt.Point.Canonical()
			if s.SideCondition == semantic.PointInValidityDomain {
				d, err := resolve(p.Domain, env)
				if err != nil {
					return nil, err
				}
				out.Domain = d.Domain
			}
		} else {
			d, err := resolve(p.Domain, env)
			if err != nil {
				return nil, err
			}
			out.Domain = d.Domain
		}
		return out, nil
	case semantic.FunctionalIdentityKind:
		o, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		l, err := resolve(p.Left, env)
		if err != nil {
			return nil, err
		}
		r, err := resolve(p.Right, env)
		if err != nil {
			return nil, err
		}
		d, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
		return semantic.FunctionalIdentity{Object: o.Object, Left: l.Transform, Right: r.Transform, Domain: d.Domain, Formula: p.Formula}, nil
	case semantic.ZeroSetPropertyKind:
		d, err := resolve(p.Domain, env)
		if err != nil {
			return nil, err
		}
		pr, err := resolve(p.Property, env)
		if err != nil {
			return nil, err
		}
		out := semantic.ZeroSetProperty{Set: d.Domain, Property: pr.ZeroSetProperty}
		if out.Property == semantic.InvariantUnderTransform {
			tr, err := resolve(p.Transform, env)
			if err != nil {
				return nil, err
			}
			out.Transform = tr.Transform
		} else {
			r, err := resolve(p.Region, env)
			if err != nil {
				return nil, err
			}
			out.Region = r.Domain
		}
		return out, nil
	case semantic.ZeroClassificationKind:
		o, err := resolve(p.Object, env)
		if err != nil {
			return nil, err
		}
		c, err := resolve(p.Classification, env)
		if err != nil {
			return nil, err
		}
		return semantic.ZeroClassification{Object: o.Object, Classification: c.ZeroClassification}, nil
	case semantic.UniversalFunctionalStatementKind:
		if p.FunctionClass == nil {
			return nil, fmt.Errorf("functional statement has no class")
		}
		return semantic.UniversalFunctionalStatement{Quantifier: semantic.ForAll, Variable: "f", FunctionClass: *p.FunctionClass, Functional: p.Functional, Predicate: p.FunctionalPredicate, TransformConvention: p.TransformConvention}, nil
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
	patterns := append(append([]Pattern(nil), c.Premises...), c.SideConditions...)
	for premiseIndex, pattern := range patterns {
		var next []matchState
		for _, state := range states {
			for _, claim := range claims {
				env := state.env.clone()
				if matchPattern(pattern, claim, env) {
					matched := append(append([]PremiseMatch(nil), state.matched...), PremiseMatch{Premise: premiseIndex, Claim: claim.ID, SideCondition: premiseIndex >= len(c.Premises)})
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
		pass := "instantiate-theorem-contract"
		for _, term := range []Term{c.Conclusion.Point} {
			if term.Transform != semantic.IdentityTransform || term.TransformParameter != "" {
				pass = "transport-theorem-contract"
			}
		}
		t := Transformation{ID: transformID, Pass: pass, From: parents[0], To: claimID, Relation: c.Relation, Provenance: c.Evidence.Source, Theorem: c.ID, Bindings: bindings, Trusted: c.Trust != UntrustedTheorem}
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
		patterns := append(append([]Pattern(nil), c.Premises...), c.SideConditions...)
		if len(patterns) == 0 || len(e.fullMatches(c)) > 0 {
			continue
		}
		state := matchState{env: make(bindingEnv)}
		used := make(map[semantic.ClaimID]bool)
		for i, p := range patterns {
			for _, claim := range e.sortedClaims() {
				if used[claim.ID] {
					continue
				}
				env := state.env.clone()
				if matchPattern(p, claim, env) {
					state.env = env
					state.matched = append(state.matched, PremiseMatch{Premise: i, Claim: claim.ID, SideCondition: i >= len(c.Premises)})
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
		for i, p := range patterns {
			if matchedIndex[i] {
				continue
			}
			o := TheoremObligation{Premise: i, Pattern: p, SideCondition: i >= len(c.Premises)}
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
