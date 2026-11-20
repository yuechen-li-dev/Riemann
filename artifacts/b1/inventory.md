# B1 ranked inventory

Revision `a172c464fd4bcfc4f53e394dcba45d742b1409a6`. Oct scores prioritize review. Proof status is source-level structural review for every record.

## 1. Cyclic resource allocation with a sharp dual

For directed odd-offset flows on an even cycle, maximum throughput and minimum dual cost agree. The deficit relative to a fixed multiple of total budget has a greatest dual value. An optimizer can satisfy the inequalities even for odd offsets beyond the chosen period.

Oct score **12** (literature, named, mathlib-only, extractable, comparison-open). Finite optimization.

[Pinned declaration Formalism.IncomingBounds.lemmaFa](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/IncomingBounds/WeightIdentities.lean#L17).

Formalization novelty **unresolved**. Independent computation **exact_finite_replay**. 405 exact two-letter primal/dual certificates, bases 3–7 and integer budgets 0–8.

- **literature** (special_case). The finite primal/dual optimization is a special structured LP. Bounds on dual coordinates and extension to odd offsets are specific to these coefficients. [Schrijver, Combinatorial Optimization, Chapter 47](https://www.cis.upenn.edu/~jean/math-deep.pdf)
- **formalization** (special_case). General LP duality is formalized there. A separate instantiation is needed for the clipped and extended dual. [Isabelle AFP, Duality of Linear Programming](https://www.isa-afp.org/entries/LP_Duality.html)
- Open. Replay the custom finite LP proof and the extension to all odd offsets.
- Open. Extend certificate experiments beyond d=2 before considering a reusable family.

## 2. Functorial resolution in characteristic zero

One rule resolves reduced finite-type schemes over characteristic-zero fields, preserving their smooth loci and producing a strict normal-crossings reduced inverse image of the nonsmooth locus. The tower is compatible with the specified smooth maps and field extensions.

Oct score **5** (literature, named, mathlib-only). Algebraic geometry.

[Pinned declaration FunctorialResolution.functorialResolutionCharZero](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/ResolutionOfSingularities/HasSimpleNormalCrossingsOfInStrictNormalCrossing.lean#L15).

Formalization novelty **already_in_mathlib_or_other_project**. Independent computation **not_run**.

- **literature** (unresolved). Known functorial resolution, with exact scope and tower semantics compared separately below. [Kollár, Lectures on Resolution of Singularities, Theorem 36](https://doi.org/10.1515/9781400837086)
- **formalization** (unresolved). The external endpoint uses reduced equidimensional schemes and a BlowUpSequenceFunctor. The corpus accepts mixed dimensions and states explicit pairwise pullback matches. Mutual implication is unresolved. [resolutionorg/hironaka, Theorem 36](https://github.com/resolutionorg/hironaka/blob/044d89c9455859446fc5e986667319f910e15f4b/Hironaka/Resolution/Algebraic/Kol07/Thm36/Theorem36.lean)
- Open. Compare coherent functoriality, empty steps, center hypotheses and treatment of mixed-dimensional schemes.
- Open. Check the transitive Lean proof against the intended statement.

## 3. Chebotarev for conjugacy-stable sets

Among the unramified primes, the primes whose Frobenius lies in a chosen union of conjugacy classes have Dirichlet density equal to the fraction of group elements in that union.

Oct score **3** (literature, comparison-open). Analytic number theory.

[Pinned declaration GroundTheory.Global.hasDirichletDensity_chebotarevSet](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/GroundTheory/GlobalFields/TateThesis/Chebotarev.lean#L198-L200).

Formalization novelty **unresolved**. Independent computation **not_run**.

- **literature** (matches). The standard Dirichlet-density theorem for a conjugacy class extends by finite disjoint union. The corpus states the union form. [Chebotarev density, class field theory course](https://kokunoyumeto.github.io/open-math-courses-public/courses/NT-CFT/the-chebotarev-density.html)
- **formalization** (stronger). The corpus quantifies over all conjugacy-stable sets, including a single class. The external header states the single-class version and its proof at this revision is sorry. A finite additivity bridge would recover the union form. [CBirkbeck/chebotarev-density at c64095e](https://github.com/CBirkbeck/chebotarev-density/blob/c64095e6cc6483b401849c7fd9182d983d3bf261/CebotarevDensity/Main.lean)
- Open. Verify the imported single-class, abelian and degree-one density proof closure and its axioms.
- Open. Compare the ideal and height-one spectrum APIs and the Dirichlet-density normalization in a Lean bridge.

## 4. Brauer induction for field automorphism groups

Each complex character of the automorphism group of a finite-dimensional field extension is an integer combination of characters induced from one-dimensional characters of elementary subgroups.

Oct score **3** (literature, comparison-open). Finite group representation theory.

[Pinned declaration GroundTheory.LocalField.IsCharacter.mem_spanZ_isElementary](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/GroundTheory/LocalFields/ArtinConductorBrauer.lean#L181-L182).

Formalization novelty **unresolved**. Independent computation **not_run**.

- **literature** (special_case). Classical Brauer induction ranges over arbitrary finite groups; this endpoint uses field automorphism groups and ordinary characters. The generic proof route assembles with monomiality. [Brauer, On Artin’s L-Series with General Group Characters (1947)](https://www.jstor.org/stable/i307231)
- **formalization** (unresolved). The checked claim described there is a consequence about values, with a distinct statement and dependency route. It does not establish availability of the same induction API. [Lean evaluation, character values in a cyclotomic field](https://lean-lang.org/eval/problems/brauer_character_in_cyclotomic/)
- Open. Check standalone replay of the generic finite-group assembly with monomiality supplied.
- Open. Compare the generic conclusion to the virtual character formulation under scalar extension.

## 5. Krasner stabilizers inside a decomposition group

An automorphism preserving the valuation and fixing a sufficiently close approximation must fix the original element as well. The conclusion is containment of stabilizers in the decomposition group.

Oct score **3** (literature, comparison-open). Valuation theory.

[Pinned declaration GlobalTheory.krasner_stabiliser](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/InertiaEigenvalues/KrasnerLemmaDecompositionGroupStabilizerVersion.lean#L18-L20).

Formalization novelty **unresolved**. Independent computation **not_run**.

- **literature** (unresolved). The usual complete-field Krasner lemma concludes containment of generated fields. The corpus isolates an ultrametric automorphism argument with a different scope. [Schütz, Algebraic Number Theory, Chapter 8](https://janschuetz.perso.math.cnrs.fr/lehre/ANT.pdf)
- **formalization** (unresolved). Mathlib concludes field membership in an ultrametric normed-field setting. Here the valuation may have a general ordered value group and the output is a subgroup stabilizer inclusion. [Mathlib.Analysis.Normed.Field.Krasner](https://leanprover-community.github.io/mathlib4_docs/Mathlib/Analysis/Normed/Field/Krasner.html)
- Open. Compare via a fixed-field correspondence on common hypotheses.
- Open. Run Lean replay on the pinned valuation imports or devise a statement bridge to a standalone valuation-preserving action.

