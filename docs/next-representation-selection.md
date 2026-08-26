# M20: representation selection beyond the exact one-radius optimum

M20 reaches **Success A** as a representation-selection milestone. It selects one next family and makes no new theorem claim:

```text
selected family: contact-adapted r=4 atomic ray
certified bound:  1+pi/4-2/pi <= c_* < 1651/1250   (unchanged)
next milestone:   M21, optimize and certify or refute this ray
```

All statements below are labeled **Structural**, **Numerical**, **Certified**, or **Speculative**. Numerical scans and utility scores never enter theorem propagation.

## 1. Exact M18 baseline

**Certified.** The reference completion is

```text
sigma_0 = c0 1_{|x|>1} dx + w0(delta_-1+delta_1),
c0 = 1+pi/4-2/pi,
w0 = 1/pi-4/pi^2+8/pi^3,
t0 = pi/2.
```

For

```text
F0(t) = 1 + 2(cos(t)-1)/t^2
          + 2(1-c0)sin(t)/t
          + 2w0 cos(t),
```

M18 proves `F0(t0)=F0'(t0)=0`, `F0''(t0)>1/8`, and whole-line nonnegativity. Numerically `c0=1.1487783910298668`, `w0=0.17103742708003555`, and `F0''(t0)=0.1289135368477682`.

## 2. First variation and active contact equations

**Structural.** Let a legal even exterior perturbation have Fourier response `H(t)`, and allow

```text
c = c0 + epsilon gamma,
w = w0 + epsilon beta.
```

The linearized transform is

```text
delta F(t) = -2 gamma sin(t)/t + 2 beta cos(t) + H(t).
```

At `a=pi/2`, preservation of the contact value and derivative gives

```text
delta F(a)=0   => gamma = a H(a)/2,
delta F'(a)=0  => beta  = gamma/a^2 + H'(a)/2.
```

Thus a perturbation can raise `c` to first order exactly when it supplies `H(pi/2)>0`. The contact value fixes the possible increase, and the derivative consumes one retuning degree. After those equations are satisfied, the remaining local margin is curvature. Without derivative retuning, the predicted contact displacement is

```text
delta t = -delta F'(a)/F0''(a).
```

Finite-difference tests agree with these formulas. The origin is not active: its M18 margin is about `0.0445`.

## 3. Pressure map and secondary weak frequencies

**Numerical.** The exact contact is the only fragile baseline region found. The next local minima are:

| frequency | `F0(t)` (approx.) | role |
|---:|---:|---|
| `pi/2` | `0` | exact certified double contact |
| `9.31` | `0.61` | first secondary minimum |
| `15.65` | `0.64` | second secondary minimum |
| `21.95` | `0.65` | third secondary minimum |

The oscillatory minima tend toward `1-2w0`, approximately `0.6579`; M18 already has a conservative certified tail margin above `1/3`. The next representation should therefore repair the active contact neighborhood rather than target a distant pocket.

The Oct artifacts are:

- [`m20_baseline_pressure_map.png`](../artifacts/m20/m20_baseline_pressure_map.png), SHA-256 `5014c88583683c84b8e41f859dea87ba8d3306b7af59c561f0066aa9e7b1099a`;
- [`m20_selected_perturbation_effect.png`](../artifacts/m20/m20_selected_perturbation_effect.png), SHA-256 `04e35e80ecea7cf0c6d745f8cdbf0f5bc0e9b0e4b92a9756c80ff376153b5a71`.

They are research evidence, not certificates.

## 4. Candidate comparison

The numerical probes use `0<=t<=80` with step `1/50`. Reported minima are off-contact mesh minima near `t=1.58`; the exact retained contact remains zero.

| family | new scalar DOF | best coarse `c` | proof route | utility | decision |
|---|---:|---:|---|---:|---|
| A: generic two active radii | 4 | `1.15474741707169` | exact contact cell, compact Taylor cover, atomic tail | 44 | reduce to D |
| B: boundary atom + uniform interval density | 3 | `1.15474825257476` | sinc/Taylor core, decaying density tail | 45 | retain as fallback |
| C: three-site lattice packet proxy | 2 | `1.15087278613226` | trigonometric compact/tail proof | 31 | reject |
| D: contact-adapted `r=4` atomic ray | 1 | `1.15474741707169` | exact Laurent-in-pi contact, inherited compact/tail strategy | 52 | **select** |

### Candidate A — two genuinely active radii

**Structural.** The M18 solution lies on the `w2=0` boundary. For a new atom of positive-side mass `m` at radius `r`,

```text
H(a)=2m cos(ra),
H'(a)=-2mr sin(ra).
```

Any phase with `cos(pi r/2)>0` permits `gamma>0`. The bounded genuine-two-radius scan selected `r=4`; both masses are positive.

**Numerical.** At `r=4`, `m=0.0038`, the scan gives `c=1.15474741707169` and an off-contact mesh minimum about `3.96e-7` near `1.58`.

**Decision.** The generic family is promising but carries unused radius and weight directions. Its useful result is the smaller Candidate D.

### Candidate B — boundary atom plus one density mode

**Structural.** A uniform positive-side mass `m` on `[15/4,17/4]`, mirrored across zero, has

```text
H(t)=2m sinc(t/4) cos(4t).
```

It is nonnegative, supported outside `(-1,1)`, positive at the contact, and decays as `O(1/t)`.

