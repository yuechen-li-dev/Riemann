package compiler

import (
	"fmt"
	"math"
	"math/cmplx"
	"sort"
	"strings"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	MellinLinearCombinationTheoremID semantic.TheoremID = "mellin-transform-linearity-on-finite-span"
	MellinConjugationTheoremID       semantic.TheoremID = "mellin-conjugation-for-real-test-functions"
	ZeroSummandCoordinateTheoremID   semantic.TheoremID = "weil-zero-summand-coordinate-lowering"
	ReflectionPairHermitianTheoremID semantic.TheoremID = "critical-reflection-pair-is-adjoint-sum"
	OuterProductPSDTheoremID         semantic.TheoremID = "outer-product-psd-rank-bound"
	PairRankSignatureTheoremID       semantic.TheoremID = "two-vector-hermitian-pair-rank-signature"
	OrbitPartitionTheoremID          semantic.TheoremID = "klein-four-orbit-partition-by-critical-reflection"
	ZeroMatrixAggregateTheoremID     semantic.TheoremID = "zero-summand-partition-into-symmetry-orbits"
)

var m8LinearAlgebraReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "Horn and Johnson, Matrix Analysis, 2nd ed., §§4.2, 7.1: Hermitian forms, outer products, rank, and two-vector restrictions"}
var m8MellinReference = semantic.Reference{Kind: semantic.StandardReference, Citation: "NIST DLMF §1.14(iv), Mellin transform: linearity and conjugation for real-valued inputs", URI: "https://dlmf.nist.gov/1.14.iv"}
var m8OrbitReference = semantic.Reference{Kind: semantic.CompilerRecord, Citation: "M8 exact Klein-four canonicalization and critical-reflection coset partition, derived from certified M3 zero-set symmetries"}

type ToyOrbitDiagnostic struct {
	Name           string                    `json:"name"`
	Inputs         string                    `json:"synthetic_inputs"`
	Matrix         [][]semantic.ComplexValue `json:"matrix"`
	Determinant    float64                   `json:"determinant"`
	Eigenvalues    []float64                 `json:"eigenvalues"`
	Classification string                    `json:"evidence_classification"`
}

type M8Result struct {
	M7                    M7Result                          `json:"-"`
	CriticalTemplate      semantic.OrbitMatrixContribution  `json:"critical_orbit_template"`
	OffCriticalTemplate   semantic.OrbitMatrixContribution  `json:"off_critical_orbit_template"`
	ZeroSide              semantic.ZeroSideMatrixAggregate  `json:"zero_side_matrix"`
	Dual                  semantic.DualMatrixRepresentation `json:"dual_representation"`
	ToyDiagnostics        []ToyOrbitDiagnostic              `json:"toy_diagnostics"`
	FinitePSDReverseProof ProofAttempt                      `json:"finite_psd_reverse_inference"`
}

