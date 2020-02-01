package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/models"
)

// ===================================== 生成搜索后取出rows
func DbQueryRows(db *sql.DB,
	query string,
	tableName string,
	pagination *models.Pagination,
	searchTerms map[string]string,
	dataModel interface{}) (
	*sql.Rows,
	error) {

	var count int
	var rowCount int
	var newQueryString string
	var newQueryStringBegin string
	var newQueryStringSearchTerms string

	if query != "" {
		newQueryStringBegin = query
	} else {
		newQueryStringBegin = "SELECT * FROM " + tableName + " WHERE 1=1"
	}

	newQueryStringSearchTerms = ""

	// searchTerms
	for k, v := range searchTerms {

		// TODO: string，int，date三种类型. 根据v的判断
		// 日期是: 2020/01/14,2020/01/23 这种格式. 用,切开两个string都是日期就是
		// https://programming.guide/go/format-parse-string-time-date-example.html
		// 用dataModel辅助判断. 因为int和float32需要可以查范围

		// 先看int
		_, errInt := strconv.Atoi(v)
		if errInt == nil {
			newQueryStringSearchTerms += " AND " + k + " = " + v
			continue
		}

		// 如果不是int, 试着切分。能切分的是range
		ranges := strings.Split(v, ",")

		// 再看日期
		if len(ranges) == 2 {
			_, errT1 := time.Parse("2006/01/02", ranges[0])
			_, errT2 := time.Parse("2006/01/02", ranges[1])
			if errT1 == nil && errT2 == nil {
				newQueryStringSearchTerms += " AND " + k + " >= DATE(" + ranges[0] + ") AND " + k + " <= DATE(" + ranges[1] + ")"
				continue
			}
		}

		// 最后作为string去like
		newQueryStringSearchTerms += " AND " + k + " LIKE '%" + v + "%'"

	}

	if pagination.OrderBy != "" {
		newQueryString += " ORDER BY " + pagination.OrderBy + " " + pagination.Order
	} else {
		newQueryString += " ORDER BY 1 DESC"
	}

	if pagination.PerPage > 0 {
		// page从0开始，在前端才转成1
		newQueryString += fmt.Sprintf(" LIMIT %d, %d", pagination.Page*pagination.PerPage, pagination.PerPage)
	}

	// pagination: 基础查询加query加order和分页
	rows, err := db.Query(newQueryStringBegin + newQueryStringSearchTerms + newQueryString)

	scanErr := db.QueryRow("SELECT COUNT(*) FROM " + tableName + " WHERE 1=1 " + newQueryStringSearchTerms).Scan(&count)

	if scanErr != nil {
		fmt.Println("scan出错", scanErr.Error())
	}

	pagination.TotalCount = count
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(pagination.PerPage)))

	rowCount = min(count-pagination.PerPage*(pagination.Page), pagination.PerPage)

	pagination.RowCount = rowCount

	fmt.Println("测试看sql语句：", newQueryStringBegin+newQueryStringSearchTerms+newQueryString)

	if err != nil {
		err = errors.New(err.Error() + ". Sql: " + newQueryStringBegin + newQueryStringSearchTerms + newQueryString)
	}

	return rows, err
}

func GetPagination(r *http.Request) models.Pagination {
	params := r.URL.Query()

	// ============ pagination (
	var pagination models.Pagination

	pagination.Page = parseInt(params.Get("page"))
	pagination.RowCount = parseInt(params.Get("rowCount"))
	pagination.PerPage = parseInt(params.Get("perPage"))
	pagination.TotalCount = parseInt(params.Get("totalCount"))
	pagination.TotalPage = parseInt(params.Get("totalPage"))
	pagination.Order = params.Get("order")
	pagination.OrderBy = params.Get("orderBy")
	// ============ pagination )

	return pagination
}

func GetSearchTerms(r *http.Request) map[string]string {

	params := r.URL.Query()
	q := params.Get("q")
	// json.Unmarshal([]byte(q), &role) // 不能这么做，因为有的搜索是范围搜索。不是严格键值对. 有的搜索需要从外表，这是没有规律的

	// 搜索用的object
	var queryObject map[string]string

	json.Unmarshal([]byte(q), &queryObject)

	// for k := range queryObject {
	// 	fmt.Println("键值对: ", k, queryObject[k])
	// }

	return queryObject
}

