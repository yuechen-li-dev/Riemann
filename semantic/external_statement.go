package semantic

import "fmt"

// ExternalStatement identifies an externally defined predicate at a pinned
// revision. Instance identifies a particular realization when nonempty.
// Its type and identity convey no trust in the external proof.
type ExternalStatement struct {
	Predicate string `json:"predicate"`
	Revision  string `json:"revision"`
	Instance  string `json:"instance,omitempty"`
}

func (ExternalStatement) isProposition()        {}
func (ExternalStatement) Kind() PropositionKind { return "external_statement" }
func (s ExternalStatement) Describe() string {
	if s.Instance != "" {
		return fmt.Sprintf("%s on %s", s.Predicate, s.Instance)
	}
	return s.Predicate
}
