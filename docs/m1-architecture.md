# M1: quantifiers, domains, and Euler-product lowering

## Scope

M1 asks whether claim strength can be computed from a small structural model
rather than hand-assigned RH-specific capabilities, and whether a change of
representation can make a useful deduction cheap while preserving an honest
proof boundary. It remains a deliberately closed vocabulary, not a parser,
symbolic set engine, or theorem prover.

## Structural claim strength

`semantic.QuantifiedStatement` contains a typed `QuantifierKind`, `Domain`, and
`Predicate`. M1's quantifiers are `ForAll`, `Exists`, and `DensityOne`. RH is now
represented directly as

```text
ForAll ρ ∈ NontrivialZeros(ζ): Re(ρ) = 1/2.
```

The density claim has the same domain and predicate but `DensityOne`
quantification. This is the semantic source of the rejection: density-one
quantification permits a density-zero exceptional set and therefore cannot
discharge a universal target.

The closed `DomainKind` vocabulary contains the complex plane, `Re(s)>1`, the
critical strip, the critical line, nontrivial zeta zeros, and nontrivial zeros
below an explicit unsigned-integer height. `semantic.IsSubset` implements only
selected known relations. In particular,

```text
ZerosBelowHeight(ζ, T) ⊆ NontrivialZeros(ζ)
```

and bounded-zero domains nest when their function agrees and their numeric
bounds order. Unknown inclusion is never guessed.

`compiler.DomainRestriction` is generic. It accepts only a universal source and
a known strict subset, preserves the predicate, emits a directional implication,
and records domain-scope loss. Thus universal truth restricts to a subset, while
the reverse traversal is rejected by both direction and domain coverage.

The M1 graph contains both (1) the conditional consequence RH implies its
bounded form, and (2) an independently certifying hypothetical bounded theorem
fixture. The fixture deliberately asserts no real zero computation. Even when
treated as trusted, it cannot discharge RH because its domain is strictly
smaller.

## Refactoring M0

`Requirements`, `Capabilities`, `PropertySet`, and the RH-specific
`ExceptionalSetSensitivity` were removed. Their job is now performed by the
quantifier/domain/predicate structure. Information-loss records survived, but
their kinds (`quantifier_weakening`, `domain_restriction`, and `approximation`)
are derived from structural endpoint differences or the relation kind; they are
diagnostics and path-history, not a second manually assigned source of truth.

The M0 directional graph, evidence classes, assumptions, obligations,
provenance, graph-owned copies, deterministic ordering, exactness boundary, and
path-wide loss checking all remain. The density-one regression is now rejected
with a structural `quantifier_mismatch` in addition to direction and recorded
path loss.

## Mathematical object and representations

`semantic.Function` identifies `RiemannZeta` independently of representation.
Two `semantic.Representation` values denote it on `Re(s)>1`:

- the Dirichlet series `ζ(s) = Σ n^-s`, exposing additive summation structure;
- the Euler product `ζ(s) = ∏p(1-p^-s)^-1`, exposing prime/multiplicative
  structure and a factor-wise nonvanishing argument.

Their formulas, validity domains, object identity, and affordances are typed
metadata. The representation-change transformation is equivalent only on the
stated domain and has the imported Euler identity as an explicit obligation.
Its provenance includes both the Dirichlet representation and that identity.

This is the M1 case study for a central compiler claim: equivalent
representations can have dramatically different costs for deriving a target
property. Nonvanishing is close to the surface in the Euler representation and
is not a cheap additive deduction from the Dirichlet-series metadata.

## Trusted boundary and exact derivation

The compiler imports, rather than proves, these exact facts:

1. the Dirichlet-series definition on `Re(s)>1`;
2. the Euler-product identity on `Re(s)>1`;
3. absolute convergence of the Euler product there;
4. every Euler factor is finite and nonzero there;
5. the standard nonzero-limit theorem for an absolutely convergent infinite
   product of nonzero factors.

The main external references are [NIST DLMF 25.2.1](https://dlmf.nist.gov/25.2.E1),
[NIST DLMF 25.2.11](https://dlmf.nist.gov/25.2.E11), and
[NIST DLMF §25.10(i)](https://dlmf.nist.gov/25.10.i). DLMF §25.2(iv) cites
Tom M. Apostol, *Introduction to Analytic Number Theory* (1976), §11.5 and
p. 231 for the infinite-product material. The compiler does not validate these
citations or reproduce the convergence proof.

The result has `DerivedEvidence` and exactly these required inputs:

```text
Euler representation of ζ on Re(s)>1
+ absolute convergence on Re(s)>1
+ nonzero Euler factors on Re(s)>1
+ nonzero-limit theorem for the infinite product
⇒ ForAll s ∈ HalfPlane(Re>1): ζ(s) ≠ 0.
```

The identity theorem is already an explicit obligation of certification for the
Euler representation. The nonvanishing pass has one primary representation
input and three additional premises. Certification recursively requires all of
them. If the infinite-product theorem is changed to unverified evidence, the
conclusion remains uncertified and reports that precise open input.

This minimal `Premises` extension avoids pretending that one source proves a
joint consequence. It is not a general hypergraph: transformations retain one
primary source/result edge, and additional premises do not become independently
traversable implication edges.

## Why the result is not RH

Zero-freeness has universal quantification, but it asserts the predicate
`ζ(s)≠0` on `HalfPlane(Re>1)`. RH asserts `Re(ρ)=1/2` on
`NontrivialZeros(ζ)`. The source neither covers the target domain nor states the
same predicate, so the graph reports domain and predicate mismatches without an
RH-specific special case.

## Evidence and determinism

The M0 evidence classes remain: definition, known theorem, derived result,
numerical experiment, and unverified conjecture. Numerical and conjectural
evidence cannot certify an exact claim. No numerical or Oct experiment was
needed for M1 because the acceptance question is structural and the relevant
analytic theorem is explicitly imported.

Both reports preserve insertion order and use dedicated JSON DTOs. The M1 schema
is `riemann.semantic-graph.m1`; quantified claims expose their quantifier,
domain, and predicate fields, while representations and analytic facts expose
their corresponding typed records.

## Architectural findings

Typed quantifier/domain structure explains the motivating invalid inferences
better than capability bits and supports a generic restriction pass. The small
multi-premise extension is sufficient for the Euler argument. The remaining
awkwardness is that theorem imports are still trusted annotations, subset
knowledge is compiled into a closed function, and proposition-specific
derivation validation lives mostly in typed passes rather than in a general rule
schema. Those are useful boundaries for the next experiment, not reasons to
expand M1 into a theorem prover.
