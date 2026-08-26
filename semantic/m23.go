package semantic

import (
	"fmt"
	"math/big"
)

// CertifiedScalarInterval is the report-facing decimal rendering of a
// directed rational interval.  The verifier never parses these strings.
type CertifiedScalarInterval struct {
	Lower string `json:"lower"`
	Upper string `json:"upper"`
}

type TwoAtomOptimizerSystem struct {
	Q1                 string `json:"Q1"`
	Q2                 string `json:"Q2"`
	Q3                 string `json:"Q3"`
	EnvelopeDerivation string `json:"envelope_derivation"`
	LagrangeDerivation string `json:"lagrange_derivation"`
}

type CertifiedOptimizerBox struct {
	Mass                CertifiedScalarInterval    `json:"m"`
	Radius              CertifiedScalarInterval    `json:"r"`
	Contact             CertifiedScalarInterval    `json:"t"`
	C                   CertifiedScalarInterval    `json:"c"`
	JacobianDeterminant CertifiedScalarInterval    `json:"jacobian_determinant"`
	KrawczykImage       [3]CertifiedScalarInterval `json:"krawczyk_image"`
	InitialBox          string                     `json:"initial_box"`
	Method              string                     `json:"method"`
	ExistsUnique        bool                       `json:"exists_unique"`
}

type LocalOptimalityEvidence struct {
	FixedContactCurvature CertifiedScalarInterval `json:"fixed_contact_curvature"`
	MovableCurvature      CertifiedScalarInterval `json:"movable_contact_curvature"`
	EnvelopeSecond        CertifiedScalarInterval `json:"envelope_second_derivative"`
	TailLower             CertifiedScalarInterval `json:"tail_lower"`
	BranchJacobian        CertifiedScalarInterval `json:"branch_jacobian"`
	ContactNeighborhood   string                  `json:"contact_neighborhood"`
	RadiusSignStructure   string                  `json:"radius_sign_structure"`
	ParameterLegality     string                  `json:"parameter_legality"`
	CertifiedLocalMaximum bool                    `json:"certified_local_maximum"`
}

// m23Jet is a deliberately three-variable, first-order interval jet.  It is
// not a generic solver abstraction: its only client is the M23 Q-system.
type m23Jet struct {
	v rationalInterval
	d [3]rationalInterval // m,r,t
}

func m23Zero() rationalInterval { return riPoint(big.NewRat(0, 1)) }
func m23Point(q *big.Rat) m23Jet {
	return m23Jet{v: riPoint(q), d: [3]rationalInterval{m23Zero(), m23Zero(), m23Zero()}}
}
func m23Interval(x rationalInterval) m23Jet {
	return m23Jet{v: x, d: [3]rationalInterval{m23Zero(), m23Zero(), m23Zero()}}
}
func m23Var(x rationalInterval, k int) m23Jet {
	z := m23Interval(x)
	z.d[k] = riPoint(big.NewRat(1, 1))
	return z
}
func m23Add(x, y m23Jet) m23Jet {
	z := m23Jet{v: riAdd(x.v, y.v)}
	for k := range z.d {
		z.d[k] = riAdd(x.d[k], y.d[k])
	}
	return z
}
func m23Neg(x m23Jet) m23Jet {
	z := m23Jet{v: riNeg(x.v)}
	for k := range z.d {
		z.d[k] = riNeg(x.d[k])
	}
	return z
}
func m23Sub(x, y m23Jet) m23Jet { return m23Add(x, m23Neg(y)) }
func m23Mul(x, y m23Jet) m23Jet {
	z := m23Jet{v: riMul(x.v, y.v)}
	for k := range z.d {
		z.d[k] = riAdd(riMul(x.d[k], y.v), riMul(x.v, y.d[k]))
	}
	return z
}
func m23Inv(x m23Jet) m23Jet {
	z := m23Jet{v: riReciprocal(x.v)}
	sq := riMul(z.v, z.v)
	for k := range z.d {
		z.d[k] = riNeg(riMul(x.d[k], sq))
	}
	return z
}
func m23Div(x, y m23Jet) m23Jet { return m23Mul(x, m23Inv(y)) }

