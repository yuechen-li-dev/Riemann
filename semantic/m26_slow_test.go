//go:build slow

package semantic

import "testing"

func TestM26SlowExactWholeLineReplay(t *testing.T) {
	w, err := ReplayM26ReleasedWitness()
	if err != nil {
		t.Fatal(err)
	}
	if !w.WholeLine {
		t.Fatal("M26 exact replay lost whole-line coverage")
	}
}
