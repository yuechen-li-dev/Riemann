package compiler

import (
	"fmt"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M5SpanDefinitionID          semantic.ClaimID = "m5.weil-span.definition"
	M5QuadraticStructureID      semantic.ClaimID = "m5.weil-quadratic-form.structure"
	M5HermitianFormID           semantic.ClaimID = "m5.weil-hermitian-form.definition"
	M5HermitianMatrixID         semantic.ClaimID = "m5.weil-hermitian-matrix.definition"
	M5HermitianPropertyID       semantic.ClaimID = "m5.weil-matrix.hermitian"
	M5QuadraticMatrixIdentityID semantic.ClaimID = "m5.weil-matrix.quadratic-identity"
	M5SpanPositivityID          semantic.ClaimID = "m5.weil-span.positivity"
	M5CoordinatePositivityID    semantic.ClaimID = "m5.weil-coordinate.positivity"
	M5MatrixPSDID               semantic.ClaimID = "m5.weil-matrix.psd"
	M5ToyDiagonalID             semantic.ClaimID = "m5.toy-indefinite.diagonal-nonnegative"
	M5ToyDiagonalPSDID          semantic.ClaimID = "m5.toy-indefinite.psd"
	M5ToyExactPSDID             semantic.ClaimID = "m5.toy-polarization.psd"
	M5ToyNumericalPSDID         semantic.ClaimID = "m5.toy-polarization.numerical-psd-estimate"

	WeilQuadraticStructureTheoremID semantic.TheoremID = "weil-quadratic-form-on-complex-span"
	ComplexPolarizationTheoremID    semantic.TheoremID = "complex-polarization-identity"
	FiniteBasisLoweringTheoremID    semantic.TheoremID = "finite-basis-hermitian-matrix"
	QuadraticCoordinateTheoremID    semantic.TheoremID = "finite-basis-quadratic-coordinate-identity"
	FiniteSpanCoordinateTheoremID   semantic.TheoremID = "finite-span-coordinate-positivity-equivalence"
	FiniteHermitianPSDTheoremID     semantic.TheoremID = "finite-hermitian-psd-equivalence"
	HermitianConstructionTheoremID  semantic.TheoremID = "polarized-form-matrix-is-hermitian"

	RefineWeilQuadraticID semantic.TransformationID = "refine-weil-functional-to-quadratic-form"
	PolarizeWeilID        semantic.TransformationID = "polarize-weil-quadratic-form"
	LowerHermitianBasisID semantic.TransformationID = "lower-hermitian-form-to-ordered-basis"
	RecoverQuadraticID    semantic.TransformationID = "recover-quadratic-as-coordinate-form"
	DeriveHermitianID     semantic.TransformationID = "derive-hermitian-matrix-structure"
	RestrictWeilSpanID    semantic.TransformationID = "restrict-weil-positivity-to-finite-span"
	SpanToCoordinatesID   semantic.TransformationID = "identify-span-positivity-with-coordinate-positivity"
	CoordinatesToPSDID    semantic.TransformationID = "identify-coordinate-positivity-with-matrix-psd"
)

var polarizationReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "Encyclopedia of Mathematics, Polarization identity (complex quadratic forms); converted to the conjugate-linear-first convention",
	URI:      "https://encyclopediaofmath.org/wiki/Polarization_identity",
}

var finitePSDReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "P. Shor, MIT 8.370/18.435 Lecture 12: Hermitian positive semidefinite is equivalent to v*Av >= 0 for every vector v",
	URI:      "https://math.mit.edu/~shor/435-LN/Lecture_12.pdf",
}

var hermitianBasisReference = semantic.Reference{
	Kind:     semantic.StandardReference,
	Citation: "D. Vogan, MIT notes, Hermitian Forms, Proposition 2.1: a form is Hermitian iff its matrix in a basis is Hermitian",
	URI:      "https://math.mit.edu/~dav/hermitian.pdf",
}

type M5Result struct {
	M4                    M4Result
	Graph                 *Graph
	Registry              *ContractRegistry
	Span                  semantic.FiniteSpan
	Form                  semantic.HermitianForm
	Matrix                semantic.HermitianMatrix
	Combination           semantic.FiniteLinearCombination
	FullToSpan            ProofAttempt
	FamilyToSpan          ProofAttempt
	DiagonalToPSD         ProofAttempt
	ApproximateToExactPSD ProofAttempt
	MatrixPSDToRH         ProofAttempt
	HermitianCertified    bool
	HermitianDiagnostics  []Diagnostic
}

