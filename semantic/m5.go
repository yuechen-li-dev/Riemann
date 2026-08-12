package semantic

import (
	"fmt"
	"sort"
	"strings"
)

// M5 is intentionally a finite-linear-algebra IR, not a general expression
// language. Scalar fields and coefficient vectors are therefore closed and
// purpose-specific.
type ScalarField string

const (
	RealField    ScalarField = "real"
	ComplexField ScalarField = "complex"
)

func (f ScalarField) Valid() bool { return f == RealField || f == ComplexField }

type BasisMember struct {
	Function                 TestFunction `json:"function"`
	AdmissibilityCertificate ClaimID      `json:"admissibility_certificate"`
}

// OrderedBasis owns coordinate order. Key is order-sensitive; the containing
// FiniteSpan has a separate, order-insensitive mathematical identity.
type OrderedBasis struct {
	ID               string        `json:"id"`
	Members          []BasisMember `json:"members"`
	CoefficientField ScalarField   `json:"coefficient_field"`
	ParentClass      FunctionClass `json:"parent_function_class"`
}

func (b OrderedBasis) Validate() error {
	if strings.TrimSpace(b.ID) == "" || !b.CoefficientField.Valid() || len(b.Members) == 0 {
		return fmt.Errorf("invalid ordered basis")
	}
	if err := b.ParentClass.Validate(); err != nil || b.ParentClass.Kind != WeilNiceFunctionClass {
		return fmt.Errorf("ordered basis has invalid parent function class")
	}
	seen := map[string]bool{}
	for _, member := range b.Members {
		if err := member.Function.Validate(); err != nil {
			return err
		}
		if member.AdmissibilityCertificate == "" {
			return fmt.Errorf("basis member %s has no admissibility certificate", member.Function.Symbol)
		}
		if member.Function.DeclaredClass != b.ParentClass.ID || member.Function.TransformConvention != b.ParentClass.TransformConvention || !containsAllConstraints(member.Function.RequiredAttributes, b.ParentClass.Constraints) {
			return fmt.Errorf("basis member %s is incompatible with parent class", member.Function.Symbol)
		}
		if seen[member.Function.Key()] {
			return fmt.Errorf("duplicate ordered-basis member")
		}
		seen[member.Function.Key()] = true
	}
	return nil
}

func (b OrderedBasis) Key() string {
	parts := make([]string, len(b.Members))
	for i, member := range b.Members {
		parts[i] = fmt.Sprintf("%s@%s", member.Function.Key(), member.AdmissibilityCertificate)
	}
	return fmt.Sprintf("%s|%s|%s|%s", b.ID, b.CoefficientField, b.ParentClass.Key(), strings.Join(parts, ";"))
}

type FiniteSpan struct {
	ID               string        `json:"id"`
	Basis            OrderedBasis  `json:"ordered_basis"`
	CoefficientField ScalarField   `json:"coefficient_field"`
	ParentClass      FunctionClass `json:"parent_function_class"`
}

func (s FiniteSpan) Validate() error {
	if strings.TrimSpace(s.ID) == "" || !s.CoefficientField.Valid() {
		return fmt.Errorf("invalid finite span")
	}
	if err := s.Basis.Validate(); err != nil {
		return err
	}
	if s.CoefficientField != s.Basis.CoefficientField || s.ParentClass.Key() != s.Basis.ParentClass.Key() {
		return fmt.Errorf("span field or parent class disagrees with basis")
	}
	return nil
}

// Key is invariant under a basis permutation: it identifies the mathematical
// span, while OrderedBasis.Key identifies a coordinate system on that span.
func (s FiniteSpan) Key() string {
	members := make([]string, len(s.Basis.Members))
	for i, member := range s.Basis.Members {
		members[i] = member.Function.Key()
	}
	sort.Strings(members)
	return fmt.Sprintf("%s|%s|%s|%s", s.ID, s.CoefficientField, s.ParentClass.Key(), strings.Join(members, ";"))
}

