# M19 parser-backed Firmament syntax audit

Status terms below mean: **production** participates in a real build/materialization path; **bounded production** is real but only for admitted cases; **compile-time** is erased before geometry AIR; **trace-only** records diagnostics/plans but is not geometry authority; **experimental** is the Preview 2 panel/assembly route exercised independently.

| Construct | Parser / semantic node | Accepted shape | Validation and lowering | Status | Mathematical interpretation | Classification |
|---|---|---|---|---|---|---|
| Legacy document | `FirmamentTopLevelParser` | `firmament`, `model`, `ops`; optional `schema`, `pmi` | Schema/root checks, then legacy op lowering | Production compatibility | Ordered construction program | CAD-specific |
| V2 canonical document | `FirmamentV2Parser` / `FirmamentV2Document` | Canonical declarations at root | Parser validation then frontend/compiler routes | Production, route-dependent | Named declarative module | Domain-general shell |
| Primitive solids | `Box/Cylinder/RoundedBox/Frustum/...Declaration` | Named dimensions/properties | Feature lowering and kernel factory | Production | Elementary exact sets | Potentially mathematical |
| Compatibility primitives | `Cone/Sphere/Torus/...` | Direct records/solid declarations | Compatibility lowering | Production compatibility | Analytic sets/surfaces | Potentially mathematical |
| Standard/imported parts | `StandardPart`, `ExactCoaxialPart`, `LibraryPart` | Catalog/import parameters | Part resolver/importer | Bounded production | None beyond reusable object | CAD-specific |
| Inline STEP | `InlineStepDeclaration` | Embedded/imported STEP payload | STEP import to BRep | Production, imported uncertainty | External geometric witness | Unclear |
| Scalar declaration | `LetDeclaration`, expression nodes | literal/reference/dotted ref; `+ - * /` | Compile-time typed evaluation | Compile-time | Scalar constants and derived values | Potentially mathematical |
| Units/tolerance | scalar type/value/tolerance nodes | `mm`, `deg`; bilateral/asymmetric tolerances | Dimension checks; arithmetic drops tolerance | Compile-time / tolerance metadata | Dimensional quantities | Domain-general |
| Record | record/static record nodes | named fields and member references | Compile-time expansion | Compile-time | Structured data | Domain-general |
| Static array/table | static authoring nodes | arrays, records, tables | Compile-time expansion/iteration | Compile-time | Finite indexed data | Domain-general |
| Enum / `Match` / `Require` | static semantic nodes | finite alternatives and compile-time guards | Static evaluation | Compile-time | Case split / precondition | Domain-general |
| Template / pattern | template/instantiation/pattern nodes | reusable parameterized authoring | Compile-time specialization and erasure | Compile-time | Finite family generator | Potentially mathematical |
| Concept definition | `ConceptDefinitionIr` | structural field requirements | Concept type resolution | Compile-time | Schema/proposition-shaped metadata | Potentially mathematical |
| Concept Struct | `ConceptStructInstanceIr` | typed scalar/spatial fields and members | Requires current CAD materialization wrapper; erased before Feature AIR | Compile-time | Witness/certificate metadata container | Potentially mathematical |
| Construction plane | construction plane node | named spatial frame/reference | Construction lowering | Bounded production | Coordinate chart/reference boundary | Domain-general |
| Path | path/line/arc/close nodes | named points, guides, segments | Construction validation | Bounded production | Piecewise elementary curve | Potentially mathematical |
| Profile / loops | profile/loop/segment nodes | bounded planar regions | Construction lowering to sketch/profile geometry | Bounded production | Domain/region boundary | Potentially mathematical |
| Extrude / revolve | construction operations | profile plus direction/axis/extent | Construction AIR to BRep | Bounded production | Product or orbit construction | Potentially mathematical |
| Section transition | section-transition nodes | ordered profiles/sections | Specialized construction lowering | Bounded production | Interpolating family | Potentially mathematical |
| Compose | compose nodes | named construction composition | Construction resolution | Bounded production | Object composition | Domain-general |
| Selection/reference | selection and semantic reference nodes | named/dotted semantic targets | Resolver/topology binding | Production within supported features | Named subobject | Domain-general |
| Aliases/exposure | semantic alias/expose nodes | named semantic topology references | Binding/exposure; canonical restrictions apply | Production, route-dependent | Stable role/name | Domain-general |
| Modify | `ModifyDeclaration` | holes, finishes, regions/cuts | Feature lowering | Production | Transformation of a set | CAD-specific implementation |
| Semantic hole | shaft/counterbore/countersink records | axis/point, diameter, blind/through-all | Feature operation and semantic binding | Production | Removed cylindrical region | CAD-specific |
| Edge finish | chamfer/fillet finish nodes | selected edge plus dimensions | Feature lowering, bounded kernels | Bounded production | Local offset/rounding transform | CAD-specific |
| Draft | legacy `Draft` op | face/direction/angle parameters | Legacy feature lowering | Bounded production | Affine-ish boundary change | CAD-specific |
| Boolean | legacy `Add/Subtract/Intersect` | operand references | Bounded boolean kernels | Bounded production | Set union/difference/intersection | Domain-general concept, CAD-bounded implementation |
| Recognition/replacement | `Recognize`, `Replace`, `Match` routes | named shape/feature matching | Recognition trace and bounded replacement | Route-dependent | Pattern recognition/rewrite | Domain-general idea |
| Lattice | lattice declaration | CAD lattice parameters | Feature/materialization path | Bounded/experimental | Discrete repeated set | Potentially mathematical |
| PMI | PMI declaration/value nodes | diameter, datum, distance, flatness, parallel, perpendicular, coplanar | Semantic PMI and limited STEP emission | Production subset | Annotation/constraint | CAD-specific vocabulary |
| Assertion | `Assert`, volume assertion, legacy expectations | volume/topology/selectability/manifold checks | Build-time diagnostic | Bounded production | Runtime proposition check | Potentially mathematical |
| Panel | `FirmamentPanelCompiler`, `PanelIr` | support, orientation, thickness, material | Panel compilation then thin-solid materialization | Experimental real path | Bounded parameter-domain patch | Potentially mathematical |
| Parametric surface | `ParametricSurfaceIr` | `X/Y/Z` expressions over rectangular `U/V` | First-jet evaluation, then adaptive sampled B-spline | Potential reasoning IR before materialization | Parametric map | Potentially mathematical |
| Named panel surfaces | typed surface IR nodes | hyperbolic/parabolic/elliptic, helicoid, ruled, boundary, section | Surface-specific compilation/materialization | Experimental | Analytic/piecewise patch | Potentially mathematical |
| Panel boundary | line/arc/circle boundary IR | ordered four-edge boundary | Stable identities; thin BRep topology | Experimental | Domain boundary | Domain-general |
| Panel network | panel adjacency validator | shared boundary relationships | Sampled G0 check; G1 unsupported | Experimental | Patch complex | Potentially mathematical |
| Assembly tree | `AssemblyM0Parser`, instance nodes | `<Assembly>`, `<Part>`, `<Panel>` hierarchy | Identity/path resolution and placement solve | Experimental real path | Named object graph | Domain-general shell |
| Transform | assembly transform nodes | translation/rotation/imported/mate-derived | Transform resolution | Experimental/bounded | Coordinate change | Domain-general |
| Interface/capability | assembly interface nodes | named role and capability | Stable semantic interface | Experimental | Typed port/role | Domain-general |
| Mate/constraint | mate nodes and solver constraints | axis/plane/point coincidence/alignment/offset | Mechanical DOF solver | Experimental | Geometric relation | CAD-specific implementation |
| Relation/tolerance stack | assembly relation nodes | assembly measurement relations | Stack evaluation | Experimental | Constraint network | Potentially mathematical, CAD vocabulary |
| Assembly semantic value | semantic axis/plane/point/dimension | exposed named values | Materialization exports only spatial subset | Experimental | Named geometric datum | Domain-general but narrow |
| Drawing concept/template | `FirmamentDrawingCompiler` parsed concept/template records | `Concept Drawing`, `Template <...> Drawing`, application or literal | Erases drawing declarations, parses geometry through V2, emits `DrawingIr` | Production side language | Structured projection/report | CAD-specific document layer |
| Drawing view/annotation | drawing view/metadata nodes | source, direction, projection, hidden lines, PMI, notes, table, BOM | Projection/layout to SVG, PDF, PPTX | Production side language | View/projection of an object | Domain-general idea, CAD output |
| Linear-elastic analysis | `FirmamentAnalysisCompiler` / `LinearElasticAnalysisIr` | `analysis LinearElastic`, body/resource, material, constraints, loads, lattice, results | Strips analysis block, invokes V2 or bounded body recognizer, lowers to continuum IR | Bounded production side language | PDE/continuum problem | Potentially mathematical, engineering vocabulary |

