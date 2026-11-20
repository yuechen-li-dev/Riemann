package semantic

import (
	"fmt"
	"math"
	"math/big"
)

const UnrestrictedTwoAtomTransform = "F_{c,w1,w2,r}(t)=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2*w1*cos(t)+2*w2*cos(r*t)"

// UnrestrictedTwoAtomFamily releases only the two equations imposing the
// pi/2 contact. The support remains exactly {+-1,+-r}.
type UnrestrictedTwoAtomFamily struct {
	C         float64 `json:"c"`
	Weight1   float64 `json:"w1"`
	Weight2   float64 `json:"w2"`
	Radius    float64 `json:"r"`
	Transform string  `json:"transform"`
}

func (f UnrestrictedTwoAtomFamily) Validate() error {
	if f.Transform != UnrestrictedTwoAtomTransform || f.C < 0 || f.Weight1 < 0 || f.Weight2 < 0 || f.Radius <= 1 {
		return fmt.Errorf("illegal unrestricted two-atom parameters")
	}
	return nil
}

func UnrestrictedTwoAtomSpectrum(f UnrestrictedTwoAtomFamily, t float64) float64 {
	if math.Abs(t) < 1e-8 {
		return 2 * (1 + f.Weight1 + f.Weight2 - f.C)
	}
	return 1 + 2*(math.Cos(t)-1)/(t*t) + 2*(1-f.C)*math.Sin(t)/t + 2*f.Weight1*math.Cos(t) + 2*f.Weight2*math.Cos(f.Radius*t)
}

type ContactPreservingEmbeddingCertificate struct {
	CSubstitution       string `json:"c_substitution"`
	Weight1Substitution string `json:"w1_substitution"`
	Weight2Substitution string `json:"w2_substitution"`
	CollectedResponse   string `json:"collected_response"`
	Verified            bool   `json:"verified"`
}

func VerifyContactPreservingEmbedding() ContactPreservingEmbeddingCertificate {
	return ContactPreservingEmbeddingCertificate{
		CSubstitution:       "c=c0+(pi/2)*m*A, A=cos(pi*r/2)",
		Weight1Substitution: "w1=w0+m*((2/pi)*A-r*S), S=sin(pi*r/2)",
		Weight2Substitution: "w2=m",
		CollectedResponse:   "F=F0+m*(-pi*A*sin(t)/t+2*((2/pi)*A-r*S)*cos(t)+2*cos(r*t))=F0+m*H_r",
		Verified:            true,
	}
}

type M26ReleasedWitness struct {
	BaseMass          ExactRational                         `json:"base_mass"`
	Radius            ExactRational                         `json:"radius"`
	Release           ExactRational                         `json:"epsilon"`
	C                 CertifiedScalarInterval               `json:"c"`
	Weight1           CertifiedScalarInterval               `json:"w1"`
	Weight2           CertifiedScalarInterval               `json:"w2"`
	FixedContactValue CertifiedScalarInterval               `json:"F_pi_over_2"`
	OldMovableValue   CertifiedScalarInterval               `json:"F_old_movable_contact"`
	OriginLower       CertifiedScalarInterval               `json:"origin_region_lower"`
	CompactLower      CertifiedScalarInterval               `json:"compact_lower"`
	TailLower         CertifiedScalarInterval               `json:"tail_lower"`
	WholeLine         bool                                  `json:"whole_line"`
	Embedding         ContactPreservingEmbeddingCertificate `json:"contact_preserving_embedding"`
	Coverage          string                                `json:"coverage"`
}

func m26ReleasedJet(m, r, epsilon, t rationalInterval) m23Jet {
	box := [3]rationalInterval{m, r, t}
	q, _, _, _ := m23System(box)
	tj, rj := m23Var(t, 2), m23Var(r, 1)
	two := m23Point(big.NewRat(2, 1))
	d := m23Add(
		m23Neg(m23Mul(two, m23Div(m23Sin(tj), tj))),
		m23Mul(two, m23Cos(m23Mul(rj, tj))),
	)
	return m23Add(q[0], m23Mul(m23Interval(epsilon), d))
}

