# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M12 preserves the structural identity
`G=P+Q`, with `P` positive semidefinite and low-rank and `Q` controlled by its
positive index, and reconstructs the parameterized rank-trace inequality
`||P+Q||_F^2 >= c tr(P)-(c^2/4)r+2c tr(Q)-c^2b` for `c>0`.  Fusing that finite
theorem with M8-M11's orbit, window, and moment IR reproduces the bandwidth-one
lower bounds `2/3` for simple critical-line zeros and distinct critical-line
zeros, and `5/6` for all distinct zeros. The Montgomery-Taylor 67.25% window
optimization is deliberately absent, and RH remains unresolved.

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
[docs/m12-architecture.md](docs/m12-architecture.md) for the finite theorem,
source normalization, coefficient counterexamples, Oct experiment, and the
M8-M11 representation fusion. Earlier milestone documents remain available.
