# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M1 gives quantifiers and domains typed
semantics, models two representations of the Riemann zeta function, and derives
its zero-free half-plane from an Euler-product representation plus explicit
trusted analytic premises.

```powershell
go run ./cmd/riemann
go run ./cmd/riemann --json
go test ./...
```

The demonstration also rejects three invalid attempts to prove RH: bounded
zero verification, density-one information, and zero-exclusion on `Re(s)>1`.
See [docs/m1-architecture.md](docs/m1-architecture.md) for the semantics, trust
boundary, and M1 research findings. [docs/m0-architecture.md](docs/m0-architecture.md)
records the earlier provisional design.
