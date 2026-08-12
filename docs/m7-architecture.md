# M7: certified Archimedean enclosure and finite Weil positivity

M7 answers its research question affirmatively. The existing `(f_2,f_3)`
matrix now has theorem-backed entry intervals, its real-symmetric `2 by 2`
principal minors are strictly positive by interval bounds, and positivity of the
Weil functional on `span{f_2,f_3}` is an exact downstream theorem. Universal
Weil positivity and RH remain unresolved: `function_space_restriction` is still
on the universal-to-finite path and no finite calculation can erase it.

## Exact integral and normalization

M7 keeps Lagarias's Mellin convention and real-place term unchanged. For
`h=f_q*tilde(f_r)`, put

```text
a=4 log(q), b=4 log(r), h(1)=h0=min(a,b),
U=(a+b)/2=2 log(qr), D=|a-b|/2=2|log(q/r)|,
H(u)=h0                 for 0 <= u <= D,
     =U-u               for D < u < U,
     =0                 for u >= U,
h(e^u)=e^(-u/2)H(u).
```

Lagarias's definition is

```text
W_inf(h)=(gamma+log(pi))*h(1)
 + integral_1^infinity
   [h(x)+tilde(h)(x)-2*x^-2*h(1)]*x*dx/(x^2-1).
```

After `u=log(x)`, the exact integrand is

```text
F(u)=2*(exp(3u/2)*H(u)-h0)/expm1(2u),  u >= 0.
```

The apparent singularity at zero is removable. On the plateau M7 uses
`R(u)=(3/4)*phi(3u/2)/phi(2u)`, where `phi(z)=expm1(z)/z`, and
`F(u)=2*h0*R(u)`. For a diagonal entry, where `H(u)=h0-u`
immediately to the right of zero, it uses
`F(u)=2*h0*R(u)-2*S(u)` with
`S(u)=exp(3u/2)/(2*phi(2u))`. These forms remove cancellation without changing
the imported formula. `F` is continuous, piecewise analytic, has derivative
kinks at `D` and `U`, and decays as `-2*h0*exp(-2u)` after support. No hidden
normalization or sign bug was found.

Source: Jeffrey Lagarias, *The Riemann Hypothesis: Arithmetic and Geometry*,
section 3, Theorem 3.1 and Remark 3.3:
<https://websites.umich.edu/~lagarias/doc/mt-holyoke-rev.pdf>.

## Certified finite integral and tail

The backend uses 192-bit `math/big.Float` endpoints. Every primitive result is
created with `ToNegativeInf` or `ToPositiveInf`; Go documents that a result
operand's precision and rounding mode control every arithmetic operation:
<https://pkg.go.dev/math/big>.

`exp`, `phi`, and `log` are enclosed by positive Taylor partial sums plus an
explicit geometric majorant for the remaining terms. Logarithms are power-of-2
range reduced so that the atanh-series argument is at most `1/3`. Euler's
constant and pi use exact decimal-prefix brackets from NIST DLMF 5.2.3 and
3.12.1: <https://dlmf.nist.gov/5.2.E3> and
<https://dlmf.nist.gov/3.12.E1>.

The finite domain is `[0,5]`. The certified enclosures of `D` and `U` are
inserted as two-sided guard cells, so no ordinary cell crosses a formula or
smoothness boundary. On each cell `I`, natural interval evaluation proves
`F(I) subset [m_I,M_I]`; hence `|I|*[m_I,M_I]` contains the exact cell
integral. Summing these enclosures is the elementary lower/upper Darboux-sum
theorem; it needs no sampled derivative estimate. The default uses 2,048 cells
per main smooth piece (fewer for the short plateau and post-support piece).
The theorem contract follows the upper/lower Riemann-sum construction in
Walter Rudin, *Principles of Mathematical Analysis*, Chapter 6. The Taylor
remainder bound is derived term-by-term from Taylor's theorem and the decreasing
geometric majorant once the successive-term ratio is below one.

The remaining infinite tail is not truncated. Since `H(u)=0` for `u>=U` and
`5>U` for every current pair,

```text
integral_5^infinity F(u)du = h0*log(1-exp(-10)).
```

Taylor interval evaluation encloses this closed form. A separate
`analytic_tail_bound` object records its start, endpoints, derivation,
exactness, and provenance. Removing it makes `EntryValue.Validate` reject the
otherwise rigorous quadrature record.

## Certified finite prime sums

