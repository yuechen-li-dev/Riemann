# M24: global control of the contact-preserving two-atom component

## Outcome

M24 ends in **Success C**. It finds no stronger contact branch and does not
claim that the M23 local optimizer controls its entire legal component. It does
prove that every legal feasible family point outside

```text
r in [3.503, 3.512]
```

has objective strictly below the certified M23 lower endpoint
`1.1566431645986454`. Thus this width-`0.009` radius interval is the exact
remaining obstruction. Compiled Oct continuation places the candidate active
contact in approximately `[1.7905,1.8020]` there, but that contact range is
research evidence rather than theorem state.

The certified family uses

```text
c = c0 + (pi/2) m cos(pi r/2).
```

The minus sign in the prompt's later radius-envelope display conflicts with
the stated exact family, M22/M23 code, and the supplied optimizer value. M24
retains the certified plus-sign semantics.

## Legal component and boundaries

`LegalTwoAtomRadiusComponent` is defined as the connected component containing
the M23 radius of

```text
{r>1:
  m*(r)=inf_{H_r(t)<0} F0(t)/(-H_r(t)) is finite and nonnegative,
  w1(m*(r),r)>=0,
  c(m*(r),r)>0}.
```

Support legality is exactly the two radii `1` and `r>1`, with weights `w1` and
`w2=m*`; no third atom, density, or support widening occurs. `F0` and `H_r` are
entire after assigning the removable sinc value at zero. The fixed equations
`F(pi/2)=F_t(pi/2)=0` hold structurally. M23 certifies strict mass, weight,
objective, fixed-curvature, and movable-curvature legality on its optimizer
box.

M24 does **not** prove that this component is all of `(1,infinity)`. Its exact
endpoints, possible weight boundary, possible `m*=0` boundary, and any
legality-induced disconnection remain unresolved. The Oct scan found legal
sampled points throughout `1.01<=r<=22.99`, but no sampled extent is promoted
to a component boundary. This endpoint uncertainty cannot hide a better
objective outside `[3.503,3.512]`, because the following bounds apply to every
legal feasible point without assuming component topology.

## Global radius reduction

For `1<r<=3`, `cos(pi r/2)<=0`; since `m>=0`,

```text
c <= c0 < 1.1566431645986454.
```

At the structurally preserved fixed contact, direct differentiation gives the
new M24 identity

```text
F_tt(pi/2)
  = F0_tt(pi/2)
    - 2 m cos(pi r/2) (r^2 - 1 + 8/pi^2).
```

Whole-line feasibility makes the left side nonnegative. On positive-cosine
bands this bounds the objective gain by

```text
c-c0
 <= (pi/4) F0_tt(pi/2)/(r^2-1+8/pi^2).
```

The rational interval replay proves this is below the M23 lower objective for
all `r>=181/50`. Nonpositive-cosine bands already obey `c<=c0`. This is a true
all-radius tail theorem, not a finite scan.

For the remaining compact range, Oct emits 611 fixed frequency witnesses for
the exact cells

```text
[3+k/1000, 3+(k+1)/1000],
k=0..502 and k=512..619.
```

If an emitted frequency has `H_r(t)<0`, feasibility gives
`m<=F0(t)/(-H_r(t))`. A handwritten purpose-specific rational interval replay
then bounds

```text
c0 + (pi/2) cos(pi r/2) F0(t)/(-H_r(t)).
```

It uses the envelope derivative

```text
d/dr [c0+(pi/2)F0(t)A(r)/(-H_r(t))]
 = (pi/2)F0(t)(A H_r - A' H)/H^2
```

in a centered Taylor enclosure. This preserves the cancellation that a naive
interval quotient loses near stationarity. Every generated cell is strictly
below the M23 lower target. Because these are adverse-frequency bounds, they
eliminate every possible active branch on the covered radii without first
proving which branch is active.

No component-uniform frequency compactification theorem was obtained. It is
not needed for the certified complement: compact radii use explicit adverse
frequencies and the entire radius tail uses the fixed contact.

## Branch atlas

The compiled Oct source scans all local minima of `M_r(t)` on `0<t<=40` for a
broad radius mesh, refines the lowest and next-lowest candidates, and audits
integer and half-integer radii. It recovers:

| radius range | branch | contact range | envelope behavior | proof status |
|---|---|---|---|---|
| `[3,3.62]` | `B1-movable-low-frequency` | about `2.424` down to `1.685` | sampled rise to M23, then fall | M23 box certified; M24 complement eliminated independently |
| representative broad radii | `Bremote-oscillatory` | multiple higher-frequency pockets | never sampled active near M23 | eliminated outside obstruction; root count inside unresolved |

At the M23 radius the next sampled branch has mass more than `0.01` above B1.
M23's earlier remote pockets near `9.44` and `15.54` remain far from active.
No branch-switch point was sampled near M23. None of these numerical facts is
used as a certificate.

