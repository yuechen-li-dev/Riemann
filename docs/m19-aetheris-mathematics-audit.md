# M19 — Aetheris/Firmament as a mathematical research runtime

## Verdict

**Promising research IR, usable now as a visualization backend, not a theorem-adjacent verifier.**

Aetheris accidentally contains part of a general geometric substrate: authored `ParametricSurfaceIr`, construction profiles, stable semantic identities, deterministic materialization, and structured composition are useful outside CAD. The useful mathematical representation is above BRep—principally Panel/parametric IR and parts of Construction AIR. Feature AIR is CAD intent, BRepPlan is bounded topology planning, BRep is tolerance-defined geometry, and STEP is interchange. Exact expressions, global sign, contact order, and proof evidence do not survive those lowerings.

No production dependency was added in either direction. M19 adds only experimental fixtures/probe code, generated artifacts, and reports in Riemann. Aetheris production code is unchanged.

## Evidence and experiment ledger

| Source | Parse / validation | Lowering route | BRep | STEP / trace |
|---|---|---|---|---|
| `experiments/m19/firmament/m18-safe.firmament` | Accepted by V2 Panel parser and CLI `validate` | `PanelIr` -> adaptive panel materializer | Two independent six-face thin panels in probe | Deterministic export and reimport; non-rational B-spline support; probe JSON |
| `experiments/m19/firmament/m18-tangent.firmament` | Accepted | Same route | Same topology as safe | Same preservation/loss profile |
| `experiments/m19/firmament/m18-failing.firmament` | Accepted | Same route | Same topology as safe | Same preservation/loss profile |
| `experiments/m19/firmament/contact-metadata-concept.firmament` | Accepted and Concept IR resolved | Concept -> compile-time erasure; dummy CAD box/chamfer continues through Feature AIR | CAD witness body only | STEP contains CAD feature name, not mathematical fields; build JSON retains Concept IR |
| `experiments/m19/firmament/m18-research-assembly.firmament` | Accepted by Assembly M0 parser/compiler | hierarchy/interface/mate placement solve | No linked panel geometry artifact | Deterministic inspection JSON; semantic roles survive, two instances remain unconstrained |

The main CLI validates panel-only sources but `build` throws `Index was out of range. Must be non-negative and less than the size of the collection.` because that path indexes a solid list that is empty for a Panel document. The experimental probe therefore invokes the real panel compiler/materializer directly. This is an implementation gap, not a reason to pretend the CLI produced a BRep.

CLI validation returned exit code 0 for all three Panel fixtures (`0 fatal, 2 warning` each) and the Concept fixture (`0 fatal, 6 warning`). `aetheris asm inspect` returned exit code 0 for the Assembly fixture. Repeated probe runs produced byte-identical JSON; each spectral STEP export was identical to a second in-process export. The materialized spectral and zero panels each have 6 faces; spectral topology is 6 faces, 12 edges, and 8 vertices in all three cases.

## Layer ownership and authority

| Layer | What it owns now | Authority | Mathematical fit | Important erasure/loss |
|---|---|---|---|---|
| Firmament | Authored CAD declarations, compile-time structure, Panel expressions, names | Parser/source authority | Good shell; two incompatible expression grammars | No general analytic scalar/function language |
| Concept IR | Typed structural metadata and spatial values | Compile-time semantic authority | Good witness/annotation envelope | Explicitly erased before Feature AIR |
| Feature AIR | CAD feature intent and dependencies | Production for supported CAD routes; trace summaries elsewhere | Poor for theorem intent | Domain meaning narrowed to CAD operations |
| Construction AIR | Profiles, curves, frames, extrude/revolve/sections/composition | Bounded production; APIs fragmented/internal | Best shared geometric-construction candidate | No general analytic function or certified predicate layer |
| Panel / Surfacing IR | Bounded parametric/named patches, first jets, panel roles | Experimental real execution path | Best existing M18 representation | Approximation at materialization; manufacturing semantics attached |
| BRepPlan | Planned topology, stable IDs, provenance, roles | Authoritative only where a realization plan exists; otherwise trace-only | Useful realization manifest for admitted CAD features | No signed-side, tangency, or contact-order obligation |
| Kernel/Core | Analytic CAD geometry, BRep, factories, bounded queries | Production geometry authority | Reusable numerical geometry | Doubles/tolerances; query coverage is narrow |
| BRep | Realized topology and geometry bindings | Production result | Useful topology/artifact layer | Generating math expression and proof meaning gone |
| STEP AP242 | External CAD interchange, geometry/topology, limited PMI/names | Artifact, never theorem authority | Good inspection/interop | Concept fields, contact relation, formula, and certificate lost |