func CompileM5() (M5Result, error) {
	m4, err := CompileM4()
	if err != nil {
		return M5Result{}, err
	}
	g := m4.Graph
	span := m5WeilSpan(m4)
	spanRef := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M5 ordered basis and finite-span construction with per-member M4 admissibility certificates"}
	spanClaim := semantic.Claim{ID: M5SpanDefinitionID, Proposition: semantic.FiniteSpanDefinition{Span: span}, Evidence: []semantic.Evidence{{Kind: semantic.DefinitionEvidence, Source: spanRef}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: spanRef}}
	if err := g.AddClaim(spanClaim); err != nil {
		return M5Result{}, err
	}

	registry, err := m5TheoremRegistry()
	if err != nil {
		return M5Result{}, err
	}
	structure := semantic.QuadraticFormStructure{Functional: semantic.WeilZetaQuadraticFunctional, DomainSpan: span, CoefficientField: semantic.ComplexField, Laws: []semantic.QuadraticLaw{semantic.AbsoluteSquareHomogeneity, semantic.ParallelogramLaw, semantic.RealValuedDiagonal}, Theorem: WeilQuadraticStructureTheoremID}
	structureRef := semantic.Reference{Kind: semantic.StandardReference, Citation: "Lagarias's Weil functional construction, restricted to a complex span, together with the standard Hermitian quadratic-form hypotheses", URI: lagariasWeilReference.URI}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5QuadraticStructureID, Proposition: structure, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: structureRef, Note: "certifies only the quadratic-form laws, not nonnegativity"}}, Exactness: semantic.Exact}, WeilFunctionalDefinitionID, []semantic.ClaimID{M5SpanDefinitionID}, RefineWeilQuadraticID, "refine-quadratic-functional", Implies, WeilQuadraticStructureTheoremID, structureRef); err != nil {
		return M5Result{}, err
	}

	functionalClaim, _ := g.Claim(WeilFunctionalDefinitionID)
	functional := functionalClaim.Proposition.(semantic.FunctionalDefinition).Functional
	form, err := PolarizeQuadraticFunctional(structure, functional, ComplexPolarizationTheoremID)
	if err != nil {
		return M5Result{}, err
	}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5HermitianFormID, Proposition: semantic.HermitianFormDefinition{Form: form}, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: polarizationReference}}, Exactness: semantic.Exact}, M5QuadraticStructureID, nil, PolarizeWeilID, "polarize-quadratic-functional", Equivalent, ComplexPolarizationTheoremID, polarizationReference); err != nil {
		return M5Result{}, err
	}

	matrix, err := LowerHermitianFormToMatrix(form, span.Basis, FiniteBasisLoweringTheoremID)
	if err != nil {
		return M5Result{}, err
	}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5HermitianMatrixID, Proposition: semantic.HermitianMatrixDefinition{Matrix: matrix}, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: hermitianBasisReference, Note: "structural entries; no numerical evaluation and no PSD claim"}}, Exactness: semantic.Exact}, M5HermitianFormID, []semantic.ClaimID{M5SpanDefinitionID}, LowerHermitianBasisID, "lower-form-to-ordered-basis", Equivalent, FiniteBasisLoweringTheoremID, hermitianBasisReference); err != nil {
		return M5Result{}, err
	}

	hermitianSpan := semantic.CloneFiniteSpan(span)
	hermitian := semantic.MatrixProperty{MatrixID: matrix.ID, SourceFunctional: semantic.WeilZetaQuadraticFunctional, DomainSpan: &hermitianSpan, Property: semantic.MatrixHermitian, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "G_ij=conjugate(G_ji) follows from B(f_i,f_j)=conjugate(B(f_j,f_i))"}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5HermitianPropertyID, Proposition: hermitian, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: polarizationReference, Note: "Hermitian by construction; does not assert PSD"}}, Exactness: semantic.Exact}, M5HermitianMatrixID, nil, DeriveHermitianID, "derive-hermitian-structure", Implies, HermitianConstructionTheoremID, polarizationReference); err != nil {
		return M5Result{}, err
	}

	combination := m5GenericCombination(span)
	identity := semantic.QuadraticMatrixIdentity{Functional: semantic.WeilZetaQuadraticFunctional, Combination: combination, MatrixID: matrix.ID, FormID: form.ID, Identity: "Q(sum_i c_i f_i)=c* G c", Theorem: QuadraticCoordinateTheoremID}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5QuadraticMatrixIdentityID, Proposition: identity, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: hermitianBasisReference}}, Exactness: semantic.Exact}, M5HermitianMatrixID, []semantic.ClaimID{M5QuadraticStructureID}, RecoverQuadraticID, "recover-quadratic-coordinate-identity", Implies, QuadraticCoordinateTheoremID, hermitianBasisReference); err != nil {
		return M5Result{}, err
	}

	spanPositivity := semantic.FiniteSpanFunctionalStatement{Quantifier: semantic.ForAll, Variable: "f", Span: span, Functional: semantic.WeilZetaQuadraticFunctional, Predicate: semantic.FunctionalNonnegative, TransformConvention: semantic.LagariasMellinConvention}
	certificates := m5BasisCertificates(span.Basis)
	parents := append([]semantic.ClaimID{WeilPositivityID}, certificates...)
	restrictionRef := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "ForAll restriction from the Weil-nice class to a certified finite complex span"}
	spanPositivityClaim := semantic.Claim{ID: M5SpanPositivityID, Proposition: spanPositivity, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: restrictionRef, Note: "strict function-space coverage loss"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: RestrictWeilSpanID, Source: restrictionRef}}
	if err := g.AddClaim(spanPositivityClaim); err != nil {
		return M5Result{}, err
	}
	if err := g.AddTransformation(Transformation{ID: RestrictWeilSpanID, Pass: "restrict-function-space-to-finite-span", From: WeilPositivityID, Premises: certificates, To: M5SpanPositivityID, Relation: Implies, Losses: []InformationLoss{{Kind: FunctionSpaceRestriction, Reason: "source covers every Weil-nice test function; conclusion covers one finite-dimensional admissible span"}}, Provenance: restrictionRef}); err != nil {
		return M5Result{}, err
	}

	coordinates := semantic.CoordinateQuadraticPositivity{MatrixID: matrix.ID, Span: span, CoefficientField: semantic.ComplexField, Expression: "c* G c >= 0 for every coefficient vector c"}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5CoordinatePositivityID, Proposition: coordinates, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: finitePSDReference}}, Exactness: semantic.Exact}, M5SpanPositivityID, []semantic.ClaimID{M5QuadraticMatrixIdentityID}, SpanToCoordinatesID, "identify-span-and-coordinate-positivity", Equivalent, FiniteSpanCoordinateTheoremID, finitePSDReference); err != nil {
		return M5Result{}, err
	}

	psdSpan := semantic.CloneFiniteSpan(span)
	psd := semantic.MatrixProperty{MatrixID: matrix.ID, SourceFunctional: semantic.WeilZetaQuadraticFunctional, DomainSpan: &psdSpan, Property: semantic.MatrixPositiveSemidefinite, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "c* G c >= 0 for every c in C^n"}
	if err := addDerivedClaimAndTransformation(g, semantic.Claim{ID: M5MatrixPSDID, Proposition: psd, Evidence: []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: finitePSDReference, Note: "equivalent finite-span target; not an assertion that the open target holds"}}, Exactness: semantic.Exact}, M5CoordinatePositivityID, []semantic.ClaimID{M5HermitianPropertyID}, CoordinatesToPSDID, "identify-coordinate-positivity-with-psd", Equivalent, FiniteHermitianPSDTheoremID, finitePSDReference); err != nil {
		return M5Result{}, err
	}

	if err := addM5BoundaryFixtures(g); err != nil {
		return M5Result{}, err
	}
	hermitianCertified, hermitianDiagnostics := g.Certify(M5HermitianPropertyID)
	return M5Result{
		M4: m4, Graph: g, Registry: registry, Span: span, Form: form, Matrix: matrix, Combination: combination,
		FullToSpan:            g.CheckDischarge(WeilPositivityID, M5SpanPositivityID),
		FamilyToSpan:          g.CheckDischarge(FiniteWeilPositivityID, M5SpanPositivityID),
		DiagonalToPSD:         g.CheckDischarge(M5ToyDiagonalID, M5ToyDiagonalPSDID),
		ApproximateToExactPSD: g.AttemptProof(M5ToyNumericalPSDID, M5ToyExactPSDID),
		MatrixPSDToRH:         g.AttemptProof(M5MatrixPSDID, RHClaimID),
		HermitianCertified:    hermitianCertified, HermitianDiagnostics: hermitianDiagnostics,
	}, nil
}

