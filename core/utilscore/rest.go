package utilscore

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"net/url"
	"reflect"
	"strconv"
)

func StructToURLParams(params interface{}) (url.Values, error) {
	values := url.Values{}
	if params == nil {
		return values, nil
	}

	paramsMap := make(map[string]interface{})
	err := mapstructure.Decode(params, &paramsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to decode params: %w", err)
	}

	for key, value := range paramsMap {
		if value == nil || fmt.Sprintf("%v", value) == "<nil>" {
			continue
		}
		strValue := formatValue(value)
		if strValue != "" && strValue != "0" && strValue != "false" {
			values.Add(key, strValue)
		}
	}
	return values, nil
}

func formatValue(value interface{}) string {
	// Проверяем, является ли value указателем
	val := reflect.ValueOf(value)
	if val.Kind() == reflect.Ptr && !val.IsNil() {
		// Получаем значение, на которое указывает указатель
		val = val.Elem()
	}

	// Теперь val содержит фактическое значение, даже если изначально был передан указатель
	// Продолжаем форматирование val как обычно
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	// Добавьте другие типы по необходимости
	default:
		return fmt.Sprintf("%v", val)
	}
}
