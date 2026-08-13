//go:build slow

package compiler

import "testing"

// TestM18SlowWholeLineProofReplay runs the production M17 upper/witness replay
// followed by the equality-aware exact M18 endpoint verifier.
func TestM18SlowWholeLineProofReplay(t *testing.T) {
	if _, err := CompileM18(); err != nil {
		t.Fatal(err)
	}
}
