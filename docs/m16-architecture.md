# M16: two-radius completion beyond 9/8

M16 reaches **Success A**. It certifies

```text
573/500 <= c_* < 1651/1250,
849/1250 < J_* <= 427/500.
```

The CGdL endpoints remain strict. RH remains unresolved.

## Exact family and normalization

With `nu=delta_0+|x|1_[-1,1]dx`, weak duality asks for
`P=nu-c dx+sigma` positive definite, where `sigma` is nonnegative and supported
outside `(-1,1)`. Cancelling `-c dx` outside the data interval forces

```text
sigma = c 1_{|x|>1}dx
      + w1(delta_-r1+delta_r1)
      + w2(delta_-r2+delta_r2),
c>0, 1<=r1<=r2, w1,w2>=0.
```

The atom weights have no further dual normalization. `w1+w2=c-1` is M15's
origin-saturation choice, not a general requirement. Radius exchange is
canonicalized; coincident radii collapse by adding weights. For `t=2*pi*xi`,

```text
hat(P)(xi) = 1 + 2(cos(t)-1)/t^2 + 2(1-c)sin(t)/t
             + 2w1 cos(r1 t) + 2w2 cos(r2 t).
```

## Exact origin analysis

Writing `hat(P)=a0+a2 t^2+a4 t^4+a6 t^6+...` gives

```text
a0 = 2(1+w1+w2-c),
a2 = c/3-1/4-w1 r1^2-w2 r2^2,
a4 = -1/360+(1-c)/60+(w1 r1^4+w2 r2^4)/12,
a6 = 1/20160+(c-1)/2520-(w1 r1^6+w2 r2^6)/360.
```

Local nonnegativity requires `a0>=0`; only when `a0=0` does `a2` decide the
local sign. On that saturated branch, `w1+w2=c-1` and `rj>=1` imply

```text
a2 <= c/3-1/4-(c-1) = (9-8c)/12.
```

Thus every origin-saturated nonnegative two-radius completion has `c<=9/8`.
In fact, `a2=0` together with `a0>=0` gives
`c-1<=w1+w2<=sum(wj rj^2)=c/3-1/4`, so quadratic cancellation above `9/8` is
impossible even with slack. This does not eliminate positive-slack candidates
with a negative quadratic term, because their constant term controls a
neighborhood of the origin.

## Candidate and whole-line proof

Oct's locally filtered `(r1,r2)=(1,2)` scan found a boundary near `c=1.147`.
The retained exact candidate is

```text
c=573/500, r1=1, r2=2, w1=21/125, w2=1/1000.
a0=23/500, a2=-1/25, a4=911/90000, a6=-451/840000.
```

Positive `a0` makes the negative quadratic term legal. The Go verifier covers
`0<=|t|<=4` using degree-40 Taylor enclosures at the exact centers `k/1000`,
`0<=k<=4000`. Taylor's theorem bounds each trigonometric remainder. The exact
derivative bound

```text
|hat(P)'(t)| <= c-2/3+2(w1 r1+w2 r2)=1229/1500
```

turns the center bounds into a cover of the compact interval. For `|t|>=4`,

```text
hat(P) >= 1-2(w1+w2)-4/t^2-2(c-1)/|t| >= 339/1000.
```

The regions meet at `4`; no middle region is needed. The proof object is
[`compiler/artifacts/m16_two_radius_witness.json`](../compiler/artifacts/m16_two_radius_witness.json).
The same exact compact/tail verifier also accepts the degeneration `w2=0`
(with tail anchor `341/1000`). Thus the certified gain does not require the
second atom; it requires leaving the saturated slice that defined M15's family.

## Oct evidence and plot

[`m16_two_radius_completion.octest`](../experiments/m16_two_radius_completion.octest)
reconstructs the transform and coefficients, checks degenerations and illegal
parameters, scans through `t=100`, and finds that `c=1.148` has a negative
pocket near `t=1.54`. Interpreted and compiled modes passed 6/6; compiled mode
used six compiled cases and zero fallback. These probes are not theorem evidence.

The plot
[`artifacts/m16/m16_two_radius_spectrum.png`](../artifacts/m16/m16_two_radius_spectrum.png)
was produced on its first run and reported unchanged on its deterministic
rerun, with identity
`M16TwoRadiusPlot.PlotCandidateAndRegionHandoff:m16/m16_two_radius_spectrum.png`,
11,782 bytes, and SHA-256
`942e1e969e550d391b6fe44e62e319eb7233216efc8299fb7edd508fef7ffc46`.
It isolated the shallow minimum near `1.54` and supported `t=4` as the proof
handoff. It is not part of the proof.

## Provenance, comparison, and lesson

Imported mathematics: EP3.1 weak duality and normalization; CGdL's strict
primal certificate; Bochner--Schwartz; and Taylor's theorem. Newly derived:
the unsaturated two-radius semantics, exact local coefficients, the saturated
`9/8` ceiling, and the `573/500` whole-line certificate.

A focused check of [CGdL](https://arxiv.org/abs/1810.08843),
[Carneiro--Milinovich--Ramos](https://arxiv.org/abs/2310.01913), and
[Das--Ismoilov--Ramos](https://arxiv.org/abs/2502.05106) found the relevant
Fourier extremal formulations and primal scale, but no matching two-atom
completion or `573/500` value. This is not exhaustive and is no priority claim.

Anthropic's `0.68185` stays inside `(0.6792,0.854]`; it is neither confirmed nor
contradicted, and its unnamed class is still not sourced as EP3.1.

Compiler Theory answer: the broader representation expresses certified
information beyond M15, but **no**, the additional radius itself is not
essential—the certified `w2=0` degeneration proves that. The useful distinction
is saturation: origin saturation kills every nonnegative exterior-atom family
above `9/8`, while positive origin slack and a strictly positive asymptotic
trigonometric polynomial remove that obstruction.

One proposed next milestone: **M17**, characterize and optimize the unsaturated
one-radius boundary-atom family analytically, correcting the representation
boundary exposed by M16 before adding further radii.