## Construct-to-mathematics mapping

| Firmament construct | Current meaning | Possible mathematical meaning | Faithfulness | Missing semantics | Recommended owner |
|---|---|---|---|---|---|
| Concept Struct | Compile-time typed CAD semantic schema/instance | witness, parameter set, certificate descriptor | High until erasure | durable downstream metadata/query handle | Firmament/Concept IR for metadata; proof stays Riemann |
| Assembly | part/panel hierarchy and placement constraints | object graph of spectral object, boundary, marker, symmetry object | Medium as composition; low as reasoning | generic semantic payloads and geometry linkage | Assembly for composition only |
| Profile/path | sketch and bounded region construction | piecewise analytic curve/domain | Medium-high for elementary segments | arbitrary analytic curve and certified predicates | Construction IR / KernelSDK |
| Parametric Panel | manufacturable bounded surface with thickness/material/orientation | bounded graph/field patch or contact neighborhood | High before materialization, low afterward | exactness kind, certified approximation, non-CAD patch abstraction | General surfacing IR above Panel |
| Mate/constraint | mechanical DOF removal | incidence/coincidence relation | Medium for spatial relations | generic predicates and proof meaning | Kernel geometry relations; not theorem layer |
| Transform/pattern | placement, repeat, mirror | symmetry/group action | High geometrically | preserve generator/action semantically through lowering | Construction/Assembly IR |
| Boolean | bounded CAD set operation | set algebra | Conceptually high, operationally bounded | general robust/certified intersection | KernelSDK |
| Selection/reference | stable topology/semantic binding | named subobject | High when identity survives | bindings to authored analytic objects | shared semantic layer |
| PMI | engineering annotation | human-facing theorem note | Low | arbitrary typed metadata and round trip | Do not overload PMI; Riemann/report side |
| BRepPlan | expected CAD realization topology | realization certificate | Medium for topology only | predicates, sidedness, contact expectations | Aetheris planning, if generalized |
| Feature | manufacturing/modeling operation | graph/restrict/reflect/intersect intent | Low | domain-neutral operation vocabulary | Do not reuse Feature AIR for math |

## Assembly audit

Assembly has stable identity and paths, nested assembly/part/panel instances, definition identities/templates, explicit/imported/mate-derived transforms, interfaces with roles/capabilities, exposed semantics, relations, tolerance stackups, and a constraint solver. Mates lower to axis coincidence/alignment, plane coincidence, point coincidence, and offset-along-axis. This is a real structured-composition abstraction, not merely a flat part list.

The calibration assembly represented a spectral panel, zero-reference panel, and contact-marker part with stable roles. The root and spectral object became anchored, the explicit contact-marker transform resolved, and the `PlaneCoincident` mate validated. The zero reference nevertheless remained underconstrained with all six reported degrees of freedom, useful evidence that a valid semantic mate is not itself a fully solved placement. However, the semantic value parser is limited to axis, plane, point, and dimension, and a panel definition remains a string rather than a link to `PanelIr` geometry.

**Answer:** Assembly generalizes cleanly as a named spatial object graph and composition shell. It does not generalize cleanly as a mathematical reasoning graph because its relations, payloads, solver, and materialization remain placement/mechanical. It is reusable without renaming for artifact composition, not as the home of contact theorems.

## Concept Struct audit

Concept definitions are structural field schemas. Instances support scalar, enum, and substantial spatial value types, nesting/member references, grids, regions, point sets, matching, and compile-time requirements. The M19 fixture encoded `contactFrequency`, `expectedMultiplicity`, `candidateC`, `candidateWeight`, `certificateStatus`, and `isCertified` using existing grammar.

