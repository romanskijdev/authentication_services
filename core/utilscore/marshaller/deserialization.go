package marshallerutils

import (
	"math/big"
	"net"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Deserialization struct{}

func InitDeserializationUtils() *Deserialization {
	return &Deserialization{}
}

func (m *Deserialization) OptionalIP(pbStr *wrapperspb.StringValue) *net.IP {
	if pbStr != nil {
		ip := net.ParseIP(pbStr.Value)
		if ip != nil {
			return &ip
		}
	}
	return nil
}

func (m *Deserialization) OptionalBool(pbBool *wrapperspb.BoolValue) *bool {
	if pbBool != nil {
		return &pbBool.Value
	}
	return nil
}

func (m *Deserialization) OptionalInt64(pbInt64 *wrapperspb.Int64Value) *int64 {
	if pbInt64 != nil {
		return &pbInt64.Value
	}
	return nil
}

// Вспомогательная функция для преобразования google.protobuf.UInt64Value в *big.Int
func (m *Deserialization) OptionalBigInt(pbUInt64 *wrapperspb.UInt64Value) *big.Int {
	if pbUInt64 != nil {
		// Создаем новый big.Int и устанавливаем его значение
		return big.NewInt(0).SetUint64(pbUInt64.Value)
	}
	return nil
}

// Вспомогательная функция для *string поля
func (m *Deserialization) OptionalString(pbStr *wrapperspb.StringValue) *string {
	if pbStr != nil {
		return &pbStr.Value
	}
	return nil
}

// Вспомогательная функция для *float64 поля
func (m *Deserialization) OptionalFloat64(pbFloat64 *wrapperspb.DoubleValue) *float64 {
	if pbFloat64 != nil {
		return &pbFloat64.Value
	}
	return nil
}

// Вспомогательная функция для *uint64 поля
func (m *Deserialization) OptionalUint64(pbUint64 *wrapperspb.UInt64Value) *uint64 {
	if pbUint64 != nil {
		return &pbUint64.Value
	}
	return nil
}

// Вспомогательная функция для *uint32 поля
func (m *Deserialization) OptionalUint32(pbUint8 *wrapperspb.UInt32Value) *uint32 {
	if pbUint8 != nil {
		value := uint32(pbUint8.Value)
		return &value
	}
	return nil
}

// Вспомогательная функция для *time.Time поля
func (m *Deserialization) OptionalTime(pbTime *timestamppb.Timestamp) *time.Time {
	if pbTime != nil {
		val := pbTime.AsTime()
		return &val
	}
	return nil
}

// Вспомогательная функция для *wrapperspb.StringValue поля
func (m *Deserialization) OptionalStringTimeOnlyDate(pbTime *wrapperspb.StringValue) *string {
	if pbTime != nil {
		// Парсим строку в time.Time
		parsedTime, err := time.Parse("2006-01-02", pbTime.GetValue())
		if err != nil {
			logrus.Errorln("🔴 error OptionalStringTimeOnlyDate:parsedTime: ", err)
			// Обработка ошибки парсинга, можно вернуть nil или пустую строку
			return nil
		}
		// Форматируем дату в строку формата "yyyy-mm-dd"
		val := parsedTime.Format("2006-01-02")
		return &val
	}
	return nil
}

// Вспомогательная функция для *decimal.Decimal поля
func (m *Deserialization) OptionalDecimal(pbDecimal *wrapperspb.StringValue) *decimal.Decimal {
	if pbDecimal != nil {
		dec, err := decimal.NewFromString(pbDecimal.Value)
		if err != nil {
			return nil
		}
		return &dec
	}
	return nil
}
