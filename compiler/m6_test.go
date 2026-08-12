package compiler

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM6AdmissibleLogBoxBasisAndExactMellinEndpoints(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Evaluation.Matrix.Basis.Members) != 2 {
		t.Fatal("wrong basis dimension")
	}
	contract, ok := r.Registry.Contract(WeilCrossEntryTheoremID)
	if !ok || contract.Trust != TrustedExternalTheorem || contract.Citation != "Lagarias §3, Theorem 3.1 and Remark 3.3" {
		t.Fatalf("missing entry-evaluation theorem contract: %+v", contract)
	}
	for _, member := range r.Evaluation.Matrix.Basis.Members {
		q, err := logBoxParameter(member.Function)
		if err != nil || q < 2 {
			t.Fatalf("basis is not evaluable log-box: %+v %v", member, err)
		}
		claim, ok := r.M5.Graph.Claim(member.AdmissibilityCertificate)
		if !ok {
			t.Fatal("missing admissibility claim")
		}
		p, ok := claim.Proposition.(semantic.TestFunctionAdmissibility)
		if !ok || p.Function.Key() != member.Function.Key() {
			t.Fatal("certificate does not identify basis function")
		}
		if certified, _ := r.M5.Graph.Certify(claim.ID); !certified {
			t.Fatal("basis admissibility not certified")
		}
	}
	for _, entry := range r.Evaluation.Matrix.Entries {
		if len(entry.TransformEvaluations) != 4 {
			t.Fatalf("missing Mellin records: %+v", entry.TransformEvaluations)
		}
		for _, m := range entry.TransformEvaluations {
			if m.Convention != semantic.LagariasMellinConvention || m.Value.Kind != semantic.ExactValue {
				t.Fatalf("bad Mellin evidence: %+v", m)
			}
		}
	}
}

func TestM6ContributionEvaluationAndVisibility(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range r.Evaluation.Matrix.Entries {
		if entry.Value.Kind != semantic.ApproximateValue || !entry.Value.DefinitionExact {
			t.Fatalf("entry state conflates definition and value: %+v", entry.Value)
		}
		if entry.Value.Metadata.SemanticTargetID == "" || entry.Value.Metadata.Backend == "" || entry.TransformConvention != semantic.LagariasMellinConvention || len(entry.TheoremProvenance) < 3 {
			t.Fatal("entry provenance incomplete")
		}
		for _, c := range entry.Contributions {
			switch c.SourceKind {
			case semantic.ZeroContribution:
				if c.Value.Kind != semantic.UnevaluatedExactDefinition {
					t.Fatal("zero-side alternate representation was flattened")
				}
			case semantic.EndpointContribution:
				if c.Value.Kind != semantic.ExactValue {
					t.Fatal("cheap exact endpoint was rounded")
				}
			case semantic.PrimePowerContribution:
				if c.Value.Kind != semantic.ApproximateValue || c.Value.Metadata.Truncation == nil || !c.Value.Metadata.Truncation.SupportExhaustive || c.Value.Metadata.Truncation.RemainderStatus != "exactly_zero_by_compact_support" {
					t.Fatalf("prime truncation hidden: %+v", c.Value)
				}
			case semantic.ArchimedeanContribution:
				if c.Value.Kind != semantic.ApproximateValue || c.Value.Metadata.Quadrature == nil || c.Value.Metadata.Quadrature.ErrorRigorous {
					t.Fatalf("quadrature semantics hidden or overstated: %+v", c.Value)
				}
			}
		}
	}
}

func TestM6FirstActualMatrixValues(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	want := []float64{0.0918849012616, -0.0113156378917, -0.0113156378917, 0.0877271147825}
	for i, e := range r.Evaluation.Matrix.Entries {
		if math.Abs(e.Value.Approximate.Real-want[i]) > 1e-11 {
			t.Fatalf("entry %d got %.15g want %.15g", i, e.Value.Approximate.Real, want[i])
		}
	}
	if r.Evaluation.Matrix.Entries[0].Contributions[1].Value.Exact.Real.Expression != "18/1" || r.Evaluation.Matrix.Entries[1].Contributions[1].Value.Exact.Real.Expression != "32/1" {
		t.Fatal("endpoint normalization drift")
	}
}

