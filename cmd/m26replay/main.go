// Command m26replay explicitly reconstructs the persisted M26 whole-line
// proof. Ordinary tests validate the immutable manifest instead.
package main

import (
	"fmt"
	"os"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func main() {
	w, err := semantic.ReplayM26ReleasedWitness()
	if err != nil {
		fmt.Fprintln(os.Stderr, "m26 replay:", err)
		os.Exit(1)
	}
	fmt.Printf("M26 exact replay: whole_line=%t c=[%s,%s] compact_lower=%s tail_lower=%s\n", w.WholeLine, w.C.Lower, w.C.Upper, w.CompactLower.Lower, w.TailLower.Lower)
}