func CompileM8() (M8Result, error) {
	m7, err := CompileM7()
	if err != nil {
		return M8Result{}, err
	}
	basis := m7.Evaluation.Matrix.Basis
	critical, err := semantic.NewZeroOrbit("critical-orbit-template", semantic.PointOnCriticalLine("rho_c"), semantic.CriticalLineOrbit, 1, []semantic.TheoremID{ConjugationInvariantTheoremID, ComposeInvarianceTheoremID}, dlmfZerosReference)
	if err != nil {
		return M8Result{}, err
	}
	off, err := semantic.NewZeroOrbit("off-critical-orbit-template", semantic.Point("rho_o"), semantic.OffCriticalOrbit, 1, []semantic.TheoremID{ConjugationInvariantTheoremID, ComposeInvarianceTheoremID}, dlmfZerosReference)
	if err != nil {
		return M8Result{}, err
	}
	criticalMatrix, err := LowerZeroOrbitToMatrix(critical, basis)
	if err != nil {
		return M8Result{}, err
	}
	offMatrix, err := LowerZeroOrbitToMatrix(off, basis)
	if err != nil {
		return M8Result{}, err
	}
	aggregate := semantic.ZeroSideMatrixAggregate{
		ID: "zero-side[" + m7.Evaluation.Matrix.ID + "]", SemanticMatrixID: m7.Evaluation.Matrix.ID, Basis: semantic.CloneOrderedBasis(basis), ZeroDomain: semantic.NontrivialZeros(semantic.RiemannZeta),
		OrbitContributions: nil, ContributionTemplates: []semantic.OrbitMatrixContribution{criticalMatrix, offMatrix}, SummationConvention: []string{"zeros counted with multiplicity", "symmetric limiting order inherited from the M4 explicit formula"}, Formula: "G=sum_over_zero_orbits G_O",
		CriticalAggregate: "P=sum_{O: critical_line} G_O", OffCriticalAggregate: "Q=sum_{O: off_critical} G_O", SplitIdentity: "G=P+Q", SourceFunctional: semantic.WeilZetaQuadraticFunctional, TransformConvention: semantic.LagariasMellinConvention,
		TheoremProvenance: []semantic.TheoremID{WeilExplicitFormulaTheoremID, ZeroMatrixAggregateTheoremID, OrbitPartitionTheoremID}, Provenance: lagariasWeilReference,
	}
	if err := aggregate.Validate(); err != nil {
		return M8Result{}, err
	}
	dual := semantic.DualMatrixRepresentation{SemanticMatrixID: m7.Evaluation.Matrix.ID, ZeroSideAggregateID: aggregate.ID, ExplicitFormulaMatrixID: m7.Evaluation.Matrix.ID, ExplicitValueEvidence: semantic.CertifiedInterval, IdentityTheorem: WeilExplicitFormulaTheoremID, Identity: "zero-side orbit sum = polarized explicit-formula matrix", NumericalIdentification: false, Provenance: lagariasWeilReference}
	// This deliberately asks the graph for an invalid strengthening. The only
	// available path still carries function_space_restriction and says nothing
	// about the signs of individual summands.
	reverse := m7.M6.M5.Graph.AttemptProof(M5MatrixPSDID, CriticalReflectionInvariantID)
	return M8Result{M7: m7, CriticalTemplate: criticalMatrix, OffCriticalTemplate: offMatrix, ZeroSide: aggregate, Dual: dual, ToyDiagnostics: m8ToyDiagnostics(), FinitePSDReverseProof: reverse}, nil
}

// LowerZeroPointSummand is generic in the basis and point. With v_i=M[f_i](p)
// and w_i=M[f_i](1-conjugate(p)), coefficient comparison gives
// K_ij=conjugate(w_i)v_j. This row/column order follows c* K c exactly.
func LowerZeroPointSummand(basis semantic.OrderedBasis, point semantic.PointExpr, multiplicity int) (semantic.PointMatrixContribution, error) {
	if err := basis.Validate(); err != nil {
		return semantic.PointMatrixContribution{}, err
	}
	if err := point.Validate(); err != nil {
		return semantic.PointMatrixContribution{}, err
	}
	if multiplicity < 1 {
		return semantic.PointMatrixContribution{}, fmt.Errorf("zero multiplicity must be positive")
	}
	v, err := semantic.NewBasisEvaluationVector(basis, point, MellinLinearCombinationTheoremID, m8MellinReference)
	if err != nil {
		return semantic.PointMatrixContribution{}, err
	}
	r := point.Apply(semantic.CriticalReflectionTransform)
	w, err := semantic.NewBasisEvaluationVector(basis, r, MellinLinearCombinationTheoremID, m8MellinReference)
	if err != nil {
		return semantic.PointMatrixContribution{}, err
	}
	n := len(basis.Members)
	entries := make([]semantic.SymbolicMatrixEntry, 0, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			expr := fmt.Sprintf("%d*conjugate(%s)*%s", multiplicity, w.Entries[i].Expression, v.Entries[j].Expression)
			entries = append(entries, semantic.SymbolicMatrixEntry{Row: i, Column: j, Expression: expr})
		}
	}
	fixed := point.Key() == r.Key()
	classification := semantic.ContributionClassification{RankUpperBound: 1, Theorems: []semantic.TheoremID{ZeroSummandCoordinateTheoremID}}
	if fixed {
		classification.Hermitian = true
		classification.PositiveSemidefinite = true
		classification.RankOneIfNonzero = true
		classification.Theorems = append(classification.Theorems, OuterProductPSDTheoremID)
	}
	k := semantic.PointMatrixContribution{ID: "K[" + point.Key() + "]", Point: point.Canonical(), ReflectedPoint: r, Evaluation: v, ReflectedEvaluation: w, Rows: n, Columns: n, Entries: entries,
		QuadraticIdentity: "Q_rho(sum_i c_i f_i)=c* K(rho) c", SourceSummand: "M[f](rho) * conjugate(M[f](1-conjugate(rho)))", Orientation: "K_ij=conjugate(w_i)*v_j; c* K c=(c^T v)*conjugate(c^T w)", Multiplicity: multiplicity, Classification: classification,
		TheoremProvenance: []semantic.TheoremID{WeilExplicitFormulaTheoremID, MellinLinearCombinationTheoremID, ZeroSummandCoordinateTheoremID}}
	return k, k.Validate()
}