func (s FiniteSpan) Describe() string {
	names := make([]string, len(s.Basis.Members))
	for i, member := range s.Basis.Members {
		names[i] = member.Function.Symbol
	}
	return fmt.Sprintf("%s-span span(%s) over %s", s.ID, strings.Join(names, ", "), s.CoefficientField)
}

type FiniteSpanDefinition struct {
	Span FiniteSpan `json:"span"`
}

func (FiniteSpanDefinition) isProposition()        {}
func (FiniteSpanDefinition) Kind() PropositionKind { return FiniteSpanDefinitionKind }
func (p FiniteSpanDefinition) Describe() string {
	return p.Span.Describe() + " is a typed finite admissible span"
}

type Coefficient struct {
	Symbol string `json:"symbol"`
}

type CoefficientVector struct {
	ID      string        `json:"id"`
	Field   ScalarField   `json:"field"`
	Entries []Coefficient `json:"entries"`
}

func (v CoefficientVector) Validate(dimension int, field ScalarField) error {
	if v.ID == "" || v.Field != field || len(v.Entries) != dimension {
		return fmt.Errorf("coefficient vector does not match span coordinates")
	}
	seen := map[string]bool{}
	for _, entry := range v.Entries {
		if entry.Symbol == "" || seen[entry.Symbol] {
			return fmt.Errorf("invalid coefficient vector entry")
		}
		seen[entry.Symbol] = true
	}
	return nil
}

type FiniteLinearCombination struct {
	ID           string            `json:"id"`
	Span         FiniteSpan        `json:"span"`
	Coefficients CoefficientVector `json:"coefficients"`
}

func (c FiniteLinearCombination) Validate() error {
	if c.ID == "" {
		return fmt.Errorf("linear combination has no identity")
	}
	if err := c.Span.Validate(); err != nil {
		return err
	}
	return c.Coefficients.Validate(len(c.Span.Basis.Members), c.Span.CoefficientField)
}

func (c FiniteLinearCombination) Describe() string {
	terms := make([]string, len(c.Span.Basis.Members))
	for i, member := range c.Span.Basis.Members {
		terms[i] = c.Coefficients.Entries[i].Symbol + " " + member.Function.Symbol
	}
	return strings.Join(terms, " + ")
}

type QuadraticLaw string

const (
	AbsoluteSquareHomogeneity QuadraticLaw = "Q(lambda*f)=|lambda|^2 Q(f)"
	ParallelogramLaw          QuadraticLaw = "Q(f+g)+Q(f-g)=2Q(f)+2Q(g)"
	RealValuedDiagonal        QuadraticLaw = "Q(f) is real-valued"
)

type QuadraticFormStructure struct {
	Functional       FunctionalID   `json:"functional"`
	DomainSpan       FiniteSpan     `json:"domain_span"`
	CoefficientField ScalarField    `json:"coefficient_field"`
	Laws             []QuadraticLaw `json:"laws"`
	Theorem          TheoremID      `json:"theorem"`
}

func (QuadraticFormStructure) isProposition()        {}
func (QuadraticFormStructure) Kind() PropositionKind { return QuadraticFormStructureKind }
func (p QuadraticFormStructure) Validate() error {
	if p.Functional == "" || p.CoefficientField != ComplexField || p.Theorem == "" {
		return fmt.Errorf("invalid complex quadratic-form structure")
	}
	if err := p.DomainSpan.Validate(); err != nil {
		return err
	}
	want := map[QuadraticLaw]bool{AbsoluteSquareHomogeneity: false, ParallelogramLaw: false, RealValuedDiagonal: false}
	for _, law := range p.Laws {
		if _, ok := want[law]; !ok || want[law] {
			return fmt.Errorf("unknown or duplicate quadratic law")
		}
		want[law] = true
	}
	for _, present := range want {
		if !present {
			return fmt.Errorf("quadratic-form prerequisite is missing")
		}
	}
	return nil
}
func (p QuadraticFormStructure) Describe() string {
	return fmt.Sprintf("%s satisfies the complex quadratic-form laws on %s", p.Functional, p.DomainSpan.Describe())
}
func (p QuadraticFormStructure) Key() string {
	return fmt.Sprintf("quadratic-structure|%s|%s|%s|%s", p.Functional, p.DomainSpan.Key(), p.CoefficientField, p.Theorem)
}