Those values survived parsing, Concept resolution, exposed Concept values, and the CLI build JSON. They were absent from Feature AIR/BRepPlan and from STEP. Only the dummy CAD feature name used to satisfy current materialization (`ParserRouteWitness`) reached STEP. Assembly materialization also ignores general float/int/string/bool Concept values and only converts its recognized spatial subset.

**Answer:** mathematical semantics can ride inside Concept Structs for source-time validation and trace/reporting, so metadata alone does not require new grammar. Concept Struct cannot avoid a language/lowering extension when a durable, queryable analytic object or predicate is required. The current CAD-wrapper requirement and deliberate erasure make it unsuitable as the sole carrier of proof-relevant meaning.

## Surfacing and Panel audit

Panels support parametric and named mathematical surfaces, rectangular domains, first derivatives, orientation, material/thickness, provenance, ordered stable edges/corners, boundary/section/ruled constructions, and thin-solid topology generation. They can naturally author bounded graph patches, scalar-field graphs, parameter-domain neighborhoods, and some piecewise patches.

The abstraction fights pure mathematics in four places:

1. Panel requires engineering orientation, thickness, and material.
2. General parametric expressions have no named parameters or `pi` and are isolated from normal Firmament declarations.
3. Materialization replaces the authored surface with an adaptive sampled degree-1 non-rational B-spline.
4. Continuity and topology checks are numerical: network G0 is sampled and G1 is unsupported.

**Answer:** the underlying bounded parametric-patch IR generalizes cleanly; `Panel` itself should retain its CAD meaning. A future general surface-patch abstraction could be the small common base, with Panel adding fabrication semantics. M19 does not implement that split.

## Expression and parameter audit

Normal Firmament expressions provide only compile-time arithmetic over literals and named/dotted values using `+ - * /`, parentheses, and `mm`/`deg`. There is no `pi`, trigonometry, power, conditionals, functions, or free variables. Templates are finite compile-time specialization. Consequently neither

```text
1 + pi/4 - 2/pi
```

nor a reusable `F(c,r,w,t)` is expressible in the normal language.

Panel expressions provide `u`, `v`, integer powers, `sin`, and `cos`, so a fixed numeric specialization of the M18 graph is expressible. Because `pi` and normal `let` references are absent, M19 substituted decimal constants. That is a visualization specialization, not an exact symbolic encoding of the M18 formula or family.

## Curve and surface capability matrix

| Capability | Current status | Needed for Riemann | Classification |
|---|---|---|---|
| Exact-form line/circle/elementary quadric | Yes, double coefficients | Useful | Already exists |
| Ellipse/hyperbola | Represented; limited queries | Useful | Already exists, query gap |
| Non-rational B-spline curve/surface | Yes | Visualization/construction | Already exists |
| Rational NURBS | No general native authored/core representation | Optional | Genuinely missing |
| Extrusion/revolution surface | Yes | Useful for symmetry | Already exists |
| Ruled/section/boundary surface | Panel route | Useful | Exists in experimental surfacing |
| Generic parametric surface | `ParametricSurfaceIr`, not core Kernel surface | Essential candidate | Exists internally/supporting API; needs generalization/exposure |
| Generic parametric curve | No equivalent public abstraction | Important | Genuinely missing |
| Parametric domain | Rectangular Panel domain | Essential | Partial |
| Trimming/pcurves | Bounded CAD/STEP topology support | Useful | Needs generalization |
| First derivatives | Panel expression jets; selected analytic curves/surfaces | Essential | Partial |
| Second derivatives/curvature | No general API | Needed for strict local contact | Genuinely missing |

## EXACTNESS AUDIT

