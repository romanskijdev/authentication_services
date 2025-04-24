package utilscore

import "github.com/shopspring/decimal"

// Вспомогательная функция для создания указателя на bool
func PointerToBool(b bool) *bool {
	return &b
}

func PointerToDecimal(i decimal.Decimal) *decimal.Decimal {
	return &i
}

func PointerToUint32(i uint32) *uint32 {
	return &i
}

func PointerToUint64(i uint64) *uint64 {
	return &i
}

func PointerToString(s string) *string {
	return &s
}

func PointerToFloat64(f float64) *float64 {
	return &f
}

func PointerToInt(i int) *int {
	return &i
}
func PointerToInt64(i int64) *int64 {
	return &i
}

func PointerToInt32(i int32) *int32 {
	return &i
}
