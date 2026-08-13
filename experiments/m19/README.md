# M19 Aetheris calibration experiments

These experiments require the Preview 2 Aetheris checkout as a sibling of this repository (`../Aetheris`). They are diagnostic only and create no production dependency.

Run the parser/materialization/STEP round-trip probe from the Riemann root:

```powershell
dotnet run --project experiments/m19/AetherisM19Probe/AetherisM19Probe.csproj -- experiments/m19/firmament artifacts/m19
```

The probe parses all three Panel cases through `FirmamentV2Parser`, evaluates the authored `ParametricSurfaceIr`, materializes both panels through `SurfacePatchBrepMaterializer`, exports STEP twice to check determinism, and reimports the first export. It also parses/resolves the Concept Struct and compiles the Assembly M0 fixture. Output paths in `m19-probe-report.json` are repository-relative and `generatedUtc` is fixed, so repeated identical runs produce byte-identical reports.

Validate the Firmament V2 fixtures using the CLI from the Aetheris root:

```powershell
dotnet run --project Aetheris.CLI -- validate ../Riemann/experiments/m19/firmament/m18-safe.firmament
dotnet run --project Aetheris.CLI -- validate ../Riemann/experiments/m19/firmament/m18-tangent.firmament
dotnet run --project Aetheris.CLI -- validate ../Riemann/experiments/m19/firmament/m18-failing.firmament
dotnet run --project Aetheris.CLI -- validate ../Riemann/experiments/m19/firmament/contact-metadata-concept.firmament
dotnet run --project Aetheris.CLI -- asm inspect ../Riemann/experiments/m19/firmament/m18-research-assembly.firmament --json
```

The panel-only fixtures validate, but Preview 2's ordinary `build` command currently indexes an empty solid collection. The probe uses the actual independent Panel execution path and labels every result `ApproximateVisualization`; it does not reinterpret samples as theorem evidence.

`concept-metadata.step` was generated separately with:

```powershell
dotnet run --project Aetheris.CLI -- build ../Riemann/experiments/m19/firmament/contact-metadata-concept.firmament --output ../Riemann/artifacts/m19/concept-metadata.step --json
```

It exists to test Concept-to-STEP semantic preservation. The mathematical fields are absent; only the CAD carrier's feature name survives.
