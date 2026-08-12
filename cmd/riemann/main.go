package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yuechen-li-dev/Riemann/compiler"
)

func main() {
	jsonOutput := flag.Bool("json", false, "emit the deterministic machine-readable proof graph")
	missingPremise := flag.Bool("missing-premise", false, "demonstrate an unresolved M3 completion-factor side condition")
	flag.Parse()
	var result compiler.M3Result
	var err error
	if *missingPremise {
		result, err = compiler.CompileM3WithOptions(compiler.M3Options{OmitCompletionFactor: true})
	} else {
		result, err = compiler.CompileM3()
	}
	if err != nil {
		fail(err)
	}
	if *jsonOutput {
		output, err := compiler.M3JSONReport(result)
		if err != nil {
			fail(err)
		}
		_, _ = os.Stdout.Write(output)
		return
	}
	fmt.Print(compiler.M3HumanReport(result))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "riemann:", err)
	os.Exit(1)
}
