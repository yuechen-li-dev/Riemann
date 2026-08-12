package semantic

import (
	"fmt"
	"sort"
	"strings"
)

// M8 keeps geometric orbit cardinality and analytic zero multiplicity distinct.
type ZeroOrbitClass string

const (
	CriticalLineOrbit ZeroOrbitClass = "critical_line"
	OffCriticalOrbit  ZeroOrbitClass = "off_critical"
)

type OrbitPoint struct {
	Point       PointExpr        `json:"point"`
	GeneratedBy []PointTransform `json:"generated_by"`
}

type ZeroOrbit struct {
	ID                    string         `json:"id"`
	Representative        PointExpr      `json:"representative"`
	Classification        ZeroOrbitClass `json:"classification"`
	TransformedPoints     []OrbitPoint   `json:"distinct_transformed_points"`
	DistinctLocationCount int            `json:"distinct_location_count"`
	ZeroMultiplicity      int            `json:"zero_multiplicity"`
	SymmetryTheorems      []TheoremID    `json:"symmetry_theorems"`
	Provenance            Reference      `json:"provenance"`
}

func NewZeroOrbit(id string, representative PointExpr, class ZeroOrbitClass, multiplicity int, theorems []TheoremID, provenance Reference) (ZeroOrbit, error) {
	byKey := map[string]OrbitPoint{}
	for _, transform := range []PointTransform{IdentityTransform, ConjugateTransform, OneMinusTransform, CriticalReflectionTransform} {
		point := representative.Apply(transform).Canonical()
		key := point.Key()
		item, ok := byKey[key]
		if !ok {
			item = OrbitPoint{Point: point}
		}
		item.GeneratedBy = append(item.GeneratedBy, transform)
		byKey[key] = item
	}
	keys := make([]string, 0, len(byKey))
	for key := range byKey {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	points := make([]OrbitPoint, 0, len(keys))
	for _, key := range keys {
		points = append(points, byKey[key])
	}
	orbit := ZeroOrbit{ID: id, Representative: representative.Canonical(), Classification: class, TransformedPoints: points, DistinctLocationCount: len(points), ZeroMultiplicity: multiplicity, SymmetryTheorems: append([]TheoremID(nil), theorems...), Provenance: provenance}
	return orbit, orbit.Validate()
}

func (o ZeroOrbit) Validate() error {
	if strings.TrimSpace(o.ID) == "" || o.ZeroMultiplicity < 1 || len(o.SymmetryTheorems) == 0 || o.Provenance.Citation == "" || o.DistinctLocationCount != len(o.TransformedPoints) {
		return fmt.Errorf("invalid zero orbit identity, multiplicity, or provenance")
	}
	if err := o.Representative.Validate(); err != nil {
		return err
	}
	if o.Classification != CriticalLineOrbit && o.Classification != OffCriticalOrbit {
		return fmt.Errorf("invalid zero orbit classification")
	}
	if o.Classification == CriticalLineOrbit && o.Representative.Location != CriticalLinePoint {
		return fmt.Errorf("critical orbit representative is not typed on the critical line")
	}
	if o.Classification == OffCriticalOrbit && o.Representative.Location == CriticalLinePoint {
		return fmt.Errorf("off-critical orbit representative is typed on the critical line")
	}
	seen := map[string]bool{}
	for _, item := range o.TransformedPoints {
		if err := item.Point.Validate(); err != nil {
			return err
		}
		if seen[item.Point.Key()] || len(item.GeneratedBy) == 0 {
			return fmt.Errorf("duplicate or unexplained orbit point")
		}
		seen[item.Point.Key()] = true
	}
	return nil
}

func (o ZeroOrbit) Key() string {
	points := make([]string, len(o.TransformedPoints))
	for i, p := range o.TransformedPoints {
		points[i] = p.Point.Key()
	}
	return fmt.Sprintf("%s|%s|%s|locations=%s|multiplicity=%d", o.ID, o.Representative.Key(), o.Classification, strings.Join(points, ";"), o.ZeroMultiplicity)
}

type SymbolicEvidenceStatus string

const ExactSymbolicEvidence SymbolicEvidenceStatus = "exact_symbolic"

type BasisEvaluationEntry struct {
	Index      int          `json:"index"`
	Function   TestFunction `json:"function"`
	Expression string       `json:"expression"`
}

type BasisEvaluationVector struct {
	ID                  string                 `json:"id"`
	Basis               OrderedBasis           `json:"basis"`
	Point               PointExpr              `json:"point"`
	TransformConvention TransformConventionID  `json:"transform_convention"`
	Dimension           int                    `json:"dimension"`
	Entries             []BasisEvaluationEntry `json:"entries"`
	EvidenceStatus      SymbolicEvidenceStatus `json:"evidence_status"`
	TheoremProvenance   []TheoremID            `json:"theorem_provenance"`
	Provenance          Reference              `json:"provenance"`
}

func NewBasisEvaluationVector(basis OrderedBasis, point PointExpr, theorem TheoremID, provenance Reference) (BasisEvaluationVector, error) {
	entries := make([]BasisEvaluationEntry, len(basis.Members))
	for i, member := range basis.Members {
		entries[i] = BasisEvaluationEntry{Index: i, Function: CloneTestFunction(member.Function), Expression: fmt.Sprintf("M[%s](%s)", member.Function.Symbol, point.Describe())}
	}
	v := BasisEvaluationVector{ID: "v[" + basis.ID + ";" + point.Key() + "]", Basis: CloneOrderedBasis(basis), Point: point.Canonical(), TransformConvention: basis.ParentClass.TransformConvention, Dimension: len(entries), Entries: entries, EvidenceStatus: ExactSymbolicEvidence, TheoremProvenance: []TheoremID{theorem}, Provenance: provenance}
	return v, v.Validate()
}

func (v BasisEvaluationVector) Validate() error {
	if v.ID == "" || v.Dimension != len(v.Entries) || v.Dimension == 0 || v.EvidenceStatus != ExactSymbolicEvidence || len(v.TheoremProvenance) == 0 || v.Provenance.Citation == "" {
		return fmt.Errorf("invalid basis evaluation vector")
	}
	if err := v.Basis.Validate(); err != nil {
		return err
	}
	if err := v.Point.Validate(); err != nil {
		return err
	}
	if v.TransformConvention != v.Basis.ParentClass.TransformConvention {
		return fmt.Errorf("evaluation vector transform convention disagrees with basis")
	}
	for i, e := range v.Entries {
		if e.Index != i || e.Function.Key() != v.Basis.Members[i].Function.Key() || e.Expression == "" {
			return fmt.Errorf("evaluation vector entry lost coordinate identity")
		}
	}
	return nil
}
func (v BasisEvaluationVector) Key() string {
	return v.Basis.Key() + "|" + v.Point.Key() + "|" + string(v.TransformConvention)
}

type SymbolicMatrixEntry struct {
	Row        int    `json:"row"`
	Column     int    `json:"column"`
	Expression string `json:"expression"`
}

type ContributionClassification struct {
	Hermitian            bool        `json:"hermitian_certified"`
	PositiveSemidefinite bool        `json:"positive_semidefinite_certified"`
	RankUpperBound       int         `json:"rank_upper_bound"`
	RankOneIfNonzero     bool        `json:"rank_one_if_evaluation_nonzero"`
	Indefinite           bool        `json:"indefinite_certified"`
	IndefiniteCondition  string      `json:"indefinite_condition,omitempty"`
	DegenerateCondition  string      `json:"degenerate_condition,omitempty"`
	Theorems             []TheoremID `json:"theorems"`
}

type PointMatrixContribution struct {
	ID                  string                     `json:"id"`
	Point               PointExpr                  `json:"point"`
	ReflectedPoint      PointExpr                  `json:"critical_reflection"`
	Evaluation          BasisEvaluationVector      `json:"evaluation_vector"`
	ReflectedEvaluation BasisEvaluationVector      `json:"reflected_evaluation_vector"`
	Rows                int                        `json:"rows"`
	Columns             int                        `json:"columns"`
	Entries             []SymbolicMatrixEntry      `json:"entries"`
	QuadraticIdentity   string                     `json:"quadratic_identity"`
	SourceSummand       string                     `json:"source_weil_summand"`
	Orientation         string                     `json:"orientation"`
	Multiplicity        int                        `json:"zero_multiplicity"`
	Classification      ContributionClassification `json:"classification"`
	TheoremProvenance   []TheoremID                `json:"theorem_provenance"`
}

func (k PointMatrixContribution) Validate() error {
	if k.ID == "" || k.Rows == 0 || k.Rows != k.Columns || len(k.Entries) != k.Rows*k.Columns || k.Multiplicity < 1 || k.SourceSummand != "M[f](rho) * conjugate(M[f](1-conjugate(rho)))" || k.Orientation != "K_ij=conjugate(w_i)*v_j; c* K c=(c^T v)*conjugate(c^T w)" || len(k.TheoremProvenance) == 0 {
		return fmt.Errorf("invalid point matrix contribution")
	}
	if k.ReflectedPoint.Key() != k.Point.Apply(CriticalReflectionTransform).Key() || k.Evaluation.Point.Key() != k.Point.Key() || k.ReflectedEvaluation.Point.Key() != k.ReflectedPoint.Key() {
		return fmt.Errorf("point contribution reflection/evaluation mismatch")
	}
	for n, e := range k.Entries {
		if e.Row != n/k.Columns || e.Column != n%k.Columns || e.Expression == "" {
			return fmt.Errorf("invalid symbolic contribution entry")
		}
	}
	return nil
}

type ReflectionPairContribution struct {
	ID                 string                     `json:"id"`
	Points             []PointExpr                `json:"distinct_points"`
	PointContributions []PointMatrixContribution  `json:"point_contributions"`
	Entries            []SymbolicMatrixEntry      `json:"entries"`
	Formula            string                     `json:"formula"`
	Classification     ContributionClassification `json:"classification"`
	GroupingTheorem    TheoremID                  `json:"grouping_theorem"`
}

type OrbitMatrixContribution struct {
	ID                string                       `json:"id"`
	Orbit             ZeroOrbit                    `json:"orbit"`
	Basis             OrderedBasis                 `json:"basis"`
	ReflectionPairs   []ReflectionPairContribution `json:"critical_reflection_pairs"`
	Entries           []SymbolicMatrixEntry        `json:"entries"`
	Formula           string                       `json:"formula"`
	Classification    ContributionClassification   `json:"classification"`
	TheoremProvenance []TheoremID                  `json:"theorem_provenance"`
}

type ZeroSideMatrixAggregate struct {
	ID                    string                    `json:"id"`
	SemanticMatrixID      string                    `json:"semantic_matrix_id"`
	Basis                 OrderedBasis              `json:"basis"`
	ZeroDomain            Domain                    `json:"zero_domain"`
	OrbitContributions    []OrbitMatrixContribution `json:"instantiated_orbit_contributions"`
	ContributionTemplates []OrbitMatrixContribution `json:"orbit_contribution_templates"`
	SummationConvention   []string                  `json:"symmetric_limiting_convention"`
	Formula               string                    `json:"formula"`
	CriticalAggregate     string                    `json:"critical_aggregate"`
	OffCriticalAggregate  string                    `json:"off_critical_aggregate"`
	SplitIdentity         string                    `json:"split_identity"`
	SourceFunctional      FunctionalID              `json:"source_functional"`
	TransformConvention   TransformConventionID     `json:"transform_convention"`
	TheoremProvenance     []TheoremID               `json:"theorem_provenance"`
	Provenance            Reference                 `json:"provenance"`
}

func (a ZeroSideMatrixAggregate) Validate() error {
	if a.ID == "" || a.SemanticMatrixID == "" || a.Formula != "G=sum_over_zero_orbits G_O" || a.SplitIdentity != "G=P+Q" || a.CriticalAggregate == "" || a.OffCriticalAggregate == "" || len(a.ContributionTemplates) != 2 || len(a.SummationConvention) == 0 || len(a.TheoremProvenance) == 0 || a.Provenance.Citation == "" {
		return fmt.Errorf("invalid zero-side matrix aggregate")
	}
	if err := a.Basis.Validate(); err != nil {
		return err
	}
	if err := a.ZeroDomain.Validate(); err != nil {
		return err
	}
	if a.ZeroDomain != NontrivialZeros(RiemannZeta) || a.SourceFunctional != WeilZetaQuadraticFunctional || a.TransformConvention != LagariasMellinConvention {
		return fmt.Errorf("zero-side aggregate changed M4 semantics")
	}
	return nil
}

type DualMatrixRepresentation struct {
	SemanticMatrixID        string            `json:"semantic_matrix_id"`
	ZeroSideAggregateID     string            `json:"zero_side_aggregate_id"`
	ExplicitFormulaMatrixID string            `json:"explicit_formula_matrix_id"`
	ExplicitValueEvidence   ValueEvidenceKind `json:"explicit_value_evidence"`
	IdentityTheorem         TheoremID         `json:"identity_theorem"`
	Identity                string            `json:"identity"`
	NumericalIdentification bool              `json:"numerical_identification"`
	Provenance              Reference         `json:"provenance"`
}
