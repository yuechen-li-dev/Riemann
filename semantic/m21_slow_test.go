//go:build slow

package semantic

import "testing"

func TestM21SlowWholeLineProofReplay(t *testing.T) {
	r, err := DeriveContactAdaptedAtomicRay(ExactRational{19, 5000})
	if err != nil {
		t.Fatal(err)
	}
	if err := VerifyContactAdaptedRayWholeLine(r, validM21Certificate(), true); err != nil {
		t.Fatal(err)
	}
}
