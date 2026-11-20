package byproducts

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
)

type sources struct {
	bodies map[string]string
	paths  map[string]bool
}

func Extract(archive string, catalog Catalog) ([]ExternalStatement, error) {
	if catalog.Revision != Revision || len(catalog.Candidates) == 0 || len(catalog.Candidates) > 32 {
		return nil, fmt.Errorf("invalid revision or candidate count")
	}
	want := map[string]bool{"lean-toolchain": true, "lakefile.toml": true, "lake-manifest.json": true}
	for _, r := range catalog.Candidates {
		for _, s := range reviewSpans(r) {
			if !strings.HasPrefix(s.Path, "src/") || !strings.HasSuffix(s.Path, ".lean") || path.Clean(s.Path) != s.Path {
				return nil, fmt.Errorf("invalid source path %q", s.Path)
			}
			want[s.Path] = true
		}
	}
	fs, err := readArchive(archive, want)
	if err != nil {
		return nil, err
	}
	result := make([]ExternalStatement, 0, len(catalog.Candidates))
	names, slugs := make(map[string]bool), make(map[string]bool)
	for _, r := range catalog.Candidates {
		if names[r.Declaration.Name] || slugs[r.Slug] {
			return nil, fmt.Errorf("duplicate declaration or slug %s", r.Slug)
		}
		names[r.Declaration.Name], slugs[r.Slug] = true, true
		record, err := fs.extract(r)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", r.Slug, err)
		}
		result = append(result, record)
	}
	return result, nil
}

func reviewSpans(r Review) []Span {
	spans := []Span{r.Declaration}
	spans = append(spans, r.Definitions...)
	return append(spans, r.Context...)
}

func readArchive(filename string, want map[string]bool) (sources, error) {
	fs := sources{make(map[string]string), make(map[string]bool)}
	f, err := os.Open(filename)
	if err != nil {
		return fs, err
	}
	defer f.Close()
	h := sha256.New()
	raw := io.TeeReader(f, h)
	gz, err := gzip.NewReader(raw)
	if err != nil {
		return fs, err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	prefix := "riemann-hypothesis-" + Revision + "/"
	for {
		item, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fs, err
		}
		if item.Typeflag != tar.TypeReg {
			continue
		}
		name, ok := strings.CutPrefix(item.Name, prefix)
		if !ok {
			return fs, fmt.Errorf("archive path outside pinned revision: %s", item.Name)
		}
		if fs.paths[name] {
			return fs, fmt.Errorf("duplicate archive path %s", name)
		}
		fs.paths[name] = true
		if want[name] {
			if item.Size > 2<<20 {
				return fs, fmt.Errorf("pilot source file too large: %s", name)
			}
			body, err := io.ReadAll(tr)
			if err != nil {
				return fs, err
			}
			fs.bodies[name] = string(body)
		}
	}
	if _, err := io.Copy(io.Discard, gz); err != nil {
		return fs, err
	}
	if _, err := io.Copy(io.Discard, raw); err != nil {
		return fs, err
	}
	if fmt.Sprintf("%x", h.Sum(nil)) != ArchiveSHA256 {
		return fs, fmt.Errorf("archive SHA-256 does not match pinned revision")
	}
	for name := range want {
		if _, ok := fs.bodies[name]; !ok {
			return fs, fmt.Errorf("missing archive member %s", name)
		}
	}
	if err := fs.checkToolchain(); err != nil {
		return fs, err
	}
	return fs, nil
}

func (fs sources) checkToolchain() error {
	var manifest struct {
		Packages []struct {
			Name string `json:"name"`
			Rev  string `json:"rev"`
		} `json:"packages"`
	}
	if err := json.Unmarshal([]byte(fs.bodies["lake-manifest.json"]), &manifest); err != nil {
		return err
	}
	for _, pkg := range manifest.Packages {
		if pkg.Name == "mathlib" && pkg.Rev == MathlibRevision &&
			strings.TrimSpace(fs.bodies["lean-toolchain"]) == LeanToolchain {
			return nil
		}
	}
	return fmt.Errorf("pinned Lean or Mathlib version mismatch")
}

func excerpt(s Span, src string) (string, error) {
	lines := strings.Split(src, "\n")
	if s.First < 1 || s.Last < s.First || s.Last > len(lines) || s.Last-s.First > 80 {
		return "", fmt.Errorf("invalid span %s:%d–%d", s.Path, s.First, s.Last)
	}
	part := strings.Join(lines[s.First-1:s.Last], "\n")
	if fmt.Sprintf("%x", sha256.Sum256([]byte(part))) != s.SHA256 {
		return "", fmt.Errorf("source excerpt mismatch at %s:%d", s.Path, s.First)
	}
	return part, nil
}