func TestM6StructuralHermitianAndNumericalDiagnosticStaySeparate(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	if !r.Evaluation.StructuralHermitian || r.Evaluation.Matrix.ValueSemantics != semantic.StructuralExactMatrix {
		t.Fatal("structural Hermitian certificate lost")
	}
	for _, d := range r.Evaluation.HermitianDiagnostics {
		if d.Discrepancy > 1e-14 || d.TheoremUse {
			t.Fatalf("symmetry diagnostic became theorem evidence: %+v", d)
		}
	}
}

func TestM6DirectQuadraticMatrixConsistencyIsExperimental(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Evaluation.DirectMatrixChecks) < 1 {
		t.Fatal("no end-to-end check")
	}
	for _, c := range r.Evaluation.DirectMatrixChecks {
		if c.Discrepancy > c.Tolerance || c.Evidence != "numerical_experiment_only" {
			t.Fatalf("direct/matrix mismatch or evidence escalation: %+v", c)
		}
	}
}

func TestM6WrongTransformConventionRejectsEvaluatorReuse(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	entry := r.Matrix.Entries[0]
	entry.TransformConvention = "other_mellin"
	if _, err := evaluateLogBoxEntry(entry); err == nil {
		t.Fatal("wrong entry transform accepted")
	}
	entry = r.Matrix.Entries[0]
	entry.RowFunction.TransformConvention = "other_mellin"
	if _, err := evaluateLogBoxEntry(entry); err == nil {
		t.Fatal("wrong function transform accepted")
	}
}

func TestM6ApproximateEigenvaluesCertifyNeitherPSDNorRH(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	if r.Evaluation.PSDCertified || r.Evaluation.RHCertified || r.Evaluation.EigenDiagnostic.CertifiesPSD {
		t.Fatal("approximate diagnostics crossed proof boundary")
	}
	if r.M5.MatrixPSDToRH.Accepted {
		t.Fatal("finite evaluated matrix certified RH")
	}
}

func TestM6QuadratureRefinementAndBasisOrderCounterprobes(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	p, err := newLogBoxPair(r.Evaluation.Matrix.Basis.Members[0].Function, r.Evaluation.Matrix.Basis.Members[1].Function)
	if err != nil {
		t.Fatal(err)
	}
	coarse := archimedeanValue(p, 1e-7)
	fine := archimedeanValue(p, 1e-12)
	if math.Abs(coarse-fine) > 1e-7 {
		t.Fatalf("archimedean refinement unstable: %.15g %.15g", coarse, fine)
	}
	reverse, err := newLogBoxPair(r.Evaluation.Matrix.Basis.Members[1].Function, r.Evaluation.Matrix.Basis.Members[0].Function)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(archimedeanValue(p, 1e-11)-archimedeanValue(reverse, 1e-11)) > 1e-14 {
		t.Fatal("basis reorder changed cross-entry")
	}
}

func TestM6ReportsDeterministicAndM0ThroughM5Regress(t *testing.T) {
	r, err := testM6()
	if err != nil {
		t.Fatal(err)
	}
	human := M6HumanReport(r)
	data, err := M6JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 1; i++ {
		again, err := CompileM6()
		if err != nil {
			t.Fatal(err)
		}
		got, err := M6JSONReport(again)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(data, got) || human != M6HumanReport(again) {
			t.Fatal("M6 report is nondeterministic")
		}
	}
	for _, want := range []string{"RIEMANN-M6", "unevaluated_exact_definition", "exact_value", "approximate_value", "exactly_zero_by_compact_support", "heuristic_quadrature_tolerance", "rh_certified", "function_space_restriction"} {
		if !strings.Contains(string(data)+human, want) {
			t.Fatalf("report omits %q", want)
		}
	}
	if !r.M5.M4.M3.M1.ZeroFreeCertified || !r.M5.M4.M3.StripCertified || !r.M5.M4.M3.SymmetryCertified {
		t.Fatal("M0-M5 regression")
	}
}
