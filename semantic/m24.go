package semantic

import (
	"fmt"
	"math/big"
)

// LegalTwoAtomRadiusComponent names only the connected legal radius set that
// contains the M23 root. It deliberately does not broaden the claim to every
// contact-preserving radius or to every two-atom measure.
type LegalTwoAtomRadiusComponent struct {
	Definition             string                  `json:"definition"`
	SupportLegality        string                  `json:"support_legality"`
	MassLegality           string                  `json:"mass_legality"`
	WeightLegality         string                  `json:"weight_legality"`
	ObjectiveLegality      string                  `json:"objective_legality"`
	FamilyWellDefinedness  string                  `json:"family_well_definedness"`
	ResponseSignStructure  string                  `json:"response_sign_structure"`
	CertifiedInterior      CertifiedScalarInterval `json:"certified_interior"`
	BoundaryClassification string                  `json:"boundary_classification"`
	ExactEndpointStatus    string                  `json:"exact_endpoint_status"`
}

type ContactBranch struct {
	ID               string `json:"id"`
	RadiusRange      string `json:"radius_range"`
	ContactRange     string `json:"contact_range"`
	EnvelopeBehavior string `json:"envelope_behavior"`
	ProofStatus      string `json:"proof_status"`
	NextCompetitor   string `json:"next_competitor"`
}

type GlobalEnvelopeEvidence struct {
	NegativeCosineRegion string                  `json:"negative_cosine_region"`
	FixedContactFormula  string                  `json:"fixed_contact_formula"`
	LargeRadiusCutoff    string                  `json:"large_radius_cutoff"`
	LargeRadiusUpper     CertifiedScalarInterval `json:"large_radius_upper"`
	GeneratedCellCount   int                     `json:"generated_cell_count"`
	GeneratedCellHash    string                  `json:"generated_cell_hash"`
	CompactRadiusCover   string                  `json:"compact_radius_cover"`
	BestCellUpper        CertifiedScalarInterval `json:"best_cell_upper"`
	UnresolvedRadius     CertifiedScalarInterval `json:"unresolved_radius"`
	UnresolvedContact    CertifiedScalarInterval `json:"unresolved_contact_research"`
	EvidenceBoundary     string                  `json:"evidence_boundary"`
}

type ComponentOptimumCertificate struct {
	Status                    string                      `json:"status"`
	Component                 LegalTwoAtomRadiusComponent `json:"legal_radius_component"`
	Optimizer                 CertifiedOptimizerBox       `json:"m23_optimizer"`
	Branches                  []ContactBranch             `json:"branch_atlas"`
	GlobalEvidence            GlobalEnvelopeEvidence      `json:"global_envelope_evidence"`
	M23ControlsComponent      bool                        `json:"m23_controls_component"`
	OptimizerGloballyFeasible bool                        `json:"optimizer_globally_feasible"`
	MissingTheorem            string                      `json:"missing_theorem"`
}

// m24RadiusBarrierCell is the stable runtime boundary for the boring table
// emitted by experiments/m24_two_atom_component.octest. K denotes the exact
// radius cell [3+K/1000, 3+(K+1)/1000]. Contact is an exact decimal rational.
type m24RadiusBarrierCell struct {
	K       int
	Contact string
}

const m24BarrierLogicalHash = "sha256:68592d9c19e79e237c1ce5f7f3f74aa45ec81b3011fe778604005f05abbbf4af"

func m24C0() rationalInterval {
	one := riPoint(big.NewRat(1, 1))
	return riAdd(
		riSub(one, riMul(riPoint(big.NewRat(2, 1)), riReciprocal(m18Pi))),
		riMul(m18Pi, riPoint(big.NewRat(1, 4))),
	)
}

func m24RadiusCell(k int) rationalInterval {
	denominator := big.NewRat(1000, 1)
	left := new(big.Rat).Add(big.NewRat(3, 1), new(big.Rat).Quo(big.NewRat(int64(k), 1), denominator))
	right := new(big.Rat).Add(big.NewRat(3, 1), new(big.Rat).Quo(big.NewRat(int64(k+1), 1), denominator))
	return ri(left, right)
}

