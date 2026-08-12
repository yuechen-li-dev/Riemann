package compiler

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM5GenericPolarizationAndMatrixLowering(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	if r.Form.Convention.ConjugateLinearArgument != semantic.FirstArgument || r.Form.Convention.LinearArgument != semantic.SecondArgument || !strings.Contains(r.Form.Convention.Formula, "-i*Q(f+i*g)") {
		t.Fatalf("wrong polarization convention: %+v", r.Form.Convention)
	}
	if r.Form.RecoveryIdentity != "Q(f)=B(f,f)" || r.Form.HermitianIdentity != "B(f,g)=conjugate(B(g,f))" {
		t.Fatalf("form identities missing: %+v", r.Form)
	}
	if r.Matrix.Rows != 2 || r.Matrix.Columns != 2 || len(r.Matrix.Entries) != 4 || r.Matrix.ValueSemantics != semantic.StructuralExactMatrix || !r.Matrix.HermitianByConstruction {
		t.Fatalf("bad matrix lowering: %+v", r.Matrix)
	}
	if !r.HermitianCertified || len(r.HermitianDiagnostics) != 0 {
		t.Fatalf("Hermitian construction was not certified: %+v", r.HermitianDiagnostics)
	}
	if certified, _ := r.Graph.Certify(M5MatrixPSDID); certified {
		t.Fatal("structural Hermitian construction spuriously certified PSD")
	}
	identityClaim, ok := r.Graph.Claim(M5QuadraticMatrixIdentityID)
	if !ok {
		t.Fatal("quadratic matrix identity missing")
	}
	identity := identityClaim.Proposition.(semantic.QuadraticMatrixIdentity)
	if identity.Identity != "Q(sum_i c_i f_i)=c* G c" || identity.Combination.Coefficients.Field != semantic.ComplexField || len(identity.Combination.Coefficients.Entries) != 2 {
		t.Fatalf("coordinate identity is not structural: %+v", identity)
	}
}

func TestM5PolarizationRequiresAllQuadraticLaws(t *testing.T) {
	r, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	span := m5WeilSpan(r)
	claim, _ := r.Graph.Claim(WeilFunctionalDefinitionID)
	q := claim.Proposition.(semantic.FunctionalDefinition).Functional
	structure := semantic.QuadraticFormStructure{Functional: q.ID, DomainSpan: span, CoefficientField: semantic.ComplexField, Laws: []semantic.QuadraticLaw{semantic.AbsoluteSquareHomogeneity, semantic.RealValuedDiagonal}, Theorem: WeilQuadraticStructureTheoremID}
	if _, err := PolarizeQuadraticFunctional(structure, q, ComplexPolarizationTheoremID); err == nil || !strings.Contains(err.Error(), "prerequisite") {
		t.Fatalf("missing parallelogram law was silently accepted: %v", err)
	}
}

func TestM5FamilyAndSpanAreNotConflated(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	if r.FamilyToSpan.Accepted || !hasCode(r.FamilyToSpan.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("basis-point positivity proved span positivity: %+v", r.FamilyToSpan)
	}
	family, _ := r.Graph.Claim(FiniteWeilPositivityID)
	span, _ := r.Graph.Claim(M5SpanPositivityID)
	if family.Proposition.Kind() == span.Proposition.Kind() {
		t.Fatal("finite family and finite span share a proposition kind")
	}
	if !r.FullToSpan.Accepted {
		t.Fatalf("valid universal-to-span restriction failed: %+v", r.FullToSpan)
	}
	var sawLoss bool
	for _, transformation := range r.Graph.Transformations() {
		if transformation.ID == RestrictWeilSpanID {
			sawLoss = len(transformation.Losses) == 1 && transformation.Losses[0].Kind == FunctionSpaceRestriction && len(transformation.Premises) == len(r.Span.Basis.Members)
		}
	}
	if !sawLoss {
		t.Fatal("finite-span restriction lost its admissibility premises or function-space loss")
	}
}

