package storage

import (
	"context"
	"io"

	"github.com/Oridak771/Vandal/schema"
)

// NewMySQLDatabase creates a new MySQL database.
func NewMySQLDatabase(host, port, user, password, dbname string) Database {
	return &mysqlDatabase{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

// mysqlDatabase is a basic implementation of the Database interface for MySQL.
type mysqlDatabase struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// Connect implements the Database interface.
func (d *mysqlDatabase) Connect(ctx context.Context) error {
	// To be implemented
	return nil
}

// GetSchema implements the Database interface.
func (d *mysqlDatabase) GetSchema(ctx context.Context) (*schema.Schema, error) {
	// To be implemented
	return nil, nil
}

// DumpTable implements the Database interface.
func (d *mysqlDatabase) DumpTable(ctx context.Context, tableName string) (io.Reader, error) {
	// To be implemented
	return nil, nil
}

// Restore implements the Database interface.
func (d *mysqlDatabase) Restore(ctx context.Context, in io.Reader) error {
	// To be implemented
	return nil
}
