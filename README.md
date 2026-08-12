# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M4 lowers RH from typed zero geometry to the
equivalent universal positivity of the Weil quadratic functional on a typed
test-function class. It records the Mellin convention, zero/prime/archimedean
explicit-formula structure, certified finite-family restriction, coverage
loss, and the numerical evidence boundary.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite family, and rejects both finite and numerical positivity
as proofs of RH. See [docs/m4-architecture.md](docs/m4-architecture.md) for the
criterion, normalization, trust boundary, IR, Oct experiment, and findings.
Earlier milestone documents remain available.