type FormArgument int

const (
	FirstArgument  FormArgument = 1
	SecondArgument FormArgument = 2
)

type PolarizationConvention struct {
	ID                      string       `json:"id"`
	CoefficientField        ScalarField  `json:"coefficient_field"`
	ConjugateLinearArgument FormArgument `json:"conjugate_linear_argument"`
	LinearArgument          FormArgument `json:"linear_argument"`
	Normalization           string       `json:"normalization"`
	Formula                 string       `json:"formula"`
	Theorem                 TheoremID    `json:"theorem"`
}

func ComplexConjugateFirstPolarization(theorem TheoremID) PolarizationConvention {
	return PolarizationConvention{ID: "complex-polarization-conjugate-first", CoefficientField: ComplexField, ConjugateLinearArgument: FirstArgument, LinearArgument: SecondArgument, Normalization: "1/4", Formula: "B(f,g)=1/4*(Q(f+g)-Q(f-g)-i*Q(f+i*g)+i*Q(f-i*g))", Theorem: theorem}
}

func (p PolarizationConvention) Validate() error {
	if p.ID == "" || p.CoefficientField != ComplexField || p.ConjugateLinearArgument != FirstArgument || p.LinearArgument != SecondArgument || p.Normalization != "1/4" || p.Formula == "" || p.Theorem == "" {
		return fmt.Errorf("invalid M5 polarization convention")
	}
	return nil
}

type FormContribution struct {
	SourceKind         FunctionalContributionKind `json:"source_kind"`
	RepresentationSide string                     `json:"representation_side"`
	Sign               int                        `json:"sign"`
	EntryDefinition    string                     `json:"entry_definition"`
}

type HermitianForm struct {
	ID                string                 `json:"id"`
	SourceFunctional  FunctionalID           `json:"source_functional"`
	DomainSpan        FiniteSpan             `json:"domain_span"`
	Convention        PolarizationConvention `json:"polarization_convention"`
	EntryDefinition   string                 `json:"entry_definition"`
	Contributions     []FormContribution     `json:"contributions"`
	RecoveryIdentity  string                 `json:"recovery_identity"`
	HermitianIdentity string                 `json:"hermitian_identity"`
	TheoremProvenance []TheoremID            `json:"theorem_provenance"`
}

func (h HermitianForm) Validate() error {
	if h.ID == "" || h.SourceFunctional == "" || h.EntryDefinition == "" || h.RecoveryIdentity != "Q(f)=B(f,f)" || h.HermitianIdentity != "B(f,g)=conjugate(B(g,f))" || len(h.TheoremProvenance) == 0 {
		return fmt.Errorf("invalid Hermitian form")
	}
	if err := h.DomainSpan.Validate(); err != nil {
		return err
	}
	if err := h.Convention.Validate(); err != nil {
		return err
	}
	return nil
}
func (h HermitianForm) Key() string {
	return fmt.Sprintf("%s|%s|%s|%s", h.ID, h.SourceFunctional, h.DomainSpan.Key(), h.Convention.ID)
}

type HermitianFormDefinition struct {
	Form HermitianForm `json:"form"`
}

func (HermitianFormDefinition) isProposition()        {}
func (HermitianFormDefinition) Kind() PropositionKind { return HermitianFormDefinitionKind }
func (p HermitianFormDefinition) Describe() string {
	return fmt.Sprintf("%s is the conjugate-first Hermitian polarization of %s", p.Form.ID, p.Form.SourceFunctional)
}

type MatrixValueSemantics string

