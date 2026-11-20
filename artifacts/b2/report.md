# B2 · A standalone cyclic-duality demonstration

| Evidence | Observed result |
|---|---|
| Scientific exploration | Oct solved 435 cyclic packing problems at lengths 2, 4 and 6. Interpreted and compiled exports were byte-identical. |
| Exact finite certification | Go replayed all 435 rows and the original B1 set of 405. Feasible primal and dual solutions had exactly equal objectives in every instance. |
| Universal formal verification | Both pinned Lean modules compiled in isolation. Type comparisons passed, and the official kernel replay accepted both compiled modules. |
| Axiom footprint | Both audited theorems depend on `propext`, `Classical.choice` and `Quot.sound`. |
| Mathematical contribution | A reproducible demonstration of established duality. A sparse six-vertex case reaches throughput one using offset five, while offset-one flow alone reaches at most one third. |

[Mathematics and visualization](../../docs/b2-cyclic-duality.md) · [Commands, toolchain and limitations](../../docs/b2-verification.md) · [Lean output](lean-transcript.txt) · [Exact certificates](certificates.json) · [Compiler provenance](proof-graph.json)
