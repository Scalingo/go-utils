package logger

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"
)

// FieldsFor extracts loggable fields from a struct based on the "log" tag.
// It returns a logrus.Fields map where the keys are the tag values prefixed
// with the provided prefix, and the values are the corresponding field values.
//
// If the struct has no fields with the "log" tag, it checks if the struct
// implements the fmt.Stringer interface. If it does, it adds a single field
// with the prefix as the key and the result of the String() method as the value.
// If the struct does not implement fmt.Stringer, it adds a single field with
// the prefix as the key and a default error message as the value.
//
// Parameters:
// - value: The struct to extract fields from.
// - prefix: The prefix to add to each field key.
//
// Returns:
// - logrus.Fields: A map of loggable fields.
func FieldsFor(value interface{}, prefix string) logrus.Fields {
	val := reflect.ValueOf(value)

	fields := logrus.Fields{}

	for i := 0; i < val.NumField(); i++ {
		name, found := val.Type().Field(i).Tag.Lookup("log")
		if found {
			fields[fmt.Sprintf("%s_%s", prefix, name)] = val.Field(i).Interface()
		}
	}

	if len(fields) != 0 {
		return fields
	}

	if valueStr, ok := value.(fmt.Stringer); ok {
		fields[prefix] = valueStr.String()
	} else {
		fields[prefix] = "failed to use FieldsFor on struct: Invalid type"
	}

	return fields
}

func WithStructToCtx(ctx context.Context, value interface{}, prefix string) (context.Context, logrus.FieldLogger) {
	return WithFieldsToCtx(ctx, FieldsFor(value, prefix))
}
