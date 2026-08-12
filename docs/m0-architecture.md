# M0: semantic claim IR and sound lowering

## Hypothesis and scope

The experiment asks whether a compiler-style intermediate representation can
make mathematical transformations mechanically honest. Rather than hiding a
reformulation inside prose, it records the proposition, logical direction,
assumptions, evidence, obligations, provenance, and information discarded by
each pass.

M0 is intentionally not a theorem prover or a symbolic algebra system. Its
mathematical vocabulary contains only the Riemann zeta function, RH, the
universal critical-line formulation, a density-one formulation, and a small
named-obligation escape hatch. There is no parser or surface theorem language.

## Why claims, not expressions

An expression records mathematical syntax. The unit needed for sound lowering
is an assertion in context: what is claimed, under which assumptions, supported
by what kind of evidence, and usable for which target. `semantic.Claim` is
therefore the root IR value. Its proposition is a sealed typed interface rather
than a generic symbolic tree or arbitrary string.

A claim has two property sets:

- `Requirements` says what semantic capability evidence must retain to
  discharge this claim.
- `Capabilities` says what the current representation actually retains.

M0 needs one property, `ExceptionalSetSensitivity`. Both RH representations
retain and require it. The density-one representation does not retain it.

## Relations and proof state

`compiler.Transformation` has one source and one result in M0, plus its relation,
open obligations, declared losses, and provenance. The relations mean:

- `Equivalent`: an exact identity usable in either direction. Both endpoints
  must have the same assumptions.
- `Implies`: an exact implication usable only from `From` to `To`.
- `Relaxation`: a directional weakening, also usable only from `From` to `To`.
  It is distinct from implication so reports can say that weakening was the
  purpose of the pass.
- `Approximation`: a directional approximate connection. It cannot certify an
  exact target.

The graph rejects a transformation that silently drops assumptions or fails to
declare a lost capability. M0 does not permit a pass to synthesize a capability
that its source lacks; a future recovery mechanism would need semantics stronger
than an annotation. Discharge checks also accumulate losses across the whole
path, so a later pass cannot hide an earlier loss.

Assumptions travel with claims. A directional transform may add assumptions but
cannot remove existing ones. An equivalence cannot change the assumption set.
Obligations are claim IDs on the transformation. They remain graph nodes, appear
in JSON, and block both structural use and certification until their own exact
evidence certifies them.

## Evidence and trust boundary

Evidence is one of:

- definition / semantic identity;
- known theorem / external mathematical result;
- derived result;
- numerical experiment;
- unverified conjecture.

Definition and known-theorem evidence are explicit trusted inputs. Derived
evidence is certified only through a permitted transformation from certified
premises with all obligations certified. Numerical evidence and conjectures are
retained as provenance but cannot certify an exact theorem. M0 does not validate
citations or prove imported theorems; those are an explicit trust boundary.

The RH root is marked `unverified_conjecture`, not `known_theorem`. Its normalized
and density forms are marked `derived`, with notes that these are conditional
compiler consequences, not proofs that RH is true.

## The M0 lowering

`normalize-rh` maps

```text
RiemannHypothesis(ζ)
    ⇔
every nontrivial zero ρ of ζ has Re(ρ) = 1/2.
```

This is deterministic, lossless, and reversible. The formulation follows the
[NIST Digital Library of Mathematical Functions, §25.10(i)](https://dlmf.nist.gov/25.10.i)
and the [Clay Mathematics Institute problem description](https://www.claymath.org/millennium/Riemann-Hypothesis/).

`critical-zero-density` then records the exact conditional consequence

```text
all counted nontrivial zeros are on the line
    ⇒
limiting on-line / total zero-count ratio = 1.
```

M0 leaves the counting functions abstract. Under this abstract definition, the
step is immediate: if the counts agree, their ratio is one wherever defined.
It is not recorded as an external theorem and it does not assert RH.

The converse fails. A limiting ratio of one only says that the off-line count is
lower order than the total count. A finite or infinite density-zero exceptional
set can remain. Thus a density-one result cannot exclude every exceptional zero,
while RH must. The pass records loss of `ExceptionalSetSensitivity`, and the
proof checker independently finds both the non-reversible implication and the
missing capability.

## Determinism and reports

The graph owns insertion-ordered slices and never derives report order from Go
map iteration. The CLI renders a stable human report or a versioned JSON DTO
(`riemann.semantic-graph.m0`). Provenance lineage is deterministic and can be
traced from the density claim through the zero-location claim to the RH root.

## What M0 taught us

The claim-first model was enough to expose the motivating error without symbolic
algebra. Logical direction alone rejects the reverse inference, but it does not
explain why the converse is mathematically unavailable; capability loss supplies
that explanation. Conversely, capability checking alone is not enough: it must
be accumulated across a path, not compared only at its endpoints. Obligations
and evidence classes fit naturally into the same graph and close two other easy
soundness holes.

The awkward part is the split between `Requirements` and `Capabilities`. It is
useful here, but `ExceptionalSetSensitivity` is specialized and manually
assigned. M0 also models each pass as one source and one result, and treats
definition/known-theorem evidence as trusted annotations. These choices are
provisional. A second genuinely different example is needed before generalizing
the property vocabulary, multi-premise rules, or the external-evidence boundary.

No numerical experiment was needed: the M0 question is about logical structure,
and sampling zeros would add no evidence about the invalid universal inference.
