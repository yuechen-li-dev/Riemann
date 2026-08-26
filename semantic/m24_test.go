package semantic

import (
	"strings"
	"sync"
	"testing"
)

var m24CertificateOnce sync.Once
var m24CertificateCached ComponentOptimumCertificate
var m24CertificateErr error

func testM24Certificate() (ComponentOptimumCertificate, error) {
	m24CertificateOnce.Do(func() {
		m24CertificateCached, m24CertificateErr = CertifyM24ComponentControl()
	})
	return m24CertificateCached, m24CertificateErr
}

func TestM24CertifiesAllButOneExactRadiusWindow(t *testing.T) {
	certificate, err := testM24Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(certificate.Status, "Success C") || certificate.M23ControlsComponent {
		t.Fatalf("wrong M24 status: %+v", certificate)
	}
	if certificate.GlobalEvidence.UnresolvedRadius.Lower != "3.5030000000000000" || certificate.GlobalEvidence.UnresolvedRadius.Upper != "3.5120000000000000" {
		t.Fatalf("wrong unresolved radius window: %+v", certificate.GlobalEvidence.UnresolvedRadius)
	}
	if certificate.GlobalEvidence.GeneratedCellCount != 611 || certificate.MissingTheorem == "" {
		t.Fatalf("incomplete M24 obstruction: %+v", certificate.GlobalEvidence)
	}
}

func TestM24LegalComponentDefinitionDoesNotOverclaimEndpoints(t *testing.T) {
	certificate, err := testM24Certificate()
	if err != nil {
		t.Fatal(err)
	}
	component := certificate.Component
	for _, required := range []string{"connected component", "m_*(r)", "w1", "unresolved"} {
		joined := component.Definition + component.WeightLegality + component.ExactEndpointStatus
		if !strings.Contains(joined, required) {
			t.Fatalf("legal component metadata misses %q: %+v", required, component)
		}
	}
	if strings.Contains(component.Definition, "all r>1") {
		t.Fatal("M24 promoted the implicit legal component to all radii")
	}
}

func TestM24GeneratedRadiusCoverageHasOnlyDeclaredGap(t *testing.T) {
	if len(m24RadiusBarriers) != 611 {
		t.Fatalf("generated row count = %d", len(m24RadiusBarriers))
	}
	for index, cell := range m24RadiusBarriers {
		want := index
		if index >= 503 {
			want += 9
		}
		if cell.K != want {
			t.Fatalf("row %d has K=%d, want %d", index, cell.K, want)
		}
		if cell.Contact == "" {
			t.Fatalf("row %d has no Oct-generated contact", index)
		}
	}
}

func TestM24EvidenceSeparationAndM23Regression(t *testing.T) {
	certificate, err := testM24Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if !certificate.Optimizer.ExistsUnique || certificate.Optimizer.Radius.Lower != "3.5078500000000000" {
		t.Fatalf("M23 optimizer regression: %+v", certificate.Optimizer)
	}
	if certificate.OptimizerGloballyFeasible || !strings.Contains(certificate.GlobalEvidence.EvidenceBoundary, "research evidence only") {
		t.Fatal("numerical branch evidence leaked into theorem status")
	}
	if !strings.Contains(certificate.GlobalEvidence.FixedContactFormula, "r^2-1+8/pi^2") {
		t.Fatal("large-radius analytic barrier is missing")
	}
}
