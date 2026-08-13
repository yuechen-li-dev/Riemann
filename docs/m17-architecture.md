# M17: the unsaturated one-radius family

M17 reaches **Success B** with a stronger rational witness and a narrow ceiling:

```text
2297/2000 <= C_1R <= 1+pi/4-2/pi < 1149/1000.
2297/2000 <= c_* < 1651/1250.
849/1250 < J_* <= 1703/2000.
```

The imported primal endpoints remain strict. RH remains unresolved.

## Exact family and transform

The nonnegative exterior measure is

```text
sigma = c 1_{|x|>1} dx + w(delta_-r+delta_r),
c>0, r>=1, w>=0.
```

The isolated minus sign in the request conflicts with `w>=0`, the supplied
Taylor coefficients, M16, and the final conceptual block. M17 retains the
established plus-sign semantics. With `t=2*pi*xi`,

```text
G(t) = hat(P)(t/(2*pi))
     = 1+2(cos(t)-1)/t^2+2(1-c)sin(t)/t+2w cos(rt).
```

Positive definiteness is equivalent to `G(t)>=0` on the whole real line in the
existing normalization.

## Saturation is a subtype

`UnsaturatedOneRadiusCompletion` owns `c,r,w`, support and nonnegativity status,
the transform, Taylor coefficients, and theorem provenance.
`SaturatedOneRadiusCompletion` is constructed only after checking `a0=0`.

```text
SaturatedOneRadius < UnsaturatedOneRadius.
```

At the origin,

```text
a0 = 2(1+w-c),
a2 = c/3-1/4-w r^2,
a4 = -1/360+(1-c)/60+w r^4/12,
a6 = 1/20160+(c-1)/2520-w r^6/360.
```

On the saturated subtype, `w=c-1`, so

```text
a2 <= c/3-1/4-(c-1) = (9-8c)/12.
```

This recovers M15's exact `c<=9/8` local ceiling. When `a0>0`, negative `a2`
is legal. A useful non-global quartic screen is `a4>0` and
`4 a0 a4>=a2^2`, making `a0+a2 u+a4 u^2` nonnegative for `u=t^2>=0`.

## Dimension reduction and family ceiling

For fixed `(r,w)`, every frequency with `sin(t)>0` gives

```text
c <= [1+2(cos(t)-1)/t^2+2sin(t)/t+2w cos(rt)] t/(2sin(t)).
```

The origin adds `c<=1+w`; frequencies with `sin(t)<=0` are feasibility checks.
Also `w<=1/2` is necessary: along `t=(2k+1)pi/r`, the non-atom terms tend to
`1` and the atom equals `-2w`.

At `t=pi/(2r)` the atom vanishes, giving

```text
c <= h(t) = 1+t/(2sin(t))-tan(t/2)/t,
t=pi/(2r) in (0,pi/2].
```

Writing `s=t/2`,

```text
h(t)-1 = (s^2-sin^2(s))/(2s sin(s) cos(s)).
```

Its convergent even-power expansion has strictly positive coefficients on this
interval, hence `h` is strictly increasing. Therefore

```text
C_1R <= h(pi/2) = 1+pi/4-2/pi.
```

The bound is strict for every `r>1`; only `r=1` can approach it. The classical
bound `pi<355/113` and monotonicity of `1+x/4-2/x` give
`h(pi/2)<1149/1000`. This is a family ceiling, not an EP3.1 class ceiling.

## Rational witness and whole-line proof

The retained candidate is

```text
c=2297/2000, r=1, w=171/1000,
(a0,a2,a4,a6)=(9/200,-229/6000,3239/360000,-1847/5040000).
```

It improves the M16 one-radius degeneration `(573/500,1,21/125)`, which the
new verifier reproduces without constructing a two-radius object.

For `0<=|t|<=4`, the Go verifier uses exact degree-40 Taylor enclosures at
centers spaced by `1/5000`. The exact derivative bound is

```text
|G'| <= c-2/3+2wr = 4943/6000.
```

For `|t|>=4`, uniformly in `r`,

```text
G(t) >= 1-2w-4/t^2-2(c-1)/|t| >= 267/800 > 0.
```

The regions meet at `4`. The durable witness and ceiling object is
[`compiler/artifacts/m17_one_radius_witness.json`](../compiler/artifacts/m17_one_radius_witness.json).

The 20,001-center exact replay is intentionally isolated from normal unit and
race tests:

```text
go test -tags=slow ./compiler -run M17Slow -count=1
```

Normal tests retain exact parameter/coefficient checks, proof-object structure,
failure-path checks, deterministic reporting, artifact hashing, and bound
propagation. Production `CompileM17` also performs one whole-line replay;
report serialization does not repeat it.

## Oct evidence and plot

The reduced search estimated `c≈1.1487788, r=1, w≈0.1713`, with contact near
`pi/2`. Utility scoring favored `2297/2000`: its sampled minimum is about
`0.00035437` near `t=1.57195`, with a tractable exact proof.

[`m17_one_radius_completion.octest`](../experiments/m17_one_radius_completion.octest)
passed 6/6 interpreted and 6/6 compiled with zero fallback. It independently
reconstructs coefficients, regressions, contact geometry, the minimum, and tail
samples; it is non-certifying.

The plot
[`artifacts/m17/m17_one_radius_spectrum.png`](../artifacts/m17/m17_one_radius_spectrum.png)
was produced once and unchanged on rerun. It is 11,866 bytes with SHA-256
`1f0dd33ec965e8be79acc0faa0dc4fb853cd1273999908b23bdec46975ab4e82`.
It isolated the first-oscillation pocket and supported `t=4` as the handoff.

## Interpretation

Unsaturation gains `47/2000` over the saturated `9/8` ceiling. The entire
one-radius family still has `J>=1-pi/4+2/pi≈0.85122`, far from Anthropic's
reported `0.68185`; that value remains neither confirmed nor contradicted for
EP3.1.

A focused post-certification search of arXiv:1810.08843, arXiv:2310.01913, and
arXiv:2502.05106 found the surrounding SDP, Cohn--Elkies, EP, and pair-correlation
frameworks, but not this one-atom completion or constant. This is not a priority
claim.

Compiler Theory lesson: M15's barrier belonged to an active boundary condition,
not to the representation. Family identity and saturation must be distinct typed
states.

One proposed next milestone only: **M18**, certify or refute attainment of
`1+pi/4-2/pi` at `r=1` with the tangency-determined weight, without adding a
radius or broader machinery.
