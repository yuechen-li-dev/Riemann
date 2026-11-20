package byproducts

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/yuechen-li-dev/Riemann/compiler"
	"github.com/yuechen-li-dev/Riemann/semantic"
)

// CyclicGraph records provenance through the exact verifier. The existing
// compiler has no rule importing verified external Lean objects, so that
// portion is a source trace with open obligations.
func CyclicGraph(c CyclicCertificate) (*compiler.Graph, error) {
	objective, err := ReplayCyclic(c)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}
	instance := fmt.Sprintf("l=%d d=%d certificate=%x", c.Base, c.Length, sha256.Sum256(payload))
	g := compiler.NewGraph()
	add := func(tag, scope, evidence string, kind semantic.EvidenceKind, src semantic.Reference, parents ...semantic.ClaimID) error {
		id := semantic.ClaimID("b2." + tag)
		claim := semantic.Claim{
			ID: id, Proposition: semantic.ExternalStatement{Predicate: string(id), Revision: Revision, Instance: scope},
			Exactness: semantic.Exact, Evidence: []semantic.Evidence{{Kind: kind, Source: src, Note: evidence}},
			Provenance: semantic.Provenance{Kind: semantic.AuthoredProvenance, Source: src},
		}
		if len(parents) == 0 {
			return g.AddClaim(claim)
		}
		tr := semantic.TransformationID("b2.derive." + tag)
		claim.Provenance = semantic.Provenance{Kind: semantic.DerivedProvenance, Parents: parents, Transformation: tr, Source: src}
		if err := g.AddClaim(claim); err != nil {
			return err
		}
		return g.AddTransformation(compiler.Transformation{ID: tr, Pass: "B2", From: parents[0], Premises: parents[1:], To: id, Relation: compiler.Implies, Provenance: src})
	}
	source := semantic.Reference{Kind: semantic.StandardReference, Citation: "Pinned cyclic LP", URI: Repository + "/blob/" + Revision + "/src/Formalism/IncomingBounds/IncomingBound.lean"}
	checker := semantic.Reference{Kind: semantic.StandardReference, Citation: "B2 exact Go replay", URI: "byproducts/cyclic.go"}
	proof := semantic.Reference{Kind: semantic.StandardReference, Citation: "Pinned Lean declarations", URI: Repository + "/blob/" + Revision + "/src/Formalism/IncomingBounds/FlowDegreeBounds/MaxFlowValueAttained.lean"}
	if err := add("instance", instance, "Fixed source parameters", semantic.DefinitionEvidence, source); err != nil {
		return nil, err
	}
	if err := add("primal", instance, "Checked all nonnegativity and capacity inequalities", semantic.CertifiedComputationEvidence, checker, "b2.instance"); err != nil {
		return nil, err
	}
	if err := add("dual", instance, "Checked all covering and box inequalities", semantic.CertifiedComputationEvidence, checker, "b2.instance"); err != nil {
		return nil, err
	}
	if err := add("certificate", instance, "Recomputed equal objectives "+objective.RatString(), semantic.CertifiedComputationEvidence, checker, "b2.primal", "b2.dual"); err != nil {
		return nil, err
	}
	if err := add("optimality", instance, "Weak duality for the exact certificate", semantic.DerivedEvidence, checker, "b2.certificate"); err != nil {
		return nil, err
	}
	if err := add("universal", "", "Source proposition FlowAttained and LemmaFa", semantic.UnverifiedConjectureEvidence, source); err != nil {
		return nil, err
	}
	if err := add("lean-proof", "", "Formalism.IncomingBounds.maxFlowValueAttained and strongDuality", semantic.UnverifiedConjectureEvidence, proof, "b2.universal"); err != nil {
		return nil, err
	}
	if err := add("formal-status", "", "External checking evidence belongs to the B2 Lean transcript; compiler import is open", semantic.UnverifiedConjectureEvidence, proof, "b2.lean-proof"); err != nil {
		return nil, err
	}
	return g, nil
}
