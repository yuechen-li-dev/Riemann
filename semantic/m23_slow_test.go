//go:build slow

package semantic

import "testing"

func TestM23SlowWholeLineProofReplay(t *testing.T) {
	if err := VerifyM23ExplicitWitness(M23ExplicitWitness(), true); err != nil {
		t.Fatal(err)
	}
}