func TestM5GraphCannotBypassFiniteMatrixCoverageLoss(t *testing.T) {
	m4, err := CompileM4()
	if err != nil {
		t.Fatal(err)
	}
	span := m5WeilSpan(m4)
	spanCopy := semantic.CloneFiniteSpan(span)
	ref := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "direct matrix restriction soundness test"}
	target := semantic.Claim{ID: "m5-direct-matrix", Proposition: semantic.MatrixProperty{MatrixID: "G[direct]", SourceFunctional: semantic.WeilZetaQuadraticFunctional, DomainSpan: &spanCopy, Property: semantic.MatrixPositiveSemidefinite, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "c*G c >= 0 for all c"}, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: ref}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: []semantic.ClaimID{WeilPositivityID}, Transformation: "m5-direct", Source: ref}}
	if err := m4.Graph.AddClaim(target); err != nil {
		t.Fatal(err)
	}
	if err := m4.Graph.AddTransformation(Transformation{ID: "m5-direct", Pass: "test-direct-matrix", From: WeilPositivityID, To: target.ID, Relation: Implies, Provenance: ref}); err == nil || !strings.Contains(err.Error(), "function_space_restriction") {
		t.Fatalf("direct finite matrix edge omitted coverage loss without rejection: %v", err)
	}
	if err := m4.Graph.AddTransformation(Transformation{ID: "m5-direct-with-loss", Pass: "test-direct-matrix", From: WeilPositivityID, To: target.ID, Relation: Implies, Losses: []InformationLoss{{Kind: FunctionSpaceRestriction, Reason: "test"}}, Provenance: ref}); err == nil || !strings.Contains(err.Error(), "admissibility premise") {
		t.Fatalf("direct finite matrix edge omitted admissibility without rejection: %v", err)
	}
}

func TestM5DiagonalApproximationAndRHBoundaries(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	if r.DiagonalToPSD.Accepted || !hasCode(r.DiagonalToPSD.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("nonnegative diagonal certified PSD: %+v", r.DiagonalToPSD)
	}
	if r.ApproximateToExactPSD.Accepted || !hasCode(r.ApproximateToExactPSD.Diagnostics, ApproximationBoundary) || !hasCode(r.ApproximateToExactPSD.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("approximate PSD estimate certified exact PSD: %+v", r.ApproximateToExactPSD)
	}
	if r.MatrixPSDToRH.Accepted || !hasCode(r.MatrixPSDToRH.Diagnostics, InformationLost) || !hasCode(r.MatrixPSDToRH.Diagnostics, NoEstablishedDirection) {
		t.Fatalf("finite matrix PSD certified RH or hid restriction: %+v", r.MatrixPSDToRH)
	}
	if !strings.Contains(r.MatrixPSDToRH.Diagnostics[0].Message+r.MatrixPSDToRH.Diagnostics[len(r.MatrixPSDToRH.Diagnostics)-1].Message, "function_space_restriction") {
		t.Fatalf("RH rejection does not expose the function-space boundary: %+v", r.MatrixPSDToRH.Diagnostics)
	}
}

