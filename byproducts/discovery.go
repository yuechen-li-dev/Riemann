package byproducts

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Report struct {
	Revision        string              `json:"revision"`
	ArchiveSHA256   string              `json:"archive_sha256"`
	LeanToolchain   string              `json:"lean_toolchain"`
	MathlibRevision string              `json:"mathlib_revision"`
	RankingEngine   string              `json:"ranking_engine"`
	RankerSHA256    string              `json:"ranker_sha256"`
	Statements      []ExternalStatement `json:"statements"`
}

func Discover(archive string, catalog Catalog, oct, program string) (Report, error) {
	records, err := Extract(archive, catalog)
	if err != nil {
		return Report{}, err
	}
	oct, err = exec.LookPath(oct)
	if err != nil {
		return Report{}, err
	}
	oct, err = filepath.Abs(oct)
	if err != nil {
		return Report{}, err
	}
	code, err := os.ReadFile(program)
	if err != nil {
		return Report{}, err
	}
	dir, err := os.MkdirTemp("", "riemann-b1-")
	if err != nil {
		return Report{}, err
	}
	defer os.RemoveAll(dir)
	if err := os.WriteFile(filepath.Join(dir, "main.oct"), code, 0600); err != nil {
		return Report{}, err
	}

	rows := make([]string, 0, len(records))
	input := make(map[string]ExternalStatement)
	for _, r := range records {
		if r.Declaration.Name == "Formalism.IncomingBounds.lemmaFa" {
			count, err := checkTwoLetterDuality()
			if err != nil {
				return Report{}, err
			}
			r.Computation = FiniteReplay
			r.ComputationDetail = fmt.Sprintf("%d exact two-letter primal/dual certificates, bases 3–7 and integer budgets 0–8.", count)
		}
		features := []string{"literature"}
		if namedStatement(r.Statement) {
			features = append(features, "named")
		}
		if mathlibOnly(r.DefinitionImports) {
			features = append(features, "mathlib-only")
		}
		if r.Computation == FiniteReplay {
			features = append(features, "extractable")
		}
		if r.FormalNovelty == NoveltyUnresolved {
			features = append(features, "comparison-open")
		}
		if r.ExtraProofPremise {
			features = append(features, "proof-premise")
		}
		r.RankingFeatures = features
		row := r.Slug + " " + strings.Join(features, "; ") + ";"
		rows = append(rows, row)
		input[row] = r
	}
	if err := os.WriteFile(filepath.Join(dir, "candidates.txt"), []byte(strings.Join(rows, "\n")+"\n"), 0600); err != nil {
		return Report{}, err
	}

	cmd := exec.Command(oct, "run", "main.oct")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return Report{}, fmt.Errorf("Oct ranker: %w\n%s", err, output)
	}
	manifest, err := os.ReadFile(filepath.Join(dir, "ranking.txt"))
	if err != nil {
		return Report{}, fmt.Errorf("read Oct ranking: %w", err)
	}
	ranked, err := parseRanking(string(manifest), input)
	if err != nil {
		return Report{}, err
	}
	return Report{
		Revision: Revision, ArchiveSHA256: ArchiveSHA256,
		LeanToolchain: LeanToolchain, MathlibRevision: MathlibRevision,
		RankingEngine: "Oct b1_discovery",
		RankerSHA256:  fmt.Sprintf("%x", sha256.Sum256([]byte(strings.ReplaceAll(string(code), "\r\n", "\n")))),
		Statements:    ranked,
	}, nil
}

func namedStatement(s string) bool {
	_, t, _ := strings.Cut(s, ": ")
	return len(strings.Fields(t)) == 1
}

func mathlibOnly(sets map[string][]string) bool {
	if len(sets) == 0 {
		return false
	}
	for _, modules := range sets {
		for _, m := range modules {
			if m != "Mathlib" && !strings.HasPrefix(m, "Mathlib.") {
				return false
			}
		}
	}
	return true
}

func parseRanking(output string, input map[string]ExternalStatement) ([]ExternalStatement, error) {
	lines := strings.Split(strings.ReplaceAll(strings.TrimSpace(output), "\r\n", "\n"), "\n")
	if len(lines) != len(input)*2 {
		return nil, fmt.Errorf("unexpected Oct ranking: %q", output)
	}
	var ranked []ExternalStatement
	seen := make(map[string]bool)
	for i := 0; i < len(lines); i += 2 {
		row := lines[i]
		r, ok := input[row]
		score, err := strconv.Atoi(lines[i+1])
		if !ok || seen[row] || err != nil || (len(ranked) > 0 && score > ranked[len(ranked)-1].Score) {
			return nil, fmt.Errorf("invalid Oct row %q, score %q", row, lines[i+1])
		}
		seen[row] = true
		r.Score = score
		ranked = append(ranked, r)
	}
	return ranked, nil
}

func (r Report) JSON() ([]byte, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(data, '\n'), nil
}

func (r Report) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# B1 ranked inventory\n\nRevision `%s`. Oct scores prioritize review. Proof status is source-level structural review for every record.\n\n", r.Revision)
	for i, s := range r.Statements {
		fmt.Fprintf(&b, "## %d. %s\n\n%s\n\nOct score **%d** (%s). %s.\n\n", i+1, s.Title, s.Interpretation, s.Score, strings.Join(s.RankingFeatures, ", "), s.Domain)
		fmt.Fprintf(&b, "[Pinned declaration %s](%s).\n\n", s.Declaration.Name, s.Declaration.Link())
		fmt.Fprintf(&b, "Formalization novelty **%s**. Independent computation **%s**.", s.FormalNovelty, s.Computation)
		if s.ComputationDetail != "" {
			fmt.Fprintf(&b, " %s", s.ComputationDetail)
		}
		b.WriteString("\n\n")
		for _, c := range s.Comparisons {
			fmt.Fprintf(&b, "- **%s** (%s). %s [%s](%s)\n", c.Kind, c.Fidelity, c.Finding, c.Title, c.URL)
		}
		for _, q := range s.Open {
			fmt.Fprintf(&b, "- Open. %s\n", q)
		}
		b.WriteByte('\n')
	}
	return b.String()
}
