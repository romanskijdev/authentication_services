package dbutils

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"
)

var (
	reflectionCacheQueryWhere sync.Map // Используем sync.Map для кэширования
)

func getCachedFieldsQueryWhere(t reflect.Type) []reflect.StructField {
	if fields, exists := reflectionCacheQueryWhere.Load(t); exists {
		return fields.([]reflect.StructField)
	}

	var cachedFields []reflect.StructField
	for i := 0; i < t.NumField(); i++ {
		cachedFields = append(cachedFields, t.Field(i))
	}
	reflectionCacheQueryWhere.Store(t, cachedFields)
	return cachedFields
}
func AddNonNullFieldsToQueryWhere(
	query squirrel.SelectBuilder,
	processor interface{},
	likeFields map[string]string,
	baseName *string,
) squirrel.SelectBuilder {
	val := reflect.ValueOf(processor).Elem()
	typ := val.Type()

	var or squirrel.Or
	fields := getCachedFieldsQueryWhere(typ)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		typeField := fields[i]

		// Пропускаем невалидные или пустые поля
		if !field.IsValid() || field.IsZero() {
			continue
		}

		// Безопасно проверяем на nil, только если тип позволяет
		if isNilable(field.Kind()) && field.IsNil() {
			continue
		}

		ignoreDbTag := typeField.Tag.Get("ignore_db")
		if ignoreDbTag == "true" {
			continue
		}

		ignoreReqDbTag := typeField.Tag.Get("ignore_req_db")
		if ignoreReqDbTag == "true" {
			continue
		}

		dbTag := typeField.Tag.Get("db")
		if dbTag != "" {
			if baseName != nil {
				dbTag = fmt.Sprintf("%s.%s", *baseName, dbTag)
			}

			if likeValue, ok := likeFields[dbTag]; ok {
				// Добавляем LIKE фильтр
				switch dbTag {
				case "serial_id":
					or = append(or, squirrel.Expr("CAST(serial_id AS TEXT) ILIKE ?", likeValue+"%"))
				case "telegram_id":
					or = append(or, squirrel.Expr("CAST(telegram_id AS TEXT) ILIKE ?", likeValue+"%"))
				default:
					or = append(or, squirrel.ILike{dbTag: likeValue + "%"})
				}
			} else {
				query = query.Where(processFieldCached(field, dbTag))
			}
		}
	}

	if len(or) > 0 {
		query = query.Where(or)
	}

	return query
}

// Вспомогательная функция — проверка, может ли тип быть nil
func isNilable(kind reflect.Kind) bool {
	switch kind {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return true
	default:
		return false
	}
}

func processFieldCached(field reflect.Value, dbTag string) squirrel.Sqlizer {
	fieldType := field.Type().String()
	switch fieldType {
	case "*bool":
		boolVal := *field.Interface().(*bool)
		return squirrel.Eq{dbTag: boolVal}
	case "*string":
		str := *field.Interface().(*string)
		if strings.Contains(str, ",") {
			values := strings.Split(str, ",")
			or := squirrel.Or{}
			for _, value := range values {
				or = append(or, squirrel.Eq{dbTag: value})
			}
			return or
		}
		return squirrel.Eq{dbTag: str}
	default:
		return squirrel.Eq{dbTag: field.Interface()}
	}
}
