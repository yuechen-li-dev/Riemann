package semantic

import "testing"

func orbitTestReference() Reference {
	return Reference{Kind: StandardReference, Citation: "test symmetry theorem"}
}

func TestM8OrbitCanonicalizationSeparatesMultiplicity(t *testing.T) {
	off, err := NewZeroOrbit("off", Point("rho"), OffCriticalOrbit, 3, []TheoremID{"symmetry"}, orbitTestReference())
	if err != nil {
		t.Fatal(err)
	}
	if off.DistinctLocationCount != 4 || off.ZeroMultiplicity != 3 {
		t.Fatalf("off orbit = %#v", off)
	}
	critical, err := NewZeroOrbit("critical", PointOnCriticalLine("rho"), CriticalLineOrbit, 2, []TheoremID{"symmetry"}, orbitTestReference())
	if err != nil {
		t.Fatal(err)
	}
	if critical.DistinctLocationCount != 2 || critical.ZeroMultiplicity != 2 {
		t.Fatalf("critical orbit = %#v", critical)
	}
	for _, p := range critical.TransformedPoints {
		if len(p.GeneratedBy) != 2 {
			t.Fatalf("critical degeneracy not retained: %#v", p)
		}
	}
}

func TestM8EvaluationVectorSemanticKeyIncludesBasisAndPoint(t *testing.T) {
	f := TestFunction{Symbol: "f", Kind: BasisTestFunction, DeclaredClass: WeilNiceClass().ID, RequiredAttributes: WeilNiceClass().Constraints, TransformConvention: LagariasMellinConvention}
	b1 := OrderedBasis{ID: "b1", Members: []BasisMember{{Function: f, AdmissibilityCertificate: "ok"}}, CoefficientField: ComplexField, ParentClass: WeilNiceClass()}
	b2 := b1
	b2.ID = "b2"
	v1, err := NewBasisEvaluationVector(b1, Point("rho"), "mellin", orbitTestReference())
	if err != nil {
		t.Fatal(err)
	}
	v2, err := NewBasisEvaluationVector(b2, Point("rho"), "mellin", orbitTestReference())
	if err != nil {
		t.Fatal(err)
	}
	v3, err := NewBasisEvaluationVector(b1, Point("sigma"), "mellin", orbitTestReference())
	if err != nil {
		t.Fatal(err)
	}
	if v1.Key() == v2.Key() || v1.Key() == v3.Key() {
		t.Fatal("basis or point identity was flattened out of vector key")
	}
}

func TestM8InvalidClassLocationRejected(t *testing.T) {
	if _, err := NewZeroOrbit("bad", PointOnCriticalLine("rho"), OffCriticalOrbit, 1, []TheoremID{"symmetry"}, orbitTestReference()); err == nil {
		t.Fatal("off-critical orbit accepted critical representative")
	}
}
