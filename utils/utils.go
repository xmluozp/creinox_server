package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gobuffalo/nulls"
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

	validate.RegisterCustomTypeFunc(ValidateValuer, nulls.String{}, nulls.Int{})

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

// ValidateValuer implements validator.CustomTypeFunc
func ValidateValuer(field reflect.Value) interface{} {

	if valuer, ok := field.Interface().(driver.Valuer); ok {

		val, err := valuer.Value()

		if err == nil {
			return val
		}
		// handle the error how you want

		// fmt.Println("自定义validator", field, val)
	}

	return nil
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

func ParseFlight(s string) (letters, numbers string) {

	// trim the postfix letters (isnt needed)
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] >= '0' && s[i] <= '9' {
			s = s[0 : i+1]
			break
		}
	}

	var l, n []rune
	for _, r := range s {
		switch {
		case r >= 'A' && r <= 'Z':
			l = append(l, r)
		case r >= 'a' && r <= 'z':
			l = append(l, r)
		case r >= '0' && r <= '9':
			n = append(n, r)
		}
	}
	return string(l), string(n)
}

func GetField(tag, key string, s interface{}) (reflect.Value, reflect.StructField) {
	rt := reflect.TypeOf(s)
	v := reflect.ValueOf(s)
	if rt.Kind() != reflect.Struct {
		panic("bad type")
	}
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		tagv := strings.Split(f.Tag.Get(key), ",")[0] // use split to ignore tag "options" like omitempty, etc.

		if tagv == tag {
			return v.Field(i), f
		}
	}
	return reflect.Value{}, reflect.StructField{}
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

	returnValue, _ := GetField(tag, key, s)

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
