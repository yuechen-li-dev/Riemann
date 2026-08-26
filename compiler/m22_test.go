package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM22CertifiesRadiusReleaseWithoutWideningSupport(t *testing.T) {
	r, err := testM22()
	if err != nil {
		t.Fatal(err)
	}
	if r.Family.Radius != (semantic.ExactRational{Numerator: 7, Denominator: 2}) || r.Family.Mass != (semantic.ExactRational{Numerator: 17, Denominator: 2500}) || r.Family.Weight2 != r.Family.Mass {
		t.Fatalf("wrong M22 witness: %+v", r.Family)
	}
	if !strings.Contains(r.Outcome, "Success A") || !strings.Contains(r.CertifiedBounds, "17*pi*sqrt(2)/10000") || r.RHStatus != "unresolved" {
		t.Fatalf("M22 theorem did not propagate: %+v", r)
	}
	if r.M21.Ray.Radius2 != (semantic.ExactRational{Numerator: 4, Denominator: 1}) || r.NumericalCandidate.C <= r.M21.NumericalCandidate.C {
		t.Fatal("M22 lost its M21 regression or numerical improvement")
	}
}

func TestM22OctArtifactsAreDeterministic(t *testing.T) {
	wants := map[string]string{
		"../artifacts/m22/m22_radius_response.png": "8cdcd23c4606ce248398f0eb9b00f6b5b4619325677aef093379430c91d49f35",
		"../artifacts/m22/m22_contact_motion.png":  "044871b4287c8e62dca49a6ca61cddcc2984a04222d6cd7f42da8bbac330c6d0",
		"../artifacts/m22/m22_best_spectrum.png":   "364a850d58f57f4b8ccc680a53435a5db657a2946ec6eea7651385f5a52ebfcb",
	}
	for path, want := range wants {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if got := fmt.Sprintf("%x", sha256.Sum256(data)); got != want {
			t.Fatalf("artifact %s changed: %s", path, got)
		}
	}
}

func TestM22EvidenceTypesAndActiveSystem(t *testing.T) {
	r, err := testM22()
	if err != nil {
		t.Fatal(err)
	}
	if r.Sensitivity.DeltaFAtOldContact >= 0 || r.Sensitivity.DeltaCPerRadiusAtFixedMass != 0 {
		t.Fatalf("wrong M21 radius sensitivity: %+v", r.Sensitivity)
	}
	if len(r.ActiveContactSystem) != 3 || !strings.Contains(r.ActiveContactSystem[2], "Q3") || !strings.Contains(r.ThirdConstraint, "No third near-active") {
		t.Fatal("finite optimizer or third-constraint evidence is incomplete")
	}
	if r.CertifiedWitness.Family == (semantic.ContactPreservingTwoAtomFamily{}) || r.NumericalCandidate.Evidence == "" {
		t.Fatal("numerical and certified evidence were not stored separately")
	}
	if r.NumericalContactFrequency != r.NumericalCandidate.ContactFrequency || !strings.Contains(r.CertifiedFamilyBracket, "not claimed") {
		t.Fatal("candidate frequency and certified-family status were not separated")
	}
}

func TestM22ReportsAndWitnessAreDeterministic(t *testing.T) {
	r, err := testM22()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M22JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M22JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M22 JSON report is nondeterministic")
	}
	h := M22HumanReport(r)
	for _, want := range []string{"M22 CONTACT-PRESERVING TWO-ATOM FAMILY", "EXACT FAMILY", "M21 REGRESSION", "FIRST-ORDER CONTACT MOTION", "NESTED ENVELOPE", "CERTIFIED WITNESS", "WHOLE-LINE PROOF", "UPDATED EP3.1 BOUND", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("M22 report missing %q", want)
		}
	}
	data, err := os.ReadFile("artifacts/m22_contact_preserving_two_atom_witness.json")
	if err != nil || !bytes.Equal(bytes.TrimSpace(data), bytes.TrimSpace(m22WitnessJSON)) {
		t.Fatalf("embedded M22 witness differs from disk: %v", err)
	}
}
