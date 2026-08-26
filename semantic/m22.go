package semantic

import (
	"fmt"
	"math"
	"math/big"
)

const ContactPreservingTwoAtomTransform = "F_{m,r}(t)=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2*w1*cos(t)+2*m*cos(r*t)"

// ContactPreservingTwoAtomFamily is the two-dimensional surface obtained by
// solving the two pi/2 contact equations.  It is intentionally not a generic
// measure optimizer.
type ContactPreservingTwoAtomFamily struct {
	Mass                ExactRational `json:"m"`
	Radius              ExactRational `json:"r"`
	CFormula            string        `json:"c_formula"`
	Weight1Formula      string        `json:"w1_formula"`
	Weight2             ExactRational `json:"w2"`
	Transform           string        `json:"transform"`
	ContactPreservation string        `json:"contact_preservation"`
}

type RadiusCandidate struct {
	Radius           float64 `json:"r"`
	Mass             float64 `json:"m"`
	C                float64 `json:"c"`
	Weight1          float64 `json:"w1"`
	ContactFrequency float64 `json:"t1"`
	Curvature        float64 `json:"curvature"`
	Evidence         string  `json:"evidence"`
}

type ContactMotionEvidence struct {
	DeltaCPerRadiusAtFixedMass       float64 `json:"partial_r_c"`
	DeltaWeight1PerRadiusAtFixedMass float64 `json:"partial_r_w1"`
	DeltaFAtOldContact               float64 `json:"partial_r_F_at_t1"`
	DeltaFPrimeAtOldContact          float64 `json:"partial_r_F_prime_at_t1"`
	ShiftedContactDerivative         float64 `json:"dt1_dr"`
	Conclusion                       string  `json:"conclusion"`
}

func DeriveContactPreservingTwoAtomFamily(m, r ExactRational) (ContactPreservingTwoAtomFamily, error) {
	if err := m.Validate(); err != nil || rat(m).Sign() < 0 {
		return ContactPreservingTwoAtomFamily{}, fmt.Errorf("mass must be a nonnegative rational: %v", err)
	}
	if err := r.Validate(); err != nil || rat(r).Cmp(big.NewRat(1, 1)) <= 0 {
		return ContactPreservingTwoAtomFamily{}, fmt.Errorf("second radius must satisfy r>1: %v", err)
	}
	f := ContactPreservingTwoAtomFamily{
		Mass: m, Radius: r, Weight2: m, Transform: ContactPreservingTwoAtomTransform,
		CFormula:            "c0+(pi/2)*m*cos(pi*r/2)",
		Weight1Formula:      "w0+m*((2/pi)*cos(pi*r/2)-r*sin(pi*r/2))",
		ContactPreservation: "F(pi/2)=F'(pi/2)=0 after solving the fixed nonsingular 2-by-2 system for (c,w1)",
	}
	c, w := TwoAtomParametersFloat(f)
	if !(c > 0) || w < 0 || c >= CompatiblePrimalC {
		return ContactPreservingTwoAtomFamily{}, fmt.Errorf("derived family point is illegal: c=%g w1=%g", c, w)
	}
	return f, nil
}

func (f ContactPreservingTwoAtomFamily) Validate() error {
	want, err := DeriveContactPreservingTwoAtomFamily(f.Mass, f.Radius)
	if err != nil {
		return err
	}
	if f != want {
		return fmt.Errorf("two-atom family is not the derived contact-preserving surface")
	}
	return nil
}

func TwoAtomParametersFloat(f ContactPreservingTwoAtomFamily) (c, w1 float64) {
	m := float64(f.Mass.Numerator) / float64(f.Mass.Denominator)
	r := float64(f.Radius.Numerator) / float64(f.Radius.Denominator)
	phase := math.Pi * r / 2
	c0 := 1 + math.Pi/4 - 2/math.Pi
	w0 := 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi)
	return c0 + math.Pi*m*math.Cos(phase)/2,
		w0 + m*(2*math.Cos(phase)/math.Pi-r*math.Sin(phase))
}

func TwoAtomResponseFloat(r, t float64) float64 {
	phase := math.Pi * r / 2
	c, s := math.Cos(phase), math.Sin(phase)
	sinc := 1.0
	if math.Abs(t) > 1e-10 {
		sinc = math.Sin(t) / t
	}
	return -math.Pi*c*sinc + 2*(2*c/math.Pi-r*s)*math.Cos(t) + 2*math.Cos(r*t)
}

