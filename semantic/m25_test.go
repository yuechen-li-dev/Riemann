package semantic

import (
	"strings"
	"sync"
	"testing"
)

var m25Once sync.Once
var m25Cached M25Certificate
var m25CachedErr error

func testM25Certificate() (M25Certificate, error) {
	m25Once.Do(func() { m25Cached, m25CachedErr = CertifyM25B1Control() })
	return m25Cached, m25CachedErr
}

func TestM25B1ContinuationCoversWindowUniquely(t *testing.T) {
	c, err := testM25Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Atlas.Slabs) != 90 || !c.Atlas.NoRadiusGaps || !c.Branch.ExistsUnique || !c.Branch.Smooth {
		t.Fatalf("incomplete B1 continuation: %+v", c.Branch)
	}
	if c.Atlas.Radius.Lower != "3.5030000000000000" || c.Atlas.Radius.Upper != "3.5120000000000000" {
		t.Fatalf("wrong M25 radius cover: %+v", c.Atlas.Radius)
	}
}

func TestM25B1CurvatureDerivativeAndM23Regression(t *testing.T) {
	c, err := testM25Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if !c.Branch.StrictMinimum || c.Branch.Curvature.Lower <= "0" || !strings.Contains(c.Branch.DerivativeSign, "C1'(r)>0") {
		t.Fatalf("lost B1 curvature or derivative structure: %+v", c.Branch)
	}
	if c.Optimum.Radius.Lower != "3.5078500000000000" || !c.Optimum.Unique {
		t.Fatalf("M23 optimizer regression: %+v", c.Optimum)
	}
}

func TestM25WholeLineCompetitorExclusionAndComposition(t *testing.T) {
	c, err := testM25Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if !c.Control.B1ControlsWindow || !c.Control.NoBranchSwitch || !c.Control.OptimizerFeasible || !c.Control.M23ControlsComponent {
		t.Fatalf("M25 did not close global component control: %+v", c.Control)
	}
	if len(c.Control.CompetingBands) != 6 || !strings.Contains(c.Control.BoundaryComposition, "no gap") {
		t.Fatalf("incomplete competitor/coverage manifest: %+v", c.Control)
	}
}

func TestM25GeneratedAtlasDeterminismAndEvidenceSeparation(t *testing.T) {
	if len(m25ContinuationAtlas) != 90 {
		t.Fatalf("generated slab count=%d", len(m25ContinuationAtlas))
	}
	for i, slab := range m25ContinuationAtlas {
		if slab.K != i || slab.ContactLo == "" || slab.MassLo == "" {
			t.Fatalf("bad generated slab %d: %+v", i, slab)
		}
	}
	c, err := testM25Certificate()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(c.Atlas.Method, "interval") || strings.Contains(c.Atlas.Method, "numerical proof") {
		t.Fatal("Oct discovery leaked into theorem authority")
	}
}
