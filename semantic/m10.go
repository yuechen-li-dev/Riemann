package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

// BoundaryConvention is part of the meaning of a height window.  M10 uses the
// convention of the source paper: lower < ordinate <= upper.
type BoundaryConvention string

const LeftOpenRightClosed BoundaryConvention = "left_open_right_closed"

type OrdinateConvention string

const PositiveOrdinateConvention OrdinateConvention = "positive_ordinates_only"

// HeightCoordinate carries either an exact integer fixture or a symbolic
// expression.  Symbolic coordinates make the T-parameterized theorem honest;
// concrete coordinates make boundary ownership executable without floats.
type HeightCoordinate struct {
	Expression   string `json:"expression"`
	ExactInteger *int64 `json:"exact_integer,omitempty"`
}

func SymbolicHeight(expression string) HeightCoordinate {
	return HeightCoordinate{Expression: expression}
}

func ExactHeight(value int64) HeightCoordinate {
	return HeightCoordinate{Expression: fmt.Sprintf("%d", value), ExactInteger: &value}
}

func (h HeightCoordinate) Validate() error {
	if strings.TrimSpace(h.Expression) == "" {
		return fmt.Errorf("height coordinate lacks an expression")
	}
	return nil
}

type HeightWindow struct {
	ID                 string             `json:"id"`
	Center             HeightCoordinate   `json:"center"`
	HalfWidth          HeightCoordinate   `json:"half_width"`
	Lower              HeightCoordinate   `json:"lower"`
	Upper              HeightCoordinate   `json:"upper"`
	Boundary           BoundaryConvention `json:"boundary"`
	OrdinateConvention OrdinateConvention `json:"ordinate_convention"`
}

func NewConcreteHeightWindow(id string, lower, upper int64) (HeightWindow, error) {
	if lower < 0 || upper <= lower || (upper-lower)%2 != 0 {
		return HeightWindow{}, fmt.Errorf("invalid concrete height window")
	}
	w := HeightWindow{ID: id, Center: ExactHeight(lower + (upper-lower)/2), HalfWidth: ExactHeight((upper - lower) / 2), Lower: ExactHeight(lower), Upper: ExactHeight(upper), Boundary: LeftOpenRightClosed, OrdinateConvention: PositiveOrdinateConvention}
	return w, w.Validate()
}

func (w HeightWindow) Validate() error {
	if strings.TrimSpace(w.ID) == "" || w.Boundary != LeftOpenRightClosed || w.OrdinateConvention != PositiveOrdinateConvention {
		return fmt.Errorf("invalid height-window convention")
	}
	for _, h := range []HeightCoordinate{w.Center, w.HalfWidth, w.Lower, w.Upper} {
		if err := h.Validate(); err != nil {
			return err
		}
	}
	if w.Lower.ExactInteger != nil && w.Upper.ExactInteger != nil && (*w.Lower.ExactInteger < 0 || *w.Upper.ExactInteger <= *w.Lower.ExactInteger) {
		return fmt.Errorf("invalid concrete height-window endpoints")
	}
	if w.HalfWidth.ExactInteger != nil && *w.HalfWidth.ExactInteger <= 0 {
		return fmt.Errorf("invalid concrete height-window half-width")
	}
	return nil
}

func (w HeightWindow) ContainsExactOrdinate(ordinate int64) (bool, error) {
	if err := w.Validate(); err != nil {
		return false, err
	}
	if w.Lower.ExactInteger == nil || w.Upper.ExactInteger == nil {
		return false, fmt.Errorf("symbolic height window has no executable membership test")
	}
	return ordinate > *w.Lower.ExactInteger && ordinate <= *w.Upper.ExactInteger, nil
}

type WindowZeroPoint struct {
	LocationID       string `json:"location_id"`
	Ordinate         int64  `json:"ordinate"`
	Multiplicity     int    `json:"multiplicity"`
	CriticalLine     bool   `json:"critical_line"`
	ReflectionPairID string `json:"reflection_pair_id,omitempty"`
}

type WindowZeroCounts struct {
	TotalMultiplicity          int `json:"total_multiplicity"`
	CriticalMultiplicity       int `json:"critical_multiplicity"`
	CriticalDistinctLocations  int `json:"critical_distinct_locations"`
	SimpleCriticalLocations    int `json:"simple_critical_locations"`
	OffCriticalMultiplicity    int `json:"off_critical_multiplicity"`
	OffCriticalReflectionPairs int `json:"off_critical_reflection_pairs"`
	DistinctAllLocations       int `json:"distinct_all_locations"`
	EvaluationRankDirections   int `json:"evaluation_rank_directions"`
}

