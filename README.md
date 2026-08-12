# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M5 lowers a certified finite-dimensional
restriction of the Weil quadratic functional through complex polarization to
a typed Hermitian form and ordered-basis matrix. It records the distinction
between finite families and finite spans, the identity `Q(sum c_i f_i)=c*Gc`,
entrywise explicit-formula provenance, exact versus numerical matrix evidence,
and the inherited finite function-space coverage loss.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs its structural Hermitian matrix, and
rejects basis-point, diagonal, approximate, and finite-coverage shortcuts.
See [docs/m5-architecture.md](docs/m5-architecture.md) for the convention,
theorem contracts, trust boundary, Oct experiment, and findings. Earlier
milestone documents remain available.
