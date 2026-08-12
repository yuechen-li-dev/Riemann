package semantic

import (
	"fmt"
	"sort"
	"strings"
)

// M4 deliberately models one transform convention, not Fourier analysis in
// general.  The fields make the convention independently inspectable.
type TransformConventionID string

const LagariasMellinConvention TransformConventionID = "lagarias_mellin_dx_over_x"

type TransformConvention struct {
	ID             TransformConventionID `json:"id"`
	Kind           string                `json:"kind"`
	SourceVariable string                `json:"source_variable"`
	TargetVariable string                `json:"target_variable"`
	Kernel         string                `json:"kernel"`
	Measure        string                `json:"measure"`
	Normalization  string                `json:"normalization"`
}

func LagariasMellinTransform() TransformConvention {
	return TransformConvention{ID: LagariasMellinConvention, Kind: "mellin", SourceVariable: "x in (0,infinity)", TargetVariable: "s in C", Kernel: "x^s", Measure: "dx/x", Normalization: "M[f](s)=integral_0^infinity f(x)x^s dx/x; no additional constant"}
}
func (t TransformConvention) Validate() error {
	if t.ID == "" || t.Kind == "" || t.SourceVariable == "" || t.TargetVariable == "" || t.Kernel == "" || t.Measure == "" || t.Normalization == "" {
		return fmt.Errorf("incomplete transform convention")
	}
	return nil
}

type FunctionConstraint string

const (
	ComplexValued           FunctionConstraint = "complex_valued"
	PiecewiseC2             FunctionConstraint = "piecewise_C2"
	CompactSupportPositive  FunctionConstraint = "compact_support_in_positive_reals"
	MidpointAtDiscontinuity FunctionConstraint = "midpoint_value_at_discontinuities"
)

type TestFunctionKind string

const (
	ImportedTestFunction TestFunctionKind = "imported_function"
	BasisTestFunction    TestFunctionKind = "basis_function"
)

type FunctionParameter struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// TestFunction identity includes its constructor, parameters, declared class,
// attributes, and transform.  Symbol is display identity, never the whole IR.
type TestFunction struct {
	Symbol              string                `json:"symbol"`
	Kind                TestFunctionKind      `json:"kind"`
	DeclaredClass       string                `json:"declared_class"`
	RequiredAttributes  []FunctionConstraint  `json:"required_attributes"`
	TransformConvention TransformConventionID `json:"transform_convention"`
	Parameters          []FunctionParameter   `json:"parameters"`
}

func (f TestFunction) Validate() error {
	if strings.TrimSpace(f.Symbol) == "" || (f.Kind != ImportedTestFunction && f.Kind != BasisTestFunction) || f.DeclaredClass == "" || f.TransformConvention == "" {
		return fmt.Errorf("invalid test-function identity")
	}
	seen := map[string]bool{}
	for _, p := range f.Parameters {
		if p.Name == "" || p.Value == "" || seen[p.Name] {
			return fmt.Errorf("invalid or duplicate function parameter")
		}
		seen[p.Name] = true
	}
	return nil
}
func (f TestFunction) Key() string {
	attrs := append([]FunctionConstraint(nil), f.RequiredAttributes...)
	sort.Slice(attrs, func(i, j int) bool { return attrs[i] < attrs[j] })
	params := append([]FunctionParameter(nil), f.Parameters...)
	sort.Slice(params, func(i, j int) bool {
		return params[i].Name < params[j].Name || (params[i].Name == params[j].Name && params[i].Value < params[j].Value)
	})
	var b strings.Builder
	fmt.Fprintf(&b, "%s|%s|%s|%s", f.Symbol, f.Kind, f.DeclaredClass, f.TransformConvention)
	for _, a := range attrs {
		fmt.Fprintf(&b, "|a:%s", a)
	}
	for _, p := range params {
		fmt.Fprintf(&b, "|p:%s=%s", p.Name, p.Value)
	}
	return b.String()
}

