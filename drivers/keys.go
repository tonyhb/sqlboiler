package drivers

import "fmt"

// PrimaryKey represents a primary key constraint in a database
type PrimaryKey struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
}

// ForeignKey represents a foreign key constraint in a database
type ForeignKey struct {
	Table   string             `json:"table"`
	Name    string             `json:"name"`
	Columns []ForeignKeyColumn `json:"columns"`

	ForeignTable   string             `json:"foreign_table"`
	ForeignColumns []ForeignKeyColumn `json:"foreign_columns"`
}

// ForeignKeyColumn represents information about a column for a foreign key. It is
// a subset of a Column struct.
type ForeignKeyColumn struct {
	Name     string `json:"column_name"`
	Nullable bool   `json:"column_nullable"`
	Unique   bool   `json:"column_unique"`
}

// SQLColumnDef formats a column name and type like an SQL column definition.
type SQLColumnDef struct {
	Name string
	Type string
}

// String for fmt.Stringer
func (s SQLColumnDef) String() string {
	return fmt.Sprintf("%s %s", s.Name, s.Type)
}

// SQLColumnDefs has small helper functions
type SQLColumnDefs []SQLColumnDef

// Names returns all the names
func (s SQLColumnDefs) Names() []string {
	names := make([]string, len(s))

	for i, sqlDef := range s {
		names[i] = sqlDef.Name
	}

	return names
}

// Types returns all the types
func (s SQLColumnDefs) Types() []string {
	types := make([]string, len(s))

	for i, sqlDef := range s {
		types[i] = sqlDef.Type
	}

	return types
}

// SQLColDefinitions creates a definition in sql format for a column
func SQLColDefinitions(cols []Column, names []string) SQLColumnDefs {
	ret := make([]SQLColumnDef, len(names))

	for i, n := range names {
		for _, c := range cols {
			if n != c.Name {
				continue
			}

			ret[i] = SQLColumnDef{Name: n, Type: c.Type}
		}
	}

	return ret
}
