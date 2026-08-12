# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M10 adds the paper-backed height family
`G_tilde_T=P_near+Q_near+E_far`, controls the far term in operator norm, and
uses strict thresholded inertia plus M9 accounting to compile a certified local
spectral observation into finite simple-critical and distinct-zero counts in
`(T,2T]`. The first/second-moment input needed to produce a strong thresholded
spectral observation remains deliberately open. Finite positivity and the
inherited function-space coverage loss still leave RH unresolved.

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
[docs/m10-architecture.md](docs/m10-architecture.md) for the exact window
theorem, tail estimate, counterexamples, experiment, and remaining spectral
input. Earlier milestone documents remain available.
