# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M6 evaluates a certified two-function Weil
basis through Lagarias's Mellin explicit formula. Exact definitions, exact
rational endpoint terms, support-exhaustive but floating-point prime sums,
heuristically quadrature-bounded archimedean terms, and approximate totals have
separate typed evidence states. The inherited finite function-space coverage
loss remains explicit.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, and rejects approximate PSD and finite-coverage shortcuts. See
[docs/m6-architecture.md](docs/m6-architecture.md) for the basis, formulas,
evidence lattice, evaluator metadata, Oct experiment, and findings. Earlier
milestone documents remain available.
