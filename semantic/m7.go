package semantic

import (
	"fmt"
	"math/big"
	"strings"
)

// TwoByTwoPrincipalMinorCertificate is the exact proposition established by
// the numerical backend. It is intentionally not itself a PSD proposition;
// the finite-dimensional theorem contract performs that final implication.
type TwoByTwoPrincipalMinorCertificate struct {
	MatrixID             string            `json:"matrix_id"`
	ALower               string            `json:"a_lower"`
	DLower               string            `json:"d_lower"`
	DeterminantLower     string            `json:"determinant_lower"`
	AllStrictlyPositive  bool              `json:"all_strictly_positive"`
	ArithmeticProofKinds []ProofObjectKind `json:"arithmetic_proof_kinds"`
}

func (TwoByTwoPrincipalMinorCertificate) isProposition()        {}
func (TwoByTwoPrincipalMinorCertificate) Kind() PropositionKind { return TwoByTwoMinorCertificateKind }
func (p TwoByTwoPrincipalMinorCertificate) Describe() string {
	return fmt.Sprintf("%s has certified strictly positive 2x2 principal minors", p.MatrixID)
}
func (p TwoByTwoPrincipalMinorCertificate) Validate() error {
	if p.MatrixID == "" || strings.TrimSpace(p.ALower) == "" || strings.TrimSpace(p.DLower) == "" || strings.TrimSpace(p.DeterminantLower) == "" || !p.AllStrictlyPositive || len(p.ArithmeticProofKinds) == 0 {
		return fmt.Errorf("invalid two-by-two principal-minor certificate")
	}
	for _, s := range []string{p.ALower, p.DLower, p.DeterminantLower} {
		x, _, err := big.ParseFloat(s, 10, 256, big.ToNegativeInf)
		if err != nil || x.Sign() <= 0 {
			return fmt.Errorf("principal-minor lower bound is not strictly positive")
		}
	}
	for _, k := range p.ArithmeticProofKinds {
		if !k.CertifiesInterval() {
			return fmt.Errorf("principal-minor certificate uses non-certifying arithmetic")
		}
	}
	return nil
}