func m5TheoremRegistry() (*ContractRegistry, error) {
	r := NewContractRegistry()
	contracts := []TheoremContract{
		{ID: WeilQuadraticStructureTheoremID, Premises: []Pattern{{Kind: semantic.FunctionalDefinitionKind, Exactness: semantic.Exact}, {Kind: semantic.FiniteSpanDefinitionKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.QuadraticFormStructureKind, Exactness: semantic.Exact}, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: lagariasWeilReference}, Trust: TrustedExternalTheorem, Citation: "Lagarias §3 Weil functional construction, under the explicitly recorded complex quadratic-form laws"},
		{ID: ComplexPolarizationTheoremID, Premises: []Pattern{{Kind: semantic.QuadraticFormStructureKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.HermitianFormDefinitionKind, Exactness: semantic.Exact}, Relation: Equivalent, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: polarizationReference}, Trust: TrustedExternalTheorem, Citation: polarizationReference.Citation},
		{ID: FiniteBasisLoweringTheoremID, Premises: []Pattern{{Kind: semantic.HermitianFormDefinitionKind, Exactness: semantic.Exact}, {Kind: semantic.FiniteSpanDefinitionKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.HermitianMatrixDefinitionKind, Exactness: semantic.Exact}, Relation: Equivalent, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: hermitianBasisReference}, Trust: TrustedExternalTheorem, Citation: hermitianBasisReference.Citation},
		{ID: QuadraticCoordinateTheoremID, Premises: []Pattern{{Kind: semantic.HermitianMatrixDefinitionKind, Exactness: semantic.Exact}, {Kind: semantic.QuadraticFormStructureKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.QuadraticMatrixIdentityKind, Exactness: semantic.Exact}, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: hermitianBasisReference}, Trust: TrustedExternalTheorem, Citation: "finite coordinate expansion of a Hermitian form"},
		{ID: FiniteSpanCoordinateTheoremID, Premises: []Pattern{{Kind: semantic.FiniteSpanFunctionalStatementKind, Exactness: semantic.Exact}, {Kind: semantic.QuadraticMatrixIdentityKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.CoordinatePositivityKind, Exactness: semantic.Exact}, Relation: Equivalent, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: finitePSDReference}, Trust: TrustedExternalTheorem, Citation: "finite-span vectors are exactly the coefficient vectors in an ordered basis"},
		{ID: HermitianConstructionTheoremID, Premises: []Pattern{{Kind: semantic.HermitianMatrixDefinitionKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.MatrixPropertyStatementKind, Exactness: semantic.Exact}, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: hermitianBasisReference}, Trust: TrustedExternalTheorem, Citation: hermitianBasisReference.Citation},
		{ID: FiniteHermitianPSDTheoremID, Premises: []Pattern{{Kind: semantic.CoordinatePositivityKind, Exactness: semantic.Exact}, {Kind: semantic.MatrixPropertyStatementKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.MatrixPropertyStatementKind, Exactness: semantic.Exact}, Relation: Equivalent, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: finitePSDReference}, Trust: TrustedExternalTheorem, Citation: finitePSDReference.Citation},
	}
	for _, contract := range contracts {
		if err := r.Register(contract); err != nil {
			return nil, err
		}
	}
	return r, nil
}