- Numeric coefficients and evaluation use IEEE-754 `double`. “Analytic” means an analytic representation evaluated numerically, not symbolic or algebraically exact geometry.
- General source units are dimension-checked; declaration tolerances may be recorded, but tolerance provenance is lost through arithmetic.
- Default kernel tolerances are linear `1e-6`, angular `1e-9`, and relative `1e-12` in `ToleranceContext`.
- BRep topology and point classification are tolerance-defined. Seam/ambiguous classifications are surfaced in some queries, but there is no uniform proof-status model.
- Panel parametric evaluation preserves the authored expression in double arithmetic only until materialization. Materialization reports an empirical sampled maximum residual against a requested `0.1 mm` tolerance and creates a degree-1 B-spline grid, up to `129 x 129`. This is not an interval or rigorous approximation certificate.
- Imported geometry carries no theorem-level exactness provenance. Analytic STEP entities may round-trip as analytic types, but that does not make their decimal coefficients symbolic.
- STEP numerics are emitted with the formatting pattern `0.###############`, approximately fifteen decimal digits. STEP representation and import introduce format/round-trip semantics, not a proof obligation.
- Intersections, containment, and topology guards use tolerance tests and bounded algorithms. Some boolean diagnostics recognize special tangent configurations, but no general certified contact predicate exists.

Therefore Riemann cannot currently ask one API whether a result is `Exact`, `CertifiedApproximate`, or `VisualApproximation`. For the M19 surface path the only truthful classification is `ApproximateVisualization`.

## Query and predicate inventory

| Predicate/query | Status | Semantics |
|---|---|---|
| Curve/surface evaluation | Analytic primitives and B-splines | Double evaluation |
| Tangent/normal | Available for several curves/surfaces and Panel first jets | Double differential result |
| Point on face/domain | Plane/cylinder/sphere subset | Tolerance-based; cone/torus/B-spline unsupported in general face-domain query |
| Ray cast / point containment | Available for bounded admitted shapes | Numerical/tolerance-based |
| Body interference | Sound bounded proof only for closed convex planar BReps | `DisjointOrTouching`, `Interfering`, or `Unsupported`; touching not separated from disjoint |
| Intersection | Specialized algorithms/boolean cases; e.g. transverse cone-plane | No general curve/surface intersection API |
| Distance / closest point | No general public query | Missing |
| Sidedness of whole curve/surface | No | Missing |
| Coincidence/overlap | Specialized/tolerance-based | Partial |
| Adjacency/manifold/orientation | BRep topology/validators | Available, tolerance/topology-defined |
| Curvature | No general public surface/curve query | Missing |
| Degeneracy/topology change | Selected validators/diagnostics | Partial and route-specific |
| Contact order/multiplicity | No | Missing |

### Tangency, multiplicity, and nonpenetration

Current APIs cannot distinguish in a general way among crossing, simple contact, tangent contact, higher-order contact, and coincident contact. `F(t*)=0` can be numerically evaluated in Panel IR; `F'(t*)=0` can be estimated or obtained from the first jet; `F''(t*)>0` requires an external finite difference because the surface expression API has no second jet. None of those evaluations constitutes a contact query between realized bodies.

Exact contact multiplicity belongs in Riemann unless Aetheris later gains a general differential/algebraic contact-order API. Aetheris should own generic jets and certified geometric predicates, not the theorem statement “multiplicity exactly two.”

The global predicate `F(t) >= 0` has no current geometric equivalent. There is no whole-domain side-of-plane query and no certified “no intersection except tangency” operation. Sampling the Panel is insufficient for theorem use.

### Symmetry

Reflection, rotation, circular/linear/mirror patterns, axis-based mates, and surfaces of revolution exist. They preserve generator/reference intent in some construction and assembly stages, but final BRep/STEP primarily preserve realized geometry and selected names. Symmetry is therefore strongest before BRep and can become merely materialized repetition afterward.

## M18 calibration

All three cases use the same `ParametricSurface` route over `u in [1.2, 1.9]`, with `v` extruding the graph into a narrow surface and a separate zero panel. The source specializes constants numerically and is labeled `M19-ApproximateVisualization`.

| Case | Analytic expectation | Probe at `t = pi/2` | Panel materialization | What Aetheris can conclude |
|---|---|---|---|---|
| Safe | strict separation | normalized `F = +0.0012732395447353428`; scaled height `+0.025464790894706855 mm` | approximated non-rational B-spline, `17 x 17`, sampled residual `0.0005211246 mm` | sampled/pointwise positive visualization only |
| Tangent | double contact | normalized `F ~= 3.2626e-16`; height `~= 6.525e-15 mm`; finite-difference first derivative `~-1.51e-10`, second `+0.77987` | same class/grid, sampled residual `0.0005214733 mm` | visually tangent and numerically consistent; no contact/multiplicity proof |
| Failing | penetration/sign crossing | normalized `F = -0.00127323954473469`; height `-0.025464790894693803 mm` | same class/grid, sampled residual `0.0005218219 mm` | sampled/pointwise penetration visualization only |

