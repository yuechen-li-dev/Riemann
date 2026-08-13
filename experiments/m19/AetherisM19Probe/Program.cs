using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using Aetheris.Kernel.Core.Step242;
using Aetheris.Kernel.Firmament.Assembly;
using Aetheris.Kernel.Firmament.FirmamentV2;
using Aetheris.Surfacing;

if (args.Length != 2)
{
    Console.Error.WriteLine("usage: AetherisM19Probe <fixture-directory> <output-directory>");
    return 2;
}

var fixtureDirectory = Path.GetFullPath(args[0]);
var outputDirectory = Path.GetFullPath(args[1]);
Directory.CreateDirectory(outputDirectory);
var contact = Math.PI / 2d;
var cases = new[] { "safe", "tangent", "failing" };
var reports = new List<object>();

foreach (var caseName in cases)
{
    var sourcePath = Path.Combine(fixtureDirectory, $"m18-{caseName}.firmament");
    var source = File.ReadAllText(sourcePath);
    var parse = FirmamentV2Parser.Parse(source, fixtureDirectory);
    if (!parse.IsSuccess || parse.Document?.Panels is not { Count: 2 } panels)
        throw new InvalidOperationException($"{caseName} parse failed: {string.Join("; ", parse.Diagnostics)}");

    var spectral = panels.Single(panel => panel.StableId == "panel:spectral_surface");
    var zero = panels.Single(panel => panel.StableId == "panel:zero_boundary");
    var u = (contact - spectral.ParameterDomain.U.Minimum) /
            (spectral.ParameterDomain.U.Maximum - spectral.ParameterDomain.U.Minimum);
    const double v = 0.5d;
    const double h = 1e-5d;
    var atContact = spectral.SurfaceConstruction.Evaluate(u, v);
    var left = spectral.SurfaceConstruction.Evaluate(u - h, v);
    var right = spectral.SurfaceConstruction.Evaluate(u + h, v);
    var derivativePerNormalizedU = (right.Z - left.Z) / (2d * h);
    var secondPerNormalizedUSquared = (right.Z - 2d * atContact.Z + left.Z) / (h * h);

    var firstExport = ExportPanel(spectral);
    var secondExport = ExportPanel(spectral);
    var zeroExport = ExportPanel(zero);
    var firstStep = firstExport.Step;
    var secondStep = secondExport.Step;
    var zeroStep = zeroExport.Step;
    var stepPath = Path.Combine(outputDirectory, $"m18-{caseName}-spectral.step");
    File.WriteAllText(stepPath, firstStep);
    File.WriteAllText(Path.Combine(outputDirectory, $"m18-{caseName}-zero.step"), zeroStep);
    var reimport = Step242Importer.ImportBody(firstStep);

    reports.Add(new
    {
        caseName,
        sourcePath = ReportPath(sourcePath),
        parseResult = "accepted",
        validationResult = parse.Diagnostics,
        loweringRoute = "FirmamentV2Parser -> PanelIr/ParametricSurfaceIr -> sampled non-rational BSpline -> thin BRep -> STEP AP242",
        classification = "ApproximateVisualization",
        atContact = new { t = contact, scaledZMm = atContact.Z, fApprox = atContact.Z / 20d },
        finiteDifference = new { derivativePerNormalizedU, secondPerNormalizedUSquared },
        panel = new
        {
            spectral.StableId,
            materialization = spectral.SurfaceConstruction.MaterializationKind.ToString(),
            spectral.SurfaceConstruction.Approximation,
            boundaryIds = spectral.BoundaryEdges.Select(edge => edge.StableId).ToArray()
        },
        brepResult = new
        {
            route = "thin panel emitted",
            firstExport.FaceCount,
            firstExport.EdgeCount,
            firstExport.VertexCount,
            zeroFaceCount = zeroExport.FaceCount
        },
        stepResult = new
        {
            path = ReportPath(stepPath),
            deterministic = firstStep == secondStep,
            sha256 = Sha(firstStep),
            reimported = reimport.IsSuccess,
            reimportDiagnostics = reimport.Diagnostics.Select(diagnostic => diagnostic.Message).ToArray(),
            reimportSurfaceKinds = reimport.IsSuccess
                ? reimport.Value.Geometry.Surfaces.Select(surface => surface.Value.Kind.ToString()).Distinct().Order().ToArray()
                : []
        },
        traceResult = "Panel path has no Feature-AIR/Construction-AIR trace; PanelIr is the retained semantic surface layer"
    });
}