func ContactPreservingTwoAtomSpectrum(f ContactPreservingTwoAtomFamily, t float64) float64 {
	m := float64(f.Mass.Numerator) / float64(f.Mass.Denominator)
	r := float64(f.Radius.Numerator) / float64(f.Radius.Denominator)
	c0 := 1 + math.Pi/4 - 2/math.Pi
	w0 := 1/math.Pi - 4/(math.Pi*math.Pi) + 8/(math.Pi*math.Pi*math.Pi)
	return OneRadiusSpectrum(c0, w0, t) + m*TwoAtomResponseFloat(r, t)
}

func ContactPreservingTwoAtomSpectrumPrime(f ContactPreservingTwoAtomFamily, t float64) float64 {
	c, w := TwoAtomParametersFloat(f)
	m := float64(f.Mass.Numerator) / float64(f.Mass.Denominator)
	r := float64(f.Radius.Numerator) / float64(f.Radius.Denominator)
	return -2*math.Sin(t)/(t*t) - 4*(math.Cos(t)-1)/(t*t*t) +
		2*(1-c)*(math.Cos(t)/t-math.Sin(t)/(t*t)) - 2*w*math.Sin(t) - 2*m*r*math.Sin(r*t)
}

func ContactPreservingTwoAtomSpectrumSecond(f ContactPreservingTwoAtomFamily, t float64) float64 {
	c, w := TwoAtomParametersFloat(f)
	m := float64(f.Mass.Numerator) / float64(f.Mass.Denominator)
	r := float64(f.Radius.Numerator) / float64(f.Radius.Denominator)
	return -2*math.Cos(t)/(t*t) + 8*math.Sin(t)/(t*t*t) + 12*(math.Cos(t)-1)/(t*t*t*t) +
		2*(1-c)*(-math.Sin(t)/t-2*math.Cos(t)/(t*t)+2*math.Sin(t)/(t*t*t)) -
		2*w*math.Cos(t) - 2*m*r*r*math.Cos(r*t)
}

func DeriveM21ContactMotionEvidence(m, t1 float64) ContactMotionEvidence {
	r := 4.0
	phase := math.Pi * r / 2
	cPhase, sPhase := math.Cos(phase), math.Sin(phase)
	cPrime := -(math.Pi / 2) * sPhase
	bPrime := -2*sPhase - (math.Pi/2)*r*cPhase
	h := TwoAtomResponseFloat(r, t1)
	hr := -math.Pi*cPrime*math.Sin(t1)/t1 + 2*bPrime*math.Cos(t1) - 2*t1*math.Sin(r*t1)
	ht := -math.Pi*cPhase*(t1*math.Cos(t1)-math.Sin(t1))/(t1*t1) -
		2*(2*cPhase/math.Pi-r*sPhase)*math.Sin(t1) - 2*r*math.Sin(r*t1)
	htr := -math.Pi*cPrime*(t1*math.Cos(t1)-math.Sin(t1))/(t1*t1) -
		2*bPrime*math.Sin(t1) - 2*math.Sin(r*t1) - 2*r*t1*math.Cos(r*t1)
	f, _ := DeriveContactPreservingTwoAtomFamily(ExactRational{int64(math.Round(m * 1e17)), 100000000000000000}, ExactRational{4, 1})
	curvature := ContactPreservingTwoAtomSpectrumSecond(f, t1)
	dmDr := -m * hr / h
	dtDr := -(ht*dmDr + m*htr) / curvature
	return ContactMotionEvidence{
		DeltaCPerRadiusAtFixedMass: 0, DeltaWeight1PerRadiusAtFixedMass: -2 * math.Pi * m,
		DeltaFAtOldContact: m * hr, DeltaFPrimeAtOldContact: m * htr, ShiftedContactDerivative: dtDr,
		Conclusion: "decreasing r lifts the old M21 t1 contact; adjusting m by dm/dr=-m*(partial_r H)/H preserves its value, and the derivative equation moves the double contact smoothly",
	}
}

func DeriveM21Subfamily(m ExactRational) (ContactPreservingTwoAtomFamily, error) {
	return DeriveContactPreservingTwoAtomFamily(m, ExactRational{4, 1})
}

// There are no singular radius branches: the coefficient matrix of the two
// contact equations is independent of r and has determinant 8/pi.  Radius
// values only change the right-hand side through sine and cosine.
func ContactPreservationSingularRadii() string {
	return "none for finite r: determinant 8/pi; r<=1 is excluded by support legality, not singularity"
}

var m22Sqrt2 = rationalInterval{
	lo: big.NewRat(470832, 332929),
	hi: big.NewRat(665857, 470832),
}

