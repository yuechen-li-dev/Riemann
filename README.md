# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M7 evaluates the existing certified
two-function Weil basis through Lagarias's Mellin explicit formula with
theorem-backed directed intervals. Exact rational endpoints, certified finite
prime sums, piecewise Archimedean enclosures, and analytic tails retain separate
proof objects. The finite matrix and finite-span positivity are certified; the
inherited finite function-space coverage loss keeps RH unresolved.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, and rejects approximate PSD and finite-coverage shortcuts. See
[docs/m7-architecture.md](docs/m7-architecture.md) for the exact integral,
certified evaluator, interval matrix theorem, evidence boundaries, and findings.
Earlier milestone documents remain available.