var importLine = regexp.MustCompile(`(?m)^import ([^\s]+)[ \t]*(?:--[^\n]*)?\r?$`)
var proofEnd = regexp.MustCompile(`(?m)^(?:theorem|lemma|def|private|end|variable|section|namespace|set_option)\b`)

func (fs sources) imports(s Span) ([]string, error) {
	var imports []string
	for _, match := range importLine.FindAllStringSubmatch(fs.bodies[s.Path], -1) {
		module := match[1]
		if module != "Mathlib" && !strings.HasPrefix(module, "Mathlib.") && !fs.paths["src/"+strings.ReplaceAll(module, ".", "/")+".lean"] {
			return nil, fmt.Errorf("%s imports unresolved module %s", s.Path, module)
		}
		imports = append(imports, module)
	}
	if len(imports) == 0 {
		return nil, fmt.Errorf("no imports extracted for %s", s.Path)
	}
	return imports, nil
}

func (fs sources) extract(r Review) (ExternalStatement, error) {
	var parts []string
	for _, s := range reviewSpans(r) {
		part, err := excerpt(s, fs.bodies[s.Path])
		if err != nil {
			return ExternalStatement{}, err
		}
		parts = append(parts, part)
	}
	header, opener, ok := strings.Cut(parts[0], " :=")
	if !ok {
		return ExternalStatement{}, fmt.Errorf("missing proof delimiter")
	}
	src := fs.bodies[r.Declaration.Path]
	if !declarationMatches(r.Declaration.Name, header, strings.Join(strings.Split(src, "\n")[:r.Declaration.First-1], "\n")+"\n") {
		return ExternalStatement{}, fmt.Errorf("theorem or namespace does not match %s", r.Declaration.Name)
	}
	parts[0] = header
	if err := validateReview(r, strings.Join(parts, "\n")); err != nil {
		return ExternalStatement{}, err
	}
	imports, err := fs.imports(r.Declaration)
	if err != nil {
		return ExternalStatement{}, err
	}
	defImports := make(map[string][]string)
	for i, def := range r.Definitions {
		part := parts[i+1]
		name := regexp.MustCompile(`\b` + regexp.QuoteMeta(def.Name) + `\b`)
		if !name.MatchString(strings.SplitN(part, "\n", 2)[0]) {
			return ExternalStatement{}, fmt.Errorf("definition anchor mismatch for %s", def.Name)
		}
		referenced := false
		for j, other := range parts {
			if j != i+1 && name.MatchString(other) {
				referenced = true
				break
			}
		}
		if !referenced {
			return ExternalStatement{}, fmt.Errorf("definition %s has no use in reviewed excerpts", def.Name)
		}
		if _, ok := defImports[def.Path]; ok {
			continue
		}
		di, err := fs.imports(def)
		if err != nil {
			return ExternalStatement{}, err
		}
		defImports[def.Path] = di
	}
	lines := strings.Split(src, "\n")
	body := opener + "\n" + strings.Join(lines[r.Declaration.Last:], "\n")
	if end := proofEnd.FindStringIndex(body); end != nil {
		body = body[:end[0]]
	}
	for _, use := range r.ProofUses {
		if use == "" || !regexp.MustCompile(`\b`+regexp.QuoteMeta(use)+`\b`).MatchString(body) {
			return ExternalStatement{}, fmt.Errorf("proof use %s absent from selected body", use)
		}
	}
	id := fmt.Sprintf("lean:%s:%s:%x", Revision, r.Declaration.Name, sha256.Sum256([]byte(strings.Join(parts, "\n"))))
	return ExternalStatement{ID: id, Repository: Repository, Revision: Revision, Statement: header, SourceFidelity: Matches, Review: r,
		Imports: imports, DefinitionImports: defImports, Computation: NotRun}, nil
}

// A Lean declaration can have a dotted local name inside its namespace.
func declarationMatches(name, header, prefix string) bool {
	for i := strings.LastIndexByte(name, '.'); i > 0; i = strings.LastIndexByte(name[:i], '.') {
		decl := regexp.MustCompile(`^theorem ` + regexp.QuoteMeta(name[i+1:]) + `[\s:{]`)
		if decl.MatchString(header) && strings.Contains("\n"+prefix, "\nnamespace "+name[:i]+"\n") {
			return true
		}
	}
	return false
}
