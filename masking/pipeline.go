package masking

import (
	"context"

	vandalv1alpha1 "github.com/Oridak771/Vandal/apis/v1alpha1"
	"github.com/Oridak771/Vandal/storage"
	"golang.org/x/sync/errgroup"
)

// Pipeline defines the interface for a masking pipeline.
type Pipeline interface {
	// Run executes the masking pipeline.
	Run(ctx context.Context) error
}

// NewPipeline creates a new masking pipeline.
func NewPipeline(db storage.Database, masker Masker, rules []vandalv1alpha1.MaskingRule) Pipeline {
	return &pipeline{
		db:     db,
		masker: masker,
		rules:  rules,
	}
}

// pipeline is a basic implementation of the Pipeline interface.
type pipeline struct {
	db     storage.Database
	masker Masker
	rules  []vandalv1alpha1.MaskingRule
}

// Run implements the Pipeline interface.
func (p *pipeline) Run(ctx context.Context) error {
	// 1. Get the database schema.
	schema, err := p.db.GetSchema(ctx)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	for _, table := range schema.Tables {
		table := table // https://golang.org/doc/faq#closures_and_goroutines
		g.Go(func() error {
			// 1. Create a dump of the table.
			dumpReader, err := p.db.DumpTable(ctx, table.Name)
			if err != nil {
				return err
			}

			// 2. Mask the data.
			maskedReader, err := p.masker.Mask(dumpReader, p.rules, schema)
			if err != nil {
				return err
			}

			// 3. Restore the masked data.
			return p.db.Restore(ctx, maskedReader)
		})
	}

	return g.Wait()
}
