package utils

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"reflect"
	"regexp"
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
		Log(err, "SendJson出错 ", data.Info, status)
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

// 一些reflect返回的时候，类型是interface{}，但error有可能是nil，interface不能是nil，所以需要转一下
func ParseError(err reflect.Value) error {

	myErr, ok := err.Interface().(error)

	if myErr == nil || !ok {
		return nil
	}

	return myErr
}

func ErrorFromString(message string) error {

	return errors.New(message)
}

func parseInt(s string) int {

	returnValue, err := strconv.Atoi(s)

	if err != nil {
		return 0
	} else {
		return returnValue
	}
}

// 把文本截成文本，数字的格式
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

// 把文本截成文本，数字，文本，数字……的格式
func ParseFlightSlice(s string) (returnSlice []string) {

	if s == "" {
		return nil
	}

	var str []rune

	isLastDigit := s[0] >= '0' && s[0] <= '9'

	for _, r := range s {

		// 当前是字母还是数字
		isThisDigit := r >= '0' && r <= '9'

		// 字母和数字之间发生了切换
		if isLastDigit != isThisDigit {

			// 推入
			returnSlice = append(returnSlice, string(str))
			str = nil
		}

		str = append(str, r)

		// 跟上
		isLastDigit = isThisDigit
		// 如果切换就存入
	}

	// 最后一组str如果有，直接存入
	if len(str) > 0 {
		returnSlice = append(returnSlice, string(str))
	}

	return returnSlice
}

// 阿拉伯数字转人民币大写
// from: https://www.jianshu.com/p/f6367c747798
func FormatConvertNumToCny(num float32) string {

	num64 := float64(num)

	strnum := strconv.FormatFloat(num64*100, 'f', 0, 64)
	sliceUnit := []string{"仟", "佰", "拾", "亿", "仟", "佰", "拾", "万", "仟", "佰", "拾", "元", "角", "分"}
	// log.Println(sliceUnit[:len(sliceUnit)-2])
	s := sliceUnit[len(sliceUnit)-len(strnum) : len(sliceUnit)]
	upperDigitUnit := map[string]string{"0": "零", "1": "壹", "2": "贰", "3": "叁", "4": "肆", "5": "伍", "6": "陆", "7": "柒", "8": "捌", "9": "玖"}
	str := ""
	for k, v := range strnum[:] {
		str = str + upperDigitUnit[string(v)] + s[k]
	}
	reg, err := regexp.Compile(`零角零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零角`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零分$`)
	str = reg.ReplaceAllString(str, "整")

	reg, err = regexp.Compile(`零[仟佰拾]`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零{2,}`)
	str = reg.ReplaceAllString(str, "零")

	reg, err = regexp.Compile(`零亿`)
	str = reg.ReplaceAllString(str, "亿")

	reg, err = regexp.Compile(`零万`)
	str = reg.ReplaceAllString(str, "万")

	reg, err = regexp.Compile(`零*元`)
	str = reg.ReplaceAllString(str, "元")

	reg, err = regexp.Compile(`亿零{0, 3}万`)
	str = reg.ReplaceAllString(str, "^元")

	reg, err = regexp.Compile(`零元`)
	str = reg.ReplaceAllString(str, "零")
	if err != nil {
		log.Fatal(err)
	}
	return str
}

func copyMap(originalMap map[string]string) map[string]string {

	// Create the target map
	targetMap := make(map[string]string)

	// Copy from the original map to the target map
	for key, value := range originalMap {
		targetMap[key] = value
	}

	return targetMap
}

// 把字符串以逗号切开，查找包含关系。用来判断col是不是有这个tag
func CheckCol(str string, target string) bool {

	split := strings.Split(str, ",")

	for _, a := range split {
		if strings.Trim(a, " ") == strings.Trim(target, " ") {
			return true
		}
	}
	return false
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

func GetPrintSourceFromInterface(item interface{}) (map[string]interface{}, error) {

	byteArray, err := json.Marshal(item)
	var m map[string]interface{}
	err = json.Unmarshal(byteArray, &m)

	// 200819: 因为打印模板需要判断第一行。如果清掉了，列名就丢失了
	m = ClearNil(m)

	return m, err
}

// key: 在dataSource里是什么元素； subitemKey：list中被修改的元素的名称
func ModifyDataSourceList(ds map[string]interface{},
	key string,
	subitemKey string,
	callback func(map[string]interface{}) string) map[string]interface{} {

	if list, ok := ds[key].([]interface{}); ok {
		for i := 0; i < len(list); i++ {
			if subitem, ok := list[i].(map[string]interface{}); ok {
				subitem[subitemKey] = callback(subitem)
			}
		}
		ds[key] = list
	}

	return ds
}

// 迭代清除map里面所有空值
func ClearNil(m map[string]interface{}) map[string]interface{} {
	for columnName := range m {
		switch field := m[columnName].(type) {
		case map[string]interface{}:

			m[columnName] = ClearNil(field)
			if len(field) == 0 {
				delete(m, columnName)
			}

		case []interface{}:

			// 清除列表需要迭代每一项. 保留第一项因为打印模板需要取第一项的列名
			for i := 1; i < len(field); i++ {
				if f, ok := field[i].(map[string]interface{}); ok {
					field[i] = ClearNil(f)
				}
			}
			// 清除完覆盖回去
			m[columnName] = field

			if len(field) == 0 {
				delete(m, columnName)
			}

		case string:

			if len(field) == 0 {
				delete(m, columnName)
			}

		case nulls.Nulls:
			delete(m, columnName)

		case nil:
			delete(m, columnName)

		default:
		}
	}

	return m
}

func RandomString(length int) string {

	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()

	return str
}

// 运行命令行来进行pdf转换
func interactiveToexec(commandName string, params []string) (string, bool) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	log.Println(cmd.Args)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	log.Printf("%s\n", w)
	if err != nil {
		log.Println("Error: <", err, "> when exec command read out buffer")
		return "", false
	} else {
		return string(buf), true
	}
}

func FormatDate(t time.Time) string {
	return t.Format("2006年01月02日")
}

func FormatDateTime(t time.Time) string {
	return t.Format("2006年01月02日15:04")
}

func Log(err error, a ...interface{}) {

	if err != nil {
		fmt.Println(time.Now().Format(time.RFC850), err.Error(), a)
	} else {
		fmt.Println(time.Now().Format(time.RFC850), "Works Fine: ", a)
	}

}