var conceptPath = Path.Combine(fixtureDirectory, "contact-metadata-concept.firmament");
var conceptParse = FirmamentV2Parser.Parse(File.ReadAllText(conceptPath), fixtureDirectory);
var concept = conceptParse.Document?.ConceptIr;

var assemblyPath = Path.Combine(fixtureDirectory, "m18-research-assembly.firmament");
var assemblyParse = new AssemblyM0Parser().ParseFile(assemblyPath);
var assemblyCompile = assemblyParse.Source is null
    ? null
    : new AssemblyM0Compiler().Compile(assemblyParse.Source, assemblyParse.ElapsedMilliseconds);

var report = new
{
    generatedUtc = DateTimeOffset.UnixEpoch,
    aetherisVersion = typeof(FirmamentV2Parser).Assembly.GetName().Version?.ToString(),
    cases = reports,
    conceptStruct = new
    {
        sourcePath = ReportPath(conceptPath),
        parsed = conceptParse.IsSuccess,
        conceptParse.Diagnostics,
        erasureStatus = concept?.ErasureStatus,
        concepts = concept?.Concepts.Select(item => item.Name).ToArray(),
        structs = concept?.Structs.Select(item => new { item.Name, item.ErasureStatus, members = item.Members.Keys.Order().ToArray() }).ToArray(),
        exposed = concept?.MaterializedStruct.ExposedMembers.Select(member => new { member.Name, type = member.Type.ToString(), phase = member.Phase.ToString(), category = member.MaterializationCategory.ToString() }).ToArray()
    },
    assembly = new
    {
        sourcePath = ReportPath(assemblyPath),
        parsed = assemblyParse.IsSuccess,
        parseDiagnostics = assemblyParse.Diagnostics,
        compiled = assemblyCompile?.IsSuccess,
        compileDiagnostics = assemblyCompile?.Diagnostics,
        instances = assemblyCompile?.Ir?.Instances.Select(instance => new { instance.StableId, path = instance.Path.ToString(), kind = instance.Kind.ToString(), semantics = instance.SemanticRoot.ExposedMembers.Keys.Order().ToArray() }).ToArray(),
        mates = assemblyCompile?.Ir?.Mates.Select(mate => new { mate.Name, mate.ValidationStatus }).ToArray(),
        placements = assemblyCompile?.Ir?.Placements.Select(placement => new
        {
            placement.InstanceStableId,
            status = placement.Status.ToString(),
            placement.FreeTranslations,
            placement.FreeRotations,
            authority = placement.Authority.ToString()
        }).ToArray()
    }
};

var json = JsonSerializer.Serialize(report, new JsonSerializerOptions { WriteIndented = true });
File.WriteAllText(Path.Combine(outputDirectory, "m19-probe-report.json"), json + Environment.NewLine);
Console.WriteLine(json);
return 0;

static (string Step, int FaceCount, int EdgeCount, int VertexCount) ExportPanel(PanelIr panel)
{
    var materialized = SurfacePatchBrepMaterializer.Materialize(
        panel.SurfaceConstruction.Support,
        panel.SurfaceConstruction.Evaluate,
        panel.Thickness ?? 0.1d);
    if (materialized.Body is null)
        throw new InvalidOperationException(string.Join("; ", materialized.Diagnostics.Select(diagnostic => diagnostic.Message)));
    var export = Step242Exporter.ExportBody(materialized.Body);
    if (!export.IsSuccess)
        throw new InvalidOperationException(string.Join("; ", export.Diagnostics.Select(diagnostic => diagnostic.Message)));
    return (
        export.Value,
        materialized.Body.Topology.Faces.Count(),
        materialized.Body.Topology.Edges.Count(),
        materialized.Body.Topology.Vertices.Count());
}

static string Sha(string text) => Convert.ToHexString(SHA256.HashData(Encoding.UTF8.GetBytes(text)));

static string ReportPath(string path) =>
    Path.GetRelativePath(Environment.CurrentDirectory, path).Replace('\\', '/');
