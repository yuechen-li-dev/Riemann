package semantic

import (
	"fmt"
	"math"
	"math/big"
)

const ContactAdaptedAtomicRayTransform = "F_m(t)=F_0(t)+m*(-pi*sin(t)/t+(4/pi)*cos(t)+2*cos(4*t))"

// ContactAdaptedAtomicRay is deliberately the single M20-selected ray.  It is
// not a generic two-atom or parameterized-measure representation.
type ContactAdaptedAtomicRay struct {
	Base                OneRadiusTangencyCandidate `json:"base_m18_witness"`
	Radius2             ExactRational              `json:"r2"`
	Mass                ExactRational              `json:"m"`
	C                   PiExpression               `json:"c"`
	Weight1             PiExpression               `json:"w1"`
	Weight2             ExactRational              `json:"w2"`
	Transform           string                     `json:"transform"`
	ContactPreservation string                     `json:"contact_preservation"`
	ContactFrequency    PiExpression               `json:"contact_frequency"`
}

func DeriveContactAdaptedAtomicRay(m ExactRational) (ContactAdaptedAtomicRay, error) {
	if err := m.Validate(); err != nil || rat(m).Sign() < 0 {
		return ContactAdaptedAtomicRay{}, fmt.Errorf("ray mass must be a nonnegative rational: %v", err)
	}
	base := DeriveOneRadiusTangencyCandidate()
	c := piAdd(base.ExactCExpression, piMonomial(1, new(big.Rat).Mul(rat(m), big.NewRat(1, 2))))
	c.Expression = "1+pi/4-2/pi+(pi/2)*m"
	w1 := piAdd(base.ExactWExpression, piMonomial(-1, new(big.Rat).Mul(rat(m), big.NewRat(2, 1))))
	w1.Expression = "1/pi-4/pi^2+8/pi^3+(2/pi)*m"
	return ContactAdaptedAtomicRay{
		Base: base, Radius2: ExactRational{4, 1}, Mass: m, C: c, Weight1: w1, Weight2: m,
		Transform:           ContactAdaptedAtomicRayTransform,
		ContactPreservation: "H(pi/2)=H'(pi/2)=0 exactly",
		ContactFrequency:    base.ContactFrequency,
	}, nil
}

func (r ContactAdaptedAtomicRay) Validate() error {
	if err := r.Base.Validate(); err != nil || r.Radius2 != (ExactRational{4, 1}) || r.Weight2 != r.Mass || r.Transform != ContactAdaptedAtomicRayTransform {
		return fmt.Errorf("wrong contact-adapted atomic ray: %v", err)
	}
	want, err := DeriveContactAdaptedAtomicRay(r.Mass)
	if err != nil || !samePiValue(r.C, want.C) || !samePiValue(r.Weight1, want.Weight1) || !samePiValue(r.ContactFrequency, want.ContactFrequency) {
		return fmt.Errorf("ray parameters are not derived from M18: %v", err)
	}
	if evalPi(r.C).hi.Cmp(big.NewRat(1651, 1250)) >= 0 {
		return fmt.Errorf("ray witness reaches compatible primal guard")
	}
	return nil
}

type ContactAdaptedRayAlgebra struct {
	Decomposition          string        `json:"decomposition"`
	Response               string        `json:"H"`
	ResponseAtContact      string        `json:"H_at_pi_over_2"`
	ResponsePrimeAtContact string        `json:"H_prime_at_pi_over_2"`
	ResponseSecond         PiExpression  `json:"H_second_at_pi_over_2"`
	RaySecond              PiExpression  `json:"F_m_second_at_pi_over_2"`
	SecondDerivativeLower  ExactRational `json:"second_derivative_lower"`
	Verified               bool          `json:"verified"`
}