Support is decided before transcendental evaluation. Integer comparisons decide
whether `p^n` is in the plateau, slope, support endpoint, or exterior. Thus
support exhaustiveness and numerical scalar certification are independent
facts. Every nonzero term
`2*log(p)*(p^n)^(-1/2)*H(n log p)` is evaluated with directed Taylor log/exp
intervals; the endpoint `p^n=(qr)^2` is exactly zero. The finite sum has
`outward_rounded_finite_sum` provenance. This was materially simpler than the
Archimedean enclosure. Exact rational endpoint contributions from M6 remain
exact.

## Results

Default certified intervals are:

| entry | prime power | Archimedean | total |
|---|---:|---:|---:|
| G11 | `[9.61161163657607, 9.61161163657607]` | `[8.27768573102409, 8.31532379043269]` | `[0.073064572991247, 0.110702632399846]` |
| G12=G21 | `[21.1614039044739, 21.1614039044739]` | `[10.8317881833572, 10.8680325371674]` | `[-0.0294364416412699, 0.0068079121688862]` |
| G22 | `[40.9031451576265, 40.9031451576265]` | `[15.8435338368135, 15.9525324298865]` | `[0.0332113013758238, 0.142209894448814]` |

The displayed prime endpoints coincide at 15 digits; deterministic JSON keeps
the distinct 192-bit outward decimal endpoints. All M6 approximate values lie
inside the corresponding intervals. `G21` reuses the canonical `G12`
enclosure by the already-certified Hermitian construction, not by numerical
overlap.

For `G=[[a,b],[b,d]]`, a real-symmetric matrix is positive semidefinite iff its
principal minors `a`, `d`, and `ad-b^2` are nonnegative. See the Hermitian
matrix entry in the Encyclopedia of Mathematics:
<https://encyclopediaofmath.org/wiki/Hermitian_matrix>. M7 proves strict bounds:

```text
lower(a) > 0,
lower(d) > 0,
ad-b^2 in [0.0015600654570083, 0.0159434103780462].
```

Therefore the finite matrix is positive definite (and hence PSD). Via the M5
coordinate identity, `Q_W(f)>=0` for every `f` in `span{f_2,f_3}` is certified.
This is not universal positivity and does not imply RH.

## Counterexample-first validation

Tests reject a tolerance-only proof kind, reject missing tail metadata, keep
prime support separate from scalar evidence, check log intervals against a
reference value, exercise outward arithmetic on exact rationals, verify all M6
approximations are enclosed, preserve basis-swap identity, expose every
breakpoint guard, and confirm a deliberately wide matrix stays open. Reducing
the backend to 24 bits and 16 cells per piece also leaves PSD open. Changing
the normal subdivision count still produces a valid enclosure.

## Architecture answers

1. The existing evaluator upgraded cleanly once its removable singularity was
   algebraically normalized; the approximate path remains unchanged.
2. The Archimedean cell-range enclosure dominates complexity. Constants,
   finite prime sums, and the tail are small in comparison.
3. Finite prime sums are easier to certify because integer support logic
   precedes finitely many positive transcendental evaluations.
4. The value lattice remains adequate. M7 adds only recognized local proof
   kinds and `certified_computation` claim evidence.
5. The smallest suitable home is an internal Go helper beside the evaluator;
   a generic interval library or Oct bridge would be needless scope.
6. Certified execution composes with theorem contracts: certified arithmetic
   creates an exact PSD premise, and existing finite-coordinate equivalences
   produce exact span positivity.
7. The `2 by 2` PSD proof needs only principal minors, not spectral machinery.
8. No hidden normalization bug appeared. Cancellation at `u=0` was an
   implementation hazard, not a formula error.
9. The matrix IR is numerically ready for a later zero-side structural
   decomposition, but M7 does not begin that work.
10. Compiler Theory: theorem-producing execution requires proof-carrying error
    control at every lowering boundary. Approximate agreement is a test;
    enclosure comes from range theorems, analytic remainders, and directed
    arithmetic; exact downstream theorems are allowed only after the evidence
    join contains no approximate component.

No M7 Oct experiment was needed. Octxiliary and `when utility` were not used.
The direct cancellation-safe Darboux strategy was clearly the smallest rigorous
route, so research scheduling did not require a utility branch.

## One next milestone

The smallest sensible next milestone is **M8: zero-side structural orbit
decomposition for the already-certified finite matrix IR**, keeping numerical
PSD and universal RH claims separate.