func m23TrigRange(x rationalInterval, cosine bool) rationalInterval {
	mid := new(big.Rat).Quo(new(big.Rat).Add(x.lo, x.hi), big.NewRat(2, 1))
	rad := new(big.Rat).Quo(new(big.Rat).Sub(x.hi, x.lo), big.NewRat(2, 1))
	lo, hi := trigTaylorBounds(mid, 80, cosine)
	lo.Sub(lo, rad)
	hi.Add(hi, rad) // |sin'|,|cos'| <= 1
	if lo.Cmp(big.NewRat(-1, 1)) < 0 {
		lo = big.NewRat(-1, 1)
	}
	if hi.Cmp(big.NewRat(1, 1)) > 0 {
		hi = big.NewRat(1, 1)
	}
	return ri(lo, hi)
}
func m23Sin(x m23Jet) m23Jet {
	s, c := m23TrigRange(x.v, false), m23TrigRange(x.v, true)
	z := m23Jet{v: s}
	for k := range z.d {
		z.d[k] = riMul(c, x.d[k])
	}
	return z
}
func m23Cos(x m23Jet) m23Jet {
	c, s := m23TrigRange(x.v, true), m23TrigRange(x.v, false)
	z := m23Jet{v: c}
	for k := range z.d {
		z.d[k] = riNeg(riMul(s, x.d[k]))
	}
	return z
}

func m23System(box [3]rationalInterval) ([3]m23Jet, m23Jet, m23Jet, m23Jet) {
	m, r, t := m23Var(box[0], 0), m23Var(box[1], 1), m23Var(box[2], 2)
	one, two := m23Point(big.NewRat(1, 1)), m23Point(big.NewRat(2, 1))
	pi := m23Interval(m18Pi)
	a := m23Mul(pi, m23Point(big.NewRat(1, 2)))
	phase := m23Mul(a, r)
	A, S := m23Cos(phase), m23Sin(phase)
	Ap := m23Neg(m23Mul(a, S))
	B := m23Sub(m23Mul(m23Div(two, pi), A), m23Mul(r, S))
	Bp := m23Sub(m23Neg(m23Mul(two, S)), m23Mul(m23Mul(a, r), A))
	sint, cost := m23Sin(t), m23Cos(t)
	sinc := m23Div(sint, t)
	rt := m23Mul(r, t)
	sinrt, cosrt := m23Sin(rt), m23Cos(rt)
	h := m23Add(m23Add(m23Neg(m23Mul(m23Mul(pi, A), sinc)), m23Mul(m23Mul(two, B), cost)), m23Mul(two, cosrt))
	hr := m23Add(m23Add(m23Neg(m23Mul(m23Mul(pi, Ap), sinc)), m23Mul(m23Mul(two, Bp), cost)), m23Neg(m23Mul(m23Mul(two, t), sinrt)))

	c0 := m23Add(m23Sub(one, m23Div(two, pi)), m23Mul(pi, m23Point(big.NewRat(1, 4))))
	w0 := m23Add(m23Sub(m23Div(one, pi), m23Div(m23Point(big.NewRat(4, 1)), m23Mul(pi, pi))), m23Div(m23Point(big.NewRat(8, 1)), m23Mul(m23Mul(pi, pi), pi)))
	f0 := m23Add(m23Add(one, m23Div(m23Mul(two, m23Sub(cost, one)), m23Mul(t, t))), m23Add(m23Div(m23Mul(m23Mul(two, m23Sub(one, c0)), sint), t), m23Mul(m23Mul(two, w0), cost)))
	f0t := m23Add(m23Add(m23Div(m23Neg(m23Mul(two, sint)), m23Mul(t, t)), m23Div(m23Neg(m23Mul(m23Point(big.NewRat(4, 1)), m23Sub(cost, one))), m23Mul(m23Mul(t, t), t))), m23Add(m23Mul(m23Mul(two, m23Sub(one, c0)), m23Sub(m23Div(cost, t), m23Div(sint, m23Mul(t, t)))), m23Neg(m23Mul(m23Mul(two, w0), sint))))
	ht := m23Add(m23Add(m23Neg(m23Div(m23Mul(m23Mul(pi, A), m23Sub(m23Mul(t, cost), sint)), m23Mul(t, t))), m23Neg(m23Mul(m23Mul(two, B), sint))), m23Neg(m23Mul(m23Mul(two, r), sinrt)))
	q1 := m23Add(f0, m23Mul(m, h))
	q2 := m23Add(f0t, m23Mul(m, ht))
	q3 := m23Sub(m23Mul(Ap, h), m23Mul(A, hr))
	return [3]m23Jet{q1, q2, q3}, h, hr, A
}

