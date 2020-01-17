package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

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
		newQueryStringBegin = "SELECT * FROM " + tableName
	}

	newQueryStringSearchTerms = " WHERE 1=1"

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
	}

	if pagination.PerPage >= 0 {
		// page从0开始，在前端才转成1
		newQueryString += fmt.Sprintf(" LIMIT %d, %d", pagination.Page*pagination.PerPage, pagination.PerPage)
	}

	// pagination: 基础查询加query加order和分页
	rows, err := db.Query(newQueryStringBegin + newQueryStringSearchTerms + newQueryString)

	db.QueryRow("SELECT COUNT(*) FROM " + tableName + newQueryStringSearchTerms).Scan(&count)

	pagination.TotalCount = count
	pagination.TotalPage = int(math.Ceil(float64(count) / float64(pagination.PerPage)))

	rowCount = min(count-pagination.PerPage*(pagination.Page), pagination.PerPage)

	pagination.RowCount = rowCount

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
	// fmt.Println("get roles: ", page, reflect.TypeOf(q), "searchTerms:", role)

	// 搜索用的object
	var queryObject map[string]string

	json.Unmarshal([]byte(q), &queryObject)
	// fmt.Println("get map: ", queryObject)

	// for k := range queryObject {
	// 	fmt.Println("键值对: ", k, queryObject[k])
	// }

	return queryObject
}

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