func (c WindowZeroCounts) Validate() error {
	values := []int{c.TotalMultiplicity, c.CriticalMultiplicity, c.CriticalDistinctLocations, c.SimpleCriticalLocations, c.OffCriticalMultiplicity, c.OffCriticalReflectionPairs, c.DistinctAllLocations, c.EvaluationRankDirections}
	for _, v := range values {
		if v < 0 {
			return fmt.Errorf("negative zero count")
		}
	}
	if c.TotalMultiplicity != c.CriticalMultiplicity+c.OffCriticalMultiplicity || c.SimpleCriticalLocations > c.CriticalDistinctLocations || c.CriticalDistinctLocations > c.CriticalMultiplicity || c.DistinctAllLocations != c.CriticalDistinctLocations+2*c.OffCriticalReflectionPairs || c.DistinctAllLocations > c.TotalMultiplicity || c.EvaluationRankDirections > c.CriticalDistinctLocations {
		return fmt.Errorf("conflated zero-count notions")
	}
	if 2*c.OffCriticalReflectionPairs > c.OffCriticalMultiplicity {
		return fmt.Errorf("reflection-pair count exceeds off-critical multiplicity")
	}
	return nil
}

// CountWindowZeros owns each point by ordinate.  Positive-ordinate windows do
// not acquire the conjugate point at negative ordinate.  Reflection partners
// at the same ordinate remain two locations but one unordered pair.
func CountWindowZeros(window HeightWindow, points []WindowZeroPoint) (WindowZeroCounts, error) {
	locations := map[string]WindowZeroPoint{}
	pairs := map[string][]WindowZeroPoint{}
	for _, point := range points {
		inside, err := window.ContainsExactOrdinate(point.Ordinate)
		if err != nil {
			return WindowZeroCounts{}, err
		}
		if !inside {
			continue
		}
		if point.LocationID == "" || point.Multiplicity < 1 {
			return WindowZeroCounts{}, fmt.Errorf("invalid zero point")
		}
		if _, duplicate := locations[point.LocationID]; duplicate {
			return WindowZeroCounts{}, fmt.Errorf("duplicate geometric zero location")
		}
		if !point.CriticalLine && point.ReflectionPairID == "" {
			return WindowZeroCounts{}, fmt.Errorf("off-critical point lacks reflection-pair ownership")
		}
		if point.CriticalLine && point.ReflectionPairID != "" {
			return WindowZeroCounts{}, fmt.Errorf("critical point cannot own an off-critical reflection pair")
		}
		locations[point.LocationID] = point
		if point.ReflectionPairID != "" {
			pairs[point.ReflectionPairID] = append(pairs[point.ReflectionPairID], point)
		}
	}
	for id, pair := range pairs {
		if len(pair) != 2 || pair[0].Multiplicity != pair[1].Multiplicity {
			return WindowZeroCounts{}, fmt.Errorf("reflection pair %q lacks exactly two equal-multiplicity partners", id)
		}
	}
	counts := WindowZeroCounts{DistinctAllLocations: len(locations), OffCriticalReflectionPairs: len(pairs)}
	for _, point := range locations {
		counts.TotalMultiplicity += point.Multiplicity
		if point.CriticalLine {
			counts.CriticalMultiplicity += point.Multiplicity
			counts.CriticalDistinctLocations++
			if point.Multiplicity == 1 {
				counts.SimpleCriticalLocations++
			}
		} else {
			counts.OffCriticalMultiplicity += point.Multiplicity
		}
	}
	return counts, counts.Validate()
}

type ThresholdComparison string

const StrictlyAboveThreshold ThresholdComparison = "strictly_above"

type ThresholdScaleKind string

const (
	AbsoluteThresholdScale   ThresholdScaleKind = "absolute"
	FarOperatorBoundScale    ThresholdScaleKind = "far_operator_bound"
	AsymptoticThresholdScale ThresholdScaleKind = "asymptotic_parameterized"
)

type ExactRational struct {
	Numerator   int64 `json:"numerator"`
	Denominator int64 `json:"denominator"`
}

func (r ExactRational) Validate() error {
	if r.Denominator <= 0 {
		return fmt.Errorf("invalid exact rational")
	}
	return nil
}

func (r ExactRational) GreaterThan(other ExactRational) bool {
	left := new(big.Int).Mul(big.NewInt(r.Numerator), big.NewInt(other.Denominator))
	right := new(big.Int).Mul(big.NewInt(other.Numerator), big.NewInt(r.Denominator))
	return left.Cmp(right) > 0
}

