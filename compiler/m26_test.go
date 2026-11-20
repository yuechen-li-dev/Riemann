package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestM26ReleasesOnlyForcedContactAndImproves(t *testing.T) {
	r, err := testM26()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Witness.WholeLine || !strings.Contains(r.CentralAnswer, "yes") || !strings.Contains(r.Support, "+-1") || strings.Contains(r.Support, "r1") {
		t.Fatalf("wrong M26 scope or result: %+v", r)
	}
	if !strings.Contains(r.BestCertifiedJ, "0.8432725285110000") || !strings.Contains(r.BestCertifiedJ, "J=2-c") {
		t.Fatalf("stale J propagation: %s", r.BestCertifiedJ)
	}
}

func TestM26ReportsAreDeterministicAndArtifactMatches(t *testing.T) {
	r, err := testM26()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M26JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := M26JSONReport(r)
	if !bytes.Equal(a, b) {
		t.Fatal("M26 JSON report is nondeterministic")
	}
	for _, want := range []string{"CENTRAL ANSWER", "MINIMAL RELEASE / EMBEDDING", "CERTIFIED STRICT WITNESS", "ACTIVE-SET DISCOVERY", "THEOREM BOUNDARY"} {
		if !strings.Contains(M26HumanReport(r), want) {
			t.Fatalf("M26 human report misses %q", want)
		}
	}
	data, err := os.ReadFile("../" + r.Experiment.Artifact)
	if err != nil {
		t.Fatal(err)
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != r.Experiment.ArtifactSHA256 {
		t.Fatalf("M26 seed artifact changed: %s", got)
	}
}
