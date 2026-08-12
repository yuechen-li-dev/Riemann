# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M0 is a deliberately small Go vertical slice:
it represents the Riemann Hypothesis, normalizes it to the universal
zero-location statement, lowers that statement to a density-one consequence,
and rejects the invalid reverse inference.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go test ./...
```

The human report explains the lowering and rejection. `--json` emits the same
claim/transformation graph in a deterministic machine-readable form. See
[docs/m0-architecture.md](docs/m0-architecture.md) for the design and research
findings.