func m23Rat(s string) *big.Rat {
	q, ok := new(big.Rat).SetString(s)
	if !ok {
		panic(s)
	}
	return q
}
func m23Box(lo, hi string) rationalInterval { return ri(m23Rat(lo), m23Rat(hi)) }
func m23Decimal(q *big.Rat) string          { return q.FloatString(16) }
func m23Display(x rationalInterval) CertifiedScalarInterval {
	return CertifiedScalarInterval{m23Decimal(x.lo), m23Decimal(x.hi)}
}
func m23Scale(q *big.Rat, x rationalInterval) rationalInterval { return riMul(riPoint(q), x) }
func m23ContainsInterior(outer, inner rationalInterval) bool {
	return inner.lo.Cmp(outer.lo) > 0 && inner.hi.Cmp(outer.hi) < 0
}

func m23Det(j [3][3]rationalInterval) rationalInterval {
	a := riMul(j[0][0], riSub(riMul(j[1][1], j[2][2]), riMul(j[1][2], j[2][1])))
	b := riMul(j[0][1], riSub(riMul(j[1][0], j[2][2]), riMul(j[1][2], j[2][0])))
	c := riMul(j[0][2], riSub(riMul(j[1][0], j[2][1]), riMul(j[1][1], j[2][0])))
	return riAdd(riSub(a, b), c)
}

func m23Second(box [3]rationalInterval) m23Jet {
	m, r, t := m23Var(box[0], 0), m23Var(box[1], 1), m23Var(box[2], 2)
	one, two := m23Point(big.NewRat(1, 1)), m23Point(big.NewRat(2, 1))
	pi := m23Interval(m18Pi)
	a := m23Mul(pi, m23Point(big.NewRat(1, 2)))
	A, S := m23Cos(m23Mul(a, r)), m23Sin(m23Mul(a, r))
	B := m23Sub(m23Mul(m23Div(two, pi), A), m23Mul(r, S))
	c0 := m23Add(m23Sub(one, m23Div(two, pi)), m23Mul(pi, m23Point(big.NewRat(1, 4))))
	w0 := m23Add(m23Sub(m23Div(one, pi), m23Div(m23Point(big.NewRat(4, 1)), m23Mul(pi, pi))), m23Div(m23Point(big.NewRat(8, 1)), m23Mul(m23Mul(pi, pi), pi)))
	c := m23Add(c0, m23Mul(m23Mul(a, m), A))
	w := m23Add(w0, m23Mul(m, B))
	s, co := m23Sin(t), m23Cos(t)
	t2, t3, t4 := m23Mul(t, t), m23Mul(m23Mul(t, t), t), m23Mul(m23Mul(t, t), m23Mul(t, t))
	v := m23Add(m23Add(m23Neg(m23Div(m23Mul(two, co), t2)), m23Div(m23Mul(m23Point(big.NewRat(8, 1)), s), t3)), m23Div(m23Mul(m23Point(big.NewRat(12, 1)), m23Sub(co, one)), t4))
	inner := m23Add(m23Neg(m23Div(s, t)), m23Add(m23Neg(m23Div(m23Mul(two, co), t2)), m23Div(m23Mul(two, s), t3)))
	v = m23Add(v, m23Mul(m23Mul(two, m23Sub(one, c)), inner))
	v = m23Add(v, m23Neg(m23Mul(m23Mul(two, w), co)))
	return m23Add(v, m23Neg(m23Mul(m23Mul(m23Mul(two, m), m23Mul(r, r)), m23Cos(m23Mul(r, t)))))
}