type FunctionClassKind string

const (
	WeilNiceFunctionClass FunctionClassKind = "weil_nice"
	FiniteFunctionClass   FunctionClassKind = "finite_family"
)

type FunctionClass struct {
	ID                  string                `json:"id"`
	Kind                FunctionClassKind     `json:"kind"`
	Constraints         []FunctionConstraint  `json:"constraints"`
	TransformConvention TransformConventionID `json:"transform_convention"`
	Members             []TestFunction        `json:"members"`
}

func WeilNiceClass() FunctionClass {
	return FunctionClass{ID: "weil-nice-zeta", Kind: WeilNiceFunctionClass, Constraints: []FunctionConstraint{ComplexValued, PiecewiseC2, CompactSupportPositive, MidpointAtDiscontinuity}, TransformConvention: LagariasMellinConvention}
}
func (c FunctionClass) Validate() error {
	if c.ID == "" || c.TransformConvention == "" || (c.Kind != WeilNiceFunctionClass && c.Kind != FiniteFunctionClass) {
		return fmt.Errorf("invalid function class")
	}
	if c.Kind == WeilNiceFunctionClass && len(c.Members) != 0 {
		return fmt.Errorf("full function class cannot enumerate members")
	}
	if c.Kind == FiniteFunctionClass && len(c.Members) == 0 {
		return fmt.Errorf("finite function family is empty")
	}
	seen := map[string]bool{}
	for _, f := range c.Members {
		if err := f.Validate(); err != nil {
			return err
		}
		if seen[f.Key()] {
			return fmt.Errorf("duplicate finite-family member")
		}
		if f.TransformConvention != c.TransformConvention {
			return fmt.Errorf("finite-family member uses an incompatible transform")
		}
		seen[f.Key()] = true
	}
	return nil
}
func (c FunctionClass) Key() string {
	constraints := append([]FunctionConstraint(nil), c.Constraints...)
	sort.Slice(constraints, func(i, j int) bool { return constraints[i] < constraints[j] })
	members := make([]string, 0, len(c.Members))
	for _, f := range c.Members {
		members = append(members, f.Key())
	}
	sort.Strings(members)
	return fmt.Sprintf("%s|%s|%s|%s|%s", c.ID, c.Kind, c.TransformConvention, joinConstraints(constraints), strings.Join(members, ";"))
}
func (c FunctionClass) Describe() string {
	if c.Kind == FiniteFunctionClass {
		names := make([]string, 0, len(c.Members))
		for _, f := range c.Members {
			names = append(names, f.Symbol)
		}
		sort.Strings(names)
		return fmt.Sprintf("finite family {%s}", strings.Join(names, ", "))
	}
	return "Weil-nice test functions (piecewise C2, compactly supported in (0,infinity), with midpoint values at discontinuities)"
}
func joinConstraints(in []FunctionConstraint) string {
	out := make([]string, len(in))
	for i, v := range in {
		out[i] = string(v)
	}
	return strings.Join(out, ",")
}

func CloneTestFunction(f TestFunction) TestFunction {
	f.RequiredAttributes = append([]FunctionConstraint(nil), f.RequiredAttributes...)
	f.Parameters = append([]FunctionParameter(nil), f.Parameters...)
	return f
}
func CloneFunctionClass(c FunctionClass) FunctionClass {
	c.Constraints = append([]FunctionConstraint(nil), c.Constraints...)
	c.Members = append([]TestFunction(nil), c.Members...)
	for i := range c.Members {
		c.Members[i] = CloneTestFunction(c.Members[i])
	}
	return c
}
func CloneQuadraticFunctional(q QuadraticFunctional) QuadraticFunctional {
	q.Contributions = append([]FunctionalContribution(nil), q.Contributions...)
	for i := range q.Contributions {
		if q.Contributions[i].Aggregate != nil {
			a := CloneAggregate(*q.Contributions[i].Aggregate)
			q.Contributions[i].Aggregate = &a
		}
	}
	return q
}

