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
func DbQueryRows_Customized(db *sql.DB,
	query string,
	tableName string,
	pagination *models.Pagination,
	searchTerms map[string]string,
	dataModel interface{},
	stringJoin string,
	stringAfterWhere string) (
	*sql.Rows,
	error) {

	var count int
	var rowCount int
	var newQueryString string
	var newQueryStringBegin string
	var newQueryStringSearchTerms string

	selectTable := tableName + " mainTable"
	selectColumns := "mainTable.*"

	// 默认排序列。暂时做成只能有一个排序列，多出来的直接覆盖。
	defaultOrderByCol := ""
	// -----------------------
	// 检查字段： 循环类型里的field，根据field的类型来决定字段要不要保留，如何处理join

	rt := reflect.TypeOf(dataModel)
	for i := 0; i < rt.NumField(); i++ {

		// 判断field是不是ref
		f := rt.Field(i)
		_, isRef := f.Tag.Lookup("ref")

		if isRef {
			refTableName := strings.Split(f.Tag.Get("ref"), ",")[0]
			refColumn := strings.Split(f.Tag.Get("ref"), ",")[1]

			refTableNameUnique := refTableName + strconv.Itoa(i)

			selectColumns += ", " + refTableNameUnique + ".*"
			selectTable += fmt.Sprintf(" LEFT JOIN %s %s ON %s.%s = %s.%s", refTableName, refTableNameUnique, "mainTable", refColumn, refTableNameUnique, "id")
		}

		// 判断是不是用来排序的字段，如果有就作为默认排序 (如果前端没有排序，就用它排序)
		col, isCol := f.Tag.Lookup("col")
		if isCol && CheckCol(col, "orderByDesc") {
			defaultOrderByCol = strings.Split(f.Tag.Get("json"), ",")[0] + " DESC"
		} else if isCol && CheckCol(col, "orderByAsc") {
			defaultOrderByCol = strings.Split(f.Tag.Get("json"), ",")[0] + " ASC"
		}

		// 判断是不是 ext，关联表查询用
	}

	selectTable += stringJoin

	if query != "" {
		newQueryStringBegin = query
	} else {
		newQueryStringBegin = "SELECT " + selectColumns + " FROM " + selectTable + " WHERE 1=1"
	}

	newQueryStringSearchTerms = ""

	// searchTerms
	fmt.Println("searchTerms of dbUtils", searchTerms)

	// 用param传参数用
	var keywords []interface{}

	// 判断查询条件：先看model里面有没有同样名称的字段，如果没有就忽略
	for k, v := range searchTerms {

		// TODO: string，int，date三种类型. 根据v的判断
		// 日期是: 2020/01/14,2020/01/23 这种格式. 用,切开两个string都是日期就是
		// https://programming.guide/go/format-parse-string-time-date-example.html
		// 用dataModel辅助判断. 因为int和float32需要可以查范围
		if v == "" {
			continue
		}

		_, field := GetField(k, "json", dataModel)

		if field.Type == nil {
			continue
		}

		// 如果字段类型是文字，就直接搜索like
		if field.Type.String() != "nulls.String" {

			// 看bool: 前端用1和0代替true和false就可以了

			// 看int:
			_, errInt := strconv.Atoi(v)
			k := "mainTable." + k

			if errInt == nil {
				newQueryStringSearchTerms += " AND " + k + " = " + v
				continue
			}

			// 如果不是int, 试着切分。能切分的是range，搜索范围
			ranges := strings.Split(v, ",")

			if len(ranges) == 2 {

				// 试图转为日期
				t1GoFormat, errT1 := time.Parse("2006/01/02", ranges[0])
				t2GoFormat, errT2 := time.Parse("2006/01/02", ranges[1])

				t1 := t1GoFormat.Format("2006-01-02 15:04:05")
				t2 := t2GoFormat.Format("2006-01-02 15:04:05")

				if errT1 == nil && errT2 == nil {
					newQueryStringSearchTerms += " AND DATE(" + k + ") >= DATE('" + t1 + "') AND " + k + " <= DATE('" + t2 + "')"
					continue
				} else {

					// 如果不是日期，则是数字范围
					newQueryStringSearchTerms += " AND " + k + " >= " + ranges[0] + " AND " + k + " <= " + ranges[1]
					continue
				}
			}
		} else {

			// 看是不是内定的keywords模糊搜索, 如果是的话就生成一串or，贴在keywords后面
			blurkeywordsStr, isBlurKeyword := field.Tag.Lookup("keywords")

			if isBlurKeyword {
				newQueryStringSearchTerms += " AND (1 = 0"
				blurkeywords := strings.Split(blurkeywordsStr, "|")

				for _, colName := range blurkeywords {
					newQueryStringSearchTerms += " OR mainTable." + colName + " LIKE ?"
					keywords = append(keywords, "%"+v+"%")
				}

				newQueryStringSearchTerms += ")"

			} else {
				// 最后作为string去like
				newQueryStringSearchTerms += " AND mainTable." + k + " LIKE ?"
				keywords = append(keywords, "%"+v+"%")
			}
		}
	}

	if pagination.OrderBy != "" {
		newQueryString += " ORDER BY mainTable." + pagination.OrderBy + " " + pagination.Order
	} else if defaultOrderByCol != "" {
		newQueryString += " ORDER BY mainTable." + defaultOrderByCol + ", mainTable.id DESC"
	} else {
		newQueryString += " ORDER BY mainTable.id DESC"
	}

	if pagination.PerPage > 0 {
		// page从0开始，在前端才转成1
		newQueryString += fmt.Sprintf(" LIMIT %d, %d", pagination.Page*pagination.PerPage, pagination.PerPage)
	}

	newQueryStringSearchTerms += stringAfterWhere

	fmt.Println("keywords", keywords)
	//  ===================================== pagination: 基础查询加query加order和分页
	rows, err := db.Query(newQueryStringBegin+newQueryStringSearchTerms+newQueryString, keywords...)

	scanErr := db.QueryRow("SELECT COUNT(*) FROM "+tableName+" mainTable WHERE 1=1 "+newQueryStringSearchTerms, keywords...).Scan(&count)

	if scanErr != nil {
		fmt.Println("scan出错", scanErr.Error())
	}

	pagination.TotalCount = count
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(pagination.PerPage)))

	rowCount = min(count-pagination.PerPage*(pagination.Page), pagination.PerPage)

	pagination.RowCount = rowCount

	fmt.Println("------------------")
	fmt.Println("查询运行的sql语句：", newQueryStringBegin+newQueryStringSearchTerms+newQueryString)
	fmt.Println("------------------")

	if err != nil {
		err = errors.New(err.Error() + ". Sql: " + newQueryStringBegin + newQueryStringSearchTerms + newQueryString)
	}

	return rows, err
}

