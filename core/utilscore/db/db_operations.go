package dbutils

import (
	"github.com/Masterminds/squirrel"
)

// SetterLimitAndOffsetQuery применяет limit и offset к переданному SelectBuilder, если они не nil.
func SetterLimitAndOffsetQuery(query squirrel.SelectBuilder, offset *uint64, limit *uint64) squirrel.SelectBuilder {
	if limit != nil {
		query = query.Limit(*limit)
	}
	if offset != nil {
		query = query.Offset(*offset)
	}
	return query
}
