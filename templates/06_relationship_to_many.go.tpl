{{- if .Table.IsJoinTable -}}
{{- else -}}
	{{- range $rel := .Table.ToManyRelationships -}}
		{{- $ltable := $.Aliases.Table $rel.Table -}}
		{{- $ftable := $.Aliases.Table $rel.ForeignTable -}}
		{{- $relAlias := $.Aliases.ManyRelationship $rel.ForeignTable $rel.Name $rel.JoinTable $rel.JoinLocalFKeyName -}}
		{{- $schemaForeignTable := .ForeignTable | $.SchemaTable -}}
		{{- $canSoftDelete := (getTable $.Tables .ForeignTable).CanSoftDelete }}
// {{$relAlias.Local}} retrieves all the {{.ForeignTable | singular}}'s {{$ftable.UpPlural}} with an executor
{{- if not (eq $relAlias.Local $ftable.UpPlural)}} via {{ range $col := $rel.ForeignColumns}} {{ $col.Name }} {{- end }}columns{{- end}}.
func (o *{{$ltable.UpSingular}}) {{$relAlias.Local}}(mods ...qm.QueryMod) {{$ftable.DownSingular}}Query {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

		{{if $rel.ToJoinTable -}}
	queryMods = append(queryMods,
		{{$schemaJoinTable := $rel.JoinTable | $.SchemaTable -}}
		qm.InnerJoin(`
			{{$schemaJoinTable}} ON
			{{- range $idx, $c := $rel.ForeignColumns -}}
				{{$schemaForeignTable}}.{{$rel.ForeignColumns[$idx].Name | $.Quotes}} = {{$schemaJoinTable}}.{{$rel.JoinForeignColumns[$idx].Name | $.Quotes}}
				{{ if $idx < len $rel }} AND {{ end -}}
			{{- end }}
			`),
		{{- range $idx, $col := $fkey.JoinLocalColumns -}}
		qm.Where("{{$schemaJoinTable}}.{{$rel.JoinLocalColumns[$idx].Name | $.Quotes}}=?", o.{{$ltable.Columns[$idx].Name $rel.Columns[$idx].Name}}),
		{{- end }}
	)
		{{else -}}
	queryMods = append(queryMods,
		{{- range $idx, $col := $fkey.ForeignColumns -}}
		qm.Where("{{$schemaForeignTable}}.{{$rel.ForeignColumns[$idx].Name | $.Quotes}}=?", o.{{$ltable.Columns[$idx].Name $rel.Columns[$idx].Name}}),
		{{- end }}
		{{if and $.AddSoftDeletes $canSoftDelete -}}
		qmhelper.WhereIsNull("{{$schemaForeignTable}}.{{"deleted_at" | $.Quotes}}"),
		{{- end}}
	)
		{{end}}

	query := {{$ftable.UpPlural}}(queryMods...)
	queries.SetFrom(query.Query, "{{$schemaForeignTable}}")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"{{$schemaForeignTable}}.*"})
	}

	return query
}

{{end -}}{{- /* range relationships */ -}}
{{- end -}}{{- /* if isJoinTable */ -}}