func m24BaseAndResponse(r, t rationalInterval) (rationalInterval, rationalInterval, rationalInterval, rationalInterval, rationalInterval) {
	one := riPoint(big.NewRat(1, 1))
	two := riPoint(big.NewRat(2, 1))
	piHalf := riMul(m18Pi, riPoint(big.NewRat(1, 2)))
	phase := riMul(piHalf, r)
	phaseCos := m23TrigRange(phase, true)
	phaseSin := m23TrigRange(phase, false)
	phaseCosPrime := riNeg(riMul(piHalf, phaseSin))
	sinT := m23TrigRange(t, false)
	cosT := m23TrigRange(t, true)
	sincT := riMul(sinT, riReciprocal(t))
	coefficient := riSub(
		riMul(riMul(two, phaseCos), riReciprocal(m18Pi)),
		riMul(r, phaseSin),
	)
	coefficientPrime := riSub(
		riNeg(riMul(two, phaseSin)),
		riMul(riMul(piHalf, r), phaseCos),
	)
	response := riAdd(
		riAdd(
			riNeg(riMul(riMul(m18Pi, phaseCos), sincT)),
			riMul(riMul(two, coefficient), cosT),
		),
		riMul(two, m23TrigRange(riMul(r, t), true)),
	)
	responsePrime := riAdd(
		riAdd(
			riNeg(riMul(riMul(m18Pi, phaseCosPrime), sincT)),
			riMul(riMul(two, coefficientPrime), cosT),
		),
		riNeg(riMul(riMul(two, t), m23TrigRange(riMul(r, t), false))),
	)
	w0 := riAdd(
		riSub(riReciprocal(m18Pi), riMul(riPoint(big.NewRat(4, 1)), riPow(m18Pi, -2))),
		riMul(riPoint(big.NewRat(8, 1)), riPow(m18Pi, -3)),
	)
	base := riAdd(
		riAdd(
			one,
			riMul(two, riMul(riSub(cosT, one), riPow(t, -2))),
		),
		riAdd(
			riMul(riMul(two, riSub(one, m24C0())), sincT),
			riMul(riMul(two, w0), cosT),
		),
	)
	return base, response, responsePrime, phaseCos, phaseCosPrime
}

func m24PointBarrierUpper(cell m24RadiusBarrierCell) (rationalInterval, error) {
	if cell.K < 0 || cell.K >= 620 {
		return rationalInterval{}, fmt.Errorf("M24 radius cell index out of range: %d", cell.K)
	}
	t := m23Rat(cell.Contact)
	radiusCell := m24RadiusCell(cell.K)
	base, response, responsePrime, phaseCos, phaseCosPrime := m24BaseAndResponse(radiusCell, riPoint(t))
	if response.hi.Sign() >= 0 {
		return rationalInterval{}, fmt.Errorf("M24 cell %d test response does not stay negative", cell.K)
	}
	if base.lo.Sign() < 0 {
		return rationalInterval{}, fmt.Errorf("M24 cell %d base spectrum enclosure is negative", cell.K)
	}
	// On the compact cover r is in [3,181/50], where cos(pi*r/2)>=0.
	// Cell zero includes the exact zero at r=3; every later cell must certify
	// strict positivity directly.
	if phaseCos.hi.Sign() < 0 || (cell.K > 0 && phaseCos.lo.Sign() <= 0) {
		return rationalInterval{}, fmt.Errorf("M24 cell %d has the wrong mass or phase sign", cell.K)
	}
	centerRadius := riPoint(new(big.Rat).Quo(new(big.Rat).Add(radiusCell.lo, radiusCell.hi), big.NewRat(2, 1)))
	centerBase, centerResponse, _, centerCos, _ := m24BaseAndResponse(centerRadius, riPoint(t))
	centerMass := riMul(centerBase, riReciprocal(riNeg(centerResponse)))
	centerObjective := riAdd(
		m24C0(),
		riMul(riMul(m18Pi, riPoint(big.NewRat(1, 2))), riMul(centerCos, centerMass)),
	)
	// With t fixed at the Oct-generated cell witness,
	// d[c0+(pi/2)F0(t)A(r)/(-H(r,t))]/dr
	//   =(pi/2)F0(t)(A*H_r-A'*H)/H^2.
	numerator := riSub(riMul(phaseCos, responsePrime), riMul(phaseCosPrime, response))
	derivative := riMul(
		riMul(riMul(m18Pi, riPoint(big.NewRat(1, 2))), base),
		riMul(numerator, riReciprocal(riMul(response, response))),
	)
	derivativeAbs := new(big.Rat).Abs(derivative.lo)
	if absHi := new(big.Rat).Abs(derivative.hi); absHi.Cmp(derivativeAbs) > 0 {
		derivativeAbs = absHi
	}
	errorRadius := riPoint(new(big.Rat).Mul(derivativeAbs, big.NewRat(1, 2000)))
	errorInterval := ri(new(big.Rat).Neg(new(big.Rat).Set(errorRadius.hi)), new(big.Rat).Set(errorRadius.hi))
	return riAdd(centerObjective, errorInterval), nil
}

