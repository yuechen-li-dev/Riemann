package compiler

import (
	"fmt"
	"math/big"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

const (
	M7PSDCertificateID     semantic.ClaimID          = "m7.weil-matrix.psd-certificate"
	M7TwoByTwoPSDTheoremID semantic.TheoremID        = "real-symmetric-2x2-principal-minor-criterion"
	M7CertifyPSDID         semantic.TransformationID = "certified-interval-principal-minors-to-psd"
)

type M7Options struct {
	PrecisionBits  uint
	PanelsPerPiece int
	OmitTail       bool
}

type CertifiedPSDResult struct {
	DiagonalA           semantic.EntryValue `json:"diagonal_a"`
	DiagonalD           semantic.EntryValue `json:"diagonal_d"`
	Determinant         semantic.EntryValue `json:"determinant"`
	APositive           bool                `json:"a_positive"`
	DPositive           bool                `json:"d_positive"`
	DeterminantPositive bool                `json:"determinant_positive"`
	Certified           bool                `json:"certified"`
	Theorem             semantic.TheoremID  `json:"theorem"`
	Criterion           string              `json:"criterion"`
}

type M7Evaluation struct {
	Matrix                           semantic.HermitianMatrix   `json:"matrix"`
	ApproximateMatrix                semantic.HermitianMatrix   `json:"m6_approximate_matrix"`
	ValueEvidence                    semantic.ValueEvidenceKind `json:"value_evidence"`
	StructuralHermitian              bool                       `json:"structural_hermitian_certified"`
	ApproximateContained             []bool                     `json:"m6_approximate_contained"`
	PSD                              CertifiedPSDResult         `json:"finite_psd"`
	FiniteSpanPositivityCertified    bool                       `json:"finite_span_positivity_certified"`
	UniversalWeilPositivityCertified bool                       `json:"universal_weil_positivity_certified"`
	RHCertified                      bool                       `json:"rh_certified"`
	InformationLoss                  string                     `json:"information_loss"`
}

type M7Result struct {
	M6         M6Result          `json:"-"`
	Registry   *ContractRegistry `json:"-"`
	Evaluation M7Evaluation      `json:"evaluation"`
}

func CompileM7() (M7Result, error) { return CompileM7WithOptions(M7Options{}) }

func CompileM7WithOptions(options M7Options) (M7Result, error) {
	m6, err := CompileM6()
	if err != nil {
		return M7Result{}, err
	}
	return compileM7FromM6(m6, options)
}

func compileM7FromM6(m6 M6Result, options M7Options) (M7Result, error) {
	if options.PrecisionBits == 0 {
		options.PrecisionBits = 192
	}
	if options.PanelsPerPiece == 0 {
		options.PanelsPerPiece = 2048
	}
	if options.PrecisionBits < 24 {
		return M7Result{}, fmt.Errorf("certified backend requires at least 24 bits")
	}
	if options.PanelsPerPiece < 4 {
		return M7Result{}, fmt.Errorf("certified backend requires at least four panels per smooth piece")
	}
	if options.OmitTail {
		return M7Result{}, fmt.Errorf("certification rejected: analytic infinite-tail proof object is required")
	}
	ctx := newM7Context(options.PrecisionBits)
	matrix := semantic.CloneHermitianMatrix(m6.Evaluation.Matrix)
	for i := range matrix.Entries {
		// Canonicalize the real off-diagonal enclosure from G12. Hermitian
		// structure, rather than numerical overlap, justifies reuse for G21.
		if i == 2 {
			matrix.Entries[i].Value = structuralConjugateValue(matrix.Entries[1].Value, entryTarget(matrix.Entries[i], "total"))
			for j := range matrix.Entries[i].Contributions {
				switch matrix.Entries[i].Contributions[j].SourceKind {
				case semantic.PrimePowerContribution:
					matrix.Entries[i].Contributions[j].Value = structuralConjugateValue(matrix.Entries[1].Contributions[j].Value, entryTarget(matrix.Entries[i], "prime-power"))
				case semantic.ArchimedeanContribution:
					matrix.Entries[i].Contributions[j].Value = structuralConjugateValue(matrix.Entries[1].Contributions[j].Value, entryTarget(matrix.Entries[i], "archimedean"))
				}
			}
			continue
		}
		entry := &matrix.Entries[i]
		pair, e := newLogBoxPair(entry.RowFunction, entry.ColumnFunction)
		if e != nil {
			return M7Result{}, e
		}
		prime, terms, e := ctx.certifiedPrime(pair)
		if e != nil {
			return M7Result{}, e
		}
		arch, e := ctx.certifiedArchimedean(pair, options.PanelsPerPiece)
		if e != nil {
			return M7Result{}, e
		}
		primeValue := certifiedValue(ctx, entryTarget(*entry, "prime-power"), prime, semantic.OutwardRoundedFiniteSum, "finite support selected by exact integer comparisons; log and exp Taylor remainders plus directed rounding enclose every summand", semantic.EvaluationMetadata{
			Backend: "Riemann small directed big.Float kernel", BackendVersion: "m7", PrecisionBits: int(options.PrecisionBits), TransformConvention: semantic.LagariasMellinConvention,
			Truncation: &semantic.TruncationInfo{Parameter: "p^n", Bound: fmt.Sprintf("p^n < (%d*%d)^2 = %d", pair.q, pair.r, pair.q*pair.r*pair.q*pair.r), SummandDefinition: "2*log(p)*(p^n)^-1/2*H_ij(n*log(p))", EnumerationSource: "deterministic Eratosthenes sieve; region chosen by integer inequalities", TermsEvaluated: terms, RemainderStatus: "exactly_zero_by_compact_support", SupportExhaustive: true}, Provenance: entry.TheoremProvenance})
		tailSem := ctx.semanticInterval(arch.tail)
		archValue := certifiedValue(ctx, entryTarget(*entry, "archimedean"), arch.total, semantic.CertifiedQuadratureBound, "Darboux lower/upper sums from interval range extensions on every deterministic cell, separately required analytic tail, and decimal prefix brackets from NIST DLMF 5.2.3 (Euler gamma) and 3.12.1 (pi)", semantic.EvaluationMetadata{
			Backend: "Riemann piecewise interval Darboux evaluator", BackendVersion: "m7", PrecisionBits: int(options.PrecisionBits), TransformConvention: semantic.LagariasMellinConvention,
			Quadrature: &semantic.QuadratureInfo{Method: "piecewise interval Darboux enclosure", DomainHandling: "u=log(x), cancellation-safe phi ratios on [0,5], then closed-form tail", ErrorRigorous: true, Partitions: arch.panels, Breakpoints: arch.breakpoints, RemainderTheorem: "for every cell I, range enclosure [m_I,M_I] implies |I|[m_I,M_I] encloses its Riemann integral"},
			Tail:       &semantic.TailBound{Start: "u=5", LowerBound: tailSem.RealLowerExact, UpperBound: tailSem.RealUpperExact, Derivation: "after support, integrand=-2*h(1)/(exp(2u)-1); integral from B to infinity is h(1)*log(1-exp(-2B))", ProofKind: semantic.AnalyticTailBound, Exactness: "certified_interval", Provenance: "geometric-series antiderivative and directed log/exp Taylor bounds"}, Provenance: entry.TheoremProvenance})
		endpoint := findContribution(*entry, semantic.EndpointContribution)
		exact, e := exactValueInterval(ctx, endpoint)
		if e != nil {
			return M7Result{}, e
		}
		totalI := ctx.sub(ctx.sub(exact, prime), arch.total)
		total := certifiedValue(ctx, entryTarget(*entry, "total"), totalI, semantic.OutwardRoundedArithmetic, "exact endpoint minus certified prime and archimedean intervals using directed interval arithmetic", semantic.EvaluationMetadata{Backend: "Riemann mixed certified arithmetic", BackendVersion: "m7", PrecisionBits: int(options.PrecisionBits), TransformConvention: semantic.LagariasMellinConvention, Provenance: append(append([]semantic.TheoremID(nil), entry.TheoremProvenance...), WeilCrossEntryTheoremID)})
		for j := range entry.Contributions {
			switch entry.Contributions[j].SourceKind {
			case semantic.PrimePowerContribution:
				entry.Contributions[j].Value = primeValue
			case semantic.ArchimedeanContribution:
				entry.Contributions[j].Value = archValue
			}
		}
		entry.Value = total
	}
	if err := matrix.Validate(); err != nil {
		return M7Result{}, err
	}
	contained := make([]bool, len(matrix.Entries))
	for i := range contained {
		a := m6.Evaluation.Matrix.Entries[i].Value.Approximate.Real
		iv := matrix.Entries[i].Value.Interval
		contained[i] = iv.RealLower <= a && a <= iv.RealUpper
	}
	psd, e := certifyTwoByTwo(ctx, matrix)
	if e != nil {
		return M7Result{}, e
	}
	registry, e := m7TheoremRegistry()
	if e != nil {
		return M7Result{}, e
	}
	if psd.Certified {
		if e = attachM7PSDProof(m6.M5.Graph, matrix, psd); e != nil {
			return M7Result{}, e
		}
	}
	spanCertified, _ := m6.M5.Graph.Certify(M5SpanPositivityID)
	universal, _ := m6.M5.Graph.Certify(WeilPositivityID)
	rh, _ := m6.M5.Graph.Certify(RHClaimID)
	return M7Result{M6: m6, Registry: registry, Evaluation: M7Evaluation{Matrix: matrix, ApproximateMatrix: m6.Evaluation.Matrix, ValueEvidence: semantic.CertifiedInterval, StructuralHermitian: m6.Evaluation.StructuralHermitian, ApproximateContained: contained, PSD: psd, FiniteSpanPositivityCertified: spanCertified, UniversalWeilPositivityCertified: universal, RHCertified: rh, InformationLoss: string(FunctionSpaceRestriction)}}, nil
}

func structuralConjugateValue(v semantic.EntryValue, target string) semantic.EntryValue {
	out := semantic.CloneEntryValue(v)
	if out.Metadata != nil {
		out.Metadata.SemanticTargetID = target
		if out.Metadata.Error.Notes != "" {
			out.Metadata.Error.Notes += "; "
		}
		out.Metadata.Error.Notes += "canonical G12 enclosure reused for G21 by certified Hermitian structure"
	}
	return out
}

func certifiedValue(ctx m7Context, target string, x m7Interval, kind semantic.ProofObjectKind, proof string, meta semantic.EvaluationMetadata) semantic.EntryValue {
	meta.SemanticTargetID = target
	meta.Error = semantic.ErrorSemantics{Kind: string(kind), Bound: "stored outward endpoints", ProofObjectKind: kind, ProofObject: proof}
	iv := ctx.semanticInterval(x)
	return semantic.EntryValue{Kind: semantic.CertifiedInterval, DefinitionExact: true, Interval: &iv, Metadata: &meta}
}

func findContribution(e semantic.MatrixEntry, k semantic.FunctionalContributionKind) semantic.EntryValue {
	for _, c := range e.Contributions {
		if c.SourceKind == k {
			return c.Value
		}
	}
	return semantic.EntryValue{}
}
func exactValueInterval(ctx m7Context, v semantic.EntryValue) (m7Interval, error) {
	if v.Kind != semantic.ExactValue || v.Exact == nil {
		return m7Interval{}, fmt.Errorf("expected exact endpoint")
	}
	return ctx.rational(v.Exact.Real.Numerator, v.Exact.Real.Denominator), nil
}

func certifyTwoByTwo(ctx m7Context, m semantic.HermitianMatrix) (CertifiedPSDResult, error) {
	a := intervalFromSemantic(ctx, *m.Entries[0].Value.Interval)
	b := intervalFromSemantic(ctx, *m.Entries[1].Value.Interval)
	d := intervalFromSemantic(ctx, *m.Entries[3].Value.Interval)
	det := ctx.sub(ctx.mul(a, d), ctx.mul(b, b))
	meta := semantic.EvaluationMetadata{Backend: "directed 2x2 interval arithmetic", BackendVersion: "m7", PrecisionBits: int(ctx.prec), TransformConvention: semantic.LagariasMellinConvention, Provenance: []semantic.TheoremID{M7TwoByTwoPSDTheoremID}}
	detV := certifiedValue(ctx, "G.determinant", det, semantic.OutwardRoundedArithmetic, "[a]*[d]-[b]^2 with all four product endpoints and outward rounding", meta)
	aPos, dPos, detPos := a.lo.Sign() > 0, d.lo.Sign() > 0, det.lo.Sign() > 0
	return CertifiedPSDResult{DiagonalA: m.Entries[0].Value, DiagonalD: m.Entries[3].Value, Determinant: detV, APositive: aPos, DPositive: dPos, DeterminantPositive: detPos, Certified: aPos && dPos && detPos, Theorem: M7TwoByTwoPSDTheoremID, Criterion: "a>=0, d>=0, and ad-b^2>=0 for a real-symmetric 2x2 matrix"}, nil
}

func intervalFromSemantic(ctx m7Context, x semantic.ComplexInterval) m7Interval {
	lo, _, _ := big.ParseFloat(x.RealLowerExact, 10, ctx.prec, big.ToNegativeInf)
	hi, _, _ := big.ParseFloat(x.RealUpperExact, 10, ctx.prec, big.ToPositiveInf)
	return m7Interval{lo, hi}
}

func attachM7PSDProof(g *Graph, m semantic.HermitianMatrix, p CertifiedPSDResult) error {
	prop := semantic.TwoByTwoPrincipalMinorCertificate{MatrixID: m.ID, ALower: p.DiagonalA.Interval.RealLowerExact, DLower: p.DiagonalD.Interval.RealLowerExact, DeterminantLower: p.Determinant.Interval.RealLowerExact, AllStrictlyPositive: p.Certified, ArithmeticProofKinds: []semantic.ProofObjectKind{p.DiagonalA.Metadata.Error.ProofObjectKind, p.DiagonalD.Metadata.Error.ProofObjectKind, p.Determinant.Metadata.Error.ProofObjectKind}}
	ref := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "Riemann M7 directed-interval principal-minor certificate for the real-symmetric 2x2 matrix"}
	claim := semantic.Claim{ID: M7PSDCertificateID, Proposition: prop, Evidence: []semantic.Evidence{{Kind: semantic.CertifiedComputationEvidence, Source: ref, Note: "compiler-verified theorem-backed interval inequalities a>0, d>0, ad-b^2>0"}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: ref}}
	if err := g.AddClaim(claim); err != nil {
		return err
	}
	return g.AddTransformation(Transformation{ID: M7CertifyPSDID, Pass: "apply-certified-2x2-principal-minor-criterion", From: M7PSDCertificateID, Premises: []semantic.ClaimID{M5HermitianPropertyID}, To: M5MatrixPSDID, Relation: Implies, Provenance: ref, Theorem: M7TwoByTwoPSDTheoremID, Trusted: true})
}

func m7TheoremRegistry() (*ContractRegistry, error) {
	r := NewContractRegistry()
	ref := semantic.Reference{Kind: semantic.StandardReference, Citation: "Encyclopedia of Mathematics, Hermitian matrix: a Hermitian matrix is positive semidefinite iff all principal minors are nonnegative; specialized to 2x2", URI: "https://encyclopediaofmath.org/wiki/Hermitian_matrix"}
	c := TheoremContract{ID: M7TwoByTwoPSDTheoremID, Premises: []Pattern{{Kind: semantic.TwoByTwoMinorCertificateKind, Exactness: semantic.Exact}, {Kind: semantic.MatrixPropertyStatementKind, Exactness: semantic.Exact}}, Conclusion: Pattern{Kind: semantic.MatrixPropertyStatementKind, Exactness: semantic.Exact}, Relation: Implies, Evidence: semantic.Evidence{Kind: semantic.KnownTheoremEvidence, Source: ref}, Trust: TrustedExternalTheorem, Citation: ref.Citation}
	if e := r.Register(c); e != nil {
		return nil, e
	}
	return r, nil
}
