package semantic

import "testing"

func TestFunctionClassCanonicalIdentity(t *testing.T) {
	a := WeilNiceClass()
	b := WeilNiceClass()
	b.Constraints[0], b.Constraints[3] = b.Constraints[3], b.Constraints[0]
	if a.Key() != b.Key() {
		t.Fatal("constraint order changed function-class identity")
	}
	f1 := TestFunction{Symbol: "f1", Kind: BasisTestFunction, DeclaredClass: a.ID, RequiredAttributes: a.Constraints, TransformConvention: a.TransformConvention, Parameters: []FunctionParameter{{Name: "i", Value: "1"}}}
	f2 := f1
	f2.Symbol = "f2"
	c := FunctionClass{ID: "finite", Kind: FiniteFunctionClass, Constraints: a.Constraints, TransformConvention: a.TransformConvention, Members: []TestFunction{f1, f2}}
	d := CloneFunctionClass(c)
	d.Members[0], d.Members[1] = d.Members[1], d.Members[0]
	if c.Key() != d.Key() {
		t.Fatal("member order changed finite-family identity")
	}
	d.Members[0].TransformConvention = "other"
	if c.Key() == d.Key() {
		t.Fatal("transform mismatch disappeared from identity")
	}
}

func TestFunctionalStatementValidationRejectsTransformMismatch(t *testing.T) {
	p := UniversalFunctionalStatement{Quantifier: ForAll, Variable: "f", FunctionClass: WeilNiceClass(), Functional: WeilZetaQuadraticFunctional, Predicate: FunctionalNonnegative, TransformConvention: "wrong"}
	if err := p.Validate(); err == nil {
		t.Fatal("transform convention mismatch accepted")
	}
}

func TestWeilFunctionalZeroAggregateIsTyped(t *testing.T) {
	q := QuadraticFunctional{ID: WeilZetaQuadraticFunctional, Object: RiemannZeta, InputConstruction: "h", TransformConvention: LagariasMellinConvention, Contributions: []FunctionalContribution{{Kind: ZeroContribution, RepresentationSide: "zero_side", Sign: 1, Formula: "sum", Aggregate: &Aggregate{Kind: SumOverDomain, IndexDomain: NontrivialZeros(RiemannZeta), Summand: "term", Convergence: []string{"symmetric order"}, TransformConvention: LagariasMellinConvention, TheoremLineage: []TheoremID{"explicit-formula"}, Provenance: Reference{Kind: StandardReference, Citation: "test"}}}}}
	if err := q.Validate(); err != nil {
		t.Fatal(err)
	}
	bad := CloneQuadraticFunctional(q)
	bad.Contributions[0].Aggregate.IndexDomain = CriticalStrip()
	if err := bad.Validate(); err == nil {
		t.Fatal("untyped zero aggregate domain accepted")
	}
}