func m24LargeRadiusUpper() (rationalInterval, rationalInterval, error) {
	zero := riPoint(big.NewRat(0, 1))
	contact := ri(
		new(big.Rat).Quo(m18Pi.lo, big.NewRat(2, 1)),
		new(big.Rat).Quo(m18Pi.hi, big.NewRat(2, 1)),
	)
	baseCurvature := m23Second([3]rationalInterval{zero, riPoint(big.NewRat(4, 1)), contact}).v
	if baseCurvature.lo.Sign() <= 0 {
		return rationalInterval{}, rationalInterval{}, fmt.Errorf("M24 base fixed-contact curvature is not positive")
	}
	cutoff := riPoint(big.NewRat(181, 50))
	denominator := riAdd(
		riSub(riMul(cutoff, cutoff), riPoint(big.NewRat(1, 1))),
		riMul(riPoint(big.NewRat(8, 1)), riPow(m18Pi, -2)),
	)
	// F_tt(pi/2)=F0_tt(pi/2)-2*m*cos(pi*r/2)*(r^2-1+8/pi^2).
	// For cos>0, feasibility forces the displayed gain bound; for cos<=0,
	// nonnegative mass already gives c<=c0.
	gainUpper := riMul(
		riMul(m18Pi, riPoint(big.NewRat(1, 4))),
		riMul(baseCurvature, riReciprocal(denominator)),
	)
	return riAdd(m24C0(), gainUpper), baseCurvature, nil
}

