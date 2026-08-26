# M22: release the second contact with one radius degree of freedom

## Outcome

M22 achieves Success A.  It stays in the same two-atom support class, releases
only the second radius, and certifies

```text
c_* >= 1+pi/4-2/pi+17*pi*sqrt(2)/10000
J_* <= 1-pi/4+2/pi-17*pi*sqrt(2)/10000.
```

The certified values are approximately `c=1.15633129202474` and
`J<=0.843668707975264`.  This improves M21's certified
`c=1.154747417...` and its numerical ray ceiling `1.154860934...`.

## Exact contact-preserving surface

With the repository's plus-sign Fourier convention,

```text
F(t)=1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t
     +2w1 cos(t)+2m cos(rt).
```

Writing `a=pi/2`, `C=cos(ar)`, and `S=sin(ar)`, the equations
`F(a)=F'(a)=0` give

```text
c(m,r)  = c0+a m C,
w1(m,r) = w0+m(C/a-rS)
        = w0+m((2/pi)C-rS),
w2       = m.
```

Consequently `F_{m,r}=F0+m H_r`, where

```text
H_r(t)=-pi C sin(t)/t+2((2/pi)C-rS)cos(t)+2cos(rt).
```

The coefficient matrix for `(c,w1)` is independent of `r` and has determinant
`8/pi`.  There are no finite singular radius branches.  The legal surface is
cut out by `r>1`, `m>=0`, `c(m,r)>0`, and `w1(m,r)>=0`; theorem propagation
also enforces `c<1651/1250`.  The analytic-tail route additionally requires
`3/4-(c-1)/2-2w1-2m>0`.

At `r=4`, `C=1` and `S=0`, recovering exactly
`c=c0+(pi/2)m`, `w1=w0+(2/pi)m`, and the M21 response.  This regression is a
typed constructor, not a report-only assertion.

The integer phase structure is informative: `r=2` decreases the objective;
`r=3` leaves `c=c0` and gives `w1=w0+3m`; `r=4` is M21; and `r=5` leaves
`c=c0` and gives `w1=w0-5m`.  Nearby nonintegers with positive `C` can improve
the objective, and the best local candidate is nonintegral.

## Sensitivity and contact motion

At the M21 numerical double contact, fixed-mass radius differentiation gives

```text
partial_r c       = 0
partial_r w1      ~= -0.0243301732790
partial_r F(t1)   ~= -0.00219765236555
partial_r F'(t1)  ~=  0.0246026460356.
```

Thus decreasing `r` raises the old `t1` contact.  A shift of `t1` alone cannot
preserve the value to first order because `F'(t1)=0`.  Coupling the mass by

```text
dm/dr = -m (partial_r H_r)/H_r
```

preserves `F(t1)=0`, after which

```text
dt1/dr = -(H_t dm/dr+m H_tr)/F_tt ~= -0.742981381362
```

moves the double contact smoothly.  This proves at first order that the M21
ray is not locally optimal in the contact-preserving two-atom surface.

## Nested envelope and finite optimizer

For fixed `r` and `H_r(t)<0`, define

```text
M_r(t)=F0(t)/(-H_r(t)),       m_*(r)=inf_t M_r(t),
C(r)=c0+(pi/2)cos(pi r/2)m_*(r).
```

At an interior optimum the three unknowns `(m,r,t1)` satisfy

```text
Q1=F0(t1)+m H_r(t1)=0,
Q2=F0'(t1)+m partial_t H_r(t1)=0,
Q3=Cphase'(r)H_r(t1)-Cphase(r)partial_r H_r(t1)=0,
```

where `Cphase(r)=cos(pi r/2)`.  The first two fixed-contact equations have
already eliminated `(c,w1)`.  Constraint counting is therefore: four original
scalars, minus two fixed-contact equations, leaving `(m,r)`; a movable double
contact introduces its frequency and two equations, leaving one degree; the
radius-envelope stationarity equation closes the interior candidate system.

The bounded search on `r in [3,5]` finds

```text
r*  ~= 3.507853493599018
m*  ~= 0.006995089360323918
w1  ~= 0.19136037852198656
c   ~= 1.1566432357279466
t1  ~= 1.7953133218061055
F'' ~= 0.008052746837405865.
```

The contact moves smoothly from `t1~=1.36475` at `r=4` through `1.46139` at
`r=3.875`, `1.56610` at `r=3.75`, and `1.80340` at `r=3.5`.  No widening was
needed because the radius response turns over inside the local window.  No
third near-active minimum appears: the first remote minima near `9.44` and
`15.54` are approximately `0.571` and `0.597`.

## Certified simple witness

The half-integer radius `r=7/2` is nearly optimal and has algebraic phases
`cos(7pi/4)=sqrt(2)/2`, `sin(7pi/4)=-sqrt(2)/2`.  With `m=17/2500`,

```text
c  = c0+17*pi*sqrt(2)/10000,
w1 = w0+(17/2500)*sqrt(2)*(1/pi+7/4).
```

The exact verifier encloses `sqrt(2)` by the adjacent Pell convergents
`470832/332929 < sqrt(2) < 665857/470832`.  This narrow algebraic interval is
kept local to M22 because the general `PiExpression` IR does not represent
algebraic extensions.

The whole-line certificate uses:

- an equality-aware cell `[154/100,160/100]`, where exact contact identities
  and degree-100 curvature samples prove `F''>1/1000`;
- strict degree-100 rational Taylor/Lipschitz cells covering `[0,154/100]` and
  `[160/100,4]`, with `|F'|<19/20`;
- the radius-uniform tail
  `F>=3/4-(c-1)/2-2w1-2m>1/4` on `|t|>=4`;
- evenness for negative frequencies.

The released minimum near `t=1.8034` is deliberately strict, approximately
`1.53e-4`, so it needs no equality-aware second-contact cell.  The full replay
took about 146 seconds, versus M21's reported roughly 235 seconds, without
weakening the degree or proof semantics.

## Evidence, literature, and compiler lesson

The Oct experiment passes 5/5 interpreted Facts in 5.029 seconds and 5/5
compiled Facts in 0.624 seconds with zero fallback.  Three deterministic
`[Artifact]` plots record the radius response, contact motion, and best
spectrum.  The compiled artifact spelling delegates, as documented by Oct, to
the build-time interpreter; its replay reported all three outputs unchanged.
Float searches and plots are research evidence only.

A limited post-certification screen of the primary/preprint sources
[Chirre--Goncalves--de Laat](https://arxiv.org/abs/1810.08843) and
[Das--Ismoilov--Ramos](https://arxiv.org/abs/2502.05106), plus targeted
atomic/Dirac/variable-radius searches, found SDP Gaussian-polynomial and
reproducing-kernel/Fourier-optimization approaches but no stated match for this
two-atom contact-preserving family, `r=7/2` witness, or three-equation active
system.  This is not a priority or novelty claim.

M22 teaches a concrete Compiler Theory pattern: preserve established equality
constraints structurally, expose one geometry parameter, compile the infinite
constraint into a nested scalar envelope, and keep numerical candidates,
contact-motion evidence, and certified witnesses as different types.  The only
architectural awkwardness is the deliberately local `sqrt(2)` interval layer.

RH remains unresolved.  The one proposed next milestone is M23: certify a
tight box for the three-equation variable-radius optimizer and determine
whether its two-contact configuration is the local two-atom optimum.  M22 does
not add a third atom or begin that work.
