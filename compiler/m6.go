package compiler

import (
	"fmt"
	"math"
	"math/cmplx"
	"sort"
	"strconv"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const WeilCrossEntryTheoremID semantic.TheoremID = "weil-lagarias.intersection-product-explicit-formula"

type HermitianConsistencyDiagnostic struct {
	Row         int     `json:"row"`
	Column      int     `json:"column"`
	Discrepancy float64 `json:"discrepancy"`
	TheoremUse  bool    `json:"theorem_use"`
}

type DirectMatrixCheck struct {
	Coefficients []semantic.ComplexValue `json:"coefficients"`
	Direct       semantic.ComplexValue   `json:"direct_quadratic"`
	Matrix       semantic.ComplexValue   `json:"matrix_quadratic"`
	Discrepancy  float64                 `json:"discrepancy"`
	Tolerance    float64                 `json:"declared_tolerance"`
	Evidence     string                  `json:"evidence"`
}

type EigenDiagnostic struct {
	Eigenvalues  []float64 `json:"eigenvalues"`
	Condition    float64   `json:"condition_number"`
	Evidence     string    `json:"evidence"`
	CertifiesPSD bool      `json:"certifies_psd"`
}

type MatrixEvaluation struct {
	Matrix               semantic.HermitianMatrix         `json:"matrix"`
	ValueEvidence        semantic.ValueEvidenceKind       `json:"value_evidence"`
	StructuralHermitian  bool                             `json:"structural_hermitian_certified"`
	HermitianDiagnostics []HermitianConsistencyDiagnostic `json:"hermitian_consistency"`
	DirectMatrixChecks   []DirectMatrixCheck              `json:"direct_matrix_checks"`
	EigenDiagnostic      EigenDiagnostic                  `json:"eigen_diagnostic"`
	PSDCertified         bool                             `json:"psd_certified"`
	RHCertified          bool                             `json:"rh_certified"`
}

type M6Result struct {
	M5         M5Result          `json:"-"`
	Registry   *ContractRegistry `json:"-"`
	Evaluation MatrixEvaluation  `json:"evaluation"`
}

type logBoxPair struct {
	q, r     int
	a, b     float64
	halfSum  float64
	halfDiff float64
	h0       float64
}

type entryNumeric struct {
	endpoint semantic.EntryValue
	prime    semantic.EntryValue
	arch     semantic.EntryValue
	total    semantic.EntryValue
	value    complex128
}

func CompileM6() (M6Result, error) {
	m5, err := CompileM5()
	if err != nil {
		return M6Result{}, err
	}
	registry, err := m6TheoremRegistry()
	if err != nil {
		return M6Result{}, err
	}
	matrix := semantic.CloneHermitianMatrix(m5.Matrix)
	for i := range matrix.Entries {
		numeric, err := evaluateLogBoxEntry(matrix.Entries[i])
		if err != nil {
			return M6Result{}, err
		}
		for j := range matrix.Entries[i].Contributions {
			switch matrix.Entries[i].Contributions[j].SourceKind {
			case semantic.EndpointContribution:
				matrix.Entries[i].Contributions[j].Value = numeric.endpoint
			case semantic.PrimePowerContribution:
				matrix.Entries[i].Contributions[j].Value = numeric.prime
			case semantic.ArchimedeanContribution:
				matrix.Entries[i].Contributions[j].Value = numeric.arch
			}
		}
		matrix.Entries[i].Value = numeric.total
		matrix.Entries[i].TransformEvaluations = entryMellinEvaluations(matrix.Entries[i])
	}
	if err := matrix.Validate(); err != nil {
		return M6Result{}, err
	}
	diagnostics := hermitianDiagnostics(matrix)
	checks, err := directMatrixChecks(matrix)
	if err != nil {
		return M6Result{}, err
	}
	eigen := twoByTwoEigenDiagnostic(matrix)
	return M6Result{M5: m5, Registry: registry, Evaluation: MatrixEvaluation{
		Matrix: matrix, ValueEvidence: semantic.ApproximateValue,
		StructuralHermitian:  m5.HermitianCertified,
		HermitianDiagnostics: diagnostics, DirectMatrixChecks: checks,
		EigenDiagnostic: eigen, PSDCertified: false, RHCertified: false,
	}}, nil
}

func m6TheoremRegistry() (*ContractRegistry, error) {
	r := NewContractRegistry()
	contract := TheoremContract{ID: WeilCrossEntryTheoremID, Premises: []Pattern{{Kind: semantic.FunctionalDefinitionKind, Exactness: semantic.Exact}, {Kind: semantic.ExplicitFormulaIdentityKind, Exactness: semantic.Exact}, {Kind: semantic.HermitianFormDefinitionKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.HermitianMatrixDefinitionKind, Exactness: semantic.Exact}, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: lagariasWeilReference, Note: "Lagarias Remark 3.3 identifies the polarized intersection product; Theorem 3.1 supplies its arithmetic representation"}, Trust: TrustedExternalTheorem, Citation: "Lagarias §3, Theorem 3.1 and Remark 3.3"}
	if err := r.Register(contract); err != nil {
		return nil, err
	}
	return r, nil
}

func logBoxParameter(f semantic.TestFunction) (int, error) {
	if f.TransformConvention != semantic.LagariasMellinConvention {
		return 0, fmt.Errorf("evaluator refuses transform convention %s", f.TransformConvention)
	}
	params := map[string]string{}
	for _, p := range f.Parameters {
		params[p.Name] = p.Value
	}
	if params["constructor"] != "centered_multiplicative_log_box" {
		return 0, fmt.Errorf("unsupported Weil basis constructor for %s", f.Symbol)
	}
	q, err := strconv.Atoi(params["q"])
	if err != nil || q < 2 {
		return 0, fmt.Errorf("invalid log-box q for %s", f.Symbol)
	}
	return q, nil
}

func newLogBoxPair(row, column semantic.TestFunction) (logBoxPair, error) {
	q, err := logBoxParameter(row)
	if err != nil {
		return logBoxPair{}, err
	}
	r, err := logBoxParameter(column)
	if err != nil {
		return logBoxPair{}, err
	}
	a, b := 4*math.Log(float64(q)), 4*math.Log(float64(r))
	return logBoxPair{q: q, r: r, a: a, b: b, halfSum: (a + b) / 2, halfDiff: math.Abs(a-b) / 2, h0: math.Min(a, b)}, nil
}

func (p logBoxPair) overlap(u float64) float64 {
	u = math.Abs(u)
	if u <= p.halfDiff {
		return math.Min(p.a, p.b)
	}
	if u < p.halfSum {
		return p.halfSum - u
	}
	return 0
}

func evaluateLogBoxEntry(entry semantic.MatrixEntry) (entryNumeric, error) {
	if entry.TransformConvention != semantic.LagariasMellinConvention {
		return entryNumeric{}, fmt.Errorf("entry %d,%d has incompatible transform", entry.Row, entry.Column)
	}
	pair, err := newLogBoxPair(entry.RowFunction, entry.ColumnFunction)
	if err != nil {
		return entryNumeric{}, err
	}
	endpoint := exactEndpointValue(entry, pair)
	primeValue, terms := primePowerValue(pair)
	prime := approximateValue(entryTarget(entry, "prime-power"), primeValue, semantic.EvaluationMetadata{
		Backend: "go/math float64", BackendVersion: "go.mod toolchain", PrecisionBits: 53, TransformConvention: semantic.LagariasMellinConvention,
		Truncation: &semantic.TruncationInfo{Parameter: "p^n", Bound: fmt.Sprintf("p^n <= (%d*%d)^2 = %d", pair.q, pair.r, pair.q*pair.r*pair.q*pair.r), SummandDefinition: "2*log(p)*(p^n)^-1/2*H_ij(n*log(p))", EnumerationSource: "deterministic Eratosthenes sieve over machine integers", TermsEvaluated: terms, RemainderStatus: "exactly_zero_by_compact_support", SupportExhaustive: true},
		Error:      semantic.ErrorSemantics{Kind: "floating_point_roundoff_unbounded", Notes: "finite support exhausts the mathematical prime-power sum; decimal evaluation is not certified"}, Provenance: entry.TheoremProvenance,
	})
	archValue := archimedeanValue(pair, 1e-11)
	arch := approximateValue(entryTarget(entry, "archimedean"), archValue, semantic.EvaluationMetadata{
		Backend: "go/math adaptive Simpson", BackendVersion: "deterministic recursive implementation", PrecisionBits: 53, TransformConvention: semantic.LagariasMellinConvention,
		Quadrature: &semantic.QuadratureInfo{Method: "adaptive Simpson split at trapezoid kinks", Tolerance: 1e-11, DomainHandling: "[0,(a+b)/2] in log x plus analytic infinite tail", ErrorRigorous: false},
		Error:      semantic.ErrorSemantics{Kind: "heuristic_quadrature_tolerance", Notes: "tolerance is not a theorem-backed error bound"}, Provenance: entry.TheoremProvenance,
	})
	endpointFloat := exactEndpointFloat(pair)
	totalFloat := endpointFloat - primeValue - archValue
	total := approximateValue(entryTarget(entry, "total"), totalFloat, semantic.EvaluationMetadata{
		Backend: "Riemann Go evaluator", BackendVersion: "m6", PrecisionBits: 53, TransformConvention: semantic.LagariasMellinConvention,
		Error: semantic.ErrorSemantics{Kind: "mixed_evidence_approximate", Notes: "exact endpoint plus approximate prime/archimedean values remains approximate"}, Provenance: append(append([]semantic.TheoremID(nil), entry.TheoremProvenance...), WeilCrossEntryTheoremID),
	})
	if semantic.WeakestValueKind(endpoint, prime, arch) != semantic.ApproximateValue {
		return entryNumeric{}, fmt.Errorf("mixed evidence assembly was not conservative")
	}
	return entryNumeric{endpoint: endpoint, prime: prime, arch: arch, total: total, value: complex(totalFloat, 0)}, nil
}

func entryTarget(entry semantic.MatrixEntry, part string) string {
	return fmt.Sprintf("%s[%d,%d].%s", entry.SourceForm, entry.Row, entry.Column, part)
}

func exactEndpointFloat(p logBoxPair) float64 {
	return 8 * float64((p.q*p.q-1)*(p.r*p.r-1)) / float64(p.q*p.r)
}

func exactEndpointValue(entry semantic.MatrixEntry, p logBoxPair) semantic.EntryValue {
	n := int64(8 * (p.q*p.q - 1) * (p.r*p.r - 1))
	d := int64(p.q * p.r)
	g := gcd64(n, d)
	n /= g
	d /= g
	meta := semantic.EvaluationMetadata{SemanticTargetID: entryTarget(entry, "endpoint"), Backend: "exact rational arithmetic", BackendVersion: "m6", PrecisionBits: 0, TransformConvention: semantic.LagariasMellinConvention, Error: semantic.ErrorSemantics{Kind: "exact", ProofObjectKind: semantic.IndependentExactArgument, ProofObject: "M[f_q](0)=M[f_q](1)=2(q-q^-1); exact rational multiplication"}, Provenance: append(append([]semantic.TheoremID(nil), entry.TheoremProvenance...), WeilCrossEntryTheoremID)}
	return semantic.EntryValue{Kind: semantic.ExactValue, DefinitionExact: true, Exact: &semantic.ExactComplex{Real: semantic.ExactScalar{Expression: fmt.Sprintf("%d/%d", n, d), Numerator: n, Denominator: d}, Imag: semantic.ExactScalar{Expression: "0", Numerator: 0, Denominator: 1}}, Metadata: &meta}
}

func gcd64(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	for b != 0 {
		a, b = b, a%b
	}
	if a == 0 {
		return 1
	}
	return a
}

func approximateValue(target string, value float64, meta semantic.EvaluationMetadata) semantic.EntryValue {
	meta.SemanticTargetID = target
	return semantic.EntryValue{Kind: semantic.ApproximateValue, DefinitionExact: true, Approximate: &semantic.ComplexValue{Real: value}, Metadata: &meta}
}

func primePowerValue(p logBoxPair) (float64, int) {
	limit := p.q * p.r * p.q * p.r
	sum, terms := 0.0, 0
	for _, prime := range primesThrough(limit) {
		power := prime
		for power <= limit {
			u := math.Log(float64(power))
			h := math.Exp(-u/2) * p.overlap(u)
			if h > 0 {
				sum += 2 * math.Log(float64(prime)) * h
				terms++
			}
			if power > limit/prime {
				break
			}
			power *= prime
		}
	}
	return sum, terms
}

func primesThrough(n int) []int {
	if n < 2 {
		return []int{}
	}
	sieve := make([]bool, n+1)
	out := []int{}
	for i := 2; i <= n; i++ {
		if sieve[i] {
			continue
		}
		out = append(out, i)
		if i <= n/i {
			for j := i * i; j <= n; j += i {
				sieve[j] = true
			}
		}
	}
	return out
}

func archimedeanValue(p logBoxPair, tol float64) float64 {
	const eulerGamma = 0.577215664901532860606512090082402431
	f := func(u float64) float64 {
		if math.Abs(u) < 1e-10 {
			d := 0.0
			if math.Abs(p.a-p.b) < 1e-14 {
				d = -1
			}
			return d + 1.5*p.h0
		}
		e2 := math.Exp(-2 * u)
		return (2*math.Exp(-u/2)*p.overlap(u) - 2*e2*p.h0) / (1 - e2)
	}
	points := []float64{0, p.halfDiff, p.halfSum}
	sort.Float64s(points)
	integral := 0.0
	for i := 0; i < len(points)-1; i++ {
		if points[i+1]-points[i] > 1e-14 {
			integral += adaptiveSimpson(f, points[i], points[i+1], tol/4, 18)
		}
	}
	tail := p.h0 * math.Log1p(-math.Exp(-2*p.halfSum))
	return (eulerGamma+math.Log(math.Pi))*p.h0 + integral + tail
}

func adaptiveSimpson(f func(float64) float64, a, b, eps float64, depth int) float64 {
	c := (a + b) / 2
	whole := (b - a) * (f(a) + 4*f(c) + f(b)) / 6
	var rec func(float64, float64, float64, float64, float64, float64, int) float64
	rec = func(a, b, fa, fc, fb, whole float64, d int) float64 {
		c := (a + b) / 2
		l := (a + c) / 2
		r := (c + b) / 2
		fl, fr := f(l), f(r)
		left := (c - a) * (fa + 4*fl + fc) / 6
		right := (b - c) * (fc + 4*fr + fb) / 6
		delta := left + right - whole
		if d <= 0 || math.Abs(delta) <= 15*eps {
			return left + right + delta/15
		}
		return rec(a, c, fa, fl, fc, left, d-1) + rec(c, b, fc, fr, fb, right, d-1)
	}
	return rec(a, b, f(a), f(c), f(b), whole, depth)
}

func entryMellinEvaluations(entry semantic.MatrixEntry) []semantic.TransformEvaluation {
	functions := []semantic.TestFunction{entry.RowFunction, entry.ColumnFunction}
	points := []string{"0", "1"}
	out := make([]semantic.TransformEvaluation, 0, 4)
	for _, f := range functions {
		q, _ := logBoxParameter(f)
		n := int64(2 * (q*q - 1))
		d := int64(q)
		g := gcd64(n, d)
		n /= g
		d /= g
		for _, point := range points {
			meta := semantic.EvaluationMetadata{SemanticTargetID: fmt.Sprintf("M[%s](%s)", f.Symbol, point), Backend: "exact rational arithmetic", BackendVersion: "m6", TransformConvention: semantic.LagariasMellinConvention, Error: semantic.ErrorSemantics{Kind: "exact", ProofObjectKind: semantic.IndependentExactArgument, ProofObject: "closed-form integral of centered log box"}, Provenance: []semantic.TheoremID{WeilCrossEntryTheoremID}}
			value := semantic.EntryValue{Kind: semantic.ExactValue, DefinitionExact: true, Exact: &semantic.ExactComplex{Real: semantic.ExactScalar{Expression: fmt.Sprintf("%d/%d", n, d), Numerator: n, Denominator: d}, Imag: semantic.ExactScalar{Expression: "0", Denominator: 1}}, Metadata: &meta}
			out = append(out, semantic.TransformEvaluation{InputFunction: semantic.CloneTestFunction(f), Convention: semantic.LagariasMellinConvention, Point: point, Definition: "M[f](s)=integral_0^infinity f(x)x^s dx/x", Value: value})
		}
	}
	return out
}

func entryComplex(entry semantic.MatrixEntry) complex128 {
	return complex(entry.Value.Approximate.Real, entry.Value.Approximate.Imag)
}

func hermitianDiagnostics(m semantic.HermitianMatrix) []HermitianConsistencyDiagnostic {
	out := make([]HermitianConsistencyDiagnostic, 0, m.Rows*m.Columns)
	for i := 0; i < m.Rows; i++ {
		for j := 0; j < m.Columns; j++ {
			a := entryComplex(m.Entries[i*m.Columns+j])
			b := entryComplex(m.Entries[j*m.Columns+i])
			out = append(out, HermitianConsistencyDiagnostic{Row: i, Column: j, Discrepancy: cmplx.Abs(a - cmplx.Conj(b)), TheoremUse: false})
		}
	}
	return out
}

func directMatrixChecks(m semantic.HermitianMatrix) ([]DirectMatrixCheck, error) {
	vectors := [][]complex128{{1 + 0.5i, -0.75 + 0.25i}, {0.125 - 1.25i, 1.5 + 0.75i}}
	out := make([]DirectMatrixCheck, 0, len(vectors))
	for _, c := range vectors {
		direct, err := directCombinedQuadratic(m.Basis, c, 1e-11)
		if err != nil {
			return nil, err
		}
		matrix := complex(0, 0)
		for i := 0; i < 2; i++ {
			for j := 0; j < 2; j++ {
				matrix += cmplx.Conj(c[i]) * c[j] * entryComplex(m.Entries[i*2+j])
			}
		}
		coeff := make([]semantic.ComplexValue, 2)
		for i, z := range c {
			coeff[i] = semantic.ComplexValue{Real: real(z), Imag: imag(z)}
		}
		out = append(out, DirectMatrixCheck{Coefficients: coeff, Direct: semantic.ComplexValue{Real: real(direct), Imag: imag(direct)}, Matrix: semantic.ComplexValue{Real: real(matrix), Imag: imag(matrix)}, Discrepancy: cmplx.Abs(direct - matrix), Tolerance: 1e-10, Evidence: "numerical_experiment_only"})
	}
	return out, nil
}

// directCombinedQuadratic evaluates the convolution of the sampled linear
// combination first, then applies the three explicit-formula components. It is
// intentionally separate from stored entry assembly.
func directCombinedQuadratic(basis semantic.OrderedBasis, c []complex128, tol float64) (complex128, error) {
	if len(c) != len(basis.Members) {
		return 0, fmt.Errorf("coefficient dimension mismatch")
	}
	pairs := make([][]logBoxPair, len(c))
	maxEnd := 0.0
	for i := range c {
		pairs[i] = make([]logBoxPair, len(c))
		for j := range c {
			p, err := newLogBoxPair(basis.Members[i].Function, basis.Members[j].Function)
			if err != nil {
				return 0, err
			}
			pairs[i][j] = p
			if p.halfSum > maxEnd {
				maxEnd = p.halfSum
			}
		}
	}
	combinedOverlap := func(u float64) complex128 {
		z := complex(0, 0)
		for i := range c {
			for j := range c {
				z += cmplx.Conj(c[i]) * c[j] * complex(pairs[i][j].overlap(u), 0)
			}
		}
		return z
	}
	h0 := combinedOverlap(0)
	endpoint := complex(0, 0)
	for i := range c {
		for j := range c {
			endpoint += cmplx.Conj(c[i]) * c[j] * complex(exactEndpointFloat(pairs[i][j]), 0)
		}
	}
	primeSum := complex(0, 0)
	limit := int(math.Round(math.Exp(maxEnd)))
	for _, prime := range primesThrough(limit) {
		power := prime
		for power <= limit {
			u := math.Log(float64(power))
			h := complex(math.Exp(-u/2), 0) * combinedOverlap(u)
			if cmplx.Abs(h) > 0 {
				primeSum += complex(2*math.Log(float64(prime)), 0) * h
			}
			if power > limit/prime {
				break
			}
			power *= prime
		}
	}
	derivative := complex(0, 0)
	for i := range c {
		for j := range c {
			d := 0.0
			if math.Abs(pairs[i][j].a-pairs[i][j].b) < 1e-14 {
				d = -1
			}
			derivative += cmplx.Conj(c[i]) * c[j] * complex(d, 0)
		}
	}
	integrand := func(u float64) complex128 {
		if math.Abs(u) < 1e-10 {
			return derivative + 1.5*h0
		}
		e2 := math.Exp(-2 * u)
		return (2*complex(math.Exp(-u/2), 0)*combinedOverlap(u) - 2*complex(e2, 0)*h0) / complex(1-e2, 0)
	}
	realIntegral := adaptiveSimpson(func(u float64) float64 { return real(integrand(u)) }, 0, maxEnd, tol, 18)
	imagIntegral := adaptiveSimpson(func(u float64) float64 { return imag(integrand(u)) }, 0, maxEnd, tol, 18)
	tail := complex(math.Log1p(-math.Exp(-2*maxEnd)), 0) * h0
	const eulerGamma = 0.577215664901532860606512090082402431
	arch := complex(eulerGamma+math.Log(math.Pi), 0)*h0 + complex(realIntegral, imagIntegral) + tail
	return endpoint - primeSum - arch, nil
}

func twoByTwoEigenDiagnostic(m semantic.HermitianMatrix) EigenDiagnostic {
	a := real(entryComplex(m.Entries[0]))
	b := real(entryComplex(m.Entries[1]))
	d := real(entryComplex(m.Entries[3]))
	disc := math.Sqrt((a-d)*(a-d) + 4*b*b)
	l1 := (a + d - disc) / 2
	l2 := (a + d + disc) / 2
	condition := math.Inf(1)
	if math.Abs(l1) > 0 {
		condition = math.Abs(l2 / l1)
	}
	return EigenDiagnostic{Eigenvalues: []float64{l1, l2}, Condition: condition, Evidence: "approximate diagnostic only", CertifiesPSD: false}
}
