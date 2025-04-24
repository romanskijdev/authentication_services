package dbutils

import (
	"reflect"
	"sync"

	"github.com/Masterminds/squirrel"
)

var (
	typeCache sync.Map // Используем sync.Map для кэширования
)

type fieldInfo struct {
	dbTag          string
	ignoreDb       bool
	ignoreReqDb    bool
	updatableDb    bool
	fieldInterface func(reflect.Value) interface{}
}

func getCachedFields(t reflect.Type) []fieldInfo {
	if fields, exists := typeCache.Load(t); exists {
		return fields.([]fieldInfo)
	}

	var cachedFields []fieldInfo
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		cInfo := fieldInfo{
			dbTag:       f.Tag.Get("db"),
			ignoreDb:    f.Tag.Get("ignore_db") == "true" || f.Tag.Get("db") == "",
			ignoreReqDb: f.Tag.Get("ignore_req_db") == "true",
			updatableDb: f.Tag.Get("updatable_db") != "false",
			fieldInterface: func(v reflect.Value) interface{} {
				return v.Interface()
			},
		}
		cachedFields = append(cachedFields, cInfo)
	}
	typeCache.Store(t, cachedFields)
	return cachedFields
}

func UpdateSetField(query squirrel.UpdateBuilder, processor interface{}) squirrel.UpdateBuilder {
	v := reflect.ValueOf(processor).Elem()
	t := v.Type()
	fields := getCachedFields(t)
	for i, f := range fields {
		field := v.Field(i)
		if (field.Kind() == reflect.Ptr && !field.IsNil()) || (field.Kind() != reflect.Ptr && !field.IsZero()) {
			if f.ignoreDb || f.ignoreReqDb || !f.updatableDb {
				continue
			}
			query = query.Set(f.dbTag, f.fieldInterface(field))
		}
	}
	return query
}