**Numerical.** At half-width `1/4` and `m=0.0039`, the coarse probe gives `c=1.15474825257476`, only about `8.4e-7` above A/D at this resolution.

**Decision.** This is the strongest rejected family and the natural fallback if D fails. Its tiny sampled edge does not yet pay for exact sinc derivative algebra and a new theorem scalar shape.

### Candidate C — periodic/lattice completion

**Structural.** A finite equal-weight packet at radii `3,4,5` is a legal proxy. Only the `r=4` site contributes positively to the contact value; the other two masses primarily add oscillation. A true infinite lattice would transform to a distribution rather than the pointwise Fourier density expected by the current verifier.

**Numerical.** With total positive-side mass `0.004`, the best probe is `c=1.15087278613226`.

**Decision.** Reject. The finite proxy is inefficient, while the infinite form would require a disproportionate semantic and proof expansion.

### Candidate D — contact-adapted perturbation

**Structural.** Choosing the smallest simple phase-aligned integer radius gives the exact ray

```text
r2 = 4,
w2 = m,
c  = c0 + (pi/2)m,
w1 = w0 + (2/pi)m.
```

Because `cos(2pi)=1` and `sin(2pi)=0`, these identities preserve

```text
F(pi/2)=0,
F'(pi/2)=0
```

for every `m`, not merely to first order. Its contact curvature is affine:

```text
F''(pi/2;m) = F0''(pi/2) + (2-16/pi^2-32)m.
```

The curvature forecast reaches zero near

```text
m = 0.00407681510458557,
c = 1.15518223722117.
```

**Speculative.** This is the likely family endpoint and may create a quartic contact. M20 does not certify that endpoint or exclude another active frequency before it.

**Decision.** Select. D introduces one scalar, retains exact Laurent-in-pi algebra, matches B's improvement scale, and has the most direct whole-line proof route.

## 5. `when utility`

Oct's `when utility` chose D deterministically. The weighted factors were expected improvement, analytic tractability, inverse whole-line proof cost, exact-witness potential, full-class reuse, inverse implementation cost, and conceptual clarity. Scores were D `52`, B `45`, A `44`, C `31`.

**Structural boundary.** Utility is prioritization metadata. It neither changes evidence classification nor derives a certified extremal bound.

## 6. Oct execution and regression protection

[`m20_representation_selection.octest`](../experiments/m20_representation_selection.octest) passed:

```text
interpreted: 5/5, 23.094 s
compiled:    5/5, zero fallback
artifacts:   2/2 produced through build-time-interpreted [Artifact]
compiler:    Oct dev, gooct-cli
```

The experiment checks the M18 zero-perturbation regression, bounded atomic/density/lattice probes, and the explicit utility choice. Go tests separately check exact M18 identity at zero perturbation, first variation against finite differences, contact-shift sensitivity, measure legality, the `1651/1250` primal guard, evidence separation, utility separation, deterministic reports, and artifact hashes.

## 7. Certification assessment

- **D:** analytic contact core plus the inherited equality-aware Taylor subdivision away from contact; elementary atomic tail. Lowest proof cost.
- **A:** the same route after fixing radii, but unnecessary free parameters complicate exact endpoint discovery.
- **B:** exact sinc/Taylor compact proof and a favorable decaying tail. Moderate proof cost.
- **C finite:** trigonometric polynomial plus compact cover. Moderate proof cost but low gain. **C infinite:** no current pointwise proof route.

No family was accepted on dense scanning alone.

## 8. Theorem boundary and external context

**Certified.** No accidental improvement occurred. M0-M18 semantics remain fixed and the live bracket remains

```text
1+pi/4-2/pi <= c_* < 1651/1250,
849/1250 < J_* <= 1-pi/4+2/pi.
```

Every numerical candidate is below the compatible CGdL primal guard `1.3208`; candidates at or above it are rejected by typed validation. The selected numerical point corresponds to `J` near `0.84525`, still far above Anthropic's contextual `0.68185`. That remark did not influence the search objective or utility score.

The selected atom is a legal full-EP3.1 exterior measure, and the first-variation formula applies to arbitrary legal even exterior responses. This makes D reusable research information without claiming that the two-atom family is full-class optimal.

## 9. Provenance, architecture, and Compiler Theory

Imported mathematics: the exact M18 extremal and certificate; EP3.1 support-one/sign-cutoff semantics; the CGdL primal bound; elementary Fourier transforms of symmetric atoms and uniform intervals.

Newly derived mathematics: the active-contact first variation; derivative retuning and contact-shift equations; the exact `r=4` ray; its curvature slope and predicted flattening point; positive directions for three representation types; and the secondary-frequency pressure map.

Architectural awkwardness: M18's Laurent-in-pi scalar is sufficient for D, but B would require exact sinc expressions. M20 therefore keeps floating probes behind an explicit `Numerical` evidence type instead of widening theorem algebra prematurely.

Compiler Theory lesson: representation selection is constraint compilation. Inspect where the certificate is tight, synthesize a transform response linearly independent at that constraint, and use existing slack to satisfy derivative conditions. More parameters are useful only when their Fourier images buy new control at active constraints.

## 10. One proposed next milestone

**M21 only:** optimize and certify or refute the one-parameter contact-adapted `r=4` atomic ray, including the predicted curvature-flattening endpoint and a whole-line equality-aware certificate. Do not widen the family unless this ray fails.
