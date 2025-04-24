package dbutils

import (
	"github.com/Masterminds/squirrel"
)

// BuildSelectQuery создает базовый SELECT-запрос.
func BuildSelectQuery(table string, fields []string) squirrel.SelectBuilder {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(fields...).From(table)
}

// ApplyFilters добавляет фильтры к запросу
func ApplyFilters(query squirrel.SelectBuilder, paramsFiltering interface{}, likeFields map[string]string, baseName ...string) squirrel.SelectBuilder {
	var baseNamePtr *string
	if len(baseName) > 0 {
		baseNamePtr = &baseName[0]
	} else {
		baseNamePtr = nil
	}
	return AddNonNullFieldsToQueryWhere(query, paramsFiltering, likeFields, baseNamePtr)
}
