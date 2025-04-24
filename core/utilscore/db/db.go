package dbutils

import (
	"authentication_service/core/typescore"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/Masterminds/squirrel"
	"github.com/sirupsen/logrus"
)

var (
	fieldsCache sync.Map // Используем sync.Map для кэширования
)

type cacheKey struct {
	typ    reflect.Type
	dbName string
}

// GetOptionsDB обобщённая функция для разбора опций.
// Если Filtering не передан, создаётся новый объект типа T.
// Если передан, то выполняется проверка соответствия типу *T.
func GetOptionsDB[T any](opts ...typescore.ListDbOptions) (typescore.ListDbOptions, *T, error) {
	var options typescore.ListDbOptions
	if len(opts) > 0 {
		options = opts[0]
	}

	var filter *T
	if options.Filtering == nil {
		filter = new(T)
	} else {
		var ok bool
		filter, ok = options.Filtering.(*T)
		if !ok {
			return typescore.ListDbOptions{}, nil, errors.New("expected *T")
		}
	}
	return options, filter, nil
}

func GetStructFieldsDB(processor interface{}, dbName *string) []string {
	if processor == nil || reflect.ValueOf(processor).Kind() != reflect.Ptr || reflect.ValueOf(processor).Elem().Kind() != reflect.Struct {
		logrus.Error("🛑 error GetStructFieldsDB: processor must be a non-nil pointer to a struct")
		return nil
	}

	procType := reflect.TypeOf(processor).Elem()
	name := ""
	if dbName != nil {
		name = *dbName
	}
	key := cacheKey{
		typ:    procType,
		dbName: name,
	}

	if fields, exists := fieldsCache.Load(key); exists {
		return fields.([]string)
	}

	val := reflect.ValueOf(processor).Elem()

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		typeField := procType.Field(i)
		dbTag := typeField.Tag.Get("db")
		if dbTag != "" {
			if typeField.Tag.Get("ignore_db") == "true" || typeField.Tag.Get("ignore_req_db") == "true" {
				continue
			}

			// Добавляем alias, чтобы исключить конфликты при одинаковых именах полей
			if dbName != nil {
				fullField := fmt.Sprintf("%s.%s", *dbName, dbTag)
				fields = append(fields, fullField)
			} else {
				fields = append(fields, dbTag)
			}
		}
	}

	fieldsCache.Store(key, fields)
	return fields
}

func AddNonNullFieldsToQueryUpdate(query squirrel.UpdateBuilder, paramsUpdate interface{}) squirrel.UpdateBuilder {
	v := reflect.ValueOf(paramsUpdate)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if (field.Kind() == reflect.Ptr && !field.IsNil()) || (field.Kind() != reflect.Ptr && !field.IsZero()) {
			// Проверяем тег ignore_update_db для игнорирования полей при обновлении
			ignoreUpdateDbTag := t.Field(i).Tag.Get("ignore_update_db")
			if ignoreUpdateDbTag == "true" {
				continue
			}
			ignoreDbTag := t.Field(i).Tag.Get("ignore_db")
			if ignoreDbTag == "true" {
				continue
			}
			dbTag := t.Field(i).Tag.Get("db")
			if dbTag != "system_id" && dbTag != "role" && dbTag != "created_at" {
				query = query.Set(dbTag, field.Interface())
			}
		}
	}

	return query
}

// Создаем два списка для хранения названий столбцов и соответствующих значений
func GenerateInsertRequest(query squirrel.InsertBuilder, processor interface{}, options ...typescore.InsertOptions) (*string, []interface{}, error) {
	var columns []string
	var values []interface{}
	v := reflect.ValueOf(processor).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		ignoreDbTag := t.Field(i).Tag.Get("ignore_db")
		if ignoreDbTag == "true" {
			continue
		}
		ignoreReqDB := t.Field(i).Tag.Get("ignore_req_db")
		if ignoreReqDB == "true" {
			continue
		}
		dbTag := t.Field(i).Tag.Get("db")
		if dbTag == "" {
			continue
		}
		if dbTag == "serial_id" || dbTag == "system_id" {
			continue
		}

		if (field.Kind() == reflect.Ptr && !field.IsNil()) || (field.Kind() != reflect.Ptr && !field.IsZero()) {
			columns = append(columns, dbTag)
			values = append(values, field.Interface())
		}
	}

	if len(columns) == 0 {
		logrus.Error("🛑 error GenerateInsertRequest: all fields are zero")
		return nil, nil, errors.New("all fields are zero")
	}

	// Применяем опции, если они есть
	if len(options) > 0 {
		for _, opt := range options {
			if opt.Prefix != "" {
				query = query.Prefix(opt.Prefix)
			}
			if opt.Suffix != "" {
				query = query.Suffix(opt.Suffix)
			} else if opt.IgnoreConflict {
				query = query.Suffix("ON CONFLICT DO NOTHING")
			}
		}
	}

	query = query.Columns(columns...).Values(values...)

	sql, args, err := query.ToSql()
	if err != nil {
		logrus.Error("🛑 error GenerateInsertRequest: ", err)
		return nil, nil, err
	}
	if sql == "" {
		logrus.Error("🛑 error GenerateInsertRequest: sql is empty")
		return nil, nil, errors.New("sql is empty")
	}
	return &sql, args, nil
}
