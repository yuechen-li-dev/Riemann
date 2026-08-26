package compiler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM21CertifiesSelectedRayWithoutWidening(t *testing.T) {
	r, err := testM21()
	if err != nil {
		t.Fatal(err)
	}
	if r.CertifiedWitness != (semantic.ExactRational{Numerator: 19, Denominator: 5000}) || r.Ray.Radius2 != (semantic.ExactRational{Numerator: 4, Denominator: 1}) || r.Ray.Weight2 != r.CertifiedWitness {
		t.Fatalf("wrong M21 ray witness: %+v", r.Ray)
	}
	if !strings.Contains(r.CertifiedBounds, ">1+pi/4-2/pi") || r.RHStatus != "unresolved" || !strings.Contains(r.Outcome, "Success A") {
		t.Fatalf("M21 did not propagate the certified improvement: %+v", r)
	}
	if r.M20.M18.Optimum.ExactValue.Expression != "1+pi/4-2/pi" || r.M20.SelectedFamily != "D: contact-adapted r=4 atomic ray" {
		t.Fatal("M21 changed the M18-M20 theorem chain")
	}
}

func TestM21ContactEnvelopeAndAboveBoundRejection(t *testing.T) {
	r, err := testM21()
	if err != nil {
		t.Fatal(err)
	}
	a := math.Pi / 2
	if math.Abs(semantic.RayResponse(a)) > 1e-14 {
		t.Fatal("contact response does not vanish")
	}
	above, err := semantic.DeriveContactAdaptedAtomicRay(r.Ceiling.UpperBound)
	if err != nil {
		t.Fatal(err)
	}
	tw := float64(r.Ceiling.UpperWitnessFrequency.Numerator) / float64(r.Ceiling.UpperWitnessFrequency.Denominator)
	if semantic.ContactAdaptedRaySpectrum(above, tw) >= 0 {
		t.Fatal("certified above-bound counterfrequency is not negative numerically")
	}
	if r.NumericalCandidate.Curvature <= 0 || r.NumericalCandidate.Frequency >= a {
		t.Fatal("numerical second contact has wrong order or location")
	}
}

func TestM21ReportsAndWitnessAreDeterministic(t *testing.T) {
	r, err := testM21()
	if err != nil {
		t.Fatal(err)
	}
	a, err := M21JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	b, err := M21JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("M21 JSON report is nondeterministic")
	}
	h := M21HumanReport(r)
	for _, want := range []string{"M21 CONTACT-ADAPTED RAY", "FIXED CONTACT", "RAY RESPONSE", "ENVELOPE", "CERTIFIED WITNESS", "WHOLE-LINE PROOF", "UPDATED EP3.1 BOUND", "RAY CEILING", "RH\n  unresolved", "ONE PROPOSED NEXT MILESTONE"} {
		if !strings.Contains(h, want) {
			t.Fatalf("M21 human report missing %q", want)
		}
	}
	data, err := os.ReadFile("artifacts/m21_contact_adapted_ray_witness.json")
	if err != nil || !bytes.Equal(bytes.TrimSpace(data), bytes.TrimSpace(m21WitnessJSON)) {
		t.Fatalf("embedded M21 witness differs from disk: %v", err)
	}
}

func TestM21OctArtifactsAreDeterministic(t *testing.T) {
	wants := map[string]string{
		"../artifacts/m21/m21_m18_baseline_contact_zoom.png":    "a87b0171e1619fd78ee4369c17cec87c34543b2bf1b0eec3bc107135f8465537",
		"../artifacts/m21/m21_contact_adapted_ray_pressure.png": "660f3defb5c68ac50c495c25352089044841aafa03b40d98e2957ed669309888",
		"../artifacts/m21/m21_contact_adapted_ray_envelope.png": "bdbc8a645e4be17ef405182876a4b2e03190ac50c94a9d0db79b66671315dc13",
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