const (
	StructuralExactMatrix        MatrixValueSemantics = "structurally_defined_exact"
	NumericallyApproximateMatrix MatrixValueSemantics = "numerically_evaluated_approximate"
)

type EntryContribution struct {
	SourceKind          FunctionalContributionKind `json:"source_kind"`
	RepresentationSide  string                     `json:"representation_side"`
	Sign                int                        `json:"sign"`
	PolarizedDefinition string                     `json:"polarized_definition"`
}

type MatrixEntry struct {
	Row                 int                   `json:"row"`
	Column              int                   `json:"column"`
	RowFunction         TestFunction          `json:"row_function"`
	ColumnFunction      TestFunction          `json:"column_function"`
	SourceForm          string                `json:"source_form"`
	SourceFunctional    FunctionalID          `json:"source_functional"`
	Definition          string                `json:"definition"`
	Contributions       []EntryContribution   `json:"contributions"`
	TransformConvention TransformConventionID `json:"transform_convention"`
	TheoremProvenance   []TheoremID           `json:"theorem_provenance"`
}

type HermitianMatrix struct {
	ID                      string               `json:"id"`
	SourceForm              HermitianForm        `json:"source_form"`
	Basis                   OrderedBasis         `json:"ordered_basis"`
	Rows                    int                  `json:"rows"`
	Columns                 int                  `json:"columns"`
	Entries                 []MatrixEntry        `json:"entries"`
	ValueSemantics          MatrixValueSemantics `json:"value_semantics"`
	HermitianByConstruction bool                 `json:"hermitian_by_construction"`
	LoweringTheorem         TheoremID            `json:"lowering_theorem"`
}

func (m HermitianMatrix) Validate() error {
	if m.ID == "" || m.Rows == 0 || m.Rows != m.Columns || len(m.Entries) != m.Rows*m.Columns || m.ValueSemantics != StructuralExactMatrix || !m.HermitianByConstruction || m.LoweringTheorem == "" {
		return fmt.Errorf("invalid structural Hermitian matrix")
	}
	if err := m.SourceForm.Validate(); err != nil {
		return err
	}
	if err := m.Basis.Validate(); err != nil {
		return err
	}
	if m.SourceForm.DomainSpan.Key() != (FiniteSpan{ID: m.SourceForm.DomainSpan.ID, Basis: m.Basis, CoefficientField: m.Basis.CoefficientField, ParentClass: m.Basis.ParentClass}).Key() {
		return fmt.Errorf("matrix basis does not coordinate source-form span")
	}
	for k, entry := range m.Entries {
		i, j := k/m.Columns, k%m.Columns
		if entry.Row != i || entry.Column != j || entry.RowFunction.Key() != m.Basis.Members[i].Function.Key() || entry.ColumnFunction.Key() != m.Basis.Members[j].Function.Key() || entry.SourceForm != m.SourceForm.ID || entry.SourceFunctional != m.SourceForm.SourceFunctional || entry.Definition == "" || len(entry.TheoremProvenance) == 0 {
			return fmt.Errorf("matrix entry %d has lost coordinate or provenance identity", k)
		}
	}
	return nil
}
func (m HermitianMatrix) Key() string {
	return fmt.Sprintf("%s|%s|%s|%s", m.ID, m.SourceForm.Key(), m.Basis.Key(), m.ValueSemantics)
}

type HermitianMatrixDefinition struct {
	Matrix HermitianMatrix `json:"matrix"`
}

func (HermitianMatrixDefinition) isProposition()        {}
func (HermitianMatrixDefinition) Kind() PropositionKind { return HermitianMatrixDefinitionKind }
func (p HermitianMatrixDefinition) Describe() string {
	return fmt.Sprintf("%s has entries G_ij=%s(f_i,f_j) in ordered basis %s", p.Matrix.ID, p.Matrix.SourceForm.ID, p.Matrix.Basis.ID)
}

