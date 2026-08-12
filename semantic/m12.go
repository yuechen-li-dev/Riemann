package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

// HermitianComponentDecomposition is deliberately an identity, rather than
// three matrix labels passed independently to a theorem.  M12 needs the fact
// that the moments of Total belong to the same object assembled from P and Q.
type HermitianComponentDecomposition struct {
	TotalMatrixID   string               `json:"total_matrix_id"`
	PSDMatrixID     string               `json:"psd_component_id"`
	IndexMatrixID   string               `json:"index_component_id"`
	Identity        string               `json:"identity"`
	SameDimension   bool                 `json:"same_dimension"`
	TotalHermitian  bool                 `json:"total_hermitian"`
	PSDHermitian    bool                 `json:"psd_component_hermitian"`
	QHermitian      bool                 `json:"index_component_hermitian"`
	PSD             bool                 `json:"psd_component_positive_semidefinite"`
	Evidence        SpectralEvidenceKind `json:"evidence"`
	IdentityTheorem TheoremID            `json:"identity_theorem"`
	Provenance      Reference            `json:"provenance"`
}

func (d HermitianComponentDecomposition) Validate() error {
	if d.TotalMatrixID == "" || d.PSDMatrixID == "" || d.IndexMatrixID == "" ||
		d.Identity == "" || !d.SameDimension || !d.TotalHermitian || !d.PSDHermitian ||
		!d.QHermitian || !d.PSD || d.IdentityTheorem == "" || d.Provenance.Citation == "" {
		return fmt.Errorf("invalid Hermitian component decomposition")
	}
	if d.Evidence != ExactTheoremEvidence && d.Evidence != CertifiedMinorEvidence {
		return fmt.Errorf("decomposition is not theorem-grade")
	}
	return nil
}

type PositiveScalarParameter struct {
	Symbol string        `json:"symbol"`
	Value  ExactRational `json:"exact_value"`
	Domain string        `json:"domain"`
}

func (p PositiveScalarParameter) Validate() error {
	if p.Symbol == "" || p.Domain != "c>0" {
		return fmt.Errorf("rank-trace parameter must have domain c>0")
	}
	if err := p.Value.Validate(); err != nil {
		return err
	}
	if rat(p.Value).Sign() <= 0 {
		return fmt.Errorf("rank-trace parameter must be strictly positive")
	}
	return nil
}

type FiniteRankTraceTheorem struct {
	Name            string      `json:"name"`
	Assumptions     []string    `json:"assumptions"`
	Expansion       string      `json:"frobenius_expansion"`
	VonNeumannStep  string      `json:"von_neumann_step"`
	ScalarSteps     []string    `json:"scalar_steps"`
	Conclusion      string      `json:"conclusion"`
	Specialization  string      `json:"c_2_specialization"`
	EqualityCase    string      `json:"equality_case"`
	Theorems        []TheoremID `json:"theorems"`
	CompilerDerived bool        `json:"compiler_derived_consequence"`
	Provenance      Reference   `json:"provenance"`
}

func (t FiniteRankTraceTheorem) Validate() error {
	fields := []string{t.Name, t.Expansion, t.VonNeumannStep, t.Conclusion, t.Specialization, t.EqualityCase}
	for _, field := range fields {
		if strings.TrimSpace(field) == "" {
			return fmt.Errorf("incomplete finite rank-trace theorem")
		}
	}
	if len(t.Assumptions) < 5 || len(t.ScalarSteps) < 2 || len(t.Theorems) < 2 || !t.CompilerDerived || t.Provenance.Citation == "" {
		return fmt.Errorf("unprovenanced finite rank-trace theorem")
	}
	return nil
}

type RankTraceExactResult struct {
	Applicable       bool          `json:"applicable"`
	FrobeniusSquared ExactRational `json:"frobenius_squared"`
	RightHandSide    ExactRational `json:"right_hand_side"`
	Slack            ExactRational `json:"slack"`
}

// CheckDiagonalRankTrace is an exact-rational regression oracle for diagonal
// fixtures.  It attacks coefficients and premises; it is not the evidence for
// the general theorem, whose noncommuting step is von Neumann's inequality.
func CheckDiagonalRankTrace(p, q []ExactRational, rankBudget, positiveIndexBudget int, c ExactRational) (RankTraceExactResult, error) {
	if len(p) == 0 || len(p) != len(q) || rankBudget < 0 || positiveIndexBudget < 0 || rankBudget > len(p) || positiveIndexBudget > len(p) {
		return RankTraceExactResult{}, fmt.Errorf("invalid diagonal rank-trace fixture")
	}
	if err := (PositiveScalarParameter{Symbol: "c", Value: c, Domain: "c>0"}).Validate(); err != nil {
		return RankTraceExactResult{}, err
	}
	trP, trQ, frob := new(big.Rat), new(big.Rat), new(big.Rat)
	rankP, nplusQ := 0, 0
	for i := range p {
		if err := p[i].Validate(); err != nil {
			return RankTraceExactResult{}, err
		}
		if err := q[i].Validate(); err != nil {
			return RankTraceExactResult{}, err
		}
		pv, qv := rat(p[i]), rat(q[i])
		if pv.Sign() < 0 {
			return RankTraceExactResult{}, fmt.Errorf("P is not positive semidefinite")
		}
		if pv.Sign() != 0 {
			rankP++
		}
		if qv.Sign() > 0 {
			nplusQ++
		}
		trP.Add(trP, pv)
		trQ.Add(trQ, qv)
		s := new(big.Rat).Add(pv, qv)
		frob.Add(frob, new(big.Rat).Mul(s, s))
	}
	if rankP > rankBudget {
		return RankTraceExactResult{}, fmt.Errorf("rank(P) exceeds rank budget")
	}
	if nplusQ > positiveIndexBudget {
		return RankTraceExactResult{}, fmt.Errorf("n_plus(Q) exceeds positive-index budget")
	}
	cRat := rat(c)
	c2 := new(big.Rat).Mul(cRat, cRat)
	rhs := new(big.Rat).Mul(cRat, trP)
	rhs.Sub(rhs, new(big.Rat).Mul(new(big.Rat).Quo(c2, big.NewRat(4, 1)), big.NewRat(int64(rankBudget), 1)))
	rhs.Add(rhs, new(big.Rat).Mul(new(big.Rat).Mul(big.NewRat(2, 1), cRat), trQ))
	rhs.Sub(rhs, new(big.Rat).Mul(c2, big.NewRat(int64(positiveIndexBudget), 1)))
	slack := new(big.Rat).Sub(frob, rhs)
	frobExact, err := exactRat(frob)
	if err != nil {
		return RankTraceExactResult{}, err
	}
	rhsExact, err := exactRat(rhs)
	if err != nil {
		return RankTraceExactResult{}, err
	}
	slackExact, err := exactRat(slack)
	if err != nil {
		return RankTraceExactResult{}, err
	}
	return RankTraceExactResult{Applicable: true, FrobeniusSquared: frobExact, RightHandSide: rhsExact, Slack: slackExact}, nil
}
