// Command m25replay performs the expensive exact reconstruction behind the
// persisted M25 proof manifest. Ordinary tests intentionally validate only the
// deterministic manifest and its composition contracts.
package main

import (
	"fmt"
	"os"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func main() {
	c, err := semantic.ReplayM25B1Control()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("M25 exact replay: slabs=%d B1=%t whole_line=%t component=%t\n", len(c.Atlas.Slabs), c.Branch.ExistsUnique, c.Control.OptimizerFeasible, c.Control.M23ControlsComponent)
}
