# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M17 characterizes the unsaturated one-radius
completion for the broader support-one Fourier extremal class. It proves that
origin saturation retains M15's `9/8` ceiling, then escapes it with positive
origin slack. It certifies `c_*>=2297/2000`, hence `J_*<=1703/2000`, and
combines this with the rigorous CGdL feasible bound to give
`849/1250<J_*<=1703/2000`. Anthropic's remark-level
`0.68185` remains unresolved and is not promoted to theorem evidence. The
compiler still retains M13's exact
Montgomery–Taylor derivation

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
go test -race ./...

# Replay the expensive exact M17 whole-line enclosure explicitly.
go test -tags=slow ./compiler -run M17Slow -count=1
```

The normal unit and race suites validate the M17 proof object's exact
parameters, structural invariants, rejection paths, deterministic artifact,
and bound propagation. The `slow` suite additionally replays the 20,001-cell
big-rational whole-line certificate. Report serialization never reruns that
proof. Compiler test fixtures also build the M6–M15 chain once and pass each
stage to the next; M7 receives an independent graph clone because it extends
the proof graph.

The demonstration lowers RH to universal Weil positivity, restricts that claim
to a certified finite span, constructs and evaluates its structural Hermitian
matrix, decomposes its zero-side meaning into symmetry-orbit contributions,
accounts for aggregate rank/inertia budgets, and localizes a separate
height-dependent compression. It rejects approximate, non-strict-threshold,
wrong-norm, multiplicity-as-rank, and finite-PSD reverse shortcuts. See
[docs/m15-architecture.md](docs/m15-architecture.md) for the certified
whole-line completion, exact family ceiling, Oct experiments, and plotting
artifact repair. [docs/m17-architecture.md](docs/m17-architecture.md) records
the unsaturated one-radius result and its slow proof boundary.
[docs/m14a-architecture.md](docs/m14a-architecture.md) retains
the primal/dual reconstruction. [docs/m13-architecture.md](docs/m13-architecture.md)
retains the exact window functional and M8-M13 provenance. Earlier milestone
documents remain available.