// CertifyVariableRadiusOptimizer applies a rational Krawczyk operator to one
// fixed box and proves the local branch maximum with the implicit-function
// theorem.  All transcendental enclosures use Taylor bounds plus |sin'|,|cos'|<=1.
func CertifyVariableRadiusOptimizer() (TwoAtomOptimizerSystem, CertifiedOptimizerBox, LocalOptimalityEvidence, error) {
	system := TwoAtomOptimizerSystem{
		Q1: "F0(t)+m H_r(t)", Q2: "F0'(t)+m partial_t H_r(t)", Q3: "A'(r)H_r(t)-A(r)partial_r H_r(t), A(r)=cos(pi*r/2)",
		EnvelopeDerivation: "F=0 and F_t=0 imply m'=-m H_r/H; therefore C'=(pi*m/2)*(A'H-A H_r)/H",
		LagrangeDerivation: "stationarity of c subject to F=F_t=0 gives the t multiplier zero when F_tt!=0; eliminating the remaining multiplier yields c_r F_m-c_m F_r=(pi*m/2)*(A'H-A H_r)=0",
	}
	initial := "m in [0.0069,0.0071], r in [3.50,3.52], t in [1.78,1.81]"
	box := [3]rationalInterval{m23Box("0.00699508", "0.00699510"), m23Box("3.507850", "3.507857"), m23Box("1.795310", "1.795317")}
	center := [3]rationalInterval{}
	for k := range box {
		mid := new(big.Rat).Quo(new(big.Rat).Add(box[k].lo, box[k].hi), big.NewRat(2, 1))
		center[k] = riPoint(mid)
	}
	qc, _, _, _ := m23System(center)
	qx, h, hr, A := m23System(box)
	var j [3][3]rationalInterval
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			j[i][k] = qx[i].d[k]
		}
	}
	// Rounded inverse of the binary64 center Jacobian.  It is only a
	// preconditioner; the rational inclusion below is the proof.
	rs := [][]string{{"6.569474360", "-0.911678766", "-0.008731100"}, {"-793.969240442", "85.044206178", "0.814463856"}, {"-342.874685773", "36.726240747", "-0.831114972"}}
	var R [3][3]*big.Rat
	for i := 0; i < 3; i++ {
		for k := 0; k < 3; k++ {
			R[i][k] = m23Rat(rs[i][k])
		}
	}
	var image [3]rationalInterval
	for i := 0; i < 3; i++ {
		base := center[i]
		for k := 0; k < 3; k++ {
			base = riSub(base, m23Scale(R[i][k], qc[k].v))
		}
		corr := m23Zero()
		for k := 0; k < 3; k++ {
			coef := riPoint(big.NewRat(0, 1))
			if i == k {
				coef = riPoint(big.NewRat(1, 1))
			}
			for l := 0; l < 3; l++ {
				coef = riSub(coef, m23Scale(R[i][l], j[l][k]))
			}
			delta := riSub(box[k], center[k])
			corr = riAdd(corr, riMul(coef, delta))
		}
		image[i] = riAdd(base, corr)
		if !m23ContainsInterior(box[i], image[i]) {
			return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("M23 Krawczyk image %d is not interior: [%s,%s]", i, image[i].lo, image[i].hi)
		}
	}
	det := m23Det(j)
	if det.lo.Sign() <= 0 && det.hi.Sign() >= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("M23 Jacobian determinant contains zero")
	}
	second := m23Second(box).v
	if second.lo.Sign() <= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("movable curvature not positive")
	}
	// Fixed-contact curvature is evaluated on the full parameter box and the
	// directed pi/2 interval.
	fixedBox := box
	fixedBox[2] = ri(new(big.Rat).Quo(m18Pi.lo, big.NewRat(2, 1)), new(big.Rat).Quo(m18Pi.hi, big.NewRat(2, 1)))
	fixedSecond := m23Second(fixedBox).v
	if fixedSecond.lo.Sign() <= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("fixed curvature not positive")
	}
	// IFT branch derivatives.  Here H_r is the radius derivative (hr), while
	// qx[1].d[1]=F_tr and qx[1].d[2]=F_tt.
	mp := riNeg(riMul(riMul(box[0], hr.v), riReciprocal(h.v)))
	tp := riNeg(riMul(riAdd(riMul(qx[1].d[0], mp), qx[1].d[1]), riReciprocal(qx[1].d[2])))
	gp := riAdd(qx[2].d[1], riMul(qx[2].d[2], tp))
	piHalf := riMul(m18Pi, riPoint(big.NewRat(1, 2)))
	csecond := riMul(riMul(riMul(piHalf, box[0]), gp), riReciprocal(h.v))
	if csecond.hi.Sign() >= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("radius envelope second derivative not negative: %s", csecond.hi)
	}
	branchDet := riMul(h.v, second)
	if branchDet.lo.Sign() <= 0 && branchDet.hi.Sign() >= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("branch Jacobian singular")
	}
	c0 := riAdd(riSub(riPoint(big.NewRat(1, 1)), riMul(riPoint(big.NewRat(2, 1)), riReciprocal(m18Pi))), riMul(m18Pi, riPoint(big.NewRat(1, 4))))
	cI := riAdd(c0, riMul(riMul(piHalf, box[0]), A.v))
	// Legality and tail use the same structural formulas as M22.
	phase := riMul(piHalf, box[1])
	cp, sp := m23TrigRange(phase, true), m23TrigRange(phase, false)
	w0 := riAdd(riSub(riReciprocal(m18Pi), riMul(riPoint(big.NewRat(4, 1)), riPow(m18Pi, -2))), riMul(riPoint(big.NewRat(8, 1)), riPow(m18Pi, -3)))
	wI := riAdd(w0, riMul(box[0], riSub(riMul(riPoint(big.NewRat(2, 1)), riMul(cp, riReciprocal(m18Pi))), riMul(box[1], sp))))
	tail := riSub(riSub(riPoint(big.NewRat(3, 4)), riMul(riSub(cI, riPoint(big.NewRat(1, 1))), riPoint(big.NewRat(1, 2)))), riMul(wI, riPoint(big.NewRat(2, 1))))
	tail = riSub(tail, riMul(box[0], riPoint(big.NewRat(2, 1))))
	if box[0].lo.Sign() <= 0 || box[1].lo.Cmp(big.NewRat(1, 1)) <= 0 || cI.lo.Sign() <= 0 || wI.lo.Sign() <= 0 || tail.lo.Sign() <= 0 {
		return system, CertifiedOptimizerBox{}, LocalOptimalityEvidence{}, fmt.Errorf("M23 box legality or tail failed")
	}
	out := CertifiedOptimizerBox{Mass: m23Display(box[0]), Radius: m23Display(box[1]), Contact: m23Display(box[2]), C: m23Display(cI), JacobianDeterminant: m23Display(det), InitialBox: initial, Method: "three-variable rational Krawczyk inclusion with degree-80 Taylor trigonometric ranges", ExistsUnique: true}
	for k := 0; k < 3; k++ {
		out.KrawczykImage[k] = m23Display(image[k])
	}
	evidence := LocalOptimalityEvidence{FixedContactCurvature: m23Display(fixedSecond), MovableCurvature: m23Display(second), EnvelopeSecond: m23Display(csecond), TailLower: m23Display(tail), BranchJacobian: m23Display(branchDet), ContactNeighborhood: "on 1.795310<=t<=1.795317, F_tt>0; hence F_t is strictly increasing and the double zero is the unique zero there", RadiusSignStructure: "the IFT branch exists because H*F_tt excludes zero; C''<0 throughout the isolating cylinder, so C' changes from positive to negative at the unique stationary radius", ParameterLegality: "box-wide r>1, m>0, c>0, w1>0; fixed contact equations are structural", CertifiedLocalMaximum: true}
	return system, out, evidence, nil
}

