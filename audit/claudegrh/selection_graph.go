package claudegrh

import (
	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

const admissible = "Formalism/LocalToGlobal/GroupLattices/AdmissibleAtGeneratorPlaces.lean"
const service = "Formalism/SharpRestriction/SmoothDesignExistence/SumOverPrimesHypotheses.lean"

var (
	selection   = statement{"a2.eventual_modelsOK_exists_LC_cmpLite_and_siteAdm", all}
	shape       = statement{"a2.eventual_orbit_shape_same_family_and_orbit_letter_bounds", all}
	commonLC    = statement{"a2.exists_LC_cmpV_and_sysTw", all}
	cmpOnly     = statement{"a2.exists_LC_cmpV", all}
	sysOnly     = statement{"a2.exists_LC_sysTw", all}
	serviceData = statement{"a2.same_selection_pdv_dgam_and_service_auxiliaries", all}
	columnBound = statement{"a2.conclusion1_on_same_family_tail", all}
	fill        = statement{"a2.parameter_fills_and_cs_drop_same_family", all}
	gainBound   = statement{"a2.gainLe_on_same_family_tail", all}
	offInputs   = statement{"a2.off_twist_period_integrality_ties_membership_and_errors", all}
	offPeriod   = statement{"a2.unrounded_selected_period_off_level_and_generators", all}
	otherPlaces = statement{"a2.level_generator_and_family_threshold_inputs", all}
)
var selectionRules = []rule{
	{commonLC, []statement{selection, shape}, "Development/AbelianVarietyTheory/Proofs/SatakeParameters/UnitaryWeightModels.lean#L25", "rowTwU_of_models (witness used by SelectionServed.twU)"},
	{cmpOnly, []statement{commonLC}, thirdPlace + "GServT.lean#L15", "first existential projection of twU"},
	{sysOnly, []statement{commonLC}, thirdPlace + "GServT.lean#L16", "second existential projection of twU"},
	{columnBound, []statement{commonLC, serviceData}, service + "#L187", "sum_over_primes_service and Serv projection"},
	{gainBound, []statement{fill}, "Formalism/SharpRestriction/Divisibility/WeightDegreeSizes.lean#L931", "gainLe_of_fills"},
	{offPeriod, []statement{columnBound, gainBound, offInputs}, admissible + "#L171", "admissible_off_twist_places on the common family tail"},
	{periods, []statement{offPeriod, otherPlaces}, admissible + "#L225", "admissible_level_two and famU"},
}

// SimultaneousGraph expands the period premise in the A1 source graph.
func SimultaneousGraph() (*compiler.Graph, error) { return proofGraph(true) }

func addSelectionTrace(g *compiler.Graph) error {
	for _, x := range []struct {
		s    statement
		path string
	}{
		{selection, "Development/AbelianVarietyTheory/Definitions/GaloisGeometry/ModelsOK.lean#L56"},
		{shape, "Development/AbelianVarietyTheory/Proofs/SatakeParameters/RowTwUHolds.lean#L21"}, {serviceData, service + "#L187"},
		{fill, admissible + "#L293"}, {offInputs, admissible + "#L173"}, {otherPlaces, admissible + "#L214"},
	} {
		src := ref(x.path, string(x.s.predicate))
		if err := g.AddClaim(semantic.Claim{ID: x.s.id(), Proposition: x.s.proposition(), Exactness: semantic.Exact, Evidence: []semantic.Evidence{{Kind: semantic.UnverifiedConjectureEvidence, Source: src, Note: "A2 open source premise"}}, Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: src}}); err != nil {
			return err
		}
	}
	for _, r := range selectionRules {
		if err := r.applyIn(g, "A2", r.inputs...); err != nil {
			return err
		}
	}
	return nil
}
