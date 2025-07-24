package schema

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Schema defines the structure of a database schema.
type Schema struct {
	Tables []Table
}

// Table defines the structure of a database table.
type Table struct {
	Name    string
	Columns []Column
}

// Column defines the structure of a database column.
type Column struct {
	Name    string
	Type    string
	IsNullable bool
	IsPrimaryKey bool
	IsForeignKey bool
	ForeignKeyTable string
	ForeignKeyColumn string
}

// GetSchema fetches the schema of a PostgreSQL database.
func GetSchema(host, port, user, password, dbname string) (*Schema, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Get tables
	rows, err := db.Query("SELECT tablename FROM pg_catalog.pg_tables WHERE schemaname != 'pg_catalog' AND schemaname != 'information_schema'")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []Table
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, Table{Name: tableName})
	}

	// Get columns for each table
	for i, table := range tables {
		rows, err := db.Query(`
			SELECT
				c.column_name,
				c.data_type,
				c.is_nullable,
				(SELECT COUNT(*) FROM information_schema.table_constraints tc JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name WHERE tc.table_name = c.table_name AND kcu.column_name = c.column_name AND tc.constraint_type = 'PRIMARY KEY') > 0 AS is_primary_key,
				(SELECT COUNT(*) FROM information_schema.table_constraints tc JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name WHERE tc.table_name = c.table_name AND kcu.column_name = c.column_name AND tc.constraint_type = 'FOREIGN KEY') > 0 AS is_foreign_key,
				(SELECT ccu.table_name FROM information_schema.referential_constraints rc JOIN information_schema.key_column_usage kcu ON rc.constraint_name = kcu.constraint_name JOIN information_schema.constraint_column_usage ccu ON rc.unique_constraint_name = ccu.constraint_name WHERE kcu.table_name = c.table_name AND kcu.column_name = c.column_name) AS foreign_key_table,
				(SELECT ccu.column_name FROM information_schema.referential_constraints rc JOIN information_schema.key_column_usage kcu ON rc.constraint_name = kcu.constraint_name JOIN information_schema.constraint_column_usage ccu ON rc.unique_constraint_name = ccu.constraint_name WHERE kcu.table_name = c.table_name AND kcu.column_name = c.column_name) AS foreign_key_column
			FROM
				information_schema.columns c
			WHERE
				c.table_name = $1
		`, table.Name)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		var columns []Column
		for rows.Next() {
			var column Column
			var isNullable string
			if err := rows.Scan(&column.Name, &column.Type, &isNullable, &column.IsPrimaryKey, &column.IsForeignKey, &column.ForeignKeyTable, &column.ForeignKeyColumn); err != nil {
				return nil, err
			}
			column.IsNullable = (isNullable == "YES")
			columns = append(columns, column)
		}
		tables[i].Columns = columns
	}

	return &Schema{Tables: tables}, nil
}
