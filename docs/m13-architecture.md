# M13: typed test-window optimization

M13 answers the central question affirmatively. Once M12's decomposition-aware
finite theorem exists, the compiler can optimize the remaining legal analytic
representation without changing that theorem or confusing numerical research
with certification.

## Test window

The authoritative source is Section 7.1, equations (7.1)-(7.4), and Theorem D
of the released [Anthropic technical paper](https://www-cdn.anthropic.com/564f962e60643842f5fcb4a17c9dbc8f608f1c37.pdf).
The normalized support parameter is

```text
0 < lambda <= 1,       L = lambda log(T/2pi).
```

Write `v(s)=phi(Ls)^2` on `[-1/2,1/2]`. The scale-free functional is

```text
c_lambda(v) = lambda (integral v)^2
              / (integral v^2
                 + lambda^2 double_integral |s-s'|v(s)v(s') ds ds').
```

The Cauchy-Schwarz extremal equation is `v+lambda^2 T v=constant`, where
`(Tv)(s)=integral |s-s'|v(s')ds'`. Thus `v''+2lambda^2 v=0`, and positivity
selects the Montgomery-Taylor profile

```text
v_lambda*(s) = cos(sqrt(2) lambda s).
```

At finite `T`, `phi` is the square root of this profile with the paper's fixed-
width `C3` endpoint ramp. It is even, nonnegative, compactly supported,
nonincreasing in absolute value, and `C2` after mollification. The ramp changes
the coefficients by `O(1/L)` only. The IR stores these conditions, the parameter
meaning and domain, support scale, transform convention, normalization, theorem,
and provenance structurally.

## Typed scale change

M13 introduces only the adapter required here:

```text
G_hat = G_tilde/(aL).
```

`ScaleChange` stores its source, target, arithmetic factor, parameter
dependencies, and provenance. It dispatches on the typed moment kind:

```text
tr(alpha G)     = alpha tr(G),
||alpha G||_F^2 = alpha^2 ||G||_F^2.
```

An exact fixture maps trace `9` to `3` and Frobenius square `9` to `1` under
factor `1/3`; the incorrect linear Frobenius rule gives `3` and is rejected.

The bounded scalar IR contains rationals, parameters, four arithmetic
operations, powers one and two, and the few unary analytic nodes used by this
family. Exact rational evaluation rejects transcendental nodes; those require a
separate certificate. This is not a general CAS.

## M12 coefficients and objective

For `theta=lambda/sqrt(2)` and `k=sqrt(2)lambda`, M13 stores

```text
a(lambda) = sin(theta)/theta,
b(lambda) = 1/2 + sin(k)/(2k),
J(lambda) = a(lambda)^2/2
            + (2a(lambda)cos(theta)-2b(lambda))/k^2.
```

Here `a=integral v`, `b=integral v^2`, and
`J=double_integral |s-s'|v(s)v(s')ds ds'`. The expression for `J` follows from
`(Tv)''=2v`, evenness, and the endpoint value of `Tv`. Each coefficient retains
its theorem, normalization, parameter identity, and endpoint-ramp error.

After the quadratic scale change, the normalized Frobenius coefficient is

```text
(b+lambda^2 J)/(lambda a^2) = 1/c_lambda(v).
```

The unchanged M12 `c=2` expression is `4 tr(G_hat)-2N-||G_hat||_F^2`.
Substituting typed coefficient expressions constructs

```text
J_M13(lambda) = 4-2-(b+lambda^2 J)/(lambda a^2).
```

It is not a stored output constant. Passing the flat profile `a=b=1,J=1/3` to
the same constructor gives `2-1/lambda-lambda/3`, hence exactly `2/3` at
`lambda=1`. This is the M12 regression and a normalization oracle.

For the extremal profile, algebra gives

```text
c_lambda* = sqrt(2)tan(theta)/(1+theta tan(theta)),
J_M13(lambda) = 2-lambda/2-(1/sqrt(2))cot(lambda/sqrt(2)).
```

## Optimization

The lower boundary is excluded and `J(lambda)` tends to negative infinity as
`lambda` tends to zero from above. The upper boundary is included. On the full
interior,

```text
J'(lambda) = cot(lambda/sqrt(2))^2/2 > 0,
```

since `0<lambda/sqrt(2)<=1/sqrt(2)<pi/2`. There is no interior stationary point:
`lambda*=1` is the unique global maximizer by strict monotonicity.

The exact value is

```text
J(1) = 3/2-(1/sqrt(2))cot(1/sqrt(2)).
```

Oct independently scanned 10,000 equally spaced legal values
`lambda=i/10000`. It found the included endpoint with value `0.6725007037`;
nearby values are `0.672363713439` at `0.9998` and `0.672432218366` at `0.9999`.
The retained plot zooms to `[0.5,1]` so the optimizing region is visible.

```powershell
cd C:/Users/yuech/source/repos/oct
go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m13_window_optimization.octest --execution interpreted
```

The run reported `7 passed, 0 failed, 0 skipped` in about four seconds. The
[plot](../experiments/m13_window_objective.png) and scan are research evidence,
not theorem evidence. OctGo and `when utility` were not used.

## Rigorous constant

Put `x=1/sqrt(2)`, so `x^2=1/2` exactly. Alternating rational Taylor bounds for
`cos(x)` and `sin(x)/x`, followed by positive interval division, prove

```text
18940672307/28164539016
  < J(1)
  < 492457480011/732278014418.
```

Thus `J(1)>269/400=0.6725`. The display `67.25%` is a safe lower rendering,
not an input or an assertion that the exact constant equals the decimal. The
all-distinct count is `(1+J(1))/2`, so it is greater than
`669/800=0.83625`, safely displayed as `83.625%`.

## Asymptotic consequence

M10's fringe/tail bounds and M11's existing `o`, `O`, and `Eventually` bridge
remain in force. For every epsilon, M12 is applied beyond the corresponding
finite threshold; only then is epsilon sent to zero. Therefore

```text
liminf N0_simple(T,2T)/N(T,2T)
  >= 3/2-(1/sqrt(2))cot(1/sqrt(2)) > 269/400,

liminf N0_distinct(T,2T)/N(T,2T)
  >= 3/2-(1/sqrt(2))cot(1/sqrt(2)) > 269/400,

liminf N_distinct(T,2T)/N(T,2T)
  >= 5/4-(1/(2sqrt(2)))cot(1/sqrt(2)) > 669/800.
```

This reproduces Theorem D exactly and is not claimed novel. RH remains
unresolved.

## Comparison and provenance

```text
M11: total first/second moments only                 -> 1/2
M12: add G=P+Q decomposition and rank/index IR       -> 2/3
M13: optimize the legal test-window representation  -> 67.25% display
```

The final theorem records M8 zero-orbit decomposition, M9 rank/index budgets,
M10 window/fringe control, M11 moment input, M12's unchanged rank-trace theorem,
and M13 admissible-window optimization.

The paper documents a broader bandwidth-one configuration ceiling near
`0.68185`. M13 neither implements that extremal ceiling nor searches beyond the
known Montgomery-Taylor family.

## Compiler Theory consequence

The gain is safe because support legality, `phi` versus `phi^2`, linear versus
quadratic scaling, objective construction, numerical search, analytic
certification, and display rounding are distinct semantic objects. The compiler
improves the result by changing an admissible representation, not by silently
strengthening M12.

The remaining awkwardness is that M11's complete asymptotic main terms still
have a prose shell around the new typed arithmetic core. M13 types only the
normalization-sensitive portion; a general symbolic asymptotic system would be
disproportionate.

## One next milestone

M14 should certify the bandwidth-one configuration ceiling near `0.68185` as a
typed impossibility theorem for this certificate class, without searching for
new kernels.