func DbQueryRows(db *sql.DB,
	query string,
	tableName string,
	pagination *models.Pagination,
	searchTerms map[string]string,
	dataModel interface{}) (
	*sql.Rows,
	error) {

	return DbQueryRows_Customized(db, query, tableName, pagination, searchTerms, dataModel, "", "")

}

func DbQueryRow(db *sql.DB,
	query string,
	tableName string,
	id int,
	dataModel interface{}) *sql.Row {

	var newQueryString string
	var newQueryStringBegin string
	var newQueryStringSearchTerms string

	selectTable := tableName + " mainTable"
	selectColumns := "mainTable.*"

	rt := reflect.TypeOf(dataModel)
	for i := 0; i < rt.NumField(); i++ {

		// 判断field是不是ref
		f := rt.Field(i)
		_, isRef := rt.Field(i).Tag.Lookup("ref")

		if isRef {
			refTableName := strings.Split(f.Tag.Get("ref"), ",")[0]
			refColumn := strings.Split(f.Tag.Get("ref"), ",")[1]

			refTableNameUnique := refTableName + strconv.Itoa(i)

			selectColumns += ", " + refTableNameUnique + ".*"
			selectTable += fmt.Sprintf(" LEFT JOIN %s %s ON %s.%s = %s.%s", refTableName, refTableNameUnique, "mainTable", refColumn, refTableNameUnique, "id")
		}
	}

	var row *sql.Row

	if query != "" {
		newQueryStringBegin = query
	} else {
		newQueryStringBegin = "SELECT " + selectColumns + " FROM " + selectTable + " WHERE mainTable.id = ?"
	}

	fmt.Println("搜索单条数据", newQueryStringBegin)

	if id > 0 {
		row = db.QueryRow(newQueryStringBegin+newQueryStringSearchTerms+newQueryString, id)
	} else {
		row = db.QueryRow(newQueryStringBegin + newQueryStringSearchTerms + newQueryString)
	}

	return row
}

