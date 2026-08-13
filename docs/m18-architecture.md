# M18: exact one-radius attainment by tangency

M18 reaches **Success A** for the unsaturated one-radius family:

\[
C_{\rm 1R}=1+\frac{\pi}{4}-\frac{2}{\pi}.
\]

This is a family theorem, not the full EP3.1 optimum and not RH.

## Sign convention

The nonnegative exterior measure is

\[
\sigma=c\,1_{\{|x|>1\}}dx+w(\delta_{-1}+\delta_1),
\]

and the M17/M18 Fourier density is

\[
F(t)=1+\frac{2(\cos t-1)}{t^2}+2(1-c)\frac{\sin t}{t}+2w\cos t.
\]

This is the convention consistent with the requested ceiling. The sign-reversed
displayed formula in the M18 request would give a different value at `pi/2`.

## Exact contact algebra

At `a=pi/2`, `cos(a)=0` and `sin(a)=1`. The value equation is

\[
F(a)=1-\frac8{\pi^2}+\frac{4(1-c)}\pi=0,
\]

so

\[
c=1+\frac\pi4-\frac2\pi.
\]

The derivative equation then gives

\[
w=\frac1\pi-\frac4{\pi^2}+\frac8{\pi^3}.
\]

The semantic layer stores both parameters as exact Laurent polynomials in
`pi`. Rational bounds `103993/33102 < pi < 104348/33215` are used only to
certify signs and displayed intervals.

## Equality-aware whole-line proof

The contact certificate proves exactly `F(a)=F'(a)=0` and

\[
F''(a)>\frac18.
\]

The integral representation bounds `|F'''|<11/20`. Every point in
`[3/2,17/10]` is within `13/100` of `a`, hence

\[
F''(t)>\frac18-\frac{11}{20}\frac{13}{100}
=\frac{107}{2000}>0.
\]

Strict convexity and the exact double root prove nonnegativity on the contact
region. The verifier labels this region
`NonnegativeCellWithCertifiedContact`; it never accepts an interval merely
because that interval contains zero.

The remaining compact regions `[0,3/2]` and `[17/10,4]` are covered by exact
degree-40 sine/cosine Taylor enclosures on a rational `1/2500` grid and a
certified `|F'|<83/100` cell loss. At the origin,

\[
a_0=2(1+w-c)>\frac1{25},
\]

so the candidate is safely unsaturated.

For `|t|>=4`,

\[
F(t)\ge1-2w-\frac4{t^2}-\frac{2(c-1)}{|t|}
\ge\frac34-2w-\frac{c-1}{2}>\frac13.
\]

Evenness supplies the negative half-line. All region endpoints overlap, so
the certificate covers every real frequency without gaps.

## Composition and bounds

The globally nonnegative exact witness supplies the lower direction and M17's
frequency envelope supplies the matching upper direction. The compiler derives,
rather than records as an axiom, the exact family equality.

Consequently,

\[
1+\frac\pi4-\frac2\pi\le c_*<\frac{1651}{1250},
\qquad
\frac{849}{1250}<J_*\le1-\frac\pi4+\frac2\pi.
\]

## Independent Oct evidence

`experiments/m18_exact_tangency.octest` passes 6/6 interpreted and 6/6
compiled with zero fallback. It reconstructs the parameters and derivatives,
scans the tangency at spacing `1e-6`, scans `[0,2000]` at spacing `1e-3`, and
shows that `c -> c+1e-6` makes the active contact negative. None of these
samples enters the theorem route.

`experiments/m18_tangency_plot.octest` deterministically produces
`artifacts/m18/m18_exact_tangency.png` in `u=t-pi/2` coordinates. Its SHA-256
is `f37dd4e48fd82b7c4b4a2fe3097decc7338b11b8616a7ff03b810b51c3e19b3d`.

## Replaying the proof

```text
go test -tags=slow ./compiler -run M18Slow -count=1
```

The normal suite checks algebra, legality, origin behavior, contact
multiplicity, equality-aware partition semantics, deterministic reports and
artifacts, and theorem composition. The slow test replays the full exact
compact covers through the production path.
