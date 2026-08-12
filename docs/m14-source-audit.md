# M14 source audit: the claimed 0.68185 ceiling is not source-certified

> M14A continuation: [m14a-architecture.md](m14a-architecture.md) reconstructs
> the exact EP3.1 primal, a positive-definite weak dual, the exact `c=1`
> baseline witness, and the remaining whole-line/tail obstruction.  The
> `0.68185` remark remains uncertified.

M14 cannot yet be implemented as a certified theorem without an additional
authoritative source or a new proof.  This note records the exact source gap so
that the decimal in the Anthropic paper is not accidentally promoted to theorem
evidence.

## What the Anthropic paper actually states

Remark 1.1 of *More Than Two Thirds of the Zeros of the Riemann Zeta
Function Lie on the Critical Line* (August 10, 2026) says that an "explicit
extremal law on configurations" gives a ceiling of `0.68185`.  The paper does
not state that extremal law, define its admissible configuration class, give a
primal or dual problem, exhibit a witness, prove an enclosure, or cite a source
for it.  The number occurs nowhere else in the paper.

The paper's proved cap, Proposition 7.4, is different.  It says that a
compression made from test functions supported in `[-L/2,L/2]`, with
`L = lambda log(T/2pi)`, has dimension asymptotic to `lambda N(T,2T)`.  It
therefore caps the number of on-line points certifiable by that compression at
`lambda N(T,2T)` and caps the Proposition 4.4(ii) certificate at
`(2-1/lambda)N(T,2T)+o(N)`.  At `lambda=1` this is a 100% dimension cap, not
the claimed 68.185% configuration ceiling.

The same section explains the support barrier: beyond `lambda=1`, the
off-diagonal prime sums are no longer dominated by the diagonal and require
prime-pair information of Hardy-Littlewood strength, equivalently pair
correlation information for `alpha>1`.  This part is authoritative, but it does
not determine the missing numerical ceiling.

Primary source:

- Claude, *More Than Two Thirds of the Zeros of the Riemann Zeta Function Lie
  on the Critical Line*, Remark 1.1 and Sections 7.1 and 7.5,
  <https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf>.

## What CCLM Corollary 14 proves

Carneiro--Chandee--Littmann--Milinovich use the Fourier convention

```text
hat R(t) = integral_R exp(-2 pi i x t) R(x) dx.
```

Their admissible one-delta class consists of nonnegative integrable `R` whose
Fourier transform is supported in `[-1,1]`, normalized by `R(0)>=1`.  For

```text
M(R) = integral_R R(x) [1-(sin(pi x)/(pi x))^2] dx,
```

Corollary 14 proves the exact lower bound

```text
M(R) >= 2^(-1/2) cot(2^(-1/2)) - 1/2,
```

with a unique explicit extremizer.  Through Montgomery's multiplicity
inequality this is precisely the Montgomery--Taylor result reconstructed in
M13.  It proves optimality of the bandlimited nonnegative one-delta subfamily;
it does not prove the broader `0.68185` ceiling.

Primary source:

- E. Carneiro, V. Chandee, F. Littmann, and M. B. Milinovich, *Hilbert spaces
  and the pair correlation of zeros of the Riemann zeta-function*, J. Reine
  Angew. Math. 725 (2017), Corollary 14,
  <https://arxiv.org/abs/1406.5462>.

## What the broader LP source proves

Chirre--Goncalves--de Laat relax Fourier support to a Cohn--Elkies sign
condition.  In a fixed-last-sign-change normalization their extremal constant
is

```text
c_LP = inf {
  f(0) + 2 integral_0^1 x f(x) dx :
  f is even Schwartz,
  hat f(0)=1,
  hat f >= 0,
  f(x) <= 0 for |x|>=1
}.
```

Equivalently, before fixing the last sign change, their class consists of even
continuous integrable functions with `f(0)=hat f(0)=1`, `hat f>=0`, and `f`
eventually nonpositive, and their objective is

```text
Z(f) = r(f) + (2/r(f)) integral_0^r(f) x f(x) dx.
```

They rigorously verify a degree-40 semidefinite candidate using interval
arithmetic, proving `c_LP <= 1.3208`.  Thus that source supplies a legal
certificate achieving a simple-zero lower bound of at least `0.6792`.  It does
not supply the reverse inequality needed for a global ceiling, an exact
optimizer, or a certified dual configuration witness near `1.31815`.

Primary source:

- A. Chirre, F. Goncalves, and D. de Laat, *Pair correlation estimates for the
  zeros of the zeta function via semidefinite programming*, Adv. Math. 361
  (2020), Sections 1.1, 3.1, and 4.1,
  <https://arxiv.org/abs/1810.08843>.

## Consequence for M14

The requested theorem needs a certified lower bound on `c_LP` (or on the exact
authoritative equivalent), because the simple-zero certificate is `2-c_LP`.
Neither cited extremal source supplies that lower bound.  Copying `0.68185`
from Remark 1.1 would hard-code the desired conclusion and violate M14's main
acceptance criterion.  Treating the CCLM optimizer as global would silently
narrow the class back to M13 and make the ceiling tautological.

The next missing input is therefore concrete: an authoritative statement and
proof of the claimed configuration extremal law, including its admissible
class, normalization, objective, and a rigorous dual witness or equivalent
global argument.  Until that exists, no `RepresentationCeiling` object should
be added and no M14 theorem should be emitted.