func LowerZeroOrbitToMatrix(orbit semantic.ZeroOrbit, basis semantic.OrderedBasis) (semantic.OrbitMatrixContribution, error) {
	if err := orbit.Validate(); err != nil {
		return semantic.OrbitMatrixContribution{}, err
	}
	points := map[string]semantic.PointExpr{}
	for _, item := range orbit.TransformedPoints {
		points[item.Point.Key()] = item.Point
	}
	keys := make([]string, 0, len(points))
	for key := range points {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	used := map[string]bool{}
	var pairs []semantic.ReflectionPairContribution
	for _, key := range keys {
		if used[key] {
			continue
		}
		p := points[key]
		r := p.Apply(semantic.CriticalReflectionTransform)
		pairPoints := []semantic.PointExpr{p}
		used[key] = true
		if r.Key() != key {
			if _, ok := points[r.Key()]; !ok {
				return semantic.OrbitMatrixContribution{}, fmt.Errorf("orbit is not closed under critical reflection")
			}
			pairPoints = append(pairPoints, points[r.Key()])
			used[r.Key()] = true
		}
		sort.Slice(pairPoints, func(i, j int) bool { return pairPoints[i].Key() < pairPoints[j].Key() })
		pcs := make([]semantic.PointMatrixContribution, len(pairPoints))
		for i, point := range pairPoints {
			var err error
			pcs[i], err = LowerZeroPointSummand(basis, point, orbit.ZeroMultiplicity)
			if err != nil {
				return semantic.OrbitMatrixContribution{}, err
			}
		}
		n := len(basis.Members)
		entries := make([]semantic.SymbolicMatrixEntry, 0, n*n)
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				terms := make([]string, len(pcs))
				for k := range pcs {
					terms[k] = pcs[k].Entries[i*n+j].Expression
				}
				entries = append(entries, semantic.SymbolicMatrixEntry{Row: i, Column: j, Expression: strings.Join(terms, " + ")})
			}
		}
		cl := semantic.ContributionClassification{Hermitian: true, RankUpperBound: len(pairPoints), Theorems: []semantic.TheoremID{ReflectionPairHermitianTheoremID}}
		formula := "K(p)"
		if len(pairPoints) == 1 {
			cl = pcs[0].Classification
			formula = "m*conjugate(v(p))*transpose(v(p))"
		} else {
			cl.RankUpperBound = 2
			cl.IndefiniteCondition = "v(p) and v(1-conjugate(p)) are linearly independent"
			cl.DegenerateCondition = "evaluation vectors are linearly dependent; sign then depends on the real proportionality coefficient"
			cl.Theorems = append(cl.Theorems, PairRankSignatureTheoremID)
			formula = "K(p)+K(1-conjugate(p)) = a*b* + b*a*"
		}
		pairs = append(pairs, semantic.ReflectionPairContribution{ID: "reflection-pair[" + pairPoints[0].Key() + "]", Points: pairPoints, PointContributions: pcs, Entries: entries, Formula: formula, Classification: cl, GroupingTheorem: OrbitPartitionTheoremID})
	}
	n := len(basis.Members)
	entries := make([]semantic.SymbolicMatrixEntry, 0, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			terms := make([]string, len(pairs))
			for k := range pairs {
				terms[k] = pairs[k].Entries[i*n+j].Expression
			}
			entries = append(entries, semantic.SymbolicMatrixEntry{Row: i, Column: j, Expression: strings.Join(terms, " + ")})
		}
	}
	cl := semantic.ContributionClassification{Hermitian: true, RankUpperBound: smallerInt(n, orbit.DistinctLocationCount), Theorems: []semantic.TheoremID{OrbitPartitionTheoremID, ReflectionPairHermitianTheoremID}}
	if orbit.Classification == semantic.CriticalLineOrbit {
		cl.PositiveSemidefinite = true
		cl.Theorems = append(cl.Theorems, OuterProductPSDTheoremID)
	}
	return semantic.OrbitMatrixContribution{ID: "G[" + orbit.ID + "]", Orbit: orbit, Basis: semantic.CloneOrderedBasis(basis), ReflectionPairs: pairs, Entries: entries, Formula: "sum over distinct orbit points K(p), each weighted by zero multiplicity", Classification: cl, TheoremProvenance: []semantic.TheoremID{WeilExplicitFormulaTheoremID, OrbitPartitionTheoremID, ReflectionPairHermitianTheoremID}}, nil
}

func smallerInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func pairMatrix(v, w [2]complex128) [2][2]complex128 {
	var h [2][2]complex128
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			h[i][j] = cmplx.Conj(w[i])*v[j] + cmplx.Conj(v[i])*w[j]
		}
	}
	return h
}
func diag(name, inputs string, h [2][2]complex128) ToyOrbitDiagnostic {
	det := real(h[0][0]*h[1][1] - h[0][1]*h[1][0])
	tr := real(h[0][0] + h[1][1])
	disc := math.Sqrt(math.Max(0, tr*tr-4*det))
	rows := make([][]semantic.ComplexValue, 2)
	for i := 0; i < 2; i++ {
		rows[i] = make([]semantic.ComplexValue, 2)
		for j := 0; j < 2; j++ {
			rows[i][j] = semantic.ComplexValue{Real: real(h[i][j]), Imag: imag(h[i][j])}
		}
	}
	return ToyOrbitDiagnostic{Name: name, Inputs: inputs, Matrix: rows, Determinant: det, Eigenvalues: []float64{(tr - disc) / 2, (tr + disc) / 2}, Classification: "synthetic numerical experiment only"}
}
func m8ToyDiagnostics() []ToyOrbitDiagnostic {
	critical := fixedMatrix([2]complex128{1 + 2i, -1 + 0.5i})
	independent := pairMatrix([2]complex128{1, 1i}, [2]complex128{1, 1})
	dependent := pairMatrix([2]complex128{1, 1i}, [2]complex128{2, 2i})
	return []ToyOrbitDiagnostic{diag("critical-fixed-point", "v=w=(1+2i,-1+0.5i)", critical), diag("off-critical-independent", "v=(1,i), w=(1,1)", independent), diag("off-critical-dependent", "v=(1,i), w=2v", dependent)}
}

func fixedMatrix(v [2]complex128) [2][2]complex128 {
	var k [2][2]complex128
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			k[i][j] = cmplx.Conj(v[i]) * v[j]
		}
	}
	return k
}
