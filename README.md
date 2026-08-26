# Riemann

Riemann is an experimental compiler for making transformations of mathematical
claims mechanically inspectable. M24 proves that every legal point in M23's
contact-preserving two-atom radius component outside `r in [3.503,3.512]` is
strictly below the M23 optimizer's certified lower objective. It finds no
stronger branch but honestly leaves that exact radius/contact window unresolved
(Success C). The best whole-line theorem witness remains `r=7/2, m=69/10000`,
giving `c_*>=1+pi/4-2/pi+69*pi*sqrt(2)/40000` and the complementary upper bound
on `J_*`. M18 closes the unsaturated one-radius
completion for the broader support-one Fourier extremal class. It derives the
exact tangency weight `1/pi-4/pi^2+8/pi^3`, certifies the resulting transform
on the whole real line, and composes that witness with M17's upper envelope:
`C_1R=1+pi/4-2/pi`. At M18 this gave
`1+pi/4-2/pi<=c_*<1651/1250` and
`849/1250<J_*<=1-pi/4+2/pi`. Anthropic's remark-level
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

# Replay the exact M18 endpoint and inherited M17 whole-line enclosures.
go test -tags=slow ./compiler -run M18Slow -count=1

# Replay the M21 rational Taylor cover for the contact-adapted ray.
go test -tags=slow ./semantic -run M21Slow -count=1

# Replay the M22 algebraic/rational Taylor cover for r=7/2.
go test -tags=slow ./semantic -run M22Slow -count=1

# Replay the M23 improved explicit whole-line witness.
go test -tags=slow ./semantic -run M23Slow -count=1

# Replay M24's certified component complement.
go test ./semantic -run M24 -count=1
```

The normal unit and race suites validate the M18 proof object's exact
parameters and bound propagation plus M20's perturbation, M21's exact ray, and
M22's variable-radius contact surface,
evidence-separation, primal-guard, and deterministic-artifact invariants. The `slow` suite
additionally replays the inherited
M17 rational cover and M18's equality-aware endpoint cover. M20 then uses the
exact endpoint as a regression oracle, derives its contact first variation, and
selects a one-parameter phase-aligned exterior perturbation without changing
the theorem stack. M21 preserves that fixed contact, certifies a positive
rational mass, and keeps numerical candidates separate from certified
witnesses and ceiling brackets. M22 then releases only the radius, certifies an
algebraic half-integer witness, and retains the two-atom support class. Report
serialization never reruns those proofs. Compiler test fixtures also build the M6–M15 chain once and pass each
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
artifact repair. [docs/m18-architecture.md](docs/m18-architecture.md) records
the exact tangency theorem, equality-aware whole-line proof, and Oct artifact.
[docs/m17-architecture.md](docs/m17-architecture.md) records the preceding
unsaturated one-radius bracket and its slow proof boundary.
[docs/m19-aetheris-mathematics-audit.md](docs/m19-aetheris-mathematics-audit.md)
audits Aetheris Preview 2 as a possible mathematical geometry runtime; its
Firmament syntax findings are recorded in the companion
[parser-backed inventory](docs/m19-firmament-syntax-audit.md) and
[language snapshot](docs/m19-firmament-language-snapshot.md). M19 introduces
no production integration and makes no theorem claim from approximate geometry.
[M21's contact-adapted ray report](docs/m21-contact-adapted-ray.md) records the
exact response, scalar envelope, second numerical contact, rational whole-line
witness, certified ceiling rejection, and Oct artifacts.
[M22's variable-radius two-atom report](docs/m22-contact-preserving-two-atom.md)
derives the exact contact surface, first-order release, nested radius envelope,
simple algebraic witness, whole-line proof, and deterministic Oct artifacts.
[M23's certified local-optimizer report](docs/m23-variable-radius-local-optimum.md)
boxes the three-equation root, proves strict local radius optimality, and keeps
the improved explicit whole-line witness separate from optimizer evidence.
[M24's component-control report](docs/m24-two-atom-component-control.md)
uses Oct-first branch discovery, a fixed-contact large-radius theorem, and
generated adverse-frequency interval cells to isolate the only unresolved
radius/contact window without claiming globality.
[M20's representation-selection report](docs/next-representation-selection.md)
derives the contact sensitivity, compares four minimal extensions, and selects
one exact one-parameter family for M21 without changing the certified bound.
[docs/m14a-architecture.md](docs/m14a-architecture.md) retains
the primal/dual reconstruction. [docs/m13-architecture.md](docs/m13-architecture.md)
retains the exact window functional and M8-M13 provenance. Earlier milestone
documents remain available.
