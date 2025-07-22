package masking

import (
	"io"
	"os/exec"
)

// PostgresDumper defines the interface for a PostgreSQL dumper.
type PostgresDumper interface {
	// Dump creates a dump of a PostgreSQL database.
	Dump() (io.Reader, error)
}

// NewPostgresDumper creates a new PostgreSQL dumper.
func NewPostgresDumper(host, port, user, password, dbname string) PostgresDumper {
	return &postgresDumper{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

// postgresDumper is a basic implementation of the PostgresDumper interface.
type postgresDumper struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// Dump implements the PostgresDumper interface.
func (d *postgresDumper) Dump() (io.Reader, error) {
	cmd := exec.Command("pg_dump",
		"-h", d.host,
		"-p", d.port,
		"-U", d.user,
		"-d", d.dbname,
		"-F", "c", // Custom format
	)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+d.password)

	return cmd.StdoutPipe()
}

// PostgresRestorer defines the interface for a PostgreSQL restorer.
type PostgresRestorer interface {
	// Restore restores a dump of a PostgreSQL database.
	Restore(in io.Reader) error
}

// NewPostgresRestorer creates a new PostgreSQL restorer.
func NewPostgresRestorer(host, port, user, password, dbname string) PostgresRestorer {
	return &postgresRestorer{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		dbname:   dbname,
	}
}

// postgresRestorer is a basic implementation of the PostgresRestorer interface.
type postgresRestorer struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

// Restore implements the PostgresRestorer interface.
func (r *postgresRestorer) Restore(in io.Reader) error {
	cmd := exec.Command("pg_restore",
		"-h", r.host,
		"-p", r.port,
		"-U", r.user,
		"-d", r.dbname,
	)
	cmd.Env = append(cmd.Env, "PGPASSWORD="+r.password)
	cmd.Stdin = in

	return cmd.Run()
}
