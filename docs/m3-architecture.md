# M3: functional equation and typed symmetry transport

M3 answers its research question positively within the zeta fragment: symmetry
can be compiled as a typed information-transport operation. A theorem contract
binds a structural point, checks ordinary premises and typed side conditions,
applies a recorded point transform, and emits a target-point proposition whose
proof edge retains every input and theorem citation.

## Mathematical normalization and trust boundary

M3 uses Riemann's completed xi function exactly as normalized by NIST DLMF
25.4.3–25.4.4:

```text
ξ(s) = 1/2 s(s-1) Γ(s/2) π^(-s/2) ζ(s)
ξ(s) = ξ(1-s).
```

`RiemannXi` is a distinct typed object. Its `CompletedXiRepresentation`
structurally records `RiemannZeta` as its base object, not merely in the formula
string. The formula is display metadata; object identity, base-object identity,
representation kind, and domain are semantic fields.

The imported trust boundary is:

* the Dirichlet series remains valid only on `Re(s)>1` ([DLMF 25.2.1](https://dlmf.nist.gov/25.2.E1));
* the meromorphic analytic continuation is a separate representation on
  `C except {1}`, and zeta's only singularity is the simple pole at `s=1`
  ([DLMF 25.2(i)](https://dlmf.nist.gov/25.2.i));
* the xi definition and functional equation are imported from
  [DLMF 25.4.3–25.4.4](https://dlmf.nist.gov/25.4);
* conjugation symmetry, the negative-even trivial zeros, nonvanishing on
  `Re(s)=1`, critical-strip confinement, and the two stated symmetry axes are
  imported from [DLMF 25.10(i)](https://dlmf.nist.gov/25.10.i);
* M1's `Re(s)>1` zero exclusion remains derived from the Euler product and its
  trusted infinite-product inputs.

The global set-invariance and strip-confinement schemas are classified as
trusted external mathematics because the current compiler matches their typed
premises but does not itself prove universal lifting or real-part inequalities.
The composition of two invariances, equality-based zero transport, point-domain
specialization, and multiplication/division by a known finite nonzero factor
are compiler-verified rules. Reports therefore correctly say “certified
relative to trusted imported theorem contracts.”

No numerical or Oct experiment was used in M3.

## Structural points and transform composition

`semantic.PointExpr` contains a base symbol, a `PointTransform`, and one of the
two fixed-set annotations actually needed by M3. The transform is the closed
Klein-four action:

```text
identity                 s
conjugate                conjugate(s)
one_minus                1-s
one_minus_conjugate      1-conjugate(s)
```

Composition is structural and canonical: conjugation and one-minus commute,
each transform is an involution, and their composition is critical-line
reflection. This is deliberately not arbitrary symbolic algebra.

Fixed-set information is also typed. A point known to be on the critical line
canonicalizes `CriticalReflection(point)` to the point itself and canonicalizes
`OneMinus(point)` with `Conjugate(point)`. Consequently the closure API always
reports four generated transforms but deduplicates them to two locations for a
generic nonreal critical-line zero. It never calls that a quartet of four
distinct zeros.

## Transport-capable contracts and side conditions

The theorem term vocabulary now includes `point`, `zero_class`,
`point_transform`, side-condition, zero-set-property, and zero-classification
parameters. A point term can carry either a fixed transform or a bound transform
parameter. `analysis.functional-identity-zero-transport` is consequently
generic:

```text
Zero(F, P, Z)
FunctionalIdentity(F, identity, T, D)
--------------------------------------
Zero(F, T(P), Z)
```

The proof transformation is labeled `transport-theorem-contract`; its bindings
show `P` and `T`, so the movement is inspectable in human and JSON reports.

Side conditions are a distinct structural list on a theorem contract, not
string-named premises. M3's conditions are the minimum forced by the real
mathematics:

* analytic continuation is available on a typed domain;
* a bound point lies in that validity domain;
* the xi completion factor is finite and nonzero at a bound point;
* zeta's real-analytic conjugation property is available on its finite domain.

The matcher processes ordinary premises and side conditions in fixed contract
order against semantically sorted certified claims. Missing conditions become
typed `TheoremObligation` values. No conclusion is emitted until all conditions
match. There is no fuzzy text matching, theorem-ranking heuristic, registry
insertion dependence, or Go-map iteration dependence.

## Exact pointwise derivations

For a conditional source claim `Zero(ζ,ρ,nontrivial)`, conjugation is:

```text
simple pole classification
  => PointInDomain(ρ, C except {1})
analytic-continuation representation
  => AnalyticContinuationAvailable(C except {1})
real-analytic conjugation theorem
Zero(ζ,ρ,nontrivial)
  => Zero(ζ,conjugate(ρ),nontrivial).
```

The Dirichlet-series claim cannot substitute for the continuation premise; the
focused omission test leaves conjugation unresolved outside `Re(s)>1`.

Functional reflection takes the richer xi route:

```text
Zero(ζ,ρ,nontrivial)
+ pole/trivial-zero classifications
  => CompletionFactorRegularNonzero(ρ)
+ xi definition
  => Zero(ξ,ρ,nontrivial)
+ ξ(s)=ξ(1-s)
  => Zero(ξ,1-ρ,nontrivial)

Zero(ζ,ρ,nontrivial)
+ M1 zero-free Re(s)>1
+ pole/trivial-zero classifications
  => CompletionFactorRegularNonzero(1-ρ)

Zero(ξ,1-ρ,nontrivial)
+ CompletionFactorRegularNonzero(1-ρ)
+ xi definition
  => Zero(ζ,1-ρ,nontrivial).
```

The reflected completion-factor premise is important: it rules out exceptional
factor behavior without assuming the critical strip in advance. Omitting the
factor theorem leaves an explicit condition obligation and withholds the
reflected zeta zero. Conjugation and reflection then saturate to the reachable
four-transform closure. Every derived claim retains the source hypothesis and
complete theorem lineage.

## Global geometry and RH rigidity

The zero set receives typed `InvariantUnderTransform` properties. DLMF's
conjugation result yields invariance under conjugation. The xi equation, xi/zeta
relationship, M1 right-half-plane exclusion, and pole/trivial classifications
yield invariance under one-minus. The generic transform-composition contract
then derives invariance under `one_minus_conjugate`, the reflection about
`Re(s)=1/2`.

Critical-strip confinement has this visible contract instance:

```text
ζ has no zeros for Re(s)>1                         [M1]
NontrivialZeros(ζ) invariant under s -> 1-s       [M3]
trivial zeros are exactly the negative evens      [import]
ζ has no zeros on Re(s)=1                         [import]
----------------------------------------------------------
NontrivialZeros(ζ) confined to 0 < Re(s) < 1.
```

Reflection handles the left side and the `Re(s)=0` boundary; the imported
Hadamard–de la Vallée Poussin boundary theorem handles `Re(s)=1`. The old domain
algebra no longer assumes `NontrivialZeros(ζ) subset CriticalStrip`; the result
exists as a derived, inspectable zero-set property.

The geometric equivalence is registered as an exact theorem contract and used
to normalize the unproved target:

```text
Re(ρ)=1/2  <=>  1-conjugate(ρ)=ρ.
```

Both RH formulations remain uncertified. The report sharpens the remaining
defect to: exclude off-axis symmetry orbits inside the critical strip.

## Architectural findings and Compiler Theory

Typed information transport is reusable at the pointwise level. A functional
identity and zero claim use the same generic contract regardless of the bound
function or transform; zeta-specific knowledge resides in imported identities,
classification facts, and factor conditions.

Side conditions have clearly become a general semantic concept: they bind
parameters, participate in deterministic matching, block emission, and survive
as inspectable obligations. Small compositional-domain needs were real, but M3
needed only two closed domains (`C except {1}` and `Re(s)=1`), not a general set
expression algebra. Transform canonicalization is manageable because the
mathematical action is finite and closed.

The theorem-contract model remains adequate for pointwise transport. It is less
adequate for lifting a pointwise rule to universal set invariance and for proving
real-part region arithmetic: M3 still imports specialized global contracts for
those steps. Proposition-family switch duplication also grew materially across
validation, matching, instantiation, semantic keys, descriptions, and JSON.
That duplication is now costly, but replacing it in M3 would have obscured the
mathematical experiment.

Compiler Theory interpretation: point transforms are typed operands, functional
identities are transport instructions, side conditions are proof-carrying
preconditions, unresolved conditions are linkage failures, orbit saturation is
a deterministic dataflow fixed point, and canonicalization is semantic common
subexpression elimination. The result preserves the distinction between four
generated operations and the cardinality of their quotient orbit.

## Exactly one proposed next milestone

**M4 — Quantified transport lifting.** Add the smallest binder-aware rule that
lifts certified pointwise transport into set invariance and transformed-domain
facts, then use it to replace M3's specialized global invariance/confinement
contracts. This is the next demonstrated compiler gap: pointwise transport is
generic, while its universal/set-level consequence is still imported.
