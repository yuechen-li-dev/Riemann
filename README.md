# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M13 leaves M12's finite rank-trace theorem
unchanged and types the remaining analytic representation choice: the legal
Montgomery-Taylor window, its `0 < lambda <= 1` support domain, the
`G_tilde -> G_hat` normalization, and the general-window moment coefficients.
The compiler derives

```text
J(lambda) = 2 - lambda/2 - (1/sqrt(2))*cot(lambda/sqrt(2))
```

from those inputs. Its exact derivative is positive throughout the legal
domain, so `lambda=1` is the global maximizer. A rational Taylor enclosure
proves the simple-critical constant is greater than `269/400`, safely displayed
as `67.25%`; the all-distinct bound is greater than `669/800` (`83.625%`). RH
remains unresolved.

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
[docs/m13-architecture.md](docs/m13-architecture.md) for the exact window
functional, typed normalization, Oct scan/plot, global optimization proof,
rational certification, and M8-M13 provenance. Earlier milestone documents
remain available.
