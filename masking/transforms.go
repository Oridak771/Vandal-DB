package masking

import (
	"crypto/sha256"
	"fmt"
)

// Transformer defines the interface for a data transformer.
type Transformer interface {
	// Transform applies a transformation to a value.
	Transform(value string) (string, error)
}

// NewTransformer creates a new transformer for the given rule.
func NewTransformer(rule string) (Transformer, error) {
	switch rule {
	case "hash":
		return &hashTransformer{}, nil
	case "redact":
		return &redactTransformer{}, nil
	case "synthesize":
		return &synthesizeTransformer{}, nil
	default:
		return nil, fmt.Errorf("unknown transformation rule: %s", rule)
	}
}

// hashTransformer implements the Transformer interface for hashing.
type hashTransformer struct{}

// Transform implements the Transformer interface.
func (t *hashTransformer) Transform(value string) (string, error) {
	h := sha256.New()
	h.Write([]byte(value))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// redactTransformer implements the Transformer interface for redacting.
type redactTransformer struct{}

// Transform implements the Transformer interface.
func (t *redactTransformer) Transform(value string) (string, error) {
	return "REDACTED", nil
}

// synthesizeTransformer implements the Transformer interface for synthesizing.
type synthesizeTransformer struct{}

// Transform implements the Transformer interface.
func (t *synthesizeTransformer) Transform(value string) (string, error) {
	// This is a placeholder implementation.
	// In a real implementation, this would generate synthetic data.
	return "SYNTHESIZED", nil
}