// PolarizeQuadraticFunctional is the reusable functional -> form lowering. It
// refuses broad QuadraticFunctional values without the refined laws.
func PolarizeQuadraticFunctional(structure semantic.QuadraticFormStructure, functional semantic.QuadraticFunctional, theorem semantic.TheoremID) (semantic.HermitianForm, error) {
	if err := structure.Validate(); err != nil {
		return semantic.HermitianForm{}, fmt.Errorf("polarize: %w", err)
	}
	if err := functional.Validate(); err != nil {
		return semantic.HermitianForm{}, fmt.Errorf("polarize: %w", err)
	}
	if structure.Functional != functional.ID || theorem == "" {
		return semantic.HermitianForm{}, fmt.Errorf("polarize: structure does not certify source functional")
	}
	contributions := make([]semantic.FormContribution, len(functional.Contributions))
	for i, contribution := range functional.Contributions {
		contributions[i] = semantic.FormContribution{SourceKind: contribution.Kind, RepresentationSide: contribution.RepresentationSide, Sign: contribution.Sign, EntryDefinition: fmt.Sprintf("polarize[%s](f,g)", contribution.Kind)}
	}
	form := semantic.HermitianForm{ID: string(functional.ID) + ".polarized", SourceFunctional: functional.ID, DomainSpan: semantic.CloneFiniteSpan(structure.DomainSpan), Convention: semantic.ComplexConjugateFirstPolarization(theorem), EntryDefinition: "B(f,g) from the recorded four-evaluation complex polarization identity", Contributions: contributions, RecoveryIdentity: "Q(f)=B(f,f)", HermitianIdentity: "B(f,g)=conjugate(B(g,f))", TheoremProvenance: []semantic.TheoremID{structure.Theorem, theorem}}
	return form, form.Validate()
}

