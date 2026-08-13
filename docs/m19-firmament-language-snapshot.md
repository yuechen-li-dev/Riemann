# Firmament Preview 2 language snapshot

This snapshot records what the Preview 2 implementation accepts, not what older documentation claims it accepts. The code baseline audited was Aetheris commit `c53c840` (`2.0.0-preview.2`). Parser-validated sources live in `experiments/m19/firmament/`.

## There is not one Firmament parser

Preview 2 has six notable authored front doors:

1. `FirmamentTopLevelParser` parses the legacy `firmament/model/ops` JSON-or-TOON document.
2. `FirmamentV2Parser` parses the current compatibility and canonical text dialects.
3. `FirmamentPanelCompiler` parses a panel/surface document embedded in the V2 route.
4. `AssemblyM0Parser` parses XML-like assembly documents.
5. `FirmamentDrawingCompiler` parses Drawing concepts, templates/applications, and literal drawings before invoking V2 for their geometry source.
6. `FirmamentAnalysisCompiler` in `Aetheris.FEA` strips and parses a bounded `analysis LinearElastic` block, then invokes V2 or its explicit bounded body recognizer.

Consequently, syntax accepted by one lane may be rejected by another. In particular, the canonical V2 root deliberately rejects compatibility declarations such as `Solid`, `Let`, `Fill`, `Manufacturing`, `Feature`, and `Expose`.

## Current V2 document vocabulary

The V2 document model can carry solids, modify blocks, templates, PMI, recognition/replacement declarations, scalar lets and records, manufacturing concepts, feature concepts, lattice declarations, construction planes, profiles, compositions, selections, static-authoring declarations, volume assertions, template instantiations, and panels. The canonical-root dispatch recognizes:

```text
Box Cylinder RoundedBox Frustum StandardPart ExactCoaxialPart
Concept Struct Construction Plane Profile Compose EdgeFinish Record Static
Template Selection InlineStep Recognize Replace Pmi Modify Match Require Assert
```

Compatibility parsing additionally recognizes direct solid declarations and legacy primitives including `Cone`, `Sphere`, and `Torus`, plus `Let`, `Fill`, `Manufacturing`, `Feature`, and `Expose`.

### Expressions and values

General scalar declarations support literals, identifier and dotted-member references, parentheses, and binary `+`, `-`, `*`, `/`. Values are compile-time `int` or IEEE-754 `double`; units are `mm` and `deg`. Length and angle dimensions are checked. Length/angle declarations can carry bilateral or asymmetric tolerances, but arithmetic drops attached tolerance with a diagnostic.

General expressions do **not** support `pi`, trigonometric functions, powers, conditionals, user-defined functions, or free runtime parameters. Function-shaped syntax is rejected. Reusable templates are compile-time specialization, not runtime functions.

An accepted member-reference and advanced-feature fragment is:

```firmament
Struct ContactMetadataCarrier : ContactMetadata {
  Units: mm
  Box Carrier { Size: [1mm, 1mm, 1mm] }
  Modify Carrier {
    EdgeFinish ParserRouteWitness {
      Face: +Z
      Target: Boundary
      Kind: Chamfer
      Distance: 0.1mm
    }
  }
  Expose { CandidateC: M18Certificate.CandidateC }
}
```

Panel parametric expressions are a separate, richer expression grammar. They provide variables `u` and `v`, arithmetic, unary signs, integer powers, `sin`, `cos`, numeric constants, parentheses, and optional `mm`. They still provide no `pi`, named constants, `let` references, conditionals, or user functions.

### Records, static authoring, and concepts

`Record`, static arrays/records/tables, enums, `Match`, `Require`, templates, and patterns are compile-time authoring facilities. Concept definitions declare structural field requirements. `Concept Struct` instances can contain literal and spatial values including:

```text
Length Angle Bool Int Float String Enum
Point2 Point3 Vector3 Axis Plane Box2 Box3 Region2 PointSet
```

