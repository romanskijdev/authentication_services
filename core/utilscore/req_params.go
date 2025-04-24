package utilscore

import (
	errm "authentication_service/core/errmodule"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func ParseParamsGetRequest(queryParams url.Values, result interface{}) (*uint64, *uint64, map[string]string, *errm.Error) {
	likeFields := map[string]string{}
	singleValueQueryParams := make(map[string]interface{})
	offset := PointerToUint64(0)
	limit := PointerToUint64(50)
	likeFieldsMode := queryParams.Get("like_fields_mode") == "true"

	paramTypes := make(map[string]string)
	// Получаем типы полей структуры result
	resultType := reflect.TypeOf(result).Elem()
	for i := 0; i < resultType.NumField(); i++ {
		field := resultType.Field(i)
		tag := field.Tag.Get("mapstructure")
		if tag != "" {
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			paramTypes[tag] = fieldType.String()
		}
	}

	for key, values := range queryParams {
		if len(values) > 0 {
			value := values[0]
			switch key {
			case "offset":
				offset = parseUintParam(value)
			case "limit":
				limit = parseUintParam(value)
			default:
				singleValueQueryParams[key] = parseQueryParam(value, likeFieldsMode, &likeFields, key, paramTypes)
			}
		}
	}

	decoder, wEvent := setupDecoder(result)
	if wEvent != nil {
		logrus.Error("🛑 error setting up decoder: ", wEvent.Error)
		return offset, limit, likeFields, wEvent
	}

	err := decoder.Decode(singleValueQueryParams)
	if err != nil {
		logrus.Error("🛑 error decoding query params: ", err)
		return offset, limit, likeFields, errm.NewError("failed_to_decode", err)
	}

	return offset, limit, likeFields, nil
}

func parseUintParam(value string) *uint64 {
	val, err := strconv.ParseUint(value, 10, 64)
	if err == nil {
		return &val
	}
	logrus.Error("🛑 error parsing uint param: ", err)
	return nil
}

func setupDecoder(result interface{}) (*mapstructure.Decoder, *errm.Error) {
	decoderConfig := &mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			decodeHookStringSliceToString,
			decodeHookStringToDecimal,
		),
		Result: result,
	}
	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		logrus.Error("🛑 error creating decoder: ", err)
		return nil, errm.NewError("failed_to_create_decoder", err)
	}
	return decoder, nil
}

func decodeHookStringSliceToString(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if t == reflect.TypeOf(new(string)) && f.Kind() == reflect.Slice {
		slice := reflect.ValueOf(data)
		strSlice := make([]string, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			strSlice[i] = fmt.Sprint(slice.Index(i).Interface())
		}
		result := strings.Join(strSlice, ",")
		return &result, nil
	}
	return data, nil
}

func decodeHookStringToDecimal(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
	if from.Kind() == reflect.String && to == reflect.TypeOf(&decimal.Decimal{}) {
		dec, err := decimal.NewFromString(data.(string))
		if err != nil {
			logrus.Error("🛑 error parsing decimal: ", err)
			return nil, err
		}
		return &dec, nil
	}
	return data, nil
}

func parseQueryParam(value string, likeFieldsMode bool, likeFields *map[string]string, key string, paramTypes map[string]string) interface{} {
	switch value {
	case "true":
		return true
	case "false":
		return false
	case "all":
		return nil
	default:
		if strings.Contains(value, ",") {
			return strings.Split(value, ",")
		}
		if likeFieldsMode {
			(*likeFields)[key] = value
		}
		if paramTypes[key] == "decimal.Decimal" {
			if decVal, err := decimal.NewFromString(value); err == nil {
				return decVal
			} else {
				// Логирование ошибки
				logrus.Error("🛑 error parsing decimal param: ", err)
				return nil
			}
		}

		if paramTypes[key] == "uint64" {
			if uintVal, err := strconv.ParseUint(value, 10, 64); err == nil {
				logrus.Error("🛑 error parsing uint64 param: ", err)
				return uintVal
			} else {
				return nil
			}
		}

		// Проверяем тип данных из paramTypes
		if paramTypes[key] == "string" {
			return value
		}

		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
		return value
	}
}
