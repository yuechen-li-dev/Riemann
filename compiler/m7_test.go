package compiler

import (
	"bytes"
	"math"
	"strings"
	"testing"

	"github.com/yuechen-li-dev/Riemann/semantic"
)

func TestM7CertifiedEntriesEncloseM6AndRetainComponents(t *testing.T) {
	r, err := testM7()
	if err != nil {
		t.Fatal(err)
	}
	for i, e := range r.Evaluation.Matrix.Entries {
		if e.Value.Kind != semantic.CertifiedInterval || e.Value.Interval == nil || !r.Evaluation.ApproximateContained[i] {
			t.Fatalf("entry %d not certified around M6", i)
		}
		if err := e.Value.Validate(); err != nil {
			t.Fatal(err)
		}
		for _, c := range e.Contributions {
			switch c.SourceKind {
			case semantic.EndpointContribution:
				if c.Value.Kind != semantic.ExactValue {
					t.Fatal("endpoint demoted")
				}
			case semantic.PrimePowerContribution:
				if c.Value.Kind != semantic.CertifiedInterval || c.Value.Metadata.Truncation == nil || !c.Value.Metadata.Truncation.SupportExhaustive {
					t.Fatal("prime support/numeric evidence conflated")
				}
			case semantic.ArchimedeanContribution:
				q := c.Value.Metadata.Quadrature
				if c.Value.Kind != semantic.CertifiedInterval || q == nil || !q.ErrorRigorous || len(q.Breakpoints) < 4 || c.Value.Metadata.Tail == nil {
					t.Fatal("archimedean proof objects incomplete")
				}
			}
		}
		approx := r.Evaluation.ApproximateMatrix.Entries[i].Value.Approximate.Real
		if approx < e.Value.Interval.RealLower || approx > e.Value.Interval.RealUpper {
			t.Fatal("M6 value escaped enclosure")
		}
	}
}

func TestM7CertifiedLogAndOutwardArithmetic(t *testing.T) {
	c := newM7Context(128)
	l, e := c.logInt(81)
	if e != nil {
		t.Fatal(e)
	}
	ref := math.Log(81)
	s := c.semanticInterval(l)
	if !(s.RealLower <= ref && ref <= s.RealUpper) {
		t.Fatalf("log(81) not enclosed: %+v", s)
	}
	oneThird := c.rational(1, 3)
	sum := c.add(oneThird, oneThird)
	prod := c.mul(sum, c.rational(3, 2))
	iv := c.semanticInterval(prod)
	if !(iv.RealLower <= 1 && 1 <= iv.RealUpper) {
		t.Fatalf("directed arithmetic lost exact one: %+v", iv)
	}
}

func TestM7TailOmissionAndHeuristicInjectionRejected(t *testing.T) {
	if _, e := CompileM7WithOptions(M7Options{OmitTail: true}); e == nil || !strings.Contains(e.Error(), "tail") {
		t.Fatal("tail omission accepted")
	}
	r, e := testM7()
	if e != nil {
		t.Fatal(e)
	}
	v := semantic.CloneEntryValue(findContribution(r.Evaluation.Matrix.Entries[0], semantic.ArchimedeanContribution))
	v.Metadata.Tail = nil
	if v.Validate() == nil {
		t.Fatal("quadrature without tail remained certifying")
	}
	v = semantic.CloneEntryValue(findContribution(r.Evaluation.Matrix.Entries[0], semantic.ArchimedeanContribution))
	v.Metadata.Error.ProofObjectKind = "heuristic_tolerance"
	if v.Validate() == nil {
		t.Fatal("heuristic interval accepted")
	}
	g := NewGraph()
	ref := semantic.Reference{Kind: semantic.CompilerRecord, Citation: "test"}
	claim := semantic.Claim{ID: "fake-certified-computation", Proposition: semantic.NamedObligation{Name: "not a recognized numerical certificate"}, Evidence: []semantic.Evidence{{Kind: semantic.CertifiedComputationEvidence, Source: ref}}, Exactness: semantic.Exact, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: ref}}
	if err := g.AddClaim(claim); err != nil {
		t.Fatal(err)
	}
	if ok, _ := g.Certify(claim.ID); ok {
		t.Fatal("unknown certified-computation proposition laundered into theorem evidence")
	}
}