They may use members, enum values, grids, boxes, planes, axes, and compile-time matching. At present, the Concept resolution path requires a materialized `Struct`/model with a `Box` and a non-empty `Modify` operation. Concept instances are explicitly marked `CompileTimeOnlyErased`, and the Concept IR document reports `ErasedBeforeFeatureAir`.

This accepted fixture demonstrates scalar mathematical metadata without a grammar change:

```firmament
Concept ContactMetadata {
  ContactFrequency: Float
  ExpectedMultiplicity: Int
  CandidateC: Float
  CandidateWeight: Float
  CertificateStatus: String
  ExactCandidate: Bool
}

Concept Struct M18Certificate : ContactMetadata {
  ContactFrequency: Float = 1.5707963267948966
  ExpectedMultiplicity: Int = 2
  CandidateC: Float = 1.1487783910298668
  CandidateWeight: Float = 0.17103742708003555
  CertificateStatus: String = "RiemannExact_AetherisApproximateVisualization"
  ExactCandidate: Bool = true
}
```

The values survive parse/Concept IR and JSON build reporting, but not Feature AIR, BRepPlan, or STEP metadata.

### Construction and CAD features

The construction vocabulary includes named planes, named points and guides, paths with line/arc/close segments, profiles with loops and segments, extrusion, composition, holes, slots, selections, section transitions, and revolve/extrusion-like construction nodes. CAD modification vocabulary includes semantic shaft/counterbore/countersink holes, blind/through-all termination, chamfers, fillets, drafts, regions/cuts, edge finishes, booleans in the legacy operation lane, recognition/replacement, inline STEP, lattice, standard/exact-coaxial parts, and PMI.

Selections and exposed semantics use stable semantic identities and named references. They are powerful for CAD topology intent, but are not general predicates over authored mathematical objects.

### Panels and surfaces

A panel declares a support surface, front/back orientation, thickness, and material:

```firmament
Model M19PanelExample {
  Units: mm;
  Panel SpectralSurface {
    Surface: ParametricSurface {
      DomainU: [1.3, 1.85]; DomainV: [-0.05, 0.05];
      X: 20mm * u;
      Y: 20mm * v;
      Z: 20mm * (1 + 2 * (cos(u) - 1) / u^2);
    }
    Orientation: Front;
    Thickness: 0.2mm;
    Material: "M19-ApproximateVisualization";
  }
}
```

Supported surface families are `ParametricSurface`, `HyperbolicParaboloid`, `ParabolicCylinder`, `EllipticParaboloid`, `Helicoid`, `RuledSurface`, `RuledTransition`, `BoundaryPatch`, and `SectionSurface`. Boundaries can be lines, arcs, or circles; section surfaces carry section lists.

The authored `ParametricSurfaceIr` retains its expression tree and exact first-derivative program in double arithmetic over a rectangular parameter domain. Panel materialization then samples it adaptively and emits an approximated, non-rational, degree-1 tensor B-spline support (default requested tolerance `0.1 mm`, maximum grid `129 x 129`). Thus the source is parametric, but the BRep is an empirical approximation.

Panels retain stable panel, edge, and corner identities, ordered boundary edges, orientation, thickness/material, and provenance. Network validation checks sampled G0 edge coincidence; G1 is explicitly unsupported. Trimming is the four-sided panel boundary/materialization route, not a general exact trimmed-surface algebra.

### Assembly

Assembly syntax is a separate XML-like language:

```firmament
Interface ZeroContact {
  Role Spectral requires PlaneCapable;
  Role Zero requires PlaneCapable;
  Lower PlaneCoincident Spectral.ZeroPlane Zero.ZeroPlane;
}
Assembly M18ResearchObjects {
  <Assembly M18ResearchObjects>
    <Panel SpectralObject = M18SpectralPanel>
      Semantic ContactFrame { Plane ZeroPlane = [0,0,0] normal [0,0,1]; }
    </Panel>
    <Part ContactMarker = ContactMarkerPart>
      Placement LegacyExplicit = [1,0,0,0, 0,1,0,0, 0,0,1,0, 31.4159,0,0,1];
      Semantic Marker { Point Contact = [31.4159,0,0]; }
    </Part>
  </Assembly>
  Anchor: M18ResearchObjects.SpectralObject.ContactFrame;
}
```