func CloneAggregate(a Aggregate) Aggregate {
	a.Convergence = append([]string(nil), a.Convergence...)
	a.TheoremLineage = append([]TheoremID(nil), a.TheoremLineage...)
	return a
}

type FunctionalID string

const WeilZetaQuadraticFunctional FunctionalID = "weil_zeta_quadratic"

type FunctionalPredicateKind string

const FunctionalNonnegative FunctionalPredicateKind = "nonnegative"

type AggregateKind string

const SumOverDomain AggregateKind = "sum_over_domain"

type Aggregate struct {
	Kind                AggregateKind         `json:"kind"`
	IndexDomain         Domain                `json:"index_domain"`
	Summand             string                `json:"summand"`
	Convergence         []string              `json:"convergence_conditions"`
	TransformConvention TransformConventionID `json:"transform_convention"`
	TheoremLineage      []TheoremID           `json:"theorem_lineage"`
	Provenance          Reference             `json:"provenance"`
}

func (a Aggregate) Validate() error {
	if a.Kind != SumOverDomain || a.Summand == "" || a.TransformConvention == "" || len(a.Convergence) == 0 || len(a.TheoremLineage) == 0 || a.Provenance.Citation == "" {
		return fmt.Errorf("invalid aggregate semantics or provenance")
	}
	return a.IndexDomain.Validate()
}

type FunctionalContributionKind string

const (
	ZeroContribution        FunctionalContributionKind = "nontrivial_zero_aggregate"
	EndpointContribution    FunctionalContributionKind = "endpoint_normalization"
	PrimePowerContribution  FunctionalContributionKind = "prime_power"
	ArchimedeanContribution FunctionalContributionKind = "archimedean"
)

type FunctionalContribution struct {
	Kind               FunctionalContributionKind `json:"kind"`
	RepresentationSide string                     `json:"representation_side"`
	Sign               int                        `json:"sign"`
	Formula            string                     `json:"formula"`
	Aggregate          *Aggregate                 `json:"aggregate,omitempty"`
	Index              string                     `json:"index,omitempty"`
	Weight             string                     `json:"weight,omitempty"`
}
type QuadraticFunctional struct {
	ID                  FunctionalID             `json:"id"`
	Object              Function                 `json:"object"`
	InputConstruction   string                   `json:"input_construction"`
	TransformConvention TransformConventionID    `json:"transform_convention"`
	Contributions       []FunctionalContribution `json:"contributions"`
}

func (q QuadraticFunctional) Validate() error {
	if q.ID == "" || q.Object == 0 || q.InputConstruction == "" || q.TransformConvention == "" || len(q.Contributions) == 0 {
		return fmt.Errorf("invalid quadratic functional")
	}
	for _, c := range q.Contributions {
		if c.RepresentationSide == "" {
			return fmt.Errorf("functional contribution has no representation side")
		}
		if c.Sign != -1 && c.Sign != 1 {
			return fmt.Errorf("functional contribution has invalid sign")
		}
		if c.Kind == ZeroContribution && (c.Aggregate == nil || c.Aggregate.IndexDomain != NontrivialZeros(q.Object)) {
			return fmt.Errorf("zero contribution has wrong aggregate domain")
		}
		if c.Aggregate != nil {
			if err := c.Aggregate.Validate(); err != nil {
				return err
			}
		}
	}
	return nil
}

type FunctionalDefinition struct {
	Functional QuadraticFunctional `json:"functional"`
}

func (FunctionalDefinition) isProposition()        {}
func (FunctionalDefinition) Kind() PropositionKind { return FunctionalDefinitionKind }
func (p FunctionalDefinition) Describe() string {
	return fmt.Sprintf("%s is a quadratic functional with %d inspectable contributions", p.Functional.ID, len(p.Functional.Contributions))
}