type SpectralThreshold struct {
	ID           string              `json:"id"`
	Kind         ThresholdScaleKind  `json:"kind"`
	Expression   string              `json:"expression"`
	ExactValue   *ExactRational      `json:"exact_value,omitempty"`
	Comparison   ThresholdComparison `json:"comparison"`
	Dependencies map[string]string   `json:"dependencies"`
	Provenance   Reference           `json:"provenance"`
}

func (t SpectralThreshold) Validate() error {
	if t.ID == "" || t.Expression == "" || t.Comparison != StrictlyAboveThreshold || len(t.Dependencies) == 0 || t.Provenance.Citation == "" {
		return fmt.Errorf("invalid spectral threshold")
	}
	if t.Kind != AbsoluteThresholdScale && t.Kind != FarOperatorBoundScale && t.Kind != AsymptoticThresholdScale {
		return fmt.Errorf("invalid threshold scale")
	}
	if t.ExactValue != nil {
		return t.ExactValue.Validate()
	}
	return nil
}

type ThresholdedPositiveIndexClaim struct {
	MatrixID   string               `json:"matrix_id"`
	Dimension  int                  `json:"dimension"`
	Threshold  SpectralThreshold    `json:"threshold"`
	Relation   BoundRelation        `json:"relation"`
	Bound      int                  `json:"bound"`
	Evidence   SpectralEvidenceKind `json:"evidence"`
	Theorems   []TheoremID          `json:"theorems"`
	Provenance Reference            `json:"provenance"`
}

func (c ThresholdedPositiveIndexClaim) Validate() error {
	if c.MatrixID == "" || c.Dimension < 1 || !c.Relation.Valid() || c.Bound < 0 || c.Bound > c.Dimension || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid thresholded positive-index claim")
	}
	if err := c.Threshold.Validate(); err != nil {
		return err
	}
	if c.Evidence != ExactTheoremEvidence && c.Evidence != CertifiedMinorEvidence && c.Evidence != ApproximateEigenEvidence {
		return fmt.Errorf("invalid spectral evidence")
	}
	if c.Evidence != ApproximateEigenEvidence && len(c.Theorems) == 0 {
		return fmt.Errorf("exact thresholded count lacks theorem provenance")
	}
	return nil
}

func (c ThresholdedPositiveIndexClaim) ExactTheoremPremise() bool {
	return c.Evidence == ExactTheoremEvidence || c.Evidence == CertifiedMinorEvidence
}

func ExactDiagonalThresholdedPositiveIndex(eigenvalues []ExactRational, threshold ExactRational) (int, error) {
	if err := threshold.Validate(); err != nil {
		return 0, err
	}
	count := 0
	for _, eigenvalue := range eigenvalues {
		if err := eigenvalue.Validate(); err != nil {
			return 0, err
		}
		if eigenvalue.GreaterThan(threshold) {
			count++
		}
	}
	return count, nil
}

type WindowCompression struct {
	ID                    string            `json:"id"`
	Window                HeightWindow      `json:"window"`
	LocalizationWindow    HeightWindow      `json:"localization_window"`
	BasisFamily           string            `json:"basis_family"`
	BasisHeightDependency string            `json:"basis_height_dependency"`
	DimensionExpression   string            `json:"dimension_expression"`
	MatrixID              string            `json:"matrix_id"`
	NearMatrixID          string            `json:"near_matrix_id"`
	FarMatrixID           string            `json:"far_matrix_id"`
	ZeroSideIdentity      string            `json:"zero_side_identity"`
	ExplicitFormula       string            `json:"explicit_formula_representation"`
	Normalization         string            `json:"normalization"`
	Threshold             SpectralThreshold `json:"threshold"`
	Provenance            Reference         `json:"provenance"`
}

func (c WindowCompression) Validate() error {
	if c.ID == "" || c.BasisFamily == "" || c.BasisHeightDependency == "" || c.DimensionExpression == "" || c.MatrixID == "" || c.NearMatrixID == "" || c.FarMatrixID == "" || c.ZeroSideIdentity != c.MatrixID+"="+c.NearMatrixID+"+"+c.FarMatrixID || c.ExplicitFormula == "" || c.Normalization == "" || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid window compression")
	}
	if err := c.Window.Validate(); err != nil {
		return err
	}
	if err := c.LocalizationWindow.Validate(); err != nil {
		return err
	}
	return c.Threshold.Validate()
}