// m22ParameterIntervals is deliberately limited to the simple certified
// radius 7/2.  Numerical radii never enter exact theorem checking.
func m22ParameterIntervals(f ContactPreservingTwoAtomFamily) (rationalInterval, rationalInterval, error) {
	if f.Radius != (ExactRational{7, 2}) {
		return rationalInterval{}, rationalInterval{}, fmt.Errorf("exact M22 verifier supports the certified radius 7/2 only")
	}
	base := DeriveOneRadiusTangencyCandidate()
	mI := riPoint(rat(f.Mass))
	sqrtHalf := riMul(m22Sqrt2, riPoint(big.NewRat(1, 2)))
	deltaC := riMul(riMul(evalPi(piMonomial(1, big.NewRat(1, 2))), mI), sqrtHalf)
	cI := riAdd(evalPi(base.ExactCExpression), deltaC)
	// At r=7/2, cos(7pi/4)=sqrt(2)/2 and sin(7pi/4)=-sqrt(2)/2.
	coefficient := riAdd(riMul(m22Sqrt2, riReciprocal(evalPi(piMonomial(1, big.NewRat(1, 1))))), riMul(m22Sqrt2, riPoint(big.NewRat(7, 4))))
	wI := riAdd(evalPi(base.ExactWExpression), riMul(mI, coefficient))
	return cI, wI, nil
}

type TwoAtomCompactPiece struct {
	Left         ExactRational   `json:"left"`
	Right        ExactRational   `json:"right"`
	GridStep     ExactRational   `json:"grid_step"`
	TaylorDegree int             `json:"taylor_degree"`
	Kind         CompactCellKind `json:"kind"`
}

type TwoAtomWholeLineCertificate struct {
	ID                    string                `json:"id"`
	OriginalContact       TwoAtomCompactPiece   `json:"original_contact_piece"`
	ContactGridStep       ExactRational         `json:"contact_curvature_grid_step"`
	ContactThirdAbsUpper  ExactRational         `json:"contact_third_abs_upper"`
	ContactConvexityLower ExactRational         `json:"contact_convexity_lower"`
	StrictPieces          []TwoAtomCompactPiece `json:"strict_compact_pieces"`
	LipschitzUpper        ExactRational         `json:"lipschitz_upper"`
	TailStart             ExactRational         `json:"tail_start"`
	TailLower             ExactRational         `json:"tail_lower"`
	Coverage              string                `json:"coverage"`
	WholeLine             bool                  `json:"whole_line"`
}

type CertifiedTwoAtomWitness struct {
	Family    ContactPreservingTwoAtomFamily `json:"family"`
	WholeLine TwoAtomWholeLineCertificate    `json:"whole_line"`
	CExact    string                         `json:"c_exact"`
	JExact    string                         `json:"J_exact"`
	Verified  bool                           `json:"verified"`
}

func twoAtomPointInterval(f ContactPreservingTwoAtomFamily, t *big.Rat, degree int) rationalInterval {
	cI, wI, _ := m22ParameterIntervals(f)
	if t.Sign() == 0 {
		return riMul(riPoint(big.NewRat(2, 1)), riAdd(riSub(riAdd(riPoint(big.NewRat(1, 1)), wI), cI), riPoint(rat(f.Mass))))
	}
	cosLo, cosHi := trigTaylorBounds(t, degree, true)
	sinLo, sinHi := trigTaylorBounds(t, degree-1, false)
	rt := new(big.Rat).Mul(t, rat(f.Radius))
	cosRLo, cosRHi := trigTaylorBounds(rt, degree, true)
	cosI, sinI, cosRI := ri(cosLo, cosHi), ri(sinLo, sinHi), ri(cosRLo, cosRHi)
	tI, t2I := riPoint(t), riPoint(new(big.Rat).Mul(t, t))
	value := riPoint(big.NewRat(1, 1))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riSub(cosI, riPoint(big.NewRat(1, 1))), riReciprocal(t2I))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riMul(riSub(riPoint(big.NewRat(1, 1)), cI), sinI), riReciprocal(tI))))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(wI, cosI)))
	return riAdd(value, riMul(riPoint(new(big.Rat).Mul(rat(f.Mass), big.NewRat(2, 1))), cosRI))
}

