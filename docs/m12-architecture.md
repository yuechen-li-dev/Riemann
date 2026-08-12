# M12: rank-trace structure beyond the first-two-moment ceiling

M12 answers the central question affirmatively.  Total first and second moments
alone give M11's sharp generic positive-index theorem and the `1/2`
simple-critical stage.  Retaining the component identity

```text
G = P + Q
```

adds information that those moments do not contain.  The compiler now consumes
the same total matrix together with `P >= 0`, `rank(P) <= r`, and
`n_plus(Q) <= b`.

## Generic finite theorem

For finite same-size Hermitian matrices `P,Q`, assume `P` is positive
semidefinite, `rank(P)<=r`, and `n_plus(Q)<=b`.  For every real `c>0`,

```text
||P+Q||_F^2
  >= c tr(P) - (c^2/4) r + 2c tr(Q) - c^2 b.
```

Here `||G||_F^2=tr(G^2)`.  At `c=2`,

```text
r >= 2 tr(P) + 4 tr(Q) - 4b - ||P+Q||_F^2
  = 4 tr(G) - 2 tr(P) - 4b - ||G||_F^2.
```

This is Lemma 3.2 of the released Anthropic paper.  M12 reconstructs its proof:
write `Q=Q_plus-Q_minus`, discard the nonnegative `tr(P Q_plus)` term, use von
Neumann's trace inequality to align the eigenvalues of `P` and `Q_minus`, and
apply `x^2>=cx-c^2/4` and `q^2>=2cq-c^2`.  The result is compiler-derived from
those imported standard matrix facts; merely passing the Oct experiment is not
theorem evidence.

Equality holds for `P=(c/2) Pi_1`, `Q=c Pi_2`, where the orthogonal projections
have ranks `r,b` and orthogonal ranges.  Thus both penalty coefficients are
sharp.  The theorem permits `tr(Q)` of either sign.

## Premise attacks

Exact fixtures preserve the principal failure modes:

- Without PSD, `c=2`, `P=[-10]`, `Q=[12]`, `r=b=1` gives Frobenius square `4`
  and false right side `23`.
- `P=[1],Q=[0]` is equality and falsifies any smaller rank penalty.
- `P=[0],Q=[2]` is equality and falsifies any smaller positive-index penalty.
- Using negative rather than positive index fails on `Q=[2]`.
- `P=[1],Q=[-2]` confirms that negative `tr(Q)` is legal.
- The dependent block `a=b=(1,0)` gives `ab*+ba*=diag(2,0)` and remains safe:
  it has positive index one without an independence premise.
- The same `G=I_2` has decompositions `(P=I_2,Q=0)` and `(P=0,Q=I_2)` with
  different component ranks.  This demonstrates why decomposition provenance
  is mathematical data.

The exact-rational Go helper handles diagonal regression fixtures.  It does not
claim to prove the noncommuting theorem; von Neumann's inequality is the
authoritative bridge.

## Zeta instantiation

M8 supplies the orbit decomposition.  Critical blocks are PSD rank-at-most-one
outer products.  Every off-critical reflection-pair block has positive index at
most one, including dependent or degenerate cases.  M9 aggregates these into
`rank(P)<=s1+s2` and `n_plus(Q)<=p`.  M10 supplies the near/far restriction,
trace/operator tail control, and the `O(sqrt(T) log T)` fringe.  M11's Theorem
5.8 moments are reused after the scale-sensitive paper normalization
`G_hat=G_tilde/(aL)`:

```text
tr(G_hat) = (1+o(1)) N,
||G_hat||_F^2 <= (1/lambda + lambda/3 + o(1)) N.
```

For all critical locations on the rank side, `tr(P)<=N_on` and
`N(I')>=N_on+2p` give the finite bound

```text
rank(P) >= 4 tr(A_hat) - 2 N(I') - ||A_hat||_F^2.
```

For simple zeros, regroup `A_hat=P_simple+Q_prime`.  Then
`rank(P_simple)<=s1`, `tr(P_simple)<=s1`, and
`n_plus(Q_prime)<=s2+p`, yielding

```text
s1 >= 4 tr(A_hat) - 2 N(I') - ||A_hat||_F^2.
```

The corresponding all-distinct bound is

```text
s1+s2+p >= (4 tr(A_hat)-N(I')-||A_hat||_F^2)/2.
```

After the M10 tail/fringe transfer and the M11 moment substitution,

```text
H(lambda) = 2 - 1/lambda - lambda/3,
H(1) = 2/3.
```

Therefore

```text
liminf N0_simple(T,2T)/N(T,2T)   >= 2/3,
liminf N0_distinct(T,2T)/N(T,2T) >= 2/3,
liminf N_distinct(T,2T)/N(T,2T)  >= 5/6.
```

These reconstruct the paper's Theorems A-C.  They are not claimed novel.  RH
remains unresolved.

## Research and evidence boundary

`experiments/m12_rank_trace.octest` contains twelve deterministic exact-integer
facts for equality, rejected coefficients, missing PSD, signed trace, dependent
blocks, endpoint budgets, parameter legality, and same-total-matrix ambiguity.
Run it with:

```powershell
go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m12_rank_trace.octest --execution compiled --json
```

OctGo is not used: no natural stable scalar seam beyond the existing Go exact
helper emerged.  `when utility` is not used because the paper statement and the
finite algebra identify the von Neumann route directly.

## Compiler Theory consequence

M11's information ceiling is real for an IR containing only total moments. M12
breaks it by fusing five representations: M8 component identity, M9 rank/index
budgets, M10 localization, M11 total moments, and M12 rank-trace algebra.  The
same total matrix can support incompatible component ranks, so the improvement
cannot be recovered by manipulating the first two moments harder.

The main architectural awkwardness is normalization: the rank-trace inequality
is not scale invariant, while M11 records `G_tilde` and the paper applies the
count theorem to `G_hat`. M12 keeps the adapter explicit. A later typed
scale-change object would be safer than growing string-valued symbolic algebra.

## One next milestone

M13 should type and verify the one-variable test-window functional and reproduce
the Montgomery-Taylor 67.25% optimization without changing M12's finite theorem.