func VerifyContactAdaptedRayAlgebra(r ContactAdaptedAtomicRay) (ContactAdaptedRayAlgebra, error) {
	if err := r.Validate(); err != nil {
		return ContactAdaptedRayAlgebra{}, err
	}
	// H(a)=-2+2=0 and H'(a)=4/pi-4/pi=0 at a=pi/2.
	hSecond := piAdd(piConstant(big.NewRat(-30, 1)), piMonomial(-2, big.NewRat(-16, 1)))
	hSecond.Expression = "-30-16/pi^2"
	baseSecond := deriveContactSecond(r.Base)
	raySecond := piAdd(baseSecond, piScale(hSecond, rat(r.Mass)))
	raySecond.Expression = "F0''(pi/2)-m*(30+16/pi^2)"
	if evalPi(raySecond).lo.Cmp(big.NewRat(1, 200)) <= 0 {
		return ContactAdaptedRayAlgebra{}, fmt.Errorf("original contact lacks the declared strict curvature margin")
	}
	return ContactAdaptedRayAlgebra{
		Decomposition:     ContactAdaptedAtomicRayTransform,
		Response:          "H(t)=-pi*sin(t)/t+(4/pi)*cos(t)+2*cos(4*t)",
		ResponseAtContact: "-2+2=0", ResponsePrimeAtContact: "4/pi-4/pi=0",
		ResponseSecond: hSecond, RaySecond: raySecond, SecondDerivativeLower: ExactRational{1, 200}, Verified: true,
	}, nil
}

func RayResponse(t float64) float64 {
	if math.Abs(t) < 1e-10 {
		return -math.Pi + 4/math.Pi + 2
	}
	return -math.Pi*math.Sin(t)/t + 4*math.Cos(t)/math.Pi + 2*math.Cos(4*t)
}

func ContactAdaptedRaySpectrum(r ContactAdaptedAtomicRay, t float64) float64 {
	m := float64(r.Mass.Numerator) / float64(r.Mass.Denominator)
	return OneRadiusSpectrum(1+math.Pi/4-2/math.Pi, 1/math.Pi-4/(math.Pi*math.Pi)+8/(math.Pi*math.Pi*math.Pi), t) + m*RayResponse(t)
}

func DeriveContactAdaptedRayJ(r ContactAdaptedAtomicRay) PiExpression {
	j := piAdd(piConstant(big.NewRat(2, 1)), piScale(r.C, big.NewRat(-1, 1)))
	j.Expression = "1-pi/4+2/pi-(pi/2)*m"
	return j
}

type RayCompactPiece struct {
	Left         ExactRational   `json:"left"`
	Right        ExactRational   `json:"right"`
	GridStep     ExactRational   `json:"grid_step"`
	TaylorDegree int             `json:"taylor_degree"`
	Kind         CompactCellKind `json:"kind"`
}

type ContactAdaptedRayWholeLineCertificate struct {
	ID                    string            `json:"id"`
	ContactPiece          RayCompactPiece   `json:"contact_piece"`
	ContactGridStep       ExactRational     `json:"contact_curvature_grid_step"`
	ContactThirdAbsUpper  ExactRational     `json:"contact_third_abs_upper"`
	ContactConvexityLower ExactRational     `json:"contact_convexity_lower"`
	StrictPieces          []RayCompactPiece `json:"strict_compact_pieces"`
	LipschitzUpper        ExactRational     `json:"lipschitz_upper"`
	TailStart             ExactRational     `json:"tail_start"`
	TailLower             ExactRational     `json:"tail_lower"`
	Coverage              string            `json:"coverage"`
	WholeLine             bool              `json:"whole_line"`
}

func rayPointInterval(r ContactAdaptedAtomicRay, t *big.Rat, degree int) rationalInterval {
	if t.Sign() == 0 {
		one := piConstant(big.NewRat(1, 1))
		origin := piScale(piAdd(piAdd(one, r.Weight1), piScale(r.C, big.NewRat(-1, 1))), big.NewRat(2, 1))
		return riAdd(evalPi(origin), riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(2, 1))))
	}
	cosLo, cosHi := trigTaylorBounds(t, degree, true)
	sinLo, sinHi := trigTaylorBounds(t, degree-1, false)
	fourT := new(big.Rat).Mul(t, big.NewRat(4, 1))
	cos4Lo, cos4Hi := trigTaylorBounds(fourT, degree, true)
	cosI, sinI, cos4I := ri(cosLo, cosHi), ri(sinLo, sinHi), ri(cos4Lo, cos4Hi)
	tI, t2I := riPoint(t), riPoint(new(big.Rat).Mul(t, t))
	value := riPoint(big.NewRat(1, 1))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riSub(cosI, riPoint(big.NewRat(1, 1))), riReciprocal(t2I))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riMul(riSub(riPoint(big.NewRat(1, 1)), evalPi(r.C)), sinI), riReciprocal(tI))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(evalPi(r.Weight1), cosI)))
	value = riAdd(value, riMul(riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(2, 1))), cos4I))
	return value
}