type FiniteSpanFunctionalStatement struct {
	Quantifier          QuantifierKind          `json:"quantifier"`
	Variable            string                  `json:"variable"`
	Span                FiniteSpan              `json:"span"`
	Functional          FunctionalID            `json:"functional"`
	Predicate           FunctionalPredicateKind `json:"predicate"`
	TransformConvention TransformConventionID   `json:"transform_convention"`
}

func (FiniteSpanFunctionalStatement) isProposition()        {}
func (FiniteSpanFunctionalStatement) Kind() PropositionKind { return FiniteSpanFunctionalStatementKind }
func (p FiniteSpanFunctionalStatement) Validate() error {
	if p.Quantifier != ForAll || p.Variable == "" || p.Functional == "" || p.Predicate != FunctionalNonnegative || p.TransformConvention != p.Span.ParentClass.TransformConvention {
		return fmt.Errorf("invalid finite-span functional statement")
	}
	return p.Span.Validate()
}
func (p FiniteSpanFunctionalStatement) Describe() string {
	return fmt.Sprintf("for every %s in %s: %s(%s) >= 0", p.Variable, p.Span.Describe(), p.Functional, p.Variable)
}
func (p FiniteSpanFunctionalStatement) Key() string {
	return fmt.Sprintf("span-functional|%s|%s|%s|%s", p.Span.Key(), p.Functional, p.Predicate, p.TransformConvention)
}

type CoordinateQuadraticPositivity struct {
	MatrixID         string      `json:"matrix_id"`
	Span             FiniteSpan  `json:"span"`
	CoefficientField ScalarField `json:"coefficient_field"`
	Expression       string      `json:"expression"`
}

func (CoordinateQuadraticPositivity) isProposition()        {}
func (CoordinateQuadraticPositivity) Kind() PropositionKind { return CoordinatePositivityKind }
func (p CoordinateQuadraticPositivity) Validate() error {
	if p.MatrixID == "" || p.CoefficientField != p.Span.CoefficientField || p.Expression != "c* G c >= 0 for every coefficient vector c" {
		return fmt.Errorf("invalid coordinate positivity statement")
	}
	return p.Span.Validate()
}
func (p CoordinateQuadraticPositivity) Describe() string {
	return fmt.Sprintf("for every c in %s^%d: c* %s c >= 0", p.CoefficientField, len(p.Span.Basis.Members), p.MatrixID)
}
func (p CoordinateQuadraticPositivity) Key() string {
	return fmt.Sprintf("coordinate-positive|%s|%s|%s", p.MatrixID, p.Span.Key(), p.CoefficientField)
}

type QuadraticMatrixIdentity struct {
	Functional  FunctionalID            `json:"functional"`
	Combination FiniteLinearCombination `json:"linear_combination"`
	MatrixID    string                  `json:"matrix_id"`
	FormID      string                  `json:"form_id"`
	Identity    string                  `json:"identity"`
	Theorem     TheoremID               `json:"theorem"`
}

func (QuadraticMatrixIdentity) isProposition()        {}
func (QuadraticMatrixIdentity) Kind() PropositionKind { return QuadraticMatrixIdentityKind }
func (p QuadraticMatrixIdentity) Validate() error {
	if p.Functional == "" || p.MatrixID == "" || p.FormID == "" || p.Identity != "Q(sum_i c_i f_i)=c* G c" || p.Theorem == "" {
		return fmt.Errorf("invalid quadratic matrix identity")
	}
	return p.Combination.Validate()
}
func (p QuadraticMatrixIdentity) Describe() string {
	return fmt.Sprintf("%s(%s) = c* %s c", p.Functional, p.Combination.Describe(), p.MatrixID)
}
func (p QuadraticMatrixIdentity) Key() string {
	return fmt.Sprintf("quadratic-matrix|%s|%s|%s|%s", p.Functional, p.Combination.Span.Key(), p.MatrixID, p.Theorem)
}

type MatrixPropertyKind string

