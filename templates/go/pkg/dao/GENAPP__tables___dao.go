package dao

import (
	"time"

	"github.com/guregu/null"
	"github.com/jmoiron/sqlx"
	"github.com/Masterminds/squirrel"
)

type {{ table.name|string.capitalize }} struct {
	{{~ for column in table.columns ~}}
	C_{{ column.name }} {{ column.go_type }} `db:"{{ column.name }}" json:"{{ column.name }}"`
	{{~ end ~}}
}

type {{ table.name|string.capitalize }}PaginatedResponse struct {
	Total uint64
	Data []{{ table.name|string.capitalize }}
}

func (d DAO) {{ table.name|string.capitalize }}GetMany(where squirrel.Sqlizer, p Pagination) (*{{ table.name|string.capitalize }}PaginatedResponse, error) {
	var filter string
	var args []interface{}
	if where != nil {
		var err error
		filter, args, err = where.ToSql()
		if err != nil {
			return nil, err
		}

		filter = "WHERE " + filter
	}

	query := fmt.Sprintf(`
SELECT
  {{~ for column in table.columns ~}}
  "{{ column.name }}",
  {{~ end ~}}
  COUNT() OVER () AS __total
FROM
  "{{table.name}}"
WHERE
  %s
ORDER BY
  %s
OFFSET
  %d
LIMIT
  %d`, filter, p.Order, p.Offset, p.Limit)
	rows, err := d.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	var response {{ table.name|string.capitalize }}PaginatedResponse
	for rows.Next() {
		var row struct {
			{{ table.name|string.capitalize }}
			Total uint64 `db:"__total"`
		}
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}

		response.Total = row.Total
		response.Data = append(response.Data, row.{{ table.name|string.capitalize }})
	}

	return &response, err
}

func (d DAO) {{ table.name|string.capitalize }}Insert(body *{{ table.name|string.capitalize }}) error {
	row := d.db.QueryRowx(`
INSERT INTO {{ table.name }} (
  {{~ for column in table.columns ~}}
  {{~ if column.auto_increment
         continue
        end ~}}
  "{{ column.name }}"{{ if !for.last }},{{ end }}
  {{~ end ~}})
VALUES (
  {{~ for column in table.columns ~}}
  {{~ if column.auto_increment
         continue
      end ~}}
  ${{ for.index + 1 }}{{ if !for.last }}, {{ end }}
  {{~ end ~}})
RETURNING {{ if table.primary_key.is_some }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }}
`, {{~ for column in table.columns ~}}{{~ if column.auto_increment
		continue
		end ~}}body.C_{{ column.name }}{{ if !for.last }}, {{ end }}{{ end }})
	return row.Scan(&body.C_{{ if table.primary_key.is_some }}{{ table.primary_key.value.column }}{{ else }}{{ table.columns[0].name }}{{ end }})
}

{{ if table.primary_key.is_some }}
func (d DAO) {{ table.name|string.capitalize }}Update(key {{ table.primary_key.go_type }}, body {{ table.name|string.capitalize }}) err {
	_, err := d.db.Exec(`
UPDATE
  "{{ table.name }}"
SET
  {{~ for column in table.columns ~}}
  "{{column.name}}" = ${{ for.index + 2 }}{{ if !for.last }},{{ end }}
  {{~ end ~}}
WHERE
  {{ table.primary_key.name }} = $1
`, id, {{ for column in table.columns }}body.C_{{ column.name }}{{ if !for.last }},{{ end }}{{ end }})
	return err
}

func (d DAO) {{ table.name|string.capitalize }}Delete(key {{ table.primary_key.go_type }}) error {
	_, err := d.db.Exec(`DELETE FROM "{{ table.name }}" WHERE "{{ table.primary_key.value.column }}" = $1`, key)
	return err
}
{{ end }}