## Curve and surface inventory

| Representation | Constructible | Query support | STEP | Exactness status |
|---|---:|---|---|---|
| 3D line | Yes | evaluation/tangent and many bounded queries | Import/export | Analytic form, double coefficients |
| Circle / arc | Yes | evaluation/tangent; bounded intersections | Import/export | Analytic form, tolerance comparisons |
| Ellipse | Yes | evaluation | Import/export | Analytic form, limited queries |
| Hyperbola | Yes | evaluation/derivatives/tangent | Import/export | Analytic form, limited queries |
| Non-rational B-spline curve | Yes | polynomial evaluation/tangent | Import/export | Exact for stored doubles; not symbolic |
| Rational NURBS curve | No general core model | No | STEP coverage is not a native authored rational model | Missing |
| Plane/cylinder/cone/sphere/torus | Yes | evaluation/normal; query support varies | Import/export | Analytic form, double/tolerance semantics |
| Linear extrusion / revolution surface | Yes | evaluation/normal | Import/export | Analytic generating form where retained |
| Non-rational B-spline surface | Yes | evaluation/normal | Import/export | Exact for stored control net; often approximated input |
| Generic callback/analytic surface in core | No | No | No | Missing |
| Panel parametric surface | Yes in Surfacing IR | expression evaluation and first jet | Approximate B-spline export | Authored double expression; sampled materialization |
| Ruled/boundary/section panel surfaces | Yes in Surfacing IR | surface-specific, limited | Materialized B-spline/solid | Route-dependent approximation |
| Pcurves/general exact trims | STEP/topology support is bounded | No general public reasoning API | Partial | CAD topology, tolerance-defined |

