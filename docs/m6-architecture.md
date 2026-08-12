# M6: typed Weil matrix entry evaluation

M6 answers its central question affirmatively: an exact Weil-form entry can be
lowered to actual values while exact definition, exact scalar evaluation,
rigorous enclosure, ordinary approximation, and no evaluation remain distinct.
The first matrix is useful computationally but does not certify PSD or RH.

## Basis and source

For integer `q >= 2`, M6 uses

```text
f_q(x) = x^(-1/2),  q^(-2) < x < q^2,
         0,          otherwise,
```

with half the one-sided value at both endpoints. The ordered basis is `(f_2,
f_3)`. It is two-dimensional, has a nontrivial cross entry, and is certified
admissible by a direct constraint check against Jeffrey Lagarias, *The Riemann
Hypothesis: Arithmetic and Geometry*, §3: nice functions are piecewise `C2`,
compactly supported in `(0,infinity)`, and take midpoint values at
discontinuities. The same source gives Theorem 3.1 (the explicit formula),
Theorem 3.2 (Weil positivity), Remark 3.3 (the polarized intersection product),
and the convention

```text
M[f](s) = integral_0^infinity f(x) x^s dx/x.
```

Source: <https://websites.umich.edu/~lagarias/doc/mt-holyoke-rev.pdf>, pages
5–7 of the PDF. No computationally friendlier replacement formula is imported.

The basis is motivated by four facts. Its Mellin endpoints are exact rationals,
its convolution is a trapezoid in log coordinates, compact support makes every
prime-power sum finite, and unequal widths make the off-diagonal nonzero. A
large-scale basis search was neither needed nor performed.

## Closed form used by the evaluator

Write `a=4 log(q)`, `b=4 log(r)`, and

```text
h_qr(e^u) = e^(-u/2) H_qr(u),
H_qr(u) = min(a,b)                         if |u| <= |a-b|/2,
          (a+b)/2 - |u|                    if |a-b|/2 < |u| < (a+b)/2,
          0                                otherwise.
```

This is `f_q * tilde(f_r)`. The endpoint contribution is retained exactly:

```text
M[f_q](0) = M[f_q](1) = 2(q-q^-1),
M[h_qr](0) + M[h_qr](1)
  = 8(q^2-1)(r^2-1)/(qr).
```

The prime contribution is evaluated as the exact finite definition

```text
2 sum_(p^n <= (qr)^2) log(p) (p^n)^(-1/2) H_qr(n log p).
```

The support proves that the omitted remainder is exactly zero. The returned
decimal is nevertheless approximate because logarithms, square roots, and
summation use binary64. Metadata records the bound, exact summand, deterministic
Eratosthenes enumeration, term count, 53-bit precision, and remainder status.

The real-place term preserves Lagarias's definition

```text
W_inf(h) = (EulerGamma+log(pi))h(1)
  + integral_1^infinity
      [h(x)+tilde(h)(x)-2 x^-2 h(1)] x dx/(x^2-1).
```

The evaluator changes variables to `u=log(x)`, splits adaptive Simpson
quadrature at both trapezoid kinks, and evaluates the post-support infinite tail
analytically as `h(1) log(1-exp(-2U))`, `U=(a+b)/2`. The recorded `1e-11`
tolerance is heuristic and is explicitly not a proof bound.

## Value evidence semantics

`EntryValue` has four disjoint states:

- `unevaluated_exact_definition`: the semantic target is exact but has no
  scalar;
- `approximate_value`: a numerical scalar with backend, precision, and
  non-rigorous error semantics;
- `certified_interval`: an enclosure requiring a recognized proof object;
- `exact_value`: a symbolic/exact scalar requiring an independent exact
  argument.

Legal upgrades are unevaluated to approximate, certified interval, or exact;
and certified interval to exact when the independent exact argument is present.
Approximate to certified or exact is illegal regardless of precision. Mixed
arithmetic takes the weakest evidence: exact plus interval is interval, while
any approximate component makes the total approximate. A tested toy interval
path exercises certification support; no actual Weil total is mislabeled as
certified.

The matrix keeps `structurally_defined_exact` and
`HermitianByConstruction=true` even though every total value is approximate.
Numerical symmetry discrepancies are diagnostics with `theorem_use=false`.

## First evaluated matrix

At 53-bit precision and `1e-11` heuristic quadrature tolerance, the separately
evaluated components are:

| entry | endpoint (exact) | prime (approx.) | archimedean (approx.) | total (approx.) |
|---|---:|---:|---:|---:|
| G11 | 18 | 9.61161163658 | 8.29650346216 | 0.0918849012616 |
| G12 | 32 | 21.1614039045 | 10.8499117334 | -0.0113156378917 |
| G21 | 32 | 21.1614039045 | 10.8499117334 | -0.0113156378917 |
| G22 | 512/9 | 40.9031451576 | 15.8980166165 | 0.0877271147825 |

Thus

```text
G ~= [[ 0.0918849012616, -0.0113156378917],
      [-0.0113156378917,  0.0877271147825]].
```

The maximum numerical Hermitian discrepancy is zero for this real centered
basis. Two complex-coefficient direct-functional probes first assemble the
combined convolution and then evaluate it; their discrepancies from `c*Gc`
are about `6.51e-11` and `8.02e-11`, within the declared `1e-10` experimental
tolerance. They are evidence about implementation consistency, not theorem
premises.

The optional two-by-two diagnostic gives approximate eigenvalues
`0.0783009892043` and `0.101311026840`, with condition number about
`1.29386649989`. `certifies_psd` remains false. No inertia, rank, or orbit
theorem machinery was added. The zero-side aggregate remains a theorem-linked,
unevaluated alternate representation; no truncated-zero cross-check was done.

## Oct experiment

- Path: `experiments/m6_weil_logbox_matrix.octest`
- Command:

  ```powershell
  C:\Users\yuech\source\repos\oct\cmd\oct\oct.exe test `
    C:\Users\yuech\source\repos\Riemann\experiments\m6_weil_logbox_matrix.octest `
    --execution compiled
  ```

- Basis: `(f_2,f_3)` above.
- Configuration: Oct `Float`, explicit prime-power list through 81, independent
  200,000-panel composite Simpson quadrature, and a 20,000-versus-200,000 panel
  refinement probe.
- Precision: Oct `Float` (binary64 in this local backend).
- Truncation: `(qr)^2`, support-exhaustive for each pair; quadrature domain
  handled in log coordinates with the same analytic tail.
- Output: five passed, zero failed, five compiled, zero interpreted fallback;
  compiler identity `oct dev`; wall time about 3.3 seconds.
- Limits: one bounded five-fact file, no external data, no hosted execution ID,
  and no interval or portability guarantee.
- Octxiliary: not used. Direct local invocation was the smallest integration;
  no Oct payload is accepted as compiler evidence.
- Why experimental: binary64 assertions and quadrature refinement provide no
  theorem-backed error enclosure. Agreement with Go is an independent
  implementation check, not proof.

`when utility` was not used. Once the source convention was inspected, the
centered log-box basis dominated the alternatives on admissibility,
tractability, finite prime support, exact endpoints, and off-diagonal value.

## Bugs, awkwardness, and compiler-theory finding

The experiments found no sign error, but they forced two important corrections
to the initial design: finite support proves a zero prime remainder without
making the floating-point sum exact, and direct quadratic validation must
evaluate the combined convolution rather than merely re-sum stored entries.

The main architectural awkwardness remains proposition-family plumbing: claim
validation, semantic keys, cloning, theorem patterns, and JSON rendering need
manual extension. M6 avoided inventing a numerical DSL, but typed value records
now make that duplication more visible. Evaluator configuration is also not yet
a reusable theorem-recognized error-bound object.

M6 teaches a compiler-theory lesson: mathematical execution is an
evidence-refinement problem, not just expression evaluation. Exact denotation,
algorithm applicability, scalar production, error enclosure, and proof use are
separate judgments. Provenance must flow through arithmetic joins, and the join
must be monotone toward weaker evidence; otherwise exact subexpressions can
launder an approximate result.

The evaluated matrix IR is ready to be consumed by later decomposition work in
the representational sense: entries, contributions, transforms, basis pairs,
and diagnostics are typed. It is not ready to support certified spectral claims
because the real-place quadrature remains heuristic.

## One next milestone

The smallest sensible next milestone is **M7: certified archimedean enclosure
for the existing `(f_2,f_3)` entries**. It should add a theorem-backed interval
bound for the split log-coordinate integral and propagate it through the already
finite prime sums, without beginning zero-orbit inertia or rank analysis.