// LowerHermitianFormToMatrix is generic over the source functional: it only
// needs a validated form and an ordered basis for the same span.
func LowerHermitianFormToMatrix(form semantic.HermitianForm, basis semantic.OrderedBasis, theorem semantic.TheoremID) (semantic.HermitianMatrix, error) {
	if err := form.Validate(); err != nil {
		return semantic.HermitianMatrix{}, fmt.Errorf("matrix lowering: %w", err)
	}
	if err := basis.Validate(); err != nil {
		return semantic.HermitianMatrix{}, fmt.Errorf("matrix lowering: %w", err)
	}
	if theorem == "" || form.DomainSpan.Key() != (semantic.FiniteSpan{ID: form.DomainSpan.ID, Basis: basis, CoefficientField: basis.CoefficientField, ParentClass: basis.ParentClass}).Key() {
		return semantic.HermitianMatrix{}, fmt.Errorf("matrix lowering: basis does not coordinate form domain")
	}
	n := len(basis.Members)
	entries := make([]semantic.MatrixEntry, 0, n*n)
	for i, row := range basis.Members {
		for j, column := range basis.Members {
			parts := make([]semantic.EntryContribution, len(form.Contributions))
			for k, contribution := range form.Contributions {
				parts[k] = semantic.EntryContribution{SourceKind: contribution.SourceKind, RepresentationSide: contribution.RepresentationSide, Sign: contribution.Sign, PolarizedDefinition: contribution.EntryDefinition + "(" + row.Function.Symbol + "," + column.Function.Symbol + ")"}
			}
			entries = append(entries, semantic.MatrixEntry{Row: i, Column: j, RowFunction: semantic.CloneTestFunction(row.Function), ColumnFunction: semantic.CloneTestFunction(column.Function), SourceForm: form.ID, SourceFunctional: form.SourceFunctional, Definition: fmt.Sprintf("%s(%s,%s)", form.ID, row.Function.Symbol, column.Function.Symbol), Contributions: parts, TransformConvention: row.Function.TransformConvention, TheoremProvenance: append(append([]semantic.TheoremID(nil), form.TheoremProvenance...), theorem)})
		}
	}
	matrix := semantic.HermitianMatrix{ID: "G[" + basis.ID + "]", SourceForm: semantic.CloneHermitianForm(form), Basis: semantic.CloneOrderedBasis(basis), Rows: n, Columns: n, Entries: entries, ValueSemantics: semantic.StructuralExactMatrix, HermitianByConstruction: true, LoweringTheorem: theorem}
	return matrix, matrix.Validate()
}

func m5WeilSpan(m4 M4Result) semantic.FiniteSpan {
	members := []semantic.BasisMember{{Function: semantic.CloneTestFunction(m4.FiniteFamily.Members[0]), AdmissibilityCertificate: WeilF1AdmissibleID}, {Function: semantic.CloneTestFunction(m4.FiniteFamily.Members[1]), AdmissibilityCertificate: WeilF2AdmissibleID}}
	basis := semantic.OrderedBasis{ID: "m5-weil-demo-basis", Members: members, CoefficientField: semantic.ComplexField, ParentClass: semantic.WeilNiceClass()}
	return semantic.FiniteSpan{ID: "m5-weil-demo-span", Basis: basis, CoefficientField: semantic.ComplexField, ParentClass: semantic.WeilNiceClass()}
}

