package claudegrh

import (
	"encoding/json"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestSelectionCertificates(t *testing.T) {
	b, err := os.ReadFile("testdata/selection.json")
	if err != nil {
		t.Fatal(err)
	}
	var boxes []SelectionBox
	if err = json.Unmarshal(b, &boxes); err != nil || len(boxes) != 2 {
		t.Fatal(err)
	}
	for i, box := range boxes {
		got, err := ReplaySelection(box)
		if err != nil || got.Outcome != ReducedCounterexample || !got.SeparateAtEveryPlace || got.JointAtEveryPlace != (i == 1) || len(got.Common) != 0 {
			t.Fatal(got, err)
		}
		again, _ := ReplaySelection(box)
		if i == 1 && !reflect.DeepEqual(got.PlaceChoices, []string{"A", "B"}) {
			t.Fatal(got)
		}
		if !reflect.DeepEqual(again, got) {
			t.Fatal("unstable replay")
		}
		for _, c := range got.Checks {
			if (c.FirstFailure == "") != (c.Comparison && c.System) {
				t.Fatal(c)
			}
		}
		t.Logf("%s %+v", got.Outcome, got.Checks)
	}
	boxes[1].Candidates = append(boxes[1].Candidates, ChainCandidate{"C", []string{"6"}, [][]string{{"1"}, {"1"}}})
	r, e := ReplaySelection(boxes[1])
	if e != nil || r.Outcome != FiniteVerified || !reflect.DeepEqual(r.Common, []string{"C"}) {
		t.Fatal(r, e)
	}
}

func TestSelectionLossAndDegeneracy(t *testing.T) {
	box := SelectionBox{[]uint32{2}, "4", []string{"1"}, []int64{1}, []ChainCandidate{{"A", []string{"2", "2"}, [][]string{{"1", "1"}}}}}
	check := func(want Outcome) {
		t.Helper()
		r, e := ReplaySelection(box)
		if e != nil || r.Outcome != want {
			t.Fatal(r, e)
		}
	}
	check(FiniteVerified)
	box.Low[0] = "1000000000000000001/1000000000000000000"
	check(Inconclusive)
	box.Candidates[0].Bounds[0][0] = "2"
	check(Inconclusive) // total suffices, first stage does not
	box.Column = "0"
	box.Candidates[0].Scalars[0] = "0"
	check(FiniteVerified)
	box.Candidates = nil
	check(Inconclusive)
	box.Primes[0] = 4
	if _, e := ReplaySelection(box); e == nil {
		t.Fatal("accepted composite")
	}
}

func TestJointWitnessGraph(t *testing.T) {
	g, e := SimultaneousGraph()
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []statement{commonLC, columnBound, gainBound, periods, dual, rh} {
		if ok, _ := g.Certify(s.id()); ok {
			t.Fatal(s)
		}
	}
	for _, args := range [][]statement{{cmpOnly, sysOnly}, {cmpOnly, serviceData}, {sysOnly, serviceData}, {cmpOnly, sysOnly, serviceData}} {
		if e := selectionRules[3].apply(g, args...); e == nil || !strings.Contains(e.Error(), "requires") {
			t.Fatal("accepted separate existential witness", args)
		}
	}
	for _, s := range []statement{cmpOnly, sysOnly} {
		if g.AttemptProof(s.id(), commonLC.id()).Accepted {
			t.Fatal(s)
		}
	}
	anc, e := g.Lineage(rh.id())
	if e != nil {
		t.Fatal(e)
	}
	for _, s := range []statement{commonLC, columnBound, gainBound, offPeriod} {
		found := false
		for _, a := range anc {
			found = found || a == s.id()
		}
		if !found {
			t.Fatal(s)
		}
	}
}
