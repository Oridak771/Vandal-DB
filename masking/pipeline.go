package masking

import (
	"context"

phantomdbv1alpha1 "github.com/vandal-db/vandal-db/apis/v1alpha1"
)

// Pipeline defines the interface for a masking pipeline.
type Pipeline interface {
	// Run executes the masking pipeline.
	Run(ctx context.Context) error
}

// NewPipeline creates a new masking pipeline.
func NewPipeline(dumper PostgresDumper, masker Masker, restorer PostgresRestorer, rules []phantomdbv1alpha1.MaskingRule) Pipeline {
	return &pipeline{
		dumper:   dumper,
		masker:   masker,
		restorer: restorer,
		rules:    rules,
	}
}

// pipeline is a basic implementation of the Pipeline interface.
type pipeline struct {
	dumper   PostgresDumper
	masker   Masker
	restorer PostgresRestorer
	rules    []phantomdbv1alpha1.MaskingRule
}

// Run implements the Pipeline interface.
func (p *pipeline) Run(ctx context.Context) error {
	// 1. Create a dump of the database.
	dumpReader, err := p.dumper.Dump()
	if err != nil {
		return err
	}

	// 2. Mask the data.
	maskedReader, err := p.masker.Mask(dumpReader, p.rules)
	if err != nil {
		return err
	}

	// 3. Restore the masked data.
	return p.restorer.Restore(maskedReader)
}
