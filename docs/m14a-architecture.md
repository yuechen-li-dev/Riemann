# M14A architecture: exact class, weak dual, and the remaining witness gap

M14A reaches **Success C**, not the requested numerical-ceiling theorem.  It
reconstructs the broad support-one primal problem, derives a checkable weak-dual
contract, certifies one exact but weak dual witness, and prevents numerical or
finite-truncation evidence from being promoted to a ceiling.  It does not claim
that Anthropic's `0.68185` is correct.

## Source map

| Source | Class and normalization | Objective | Certified result |
|---|---|---|---|
| Anthropic, Remark 1.1 | Unspecified “bandwidth-one data” and configurationwise optimization | Unstated | Only the decimal `0.68185`; no witness or proof |
| CCLM, Corollary 14 | Montgomery–Taylor exact-bandlimited subfamily | The M13 one-parameter bound | Exact subfamily optimum reconstructed in M13 |
| Chirre–Gonçalves–de Laat | `A_LP`: even continuous integrable `f`, `f(0)=hat f(0)=1`, `hat f>=0`, and `f` eventually nonpositive | `Z(f)=r+2/r integral_0^r x f(x) dx` | A rigorously feasible degree-40 certificate with `Z(f)<1.3208` |
| Ramos, EP3.1 | even continuous `g`, `g,hat g in L1`, `g>=0`, `hat g(alpha)<=0` for `|alpha|>=1`; homogeneous `g(0)=1` | `Phi_nu(g)/g(0)`, `nu=delta_0+|alpha|dalpha` on `[-1,1]` | Exact infinite-class formulation and equivalence with the LP sign problem; no sharp constant |

The compiler uses EP3.1 as the broad authoritative class naturally induced by
the support-one pair-correlation data.  The released Anthropic paper does not
prove that this is exactly the unnamed class behind its remark, so that identity
remains an explicit source ambiguity.

## Typed primal problem

With

```text
hat g(alpha) = integral_R g(x) exp(-2*pi*i*x*alpha) dx,
nu = delta_0 + |alpha| 1_[−1,1](alpha) dalpha,
```

the class consists of even continuous `g` such that `g` and `hat g` are
integrable, `g>=0`, `g(0)=1`, and `hat g(alpha)<=0` outside the exact radius-one
data interval.  The extremal values are

```text
c_* = inf_g Phi_nu(g)/g(0),
J_* = sup_g (2-Phi_nu(g)/g(0)) = 2-c_*.
```

This is a sign cutoff, not compact Fourier support.  M13's
Montgomery–Taylor function has exact Fourier support in `[-1,1]`; exact zero
tail implies the required nonpositive tail, which mechanically proves its
membership after the explicit homogeneous normalization.

The CGdL scaling `g(x)=hat f(x/r(f))` sends their LP class into this class and
turns `Z(f)` into `Phi_nu(g)/g(0)`.  Their certified feasible result therefore
gives `c_*<1.3208`, or `J_*>0.6792`.  It is a primal lower bound on the attainable
simple-zero proportion, not an upper ceiling.

## Weak dual and exact baseline

A number `c` is a valid lower bound on `c_*` if there is a nonnegative tempered
measure `sigma` supported outside `(-1,1)` for which

```text
P = nu - c dx + sigma
```

is a positive-definite tempered distribution.  For every primal-feasible `g`,
positive definiteness and the tail sign give

```text
0 <= <P,hat g>
   = Phi_nu(g)-c g(0)+<sigma,hat g>
<= Phi_nu(g)-c g(0),
```

so `J(g)<=2-c`.  This uses weak duality only; M14A assumes neither strong
duality nor numerical primal/dual agreement.

The exact baseline takes

```text
c = 1,
sigma = 1_{|x|>1} dx,
P = delta_0-(1-|x|)_+ dx,
hat P(xi) = 1-(sin(pi xi)/(pi xi))^2 >= 0.
```

Thus `c_*>=1` and `J_*<=1`.  Combined with CGdL:

```text
1 <= c_* < 1.3208,
0.6792 < J_* <= 1.
```

This bracket contains `0.68185` but does not support it.

## Oct-first exploration

The main Oct experiment studies the deliberately small completion family

```text
sigma_c = c 1_{|x|>1}dx + (c-1)(delta_-1+delta_1).
```

Its Fourier density has local coefficient `(9-8c)/12`, proving that this
restricted family cannot pass the origin test for `c>9/8`.  Oct found the same
`c=9/8` boundary by bounded scan and rejected `c=1.126`.  This is useful
counterexample evidence, not a full-class dual witness.

Runs from `C:/Users/yuech/source/repos/oct`:

```text
go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m14a_dual_completion.octest --json
  interpreted: 6 passed, 0 failed, 4.811 s

go run ./cmd/oct test C:/Users/yuech/source/repos/Riemann/experiments/m14a_dual_completion_compiled.octest --execution compiled --json
  compiled: 3 passed, 0 failed, 0 fallback, 0.633 s

go run ./cmd/oct artifact C:/Users/yuech/source/repos/Riemann/experiments/m14a_plot_artifact.octest --output-root C:/Users/yuech/source/repos/Riemann --json
  build-time interpreted artifact: 1 passed, 0 failed, 0.395 s
```

`[Artifact]` resolves the immediate plotting friction: it runs `PlotLine` even
though the compiled Fact lane cannot lower it.  The probe also records a future
Oct issue: the PNG is written, but the artifact command reports `0 produced`
and does not attribute the Plot output to `--output-root`; an explicit path was
required.  The plot is therefore a local file side effect, not a retrievable
artifact ID.  Python was used only as an independent implementation oracle and
does not contribute theorem evidence.

## Exact obstruction and compiler lesson

The audited sources supply no outside measure near `c=1.31815`.  A finite
frequency grid does not prove positivity on the whole line, and a finite basis
or SDP truncation does not control omitted directions.  The smallest remaining
mathematical obstruction is therefore concrete:

```text
construct sigma >= 0 outside (-1,1)
and certify P=nu-c dx+sigma positive definite on all of R,
with exact tail/completeness control and useful c>1.
```

M14A makes evidence levels first-class, so grid positivity, missing tails, wrong
normalization, and support beyond one are rejected before theorem construction.
This is the Compiler Theory result: search-space pruning needs a whole-class
dual object; numerical convergence cannot serve as negative knowledge.

A generic SDP framework and a first-class `RepresentationCeiling` are not yet
justified because no useful ceiling exists.  The narrow `DualCompletionWitness`
and `CertifiedExtremalBound` types pay for themselves by making the missing
obligations executable.

## One next milestone

M15: build a purpose-specific verifier for periodic/tempered positive-definite
completions of EP3.1, including exact whole-line and tail control, and use it to
search for the first nontrivial dual witness beyond the analytic `c=1` baseline.

RH remains unresolved.