type NearFarZeroDecomposition struct {
	MatrixID              string      `json:"matrix_id"`
	NearMatrixID          string      `json:"near_matrix_id"`
	CriticalNearMatrixID  string      `json:"critical_near_matrix_id"`
	OffCriticalNearMatrix string      `json:"off_critical_near_matrix_id"`
	FarMatrixID           string      `json:"far_matrix_id"`
	Identity              string      `json:"identity"`
	MembershipRule        string      `json:"membership_rule"`
	Theorems              []TheoremID `json:"theorems"`
	Provenance            Reference   `json:"provenance"`
}

func (d NearFarZeroDecomposition) Validate() error {
	if d.MatrixID == "" || d.NearMatrixID == "" || d.CriticalNearMatrixID == "" || d.OffCriticalNearMatrix == "" || d.FarMatrixID == "" || d.Identity == "" || d.MembershipRule == "" || len(d.Theorems) == 0 || d.Provenance.Citation == "" {
		return fmt.Errorf("invalid near/far decomposition")
	}
	return nil
}

type MatrixNormKind string

const OperatorNorm MatrixNormKind = "operator_norm"

type FarZeroContributionBound struct {
	MatrixID            string         `json:"matrix_id"`
	Norm                MatrixNormKind `json:"norm"`
	BoundSymbol         string         `json:"bound_symbol"`
	BoundExpression     string         `json:"bound_expression"`
	Assumptions         []string       `json:"assumptions"`
	AsymptoticStatement string         `json:"asymptotic_statement"`
	Uniformity          string         `json:"uniformity"`
	ExactOrTrusted      bool           `json:"exact_or_trusted"`
	Theorems            []TheoremID    `json:"theorems"`
	Provenance          Reference      `json:"provenance"`
}

func (b FarZeroContributionBound) Validate() error {
	if b.MatrixID == "" || b.Norm != OperatorNorm || b.BoundSymbol == "" || b.BoundExpression == "" || len(b.Assumptions) == 0 || b.AsymptoticStatement == "" || b.Uniformity == "" || !b.ExactOrTrusted || len(b.Theorems) == 0 || b.Provenance.Citation == "" {
		return fmt.Errorf("invalid far-zero contribution bound")
	}
	return nil
}

type ThresholdPerturbationContract struct {
	ID            TheoremID      `json:"id"`
	Statement     string         `json:"statement"`
	ThresholdRule string         `json:"threshold_rule"`
	Comparison    string         `json:"comparison"`
	Norm          MatrixNormKind `json:"norm"`
	ExactRequired bool           `json:"exact_required"`
	Provenance    Reference      `json:"provenance"`
}

func (c ThresholdPerturbationContract) Validate() error {
	if c.ID == "" || c.Statement == "" || c.ThresholdRule == "" || c.Comparison != "strict" || c.Norm != OperatorNorm || !c.ExactRequired || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid threshold perturbation contract")
	}
	return nil
}

type FiniteWindowCountingTheorem struct {
	Name                        string      `json:"name"`
	Assumptions                 []string    `json:"assumptions"`
	ThresholdTransfer           string      `json:"threshold_transfer"`
	M9Accounting                string      `json:"m9_accounting"`
	EnlargedWindowBound         string      `json:"enlarged_window_bound"`
	TargetWindowBound           string      `json:"target_window_bound"`
	DistinctZeroBound           string      `json:"distinct_zero_bound"`
	CountConversion             string      `json:"count_conversion"`
	RemainingInput              string      `json:"remaining_input"`
	Theorems                    []TheoremID `json:"theorems"`
	KnownLiteratureResult       bool        `json:"known_literature_result"`
	StructurallyReconstructed   bool        `json:"structurally_reconstructed"`
	AsymptoticProportionDerived bool        `json:"asymptotic_proportion_derived"`
	Provenance                  Reference   `json:"provenance"`
}

func (t FiniteWindowCountingTheorem) Validate() error {
	if t.Name == "" || len(t.Assumptions) == 0 || t.ThresholdTransfer == "" || t.M9Accounting == "" || t.EnlargedWindowBound == "" || t.TargetWindowBound == "" || t.DistinctZeroBound == "" || t.CountConversion == "" || t.RemainingInput == "" || len(t.Theorems) == 0 || !t.KnownLiteratureResult || !t.StructurallyReconstructed || t.AsymptoticProportionDerived || t.Provenance.Citation == "" {
		return fmt.Errorf("invalid finite window-counting theorem")
	}
	return nil
}
