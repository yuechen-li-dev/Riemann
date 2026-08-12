package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yuechen-li-dev/Riemann/compiler"
)

func main() {
	jsonOutput := flag.Bool("json", false, "emit the deterministic machine-readable proof graph")
	flag.Parse()
	result, err := compiler.CompileM1()
	if err != nil {
		fail(err)
	}
	if *jsonOutput {
		output, err := compiler.M1JSONReport(result)
		if err != nil {
			fail(err)
		}
		_, _ = os.Stdout.Write(output)
		return
	}
	fmt.Print(compiler.M1HumanReport(result))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "riemann:", err)
	os.Exit(1)
}
