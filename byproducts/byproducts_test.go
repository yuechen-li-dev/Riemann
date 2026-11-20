package byproducts

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func fixture(src string) (sources, Review) {
	lines := strings.Split(src, "\n")
	line := lines[2]
	span := Span{"Pilot.claim", "src/Pilot.lean", 3, 3, fmt.Sprintf("%x", sha256.Sum256([]byte(line)))}
	record := Review{
		Slug: "pilot", Title: "fixture", Domain: "optimization", Declaration: span,
		Quantifiers: []string{"all instances"}, Interpretation: "fixture", HumanReview: "test",
		HypothesisStatus: Conditional, Assumptions: []string{"nonnegative"},
		Independence: NoRHInReviewedStatement, MathNovelty: KnownMath, FormalNovelty: NoveltyUnresolved,
		Verification: StructurallyReviewed,
		Comparisons:  []Comparison{{"literature", "reference", "https://example.org", Matches, "fixture"}},
		Open:         []string{"kernel replay"},
	}
	return sources{map[string]string{span.Path: src}, nil}, record
}

func TestEvidenceCannotBePromoted(t *testing.T) {
	good := "import Mathlib\nnamespace Pilot\ntheorem claim : (0 : Nat) ≤ 1 := by omega\nend Pilot\n"
	cases := []struct {
		name          string
		src           string
		edit          func(*Review)
		errorContains string
	}{
		{"RH hypothesis", "import Mathlib\nnamespace Pilot\ntheorem claim (h : RiemannHypothesis) : True := trivial\nend Pilot\n", nil, "RH independence"},
		{"conditional", "import Mathlib\nnamespace Pilot\ntheorem claim : False → True := by simp\nend Pilot\n", func(r *Review) { r.HypothesisStatus = Unconditional; r.Assumptions = nil }, "unconditional"},
		{"known without comparison", good, func(r *Review) { r.FormalNovelty = ExistingFormalization }, "needs a comparison"},
		{"invented novelty", good, func(r *Review) { r.FormalNovelty = "possibly_new_formalization" }, "novelty claims"},
		{"source mismatch", good, func(r *Review) { r.Declaration.SHA256 = strings.Repeat("f", 64) }, "source excerpt mismatch"},
		{"kernel replay", good, func(r *Review) { r.Verification = KernelReplayed }, "formal replay workflow"},
		{"missing import", strings.Replace(good, "import Mathlib", "import Absent.Module", 1), nil, "unresolved module"},
		{"missing proof use", good, func(r *Review) { r.ProofUses = []string{"absent"} }, "absent from selected body"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			fs, r := fixture(c.src)
			if c.edit != nil {
				c.edit(&r)
			}
			_, err := fs.extract(r)
			if err == nil || !strings.Contains(err.Error(), c.errorContains) {
				t.Fatalf("got %v, want %q", err, c.errorContains)
			}
		})
	}
}

func TestPinnedRevisionAndDistinctIdentities(t *testing.T) {
	if _, err := ReadCatalog(strings.NewReader(`{"revision":"bad","candidates":[]}`)); err == nil || !strings.Contains(err.Error(), "revision") {
		t.Fatalf("revision mismatch: %v", err)
	}
	source := func(ty string) string {
		return "import Mathlib\nnamespace Pilot\ntheorem claim : " + ty + " := by decide\nend Pilot\n"
	}
	a, ra := fixture(source("True"))
	b, rb := fixture(source("1 = 1"))
	one, err := a.extract(ra)
	if err != nil {
		t.Fatal(err)
	}
	two, err := b.extract(rb)
	if err != nil {
		t.Fatal(err)
	}
	if one.ID == two.ID {
		t.Fatal("different statements shared an identity")
	}
}

func TestOctRankingProtocol(t *testing.T) {
	input := map[string]ExternalStatement{"alpha named;": {}, "beta literature;": {}}
	ranked, err := parseRanking("alpha named;\r\n2\r\nbeta literature;\r\n1\r\n", input)
	if err != nil || ranked[0].Score != 2 {
		t.Fatalf("valid ranking: %v", err)
	}
	for _, bad := range []string{"alpha named;\n2\nalpha named;\n1", "alpha named;\n0\nbeta literature;\n1", "alpha named;\n2"} {
		if _, err := parseRanking(bad, input); err == nil {
			t.Fatalf("accepted %q", bad)
		}
	}
}

func TestTwoLetterCertificates(t *testing.T) {
	if n, err := checkTwoLetterDuality(); n != 405 || err != nil {
		t.Fatalf("%d certificates: %v", n, err)
	}
}

func TestB1PilotReplay(t *testing.T) {
	archive, oct := os.Getenv("A1_LEAN_ARCHIVE"), os.Getenv("B1_OCT")
	if archive == "" || oct == "" {
		t.Skip("set A1_LEAN_ARCHIVE and B1_OCT for corpus and Oct replay")
	}
	f, err := os.Open("b1_candidates.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	catalog, err := ReadCatalog(f)
	if err != nil {
		t.Fatal(err)
	}
	r, err := Discover(archive, catalog, oct, "../experiments/b1_discovery/main.oct")
	if err != nil {
		t.Fatal(err)
	}
	data, err := r.JSON()
	if err != nil {
		t.Fatal(err)
	}
	for name, got := range map[string]string{"inventory.json": string(data), "inventory.md": r.Markdown()} {
		want, err := os.ReadFile("../artifacts/b1/" + name)
		if err != nil || strings.ReplaceAll(string(want), "\r\n", "\n") != got {
			t.Fatalf("%s differs from replay: %v", name, err)
		}
	}
}
