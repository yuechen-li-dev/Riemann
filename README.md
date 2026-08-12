# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M11 adds typed trace, Frobenius-square, and
dimension claims; proves the generic finite bound
`n_plus^theta(G) >= ceil((A-d*theta)^2/B)` under its explicit positivity and
evidence premises; and instantiates it with the paper's first two prime-side
moments. At bandwidth one this yields
`n_plus^theta(G_tilde_T) >= (3/4-o(1))N(T,2T)`. Reusing M10 then reconstructs
the earlier half-stage result for simple critical zeros and the three-quarter
result for distinct zeros. The later rank-trace optimization is deliberately
absent, and RH remains unresolved.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, decomposes its zero-side meaning into symmetry-orbit contributions,
accounts for aggregate rank/inertia budgets, and localizes a separate
height-dependent compression. It rejects approximate, non-strict-threshold,
wrong-norm, multiplicity-as-rank, and finite-PSD reverse shortcuts. See
[docs/m11-architecture.md](docs/m11-architecture.md) for the exact moment
theorem, source normalization, threshold scaling, counterexamples, experiment,
and M10 composition. Earlier milestone documents remain available.
