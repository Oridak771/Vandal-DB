package masking

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
	case "creditCard":
		return &creditCardTransformer{}, nil
	case "name":
		return &nameTransformer{}, nil
	case "address":
		return &addressTransformer{}, nil
	case "dateTime":
		return &dateTimeTransformer{}, nil
	case "null":
		return &nullTransformer{}, nil
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

// creditCardTransformer implements the Transformer interface for credit card numbers.
type creditCardTransformer struct{}

// Transform implements the Transformer interface.
func (t *creditCardTransformer) Transform(value string) (string, error) {
	if len(value) < 4 {
		return "****", nil
	}
	return fmt.Sprintf("****-****-****-%s", value[len(value)-4:]), nil
}

// nameTransformer implements the Transformer interface for names.
type nameTransformer struct{}

// Transform implements the Transformer interface.
func (t *nameTransformer) Transform(value string) (string, error) {
	return gofakeit.Name(), nil
}

// addressTransformer implements the Transformer interface for addresses.
type addressTransformer struct{}

// Transform implements the Transformer interface.
func (t *addressTransformer) Transform(value string) (string, error) {
	return gofakeit.Address().Address, nil
}

// dateTimeTransformer implements the Transformer interface for date/time values.
type dateTimeTransformer struct{}

// Transform implements the Transformer interface.
func (t *dateTimeTransformer) Transform(value string) (string, error) {
	return gofakeit.Date().Format(time.RFC3339), nil
}

// nullTransformer implements the Transformer interface for nullifying values.
type nullTransformer struct{}

// Transform implements the Transformer interface.
func (t *nullTransformer) Transform(value string) (string, error) {
	return "", nil
}
