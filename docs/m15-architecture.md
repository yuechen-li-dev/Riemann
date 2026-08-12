# M15: a whole-line boundary-atom completion

M15 reaches **Success A**. It certifies

```text
9/8 <= c_* < 1651/1250,
849/1250 < J_* <= 7/8.
```

The CGdL endpoints are strict. The new `9/8` endpoint is exact; no grid or plot
participates in its proof. RH remains unresolved.

## Objective-sign audit

The authoritative EP3.1/M14A contract is

```text
nu = delta_0 + |alpha| 1_[-1,1](alpha) dalpha,
c(g) = [hat g(0) + integral_-1^1 |alpha| hat g(alpha) dalpha]/g(0).
```

The M15 request displays a minus sign once, but retains M14A's dual, CGdL
mapping, and scale near `1.318`. Those data require the plus sign. M15 continues
that explicit repository/EP3.1 contract and claims nothing for a minus variant.

## Candidate family and support

For `c>=1`, take

```text
sigma_c = c 1_{|x|>1} dx + (c-1)(delta_-1+delta_1).
```

Its density and atom masses are nonnegative. The density lies in `|x|>1` and
the atoms are on the legal boundary, so the measure is supported outside
`(-1,1)`. Cancellation gives

```text
P_c = delta_0 + (|x|-c)1_[-1,1] dx
      + (c-1)(delta_-1+delta_1).
```

With `t=2*pi*xi`,

```text
hat P_c(xi)
 = 1 + 2(cos(t)-1)/t^2 + 2(c-1)(cos(t)-sin(t)/t).
```

The origin expansion is `((9-8c)/12)t^2+O(t^4)`, so positivity forces
`c<=9/8` in this family.

## Whole-line endpoint proof

At `c=9/8`, set

```text
G(t) = 4t^2 hat P_(9/8)(t/(2*pi))
     = t^2(4+cos(t))-t sin(t)-8(1-cos(t)).
```

For `0<=|t|<=pi`, coefficient collection gives

```text
G(t) = sum_(n>=3) (-1)^(n-1)
       4(n-2)(n+1)t^(2n)/(2n)!.
```

The first coefficient is `1/45`. Consecutive absolute terms have ratio

```text
((n-1)(n+2))/((n-2)(n+1))
* t^2/((2n+2)(2n+1)) <= 25/56 < 1.
```

This uses `pi^2<10`, the first factor at most `5/2` (after clearing
denominators the difference is `(3n+2)(n-3)>=0`), and denominator at least
`56`. Thus the alternating series is nonnegative.

For `|t|>=pi`, evenness and elementary trigonometric bounds give
`G(t)>=3t^2-t-16`. Since `t>=pi>3`, this exceeds its value `8` at `t=3`, and
its derivative is positive thereafter. The two regions cover every real
frequency. Therefore `hat P>=0` globally, so `P` is positive definite by
Bochner--Schwartz. Weak duality yields `c_*>=9/8` and `J_*<=7/8`.

The compiler consumes the deterministic proof object
[`compiler/artifacts/m15_boundary_atom_witness.json`](../compiler/artifacts/m15_boundary_atom_witness.json)
and checks its exact rational obligations.

## Oct experiment

[`m15_boundary_atom_completion.octest`](../experiments/m15_boundary_atom_completion.octest)
independently reconstructs the density, derives `1.125`, scans through
frequency `100`, and finds a negative direction for `c=1.126`:

```text
interpreted: 4 passed, 0 failed, 2131 ms
compiled:    4 passed, 0 failed, 629 ms, 4 compiled, 0 fallback
```

The research-only plot is retained at
[`artifacts/m15/m15_boundary_atom_spectrum.png`](../artifacts/m15/m15_boundary_atom_spectrum.png).

## OCT TOOLING

During `[Artifact]`, convenience `PlotLine` called `plotrender.Render` directly.
It never called `ArtifactWriteCapability.StageArtifactOutput`, so the renderer
wrote an ambient path while the publisher had nothing to count.

`PlotLine` now treats its artifact-phase path as relative to `--output-root`,
rejects absolute/traversing paths, stages and atomically publishes the PNG, and
reports package, function, source, kind, execution, path, deterministic
identity, size, and SHA-256. Failure remains failure and does not publish.
Ordinary run/test plotting is unchanged. `--execution compiled` for artifacts
remains an explicitly reported build-time-interpreter delegation.

Before:

```text
Result: 1 artifact(s) passed, 0 failed
Outputs: 0 produced, 0 unchanged
```

After on the M15 plot:

```text
PRODUCED m15/m15_boundary_atom_spectrum.png
Outputs: 1 produced, 0 unchanged
Result: 1 artifact(s) passed, 0 failed
identity: M15BoundaryAtomPlot.PlotCertifiedEndpointSpectrum:m15/m15_boundary_atom_spectrum.png
sha256: 3efec504a953a03a8c5b3792033a233354212d513d1dd67d52d6c2453657ffe9
```

Regressions cover count/root, deterministic attribution, multiple noncolliding
tests, Fact exclusion, failure without publication, absolute-path rejection,
and both requested modes. The repair materially helped by retaining the
endpoint profile without file hunting.

## Provenance and status

Imported: EP3.1; CGdL's primal certificate; positive-definite distribution
theory; and the classical `c=1` completion.

Research-derived: the typed boundary-atom family; its exact `9/8` ceiling; the
endpoint identity; the alternating-series/coercive-tail proof; and the bracket.

A focused post-certification search checked the EP3.1 thesis and principal
Fourier-optimization papers and found neither `9/8` nor this completion. This
supports only **potential novelty**, not a priority claim.

Anthropic's `0.68185` remains inside `(0.6792,0.875]`, so it is neither
confirmed nor contradicted. Its class identity remains unsourced.

One proposed next milestone: **M16**, investigate a two-radius symmetric
exterior-atom completion with the same origin and whole-line proof discipline.

RH remains unresolved.
