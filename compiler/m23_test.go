package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestM23OctArtifactsAreDeterministic(t *testing.T) {
	wants := map[string]string{
		"../artifacts/m23/m23_radius_response.png":    "0d4abaa3b1488ce9511066e5e5aa2ccd58964819e51b352a9a0189f1db53c5c7",
		"../artifacts/m23/m23_contact_branch.png":     "e4521aef10b9ca96323cdc8373e6c9a0a2028dfe53692f853ca6ff4208bf2935",
		"../artifacts/m23/m23_optimizer_spectrum.png": "34c9d21307c32f05ba431139bf30c91d6ea782d0b9a7b50ea931fe5f897926bd",
	}
	for p, w := range wants {
		data, err := os.ReadFile(p)
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != w {
			t.Fatalf("artifact %s changed: %s", p, got)
		}
	}
}

func TestM23CertifiesLocalOptimizerAndSeparateWitness(t *testing.T) {
	r, err := testM23()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(r.Outcome, "Success A") || !r.Optimizer.ExistsUnique || !r.LocalOptimality.CertifiedLocalMaximum {
		t.Fatalf("local optimizer missing: %+v", r)
	}
	if !r.ExplicitWitness.ImprovesM22 || r.ExplicitWitness.Witness.Family.Mass.Numerator != 69 || r.ExplicitWitness.Witness.Family.Mass.Denominator != 10000 {
		t.Fatal("global witness not separated or improved")
	}
	if !strings.Contains(r.CertifiedBounds, "69*pi*sqrt(2)/40000") || r.RHStatus != "unresolved" {
		t.Fatal("wrong theorem propagation")
	}
}

func TestM23BoxWidthsJacobianCurvaturesAndTail(t *testing.T) {
	r, err := testM23()
	if err != nil {
		t.Fatal(err)
	}
	if r.Optimizer.JacobianDeterminant.Lower <= "0" || r.LocalOptimality.MovableCurvature.Lower <= "0" || r.LocalOptimality.FixedContactCurvature.Lower <= "0" || r.LocalOptimality.TailLower.Lower <= "0" {
		t.Fatalf("missing positive interval: %+v", r.LocalOptimality)
	}
	if !strings.HasPrefix(r.LocalOptimality.EnvelopeSecond.Upper, "-") {
		t.Fatal("radius local maximum sign not certified")
	}
}

func TestM23ReportsAreDeterministicAndComplete(t *testing.T) {
	r, err := testM23()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M23JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M23JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M23 JSON is nondeterministic")
	}
	h := M23HumanReport(r)
	for _, want := range []string{"Q SYSTEM", "CERTIFIED OPTIMIZER BOX", "LOCAL OPTIMALITY", "EXPLICIT GLOBAL WITNESS", "UPDATED EP3.1 BOUNDS", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("report missing %q", want)
		}
	}
}
