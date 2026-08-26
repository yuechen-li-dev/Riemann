package semantic

import (
	"fmt"
	"math"
	"strings"
)

const (
	M20BaselineC       = 1 + math.Pi/4 - 2/math.Pi
	M20BaselineW       = 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi)
	M20Contact         = math.Pi / 2
	CompatiblePrimalC  = 1651.0 / 1250.0
	m20FiniteTolerance = 1e-12
)

// EvidenceClass prevents exploratory M20 scans from silently entering the
// theorem-bearing evidence vocabulary used by earlier milestones.
type EvidenceClass string

const (
	StructuralEvidence  EvidenceClass = "Structural"
	NumericalEvidence   EvidenceClass = "Numerical"
	CertifiedEvidence   EvidenceClass = "Certified"
	SpeculativeEvidence EvidenceClass = "Speculative"
)

type M20Probe struct {
	Family         string        `json:"family"`
	CandidateC     float64       `json:"candidate_c"`
	SampledMinimum float64       `json:"sampled_minimum"`
	MinimumAt      float64       `json:"minimum_at"`
	Evidence       EvidenceClass `json:"evidence"`
	TheoremClaim   bool          `json:"theorem_claim"`
}

func (p M20Probe) Validate() error {
	if strings.TrimSpace(p.Family) == "" || math.IsNaN(p.CandidateC) || math.IsInf(p.CandidateC, 0) || math.IsNaN(p.SampledMinimum) || math.IsInf(p.SampledMinimum, 0) {
		return fmt.Errorf("incomplete M20 candidate probe")
	}
	switch p.Evidence {
	case StructuralEvidence, NumericalEvidence, CertifiedEvidence, SpeculativeEvidence:
	default:
		return fmt.Errorf("unknown evidence classification %q", p.Evidence)
	}
	if p.CandidateC >= CompatiblePrimalC {
		return fmt.Errorf("candidate c=%g reaches or exceeds compatible primal guard 1651/1250", p.CandidateC)
	}
	if p.TheoremClaim && p.Evidence != CertifiedEvidence {
		return fmt.Errorf("%s evidence cannot carry a theorem claim", p.Evidence)
	}
	return nil
}

type ContactVariation struct {
	Mode                 string  `json:"mode"`
	Mass                 float64 `json:"mass"`
	ResponseAtContact    float64 `json:"H_at_contact"`
	ResponseDerivative   float64 `json:"H_prime_at_contact"`
	CIncrease            float64 `json:"gamma"`
	BoundaryWeightChange float64 `json:"beta"`
	PreservesValue       bool    `json:"preserves_contact_value"`
	PreservesDerivative  bool    `json:"preserves_contact_derivative"`
}

func OneRadiusSpectrum(c, w, t float64) float64 {
	if math.Abs(t) < 1e-8 {
		return 2 * (1 + w - c)
	}
	return 1 + 2*(math.Cos(t)-1)/(t*t) + 2*(1-c)*math.Sin(t)/t + 2*w*math.Cos(t)
}

func M20BaselineSpectrum(t float64) float64 { return OneRadiusSpectrum(M20BaselineC, M20BaselineW, t) }

func m20BaseDerivative(t float64) float64 {
	c, w := M20BaselineC, M20BaselineW
	return -2*math.Sin(t)/(t*t) - 4*(math.Cos(t)-1)/(t*t*t) +
		2*(1-c)*(math.Cos(t)/t-math.Sin(t)/(t*t)) - 2*w*math.Sin(t)
}

func M20BaselineCurvature() float64 {
	t, c, w := M20Contact, M20BaselineC, M20BaselineW
	return -2*math.Cos(t)/(t*t) + 8*math.Sin(t)/(t*t*t) + 12*(math.Cos(t)-1)/(t*t*t*t) +
		2*(1-c)*(-math.Sin(t)/t-2*math.Cos(t)/(t*t)+2*math.Sin(t)/(t*t*t)) - 2*w*math.Cos(t)
}

// ContactVariationFromResponse solves the two linearized active-contact
// equations.  If H is the new exterior mode and c->c+gamma, w->w+beta,
// then delta F(a)=-2 gamma/a+H(a) and
// delta F'(a)=2 gamma/a^2-2 beta+H'(a), a=pi/2.
func ContactVariationFromResponse(mode string, mass, h, hPrime float64) (ContactVariation, error) {
	if mass < 0 || math.IsNaN(mass) || math.IsInf(mass, 0) {
		return ContactVariation{}, fmt.Errorf("perturbation mass must be finite and nonnegative")
	}
	a := M20Contact
	gamma := a * h / 2
	beta := gamma/(a*a) + hPrime/2
	v := ContactVariation{Mode: mode, Mass: mass, ResponseAtContact: h, ResponseDerivative: hPrime,
		CIncrease: gamma, BoundaryWeightChange: beta}
	v.PreservesValue = math.Abs(-2*gamma/a+h) < m20FiniteTolerance
	v.PreservesDerivative = math.Abs(2*gamma/(a*a)-2*beta+hPrime) < m20FiniteTolerance
	return v, nil
}