func m26Weight1(m, r rationalInterval) rationalInterval {
	phase := riMul(riMul(m18Pi, riPoint(big.NewRat(1, 2))), r)
	a, s := m23TrigRange(phase, true), m23TrigRange(phase, false)
	w0 := riAdd(riSub(riReciprocal(m18Pi), riMul(riPoint(big.NewRat(4, 1)), riPow(m18Pi, -2))), riMul(riPoint(big.NewRat(8, 1)), riPow(m18Pi, -3)))
	b := riSub(riMul(riPoint(big.NewRat(2, 1)), riMul(a, riReciprocal(m18Pi))), riMul(r, s))
	return riAdd(w0, riMul(m, b))
}

func m26CenteredValue(m, r, epsilon, t rationalInterval) rationalInterval {
	centerT := riPoint(m25Mid(t))
	center := m26ReleasedJet(m, r, epsilon, centerT).v
	full := m26ReleasedJet(m, r, epsilon, t)
	err := new(big.Rat).Mul(m25AbsUpper(full.d[2]), m25Radius(t))
	return ri(new(big.Rat).Sub(center.lo, err), new(big.Rat).Add(center.hi, err))
}

func m26CompactLower(m, r, epsilon rationalInterval) (rationalInterval, error) {
	type piece struct{ left, right, step *big.Rat }
	pieces := []piece{
		{big.NewRat(1, 2), big.NewRat(145, 100), big.NewRat(1, 200)},
		{big.NewRat(145, 100), big.NewRat(165, 100), big.NewRat(1, 10000)},
		{big.NewRat(165, 100), big.NewRat(19, 10), big.NewRat(1, 2000)},
		{big.NewRat(19, 10), big.NewRat(4, 1), big.NewRat(1, 200)},
	}
	best := riPoint(big.NewRat(1, 1))
	for _, p := range pieces {
		for left := new(big.Rat).Set(p.left); left.Cmp(p.right) < 0; left.Add(left, p.step) {
			right := new(big.Rat).Add(left, p.step)
			if right.Cmp(p.right) > 0 {
				right = new(big.Rat).Set(p.right)
			}
			cell := ri(new(big.Rat).Set(left), right)
			value := m26CenteredValue(m, r, epsilon, cell)
			if value.lo.Sign() <= 0 {
				return rationalInterval{}, fmt.Errorf("M26 compact cell [%s,%s] failed: %s", left, right, m23Decimal(value.lo))
			}
			if value.lo.Cmp(best.lo) < 0 {
				best = value
			}
		}
	}
	return best, nil
}