func twoAtomSecondPointInterval(f ContactPreservingTwoAtomFamily, t *big.Rat, degree int) rationalInterval {
	cI, wI, _ := m22ParameterIntervals(f)
	cosLo, cosHi := trigTaylorBounds(t, degree, true)
	sinLo, sinHi := trigTaylorBounds(t, degree-1, false)
	rt := new(big.Rat).Mul(t, rat(f.Radius))
	cosRLo, cosRHi := trigTaylorBounds(rt, degree, true)
	cosI, sinI, cosRI := ri(cosLo, cosHi), ri(sinLo, sinHi), ri(cosRLo, cosRHi)
	tI := riPoint(t)
	invT, invT2, invT3, invT4 := riReciprocal(tI), riPow(tI, -2), riPow(tI, -3), riPow(tI, -4)
	value := riMul(riPoint(big.NewRat(-2, 1)), riMul(cosI, invT2))
	value = riAdd(value, riMul(riPoint(big.NewRat(8, 1)), riMul(sinI, invT3)))
	value = riAdd(value, riMul(riPoint(big.NewRat(12, 1)), riMul(riSub(cosI, riPoint(big.NewRat(1, 1))), invT4)))
	inner := riAdd(riMul(riPoint(big.NewRat(-1, 1)), riMul(sinI, invT)), riMul(riPoint(big.NewRat(-2, 1)), riMul(cosI, invT2)))
	inner = riAdd(inner, riMul(riPoint(big.NewRat(2, 1)), riMul(sinI, invT3)))
	value = riAdd(value, riMul(riPoint(big.NewRat(2, 1)), riMul(riSub(riPoint(big.NewRat(1, 1)), cI), inner)))
	value = riAdd(value, riMul(riPoint(big.NewRat(-2, 1)), riMul(wI, cosI)))
	r2m2 := new(big.Rat).Mul(big.NewRat(-2, 1), new(big.Rat).Mul(rat(f.Mass), new(big.Rat).Mul(rat(f.Radius), rat(f.Radius))))
	return riAdd(value, riMul(riPoint(r2m2), cosRI))
}

func verifyTwoAtomPiece(f ContactPreservingTwoAtomFamily, p TwoAtomCompactPiece, lipschitz ExactRational) error {
	left, right, step := rat(p.Left), rat(p.Right), rat(p.GridStep)
	span := new(big.Rat).Sub(right, left)
	q := new(big.Rat).Quo(span, step)
	if left.Sign() < 0 || right.Cmp(left) <= 0 || step.Sign() <= 0 || !q.IsInt() || p.Kind != StrictPositiveCell || p.TaylorDegree < 80 {
		return fmt.Errorf("invalid strict two-atom piece")
	}
	loss := new(big.Rat).Mul(rat(lipschitz), new(big.Rat).Quo(step, big.NewRat(2, 1)))
	for k, count := int64(0), q.Num().Int64(); k <= count; k++ {
		t := new(big.Rat).Add(left, new(big.Rat).Mul(big.NewRat(k, 1), step))
		if twoAtomPointInterval(f, t, p.TaylorDegree).lo.Cmp(loss) <= 0 {
			return fmt.Errorf("strict two-atom piece [%s,%s] fails at %s", left, right, t)
		}
	}
	return nil
}

