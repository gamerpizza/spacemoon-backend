package postgre_sql

import "errors"

var errComplexKey = errors.New("complex key detected in qerrors")

func CreateQueryFromkey(key map[string]any) (string, []any, error) {
	var query = ""
	vals := make([]any, 0)

	for i, k := range key {
		query = query + i + " = ? AND "
		vals = append(vals, k)
	}

	query = query[:len(query)-5]

	return query, vals, nil
}