## Parser and documentation discrepancies

- The apparent language is a federation of parser lanes; a feature documented without its lane is ambiguous.
- V2 canonical syntax is intentionally narrower than compatibility syntax.
- General expression syntax is much smaller than the Panel expression syntax.
- `validate` accepts a panel-only document, while CLI `build` currently assumes a solid document and crashes. Validation therefore does not prove that the main CLI materialization route accepts the document shape.
- Concept metadata is visible in Concept IR/build JSON but is deliberately erased before Feature AIR; its appearance in a frontend report is not downstream preservation.
- Assembly `<Panel Definition="...">` does not compile or attach the named Firmament panel geometry.
- Panel G1 continuity is represented as a desired notion in the architecture but is explicitly unsupported by the current validator.
- STEP accepts rich geometry, but custom mathematical Concept fields are not automatically mapped to STEP names or PMI.

## Extension points actually present

The practical extension seams are new V2 declaration/parser branches, new Concept value/materialization mappings, new construction/feature lowerers, new Panel surface kinds, new kernel factories/queries, and new STEP semantic mappings. None is a declarative plugin mechanism: each requires implementation changes and tests. Internal AIR/BRepPlan types are not a stable third-party extension API.

## Executable witnesses

- `experiments/m19/firmament/m18-safe.firmament`
- `experiments/m19/firmament/m18-tangent.firmament`
- `experiments/m19/firmament/m18-failing.firmament`
- `experiments/m19/firmament/contact-metadata-concept.firmament`
- `experiments/m19/firmament/m18-research-assembly.firmament`
- `experiments/m19/AetherisM19Probe/Program.cs`

The probe report is `artifacts/m19/m19-probe-report.json`.
