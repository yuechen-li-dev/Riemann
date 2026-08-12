# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M2 represents imported mathematics as typed
theorem contracts, unifies their parameters structurally, exposes unresolved
premises, and composes them deterministically. The M1 Euler-product zero-free
half-plane now emerges from generic contract instantiation rather than a custom
Euler derivation pass.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go run ./cmd/riemann --missing-premise
go test ./...
```

The demonstration also rejects three invalid attempts to prove RH: bounded
zero verification, density-one information, and zero-exclusion on `Re(s)>1`.
See [docs/m2-architecture.md](docs/m2-architecture.md) for theorem schemas,
matching, obligations, composition, trust, citations, and M2 findings. The M1
and M0 documents record the earlier designs.
