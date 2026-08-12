# M8: zero-side orbit decomposition

M8 answers its central question affirmatively, with one qualification that the
typed lowering makes unavoidable. Critical and off-critical zeros do reappear
as different finite matrix structures. The clean local unit is a **critical
reflection pair**. A full Klein-four geometric orbit contains two such units in
general, so a property of one pair must not automatically be attributed to the
full quartet.

## Unchanged functional and trust boundary

M8 uses the M4 Lagarias/Weil summand verbatim:

```text
M[f](rho) * conjugate(M[f](1-conjugate(rho)))
```

The transform remains
`M[f](s)=integral f(x)x^s dx/x`. The zero sum still counts multiplicities and
uses the symmetric limiting order imported with the explicit formula. The
source is J. C. Lagarias, *The Riemann Hypothesis: Arithmetic and Geometry*,
§3, Theorems 3.1 and 3.2:
<https://websites.umich.edu/~lagarias/doc/mt-holyoke-rev.pdf>.

The M7 interval matrix is not modified. `DualMatrixRepresentation` links its
certified explicit-formula realization to the new symbolic zero-side aggregate
using `weil-lagarias.explicit-formula`. Its
`numerical_identification` field is false: interval values did not establish
the identity.

## Typed orbit and evaluation semantics

`ZeroOrbit` records a representative, critical/off-critical classification,
the four generating transforms, a canonical quotient of distinct points, zero
multiplicity, symmetry theorems, and provenance. A generic off-critical orbit
has four locations. A critical-line orbit has two; the transform collisions are
retained in `generated_by`. `DistinctLocationCount` and `ZeroMultiplicity` are
separate fields.

`BasisEvaluationVector` owns the ordered basis, point, Mellin convention,
dimension, exact symbolic entries, evidence status, theorem lineage, and
provenance. Its semantic key includes both basis and point. For the current
basis its entries are `M[f2](rho)` and `M[f3](rho)`; no zero table is needed.

## Coordinate derivation and orientation

Put `v_i=M[f_i](rho)`, `w_i=M[f_i](1-conjugate(rho))`, and
`f=sum_i c_i f_i`. Mellin linearity gives

```text
M[f](rho) = c^T v
M[f](1-conjugate(rho)) = c^T w.
```

Therefore

```text
(c^T v) conjugate(c^T w)
  = sum_ij conjugate(c_i) conjugate(w_i) v_j c_j
  = c* K(rho) c,

K(rho)_ij = conjugate(w_i) v_j.
```

Thus `K=conjugate(w) transpose(v)`. It is not `v w*`. A focused complex test
compares the direct summand with `c* K c` and catches a transpose or misplaced
conjugation. This lowering is generic over a typed basis and point.

The Mellin linearity/conjugation source is NIST DLMF §1.14(iv):
<https://dlmf.nist.gov/1.14.iv>. The local outer-product, rank, and two-vector
facts are recorded against Horn and Johnson, *Matrix Analysis*, 2nd ed., §§4.2
and 7.1. M8 imports no global inertia theorem.

## Critical-line structure

If `rho=1/2+i gamma`, critical reflection fixes `rho`, hence `w=v` and

```text
K(rho) = conjugate(v) transpose(v) = u u*,  u=conjugate(v).
```

Each fixed-point block is Hermitian, PSD, and rank at most one. Rank is exactly
one only under the explicit premise `v != 0`. The conjugate location is a
second distinct fixed point. Consequently a full critical geometric orbit is a
sum of two PSD rank-at-most-one blocks: it is PSD and rank at most two for the
current basis, not automatically rank one.

This was the main place where orbit shape and grouping prevented an attractive
but false overstatement.

## Off-critical structure

For `p` off the line, let `r(p)=1-conjugate(p)`. The point contributions obey
`K(r(p))=K(p)*`, so the reflection-pair block is exactly

```text
H = K(p)+K(r(p)) = a b* + b a*,
```

for the corresponding conjugated evaluation vectors `a,b`. This is Hermitian
and has rank at most two. On `span(a,b)`, its determinant is the negative Gram
determinant. It has signature `(1,1)` if `a,b` are linearly independent. If
`a=lambda b`, it becomes `2 Re(lambda) b b*`: it can be positive semidefinite,
negative semidefinite, or zero. M8 therefore stores conditional indefiniteness
and the dependence degeneracy separately; it never claims unconditional
signature.