For every case, the authored expression and first jet exist in Panel IR; materialized BRep contains approximated support and stable panel topology; exported STEP reimports and contains `B_SPLINE_SURFACE_WITH_KNOTS`. The M18 formula, contact relation, candidate constant provenance, and multiplicity do not survive. The panel source name/material were not found as durable mathematical semantics in the reimported geometry.

### Topology-change observability and the useful observation

The cheap geometric observation is a **semantic/topological decoupling**: the same stable Panel edge/corner identity and the same six-face thin-solid topology persist across safe, tangent, and failing variants while the authored signed height changes from positive through zero to negative. This makes the phase transition visually obvious and reproducible, but it also proves that ordinary body topology alone does not expose the contact transition when the spectral and zero objects are materialized independently. A future signed-side/contact query should operate on the authored patch/domain before approximation, not infer the theorem from unchanged exported topology.

This is `ResearchEvidence`, not a theorem claim.

## BRepPlan, AIR, BRep, and STEP

Feature AIR is fundamentally CAD-feature-shaped; adding operations named `ConstructGraph` or `RestrictDomain` there would mix source-domain semantics into the wrong layer. Construction AIR is closer to domain-neutral geometry and could accept lowering from another frontend, but it is fragmented and lacks general analytic curves/surfaces and predicate contracts. Panel IR presently fills part of that gap.

`AirBRepPlan` records semantic roles, expected topology counts, stable IDs, provenance, and whether an actual realization plan makes it authoritative. This can act as a bounded CAD realization manifest. It cannot certify theorem-relevant sidedness, contact location/order, or expected intersection count today. Trace-only plans must not be mistaken for production authority.

BRep is valuable for topology, inspection, downstream CAD algorithms, and export. It is the wrong owner for the exact generating formula and proof evidence. STEP preserves useful analytic CAD entities when those entities exist, topology, and limited names/PMI. In the panel route it preserves the approximated B-spline, not the authored parameter expression. Concept metadata and contact semantics were absent after round trip. STEP is an external inspection artifact, never a reasoning authority.

For LLM-assisted work, deterministic JSON trace/probe reports are easiest to inspect; Firmament source is best for intent; STEP is useful for human/CAD visualization and independent geometry inspection; screenshots are communication aids only.

## Public package/API audit

Preview 2 publishes `Aetheris.Kernel.Core`, `Aetheris.Kernel.Firmament`, `Aetheris.Forge.Host`, and `Aetheris.Forge.KernelSDK`. Core and Firmament expose usable geometry, parser, and selected compiler types. KernelSDK is an extension-author façade, not a complete façade over every query. Surfacing and semantics are supporting/transitive project surfaces and should be treated as less stable than the named public façade. Feature AIR internals, `AirBRepPlan`, multiple resolvers, and substantial compiler plumbing are internal-only.

External Riemann use should therefore not bind directly to internal AIR/BRepPlan types. If M20 proceeds, the smallest public query contract should be exposed through a deliberately public geometry surface, with the experiment depending only on public Preview APIs.

## Gaps, ownership, and ranking

Scores are 1 (low) to 5 (high). Lower implementation cost and exactness risk are better.

