package marshallerutils

import (
	"net"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type Serialization struct{}

func InitSerializationUtils() *Serialization {
	return &Serialization{}
}

func (m *Serialization) StringToWrapperStringValue(s *string) *wrapperspb.StringValue {
	if s == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: *s}
}

func (m *Serialization) IPToWrapperStringValue(ip *net.IP) *wrapperspb.StringValue {
	if ip == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: ip.String()}
}

func (m *Serialization) BoolToWrapperBoolValue(b *bool) *wrapperspb.BoolValue {
	if b == nil {
		return nil
	}
	return &wrapperspb.BoolValue{Value: *b}
}

func (m *Serialization) DecimalToWrapperStringValue(d *decimal.Decimal) *wrapperspb.StringValue {
	if d == nil {
		return nil
	}
	return &wrapperspb.StringValue{Value: d.String()}
}

func (m *Serialization) Float64ToWrapperDoubleValue(f *float64) *wrapperspb.DoubleValue {
	if f == nil {
		return nil
	}
	return &wrapperspb.DoubleValue{Value: *f}
}

func (m *Serialization) Uint64ToWrapperUInt64Value(u *uint64) *wrapperspb.UInt64Value {
	if u == nil {
		return nil
	}
	return &wrapperspb.UInt64Value{Value: *u}
}

func (m *Serialization) Int64ToWrapperInt64Value(i *int64) *wrapperspb.Int64Value {
	if i == nil {
		return nil
	}
	return &wrapperspb.Int64Value{Value: *i}
}

func (m *Serialization) TimePtrToTimestampPB(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func (m *Serialization) StringTimeToWrapperOnlyDate(pbTime *string) *wrapperspb.StringValue {
	if pbTime != nil {
		// Парсим строку в time.Time
		parsedTime, err := time.Parse(time.RFC3339, *pbTime)
		if err != nil {
			logrus.Errorln("🔴 error StringTimeToWrapperOnlyDate: parsedTime: ", pbTime)
			// Обработка ошибки парсинга, можно вернуть nil или пустую строку
			return nil
		}
		// Форматируем дату в строку формата "yyyy-mm-dd"
		formattedDate := parsedTime.Format("2006-01-02")
		// Создаем и возвращаем wrapperspb.StringValue
		return wrapperspb.String(formattedDate)
	}
	return nil
}