func m5GenericCombination(span semantic.FiniteSpan) semantic.FiniteLinearCombination {
	entries := make([]semantic.Coefficient, len(span.Basis.Members))
	for i := range entries {
		entries[i] = semantic.Coefficient{Symbol: fmt.Sprintf("c%d", i+1)}
	}
	return semantic.FiniteLinearCombination{ID: "m5-generic-linear-combination", Span: semantic.CloneFiniteSpan(span), Coefficients: semantic.CoefficientVector{ID: "c", Field: span.CoefficientField, Entries: entries}}
}

func m5BasisCertificates(basis semantic.OrderedBasis) []semantic.ClaimID {
	out := make([]semantic.ClaimID, len(basis.Members))
	for i, member := range basis.Members {
		out[i] = member.AdmissibilityCertificate
	}
	return out
}

func addDerivedClaimAndTransformation(g *Graph, claim semantic.Claim, from semantic.ClaimID, premises []semantic.ClaimID, transformation semantic.TransformationID, pass string, relation Relation, theorem semantic.TheoremID, ref semantic.Reference) error {
	parents := append([]semantic.ClaimID{from}, premises...)
	claim.Assumptions = nil
	claim.Provenance = semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: transformation, Theorem: theorem, Source: ref}
	if err := g.AddClaim(claim); err != nil {
		return err
	}
	return g.AddTransformation(Transformation{ID: transformation, Pass: pass, From: from, Premises: premises, To: claim.ID, Relation: relation, Provenance: ref, Theorem: theorem, Trusted: true})
}

func addM5BoundaryFixtures(g *Graph) error {
	exactRef := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M5 exact toy matrix [[0,1],[1,0]] has nonnegative diagonal but is not thereby certified PSD"}
	claims := []semantic.Claim{
		{ID: M5ToyDiagonalID, Proposition: semantic.MatrixProperty{MatrixID: "G[toy-offdiagonal]", SourceFunctional: "toy-offdiagonal", Property: semantic.MatrixDiagonalNonnegative, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "G_11=G_22=0"}, Evidence: []semantic.Evidence{{Kind: semantic.DefinitionEvidence, Source: exactRef}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: exactRef}},
		{ID: M5ToyDiagonalPSDID, Proposition: semantic.MatrixProperty{MatrixID: "G[toy-offdiagonal]", SourceFunctional: "toy-offdiagonal", Property: semantic.MatrixPositiveSemidefinite, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "c* G c >= 0 for every c"}, Evidence: []semantic.Evidence{{Kind: semantic.UnverifiedConjectureEvidence, Source: exactRef, Note: "deliberately false target used to test that diagonal evidence cannot prove it"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: exactRef}},
		{ID: M5ToyExactPSDID, Proposition: semantic.MatrixProperty{MatrixID: "G[toy-polarization]", SourceFunctional: "toy-polarization", Property: semantic.MatrixPositiveSemidefinite, ValueSemantics: semantic.StructuralExactMatrix, Criterion: "c* G c >= 0 for every c"}, Evidence: []semantic.Evidence{{Kind: semantic.UnverifiedConjectureEvidence, Source: exactRef, Note: "exact target intentionally lacks an exact verification step"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: exactRef}},
		{ID: M5ToyNumericalPSDID, Proposition: semantic.MatrixProperty{MatrixID: "G[toy-polarization]", SourceFunctional: "toy-polarization", Property: semantic.MatrixNumericallyPSDEstimated, ValueSemantics: semantic.NumericallyApproximateMatrix, Criterion: "bounded Float probes appeared nonnegative"}, Evidence: []semantic.Evidence{{Kind: semantic.NumericalExperimentEvidence, Source: semantic.Reference{Kind: semantic.ExperimentRecord, Citation: "M5 Oct polarization toy"}, Note: "numerical sampling is non-certifying"}}, Exactness: semantic.Approximate, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: semantic.Reference{Kind: semantic.ExperimentRecord, Citation: "M5 Oct polarization toy"}}},
	}
	for _, claim := range claims {
		if err := g.AddClaim(claim); err != nil {
			return err
		}
	}
	return nil
}