| Candidate | Current classification | Riemann value | CAD value | Cost | Exactness risk | Reuse | Owner |
|---|---|---:|---:|---:|---:|---:|---|
| Whole-domain signed-side/contact classification for parametric patch vs plane | Genuinely missing | 5 | 5 | 3 | 4 | 5 | Aetheris KernelSDK/Surfacing |
| Preserve/query analytic provenance and approximation class through lowering | Needs generalization | 5 | 4 | 3 | 2 | 5 | Aetheris shared semantic geometry |
| Public generic parametric curve with jets/domain | Genuinely missing | 5 | 4 | 3 | 3 | 5 | Aetheris KernelSDK |
| Second jets and generic curvature | Genuinely missing | 5 | 4 | 3 | 4 | 5 | Aetheris KernelSDK |
| Certified intersection | Genuinely missing | 5 | 5 | 5 | 5 | 5 | Aetheris KernelSDK, later |
| Contact order/multiplicity | Genuinely missing | 5 | 3 | 4 | 5 | 4 | Generic differential part in Aetheris; theorem conclusion in Riemann |
| Arbitrary Concept metadata in assembly/materialization | Needs exposure/generalization | 3 | 4 | 2 | 1 | 4 | Firmament semantic layer |
| Exact symbolic expression algebra | Genuinely missing | 5 | 2 | 5 | 5 | 3 | Primarily Riemann; do not make CAD kernel a CAS |
| M18 theorem/certificate schema | Riemann-side only | 5 | 1 | 2 | 3 | 1 | Riemann |

Firmament is missing a shared analytic expression/function capability, named constants including `pi`, parameterized functions/families, and a general bounded parametric-curve declaration. It does **not** need zeta-specific syntax. Concept Struct is sufficient for descriptive metadata; it is insufficient for executable or durable analytic semantics.

## Reusable abstractions and architectural awkwardness

Reusable without production changes:

- stable named references and semantic identities;
- Assembly hierarchy/interfaces for spatial artifact composition;
- Concept Struct for source-time typed metadata and trace reporting;
- construction profiles/frames/extrude/revolve for elementary mathematical geometry;
- Panel `ParametricSurfaceIr` for deterministic bounded visual experiments;
- BRep/topology/STEP for inspection and interchange.

CAD-specific or unsafe to generalize directly:

- Feature AIR operation vocabulary;
- holes, edge finishes, manufacturing features, PMI taxonomy;
- Panel thickness/material/orientation as the base mathematical patch concept;
- mate DOF solver as a generic theorem constraint system;
- BRepPlan as a proof certificate;
- STEP as semantic/proof storage.

The sharpest awkwardness is that the mathematically interesting Panel expression exists, but the ordinary build path expects a CAD solid, and the real materializer immediately approximates the expression. Similarly, Concept Struct accepts the desired metadata but requires a CAD materialization witness and then erases it. These are clean audit findings: the parser is more general than the downstream CAD pipeline.

## Compiler-theory answer

Yes, partially. A CAD compiler can expose a domain-neutral middle representation without sharing the original source semantics: named coordinate frames, bounded parametric maps, profiles, composition, stable identities, derivatives, and realization provenance are geometric rather than mechanical. Aetheris has most of those pieces, but not yet as one public, exactness-aware reasoning layer. Its current compilation pipeline demonstrates a useful separation between authored intent and realized geometry, while also demonstrating that lowering can erase precisely the information mathematics needs.

The right architecture is not “reason over STEP.” It is:

```text
Riemann exact expression and proof obligations
    -> explicit, semantics-preserving bounded geometric adapter
    -> public parametric geometry + certified predicate contract
    -> optional BRep/STEP visualization artifact
```

The future proof obligation would be an explicit equivalence, owned by Riemann at its ends:

```text
F(t) >= 0 on D
    iff
the graph patch is classified on the nonnegative side of the zero plane on D
```

Aetheris would own the generic predicate result and its certificate; Riemann would own that this geometric construction represents `F` and that the returned certificate discharges the theorem obligation.

## Exactly one proposed next milestone

**M20 should prototype one public, tolerance-explicit `SignedSideQuery` for an authored bounded `ParametricSurfaceIr` against a plane, before B-spline/BRep materialization.**

The result should be `Positive`, `Negative`, `Touching`, `Crossing`, or `Unknown`, carry the parameter domain and an explicit evidence classification (`Sampled` initially or `Certified` only when justified), and never upgrade sampling to certification. Use the existing safe/tangent/failing fixtures as calibration, but make no Riemann production dependency and add no new Firmament syntax in that milestone.

This is the smallest experiment that tests whether Aetheris can move from visualization to geometric research IR. It is valuable to CAD for clearance/interference and surface-side classification, preserves the authored surface rather than reasoning over its approximation, and exposes the exactness gap without prematurely implementing symbolic algebra, general intersection, or contact multiplicity.

