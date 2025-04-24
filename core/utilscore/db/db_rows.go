package dbutils

import (
	"encoding/json"
	"errors"
	"reflect"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
)

// Кэш для функции prepareFieldMap
var prepareFieldMapCache sync.Map // Используем sync.Map для кэширования

type cachedFieldInfo struct {
	dbTag     string
	fieldPath []int
}

// ScanRowsToStructRows сканирует строки из pgx.Rows в структуру и totalCount.
func ScanRowsToStructRows(rows pgx.Rows, dest interface{}, totalCount *uint64) error {
	if dest == nil {
		return errors.New("dest must not be nil")
	}

	val := reflect.ValueOf(dest)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return errors.New("dest must be a non-nil pointer")
	}

	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return errors.New("dest must be a pointer to a struct")
	}

	typ := elem.Type()
	fieldMap := prepareFieldMap(elem, typ)

	columns := getColumnsFromRows(rows)
	args := getScanArgs(columns, fieldMap, totalCount)

	err := rows.Scan(args...)
	if err != nil {
		logrus.Error("🛑 error ScanRowsToStructRows: ", err)
		return err
	}

	return nil
}

func getColumnsFromRows(rows pgx.Rows) []string {
	fieldDescriptions := rows.FieldDescriptions()
	columns := make([]string, len(fieldDescriptions))
	for i, fd := range fieldDescriptions {
		columns[i] = string(fd.Name)
	}
	return columns
}

func getScanArgs(columns []string, fieldMap map[string]interface{}, totalCount *uint64) []interface{} {
	scanArgs := make([]interface{}, len(columns))
	for i, col := range columns {
		if col == "total_count" && totalCount != nil {
			scanArgs[i] = totalCount
		} else if field, ok := fieldMap[col]; ok {
			scanArgs[i] = field
		} else {
			dummy := new(interface{})
			scanArgs[i] = &dummy
		}
	}
	return scanArgs
}

func prepareFieldMap(val reflect.Value, typ reflect.Type) map[string]interface{} {
	if cachedFields, ok := prepareFieldMapCache.Load(typ); ok {
		fieldMap := make(map[string]interface{})
		// Используем закешированные данные
		for _, info := range cachedFields.([]cachedFieldInfo) {
			field := val.FieldByIndex(info.fieldPath)
			if field.CanAddr() {
				fieldMap[info.dbTag] = field.Addr().Interface()
			}
		}
		return fieldMap
	}

	// Если тип не закеширован, собираем информацию о полях
	var cachedFieldsInfo []cachedFieldInfo
	fieldMap := make(map[string]interface{})
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		if fieldType.Tag.Get("ignore_db") == "true" {
			continue
		}
		if tag := fieldType.Tag.Get("db"); tag != "" {
			if field.Type() == reflect.TypeOf(json.RawMessage{}) {
				fieldMap[tag] = field.Addr().Interface()
			} else if field.Kind() == reflect.Slice && field.Type().Elem() == reflect.TypeOf(json.RawMessage{}) {
				fieldMap[tag] = &json.RawMessage{}
			} else if field.Kind() == reflect.Struct && field.Type() == reflect.TypeOf(json.RawMessage{}) {
				fieldMap[tag] = &json.RawMessage{}
			} else {
				fieldMap[tag] = field.Addr().Interface()
			}
			cachedFieldsInfo = append(cachedFieldsInfo, cachedFieldInfo{
				dbTag:     tag,
				fieldPath: fieldType.Index,
			})
		}
	}

	// Сохраняем информацию в кэше
	prepareFieldMapCache.Store(typ, cachedFieldsInfo)

	return fieldMap
}
