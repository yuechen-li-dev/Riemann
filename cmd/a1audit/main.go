// Command a1audit replays a rational lattice instance or prints the A1 proof graph.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yuechen-li-dev/Riemann/audit/claudegrh"
	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

func run() error {
	path := flag.String("replay", "", "replay a JSON instance (default: discover the first averaging failure)")
	graph := flag.Bool("graph", false, "show pinned proof dependencies and their evidence status")
	flag.Parse()
	out := json.NewEncoder(os.Stdout)
	out.SetIndent("", "  ")
	if flag.NArg() != 0 || (*graph && *path != "") {
		return fmt.Errorf("choose discovery, -replay FILE, or -graph")
	}
	if *graph {
		g, err := claudegrh.ProofGraph()
		if err != nil {
			return err
		}
		return out.Encode(struct {
			Claims []semantic.Claim
			Uses   []compiler.Transformation
		}{g.Claims(), g.Transformations()})
	}
	var in claudegrh.Instance
	if *path == "" {
		var err error
		in, err = claudegrh.FirstAveragingFailure()
		if err != nil {
			return err
		}
	} else {
		f, err := os.Open(*path)
		if err != nil {
			return err
		}
		defer f.Close()
		dec := json.NewDecoder(f)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&in); err != nil {
			return err
		}
		if err := dec.Decode(new(any)); err != io.EOF {
			return fmt.Errorf("trailing data in %s", *path)
		}
	}
	result, err := claudegrh.Replay(in)
	if err != nil {
		return err
	}
	return out.Encode(struct {
		Instance claudegrh.Instance
		Result   claudegrh.Result
	}{in, result})
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