// CertifyM24ComponentControl proves that every legal family point outside one
// nine-thousandth-wide radius window is strictly below the M23 optimizer's
// certified lower objective. The remaining window is intentionally not
// promoted from numerical branch evidence to a theorem.
func CertifyM24ComponentControl() (ComponentOptimumCertificate, error) {
	_, optimizer, _, err := CertifyVariableRadiusOptimizer()
	if err != nil {
		return ComponentOptimumCertificate{}, err
	}
	target := m23Rat(optimizer.C.Lower)
	c0 := m24C0()
	if c0.hi.Cmp(target) >= 0 {
		return ComponentOptimumCertificate{}, fmt.Errorf("M24 negative-cosine barrier does not beat the optimizer lower bound")
	}
	largeUpper, _, err := m24LargeRadiusUpper()
	if err != nil {
		return ComponentOptimumCertificate{}, err
	}
	if largeUpper.hi.Cmp(target) >= 0 {
		return ComponentOptimumCertificate{}, fmt.Errorf("M24 large-radius curvature barrier is too weak: %s", largeUpper.hi)
	}

	if len(m24RadiusBarriers) != 611 {
		return ComponentOptimumCertificate{}, fmt.Errorf("M24 generated radius barrier count changed: %d", len(m24RadiusBarriers))
	}
	bestUpper := riPoint(big.NewRat(0, 1))
	for index, cell := range m24RadiusBarriers {
		wantK := index
		if index >= 503 {
			wantK += 9
		}
		if cell.K != wantK {
			return ComponentOptimumCertificate{}, fmt.Errorf("M24 generated radius coverage gap at row %d: got %d want %d", index, cell.K, wantK)
		}
		upper, err := m24PointBarrierUpper(cell)
		if err != nil {
			return ComponentOptimumCertificate{}, err
		}
		if upper.hi.Cmp(target) >= 0 {
			return ComponentOptimumCertificate{}, fmt.Errorf("M24 cell %d upper %s does not beat target %s", cell.K, m23Decimal(upper.hi), m23Decimal(target))
		}
		if upper.hi.Cmp(bestUpper.hi) > 0 {
			bestUpper = upper
		}
	}

	component := LegalTwoAtomRadiusComponent{
		Definition:             "the connected component containing the M23 radius of {r>1: m_*(r)=inf_{H_r(t)<0} F0(t)/(-H_r(t)) is finite and nonnegative, w1(m_*(r),r)>=0, and c(m_*(r),r)>0}",
		SupportLegality:        "structural radii are 1 and r with r>1; no third atom, density, or widened support",
		MassLegality:           "m_*(r)>=0 is part of component membership; M23 proves it on the optimizer box",
		WeightLegality:         "w2=m_* and w1=w0-m_*(2*cos(pi*r/2)/pi-r*sin(pi*r/2)) are nonnegative by membership; M23 proves strict positivity on its box",
		ObjectiveLegality:      "c=c0+(pi/2)*m_* cos(pi*r/2)>0 by membership; the prompt's later minus-sign display is not the certified M22/M23 family",
		FamilyWellDefinedness:  "F0 and H_r are entire after the removable sinc value at t=0; the fixed double contact at pi/2 is structural for every r>1",
		ResponseSignStructure:  "H_r(pi/2)=partial_t H_r(pi/2)=0; away from that removable ratio point, negative H_r frequencies cap feasible mass",
		CertifiedInterior:      optimizer.Radius,
		BoundaryClassification: "the exact component endpoints and possible legality-induced disconnections remain unresolved; no arbitrary numerical truncation is asserted",
		ExactEndpointStatus:    "unresolved, but irrelevant to the proved objective exclusion outside [3.503,3.512]",
	}
	branches := []ContactBranch{
		{
			ID:               "B1-movable-low-frequency",
			RadiusRange:      "research continuation on [3,3.62]; certified only in the M23 root box",
			ContactRange:     "about [1.685,2.424] over the research range and [1.795310,1.795317] in the certified root box",
			EnvelopeBehavior: "sampled envelope rises to the M23 root and falls afterward",
			ProofStatus:      "M23 local Krawczyk/IFT theorem plus M24 pointwise upper barriers outside [3.503,3.512]; global branch uniqueness is unresolved inside that window",
			NextCompetitor:   "near M23 the next sampled local minimum has a mass gap greater than 0.01; research evidence only",
		},
		{
			ID:               "Bremote-oscillatory",
			RadiusRange:      "detected across representative radii",
			ContactRange:     "multiple higher-frequency pockets, including the M23 competitors near 9.44 and 15.54",
			EnvelopeBehavior: "not sampled active near M23",
			ProofStatus:      "eliminated outside [3.503,3.512] without root counting by the stronger distinguished-frequency caps; root count inside the unresolved window is not certified",
			NextCompetitor:   "finite candidate list remains numerical inside the unresolved window",
		},
	}
	evidence := GlobalEnvelopeEvidence{
		NegativeCosineRegion: "for 1<r<=3, cos(pi*r/2)<=0, so m>=0 implies c<=c0<the M23 lower bound",
		FixedContactFormula:  "F_tt(pi/2)=F0_tt(pi/2)-2*m*cos(pi*r/2)*(r^2-1+8/pi^2)>=0",
		LargeRadiusCutoff:    "r>=181/50; positive-cosine bands obey the fixed-contact curvature cap and nonpositive-cosine bands obey c<=c0",
		LargeRadiusUpper:     m23Display(largeUpper),
		GeneratedCellCount:   len(m24RadiusBarriers),
		GeneratedCellHash:    m24BarrierLogicalHash,
		CompactRadiusCover:   "exact cells [3+k/1000,3+(k+1)/1000], k=0..502 and 512..619, each eliminated by one Oct-generated fixed contact-frequency witness",
		BestCellUpper:        m23Display(bestUpper),
		UnresolvedRadius:     CertifiedScalarInterval{Lower: "3.5030000000000000", Upper: "3.5120000000000000"},
		UnresolvedContact:    CertifiedScalarInterval{Lower: "1.7905000000000000", Upper: "1.8020000000000000"},
		EvidenceBoundary:     "the interval inequalities and fixed-contact curvature formula are theorem evidence; Oct scans, branch labels, contact ranges, competitor gaps, and plots are research evidence only",
	}
	return ComponentOptimumCertificate{
		Status:                    "Success C — global control is proved outside one exact radius window; branch/global feasibility inside it remains unresolved",
		Component:                 component,
		Optimizer:                 optimizer,
		Branches:                  branches,
		GlobalEvidence:            evidence,
		M23ControlsComponent:      false,
		OptimizerGloballyFeasible: false,
		MissingTheorem:            "on r in [3.503,3.512], prove that every negative-response frequency gives m<=m_M23(r), equivalently certify the B1 contact continuation and exclude all competing tangencies; then combine with M23 C' signs",
	}, nil
}