M23 already proves `C'>0` immediately left and `C'<0` immediately right of its
unique stationary root. M24 uses the same purpose-specific derivative only to
bound fixed-frequency proof cells. It does not claim certified B1 monotonicity
throughout `[3.503,3.512]`. The missing theorem is precisely:

> On `r in [3.503,3.512]`, certify the B1 contact continuation and prove that
> every other negative-response frequency gives a no-smaller mass cap; then
> combine B1 with the M23 stationarity signs.

## Optimizer, witness, and bounds

The M23 optimizer box is unchanged:

```text
m in [0.0069950800000000,0.0069951000000000]
r in [3.5078500000000000,3.5078570000000000]
t in [1.7953100000000000,1.7953170000000000]
c in [1.1566431645986454,1.1566433083347464].
```

M24 does not globally certify this optimizer's whole-line feasibility. The best
global witness therefore remains

```text
r=7/2, m=69/10000,
c_* >= 1+pi/4-2/pi+69*pi*sqrt(2)/40000,
J_* <= 1-pi/4+2/pi-69*pi*sqrt(2)/40000.
```

Its certified `c` lower endpoint is `1.1564423523252728`. The gap from the M23
local family candidate is about `0.00020088`. M24 changes neither whole-class
bound nor EP3.1 and RH remains unresolved.

## Active-set interpretation

```text
family freedoms:       m, r
structurally preserved F(pi/2)=0, F_t(pi/2)=0
movable contact:       F(t1)=0, F_t(t1)=0
radius stationarity:   C'(r)=0
```

This explains why a whole-line positivity problem can have a finite local
optimizer system: structural equalities and active contacts consume the free
parameters, while strict compact regions and an analytic tail remain
inequalities. The evidence is consistent with an equioscillation-like rule in
which each useful parameter is consumed by another contact or stationarity
condition; M24 does not state a general theorem.

M24 is only about the M18-contact-preserving surface. Whether a stronger
two-atom witness requires releasing the original `pi/2` contact is left as a
representation question, not solved here.

## Oct-first implementation and artifacts

Oct owns branch scanning, special-radius audits, competitor discovery, contact
witness selection, plots, and deterministic Go-table generation. The primary
compiled suite passes `5/5` with zero fallback in about `1.19s`. Artifact
publication is build-time interpreted by Oct's documented design; it publishes
three plots plus the generated table. The relocated publication took about
`87.0s`; the final deterministic replay took `91.7s` and reported all four
outputs unchanged.

Artifacts:

- `experiments/m24_two_atom_component.octest`;
- `semantic/m24_radius_barriers.generated.go`;
- `artifacts/m24/m24_global_radius_envelope.png`;
- `artifacts/m24/m24_active_contact_atlas.png`;
- `artifacts/m24/m24_competitor_mass_gap.png`.

The generated table SHA-256 is
`68592d9c19e79e237c1ce5f7f3f74aa45ec81b3011fe778604005f05abbbf4af`.
The plot hashes are recorded in compiler tests. The generated table is not
independent verification: independence comes from broad Float continuation
versus rational interval adverse-frequency replay.

Handwritten Go is limited to the stable semantic types, the fixed-contact
curvature theorem, a scalar rational interval response/derivative evaluator,
centered Taylor replay, validation, and reports. M24 adds no handwritten Go
branch finder. The family formula necessarily appears in the Oct research
specification and the pre-existing Go certification kernel; no new duplicate
optimizer logic remains.

The architectural rough edge is that Oct compiled-data publication targets
standalone `package main` readers. M24 therefore uses readable
`Artifact.WriteText` metaprogramming for the destination-specific semantic
table. Also, the rational interval kernel does not cache immutable trig point
enclosures, so replay is slower than the theorem's size suggests.

## Literature, provenance, and Compiler Theory

The requested post-result literature screen was conditional on Success A or B.
Success C neither solves the component nor finds a stronger branch, so no new
screen was performed. M23's primary/preprint screen remains the current
context; no novelty or priority claim is made.

Imported mathematics consists of M18's fixed double contact, M22's exact
family and whole-line semantics, M23's Krawczyk/IFT/local-maximum theorem and
explicit witness, Taylor's theorem, and the necessary nonnegative curvature at
a smooth global minimum. Newly derived mathematics is the fixed-contact
curvature identity, the all-radius `181/50` barrier, the centered
fixed-frequency envelope bound, and exclusion of every legal branch outside
the exact unresolved window.

For Compiler Theory, M24 shows that a global continuum question can compile
into three layers: structural invariants remove infinite tails, generated
adverse witnesses remove compact cells without identifying active branches,
and the residual box becomes a typed statement of the exact missing theorem.
An honest residual box is useful compiler output, not a numerical claim of
globality.

## One proposed next milestone

Certify the B1 contact branch and exclude all competing tangencies on
`r in [3.503,3.512]`, then combine that theorem with M23 stationarity, while
keeping the same contact-preserving two-atom family.
