//go:build slow

package compiler

import "testing"

// TestM17SlowWholeLineProofReplay is the single intentional replay of the
// 20,001-cell exact rational enclosure. Run it with:
//
//	go test -tags=slow ./compiler -run M17Slow -count=1
//
// Normal unit and race suites exercise the proof object's structure, exact
// coefficients, rejection paths, and result propagation without this replay.
func TestM17SlowWholeLineProofReplay(t *testing.T) {
	if _, err := CompileM17(); err != nil {
		t.Fatal(err)
	}
}