func TestM5ComponentProvenanceAndBasisPermutation(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	want := []semantic.FunctionalContributionKind{semantic.ZeroContribution, semantic.EndpointContribution, semantic.PrimePowerContribution, semantic.ArchimedeanContribution}
	for _, entry := range r.Matrix.Entries {
		if entry.RowFunction.Key() != r.Span.Basis.Members[entry.Row].Function.Key() || entry.ColumnFunction.Key() != r.Span.Basis.Members[entry.Column].Function.Key() || entry.TransformConvention != semantic.LagariasMellinConvention || len(entry.TheoremProvenance) < 3 {
			t.Fatalf("entry provenance incomplete: %+v", entry)
		}
		if len(entry.Contributions) != len(want) {
			t.Fatalf("entry decomposition missing: %+v", entry.Contributions)
		}
		for i, kind := range want {
			if entry.Contributions[i].SourceKind != kind {
				t.Fatalf("entry component %d got %s want %s", i, entry.Contributions[i].SourceKind, kind)
			}
		}
	}
	permuted := semantic.CloneOrderedBasis(r.Span.Basis)
	permuted.Members[0], permuted.Members[1] = permuted.Members[1], permuted.Members[0]
	if permuted.Key() == r.Span.Basis.Key() {
		t.Fatal("basis permutation did not change coordinate identity")
	}
	permutedMatrix, err := LowerHermitianFormToMatrix(r.Form, permuted, FiniteBasisLoweringTheoremID)
	if err != nil {
		t.Fatal(err)
	}
	if permutedMatrix.Entries[0].RowFunction.Key() != r.Matrix.Entries[3].RowFunction.Key() || permutedMatrix.SourceForm.Key() != r.Matrix.SourceForm.Key() || permutedMatrix.SourceForm.DomainSpan.Key() != r.Matrix.SourceForm.DomainSpan.Key() {
		t.Fatal("basis permutation did not permute coordinates while retaining the form/span identity")
	}
}

func TestM5FiniteSpanPSDEquivalenceUsesTrustedContracts(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	for _, id := range []semantic.TheoremID{ComplexPolarizationTheoremID, FiniteBasisLoweringTheoremID, QuadraticCoordinateTheoremID, FiniteSpanCoordinateTheoremID, FiniteHermitianPSDTheoremID, HermitianConstructionTheoremID} {
		contract, ok := r.Registry.Contract(id)
		if !ok || contract.Trust != TrustedExternalTheorem || contract.Evidence.Kind != semantic.KnownTheoremEvidence {
			t.Fatalf("missing trusted theorem contract %s: %+v", id, contract)
		}
	}
	if a := r.Graph.CheckDischarge(M5SpanPositivityID, M5CoordinatePositivityID); !a.Accepted {
		t.Fatalf("span-to-coordinate equivalence failed: %+v", a)
	}
	if a := r.Graph.CheckDischarge(M5CoordinatePositivityID, M5MatrixPSDID); !a.Accepted {
		t.Fatalf("coordinate-to-PSD equivalence failed: %+v", a)
	}
	if a := r.Graph.CheckDischarge(M5MatrixPSDID, M5SpanPositivityID); !a.Accepted {
		t.Fatalf("finite PSD equivalence is not reversible: %+v", a)
	}
}

func TestM5ReportsDeterministicAndM0ThroughM4Regress(t *testing.T) {
	r, err := CompileM5()
	if err != nil {
		t.Fatal(err)
	}
	wantJSON, err := M5JSONReport(r)
	if err != nil {
		t.Fatal(err)
	}
	wantHuman := M5HumanReport(r)
	for i := 0; i < 3; i++ {
		got, err := CompileM5()
		if err != nil {
			t.Fatal(err)
		}
		data, _ := M5JSONReport(got)
		if !bytes.Equal(data, wantJSON) || M5HumanReport(got) != wantHuman {
			t.Fatal("M5 output is nondeterministic")
		}
	}
	for _, text := range []string{"RIEMANN-M5", "POLARIZE", "BASIS LOWERING", "function_space_restriction", "c* G c", "RH CERTIFICATION\n  unavailable", "numerical only"} {
		if !strings.Contains(wantHuman, text) {
			t.Fatalf("human report omits %q", text)
		}
	}
	for _, text := range []string{`"schema": "riemann.semantic-graph.m5"`, `"finite_span_definition"`, `"hermitian_matrix_definition"`, `"structurally_defined_exact"`, `"numerically_evaluated_approximate"`} {
		if !strings.Contains(string(wantJSON), text) {
			t.Fatalf("JSON omits %q", text)
		}
	}
	if !r.M4.M3.M1.ZeroFreeCertified || !r.M4.M3.StripCertified || !r.M4.M3.SymmetryCertified {
		t.Fatal("M0-M4 soundness regression")
	}
}