type ExplicitTwoAtomWitness struct {
	Witness     CertifiedTwoAtomWitness `json:"witness"`
	CInterval   CertifiedScalarInterval `json:"c_interval"`
	ImprovesM22 bool                    `json:"improves_m22"`
}

// M23ExplicitWitness deliberately keeps r=7/2 and raises only the rational
// mass to 69/10000.  This retains M22's cheap Pell arithmetic while leaving a
// strict margin at the movable minimum.
func M23ExplicitWitness() ExplicitTwoAtomWitness {
	f, err := DeriveContactPreservingTwoAtomFamily(ExactRational{69, 10000}, ExactRational{7, 2})
	if err != nil {
		panic(err)
	}
	p := TwoAtomWholeLineCertificate{
		ID:              "m23-r7-over-2-m69-over-10000-degree100",
		OriginalContact: TwoAtomCompactPiece{Left: ExactRational{154, 100}, Right: ExactRational{161, 100}, GridStep: ExactRational{1, 10000}, TaylorDegree: 100, Kind: NonnegativeCellWithCertifiedContact},
		ContactGridStep: ExactRational{1, 10000}, ContactThirdAbsUpper: ExactRational{6, 5}, ContactConvexityLower: ExactRational{1, 1000},
		StrictPieces: []TwoAtomCompactPiece{
			{Left: ExactRational{0, 1}, Right: ExactRational{7, 5}, GridStep: ExactRational{1, 5000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{7, 5}, Right: ExactRational{154, 100}, GridStep: ExactRational{1, 100000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{161, 100}, Right: ExactRational{17, 10}, GridStep: ExactRational{1, 100000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{17, 10}, Right: ExactRational{19, 10}, GridStep: ExactRational{1, 200000}, TaylorDegree: 100, Kind: StrictPositiveCell},
			{Left: ExactRational{19, 10}, Right: ExactRational{4, 1}, GridStep: ExactRational{1, 5000}, TaylorDegree: 100, Kind: StrictPositiveCell},
		},
		LipschitzUpper: ExactRational{19, 20}, TailStart: ExactRational{4, 1}, TailLower: ExactRational{1, 4},
		Coverage: "evenness; adaptive strict rational Taylor cells on [0,154/100] and [161/100,4]; certified convex original-contact cell [154/100,161/100]; radius-uniform analytic tail [4,infinity)", WholeLine: true,
	}
	w := CertifiedTwoAtomWitness{Family: f, WholeLine: p, CExact: "1+pi/4-2/pi+69*pi*sqrt(2)/40000", JExact: "1-pi/4+2/pi-69*pi*sqrt(2)/40000", Verified: true}
	cI, _, _ := m22ParameterIntervals(f)
	old, _ := DeriveContactPreservingTwoAtomFamily(ExactRational{17, 2500}, ExactRational{7, 2})
	oldC, _, _ := m22ParameterIntervals(old)
	return ExplicitTwoAtomWitness{Witness: w, CInterval: m23Display(cI), ImprovesM22: cI.lo.Cmp(oldC.hi) > 0}
}

func VerifyM23ExplicitWitness(w ExplicitTwoAtomWitness, replayCompact bool) error {
	if w.Witness.CExact != "1+pi/4-2/pi+69*pi*sqrt(2)/40000" || w.Witness.JExact != "1-pi/4+2/pi-69*pi*sqrt(2)/40000" || !w.Witness.Verified || !w.ImprovesM22 {
		return fmt.Errorf("wrong M23 explicit witness metadata")
	}
	if err := VerifyTwoAtomWholeLine(w.Witness.Family, w.Witness.WholeLine, replayCompact); err != nil {
		return err
	}
	cI, _, _ := m22ParameterIntervals(w.Witness.Family)
	if cI.hi.Cmp(big.NewRat(1651, 1250)) >= 0 {
		return fmt.Errorf("M23 witness violates primal guard")
	}
	return nil
}
