package storage

import (
	"context"
	"io"
	"os/exec"

	"github.com/Oridak771/Vandal/schema"
)

// NewPostgresDatabase creates a new PostgreSQL database.
func NewPostgresDatabase(host, port, user, password, dbname string) Database {
	return &postgresDatabase{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

// postgresDatabase is a basic implementation of the Database interface for PostgreSQL.
type postgresDatabase struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// Connect implements the Database interface.
func (d *postgresDatabase) Connect(ctx context.Context) error {
	// The connection is implicitly tested by other methods.
	return nil
}

// GetSchema implements the Database interface.
func (d *postgresDatabase) GetSchema(ctx context.Context) (*schema.Schema, error) {
	return schema.GetSchema(d.host, d.port, d.user, d.password, d.dbname)
}

// DumpTable implements the Database interface.
func (d *postgresDatabase) DumpTable(ctx context.Context, tableName string) (io.Reader, error) {
	cmd := exec.CommandContext(ctx, "pg_dump",
		"-h", d.host,
		"-p", d.port,
		"-U", d.user,
		"-d", d.dbname,
		"-t", tableName,
		"-F", "c", // Custom format
	)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+d.password)

	return cmd.StdoutPipe()
}

// Restore implements the Database interface.
func (d *postgresDatabase) Restore(ctx context.Context, in io.Reader) error {
	cmd := exec.CommandContext(ctx, "pg_restore",
		"-h", d.host,
		"-p", d.port,
		"-U", d.user,
		"-d", d.dbname,
	)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+d.password)
	cmd.Stdin = in

	return cmd.Run()
}
