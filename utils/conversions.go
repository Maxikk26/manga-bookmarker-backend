package utils

import (
	"fmt"
	"github.com/dranikpg/dto-mapper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
)

var Mapper dto.Mapper

func AddConvertionFunctions() {
	Mapper.AddConvFunc(objectIDToString)
	Mapper.AddConvFunc(stringToObjectID)
	Mapper.AddConvFunc(primitiveDateTimeToTime)
	Mapper.AddConvFunc(timeToPrimitiveDateTime)
}

func objectIDToString(input primitive.ObjectID) string {
	return input.Hex()
}

func stringToObjectID(input string) primitive.ObjectID {
	// Convert string to primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(input)
	if err != nil {
		fmt.Println("Error converting string to ObjectID:", err)
		objectID = primitive.NewObjectID()
	}
	return objectID
}

func primitiveDateTimeToTime(input primitive.DateTime) time.Time {
	return input.Time()
}

func timeToPrimitiveDateTime(input time.Time) primitive.DateTime {
	return primitive.NewDateTimeFromTime(input)
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

		// Convert the first letter of the field name to lowercase
		if len(fieldName) > 0 {
			fieldName = strings.ToLower(string(fieldName[0])) + fieldName[1:]
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