It supports assembly/part/panel instances, nested hierarchy and stable paths, definitions/templates, explicit/imported/mate-derived transforms, interfaces with roles/capabilities, exposed semantics, relations, tolerance stackups, and mates that lower to axis coincidence/alignment, plane coincidence, point coincidence, and axis offset constraints. Semantic values are limited to axis, plane, point, and dimension. It can describe generic named composition, but its solver and materializer are mechanical placement systems. A `<Panel>` definition is an identity string, not a link to a compiled `PanelIr`.

### Drawing and analysis side languages

Drawing accepts `Concept Drawing`, parameterized `Template <...> Drawing`, template application, or a literal `Drawing`. A drawing body selects a geometry `Source`, orientation, views with direction/projection/hidden-line policy and PMI assignments, title/part metadata, tables, notes, metadata, and BOM emission. It lowers to `DrawingIr` and SVG/PDF/PPTX artifacts. This is a document-production side language, not geometric reasoning IR.

The FEA package accepts one bounded `analysis LinearElastic Name { ... }` block with body/resource selection, material, fixed regions, traction/force/pressure loads, lattice resolution, and requested result fields. It lowers to `LinearElasticAnalysisIr` after V2 geometry parsing or a narrow exact box/box-with-hole recognizer. It is an application-specific continuum-analysis extension rather than general Firmament expression support.

### Legacy operation document

The legacy root has exactly `firmament`, `model`, `ops`, plus optional `schema` and `pmi`. Its operation vocabulary is:

```text
Box Cylinder Cone Torus Sphere TriangularPrism HexagonalPrism
StraightSlot RoundedCornerBox SlotCut LibraryPart
Add Subtract Intersect Draft Chamfer Fillet
ExpectExists ExpectSelectable ExpectManifold
PatternLinear PatternCircular PatternMirror
```

This remains a real parser lane, but it is not evidence that identical declarations are accepted at the canonical V2 root.

## Rejections that define the boundary

- `pi`, `sqrt`, arbitrary function calls, user functions, general powers, and conditionals are absent from general Firmament expressions.
- Free parameters such as `F(c, r, w, t)` are absent; templates specialize and erase at compile time.
- `Sphere`, `Torus`, and `Cone` are not canonical-root primitives even though compatibility/legacy routes support them.
- Generic NURBS with rational weights are absent from the authored/core curve and surface models.
- Assembly semantic payloads are not arbitrary Concept values; its materialization only exposes axis, plane, point, and length-like dimensions.
- A panel-only source validates but the CLI `build` path assumes at least one solid and currently throws an index-out-of-range exception. The independent panel materializer succeeds.

## Code authority

Primary implementation sources audited:

- `Aetheris.Kernel.Firmament/FirmamentV2/FirmamentV2Parser.cs`
- `Aetheris.Kernel.Firmament/FirmamentV2/FirmamentV2Ast.cs`
- `Aetheris.Kernel.Firmament/FirmamentTopLevelParser.cs`
- `Aetheris.Kernel.Firmament/Concepts/ConceptIr.cs`
- `Aetheris.Kernel.Firmament/Assembly/AssemblyM0Parser.cs`
- `Aetheris.Kernel.Firmament/Assembly/AssemblyM0Models.cs`
- `Aetheris.Kernel.Firmament/Surfacing/FirmamentPanelCompiler.cs`
- `Aetheris.Kernel.Firmament/Surfacing/SurfaceScalarExpression.cs`
- `Aetheris.Kernel.Firmament/Surfacing/ParametricSurfaceIr.cs`
- `Aetheris.Kernel.Firmament/Surfacing/PanelFactory.cs`
- `Aetheris.Kernel.Firmament/Drawing/FirmamentDrawingCompiler.cs`
- `Aetheris.FEA/Firmament/FirmamentAnalysisCompiler.cs`
