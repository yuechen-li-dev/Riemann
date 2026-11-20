package claudegrh

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

// TestPinnedSourceUses checks proof-body call sites in the original codeload
// archive without extracting it. This replays provenance, not Lean elaboration.
func TestPinnedSourceUses(t *testing.T) {
	path := os.Getenv("A1_LEAN_ARCHIVE")
	if path == "" {
		t.Skip("set A1_LEAN_ARCHIVE to replay the pinned source trace")
	}
	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		t.Fatal(err)
	}
	const digest = "2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75"
	if got := fmt.Sprintf("%x", hash.Sum(nil)); got != digest {
		t.Fatalf("archive SHA-256 %s, expected %s", got, digest)
	}
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	gz, err := gzip.NewReader(f)
	if err != nil {
		t.Fatal(err)
	}
	defer gz.Close()
	const characterBoxes = "Development/AutomorphicArithmetic/CharacterBoxes/"
	const satake = "Development/AbelianVarietyTheory/Proofs/SatakeParameters/"
	sites := map[string]map[int]string{
		satake + "UnitaryWeightModels.lean":                              {27: "h0, LC, hLC", 29: "LC h0", 30: "E LC", 31: "hLC μ hU u w hg hQ).2"},
		satake + "LocalSystemsLowerBoundsSecondRealPlace.lean":           {53: "fun v j => if", 54: "FOrb D j", 60: "local_model_bound_place_over_second_embedding", 64: "low_term_le_second_real_place_bounds", 66: "hr2"},
		satake + "LocalModelBoundPlaceOverSecondEmbeddingStatement.lean": {12: "GlobalTheory.TwU.LowF.lowAt_F"},
		satake + "RowTwUHolds.lean":                                      {24: "rowTwU_of_models"},
		thirdPlace + "LinkROfIsLinked.lean":                              {76: "hServ_of_linkR"},
		thirdPlace + "GServT.lean":                                       {15: "∃ LC", 16: "LC)", 29: "twU := h.twU"},
		thirdPlace + "ClosedFamilyAdmissibleTwistLink.lean":              {29: "closed_servF_of_servT_of_dGamSize", 31: "admissible_family_of_inclusion"},
		service: {128: "lowS", 201: "⟨LC, hV, hSys⟩", 202: "hSys", 220: "hV hAdm"},
		"Formalism/LocalRingCardinality/DataRows.lean": {56: "muAt_conj", 60: "muAt_or"},
		topWin: {28: "Λ.Gm x ≤ x ^ ((1 : ℝ) / 2 + (eD + eL + 7) * eps)", 75: "exists_small_eps", 77: "h112 eD eL ε", 78: "refine hB"},
		characterBoxes + "ZeroFreeAllOfTopG.lean":                                               {28: "zeroFreeOn_of_tiedWin_cyc", 29: "level_gm_le_on_family_scales"},
		"Development/AutomorphicArithmetic/PrimitiveCosets/LevelGmLeOnFamilyScales.lean":        {30: "windowInequality_of_artinInputs"},
		"Formalism/QuotientMetricDivisors/DivIntoDualShift.lean":                                {144: "∀ φ ∈ N, ∀ F ∈ M w", 147: "∀ φ ∈ N', ∀ x ∈ R"},
		"Formalism/QuotientMetricDivisors/HermitianTwistDivisors/TwistedLatticeNaturality.lean": {722: "refine divisor_into_dual_lattice", 725: "refine shift_pairing_formula"},
		"Formalism/QuotientMetricDivisors/FinalW.lean":                                          {34: ":= divisor_into_dual_shift"},
		"Formalism/LocalToGlobal/IndefiniteForms/PairingConditionOfInverseMatrixIdentity.lean":  {401: "refine hD K", 404: "hL2 hτ hsh", 408: "S.RR_spec hG wf"},
		"Formalism/LocalToGlobal/GroupLattices/AdmissibleAtGeneratorPlaces.lean":                {187: "have hc := h ht hE", 193: "hgain w j hj", 236: "exact admissible_off_twist_places"},
		"Formalism/LocalToGlobal/GroupLattices/ServJointF.lean":                                 {98: "have h' := h u", 102: "exact key φ"},
		division + "ConstantGeTwentyEight.lean":                                                 {191: "intoDual divIntoDualShift latticeGlue"},
		division + "LinkDIVL2R3Col.lean":                                                        {24: "divClausesR_of wf0 wf1"},
		division + "RecBridge.lean":                                                             {637: "linkDIV_L2R_face", 639: "ServJointF.val2AllAtRF", 690: "P.wordFacts"},
		division + "RecBridgeTied.lean":                                                         {203: "linkDIVt_of_rowsF"},
		thirdPlace + "DIVRecBridgeSel.lean":                                                     {196: "linkDIVt_of_rowsT", 205: "DIV.linkDIVt_of_bridge_kernel"},
		thirdPlace + "TopShape.lean":                                                            {47: "have hDIV :=", 52: "exact UpperTheory.Top.topDesWin", 54: "hDIV hgl.1"},
		characterBoxes + "DesWin.lean":                                                          {29: "have hD := aD", 30: "aGL hT hDG", 31: "linkMu_of_linkDIV Fm hD"},
		"Formalism/SharpRestriction/GlobalSchurBounds/TopDensityNGeSixteen.lean":                {716: "obtain ⟨hOne", 723: "hOne hht"},
		thirdPlace + "TopPlug.lean":                                                             {14: "zeroFreeAll_of_topG", 15: "shape_closed_link_n_of_green hG", 18: "zero_free_all_of_green_v_win hG"},
		thirdPlace + "SiegelEisensteinSecondTermHypotheses.lean":                                {72: "every_character_of_green_v_win"},
		thirdPlace + "EveryCharacterOfLiteratureOfVWin.lean":                                    {811: "every_character_of_theta_integrals_local_matching"},
		characterBoxes + "MainConclusion.lean":                                                  {8: "every_character_of_literature_of_vWin"},
		"Supplements/Kernels/MainTheorems.lean":                                                 {7: "FinalTheorems.everyCharacter", 9: "rh_of_every allChi"},
		"Main.lean":                                                                             {4: ":= rhOfAll", 8: ":= allChi"},
	}
	tr := tar.NewReader(gz)
	for len(sites) > 0 {
		item, err := tr.Next()
		if err == io.EOF {
			t.Fatalf("missing files in archive: %v", sites)
		}
		if err != nil {
			t.Fatal(err)
		}
		name := strings.TrimPrefix(item.Name, "riemann-hypothesis-"+Revision+"/src/")
		required, ok := sites[name]
		if !ok {
			continue
		}
		body, err := io.ReadAll(tr)
		if err != nil {
			t.Fatal(err)
		}
		lines := strings.Split(string(body), "\n")
		for line, token := range required {
			if len(lines) < line || !strings.Contains(lines[line-1], token) {
				t.Errorf("%s:%d missing %q", name, line, token)
			}
		}
		delete(sites, name)
	}
}