const (
	MatrixHermitian               MatrixPropertyKind = "hermitian"
	MatrixPositiveSemidefinite    MatrixPropertyKind = "positive_semidefinite"
	MatrixDiagonalNonnegative     MatrixPropertyKind = "diagonal_nonnegative"
	MatrixNumericallyPSDEstimated MatrixPropertyKind = "numerically_psd_estimated"
)

type MatrixProperty struct {
	MatrixID         string               `json:"matrix_id"`
	SourceFunctional FunctionalID         `json:"source_functional"`
	DomainSpan       *FiniteSpan          `json:"domain_span,omitempty"`
	Property         MatrixPropertyKind   `json:"property"`
	ValueSemantics   MatrixValueSemantics `json:"value_semantics"`
	Criterion        string               `json:"criterion"`
}

func (MatrixProperty) isProposition()        {}
func (MatrixProperty) Kind() PropositionKind { return MatrixPropertyStatementKind }
func (p MatrixProperty) Validate() error {
	if p.MatrixID == "" || p.SourceFunctional == "" || p.Criterion == "" {
		return fmt.Errorf("invalid matrix property")
	}
	if p.DomainSpan != nil {
		if err := p.DomainSpan.Validate(); err != nil {
			return err
		}
	}
	switch p.Property {
	case MatrixHermitian, MatrixPositiveSemidefinite, MatrixDiagonalNonnegative:
		if p.ValueSemantics != StructuralExactMatrix {
			return fmt.Errorf("exact matrix property requires a structural exact matrix")
		}
	case MatrixNumericallyPSDEstimated:
		if p.ValueSemantics != NumericallyApproximateMatrix {
			return fmt.Errorf("numerical PSD estimate must remain approximate")
		}
	default:
		return fmt.Errorf("unknown matrix property")
	}
	return nil
}
func (p MatrixProperty) Describe() string {
	return fmt.Sprintf("%s is %s (%s)", p.MatrixID, p.Property, p.ValueSemantics)
}
func (p MatrixProperty) Key() string {
	span := "none"
	if p.DomainSpan != nil {
		span = p.DomainSpan.Key()
	}
	return fmt.Sprintf("matrix-property|%s|%s|%s|%s|%s", p.MatrixID, p.SourceFunctional, span, p.Property, p.ValueSemantics)
}

func CloneOrderedBasis(b OrderedBasis) OrderedBasis {
	b.ParentClass = CloneFunctionClass(b.ParentClass)
	b.Members = append([]BasisMember(nil), b.Members...)
	for i := range b.Members {
		b.Members[i].Function = CloneTestFunction(b.Members[i].Function)
	}
	return b
}
func CloneFiniteSpan(s FiniteSpan) FiniteSpan {
	s.Basis = CloneOrderedBasis(s.Basis)
	s.ParentClass = CloneFunctionClass(s.ParentClass)
	return s
}
func CloneHermitianForm(h HermitianForm) HermitianForm {
	h.DomainSpan = CloneFiniteSpan(h.DomainSpan)
	h.Contributions = append([]FormContribution(nil), h.Contributions...)
	h.TheoremProvenance = append([]TheoremID(nil), h.TheoremProvenance...)
	return h
}
func CloneHermitianMatrix(m HermitianMatrix) HermitianMatrix {
	m.SourceForm = CloneHermitianForm(m.SourceForm)
	m.Basis = CloneOrderedBasis(m.Basis)
	m.Entries = append([]MatrixEntry(nil), m.Entries...)
	for i := range m.Entries {
		m.Entries[i].RowFunction = CloneTestFunction(m.Entries[i].RowFunction)
		m.Entries[i].ColumnFunction = CloneTestFunction(m.Entries[i].ColumnFunction)
		m.Entries[i].Contributions = append([]EntryContribution(nil), m.Entries[i].Contributions...)
		m.Entries[i].TheoremProvenance = append([]TheoremID(nil), m.Entries[i].TheoremProvenance...)
	}
	return m
}
