package storage

import (
	"context"
	"io"

	"github.com/Oridak771/Vandal/schema"
)

// Database defines the interface for a database.
type Database interface {
	// Connect connects to the database.
	Connect(ctx context.Context) error
	// GetSchema returns the schema of the database.
	GetSchema(ctx context.Context) (*schema.Schema, error)
	// DumpTable creates a dump of a single table.
	DumpTable(ctx context.Context, tableName string) (io.Reader, error)
	// Restore restores a dump of a database.
	Restore(ctx context.Context, in io.Reader) error
}
