# A2 audit · Simultaneous admissibility

**Success B for the selected implication.** With the original point-integrality, period, column, gain and compatibility inputs, the off-twist theorem gives the asserted algebraic representative and valuation on one and the same period. Its multiplication proof, including the zero case, is reconstructed in the [target report](a2-simultaneous-admissibility-target.md). We found no counterexample to that statement. Whether the submitted construction realizes the necessary arithmetic data with sufficient gain for RH remains open in this audit.

## What is shared

| Question | Finding |
|---|---|
| Does the theorem construct a shared witness? | It constructs y for a prescribed word, place and pair of lattice vectors. Uniqueness follows from embLL being injective. It takes the tuple selector, column system and numerical records as input. |
| Do valuation and gain apply to the same object? | They share X,K,L,w and the lift indices. The column estimate applies to the selected tuple and actual vector. GainLe is a scalar inequality indexed by place and stage. It lowers the threshold already proved for that column. There is no second period to match. |
| Where is a genuine common choice required? | SelectionServed.twU requires a single indexed chain datum LC supporting both CmpV and SysTw. The service proof destructures exactly that witness. This is explicit upstream evidence in the submitted code and an open realization obligation in Riemann. The sibling pdv field similarly bundles six restrictions under one choice of ep,be,sg. [Contract](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AutomorphicArithmetic/ThetaCorrespondence/ThetaPairLifts/ThirdPlaceSeries/GServT.lean#L13) |
| Are local choices globally compatible? | For the local conclusion, one prescribed selector is used, including agreement on conjugate places. LC can assign different chains at different places. A single globally constant twist or scalar y across distinct periods would be a stronger requirement. When passing to global divisibility, A1 additionally requires the glued row lattice, local unit identities and shift localizations. The Go models do not construct those arithmetic objects. |
| Are parameter ranges nonempty? | The intersections of eventual tails and the projected lift-count intervals below are nonempty under the stated hypotheses. Existence of a family carrying **all** the period, rank, cost, chain and window data for each sufficiently small ε remains unverified here. |
| Does gain survive composition uniformly? | The conditional local sum inequality holds for every finite m and loses no extra constant in that comparison. Neither the inequality nor zero-valued periods give a positive global saving. The service bound permits C and the tail threshold to depend on the family. At the analytic endpoint eD and eL precede the choice of ε. |
| What remains? | Establish the shared local-system lower bounds inside SysTw on the actual selected rows. No source contradiction or failure of the original hypotheses was exhibited. |

## Exact arithmetic attacks

The finite model has one proposed indexed datum per candidate name at a fixed family index. At specified primes its chain is represented by multiplication by rational numbers c₁,…,cₘ on one-dimensional lattices. The actual column is multiplication by b. For all inputs and all thresholds the model valuation transfer is equivalent to vₚ(b)≥vₚ(∏cⱼ), including the zero maps. This follows by testing the input 1 and then using multiplicativity. Both zero maps satisfy transfer, while a zero chain can transfer arbitrary bounds only to a zero column.

For one complementary direction, the model system requires fₚⱼ≤vₚ(cⱼ) for every j and Σⱼ fₚⱼ−lossₚ≥lowₚ. Loss represents the total stage debit at p. Inputs explicitly record the proposed f values. Every candidate must pass both tests for each listed prime. The integer and rational comparisons are exact. The replay prints the prime, candidate, truth of each condition and first failed check. [Implementation](../audit/claudegrh/selection.go) · [Persisted inputs](../audit/claudegrh/testdata/selection.json)

| Attack | Exact result |
|---|---|
| A. Separate existential witnesses | At 2 take b=1, low=1 and loss=0. A has c=1, f=0 and fails the aggregate lower bound. B has c=2, f=1 and fails valuation transfer. Each property has a witness and the intersection is empty. `REDUCED_MODEL_COUNTEREXAMPLE`. Two candidates are minimal, since separate satisfiability on a singleton implies joint satisfiability. |
| B. Place-dependent choices | With b=6, primes (2,3), low=(1,1) and zero loss, A with c=2 passes both tests only at 2, B with c=3 only at 3. A selector taking A at 2 and B at 3 exists. A constant selection from {A,B} fails. `REDUCED_MODEL_COUNTEREXAMPLE` to requiring that constant selection. Adding a third candidate c=6 with f=1 at both places gives `FINITE_VERIFIED`. This is precisely where restricting the permitted set of chains changes the answer. |
| C. Accumulated loss | At 2 take b=4 and chain (2,2), f=(1,1), loss=1. The net bound is exactly 1. Threshold 1 passes (`FINITE_VERIFIED`). Raising it by 10⁻¹⁸ fails the aggregate check. Replacing the first f by 2 fails the first stage even when the total suffices. These failed candidates are `INCONCLUSIVE` for the source. |
| D. Degeneracy | Changing the preceding example to zero column and zero first factor makes both transfer and arbitrarily large finite first-stage lower bounds possible (`FINITE_VERIFIED` for the tested instance). An empty candidate list supplies no selection (`INCONCLUSIVE`). In the source VV=∅ triggers a separate membership and column argument. It is distinct from an empty tuple type, which already contains μ₀. The admissible reference tuple uses the base-period formula. Cleared local denominators are units at w. |
| E. Parameter intersection | Let 0<ε≤1/2 and x≥1. Put q=x^(1/2−ε/16). Since q≥1, m=⌈q⌉ lies in [q,2q] and m≥x^(1/2−ε). Thus even an integer lift count satisfies both the illustrative SizesTextAt bracket and the RestSize lower bound. This is an independent exact inequality. Full data satisfying SizesTextAt or ChainTied remain `INCONCLUSIVE`. [Size predicates with distinct scopes](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SizeWindowBrackets/SChain.lean#L207) |

The finite prime list, rational scalars, fixed finite menu of chains, specified f values and one complementary direction are modeling assumptions. The fixtures use abstract candidate labels and no Hecke characters. They omit automorphic representations, formal group hypotheses, ramification, conductors, global gluing and the analytic estimates. In particular A and B in the second fixture do not exhaust place-indexed functions. None of these fixtures has the evidence class `EXACT_COUNTEREXAMPLE`.

## Quantitative account

For p=resChar(w), on one shared record and at scale 1,

$$
s_g=\sum_{j=1}^m(g^D_{jw}+F_{jw}-E_{ujw}),\qquad
\mathrm{totS}=d+\sum_{j=1}^m(g^S_{pj}+F_{jw}-ES_{ujw}).
$$

Here d=vₚ(DGam(u)) is the additive valuation of the stabilizer factor. The error identification E=ES and gᴰ≤gˢ hold **stage by stage**. Subtraction gives a sum of nonnegative differences. The period conversion uses exactly e_w d of the available valuation allowance. Equality at every stage is legal. A prime dividing condN·g₀·g₁, a prime outside the eligible classes, or IsCs contributes zero served gain. Level and generator primes have dedicated period estimates. The class and conductor restrictions cannot be removed by treating all small primes as favorable. [Definitions](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/WeilJacquetShalikaLaw.lean#L37)

The arithmetic cost bound used later has error C·|TT|·a·La + C·m·(B₀+Nc)·x^(2ε). C must work for all sufficiently large x and indices of its family. Field degree enters other inputs, including the log-weighted prime sum with bound C·k·f·b′·log(X″(n)). The source also branches on logₚ(2a)<1 in the twisted-system comparison. At the family join the explicit lower threshold 3^(5/Fm.ez) can grow rapidly as ez shrinks. These are conditional controls with recorded parameter dependence. The sum identity alone bounds none of these costs. [Prime-sum bound](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/WeilJacquetShalikaLaw.lean#L127) · [Large-prime condition](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GrowthBounds/JointsTw.lean#L216) · [Accumulated height](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Formalism/SharpRestriction/GlobalTheoremInputs.lean#L77)

The deduction of the exponent also needs design existence, the rank and divisibility links, the remaining global comparisons, compatible primitive-character hosts and covered windows. Given these, the source fixes eD=2, eL=0 and obtains growth exponent 1/2+(eD+eL+7)ε. For every β>1/2 take 0<ε<min(ε₁,1/54,(β−1/2)/9). The exponent is below both 2/3 and β. Family and x-threshold may depend on ε. The source uses the window form of the analytic implication, then extends to imprimitive characters. The `CONDITIONAL` arithmetic deduction holds for arbitrarily small ε supplied by these hypotheses. A2 supplies no independent uniform realization as field, conductor and m vary. [Endpoint trace and hypotheses](a1-claude-grh-target.md)

## One A3 target

**Attack `GlobalTheory.TwU.LowF.lowAt_F`, the local coefficient bound used for the explicit Fn.** For each eligible μ,u,w and tied model, one local system Y5 must admit one table Fn with every required LowAt bound and

$$
 \mathrm{lowS}\ \le\ \sum_{j=1}^m\sum_{v\ \mathrm{outside\ representatives}} Fn(v,j)
                 -\sum_{j=0}^{m-1} dG(j).
$$

The source actually writes down Fn. For p below or at the size cutoff, its value is the orbit main term FOrb minus the error errF at the same second-orbit label. Above the cutoff it is zero. A3 should reconstruct the lower bound on order for each coefficient of the transformed integral input required by LowAt, with the ceiling on this real-valued difference. A failure would require an eligible second-orbit place, legal stage and genuine input of the original local model. The scalar-chain fixtures do not supply any of these. [Formula and use](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LocalSystemsLowerBoundsSecondRealPlace.lean#L53) · [Exact A3 theorem](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LocalModelBoundPlaceOverSecondEmbedding.lean#L66) · [Coefficient predicate](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Definitions/GaloisGeometry/IsPlaceCofactorValuationLowerBound.lean#L21)

At lines 83–95 the proposed A3 theorem aligns the rows of the Young diagram with the local system and invokes its single-place series estimate for every integral input series U. An exact attack can therefore fix p,u,i,t,n,r₀,r₁ and the integral parameter data first, then target the order of one specified coefficient. The stage argument is i+1 in both FOrb and errF. [Single-place reduction](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LocalModelBoundPlaceOverSecondEmbedding.lean#L83) · [Series inequality and hypotheses](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LocalModelBoundPlaceOverSecondEmbedding.lean#L38)

There is a short independent degeneracy check. At a served group index take the invariant test function identically 1 in the second-orbit sum identity. Its left side counts the complementary places and its right side counts reps. The latter is the image of the full finite group under the orbit representative map, hence contains the image of its identity. An empty complementary sum with a positive required lower bound is therefore unavailable on these hypotheses. This differs from the empty VV branch of the off-place theorem. [Orbit identity](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/PlacesOffFirstOrbits.lean#L30) · [Use on the selected system](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/Development/AbelianVarietyTheory/Proofs/SatakeParameters/LocalSystemsLowerBoundsSecondRealPlace.lean#L51)

## Reproduction and trust

Run the following from the repository root. The A2 graph expands A1's adm2 premise and preserves its downstream parent links. Its projections to separate existence are one-way. Substituting them for the joint premise fails rule matching, and the common selection, divisibility and RH applications remain uncertified. Mathematical reconstruction here is by hand. The archive replay checks the hash and source text. It does not run Lean or Comparator.

```powershell
go test ./audit/claudegrh -run 'TestSelection|TestJointWitnessGraph' -v -count=1
go run ./cmd/a1audit -a2-graph
go run ./cmd/a1audit -replay audit/claudegrh/testdata/averaging.json
go test ./...
go test -race ./...
go vet ./...
```

With `A1_LEAN_ARCHIVE` set to the original archive, run `go test ./audit/claudegrh -run TestPinnedSourceUses -count=1 -v`. A1 records the archive digest and download link. Run `oct test (Resolve-Path experiments/a2_simultaneous.octest).Path --execution compiled` and repeat with `interpreted`. These tests search the minimal one-place menu. Go separately checks both persisted certificates, including the first violated condition. [Graph](../audit/claudegrh/selection_graph.go) · [Exact replay tests](../audit/claudegrh/selection_test.go)

On Windows the full Go suite, race suite and vet passed. The hash-pinned source replay and the M25 report regression for J*=2−c* passed. Replaying the A1 fixture returned REDUCED_MODEL_COUNTEREXAMPLE with shifted value 1/2. Both A1 and A2 Oct experiments passed in interpreted and compiled modes, with zero compiled fallback using the oct dev binary fingerprinted in A1. The expanded graph reported six A1 edges and seven A2 edges, all untrusted. The two-case fixture has SHA-256 `1ff0f067c793683d94880754bb89ceef69f35239188f1fe7c9306cbe81717574` for the audited working-copy bytes.
