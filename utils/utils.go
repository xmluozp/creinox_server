package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/xmluozp/creinox_server/models"
)

var validate *validator.Validate

func SendJson(w http.ResponseWriter, status int, data models.JsonRowsReturn, err error) {
	if err != nil {
		Log(err, data.Info)
		w.WriteHeader(status)
	} else {
		w.Header().Set("Content-Type", "application/json")
	}
	json.NewEncoder(w).Encode(data)
}

func SendError(w http.ResponseWriter, status int, data models.JsonRowsReturn) {

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

				tag := strings.Split(field.Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
				// tag := field.Tag.Get("json")
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

func GetField(tag, key string, s interface{}) reflect.Value {
	rt := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tagv := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.

		if tagv == tag {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

func GetFieldName(tag, key string, s interface{}) (fieldname string) {
	rt := reflect.TypeOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == tag {
			return f.Name
		}
	}
	return ""
}

func GetFieldValue(tag, key string, s interface{}) (value interface{}) {
	// rt := reflect.TypeOf(s)
	// v := reflect.ValueOf(s)
	// if rt.Kind() != reflect.Struct {
	// 	panic("bad type")
	// }
	// for i := 0; i < rt.NumField(); i++ {
	// 	f := rt.Field(i)
	// 	tagv := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.
	// 	if tagv == tag {
	// 		return v.Field(i).Interface()
	// 	}
	// }
	returnValue := GetField(tag, key, s)

	// returnValue.Interface()
	// fmt.Println("getfield value", returnValue.Interface())

	if returnValue.IsValid() {
		return returnValue.Interface()
	} else {
		return nil
	}
}

func Log(err error, a ...interface{}) {

	if err != nil {
		fmt.Println(time.Now().Format(time.RFC850), err.Error(), a)
	} else {
		fmt.Println(time.Now().Format(time.RFC850), "Works Fine: ", a)
	}

}
