# M2: typed theorem contracts and deterministic composition

M2 answers its research question positively for the deliberately small theorem
set tested here: established mathematics can be represented as reusable typed
lowering rules, instantiated mechanically, and composed without an
Euler-specific validation pass.

## Contract representation

`compiler.TheoremContract` separates a theorem schema from every application.
A schema records a stable theorem ID, typed parameters, typed premise patterns,
a typed conclusion pattern, logical relation, evidence, trust classification,
and source citation. `ContractRegistry` is explicitly constructed and passed to
the engine. It has no global state or initialization magic and enumerates
contracts by theorem ID.

Patterns cover the existing sealed proposition vocabulary:

* representations: object, representation kind, and valid domain;
* representation identities: object, both representations, and domain;
* analytic facts: fact kind, object, and domain;
* quantified statements: quantifier, domain, predicate, and object.

Every field is a typed constant or a reference to a declared typed parameter.
Formulae, affordances, descriptions, and citations are report metadata and are
never semantic match keys. Exactness is matched explicitly.

`BindingValue` is a closed tagged union for objects, domains, representation
kinds, scalars, quantifiers, predicates, and analytic-fact kinds. A repeated
parameter binds once and subsequent occurrences must equal that value exactly.
Thus the same `D` cannot mean `Re(s)>1` in one premise and the critical strip in
another. Contract validation also rejects use of an object parameter in a
domain position before matching begins.

## Application, obligations, and composition

The engine considers certified claims in canonical semantic order. For every
contract it performs exact multi-premise matching and unification, instantiates
the conclusion, and records a `TheoremApplication` containing the schema ID,
bindings, matched claim IDs, conclusion ID, source, and proof-graph edge. The
edge carries the theorem ID, bindings, trust bit, and all premise claim IDs.

If only a consistent subset of premises is present, the engine records a
partial application. Each missing typed pattern becomes a
`TheoremObligation`; any already inferred bindings are substituted into its
description and structured pattern. No conclusion node is emitted for a
partial application, so it cannot be certified accidentally.

Composition is a deterministic fixed point over contracts sorted by theorem
ID and claims sorted by canonical semantic identity. Registration and Go map
iteration order therefore cannot affect the result or JSON. There is no theorem
ranking or heuristic search.

`semantic.SemanticKey` provides canonical identity for the sealed proposition
types. A conclusion already present reuses that claim node. A set of applied
schema/binding/premise tuples prevents repeated applications. Cyclic rules
therefore terminate on the M2 set. When two theorem instances reach the same
semantic conclusion, both `TheoremApplication` records and both transformation
edges remain; the second does not overwrite the first derivation.

The graph remains the final soundness boundary. Contract premises must already
certify, so numerical or unverified evidence cannot satisfy an exact pattern.
The graph also retains M0/M1 assumption, approximation, quantifier, domain, and
information-loss checks. Typed matching is equality, not implicit domain
restriction: a rule bound at `Re(s)>1` cannot silently instantiate at another
domain. The existing explicit `DomainRestriction` rule remains the only modeled
way to restrict a universal conclusion to a known subset.

## Imported M2 contracts

The demonstration registry contains these schemas:

1. `zeta.dirichlet-representation`: trusted definition/import of the Dirichlet
   series on `Re(s)>1`.
2. `zeta.euler-product`: zeta's Dirichlet representation on `Re(s)>1` lowers
   to zeta's Euler-product representation on that same domain. This schema is
   intentionally zeta-specific because a general Dirichlet series need not
   have an Euler product.
3. `zeta.euler-product-absolute-convergence`: trusted exact analytic fact for
   zeta on `Re(s)>1`.
4. `zeta.euler-factors-nonzero`: trusted exact factor fact for zeta there.
5. `analysis.infinite-product-nonvanishing(F,D)`: an Euler-product
   representation, absolute convergence, and nonzero factors for the same
   `F,D` imply `ForAll s in D: F(s) != 0`.

The last schema is generic over its object and domain. Its zeta instance binds
`F = RiemannZeta` and `D = Re(s)>1`. This produces the exact M1 result through
generic matching and composition. `compiler/m1.go` now supplies contract data;
it no longer constructs the Euler representation, zero-free claim, or their
proof transformations directly.

The sources are [NIST DLMF 25.2.1](https://dlmf.nist.gov/25.2.E1) for the
Dirichlet definition, [NIST DLMF 25.2.11](https://dlmf.nist.gov/25.2.E11) and
its cited source Apostol, *Introduction to Analytic Number Theory* (1976),
Theorem 11.7, for the Euler product, [NIST DLMF 1.10(ix)](https://dlmf.nist.gov/1.10.ix)
for infinite-product convergence conventions, and
[NIST DLMF 25.10(i)](https://dlmf.nist.gov/25.10.i) for the stated implication
that the product representation makes zeta nonzero on `Re(s)>1`.

These are trusted imported contracts, not internally proved results. Reports
therefore say “certified relative to trusted imported theorem contracts.”

## Demonstrations

Successful composition:

```text
zeta.dirichlet-representation
  => DirichletRepresentation(zeta, Re(s)>1)
zeta.euler-product
  => EulerProductRepresentation(zeta, Re(s)>1)
analysis.infinite-product-nonvanishing [F=zeta, D=Re(s)>1]
  + AbsoluteConvergence(zeta, Re(s)>1)
  + NonzeroFactors(zeta, Re(s)>1)
  => ForAll s in Re(s)>1: zeta(s) != 0
```

Missing-premise fixture (`OmitEulerFactorsTheorem`):

```text
INSTANTIATE analysis.infinite-product-nonvanishing
  F = RiemannZeta
  D = Re(s)>1
  EulerProductRepresentation       satisfied
  AbsoluteConvergence              satisfied
  NonzeroFactors on Re(s)>1        unresolved obligation
  conclusion                       withheld
```

Mismatch fixture:

```text
EulerProductRepresentation(zeta, Re(s)>1)
AbsoluteConvergence(zeta, CriticalStrip)
=> repeated D conflicts; no complete instance
=> obligation remains AbsoluteConvergence(zeta, Re(s)>1)
```

Focused tests also reject a different object, `DensityOne` for `ForAll`, an
unrelated predicate, a different representation, and exact-looking claims whose
only evidence is numerical.

## What M2 exposes

The main awkwardness is that patterns mirror the closed proposition union by
hand, including matching, validation, instantiation, and JSON projection. This
is small and auditable now, but each new proposition family requires edits in
several exhaustive switches. Also, canonical identity deliberately handles
syntactic structural equality, not mathematical equivalence; equal domains or
objects need explicit structural relations.

Within those limits, contracts look sufficient for ingesting larger, curated
pieces of literature whose applicability can be expressed in this vocabulary.
They are not yet sufficient for theorem statements needing binders inside
terms, side conditions, derived domain expressions, or richer equality. M2
teaches the compiler-design lesson that theorem schemas function like typed
instructions, unification is operand checking, obligations are unresolved
linkage, and saturation is deterministic lowering. The hard problem moves from
bespoke proof code into designing a precise, reviewable semantic instruction
set.
