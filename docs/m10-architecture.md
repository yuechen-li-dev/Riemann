# M10: height-window compression and thresholded inertia counting

## Result

M10 answers its finite question affirmatively. For the height-dependent Weil
compression in Claude's August 10, 2026 paper, a certified input

```text
n_plus^theta(G_tilde_T) >= L_theta
```

with `theta >= theta0` compiles to

```text
N0_simple(T,2T) >= 2*L_theta - N(T,2T) - 2*N(I'\I),
N_distinct(T,2T) >= L_theta - N(I'\I).
```

Here `I=(T,2T]` and `I'=(T-sqrt(T),2T+sqrt(T)]`. This is a
structural reconstruction of Proposition 4.5, not a novel result. It is finite
and does not use the later rank-trace or moment machinery.

## Authoritative source and proof factorization

The implementation follows Claude, *More Than Two Thirds of the Zeta Zeros on
the Critical Line* (August 10, 2026), Sections 2.2 and 4.1--4.3, especially
Propositions 4.1, 4.2 and 4.5:

https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf

The source proof factors into five compiler concerns:

1. hard target and enlarged localization windows;
2. the zero-side near/far split;
3. an operator-norm tail bound;
4. strict threshold transfer by Weyl;
5. critical/off-critical inertia accounting and conservative count conversion.

The prime-side moments that manufacture `L_theta` remain an explicit input.

## Height, boundary and orbit semantics

`HeightWindow` carries boundary and ordinate conventions as data. M10 uses
`lower < gamma <= upper` and positive ordinates only, matching the paper's
`N(T,2T)`. Thus a zero at `T` is excluded, one at `2T` is included, and its
negative-ordinate conjugate is not silently added. The target center is
`3T/2`, with half-width `T/2`; the localization half-width is
`T/2+sqrt(T)`.

The smooth `C3` taper belongs to the compression family and is distinct from
hard zero membership. `HeightCoordinate` supports symbolic theorem parameters
and exact-integer test fixtures; symbolic windows reject executable membership
queries rather than guessing.

Reflection partners `rho` and `1-conjugate(rho)` share an ordinate. They enter
or leave together, remain two geometric locations, and own one unordered pair
ID. M10 keeps separate:

```text
N(W)             total multiplicity
N0(W)            critical multiplicity
N0Distinct(W)    distinct critical locations
N0Simple(W)      simple critical locations
NOffPairs(W)     unordered off-critical reflection pairs
rank(P_near)     evaluation-vector directions
```

Multiplicity weights a block; it does not clone a direction.

## Window compression and splits

The new `WindowCompression` is separate from M7's `(f2,f3)` matrix:

```text
L = lambda*log(T/2pi)
tau_k = T + 2pi*k/L
f_{T,k}(u) = phi_T(u) exp(-i*tau_k*u)
d = floor(L*T/(2*pi))
G_tilde_T = G_T/L
```

It records zero-side and explicit-formula representations of the same matrix.
M7 remains an unchanged transitive regression case, not the asymptotic family.

The refined split is

```text
G_tilde_T = A_tilde_T + E_tilde_T
          = P_near + Q_near + E_tilde_T.
```

`A_tilde_T` contains zeros with ordinate in `I'`; `E_tilde_T` contains all
others. `P_near` is critical and near; `Q_near` is off-critical and near. Far
and off-critical are independent axes.

## Far-zero control

Proposition 4.2 imports

```text
||E_tilde_T||_op <= theta0
theta0 = 4*A0*C1^2*X^(1/2)*log(4T)/D0^2
D0=sqrt(T), X=exp(L), C1=||phi''||_1.
```

For fixed `0<lambda<=1` and taper profile,

```text
theta0 = O(log(T)*T^(lambda/2-1)) = o(1).
```

The proof uses the taper's off-real-axis `r^-2` decay and Titchmarsh's local
zero count `N(t+1)-N(t)<=A0 log(t+3)`. The typed contract retains assumptions,
uniformity, dependencies and provenance. A numerical norm estimate cannot
satisfy it.

## Strict threshold and Weyl transfer

M10 defines

```text
n_plus^theta(R) = #{i : lambda_i(R) > theta}.
```

It remains distinct from ordinary `n_plus`. The threshold is structured by
`T`, `lambda`, taper width, `A0`, and normalization. Weyl/Courant--Fischer gives

```text
theta >= ||E||_op  =>  n_plus^theta(A+E) <= n_plus(A).
```

Strictness is essential: `A=[0]`, `E=[1]`, `theta=1` defeats a non-strict
count. Likewise `theta<||E||` is unsafe. Exact rational diagonal fixtures test
equality without tolerances; approximate eigensolver results are report-only.

## M9 reuse and proof

M10 adapts M9's resource accounting:

```text
n_plus^theta(G_tilde)
  <= n_plus(A_tilde)
  <= rank(P_near) + n_plus(Q_near)
  <= s1 + s2 + p.
```

The M8 pull-back blocks supply the last line, with M9's off-critical budget
`B_off=s2+p`. Since

```text
N(I') >= s1 + 2*s2 + 2*p,
```

we obtain

```text
s1 >= 2*n_plus^theta(G_tilde) - N(I').
```

Also `n_plus^theta(G_tilde)<=rank(A_tilde)<=#Z(I')`. Removing `I'\I`
gives the target bounds. No evaluation-vector independence is assumed; the
base conversion stays `rank(P_near)<=N0Distinct<=N0`.

## Oct filtering and scope

`experiments/m10_threshold_window.octest` ran as:

```text
go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m10_threshold_window.octest
```

`StrictlyAbove(value,theta)` is implemented through the refined Concept
`StrictlyPositiveGap`, whose checked constructor admits exactly positive
`value-theta` gaps. Equality and negative gaps follow the fallible error arm,
so the experiment's strict convention is expressed at a typed Concept boundary
rather than duplicated as a raw conditional. All eight bounded fixtures passed.
They cover equality, sub-bound thresholds, near/far cancellation, small norm
with high rank, low rank with large norm, and a dimension-averaged Frobenius
substitute. The rerun compiled all eight facts with zero interpreted fallbacks.
This is counterexample filtering only, never theorem evidence. The installed
CLI did not expose the skill-documented `oct check` command; `oct test` supplied
the supported parse, type-check, compile, and execution path.

`when utility` was not used: Proposition 4.5 fixes the operator norm and strict
Weyl threshold, leaving no genuine research scheduling choice.

M10 does not derive a half-type asymptotic result. The remaining obligation is
a certified or trusted asymptotic lower bound for
`n_plus^theta(G_tilde_T)` at a threshold compatible with `theta0`. Trace,
Frobenius, rank-trace, support optimization, `2/3`, and `67.25%` are not begun.
RH remains unresolved.
