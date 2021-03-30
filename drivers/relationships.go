package drivers

// ToOneRelationship describes a relationship between two tables where the local
// table has no id, and the foreign table has an id that matches a column in the
// local table, that column can also be unique which changes the dynamic into a
// one-to-one style, not a to-many.
type ToOneRelationship struct {
	Name string `json:"name"`

	Table   string             `json:"table"`
	Columns []ForeignKeyColumn `json:"columns"`

	ForeignTable   string             `json:"foreign_table"`
	ForeignColumns []ForeignKeyColumn `json:"foreign_columns"`
}

// ToManyRelationship describes a relationship between two tables where the
// local table has no id, and either:
// 1. the foreign table has an id that matches a column in the local table.
// 2. the foreign table has columns that matches the primary key in the local table,
//    in the case of composite foreign keys.
type ToManyRelationship struct {
	Name string `json:"name"`

	Table   string             `json:"table"`
	Columns []ForeignKeyColumn `json:"columns"`

	ForeignTable   string             `json:"foreign_table"`
	ForeignColumns []ForeignKeyColumn `json:"foreign_columns"`

	ToJoinTable bool   `json:"to_join_table"`
	JoinTable   string `json:"join_table"`

	JoinLocalFKeyName string             `json:"join_local_fkey_name"`
	JoinLocalColumns  []ForeignKeyColumn `json:"join_local_columns"`

	JoinForeignFKeyName string             `json:"join_foreign_fkey_name"`
	JoinForeignColumns  []ForeignKeyColumn `json:"join_foreign_columns"`
}

// ToOneRelationships relationship lookups
// Input should be the sql name of a table like: videos
func ToOneRelationships(table string, tables []Table) []ToOneRelationship {
	localTable := GetTable(tables, table)
	return toOneRelationships(localTable, tables)
}

// ToManyRelationships relationship lookups
// Input should be the sql name of a table like: videos
func ToManyRelationships(table string, tables []Table) []ToManyRelationship {
	localTable := GetTable(tables, table)
	return toManyRelationships(localTable, tables)
}

func toOneRelationships(table Table, tables []Table) []ToOneRelationship {
	var relationships []ToOneRelationship

	for _, t := range tables {
		for _, f := range t.FKeys {

			allUnique := true
			for _, cols := range f.Columns {
				if !cols.Unique {
					allUnique = false
					break
				}
			}

			if allUnique && !t.IsJoinTable {
				relationships = append(relationships, buildToOneRelationship(table, f, t, tables))
			}

		}
	}

	return relationships
}

func toManyRelationships(table Table, tables []Table) []ToManyRelationship {
	var relationships []ToManyRelationship

	for _, t := range tables {
		for _, f := range t.FKeys {

			allUnique := true
			for _, cols := range f.Columns {
				if !cols.Unique {
					allUnique = false
					break
				}
			}

			if f.ForeignTable == table.Name && (t.IsJoinTable || !allUnique) {
				relationships = append(relationships, buildToManyRelationship(table, f, t, tables))
			}
		}
	}

	return relationships
}

func buildToOneRelationship(localTable Table, foreignKey ForeignKey, foreignTable Table, tables []Table) ToOneRelationship {
	return ToOneRelationship{
		Name:    foreignKey.Name,
		Table:   localTable.Name,
		Columns: foreignKey.ForeignColumns,

		ForeignTable:   foreignTable.Name,
		ForeignColumns: foreignKey.Columns,
	}
}

func buildToManyRelationship(localTable Table, foreignKey ForeignKey, foreignTable Table, tables []Table) ToManyRelationship {
	if !foreignTable.IsJoinTable {
		return ToManyRelationship{
			Name:           foreignKey.Name,
			Table:          localTable.Name,
			Columns:        foreignKey.ForeignColumns,
			ForeignTable:   foreignTable.Name,
			ForeignColumns: foreignKey.Columns,
			ToJoinTable:    false,
		}
	}

	relationship := ToManyRelationship{
		Table:   localTable.Name,
		Columns: foreignKey.ForeignColumns,

		ToJoinTable: true,
		JoinTable:   foreignTable.Name,

		JoinLocalFKeyName: foreignKey.Name,
		JoinLocalColumns:  foreignKey.Columns,
	}

	for _, fk := range foreignTable.FKeys {
		if fk.Name == foreignKey.Name {
			continue
		}

		relationship.JoinForeignFKeyName = fk.Name
		relationship.JoinForeignColumns = fk.Columns

		relationship.ForeignTable = fk.ForeignTable
		relationship.ForeignColumns = fk.ForeignColumns
	}

	return relationship
}
