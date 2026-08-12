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
	var err error
	if *missingPremise {
		result, compileErr := compiler.CompileM3WithOptions(compiler.M3Options{OmitCompletionFactor: true})
		if compileErr != nil {
			fail(compileErr)
		}
		if *jsonOutput {
			output, reportErr := compiler.M3JSONReport(result)
			if reportErr != nil {
				fail(reportErr)
			}
			_, _ = os.Stdout.Write(output)
			return
		}
		fmt.Print(compiler.M3HumanReport(result))
		return
	}
	result, err := compiler.CompileM14A()
	if err != nil {
		fail(err)
	}
	if *jsonOutput {
		output, err := compiler.M14AJSONReport(result)
		if err != nil {
			fail(err)
		}
		_, _ = os.Stdout.Write(output)
		return
	}
	fmt.Print(compiler.M14AHumanReport(result))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "riemann:", err)
	os.Exit(1)
}