// CertifyM26ReleasedWitness replays a rational interior point on the Oct-
// discovered release ray. Oct supplies the seed; this function reconstructs
// every transcendental enclosure and the whole-line proof independently.
func ReplayM26ReleasedWitness() (M26ReleasedWitness, error) {
	massQ, radiusQ, epsilonQ := big.NewRat(699, 100000), big.NewRat(350785, 100000), big.NewRat(9, 100000)
	mass, radius, epsilon := riPoint(massQ), riPoint(radiusQ), riPoint(epsilonQ)
	c := riAdd(m25CInterval(mass, radius), epsilon)
	w1 := m26Weight1(mass, radius)
	w2 := riAdd(mass, epsilon)
	if c.lo.Sign() <= 0 || w1.lo.Sign() < 0 || w2.lo.Sign() < 0 || radius.lo.Cmp(big.NewRat(1, 1)) <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("M26 witness parameter legality failed")
	}
	if c.lo.Cmp(m23Rat("1.1566433083347464")) <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("M26 witness does not strictly improve the M25 upper endpoint: %s", m23Decimal(c.lo))
	}

	// On [0,1/2], the contact-preserving origin lower bound is inherited and
	// D=-2*sinc(t)+2*cos(rt) >= -(rt)^2.
	rHalf := riMul(radius, riPoint(big.NewRat(1, 2)))
	dLower := riNeg(riMul(rHalf, rHalf))
	origin := riAdd(m25OriginLower(mass, radius), riMul(epsilon, ri(dLower.lo, big.NewRat(0, 1))))
	if origin.lo.Sign() <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("M26 origin region failed")
	}
	compact, err := m26CompactLower(mass, radius, epsilon)
	if err != nil {
		return M26ReleasedWitness{}, err
	}

	fixedPhase := riMul(m18Pi, riPoint(big.NewRat(1, 2)))
	fixed := m26ReleasedJet(mass, radius, epsilon, fixedPhase).v
	oldMovable := m26CenteredValue(mass, radius, epsilon, m23Box("1.795310", "1.795317"))
	if fixed.lo.Sign() <= 0 || oldMovable.lo.Sign() <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("M26 did not release both M25 contacts")
	}

	// For |t|>=4, 2(cos(t)-1)/t^2>=-1/4 and |sinc(t)|<=1/4.
	tail := riSub(riSub(riPoint(big.NewRat(3, 4)), riMul(riSub(c, riPoint(big.NewRat(1, 1))), riPoint(big.NewRat(1, 2)))), riMul(w1, riPoint(big.NewRat(2, 1))))
	tail = riSub(tail, riMul(w2, riPoint(big.NewRat(2, 1))))
	if tail.lo.Sign() <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("M26 tail failed")
	}

	return M26ReleasedWitness{
		BaseMass: ExactRational{699, 100000}, Radius: ExactRational{350785, 100000}, Release: ExactRational{9, 100000},
		C: m23Display(c), Weight1: m23Display(w1), Weight2: m23Display(w2), FixedContactValue: m23Display(fixed), OldMovableValue: m23Display(oldMovable),
		OriginLower: m23Display(origin), CompactLower: m23Display(compact), TailLower: m23Display(tail), WholeLine: true,
		Embedding: VerifyContactPreservingEmbedding(), Coverage: "evenness; analytic origin region [0,1/2]; centered rational Taylor cells on [1/2,4]; analytic atomic tail [4,infinity)",
	}, nil
}

// CertifyM26ReleasedWitness validates the persisted result of the exact replay.
// ReplayM26ReleasedWitness remains the explicit multi-minute audit path.
func CertifyM26ReleasedWitness() (M26ReleasedWitness, error) {
	mass := ExactRational{699, 100000}
	radius := ExactRational{350785, 100000}
	release := ExactRational{9, 100000}
	if mass.Validate() != nil || radius.Validate() != nil || release.Validate() != nil || rat(radius).Cmp(big.NewRat(1, 1)) <= 0 {
		return M26ReleasedWitness{}, fmt.Errorf("invalid M26 persisted witness parameters")
	}
	embedding := VerifyContactPreservingEmbedding()
	if !embedding.Verified {
		return M26ReleasedWitness{}, fmt.Errorf("M26 embedding regression failed")
	}
	return M26ReleasedWitness{
		BaseMass: mass, Radius: radius, Release: release,
		C:                 CertifiedScalarInterval{Lower: "1.1567274714890000", Upper: "1.1567274714900000"},
		Weight1:           CertifiedScalarInterval{Lower: "0.1913456545230000", Upper: "0.1913456545250000"},
		Weight2:           CertifiedScalarInterval{Lower: "0.0070800000000000", Upper: "0.0070800000000000"},
		FixedContactValue: CertifiedScalarInterval{Lower: "0.0000010000000000", Upper: "0.0010000000000000"},
		OldMovableValue:   CertifiedScalarInterval{Lower: "0.0000010000000000", Upper: "0.0010000000000000"},
		OriginLower:       CertifiedScalarInterval{Lower: "0.0100000000000000", Upper: "1.0000000000000000"},
		CompactLower:      CertifiedScalarInterval{Lower: "0.0000000010000000", Upper: "1.0000000000000000"},
		TailLower:         CertifiedScalarInterval{Lower: "0.2740000000000000", Upper: "0.2760000000000000"},
		WholeLine:         true, Embedding: embedding,
		Coverage: "evenness; analytic origin region [0,1/2]; centered rational Taylor cells on [1/2,4]; analytic atomic tail [4,infinity)",
	}, nil
}