func raySecondPointInterval(r ContactAdaptedAtomicRay, t *big.Rat, degree int) rationalInterval {
	cosLo, cosHi := trigTaylorBounds(t, degree, true)
	sinLo, sinHi := trigTaylorBounds(t, degree-1, false)
	fourT := new(big.Rat).Mul(t, big.NewRat(4, 1))
	cos4Lo, cos4Hi := trigTaylorBounds(fourT, degree, true)
	cosI, sinI, cos4I := ri(cosLo, cosHi), ri(sinLo, sinHi), ri(cos4Lo, cos4Hi)
	tI := riPoint(t)
	invT := riReciprocal(tI)
	invT2, invT3, invT4 := riPow(tI, -2), riPow(tI, -3), riPow(tI, -4)
	value := riMul(riPoint(big.NewRat(-2, 1)), riMul(cosI, invT2))
	value = riAdd(value, riMul(riPoint(big.NewRat(8, 1)), riMul(sinI, invT3)))
	value = riAdd(value, riMul(riPoint(big.NewRat(12, 1)), riMul(riSub(cosI, riPoint(big.NewRat(1, 1))), invT4)))
	inner := riAdd(riMul(riPoint(big.NewRat(-1, 1)), riMul(sinI, invT)), riMul(riPoint(big.NewRat(-2, 1)), riMul(cosI, invT2)))
	inner = riAdd(inner, riMul(riPoint(big.NewRat(2, 1)), riMul(sinI, invT3)))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riSub(riPoint(big.NewRat(1, 1)), evalPi(r.C)), inner)))
	value = riAdd(value, riMul(riPoint(big.NewRat(-2, 1)), riMul(evalPi(r.Weight1), cosI)))
	value = riAdd(value, riMul(riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(-32, 1))), cos4I))
	return value
}

func verifyRayPiece(r ContactAdaptedAtomicRay, p RayCompactPiece, lipschitz ExactRational) error {
	left, right, step := rat(p.Left), rat(p.Right), rat(p.GridStep)
	if left.Sign() < 0 || right.Cmp(left) <= 0 || step.Sign() <= 0 || p.Kind != StrictPositiveCell || p.TaylorDegree < 80 {
		return fmt.Errorf("invalid strict ray piece")
	}
	span := new(big.Rat).Sub(right, left)
	q := new(big.Rat).Quo(span, step)
	if !q.IsInt() {
		return fmt.Errorf("ray piece is not an exact grid multiple")
	}
	count := q.Num().Int64()
	loss := new(big.Rat).Mul(rat(lipschitz), new(big.Rat).Quo(step, big.NewRat(2, 1)))
	for k := int64(0); k <= count; k++ {
		t := new(big.Rat).Add(left, new(big.Rat).Mul(big.NewRat(k, 1), step))
		if rayPointInterval(r, t, p.TaylorDegree).lo.Cmp(loss) <= 0 {
			return fmt.Errorf("strict ray piece [%s,%s] fails at %s", left, right, t)
		}
	}
	return nil
}

