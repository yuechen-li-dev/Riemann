# M11: first and second moments to a thresholded positive index

M11 answers the bounded question affirmatively. For a finite Hermitian
`d x d` matrix `G`, threshold `theta >= 0`, and theorem-grade finite bounds

```text
tr(G) >= A,
tr(G^2) = ||G||_F^2 <= B,
A - d theta > 0,
B > 0,
```

it derives

```text
n_plus^theta(G) >= ceil((A-d theta)^2/B).
```

This is the paper's thresholded Cauchy-Schwarz lemma, represented as a generic
finite theorem rather than a zeta-special formula. Approximate or raw
asymptotic moments cannot satisfy its finite premises.

## Exact finite proof

Let `J={i: lambda_i>theta}`, `m=#J`, and `S=sum_{i in J} lambda_i`. Because
`theta>=0`,

```text
S >= tr(G) - theta #{i not in J} >= A-d theta > 0.
```

Negative eigenvalues in the remainder only improve this inequality. Then

```text
S^2 <= m sum_{i in J} lambda_i^2 <= m tr(G^2) <= mB.
```

The real bound follows, and integrality gives the ceiling. The executable
version uses exact rationals and rejects infeasible moment premises.

## Paper normalization and imported moments

The source uses `l=log(T/(2pi))`, `L=lambda*l`, spacing `h=2pi/L`,

```text
d=floor(T/h)=floor(LT/(2pi)),
G_tilde_T=G_T/L.
```

For fixed `0<lambda<=1` and the fixed-width taper (`w=1`), Theorem 5.8 gives

```text
tr(G_tilde_T) = L N(T,2T) (1+O(E_T)),
||G_tilde_T||_F^2
  = lambda l^2 N(T,2T) (1+lambda^2/3) (1+O(E_T)),
E_T=o(1).
```

M11 retains the error terms and adds explicit epsilon-to-eventual-finite-bound
contracts before using the finite theorem.

## Threshold scaling

M10 supplies

```text
theta0=O(l T^(lambda/2-1)).
```

The relevant check is not merely `theta0=o(1)`. Since
`d~LT/(2pi)` and `tr(G_tilde_T)~L N~LTl/(2pi)`,

```text
d theta0 / tr(G_tilde_T) = O(T^(lambda/2-1)) = o(1).
```

Thus the dimension-amplified threshold penalty is negligible for every fixed
`0<lambda<=1`.

## Asymptotic result and M10 reuse

Writing

```text
F(lambda)=lambda/(1+lambda^2/3),
```

the theorem gives

```text
n_plus^theta0(G_tilde_T) >= (F(lambda)-o(1))N(T,2T).
```

At `lambda=1`, `F(1)=3/4`. M10's existing Proposition 4.5 bridge gives

```text
N0_simple(T,2T)
  >= 2 n_plus^theta0(G_tilde_T)-N(T,2T)-2N(I'\I),
N_distinct(T,2T)
  >= n_plus^theta0(G_tilde_T)-N(I'\I).
```

The retained fringe satisfies `N(I'\I)=O(sqrt(T)log T)=o(N(T,2T))`.
Consequently

```text
liminf N0_simple(T,2T)/N(T,2T) >= 1/2,
liminf N_distinct(T,2T)/N(T,2T) >= 3/4.
```

This reconstructs the paper's earlier Cauchy-Schwarz stage. It is not claimed
as novel mathematics. M11 does not implement the stronger rank-trace lemma.

## Counterexample-first boundary

Exact fixtures reject: omitting `d theta`; allowing negative thresholds in the
same proof; using a lower rather than upper second-moment bound; inferring many
directions from positive trace alone; changing strict `>` to `>=`; and omitting
the Frobenius upper bound. Equality spectra at `theta=0` show the generic
Cauchy-Schwarz theorem is sharp from these aggregate facts.

The bounded Oct experiment is
`experiments/m11_moment_count.octest`. It uses exact integer diagonal spectra
and cross-multiplied inequalities. Its role is counterexample filtering and
sharpness exploration only.

## Compiler-theory lesson

M11 is a four-way representation fusion: analytic moments supply magnitude,
spectral partitioning turns magnitude into directions, M10's window semantics
control the threshold and fringe, and the zero-orbit representation converts
directions into zero counts. The key compiler obligation is evidence staging:
an asymptotic statement must become an eventual finite inequality before it can
enter an exact finite theorem.

RH remains unresolved.
