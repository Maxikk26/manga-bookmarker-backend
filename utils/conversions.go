package utils

import (
	"github.com/dranikpg/dto-mapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

var Mapper dto.Mapper

func AddConvertionFunctions() {
	Mapper.AddConvFunc(objectIDToString)
	Mapper.AddConvFunc(primitiveDateTimeToTime)
}

func objectIDToString(input primitive.ObjectID) string {
	return input.Hex()
}

func primitiveDateTimeToTime(input primitive.DateTime) time.Time {
	return input.Time()
}

func StructToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(obj)
	t := reflect.TypeOf(obj)

	for i := 0; i < v.NumField(); i++ {
		fieldName := t.Field(i).Name
		fieldValue := v.Field(i).Interface()

		// Skip fields with names "ID", "Id", "_id", etc.
		if strings.ToLower(fieldName) == "id" || strings.ToLower(fieldName) == "_id" {
			continue
		}

		// Check if the field is empty (zero value)
		if !isEmptyValue(v.Field(i)) {
			result[fieldName] = fieldValue
		}
	}
	return result
}

func isEmptyValue(v reflect.Value) bool {
	return v.IsZero()
}
