package semantic

import "testing"

func TestM7PrincipalMinorCertificateRequiresPositiveCertifiedBounds(t *testing.T) {
	good := TwoByTwoPrincipalMinorCertificate{MatrixID: "G", ALower: "0.1", DLower: "0.2", DeterminantLower: "0.01", AllStrictlyPositive: true, ArithmeticProofKinds: []ProofObjectKind{OutwardRoundedArithmetic}}
	if err := good.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := good
	bad.DeterminantLower = "-0.01"
	if bad.Validate() == nil {
		t.Fatal("negative lower bound accepted")
	}
	bad = good
	bad.ArithmeticProofKinds = []ProofObjectKind{"heuristic"}
	if bad.Validate() == nil {
		t.Fatal("heuristic arithmetic accepted")
	}
}
