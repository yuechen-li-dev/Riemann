package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

// SpectralMomentKind keeps the linear and quadratic spectral statistics
// distinct in the IR. Dimension is included because threshold splitting needs
// the number of eigenvalues, not merely a generic matrix "size".
type SpectralMomentKind string

const (
	Trace                SpectralMomentKind = "trace"
	FrobeniusNormSquared SpectralMomentKind = "frobenius_norm_squared"
	Dimension            SpectralMomentKind = "dimension"
)

type MomentEvidenceKind string

const (
	ExactMomentEvidence       MomentEvidenceKind = "exact"
	TrustedMomentEvidence     MomentEvidenceKind = "trusted_theorem"
	AsymptoticMomentEvidence  MomentEvidenceKind = "trusted_asymptotic"
	ApproximateMomentEvidence MomentEvidenceKind = "approximate_numerical"
)

type SpectralMomentClaim struct {
	ID          string             `json:"id"`
	MatrixID    string             `json:"matrix_id"`
	Kind        SpectralMomentKind `json:"kind"`
	Relation    BoundRelation      `json:"relation"`
	Expression  string             `json:"expression"`
	ExactValue  *ExactRational     `json:"exact_value,omitempty"`
	Evidence    MomentEvidenceKind `json:"evidence"`
	Assumptions []string           `json:"assumptions"`
	Theorems    []TheoremID        `json:"theorems"`
	Provenance  Reference          `json:"provenance"`
}

func (c SpectralMomentClaim) Validate() error {
	if c.ID == "" || c.MatrixID == "" || c.Expression == "" || !c.Relation.Valid() || c.Provenance.Citation == "" {
		return fmt.Errorf("invalid spectral-moment claim")
	}
	if c.Kind != Trace && c.Kind != FrobeniusNormSquared && c.Kind != Dimension {
		return fmt.Errorf("unknown spectral moment %q", c.Kind)
	}
	if c.Evidence != ExactMomentEvidence && c.Evidence != TrustedMomentEvidence && c.Evidence != AsymptoticMomentEvidence && c.Evidence != ApproximateMomentEvidence {
		return fmt.Errorf("unknown moment evidence %q", c.Evidence)
	}
	if c.ExactValue != nil {
		if err := c.ExactValue.Validate(); err != nil {
			return err
		}
		if c.Evidence != ExactMomentEvidence && c.Evidence != TrustedMomentEvidence {
			return fmt.Errorf("non-finite moment evidence cannot carry an exact theorem value")
		}
	}
	if c.Evidence != ApproximateMomentEvidence && len(c.Theorems) == 0 {
		return fmt.Errorf("theorem-grade moment lacks theorem provenance")
	}
	return nil
}

// FiniteTheoremPremise deliberately excludes raw asymptotic and approximate
// claims. An asymptotic statement first needs an EventuallyBound discharge.
func (c SpectralMomentClaim) FiniteTheoremPremise() bool {
	return (c.Evidence == ExactMomentEvidence || c.Evidence == TrustedMomentEvidence) && c.ExactValue != nil
}

type RemainderKind string

const LittleORemainder RemainderKind = "little_o"

type AsymptoticMomentStatement struct {
	Moment        SpectralMomentClaim `json:"moment"`
	MainTerm      string              `json:"main_term"`
	Remainder     string              `json:"remainder"`
	RemainderKind RemainderKind       `json:"remainder_kind"`
	Scale         string              `json:"scale"`
	Parameter     string              `json:"parameter"`
	Uniformity    string              `json:"uniformity"`
}

func (s AsymptoticMomentStatement) Validate() error {
	if err := s.Moment.Validate(); err != nil {
		return err
	}
	if s.Moment.Evidence != AsymptoticMomentEvidence || s.MainTerm == "" || s.Remainder == "" || s.RemainderKind != LittleORemainder || s.Scale == "" || s.Parameter == "" || s.Uniformity == "" {
		return fmt.Errorf("invalid asymptotic moment statement")
	}
	return nil
}

type EventuallyBound struct {
	SourceMomentID string        `json:"source_moment_id"`
	Epsilon        string        `json:"epsilon"`
	Threshold      string        `json:"parameter_threshold"`
	FiniteBound    string        `json:"finite_bound"`
	Relation       BoundRelation `json:"relation"`
	Theorem        TheoremID     `json:"theorem"`
	Provenance     Reference     `json:"provenance"`
}

func (b EventuallyBound) Validate() error {
	if b.SourceMomentID == "" || b.Epsilon == "" || b.Threshold == "" || b.FiniteBound == "" || !b.Relation.Valid() || b.Theorem == "" || b.Provenance.Citation == "" {
		return fmt.Errorf("invalid eventually-finite bound")
	}
	return nil
}

type ExactComplexRational struct {
	Real ExactRational `json:"real"`
	Imag ExactRational `json:"imag"`
}

func (z ExactComplexRational) Validate() error {
	if err := z.Real.Validate(); err != nil {
		return err
	}
	return z.Imag.Validate()
}

func rat(r ExactRational) *big.Rat {
	return new(big.Rat).SetFrac(big.NewInt(r.Numerator), big.NewInt(r.Denominator))
}

func exactRat(r *big.Rat) (ExactRational, error) {
	if !r.Num().IsInt64() || !r.Denom().IsInt64() {
		return ExactRational{}, fmt.Errorf("exact rational exceeds IR integer range")
	}
	return ExactRational{Numerator: r.Num().Int64(), Denominator: r.Denom().Int64()}, nil
}

