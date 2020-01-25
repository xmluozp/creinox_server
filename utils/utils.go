package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/xmluozp/creinox_server/models"
)

var validate *validator.Validate

func SendError(w http.ResponseWriter, status int, data models.JsonRowsReturn) {

	// ??
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func SendJsonError(w http.ResponseWriter, status int, data models.JsonRowsReturn) {

	// ??
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func SendSuccess(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func ValidateInputs(obj interface{}) (bool, models.JsonRowsReturn) {

	var returnValue models.JsonRowsReturn

	// https://github.com/go-playground/validator
	validate = validator.New()

	errValidation := validate.Struct(obj)

	if errValidation != nil {

		returnValue.Message = make(map[string]string)

		// TODO：根据field的名字找到json:tag然后生成map
		for _, errValidation := range errValidation.(validator.ValidationErrors) {

			t := reflect.TypeOf(obj)
			field, isFound := t.FieldByName(errValidation.Field())
			if isFound {
				tag := field.Tag.Get("json")
				errMessage := field.Tag.Get("errm")

				fmt.Println("errm", errMessage)
				returnValue.Message[tag] = errMessage
			}
		}
	}

	return errValidation == nil, returnValue
}

func ParseError(err reflect.Value) error {

	myErr, ok := err.Interface().(error)

	if myErr == nil || !ok {
		return nil
	}

	return myErr
}

func parseInt(s string) int {

	returnValue, err := strconv.Atoi(s)

	if err != nil {
		return 0
	} else {
		return returnValue
	}
}
