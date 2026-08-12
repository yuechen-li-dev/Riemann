# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M8 lowers the exact Weil zero-side summand
through the existing two-function Hermitian basis, canonicalizes critical and
off-critical zero orbits, and derives their symbolic low-rank matrix blocks.
The same semantic matrix retains M7's certified explicit-formula intervals as
an alternate theorem-linked representation. Finite positivity and the inherited
function-space coverage loss still leave RH unresolved.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, then decomposes its zero-side meaning into symmetry-orbit contributions.
It rejects approximate, finite-coverage, and finite-PSD reverse-inference
shortcuts. See [docs/m8-architecture.md](docs/m8-architecture.md) for the exact
coordinate derivation, local classifications, experiment, evidence boundaries,
and findings. Earlier milestone documents remain available.
