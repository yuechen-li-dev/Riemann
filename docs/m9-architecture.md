# M9: aggregate spectral accounting

M9 is the first deduction milestone. It consumes, rather than rebuilds, M8's
identity

```text
G = P + Q
```

and combines facts from its two theorem-linked representations. The zero side
supplies the local structure of `P` and `Q`; the explicit-formula side supplies
an accessible exact spectral fact about the same `G`. The permanent
`function_space_restriction` remains attached to the finite compression.

## Deliberately small invariant vocabulary

M9 introduces only `rank`, `positive_index`, and `negative_index`. An invariant
claim records its matrix, structural dimension, relation and natural-number
bound, evidence class, theorem lineage, and source. Certified principal-minor
and exact-theorem evidence can satisfy exact theorem premises. Approximate
eigensolver output can be reported as an experiment but is rejected as theorem
evidence. M9 does not implement a general eigensolver or spectral decomposition.

For a finite Hermitian matrix, `n_plus` and `n_minus` mean the numbers of
strictly positive and strictly negative eigenvalues, equivalently the maximal
dimensions of positive- and negative-definite subspaces. M9 does not need to
store `n_zero`: the standard identity with the dimension remains mathematical
background rather than unused IR.

## Imported finite theorems

The contract set is intentionally limited to:

1. `rank(sum A_j) <= sum rank(A_j)`;
2. `n_plus(A+B) <= n_plus(A)+n_plus(B)`;
3. the negative-index analogue, obtained by applying (2) to `-A,-B`;
4. `n_plus(A) <= rank(A)`.

Rank subadditivity is sourced to Horn and Johnson, *Matrix Analysis*, 2nd ed.,
§0.4.5. Positive-index subadditivity is sourced to the finite pull-back proof
in Lemma 3.1 of the released Anthropic paper. All contracts require finite
Hermitian matrices, and sum contracts require a shared dimension. No numerical
eigenvalue premise is accepted.

## Zero-side resource accounting

Every critical fixed-location block is a positive multiple of `u*u*`. Thus it
is PSD and has rank at most one. Rank subadditivity gives

```text
rank(P) <= sum_k rank(P_k) <= C_nz <= C_mult.
```

`C_nz` counts nonzero critical evaluation vectors; `C_mult` counts critical
zeros with multiplicity. These are deliberately not identified. If a location
has multiplicity `m`, its matrix is `m*u*u*`: multiplicity changes its positive
weight but creates no new vector direction. Repeated or collinear locations can
collapse the aggregate rank further.

Each M8 off-critical reflection-pair block is Hermitian with

```text
rank <= 2, n_plus <= 1, n_minus <= 1.
```

Positive-/negative-index and rank subadditivity therefore give

```text
rank(Q)    <= sum_j rank(Q_j)    <= 2 B_pairs,
n_plus(Q)  <= sum_j n_plus(Q_j)  <= B_pairs,
n_minus(Q) <= sum_j n_minus(Q_j) <= B_pairs.
```

These are budgets, never additive equalities. Zero vectors, dependent pairs,
overlapping subspaces, multiplicity, and cancellation can only reduce the
actual aggregate invariants.

## NEWLY DERIVED MATHEMATICAL RESULT

**Finite critical-rank theorem.** Let `G=P+Q` be finite Hermitian matrices of a
common dimension. Suppose an exact/certified fact gives `n_plus(G)=g`, the
off-critical aggregate satisfies `n_plus(Q)<=B_off`, and `P` is the critical
aggregate. Then

```text
n_plus(G) <= n_plus(P) + n_plus(Q) <= rank(P) + B_off,
```

so

```text
rank(P) >= max(0, g-B_off).
```

Since `rank(P)<=C_nz<=C_mult`, the sound count consequence is only

```text
C_mult >= C_nz >= rank(P) >= max(0, g-B_off).
```

The proof is the direct combination of positive-index subadditivity and
`n_plus(P)<=rank(P)`, followed by exact arithmetic in the natural numbers. The
component inequalities are standard. M9 newly derives their representation-
fused consequence in this compiler; it does not claim a new linear-algebra
theorem in the literature.

For the current certified M7 matrix, strict positive principal minors give
`n_plus(G)=2` without computing approximate eigenvalues. The instance is

```text
rank(P) >= max(0, 2-B_off).
```

M8 supplies no finite `B_off` for the full zero sum, so this sanity instance
does not exclude an off-critical block and does not imply RH.

## Counterexample-first research

Exact Go fixtures and the bounded Oct experiment reject three tempting
strengthenings:

- inertia is not additive: `diag(1,-1)+diag(-1,1)=0`;
- positive-definite `G` need not have `Q=0`:
  `I_2=diag(0,2)+diag(1,-1)`;
- multiplicity is not rank: `m*u*u*` has rank one for nonzero `u`, for every
  positive `m`.

The Go adversarial loop exhausts 125 exact small combinations made from one
critical outer product and one reflection-pair block. It checks
`n_plus(G)<=rank(P)+n_plus(Q)` and the rearranged lower bound using determinant
and trace signs for explicit 2-by-2 matrices. This is a tiny-matrix fixture, not
a general symbolic eigenvalue engine.

The Oct experiment is
`experiments/m9_spectral_budget.octest`, run with:

```powershell
go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m9_spectral_budget.octest --execution interpreted
```

from the Oct repository. Its eight deterministic cases passed. It covers an
independent pair, a dependent pair, cancellation, direction collapse,
multiplicity, the finite-PSD reverse counterexample, and both nonvacuous and
vacuous critical-rank budgets. Its evidence class is experimental
counterexample search only. `when utility` was not used because the
positive-index route was mathematically direct after the exact counterexamples.

## Comparison with Claude's first one-half stage

The primary paper was consulted only after the compiler-side finite inequality
had emerged. Proposition 4.5 proves the more count-specific finite inequality

```text
s_1 >= 2*n_plus^theta(G) - N_window,
```

after adding a height window and far-zero perturbation control. Appendix C.5
states that the earlier one-half result then used a thresholded
Cauchy-Schwarz lower bound from the first two trace moments. M9 does not
reproduce the percentage because the current IR lacks exactly:

1. a height-window total-zero count tied to the compression;
2. far-zero operator-norm control for the threshold `theta`;
3. an imported asymptotic first/second-moment theorem yielding
   `n_plus^theta(G) >= (3/4-o(1))*N_window`.

The stronger rank-trace Lemma 3.2, pair-correlation support optimization, and
the `2/3` / `67.25%` result were not encoded or used.

## Architectural finding

M9 exposes a useful compiler boundary: local orbit semantics naturally lower
to *upper budgets*, while a critical-count result requires an independently
observable *lower bound* on the positive index of the fused matrix. The two
representations are complementary, not interchangeable. Retaining their
separate provenance makes the new deduction inspectable and prevents the
certified finite `2x2` result from masquerading as global zero information.

RH remains unresolved.

