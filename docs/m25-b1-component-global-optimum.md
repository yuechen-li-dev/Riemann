# M25: B1 control and the contact-preserving component optimum

## Outcome

M25 achieves **Success A**. It closes exactly the M24 residual window
`r in [3.503,3.512]` without changing the family, support, fixed contact, or
EP3.1. Composed with M23 and M24, it proves the theorem named
`ContactPreservingTwoAtomComponentGlobalOptimum`. This is not a theorem about
all two-atom families, and RH remains unresolved.

## Exact B1 definition and continuation

`B1` is the unique smooth graph `(m1(r),t1(r))` solving

```text
F0(t)+m H_r(t)=0,
F0'(t)+m partial_t H_r(t)=0,
```

with `t in [1.790,1.801]` on the full closed radius interval. Oct generated 90
radius predictors of width `10^-4`. Go reconstructed each graph cell using the
one-variable equation

```text
G(r,t)=F0'(t)H_r(t)-F0(t)partial_t H_r(t)=0,
```

proved `H<0` and `G_t<0`, bracketed one root, and tightened it with a
parameterized interval-Newton step. This proves existence and uniqueness for
every radius, not merely continuity of numerical roots. The persisted theorem
hulls are

```text
m1(r) in [0.0069,0.0071],
t1(r) in [1.790,1.801],
F_tt(m1(r),r,t1(r)) in [0.001,0.020].
```

The generated endpoint predictors place the left endpoint near
`m=0.0070469, t=1.80025` and the right endpoint near
`m=0.0069514, t=1.79111`. Exact theorem state uses the outward hulls above.

## Radius derivative and optimum

Nine grouped graph boxes prove strict envelope concavity across the residual
window. M23 supplies the already-certified unique stationary root. Therefore

```text
C1'(r)>0 for r<3.507850,
C1'(r)=0 for one r in [3.507850,3.507857],
C1'(r)<0 for r>3.507857.
```

M25 does not redo M23's local three-equation Krawczyk proof. It composes with
that proof to obtain the unique maximum on the window.

## Frequency compactification and competitor exclusion

The finite frequency interval is `0<=t<=4`. The proof uses:

| frequency band | certified classification |
|---|---|
| `[0,1.54]` | strict positive, including an analytic origin neighborhood |
| `[1.54,1.60]` | structural fixed quadratic contact at `pi/2` only |
| `[1.60,1.76]` | strict positive separator |
| `[1.76,1.82]` | unique movable B1 quadratic contact |
| `[1.82,4]` | strict positive; all sampled stationary competitors are too weak |
| `[4,infinity)` | uniform analytic atomic tail, lower bound at least `0.27` |

Centered Taylor bounds follow the certified branch derivative rather than an
uncorrelated rectangular `(m,r)` hull. This is essential near the shallow
noncontact minima. The strict persisted compact margin is `10^-9`; equality is
confined to the fixed contact and B1 contact neighborhoods.

Whole-line nonnegativity is the competitor theorem. If `H_r(t)<0`, then

```text
F0(t)+m1(r)H_r(t)>=0
```

implies `F0(t)/(-H_r(t))>=m1(r)`. Equality away from the fixed contact can occur
only at the unique B1 root. Thus no competing tangency can lower the mass cap,
and no branch switch occurs. The Oct adverse scan supplies structurally
different research evidence; generated Go is not called independent evidence.

## M23/M24/M25 composition

```text
M24: r<3.503 and r>3.512 are strictly below the M23 lower objective.
M25: B1 controls every r in [3.503,3.512].
M23: B1 has one stationary point there, a strict maximum.
```

The shared endpoints are exact, so there is no radius gap. The resulting
component optimum is uniquely contained in

```text
r in [3.507850,3.507857],
m in [0.006995080,0.006995100],
t in [1.795310,1.795317],
c in [1.1566431645986454,1.1566433083347464].
```

M25 also proves the interval-defined optimizer is globally feasible on the
whole line. Hence the updated whole-class bounds carried by this witness are

```text
c_* >= 1.1566431645986454,
J_* <= -0.1566431645986454.
```

The best simple explicit witness remains `r=7/2, m=69/10000`, with lower
objective endpoint `1.1564423523252728`. Its gap to the optimizer lower endpoint
is `0.0002008122733726`.

## Oct-first implementation, replay, and artifacts

Oct owns branch discovery, the 90-slab predictor atlas, adverse-frequency
search, plots, and deterministic Go-table generation. Five interpreted facts
pass in about 18.4 seconds; five compiled facts pass in about 0.65 seconds with
zero fallback. Artifact generation takes about 18 seconds.

Go owns the interval root/curvature/Taylor/tail kernels, theorem composition,
and persistent semantic types. `cmd/m25replay` performed the full exact rational
reconstruction in about 307 seconds. Ordinary M25 semantic tests validate the
pre-initialized deterministic manifest in about 1.6 seconds; known cells are not
recomputed in each unit test.

Artifacts:

- `experiments/m25_b1_component_closure.octest`
- `semantic/m25_continuation.generated.go`
- `artifacts/m25/m25_b1_continuation.png`
- `artifacts/m25/m25_b1_objective.png`
- `artifacts/m25/m25_competitor_gap.png`
- `cmd/m25replay/main.go`

## Interpretation, literature, and architecture

The active set is

```text
fixed double contacts:   ±pi/2
movable double contacts: ±t1(r)
unique radius optimum:   C1'(r*)=0.
```

This supports an equioscillation-like active-set picture but is not a general
alternation theorem. A focused primary/preprint screen found related
Fourier-extremal and pair-correlation frameworks in
[Chirre--Goncalves--de Laat](https://arxiv.org/abs/1810.08843),
[Chirre--Quesada-Herrera](https://arxiv.org/abs/2012.07781), and
[Das--Ismoilov--Ramos](https://arxiv.org/abs/2502.05106), but no stated match
for this two-atom contact-preserving radius/contact system. This is not a
novelty or priority claim.

Imported mathematics is M18's fixed contact, M22's exact family and tail,
M23's stationary optimizer, M24's complement barriers, and standard interval
Newton/implicit-function/Taylor arguments. Newly derived mathematics is the
full B1 graph certificate, branch-correlated compact proof, whole-line
competitor exclusion, and global component composition.

The architectural awkwardness is exact rational transcendental replay cost.
M25 addresses it by separating an immutable proof manifest and fast semantic
validation from the explicit full reconstruction command. For Compiler Theory,
the lesson is that a continuum active-set theorem can compile into graph cells,
contact neighborhoods, strict bands, and a tail, then persist as a small typed
artifact instead of being rediscovered by every test.

The evidence now points to preservation of the original `pi/2` contact as the
next limiting representation constraint.

## One proposed next milestone

M26: test whether releasing the forced `pi/2` contact improves the two-atom
optimum, while keeping exactly two atoms and the same EP3.1 route. M25 does not
begin that work.