func TestM7HermitianPSDSpanAndRHBoundaries(t *testing.T) {
	r, e := testM7()
	if e != nil {
		t.Fatal(e)
	}
	if !r.Evaluation.StructuralHermitian || !r.Evaluation.PSD.APositive || !r.Evaluation.PSD.DPositive || !r.Evaluation.PSD.DeterminantPositive || !r.Evaluation.PSD.Certified {
		t.Fatalf("finite PSD not certified: %+v", r.Evaluation.PSD)
	}
	d := r.Evaluation.PSD.Determinant.Interval
	approx := 0.0918849012616*0.0877271147825 - math.Pow(-0.0113156378917, 2)
	if approx < d.RealLower || approx > d.RealUpper {
		t.Fatal("determinant enclosure is not conservative")
	}
	if !r.Evaluation.FiniteSpanPositivityCertified || r.Evaluation.UniversalWeilPositivityCertified || r.Evaluation.RHCertified || r.Evaluation.InformationLoss != "function_space_restriction" {
		t.Fatal("finite theorem crossed universal boundary")
	}
	if c, _ := r.M6.M5.Graph.Certify(M5MatrixPSDID); !c {
		t.Fatal("exact PSD theorem not attached to graph")
	}
}

func TestM7WideAndReducedPrecisionIntervalsStayOpen(t *testing.T) {
	r, e := testM7()
	if e != nil {
		t.Fatal(e)
	}
	m := semantic.CloneHermitianMatrix(r.Evaluation.Matrix)
	for _, i := range []int{0, 3} {
		m.Entries[i].Value.Interval.RealLower = -1
		m.Entries[i].Value.Interval.RealUpper = 1
		m.Entries[i].Value.Interval.RealLowerExact = "-1"
		m.Entries[i].Value.Interval.RealUpperExact = "1"
	}
	p, e := certifyTwoByTwo(newM7Context(128), m)
	if e != nil {
		t.Fatal(e)
	}
	if p.Certified {
		t.Fatal("wide interval guessed PSD")
	}
	coarse, e := CompileM7WithOptions(M7Options{PrecisionBits: 24, PanelsPerPiece: 16})
	if e != nil {
		t.Fatal(e)
	}
	if coarse.Evaluation.PSD.Certified {
		t.Fatal("artificially weak precision/partition unexpectedly certified PSD")
	}
}

func TestM7BreakpointPartitionAndBasisSwap(t *testing.T) {
	r, e := testM7()
	if e != nil {
		t.Fatal(e)
	}
	f2 := r.Evaluation.Matrix.Basis.Members[0].Function
	f3 := r.Evaluation.Matrix.Basis.Members[1].Function
	p, e := newLogBoxPair(f2, f3)
	if e != nil {
		t.Fatal(e)
	}
	rev, e := newLogBoxPair(f3, f2)
	if e != nil {
		t.Fatal(e)
	}
	ctx := newM7Context(96)
	a, e := ctx.certifiedArchimedean(p, 32)
	if e != nil {
		t.Fatal(e)
	}
	z, e := ctx.certifiedArchimedean(rev, 32)
	if e != nil {
		t.Fatal(e)
	}
	if a.total.lo.Cmp(z.total.lo) != 0 || a.total.hi.Cmp(z.total.hi) != 0 {
		t.Fatal("basis swap changed canonical enclosure")
	}
	coarse := ctx.semanticInterval(a.total)
	approx := archimedeanValue(p, 1e-12)
	if approx < coarse.RealLower || approx > coarse.RealUpper {
		t.Fatal("perturbed subdivision lost the direct numerical value")
	}
	if len(a.breakpoints) != 5 || !strings.Contains(a.breakpoints[1], "..") || !strings.Contains(a.breakpoints[2], "..") {
		t.Fatal("certified kink guards absent")
	}
}

func TestM7ReportsDeterministicAndExposeProofBoundary(t *testing.T) {
	r, e := testM7()
	if e != nil {
		t.Fatal(e)
	}
	a, e := M7JSONReport(r)
	if e != nil {
		t.Fatal(e)
	}
	b, e := M7JSONReport(r)
	if e != nil {
		t.Fatal(e)
	}
	if !bytes.Equal(a, b) {
		t.Fatal("JSON nondeterministic")
	}
	h := M7HumanReport(r)
	for _, want := range []string{"riemann.semantic-graph.m7", "certified_quadrature_bound", "analytic_tail_bound", "support_exhaustive", "finite_span_positivity_certified", "function_space_restriction", "RH\n  unresolved"} {
		if !strings.Contains(string(a)+h, want) {
			t.Fatalf("report omits %q", want)
		}
	}
}