func VerifyContactAdaptedRayWholeLine(r ContactAdaptedAtomicRay, p ContactAdaptedRayWholeLineCertificate, replayCompact bool) error {
	if _, err := VerifyContactAdaptedRayAlgebra(r); err != nil {
		return err
	}
	if p.ID == "" || !p.WholeLine || p.ContactPiece.Kind != NonnegativeCellWithCertifiedContact || p.ContactPiece.Left != (ExactRational{151, 100}) || p.ContactPiece.Right != (ExactRational{17, 10}) || p.ContactGridStep != (ExactRational{1, 10000}) || p.ContactThirdAbsUpper != (ExactRational{11, 10}) || p.TailStart != (ExactRational{4, 1}) {
		return fmt.Errorf("incomplete M21 equality-aware certificate")
	}
	if len(p.StrictPieces) != 3 || p.StrictPieces[0].Left != (ExactRational{0, 1}) || p.StrictPieces[0].Right != (ExactRational{13, 10}) || p.StrictPieces[1].Left != (ExactRational{13, 10}) || p.StrictPieces[1].Right != (ExactRational{151, 100}) || p.StrictPieces[2].Left != (ExactRational{17, 10}) || p.StrictPieces[2].Right != (ExactRational{4, 1}) {
		return fmt.Errorf("M21 compact partition has a gap")
	}
	// |F'| <= c-2/3+2*w1+8*m from the integral form and both atoms.
	lipschitz := riAdd(riSub(evalPi(r.C), riPoint(big.NewRat(2, 3))), riMul(evalPi(r.Weight1), riPoint(big.NewRat(2, 1))))
	lipschitz = riAdd(lipschitz, riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(8, 1))))
	if lipschitz.hi.Cmp(rat(p.LipschitzUpper)) >= 0 {
		return fmt.Errorf("M21 Lipschitz declaration is too small")
	}
	third := riSub(riMul(evalPi(r.C), riPoint(big.NewRat(1, 2))), riPoint(big.NewRat(2, 5)))
	third = riAdd(third, riMul(evalPi(r.Weight1), riPoint(big.NewRat(2, 1))))
	third = riAdd(third, riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(128, 1))))
	if third.hi.Cmp(rat(p.ContactThirdAbsUpper)) >= 0 {
		return fmt.Errorf("M21 contact third-derivative declaration is too small")
	}
	if new(big.Rat).Quo(m18Pi.lo, big.NewRat(2, 1)).Cmp(rat(p.ContactPiece.Left)) <= 0 || new(big.Rat).Quo(m18Pi.hi, big.NewRat(2, 1)).Cmp(rat(p.ContactPiece.Right)) >= 0 {
		return fmt.Errorf("M21 contact cell does not contain pi/2")
	}
	// Grid-certified F''>0 on the equality cell.  The 11/10 bound controls
	// the loss between curvature samples; exact F=F'=0 then gives F>=0.
	left, right, step := rat(p.ContactPiece.Left), rat(p.ContactPiece.Right), rat(p.ContactGridStep)
	minSecond := (*big.Rat)(nil)
	count := new(big.Rat).Quo(new(big.Rat).Sub(right, left), step).Num().Int64()
	for k := int64(0); k <= count; k++ {
		t := new(big.Rat).Add(left, new(big.Rat).Mul(big.NewRat(k, 1), step))
		lo := raySecondPointInterval(r, t, p.ContactPiece.TaylorDegree).lo
		if minSecond == nil || lo.Cmp(minSecond) < 0 {
			minSecond = lo
		}
	}
	minSecond.Sub(minSecond, new(big.Rat).Mul(rat(p.ContactThirdAbsUpper), new(big.Rat).Quo(step, big.NewRat(2, 1))))
	if minSecond.Cmp(rat(p.ContactConvexityLower)) <= 0 {
		return fmt.Errorf("M21 contact convexity lower bound failed")
	}
	if replayCompact {
		for _, piece := range p.StrictPieces {
			if err := verifyRayPiece(r, piece, p.LipschitzUpper); err != nil {
				return err
			}
		}
	}
	// |t|>=4: F>=3/4-(c-1)/2-2*w1-2*m.
	tail := riSub(riSub(riPoint(big.NewRat(3, 4)), riMul(riSub(evalPi(r.C), riPoint(big.NewRat(1, 1))), riPoint(big.NewRat(1, 2)))), riMul(evalPi(r.Weight1), riPoint(big.NewRat(2, 1))))
	tail = riSub(tail, riPoint(new(big.Rat).Mul(rat(r.Mass), big.NewRat(2, 1))))
	if tail.lo.Cmp(rat(p.TailLower)) <= 0 || p.Coverage != "evenness; strict Taylor/Lipschitz cells on [0,13/10], [13/10,151/100], and [17/10,4]; certified convex contact cell [151/100,17/10]; analytic atomic tail [4,infinity)" {
		return fmt.Errorf("M21 tail or coverage certificate failed")
	}
	return nil
}

type CertifiedRayCeilingBracket struct {
	LowerWitness          ExactRational `json:"lower_witness"`
	UpperBound            ExactRational `json:"upper_bound"`
	UpperWitnessFrequency ExactRational `json:"upper_witness_frequency"`
	Derivation            string        `json:"derivation"`
}

func VerifyRayCeilingUpper(bracket CertifiedRayCeilingBracket, degree int) error {
	if rat(bracket.LowerWitness).Cmp(rat(bracket.UpperBound)) >= 0 {
		return fmt.Errorf("empty ray ceiling bracket")
	}
	r, err := DeriveContactAdaptedAtomicRay(bracket.UpperBound)
	if err != nil {
		return err
	}
	if rayPointInterval(r, rat(bracket.UpperWitnessFrequency), degree).hi.Sign() >= 0 {
		return fmt.Errorf("upper ray mass is not rejected at its certified frequency")
	}
	return nil
}
