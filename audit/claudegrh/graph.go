// The audit overlay binds only the source statements listed below. Its scope
// tags distinguish a finite replay, an assumed bound and a universal premise.
package claudegrh

import (
	"fmt"
	"slices"

	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

const Revision = "a172c464fd4bcfc4f53e394dcba45d742b1409a6"
const source = "https://github.com/Deicyde/riemann-hypothesis/blob/" + Revision + "/src/"
const thirdPlace = "Development/AutomorphicArithmetic/ThetaCorrespondence/ThetaPairLifts/ThirdPlaceSeries/"
const division = "Formalism/LocalToGlobal/Divisibility/"
const topWin = "Formalism/SharpRestriction/RamificationFrobenius/TotientGaloisClasses/TopWin.lean"

type predicate string
type scope string

const (
	all       scope     = "all"
	instance  scope     = "instance"
	assuming  scope     = "under_assumption"
	pointwise predicate = "pairing_valuation"
)

type statement struct {
	predicate predicate
	scope     scope
}

func (s statement) id() semantic.ClaimID {
	return semantic.ClaimID(fmt.Sprintf("a1.%s.%s", s.predicate, s.scope))
}

func (s statement) proposition() semantic.ExternalStatement {
	p := semantic.ExternalStatement{Predicate: "a1." + string(s.predicate), Revision: Revision}
	switch s {
	case finite:
		p.Instance = "Q, row=(2,4), p=2, sigma=tau=1"
	case average:
		p.Instance = "Q, row=(1,3), p=2, sigma=tau=1"
	}
	return p
}

var (
	lifts        = statement{"theta_at_each_legal_stage", all}
	periods      = statement{"selected_tuple_period_bound_adm2", all}
	ties         = statement{"tied_columns_cut_and_lattice_membership", all}
	local        = statement{pointwise, all}
	finite       = statement{pointwise, instance}
	conditional  = statement{pointwise, assuming}
	average      = statement{"normalized_sum_valuation", instance}
	shiftData    = statement{"unit_pairing_identity_localizations_and_shift", all}
	dual         = statement{"shifted_dual_integrality", all}
	cost         = statement{"accumulated_word_height_bound", all}
	globalInputs = statement{"remaining_design_and_global_inputs", all}
	window       = statement{"eventual_tied_family_window", all}
	hosts        = statement{"primitive_character_hosts_and_compatible_windows", all}
	comparison   = statement{"chebyshev_comparison_for_tied_families", all}
	estimate     = statement{"chebyshev_growth_bound_on_selected_windows", all}
	analytic     = statement{"form_B_on_compatible_windows_and_imprimitive_extension", all}
	region       = statement{"dirichlet_right_half_plane", all}
	rh           = statement{"riemann_hypothesis", all}
)

// A step can compress intervening Lean calls. Its other hypotheses are
// explicit graph inputs, with pinned sources for the complete contracts.
type rule struct {
	output     statement
	inputs     []statement
	path, name string
}

var route = []rule{
	{local, []statement{periods, ties}, "Formalism/LocalToGlobal/GroupLattices/ServJointF.lean#L84", "val2AllAtRF"},
	{dual, []statement{local, shiftData}, "Formalism/QuotientMetricDivisors/HermitianTwistDivisors/TwistedLatticeNaturality.lean#L720", "divisor_into_dual_shift"},
	{window, []statement{dual, cost, globalInputs}, "Development/AutomorphicArithmetic/CharacterBoxes/DesWin.lean#L11", "topDesWin (via wordFacts, LinkDIV and LinkMu)"},
	{estimate, []statement{window, hosts, comparison}, topWin + "#L77", "level_of_tiedWin application of h112 at the selected family"},
	{region, []statement{estimate, hosts, analytic}, thirdPlace + "TopPlug.lean#L11", "zero_free_all_of_green_v_win via level_of_tiedWin and imprimitive"},
	{rh, []statement{region}, "Supplements/Kernels/MainTheorems.lean#L9", "rhOfAll via rh_of_every"},
}

func ref(path, name string) semantic.Reference {
	return semantic.Reference{Kind: semantic.StandardReference, Citation: name, URI: source + path}
}

// apply uses only the recorded signature, even when a caller proposes a
// differently scoped bound with the same underlying pairing predicate.
func (r rule) apply(g *compiler.Graph, args ...statement) error {
	if !slices.Equal(args, r.inputs) {
		return fmt.Errorf("%s requires %v, got %v", r.name, r.inputs, args)
	}
	parents := make([]semantic.ClaimID, len(args))
	for i, arg := range args {
		parents[i] = arg.id()
		c, ok := g.Claim(parents[i])
		if !ok || c.Proposition != arg.proposition() || len(c.Assumptions) != 0 {
			return fmt.Errorf("%s lacks premise %s", r.name, parents[i])
		}
	}
	src := ref(r.path, r.name)
	id := r.output.id()
	tr := semantic.TransformationID("use." + string(id))
	proof := semantic.TheoremID(r.name)
	if err := g.AddClaim(semantic.Claim{
		ID: id, Proposition: r.output.proposition(), Exactness: semantic.Exact,
		Evidence:   []semantic.Evidence{{Kind: semantic.DerivedEvidence, Source: src, Note: "source use; application premises open"}},
		Provenance: semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: tr, Theorem: proof, Source: src},
	}); err != nil {
		return err
	}
	return g.AddTransformation(compiler.Transformation{ID: tr, Pass: "A1", From: parents[0], Premises: parents[1:], To: id,
		Relation: compiler.Implies, Theorem: proof, Trusted: false, Provenance: src})
}