## Verification

- `dotnet build Aetheris.slnx --no-restore`: passed, 0 warnings and 0 errors.
- `dotnet test Aetheris.slnx --no-build --no-restore`: passed 2,696 tests; one test assembly reports no discoverable tests; no failures.
- `go test ./...`: passed.
- `go test -race ./...`: passed.
- M19 probe: three Panel parses/materializations/STEP round trips, Concept parse, and Assembly parse/compile passed; repeat report deterministic.
- `git diff --check`: passed in both repositories (the Riemann checkout reports only a line-ending conversion notice for `README.md`).

## Final answers to the requested decision points

1. Overall: promising research IR; visualization backend now; not theorem-adjacent.
2. Syntax: a federation of legacy, V2, Panel, Assembly, Drawing, and FEA front doors; the inventory is in `m19-firmament-syntax-audit.md` and the language snapshot.
3. Docs/parser: lane ambiguity, canonical/compatibility divergence, expression-grammar split, panel validate/build mismatch.
4. Assembly: generic composition shell, mechanical reasoning semantics.
5. Concept Struct: expressive metadata envelope, deliberately erased.
6. Panels: reusable bounded parametric core, CAD-specific wrapper and approximate realization.
7. Expressions/parameters: normal grammar is compile-time arithmetic only; Panel adds `u/v`, trig, integer powers; no free families.
8. Curves/surfaces: strong elementary CAD set plus non-rational B-splines; generic parametric surface only in Surfacing IR; no generic parametric curve/NURBS.
9. Exactness: double/tolerance-defined; Panel certificate is sampled residual, not rigorous.
10. Queries: useful evaluation/topology and bounded spatial queries; no comprehensive differential/contact suite.
11. Tangency: first-jet evidence possible before BRep; no general tangency classifier.
12. Multiplicity: absent; theorem conclusion remains Riemann-side.
13. Nonpenetration: no whole-domain sidedness predicate.
14. Topology change: not exposed by independently materialized panel topology.
15. Symmetry: good construction support, semantic preservation route-dependent.
16. BRepPlan: realization manifest for admitted CAD routes, not theorem certificate.
17. Construction AIR: strongest shared candidate, but incomplete/fragmented.
18. BRep: topology/artifact layer, not source-math authority.
19. STEP: inspection/interchange only; loses formula/contact/Concept metadata.
20. Public API: Core/Firmament/KernelSDK usable, with important internals not exposed.
21. Internal capability: Panel first jets and BRepPlan exceed the stable façade.
22. Safe: positive sampled/point evidence, visualization only.
23. Tangent: numerically consistent double contact, no certified contact/order query.
24. Failing: negative sampled/point evidence, visualization only.
25. Fixtures: five real sources plus one C# experimental probe.
26. Survived: Panel topology/approximate support; Concept IR in JSON; Assembly identities/roles; selected CAD names.
27. Lost: exact formula/provenance, metadata in geometry, contact relation/order, certification, many source names.
28. Observation: stable topology across a signed-height phase transition reveals where a pre-BRep predicate must live.
29. Reusable: references, frames/profiles, parametric IR, identities, composition shell, deterministic artifacts.
30. CAD-specific: feature vocabulary, fabrication Panel semantics, mates, PMI, BRepPlan obligations.
31. Kernel gaps: domain signed side, generic parametric curve, higher jets, certified intersection/contact.
32. Firmament gaps: shared analytic expressions, constants/functions/free parameters, bounded curve declaration.
33. Ownership: generic geometry in Aetheris; exact formula, theorem, and proof evidence in Riemann.
34. Ranking: signed-side query first; table above records the remaining candidates.
35. Concept Struct avoids parser changes only for compile-time descriptive metadata.
36. Assembly generalizes for spatial composition, not generic logical constraints.
37. Panel core generalizes; Panel manufacturing wrapper should not be redefined.
38. Use Aetheris as visualization backend now and potential geometric research IR; not a verifier.
39. Compiler theory: the middle representation is reusable only if semantic/exactness contracts survive lowering.
40. One next milestone: the pre-materialization `SignedSideQuery` prototype described above.
