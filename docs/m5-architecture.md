# M5: finite Hermitian lowering

M5 answers its central question affirmatively: once the Weil functional is
restricted to a certified finite complex span and its quadratic-form laws are
made explicit, the compiler can lower it generically to a Hermitian form and
then to an ordered-basis matrix. The lowering is computationally explorable,
but the path retains `function_space_restriction`, so finite matrix positivity
cannot certify universal Weil positivity or RH.

## Typed ladder

The new finite algebraic path is:

```text
QuadraticFunctional
  + QuadraticFormStructure
      Q(lambda f)=|lambda|^2 Q(f)
      parallelogram law
      Q(f) real
        ↓ complex polarization
HermitianForm
        ↓ ordered finite basis
HermitianMatrix
        ↓ finite-coordinate theorem
Q(sum_i c_i f_i) = c* G c
```

`OrderedBasis` has an order-sensitive coordinate identity and contains a
per-member admissibility claim ID. `FiniteSpan` has a separate mathematical
identity invariant under basis permutation. `FunctionClass{Kind:
finite_family}` remains the unordered finite set from M4 and is not accepted as
a span claim. `FiniteLinearCombination` retains its span, complex coefficient
field, dimension, and ordered coefficient vector without introducing a general
symbolic-expression engine.

## Polarization convention

M5 uses a form conjugate-linear in its first argument and linear in its second.
This makes the coordinate expansion, for `G_ij=B(f_i,f_j)`, exactly

```text
B(sum_i c_i f_i, sum_j c_j f_j) = sum_ij conjugate(c_i) G_ij c_j = c* G c.
```

For this convention the encoded identity is

```text
B(f,g) = 1/4 * (
    Q(f+g) - Q(f-g)
  - i Q(f+i g) + i Q(f-i g)
).
```

The frequently printed formula with the opposite sign on the imaginary term
uses the opposite linearity convention. M5 records the argument convention,
coefficient field, normalization, full formula, and theorem ID together, so a
sign convention cannot silently drift.