func GetPagination(r *http.Request) models.Pagination {
	params := r.URL.Query()

	// ============ pagination (
	var pagination models.Pagination

	// 防止无pagination导致数据全取
	perPagestr := params.Get("perPage")
	if perPagestr == "" {
		perPagestr = "50"
	}

	pagination.Page = parseInt(params.Get("page"))
	pagination.RowCount = parseInt(params.Get("rowCount"))
	pagination.PerPage = parseInt(perPagestr)
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

	// 循环寻找model里面的列名
	for i := 0; i < v.NumField(); i++ {

		col, isCol := t.Field(i).Tag.Lookup("col")

		if !isCol {
			continue
		}

		// 忽略不出现在json里的（通过nulls的valid来判断: 也就是说本系统不允许上传null）
		isValid := v.Field(i).FieldByName("Valid").IsValid() &&
			v.Field(i).FieldByName("Valid").Interface().(bool)

		if !isValid {
			continue
		}

		if isCol && !CheckCol(col, "default") && t.Field(i).Name != "ID" {
			values = append(values, v.Field(i))
			tagName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
			tagName = fmt.Sprintf("`%s`", tagName)
			columns = append(columns, tagName) // 等同于数据库里的column name
			questionMarks = append(questionMarks, "?")
		}
	}

	// 生成INSERT字符串
	values[0] = reflect.ValueOf("INSERT INTO " + tableName + " (" + strings.Join(columns, ",") + ") VALUES(" + strings.Join(questionMarks, ",") + ");")

	execDb := reflect.ValueOf(db).MethodByName("Exec")

	fmt.Println("SQL:insert", values[0])
	// 传入参数：字符串，字段valueof。。。
	out := execDb.Call(values)

	result, _ := out[0].Interface().(sql.Result)
	err := ParseError(out[1])

	return result, err
}

// col属性是 newtime 的，update的时候取系统时间
func DbQueryUpdate(db *sql.DB, tableName string, queryTable string, item interface{}) (sql.Result, *sql.Row, error) {
	// fmt.Println("连接数", db.Stats())
	// 获取item的值和类型
	v := reflect.ValueOf(item)
	t := reflect.TypeOf(item)

	// 动态数组：数值，字段名，问号
	var values []reflect.Value
	var columns []string

	// 第一个准备用来放string的
	values = append(values, reflect.ValueOf(""))

	// 匹配每个字段：前台的json对后台的model
	for i := 1; i < v.NumField(); i++ {

		// 1. 确认model里是不是col, 如果不是col，说明数据库里没有，不允许上传
		// 本系统有很多字段只为了显示用，是在后端生成的，数据库里没有
		col, isCol := t.Field(i).Tag.Lookup("col")

		if !isCol {
			continue
		}

		// 2. 确认前台的json在model有没有匹配的
		// 通过nulls的valid来判断, 如果是个null就跳过，而不是插入null. : 也就是说本系统没办法上传null
		isValid := v.Field(i).FieldByName("Valid").IsValid() &&
			v.Field(i).FieldByName("Valid").Interface().(bool)

		if !isValid {
			continue
		}

		if t.Field(i).Name != "ID" {

			tagName := strings.Split(t.Field(i).Tag.Get("json"), ",")[0] // use split to ignore tag "options" like omitempty, etc.
			tagName = fmt.Sprintf("`%s`", tagName)
			// 假如fk字段是 -1，就设置成null（为了补救上面那个不分青红皂白删掉null的）

			fmt.Println("dbUtils_update", t.Field(i).Name, v.Field(i), v.Field(i).FieldByName("Valid").IsValid())

			if CheckCol(col, "newtime") { // 如果每次提交都无论如何要更新时间

				columns = append(columns, tagName+"=CURRENT_TIMESTAMP")

			} else if CheckCol(col, "fk") { // 如果是个外键

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

	fmt.Println("Update 的sql语句", values[0])

	execDb := reflect.ValueOf(db).MethodByName("Exec")

	// 传入参数：字符串，字段valueof。。。
	out := execDb.Call(values)

	result, _ := out[0].Interface().(sql.Result)
	err := ParseError(out[1])

	// 取出id搜索记录
	id := v.FieldByName("ID").Interface().(nulls.Int).Int

	rowUploaded := DbQueryRow(db, "", queryTable, id, item)

	return result, rowUploaded, err
}

// called from repository. deleteName 删除的对象； queryTable是返回的对象
func DbQueryDelete(db *sql.DB, deleteName string, queryTable string, id int, dataModel interface{}) (sql.Result, *sql.Row, error) {

	rowDeleted := DbQueryRow(db, "", queryTable, id, dataModel)

	// db.QueryRow("SELECT * FROM "+tableName+" WHERE id = ?", id)

	result, err := db.Exec("DELETE FROM "+deleteName+" WHERE id = ?", id)

	fmt.Println("删除 ", "DELETE FROM "+deleteName+" WHERE id = ?", id)

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

// 把字符串以逗号切开，查找包含关系。用来判断col是不是有这个tag
func checkCol(str string, target string) bool {

	split := strings.Split(str, ",")

	for _, a := range split {
		if strings.Trim(a, " ") == strings.Trim(target, " ") {
			return true
		}
	}
	return false
}

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
