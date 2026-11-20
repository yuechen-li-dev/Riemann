# B2 · Verification and replay

**Success A and B.** Both original authored modules compiled in isolation at the pinned versions. Lean confirmed the exact proposition types and axiom footprint. Its official `leanchecker` reloaded the two generated `.olean` modules and replayed their declarations in a separate process. The finite Oct exports agreed byte for byte, and Go independently certified every row. [Lean transcript](../artifacts/b2/lean-transcript.txt) · [Mathematical statements and examples](b2-cyclic-duality.md)

## Formal verification

| Input | Verified pin |
|---|---|
| Corpus | `a172c464fd4bcfc4f53e394dcba45d742b1409a6` |
| Codeload archive SHA-256 | `2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75` |
| `IncomingBound.lean` SHA-256 | `201476c18ccbb46a0022e76b63901c7feea2a565b8f9802910d9c36285791ba6` |
| `FlowDegreeBounds/MaxFlowValueAttained.lean` SHA-256 | `18d881da800499e504baa5eb485c4cc800837fa4f14a87fa217b504ea4a903b9` |
| Official Linux Lean | `v4.33.1`, commit `819816b2e0a3bf405af45ae5c7af2491d8f5bee6` |
| Official Linux release archive SHA-256 | `890afd185370f85666025b883914ab4f4b339136f8c96167b69cfb62aecaf235` |
| Mathlib | `0df444a360eaa60ab8c11dca51a86af692955474` |

The project extracts only the two byte-checked source files into a new scratch directory. Its authored proof module imports Mathlib and the authored statement module. The probe checks `FlowAttained : Prop` and `LemmaFa : Prop` and compares them against independently written interfaces using definitional equality. It then type-checks `maxFlowValueAttained : FlowAttained` and `strongDuality : LemmaFa`. The transcript includes both complete check-printed types. [Script](../scripts/b2-lean.sh) · [Official Lean release](https://github.com/leanprover/lean4/releases/tag/v4.33.1) · [Mathlib revision](https://github.com/leanprover-community/mathlib4/tree/0df444a360eaa60ab8c11dca51a86af692955474)

Both `#print axioms` commands reported exactly `[propext, Classical.choice, Quot.sound]`, which passed the script’s explicit allowlist. Lake, the probe, the axiom filter and the subsequent replay completed with exit code 0. The supplier emitted linter and deprecation warnings for unused section variables, `push_neg` and `Set.mem_setOf_eq`. The probe ran with warnings treated as errors. An initial WSL launcher command failed because the inherited Windows path was expanded into shell syntax. Passing the pinned Linux path with `wsl ... -- env PATH=...` resolved it. [Selected build log and complete probe output](../artifacts/b2/lean-transcript.txt)

The replay uses the official Lean kernel with imported Mathlib objects. It is a separate check after elaboration and uses the same kernel implementation. The run used Mathlib’s cache and did not replay its entire transitive closure from a fresh environment or use a second kernel implementation. The B1 alias in `WeightIdentities.lean` lies outside this two-module check. [Official checker implementation](https://github.com/leanprover/lean4/blob/819816b2e0a3bf405af45ae5c7af2491d8f5bee6/src/LeanChecker.lean)

Run the formal verification from Linux or WSL with Lean 4.33.1 and its bundled Lake and leanchecker on `PATH`:

~~~sh
bash scripts/b2-lean.sh /path/to/riemann-a1-a172c464.tar.gz /path/to/new-scratch-project
~~~

The script writes a complete `verification.log`, `checked.txt`, `kernel-replay.txt`, `sources.sha256`, and `exit-code.txt` inside the scratch project. The committed transcript omits cache progress and keeps the complete probe output. The original Lean files and build products stay in scratch space. [Pinned archive](https://codeload.github.com/Deicyde/riemann-hypothesis/tar.gz/a172c464fd4bcfc4f53e394dcba45d742b1409a6)

## Computational demonstration

~~~powershell
$env:OCT_WRAPPER_PATH = 'C:/path/to/Oct/sidecars'
go run ./cmd/b2duality -oct 'C:/path/to/oct.exe'
go test ./byproducts -run 'TestB2|TestTwoLetterCertificates' -count=1
~~~

The CLI runs the same Oct test in interpreted and compiled modes in separate scratch directories and compares the exports byte for byte. A compiled run initially failed to find `octxiliary-io`. Supplying its directory via `OCT_WRAPPER_PATH` resolved the failure. Both final modes passed with zero fallback. The Oct executable identifies as `oct dev`, SHA-256 `c5c809641d5c913980b7552af0187d31dd0160741c18e68c266b99730788517b`. [Oct experiment](../experiments/b2_cyclic_duality.octest)

Oct’s fraction-free simplex starts at the slack basis, chooses pivots deterministically by Bland’s rule, and exports integer numerators and a common denominator. An inexact division, failed pivot or exhausted pivot budget yields an unsolved row. The fixed-width Oct arithmetic is scoped to the bounded search in the experiment. The rational Go verifier computes independently from `l,d,D,x,σ`. It enumerates every directed odd offset, checks variables and loads, checks every covering inequality and the dual box, and recomputes equality of objectives. The 435 cases comprise 405 B1 parameter combinations, 26 four- and six-vertex sweep points, two sparse long-arc cases, one balanced six-vertex case and one precision case above (2^{53}). The regression also runs the original B1 certificate generator and checker. Tests corrupt primal nonnegativity, capacity feasibility, dual feasibility, dual bounds, objective equality and source parameters. [Exact verifier](../byproducts/cyclic.go) · [Adversarial replay tests](../byproducts/cyclic_test.go) · [Certificates](../artifacts/b2/certificates.json)

The proof graph records primal and dual feasibility as parallel checks of a fixed problem. The joint certificate parents the finite optimality claim. The universal theorem, source proof and formal-status claims have a separate provenance chain. The serialized graph leaves formal import open in Riemann, and the Lean transcript records the observed replay independently. The compiler’s graph guard rejects an attempted finite-to-universal inference. [Graph integration and provenance](../byproducts/cyclic_graph.go) · [Serialized graph](../artifacts/b2/proof-graph.json)

Regenerate the plot using Python with Matplotlib. The script reads the JSON and converts verified rationals to display coordinates. [Plot generator](../scripts/b2-plot.py)

~~~sh
python scripts/b2-plot.py
go test ./...
go test -race ./...
go vet ./...
~~~

**B3.** Isolate the Chebotarev conjugacy-stable-set theorem. Audit its analytic closure first, then replay the single-class supplier and the finite disjoint-union step with separate axiom reports. [B1 Chebotarev comparison](b1-byproduct-inventory.md#2-chebotarev-density-for-conjugacy-stable-sets)