// ProofGraph constructs an isolated graph. Source use, even with a paper
// reconstruction of a rule, cannot certify a concrete application in Go.
func ProofGraph() (*compiler.Graph, error) {
	g := compiler.NewGraph()
	add := func(s statement, src semantic.Reference, kind semantic.EvidenceKind, note string) error {
		c := semantic.Claim{ID: s.id(), Proposition: s.proposition(), Exactness: semantic.Exact,
			Evidence: []semantic.Evidence{{Kind: kind, Source: src, Note: note}}, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: src}}
		if s.scope == assuming {
			c.Assumptions = []semantic.Assumption{{ID: "a1.assume_period_bound", Description: "Assume the required valuation for every local period."}}
			c.Evidence[0].Note = string(Conditional) + ": " + c.Assumptions[0].Description
		}
		return g.AddClaim(c)
	}
	for _, x := range []struct {
		s          statement
		path, name string
	}{
		{lifts, thirdPlace + "ThetaLiftOfReference.lean#L37", "exists_theta_at_of_third_layer (published third-layer and local Langlands inputs, legal j)"},
		{periods, "Formalism/SharpRestriction/Divisibility/LocalQuadraticParts.lean", "adm2 (word, place, selected tuple, every admissible pair)"},
		{ties, "Formalism/LocalToGlobal/GroupLattices/ServJointF.lean#L84", "val2AllAtRF remaining hypotheses (tied columns, cut and lattice membership)"},
		{shiftData, "Formalism/QuotientMetricDivisors/DivIntoDualShift.lean#L133", "DivIntoDualShift (unit scalar, identity, rounding and both localization equalities)"},
		{cost, division + "RecBridge.lean#L261", "rhs_le_record and wordFacts (compatible numerical records and uniform error constants)"},
		{globalInputs, "Development/AutomorphicArithmetic/CharacterBoxes/DesWin.lean#L11", "remaining topDesWin inputs, including design existence, HypV, LinkN/Eta/IRR, rank, LinkS/G/GLOBAL and remaining LinkDIV clauses"},
		{hosts, "Development/AutomorphicArithmetic/CharacterBoxes/ZeroFreeAllOfTopG.lean#L26", "TopG at scale 3/2, host supply for every primitive character, EFamC and Windows for the same chosen family"},
		{comparison, "Development/AutomorphicArithmetic/PrimitiveCosets/LevelGmLeOnFamilyScales.lean#L28", "level_gm_le_on_family_scales (ChebyshevInequalityOnWindows)"},
		{analytic, topWin + "#L68", "level_of_tiedWin (FormBAt and arbitrary small positive epsilon); imprimitive extension in TopPlug"},
		{conditional, "Formalism/SharpRestriction/Divisibility/LocalQuadraticParts.lean", "assumed period bound"},
	} {
		if err := add(x.s, ref(x.path, x.name), semantic.UnverifiedConjectureEvidence, "open in Riemann"); err != nil {
			return nil, err
		}
	}
	// These two observations concern different, replayed concrete fixtures.
	for _, x := range []struct {
		s       statement
		row     []string
		outcome Outcome
	}{
		{finite, []string{"2", "4"}, FiniteVerified}, {average, []string{"1", "3"}, ReducedCounterexample},
	} {
		got, err := Replay(Instance{x.row, []Step{{Prime: 2, Sigma: "1", Tau: 1}}})
		if err != nil || got.Outcome != x.outcome {
			return nil, fmt.Errorf("graph fixture %v: %v (%v)", x.row, got, err)
		}
		src := semantic.Reference{Kind: semantic.ExperimentRecord, Citation: fmt.Sprintf("A1 Go replay, p=2 sigma=tau=1 row=%v", x.row)}
		if err := add(x.s, src, semantic.CertifiedComputationEvidence, string(got.Outcome)); err != nil {
			return nil, err
		}
	}
	for _, r := range route {
		if err := r.apply(g, r.inputs...); err != nil {
			return nil, err
		}
	}
	return g, nil
}