func ContactAdaptedAtom(mass, radius float64) (ContactVariation, error) {
	if radius < 1 || math.IsNaN(radius) || math.IsInf(radius, 0) {
		return ContactVariation{}, fmt.Errorf("atom radius must be finite and exterior")
	}
	a := M20Contact
	h := 2 * mass * math.Cos(radius*a)
	hPrime := -2 * mass * radius * math.Sin(radius*a)
	return ContactVariationFromResponse(fmt.Sprintf("atom(r=%g)", radius), mass, h, hPrime)
}

func sinc(x float64) float64 {
	if math.Abs(x) < 1e-10 {
		return 1
	}
	return math.Sin(x) / x
}

func sincDerivativeAt(scale, t float64) float64 {
	x := scale * t
	if math.Abs(x) < 1e-8 {
		return 0
	}
	return scale * (x*math.Cos(x) - math.Sin(x)) / (x * x)
}

// ContactAdaptedInterval is a uniform positive-side mass on [center-halfWidth,
// center+halfWidth], mirrored to the negative axis.
func ContactAdaptedInterval(mass, center, halfWidth float64) (ContactVariation, error) {
	if mass < 0 || halfWidth <= 0 || center-halfWidth < 1 {
		return ContactVariation{}, fmt.Errorf("interval density must be nonnegative and supported outside (-1,1)")
	}
	a := M20Contact
	s := sinc(halfWidth * a)
	h := 2 * mass * s * math.Cos(center*a)
	hPrime := 2 * mass * (sincDerivativeAt(halfWidth, a)*math.Cos(center*a) - s*center*math.Sin(center*a))
	return ContactVariationFromResponse(fmt.Sprintf("uniform[%g,%g]", center-halfWidth, center+halfWidth), mass, h, hPrime)
}

func AtomExtendedSpectrum(v ContactVariation, radius, t float64) float64 {
	return OneRadiusSpectrum(M20BaselineC+v.CIncrease, M20BaselineW+v.BoundaryWeightChange, t) + 2*v.Mass*math.Cos(radius*t)
}

func IntervalExtendedSpectrum(v ContactVariation, center, halfWidth, t float64) float64 {
	return OneRadiusSpectrum(M20BaselineC+v.CIncrease, M20BaselineW+v.BoundaryWeightChange, t) + 2*v.Mass*sinc(halfWidth*t)*math.Cos(center*t)
}

func ContactAdaptedLatticePacket(mass float64) (ContactVariation, error) {
	if mass < 0 || math.IsNaN(mass) || math.IsInf(mass, 0) {
		return ContactVariation{}, fmt.Errorf("lattice-packet mass must be finite and nonnegative")
	}
	a := M20Contact
	h := 2 * mass / 3 * (math.Cos(3*a) + math.Cos(4*a) + math.Cos(5*a))
	hPrime := -2 * mass / 3 * (3*math.Sin(3*a) + 4*math.Sin(4*a) + 5*math.Sin(5*a))
	return ContactVariationFromResponse("equal packet r={3,4,5}", mass, h, hPrime)
}

func LatticePacketExtendedSpectrum(v ContactVariation, t float64) float64 {
	return OneRadiusSpectrum(M20BaselineC+v.CIncrease, M20BaselineW+v.BoundaryWeightChange, t) +
		2*v.Mass/3*(math.Cos(3*t)+math.Cos(4*t)+math.Cos(5*t))
}

func PredictContactShift(deltaDerivative float64) float64 {
	return -deltaDerivative / M20BaselineCurvature()
}

func BaselineContactDerivative() float64 { return m20BaseDerivative(M20Contact) }

type RepresentationUtility struct {
	Family         string `json:"family"`
	Improvement    int    `json:"expected_improvement"`
	Tractability   int    `json:"analytic_tractability"`
	WholeLineCost  int    `json:"whole_line_proof_cost_inverse"`
	ExactPotential int    `json:"exact_witness_potential"`
	Reuse          int    `json:"full_class_reuse"`
	Implementation int    `json:"implementation_cost_inverse"`
	Clarity        int    `json:"conceptual_clarity"`
}

func (u RepresentationUtility) Score() int {
	return 3*u.Improvement + 2*u.Tractability + 2*u.WholeLineCost + u.ExactPotential + u.Reuse + u.Implementation + u.Clarity
}