type UniversalFunctionalStatement struct {
	Quantifier          QuantifierKind          `json:"quantifier"`
	Variable            string                  `json:"variable"`
	FunctionClass       FunctionClass           `json:"function_class"`
	Functional          FunctionalID            `json:"functional"`
	Predicate           FunctionalPredicateKind `json:"predicate"`
	TransformConvention TransformConventionID   `json:"transform_convention"`
}

func (UniversalFunctionalStatement) isProposition()        {}
func (UniversalFunctionalStatement) Kind() PropositionKind { return UniversalFunctionalStatementKind }
func (p UniversalFunctionalStatement) Validate() error {
	if p.Quantifier != ForAll || p.Variable == "" || p.Functional == "" || p.Predicate != FunctionalNonnegative {
		return fmt.Errorf("invalid universal functional statement")
	}
	if err := p.FunctionClass.Validate(); err != nil {
		return err
	}
	if p.TransformConvention != p.FunctionClass.TransformConvention {
		return fmt.Errorf("functional statement transform does not match class")
	}
	return nil
}
func (p UniversalFunctionalStatement) Describe() string {
	return fmt.Sprintf("for every %s in %s: %s(%s) >= 0", p.Variable, p.FunctionClass.Describe(), p.Functional, p.Variable)
}

type TestFunctionAdmissibility struct {
	Function TestFunction  `json:"function"`
	Class    FunctionClass `json:"class"`
}

func (TestFunctionAdmissibility) isProposition()        {}
func (TestFunctionAdmissibility) Kind() PropositionKind { return TestFunctionAdmissibilityKind }
func (p TestFunctionAdmissibility) Validate() error {
	if err := p.Function.Validate(); err != nil {
		return err
	}
	if err := p.Class.Validate(); err != nil {
		return err
	}
	if p.Class.Kind != WeilNiceFunctionClass || p.Function.DeclaredClass != p.Class.ID || p.Function.TransformConvention != p.Class.TransformConvention || !containsAllConstraints(p.Function.RequiredAttributes, p.Class.Constraints) {
		return fmt.Errorf("admissibility class mismatch")
	}
	return nil
}
func (p TestFunctionAdmissibility) Describe() string {
	return fmt.Sprintf("%s is admissible in %s", p.Function.Symbol, p.Class.Describe())
}

type ExplicitFormulaIdentity struct {
	Functional          FunctionalID             `json:"functional"`
	FunctionClass       FunctionClass            `json:"function_class"`
	TransformConvention TransformConvention      `json:"transform_convention"`
	ZeroSide            Aggregate                `json:"zero_side"`
	ArithmeticSide      []FunctionalContribution `json:"arithmetic_side"`
	Theorem             TheoremID                `json:"theorem"`
}

func (ExplicitFormulaIdentity) isProposition()        {}
func (ExplicitFormulaIdentity) Kind() PropositionKind { return ExplicitFormulaIdentityKind }
func (p ExplicitFormulaIdentity) Validate() error {
	if p.Functional == "" || p.Theorem == "" || len(p.ArithmeticSide) == 0 {
		return fmt.Errorf("invalid explicit-formula identity")
	}
	if err := p.ZeroSide.Validate(); err != nil {
		return err
	}
	if err := p.FunctionClass.Validate(); err != nil {
		return err
	}
	if err := p.TransformConvention.Validate(); err != nil {
		return err
	}
	if p.TransformConvention.ID != p.FunctionClass.TransformConvention {
		return fmt.Errorf("explicit formula transform mismatch")
	}
	return nil
}

func containsAllConstraints(have, need []FunctionConstraint) bool {
	set := make(map[FunctionConstraint]bool, len(have))
	for _, v := range have {
		set[v] = true
	}
	for _, v := range need {
		if !set[v] {
			return false
		}
	}
	return true
}
func (p ExplicitFormulaIdentity) Describe() string {
	return fmt.Sprintf("the explicit formula identifies the zero-side representation of %s with its endpoint, prime-power, and archimedean representation", p.Functional)
}
