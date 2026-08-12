# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M9 combines M8's zero-orbit decomposition
`G=P+Q` with exact spectral information available from the explicit-formula
representation. It derives the reusable finite theorem
`rank(P) >= max(0,n_plus(G)-B_off)`, conservatively transports that to a
critical-contribution count, and records why the asymptotic one-half stage still
needs height-window moment and tail inputs. Finite positivity and the inherited
function-space coverage loss still leave RH unresolved.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, decomposes its zero-side meaning into symmetry-orbit contributions, and
accounts for their aggregate rank and inertia budgets. It rejects approximate,
finite-coverage, additive-inertia, multiplicity-as-rank, and finite-PSD reverse
shortcuts. See [docs/m9-architecture.md](docs/m9-architecture.md) for the exact
finite theorem, counterexamples, experiment, provenance fusion, and isolated
asymptotic obstruction. Earlier milestone documents remain available.
