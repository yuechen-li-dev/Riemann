# B3 · Verification and replay

**Verified.** The independently written Chebotarev union implication compiled in a Mathlib-only project. Lean also replayed the original nine-module degree-one branch and the two selected analytic comparison results. The bounded original single-class build stopped with concrete dependency and resource evidence. [Mathematics and statement mapping](b3-chebotarev.md)

## Pins and proof evidence

| Input | Verified pin |
|---|---|
| Submitted corpus | `a172c464fd4bcfc4f53e394dcba45d742b1409a6` |
| Codeload archive SHA-256 | `2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75` |
| Official Lean | `4.33.1`, commit `819816b2e0a3bf405af45ae5c7af2491d8f5bee6` |
| Mathlib | `0df444a360eaa60ab8c11dca51a86af692955474` |
| Conditional module SHA-256 | `24bfe23e17b36dd86a3abb71448a07efa11a9559a7160cb08db63b560b15e734` |

**Conditional module.** `B3.densityFromCentralizers` accepts exactly the source's reciprocal-centralizer density normalization for single classes and concludes the `Nat.card C / Nat.card G` density. `B3.classWeight` checked the normalization conversion by orbit-stabilizer. `B3.densityAllAbove` checks the change from “every place over this prime” to a class-tagged set under two explicit place hypotheses. The proof uses the exact Mathlib density predicate. Substituting (G=operatorname{Gal}(L/K)), the relation `LiesOver`, the predicate (e=1), and the submitted arithmetic Frobenius gives the mathematical source interface. The generic module leaves these objects as parameters. Source comparison is a header and predicate mapping because the original Chebotarev module was unavailable to import. [Standalone Lean source](../scripts/B3.lean) · [Complete checked types](../artifacts/b3/conditional-transcript.txt)

**Original degree-one branch.** A distinct scratch project built the nine submitted modules needed by `DegreeOneDensity`. The proof checker replayed the declarations of `DegreeOneDensity` from its compiled object. The probe checked the original `hasDirichletDensity_inter_degreeOnePrimes_iff` and `tendsto_primeIdealZetaSum_univ_div_univ`. There was no edit to the submitted source. The probe signatures are in the output. [Complete checked types](../artifacts/b3/degree-one-transcript.txt)

For all seven named standalone theorems and both original analytic comparisons, `#print axioms` reported only `propext`, `Classical.choice`, and `Quot.sound`. Both projects returned exit code **0** after compilation, probes, the axiom gate and separate official `leanchecker` processes. The replay processes use the Lean kernel and the pinned imported objects. B3 reused B2’s Mathlib cache and checked the pinned Git revision. Earlier conditional drafts failed elaboration before the final successful run.

## Bounded original attempt

The authenticated authored import closure contains **351** modules, with Mathlib as its only external import root. The transitive counts for `ChebotarevAbelian`, `ChebotarevCount`, and `DegreeOneDensity` are 346, 57, and 9, respectively. These are sizes of source import closures, with overlap. The manifest records each module path, SHA-256, direct imports and the four requested theorem headers. This is the module-granularity closure for building the unmodified source. It is not a declaration-minimal dependency slice. [Closure manifest](../artifacts/b3/imports.json)

The first observed hard failure was `GroundTheory/LocalFields/Interface.lean:35:36`. Lean exhausted the default 20,000 typeclass heartbeats while synthesizing `LinearOrder (Submodule ↥𝒪[K] ↥𝒪[K])` for the `FaithfulSMul` example. A focused rerun of that file with `-DsynthInstance.maxHeartbeats=400000` exited 0. The 360-second original build still ended with exit **124** before any of the submitted Chebotarev targets could be probed. The retry diagnoses a build-budget problem at that site. It supplies no evidence about the full proof’s correctness. [Failure and recovery output](../artifacts/b3/original-attempt.txt)

Two paths explain why the source hits that file:

- `ChebotarevAbelian → GlobalReciprocity.Frobenius → ExistsUnitFiltrationUnitsLeRangeNormUnits → LocalFields.Interface`.
- `ChebotarevCount → ChebotarevFrobenius → SplittingFromInequality.Decomposition → UnramifiedGalois → LocalFields.Interface`.

The remaining analytic supplier is `GroundTheory.Global.hasDirichletDensity_setOf_arithFrobAt_eq` in `ChebotarevAbelian.lean`. The original Deuring proof applies it to the cyclic extension (L/L^{\langle\sigma\rangle}). A source-level replay also needs the separate Deuring fibre count and conjugacy interfaces. B3 stopped at the requested boundary with the conditional generalization and the nine-module analytic branch replayed. [Supplier and its ray-class call](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/src/GroundTheory/GlobalFields/TateThesis/ChebotarevAbelian.lean#L90)

## Reproduction

Use a Linux or WSL shell with the official Lean 4.33.1 binaries on `PATH`, Python 3.11 or later and a pinned archive outside the repository. Each script call takes a fresh scratch directory. The optional last argument points to a pinned B2 scratch project for reuse of its Mathlib packages. Without it the script fetches the Mathlib cache.

~~~sh
bash scripts/b3-lean.sh /path/to/pinned.tar.gz /new/b3-conditional conditional
bash scripts/b3-lean.sh /path/to/pinned.tar.gz /new/b3-degree degree-one
B3_BUILD_SECONDS=360 bash scripts/b3-lean.sh /path/to/pinned.tar.gz /new/b3-original original
python3 scripts/b3-source.py /path/to/pinned.tar.gz artifacts/b3/imports.json
~~~

The scripts preserve `verification.log`, `checked.txt`, `kernel-replay.txt` and `exit-code.txt` where reached. The transcripts keep the status lines, complete check outputs and the observed failure. The corpus artifacts contain source headers and provenance. [Pinned archive](https://codeload.github.com/Deicyde/riemann-hypothesis/tar.gz/a172c464fd4bcfc4f53e394dcba45d742b1409a6)

~~~powershell
$env:OCT_WRAPPER_PATH = 'C:/path/to/Oct/sidecars'
go run ./cmd/b3chebotarev -oct 'C:/path/to/oct.exe' -archive 'C:/path/to/pinned.tar.gz'
go test ./byproducts -run TestB3 -count=1
python scripts/b3-plot.py
go test ./...
go test -race ./...
go vet ./...
~~~

The CLI runs Oct in separate scratch directories, compares interpreted and compiled output byte for byte, then checks the group identity, associativity, inverses, exhaustive class partition, closure under conjugation, orbit-stabilizer sizes, and all class masks. It recomputes weights using `math/big.Rat`. B3 shares the Oct export runner factored from B2. The optional source archive flag replays the B1 typed Chebotarev record with the existing Go authenticator. Both B2 and B3 artifacts reproduced byte for byte. B3 ran in both modes with zero fallback. Full Go tests, race tests, vet, the pinned B1/A1 source checks, M25 normalization, and A1–A4 Oct regressions passed. Oct reported `oct dev`, executable SHA-256 `c5c809641d5c913980b7552af0187d31dd0160741c18e68c266b99730788517b`. [Finite checker](../byproducts/chebotarev.go) · [Corruption tests](../byproducts/chebotarev_test.go) · [Oct search](../experiments/b3_chebotarev.octest)