A generic quartet partitions into two reflection pairs, exchanged by
conjugation. Each pair has the local classification above. The sum of the two
pair blocks is Hermitian, but M8 does not claim that the full quartet sum is
indefinite. Pair grouping is mathematically cleaner than forced quartet-block
terminology.

## Aggregate and soundness boundary

`ZeroSideMatrixAggregate` is an unevaluated semantic sum over nontrivial-zero
orbits. It retains the domain, multiplicity rule, symmetric limiting convention,
basis, transform, theorem lineage, and critical/off-critical contribution
templates. The templates are schemas, not assertions that particular zeros
exist. Classification partitions the index domain naturally:

```text
P = sum over critical orbits G_O
Q = sum over off-critical orbits G_O
G = P + Q.
```

No infinite summation and no aggregate inertia calculation is implemented.
The compiler also performs a deliberate attempted inference from M7 finite PSD
to zero-set critical-reflection structure. It is rejected because no theorem
path exists. Even a PSD total need not have PSD summands; compensation among
orbit contributions is possible.

## Oct experiment

- Path: `experiments/m8_zero_orbit_blocks.octest`
- Command:

  ```powershell
  C:\Users\yuech\source\repos\oct\cmd\oct\oct.exe test C:\Users\yuech\source\repos\Riemann\experiments\m8_zero_orbit_blocks.octest
  ```

- Execution: 3 compiled, 0 interpreted fallback; 3 passed, 0 failed.
- Inputs: synthetic vectors only; no actual zeta zeros.
- Critical output: `[[5,2.5i],[-2.5i,1.25]]`, determinant `0`, and the tested
  quadratic value `6.25`.
- Independent off-critical output: `[[2,1+i],[1-i,0]]`, determinant `-2`, hence
  one positive and one negative eigenvalue (the Go diagnostic reports
  `1±sqrt(3)`).
- Dependent output: `4*conjugate(v)transpose(v)` for `v=(1,i)`, determinant `0`.
- Evidence: numerical/synthetic experiment only. Exact compiler classification
  comes independently from the finite-algebra contracts.
- Octxiliary: not used; the experiment needed no transport or sidecar.
- `when utility`: not used. Direct coefficient expansion uniquely fixed the
  orientation, and critical-reflection pairing immediately supplied the
  adjoint, so there was no genuine scheduling choice to score.

The attempted trailing `--json` CLI placement was rejected by this local Oct
build after the tests ran; rerunning the documented plain command succeeded.
That CLI issue did not affect the mathematical experiment.

## Findings

1. Critical/off-critical geometry emerges from existing M3 symmetry, M4's
   summand, and M5 coordinate conventions; it is not a zeta-specific matrix
   constructor.
2. The orientation is phase-sensitive. Visual outer-product intuition would
   have selected the wrong row/column order.
3. Critical rank one belongs to each reflection-fixed location, not generally
   to its two-location geometric orbit.
4. Off-critical `(1,1)` signature is exact for an independent reflection pair,
   conditional on evaluation-vector independence. The full quartet needs no
   such local signature claim.
5. The architectural awkwardness is that M5's `HermitianMatrix` requires a
   globally Hermitian matrix, while an individual Weil point contribution need
   not be Hermitian. M8 therefore introduces a small symbolic contribution IR
   below the existing matrix without weakening M5 validation.
6. The matrix IR is now structurally ready to *receive* later global
   rank/inertia reasoning, but no engine for sums or inertia exists yet.
7. Compiler Theory gains a concrete lesson: symmetry is not decorative
   metadata. Compiling the group action into canonical quotient and coset
   structure determines which non-Hermitian terms may legally become a
   Hermitian block, and controls exactly which spectral facts are sound.
8. RH remains unresolved, and M7 finite positivity cannot exclude off-critical
   orbit contributions.

## One next milestone

The smallest sensible next milestone is **M9: theorem-backed rank bounds for
finite sums of the M8 orbit blocks**, without yet adding general inertia or
percentage-on-line theorems. M8 reveals rank-controlled local constituents;
sum-rank bookkeeping is the narrow missing bridge before broader inertia
questions are meaningful.