The polarization source is the [Encyclopedia of Mathematics polarization
identity](https://encyclopediaofmath.org/wiki/Polarization_identity). Its
displayed formula is translated explicitly to the conjugate-linear-first
convention. The Weil functional and its normalization continue to use the M4
Lagarias source. Polarization is unavailable unless the typed complex
homogeneity, parallelogram, and real-diagonal prerequisites are all present.

## Form, matrix, and provenance

`HermitianForm` records its source functional, domain span, convention, entry
definition, recovery and Hermitian identities, component decomposition, and
theorem lineage. `LowerHermitianFormToMatrix` is generic: it accepts a validated
form and an ordered basis for the same span. It does not mention Weil.

Each `MatrixEntry` retains its `(row function, column function)` pair, source
form, source functional, transform convention, theorem lineage, and polarized
zero/endpoint/prime-power/archimedean contributions. Thus the M4 aggregate and
explicit-formula decomposition survives entrywise lowering without a general
symbolic matrix algebra.

The form-to-matrix correspondence is sourced to [David Vogan's MIT Hermitian
Forms notes](https://math.mit.edu/~dav/hermitian.pdf), Proposition 2.1.

M5 needs two matrix value classes:

- `structurally_defined_exact`: entries are exact definitions such as
  `B(f_i,f_j)`, whether or not they have been evaluated to scalars;
- `numerically_evaluated_approximate`: floating-point experimental results.

Hermitian structure is certified from construction. PSD is a separate
proposition. A numerical PSD estimate has approximate exactness and numerical
experiment evidence, and the graph rejects its use for an exact PSD target.

## Finite positivity theorem contract

On the M5 finite span only, trusted contracts connect

```text
ForAll f in V: Q_W(f) >= 0
  ⇔ ForAll c in C^n: c* G c >= 0
  ⇔ G is positive semidefinite.
```

The second equivalence uses the standard Hermitian PSD criterion, documented
in [MIT 8.370/18.435 Lecture
12](https://math.mit.edu/~shor/435-LN/Lecture_12.pdf). The contract is exact and
finite-dimensional. It neither asserts that the open finite PSD claim holds
nor removes the earlier `function_space_restriction` edge.

The graph deliberately contains no rule from finite-family positivity or
nonnegative diagonal entries to span positivity/PSD. The exact toy matrix
`[[0,1],[1,0]]` has zero diagonal while the vector `(1,-1)` gives quadratic
value `-2`; tests require both invalid proof attempts to be rejected.

## Oct experiment

- Path: `experiments/m5_polarization_matrix.octest`
- Purpose: falsify polarization-sign, Hermitian-symmetry, quadratic-recovery,
  diagonal-sufficiency, and basis-permutation mistakes cheaply.
- Command:

  ```powershell
  C:\Users\yuech\source\repos\oct\cmd\oct\oct.exe test `
    C:\Users\yuech\source\repos\Riemann\experiments\m5_polarization_matrix.octest `
    --execution compiled
  ```

- Inputs: decimal toy coordinates for the Hermitian matrices
  `[[1,1+i],[1-i,3]]` and `[[0,1],[1,0]]`.
- Output: six passed, zero failed, compiled six, interpreted fallback zero;
  the adversarial `(1,-1)` coefficient gives `-2` for the off-diagonal toy.
- Compiler identity: local `oct dev`; measured wall time about 3.5 seconds.
- Execution ID and hosted limits: not exposed by the local CLI. The run was
  bounded to one six-fact file with no artifacts or external inputs.
- Precision: Oct `Float`; the chosen small integers are exactly representable,
  but the run supplies no interval bounds, portability claim, or exact
  verification of a numerical backend.
- Octxiliary: not used. Direct local Oct was sufficient, and no experimental
  payload needed to enter Riemann. Adding transport here would have been
  infrastructure without a consuming semantic path.
- `when utility`: not used. There was an obvious first experiment: a complex
  off-diagonal entry that distinguishes the two polarization sign conventions,
  followed by the cheapest diagonal counterexample.

The experiment is research feedback, not proof. It evaluates toy forms, not
admissible Weil basis functions or zeta explicit-formula entries.

## Architectural answers

1. `QuadraticFunctional → HermitianForm → Matrix` is reusable: both lowering
   functions are independent of RH and, at the matrix step, independent of
   Weil.
2. The required new machinery is narrow: scalar-field tags, ordered bases,
   spans, coefficient vectors, quadratic prerequisites, form/matrix records,
   matrix value semantics, and a small set of typed proposition families.
3. The M4 decomposition survives as entry components with full provenance.
4. Span and family identities are cleanly separate; only a span owns arbitrary
   coefficient vectors.
5. Matrix exactness belongs on the value representation, alongside ordinary
   claim exactness. Symbolic exact entries are not approximate merely because
   they are unevaluated.
6. Oct usefully caught the convention sign and exercised the negative
   off-diagonal direction and coordinate permutation.
7. Octxiliary was not useful enough to keep in M5 because no evaluated Weil
   matrix crossed the process boundary.
8. Utility scoring was unnecessary because the discriminating experiment was
   obvious.
9. Proposition-family switch duplication is now visibly awkward across claim
   validation, cloning, reporting, and theorem patterns. It is bad enough to
   target soon, but refactoring it during M5 would not improve the mathematical
   lowering and was therefore deferred.
10. Functional IR hid interference between basis directions. Matrix IR exposes
    off-diagonal phase-sensitive interactions, coordinate changes, Hermitian
    symmetry, and the exact location at which future evaluated contributions
    can be inspected.

## Compiler-theory finding

M5 shows that a mathematical compiler benefits from separating semantic
exactness from evaluability. An unevaluated entry can be exact because its
definition and provenance are exact; a floating-point scalar can be easier to
compute while carrying strictly weaker evidence. It also shows that restriction
loss is a path property: equivalence-preserving algebraic lowerings after the
restriction may change representation radically, but cannot restore discarded
domain coverage.

## M5 boundary

M5 does not implement inertia, rank, trace/Frobenius theorem machinery,
off-critical block classification, pair correlation, percentage-on-line
arguments, automated basis search, or autonomous research orchestration.
