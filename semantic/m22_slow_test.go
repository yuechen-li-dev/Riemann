//go:build slow

package semantic

import "testing"

func TestM22SlowWholeLineProofReplay(t *testing.T) {
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{17, 2500}, ExactRational{7, 2})
	if err != nil {
		t.Fatal(err)
	}
	w := CertifiedTwoAtomWitness{Family: f, WholeLine: validM22Certificate(), CExact: "1+pi/4-2/pi+17*pi*sqrt(2)/10000", JExact: "1-pi/4+2/pi-17*pi*sqrt(2)/10000", Verified: true}
	if err := VerifyCertifiedTwoAtomWitness(w, true); err != nil {
		t.Fatal(err)
	}
}
