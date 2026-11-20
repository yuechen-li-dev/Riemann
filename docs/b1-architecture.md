# B1 · Reproducing the extraction pilot

B1 has a separate evidence boundary. Riemann’s compiler graph and typed theorem contracts apply to established claims in its supported language. The independent `ExternalStatement` records in `byproducts/` never enter that graph. B1 follows M24–M26’s separation of discovery and exact replay, and A1–A4’s digest-pinned source checks. [Record types](../byproducts/statements.go) · [Compiler graph](../compiler/graph.go) · [Compiler contracts](../compiler/theorem.go) · [Semantic types](../semantic/types.go) · [M26 replay architecture](m26-release-forced-contact.md)

| Stage | Input and output | Evidence |
|---|---|---|
| Bounded human shortlist | Five declaration anchors, source spans, interpretations and comparison findings | Review of specific statements and prerequisites |
| Go extraction | Pinned tar.gz and shortlist to typed records | Archive hash, source and scope hashes, name anchors, cited definition mentions, direct module paths and selected proof-body references |
| Oct discovery | Small feature manifest to ordered rows and integer scores | Ranking only |
| Go replay and reporting | Exact arithmetic certificates and validated rank manifest to JSON and Markdown | Exact finite results for the two-letter LP, plus reproducible source records |
| Lean replay | Pinned Lean, Mathlib and required import closure | Future evidence for the specific declarations checked |

The streaming reader hashes the compressed archive and retains text only for the selected files. It walks archive headers to resolve file paths. It does not index every declaration or export a corpus database. The catalog binds each human review to exact source spans, including section context where needed. The stable ID includes the commit, qualified name and SHA-256 of the exact statement and reviewed source context. It is a versioned source identity, not a semantic equivalence key. The report includes the Oct program digest with normalized line endings. [Catalog](../byproducts/b1_candidates.json) · [Extraction and provenance](../byproducts/provenance.go)

Oct adds two points each for a named proposition and Mathlib-only definition files, one for a literature comparison, five for an exact finite replay and two for unresolved formalization comparisons. An additional theorem-level proof premise receives a three-point penalty. The scorer also penalizes visible RH mentions in exploratory rows. The five admitted records have passed Go’s source-level RH mention guard, so that last feature is exercised by the Oct test. Ties retain catalog order. The reader checks that the manifest returns each input exactly once with integer scores in descending order. Neither a high score nor a close literature match changes the verification field. [Oct program](../experiments/b1_discovery/main.oct) · [Runner and rank validation](../byproducts/discovery.go)

Go checks source spans and lexical uses. It cannot reconstruct elaboration, find all effective typeclass or notation dependencies, check arbitrary aliases hiding RH, settle semantic equivalence, or determine a transitive axiom footprint. Human review expands the selected propositions and records ordinary hypotheses independently of the proof environment. All five records carry `no_rh_in_reviewed_statement`, which is scoped to the statement review. Unknown proof dependencies remain open. Novelty is tracked separately for the mathematics and the formalization. The B1 validator rejects proof replay and positive novelty labels, including when they are supplied by an input catalog. [Evidence checks](../byproducts/byproducts_test.go)

## Reproduce the inventory

Use Go 1.23 or later and an Oct CLI with file-line I/O. Keep the corpus archive outside the repository. Fetch the pinned codeload archive and pass its path to the CLI. B1 requires its original gzip bytes and verifies SHA-256 `2757f699bc40beff25a67722696f55be68009e00ce6b27ad6fdc95a79df38d75`. The outputs contain only theorem headers, source anchors and independently authored review text. [Pinned source archive](https://codeload.github.com/Deicyde/riemann-hypothesis/tar.gz/a172c464fd4bcfc4f53e394dcba45d742b1409a6)

~~~powershell
$archive = 'C:/path/to/pinned-corpus.tar.gz'
$oct = 'C:/path/to/oct.exe'
go run ./cmd/riemann-byproducts -archive $archive -oct $oct -out artifacts/b1

$env:A1_LEAN_ARCHIVE = $archive
$env:B1_OCT = $oct
go test ./byproducts -run TestB1PilotReplay -count=1 -v
go test ./audit/claudegrh -run TestPinnedSourceUses -count=1 -v

& $oct test experiments/b1_discovery --execution interpreted
& $oct test experiments/b1_discovery --execution compiled
go test ./...
go test -race ./...
go vet ./...
~~~

The replay test reruns the actual Oct CLI and compares both reports to the checked-in artifacts. It normalizes Git checkout line endings in the expected files. The normal Go suite also checks unsupported evidence promotions, bad source identity, missing imports or references, invalid ranking output and the rational LP certificates. The Oct build used for B1 identifies as `oct dev`, with executable SHA-256 `c5c809641d5c913980b7552af0187d31dd0160741c18e68c266b99730788517b`. Both execution modes passed, and compiled execution had zero interpreted fallback. [Expected JSON](../artifacts/b1/inventory.json) · [Expected Markdown](../artifacts/b1/inventory.md)

## Formal verification status

The archive supplies Lean `v4.33.1` and Mathlib `0df444a360eaa60ab8c11dca51a86af692955474`. The current host returned a command-not-found for Lean, and the available Arch WSL environment yielded no Lean or Lake executable. The isolated check therefore stopped at toolchain availability. Installing the pinned Windows toolchain would require an 842 MB release archive, plus the pinned Mathlib environment for the umbrella `import Mathlib`. No `#print axioms` result or kernel replay was recorded.

The narrow next replay needs `Formalism/IncomingBounds/IncomingBound.lean` and `Formalism/IncomingBounds/FlowDegreeBounds/MaxFlowValueAttained.lean`, both at the submitted source revision. Build those modules in a temporary project at the pinned Mathlib revision. Run `#print axioms Formalism.IncomingBounds.strongDuality` and compare its type to `LemmaFa`. The corpus’s `WeightIdentities` wrapper imports several further modules, so the two-module check would apply to the supplier and any locally reconstructed alias. For Hironaka, Chebotarev, Brauer and Krasner, the exact direct proof imports are listed in the JSON, with additional transitive imports required. There was no full-corpus build. [Pinned toolchain](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/lean-toolchain) · [Pinned dependency manifest](https://github.com/Deicyde/riemann-hypothesis/blob/a172c464fd4bcfc4f53e394dcba45d742b1409a6/lake-manifest.json) · [Results and B2](b1-byproduct-inventory.md)