func VerifyTwoAtomWholeLine(f ContactPreservingTwoAtomFamily, p TwoAtomWholeLineCertificate, replayCompact bool) error {
	if err := f.Validate(); err != nil {
		return err
	}
	cI, wI, err := m22ParameterIntervals(f)
	if err != nil {
		return err
	}
	if p.ID == "" || !p.WholeLine || p.OriginalContact.Kind != NonnegativeCellWithCertifiedContact || p.TailStart != (ExactRational{4, 1}) || len(p.StrictPieces) == 0 {
		return fmt.Errorf("incomplete M22 whole-line certificate")
	}
	m22Coverage := p.Coverage == "evenness; strict rational Taylor cells on [0,154/100] and [160/100,4]; certified convex original-contact cell [154/100,160/100]; radius-uniform analytic tail [4,infinity)"
	m23Coverage := p.Coverage == "evenness; adaptive strict rational Taylor cells on [0,154/100] and [161/100,4]; certified convex original-contact cell [154/100,161/100]; radius-uniform analytic tail [4,infinity)"
	if p.OriginalContact.GridStep != p.ContactGridStep || p.OriginalContact.TaylorDegree < 80 || (!m22Coverage && !m23Coverage) {
		return fmt.Errorf("M22 contact or coverage declaration is incomplete")
	}
	cursor := big.NewRat(0, 1)
	contactLeft, contactRight := rat(p.OriginalContact.Left), rat(p.OriginalContact.Right)
	for _, piece := range p.StrictPieces {
		if cursor.Cmp(contactLeft) == 0 {
			cursor = new(big.Rat).Set(contactRight)
		}
		if rat(piece.Left).Cmp(cursor) != 0 {
			return fmt.Errorf("M22 compact partition has a gap before %s", rat(piece.Left))
		}
		cursor = rat(piece.Right)
	}
	if cursor.Cmp(contactLeft) == 0 {
		cursor = new(big.Rat).Set(contactRight)
	}
	if cursor.Cmp(rat(p.TailStart)) != 0 {
		return fmt.Errorf("M22 compact partition does not reach the tail")
	}
	// Radius-uniform derivative and tail bounds use only |sin|, |cos|<=1.
	lipschitz := riAdd(riSub(cI, riPoint(big.NewRat(2, 3))), riMul(wI, riPoint(big.NewRat(2, 1))))
	lipschitz = riAdd(lipschitz, riPoint(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Mul(rat(f.Mass), rat(f.Radius)))))
	if lipschitz.hi.Cmp(rat(p.LipschitzUpper)) >= 0 {
		return fmt.Errorf("M22 Lipschitz declaration is too small")
	}
	third := riAdd(riSub(riMul(cI, riPoint(big.NewRat(1, 2))), riPoint(big.NewRat(2, 5))), riMul(wI, riPoint(big.NewRat(2, 1))))
	r3 := new(big.Rat).Mul(rat(f.Radius), new(big.Rat).Mul(rat(f.Radius), rat(f.Radius)))
	third = riAdd(third, riPoint(new(big.Rat).Mul(big.NewRat(2, 1), new(big.Rat).Mul(rat(f.Mass), r3))))
	if third.hi.Cmp(rat(p.ContactThirdAbsUpper)) >= 0 {
		return fmt.Errorf("M22 contact third-derivative declaration is too small")
	}
	left, right, step := contactLeft, contactRight, rat(p.ContactGridStep)
	if new(big.Rat).Quo(m18Pi.lo, big.NewRat(2, 1)).Cmp(left) <= 0 || new(big.Rat).Quo(m18Pi.hi, big.NewRat(2, 1)).Cmp(right) >= 0 {
		return fmt.Errorf("M22 original-contact cell misses pi/2")
	}
	q := new(big.Rat).Quo(new(big.Rat).Sub(right, left), step)
	if !q.IsInt() {
		return fmt.Errorf("M22 contact curvature grid is inexact")
	}
	var minSecond *big.Rat
	for k, count := int64(0), q.Num().Int64(); k <= count; k++ {
		t := new(big.Rat).Add(left, new(big.Rat).Mul(big.NewRat(k, 1), step))
		lo := twoAtomSecondPointInterval(f, t, p.OriginalContact.TaylorDegree).lo
		if minSecond == nil || lo.Cmp(minSecond) < 0 {
			minSecond = lo
		}
	}
	minSecond.Sub(minSecond, new(big.Rat).Mul(rat(p.ContactThirdAbsUpper), new(big.Rat).Quo(step, big.NewRat(2, 1))))
	if minSecond.Cmp(rat(p.ContactConvexityLower)) <= 0 {
		return fmt.Errorf("M22 original-contact convexity lower bound failed")
	}
	if replayCompact {
		for _, piece := range p.StrictPieces {
			if err := verifyTwoAtomPiece(f, piece, p.LipschitzUpper); err != nil {
				return err
			}
		}
	}
	tail := riSub(riSub(riPoint(big.NewRat(3, 4)), riMul(riSub(cI, riPoint(big.NewRat(1, 1))), riPoint(big.NewRat(1, 2)))), riMul(wI, riPoint(big.NewRat(2, 1))))
	tail = riSub(tail, riPoint(new(big.Rat).Mul(rat(f.Mass), big.NewRat(2, 1))))
	if tail.lo.Cmp(rat(p.TailLower)) <= 0 {
		return fmt.Errorf("M22 radius-uniform tail bound failed")
	}
	return nil
}

func VerifyCertifiedTwoAtomWitness(w CertifiedTwoAtomWitness, replayCompact bool) error {
	if err := VerifyTwoAtomWholeLine(w.Family, w.WholeLine, replayCompact); err != nil {
		return err
	}
	if w.CExact != "1+pi/4-2/pi+17*pi*sqrt(2)/10000" || w.JExact != "1-pi/4+2/pi-17*pi*sqrt(2)/10000" || !w.Verified {
		return fmt.Errorf("wrong M22 exact theorem consequence")
	}
	cI, _, _ := m22ParameterIntervals(w.Family)
	m21, _ := DeriveContactAdaptedAtomicRay(ExactRational{19, 5000})
	if cI.lo.Cmp(evalPi(m21.C).hi) <= 0 || cI.hi.Cmp(big.NewRat(1651, 1250)) >= 0 {
		return fmt.Errorf("M22 witness does not improve M21 safely or violates the primal guard")
	}
	return nil
}