func DbQueryInsert(db *sql.DB, tableName string, item interface{}) (sql.Result, error) {

	// 获取item的值和类型
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	// 动态数组：数值，字段名，问号
	var values []reflect.Value
	var columns []string
	var questionMarks []string

	// 第一个准备用来放string的
	values = append(values, reflect.ValueOf(""))

	for i := 1; i < v.NumField(); i++ {

		col, isCol := t.Field(i).Tag.Lookup("col")
		isValid := true

		// 忽略不出现在json里的（通过nulls的valid来判断: 也就是说本系统不允许上传null）
		if v.Field(i).FieldByName("Valid").IsValid() { // if its a "nulls"
			isValid = v.Field(i).FieldByName("Valid").Interface().(bool)
		}

		if isValid && isCol && col != "default" && t.Field(i).Name != "ID" {
			values = append(values, v.Field(i))
			tagName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
			columns = append(columns, tagName)                           // 等同于数据库里的column name
			questionMarks = append(questionMarks, "?")
		}
	}

	// 生成INSERT字符串
	values[0] = reflect.ValueOf("INSERT INTO " + tableName + " (" + strings.Join(columns, ",") + ") VALUES(" + strings.Join(questionMarks, ",") + ");")
	execDb := reflect.ValueOf(db).MethodByName("Exec")

	// 传入参数：字符串，字段valueof。。。
	out := execDb.Call(values)

	result, _ := out[0].Interface().(sql.Result)
	err := ParseError(out[1])

	return result, err
}

// col属性是 newtime 的，update的时候取系统时间
func DbQueryUpdate(db *sql.DB, tableName string, item interface{}) (sql.Result, error) {

	// 获取item的值和类型
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	// 动态数组：数值，字段名，问号
	var values []reflect.Value
	var columns []string

	// 第一个准备用来放string的
	values = append(values, reflect.ValueOf(""))

	for i := 1; i < v.NumField(); i++ {
		col, isCol := t.Field(i).Tag.Lookup("col")
		isValid := true

		// 忽略不出现在json里的（通过nulls的valid来判断: 也就是说本系统不允许上传null）
		if v.Field(i).FieldByName("Valid").IsValid() { // if its a "nulls"
			isValid = v.Field(i).FieldByName("Valid").Interface().(bool)
		}

		if isValid && isCol && t.Field(i).Name != "ID" {

			tagName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.

			// 假如fk字段是 -1，就设置成null（为了补救上面那个不分青红皂白删掉null的）

			fmt.Println("dbUtils_update", t.Field(i).Name, v.Field(i), v.Field(i).FieldByName("Valid").IsValid())

			if col == "newtime" { // 如果每次提交都无论如何要更新时间

				columns = append(columns, tagName+"=CURRENT_TIMESTAMP")

			} else if col == "fk" { // 如果是个外键
				foreignKey := v.Field(i).Interface().(nulls.Int)

				// 如果外键不是0，就设置成空。
				if foreignKey.Int <= 0 {
					newFK := nulls.Int{}
					newFK.Int = 0
					newFK.Valid = false
					values = append(values, reflect.ValueOf(newFK))

				} else {
					values = append(values, v.Field(i))
				}
				columns = append(columns, tagName+"=?")
			} else {
				values = append(values, v.Field(i))
				columns = append(columns, tagName+"=?") // 等同于数据库里的column name
			}
		}
	}

	// 对应了 where id = ?
	values = append(values, v.FieldByName("ID"))

	// result, err := db.Exec("UPDATE role SET name=?, rank = ?, auth=? WHERE id=?", &item.Name, &item.Rank, &item.Auth, &item.ID)

	// 生成UPDATE字符串
	values[0] = reflect.ValueOf("UPDATE " + tableName + " SET " + strings.Join(columns, ",") + " WHERE id=?;")

	execDb := reflect.ValueOf(db).MethodByName("Exec")

	// 传入参数：字符串，字段valueof。。。
	out := execDb.Call(values)

	result, _ := out[0].Interface().(sql.Result)
	err := ParseError(out[1])

	return result, err
}

// called from repository
func DbQueryDelete(db *sql.DB, tableName string, id int) (sql.Result, *sql.Row, error) {

	rowDeleted := db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	result, err := db.Exec("DELETE FROM "+tableName+" WHERE id = ?", id)

	if err != nil {

		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1451 {
				err = errors.New("外键约束：有其他数据引用此数据。无法删除。")
			}
		} else {
			fmt.Println("删除出错", err.Error())
		}

	}

	return result, rowDeleted, err
}

// called from repository
// func DbQueryDelete_multiple(db *sql.DB, tableName string, id int) (sql.Result, error) {

// 	result, err := db.Exec("DELETE FROM "+tableName+" WHERE id = ?", id)

// 	if err != nil {

// 		if driverErr, ok := err.(*mysql.MySQLError); ok {
// 			if driverErr.Number == 1451 {
// 				err = errors.New("外键约束：有其他数据引用此数据。无法删除。")
// 			}
// 		}

// 	}

// 	return result, err
// }

//------------ private
func checkCount(row *sql.Row) (count int) {
	row.Scan(&count)
	return count
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