// ExactHermitianMoments computes tr(G) and sum_ij |G_ij|^2 exactly. For a
// Hermitian matrix the latter equals tr(G^2)=sum_i lambda_i^2.
func ExactHermitianMoments(dimension int, entries []ExactComplexRational) (ExactRational, ExactRational, error) {
	if dimension < 1 || len(entries) != dimension*dimension {
		return ExactRational{}, ExactRational{}, fmt.Errorf("invalid exact matrix shape")
	}
	zero := new(big.Rat)
	trace, frob := new(big.Rat), new(big.Rat)
	for _, z := range entries {
		if err := z.Validate(); err != nil {
			return ExactRational{}, ExactRational{}, err
		}
	}
	for i, z := range entries {
		r, im := rat(z.Real), rat(z.Imag)
		frob.Add(frob, new(big.Rat).Mul(r, r))
		frob.Add(frob, new(big.Rat).Mul(im, im))
		row, col := i/dimension, i%dimension
		if row == col {
			if im.Cmp(zero) != 0 {
				return ExactRational{}, ExactRational{}, fmt.Errorf("Hermitian diagonal is not real")
			}
			trace.Add(trace, r)
		}
		mirror := entries[col*dimension+row]
		if rat(z.Real).Cmp(rat(mirror.Real)) != 0 || rat(z.Imag).Cmp(new(big.Rat).Neg(rat(mirror.Imag))) != 0 {
			return ExactRational{}, ExactRational{}, fmt.Errorf("matrix is not Hermitian")
		}
	}
	t, err := exactRat(trace)
	if err != nil {
		return ExactRational{}, ExactRational{}, err
	}
	f, err := exactRat(frob)
	if err != nil {
		return ExactRational{}, ExactRational{}, err
	}
	return t, f, nil
}

type MomentCountResult struct {
	Applicable        bool          `json:"applicable"`
	TraceResidual     ExactRational `json:"trace_residual"`
	RealLowerBound    ExactRational `json:"real_lower_bound"`
	IntegerLowerBound int           `json:"integer_lower_bound"`
}

// ThresholdedCountFromMoments is the executable exact-rational form of
// Lemma 3.3: if theta>=0, tr(G)>=A, tr(G^2)<=B, dim(G)=d, and A-d*theta>0,
// then n_+^theta(G) >= ceil((A-d*theta)^2/B).
func ThresholdedCountFromMoments(traceLower, frobeniusSquaredUpper ExactRational, dimension int, threshold ExactRational) (MomentCountResult, error) {
	for _, r := range []ExactRational{traceLower, frobeniusSquaredUpper, threshold} {
		if err := r.Validate(); err != nil {
			return MomentCountResult{}, err
		}
	}
	if dimension < 1 || rat(threshold).Sign() < 0 {
		return MomentCountResult{}, fmt.Errorf("moment theorem requires positive dimension and nonnegative threshold")
	}
	residual := new(big.Rat).Sub(rat(traceLower), new(big.Rat).Mul(big.NewRat(int64(dimension), 1), rat(threshold)))
	residualExact, err := exactRat(residual)
	if err != nil {
		return MomentCountResult{}, err
	}
	if rat(frobeniusSquaredUpper).Sign() <= 0 {
		return MomentCountResult{}, fmt.Errorf("moment theorem requires a positive Frobenius-square upper bound")
	}
	if residual.Sign() <= 0 {
		return MomentCountResult{Applicable: false, TraceResidual: residualExact, RealLowerBound: ExactRational{Numerator: 0, Denominator: 1}}, nil
	}
	x := new(big.Rat).Quo(new(big.Rat).Mul(residual, residual), rat(frobeniusSquaredUpper))
	xExact, err := exactRat(x)
	if err != nil {
		return MomentCountResult{}, err
	}
	q, rem := new(big.Int).QuoRem(x.Num(), x.Denom(), new(big.Int))
	if rem.Sign() > 0 {
		q.Add(q, big.NewInt(1))
	}
	if !q.IsInt64() || q.Int64() > int64(dimension) {
		return MomentCountResult{}, fmt.Errorf("moment premises are infeasible for the stated dimension")
	}
	return MomentCountResult{Applicable: true, TraceResidual: residualExact, RealLowerBound: xExact, IntegerLowerBound: int(q.Int64())}, nil
}

type FiniteMomentCountTheorem struct {
	Name              string      `json:"name"`
	Assumptions       []string    `json:"assumptions"`
	Partition         string      `json:"threshold_partition"`
	TraceResidual     string      `json:"trace_residual"`
	CauchySchwarz     string      `json:"cauchy_schwarz"`
	RealConclusion    string      `json:"real_conclusion"`
	IntegerConclusion string      `json:"integer_conclusion"`
	Theorems          []TheoremID `json:"theorems"`
	Provenance        Reference   `json:"provenance"`
}

func (t FiniteMomentCountTheorem) Validate() error {
	fields := []string{t.Name, t.Partition, t.TraceResidual, t.CauchySchwarz, t.RealConclusion, t.IntegerConclusion}
	for _, f := range fields {
		if strings.TrimSpace(f) == "" {
			return fmt.Errorf("incomplete finite moment theorem")
		}
	}
	if len(t.Assumptions) == 0 || len(t.Theorems) == 0 || t.Provenance.Citation == "" {
		return fmt.Errorf("unprovenanced finite moment theorem")
	}
	return nil
}
