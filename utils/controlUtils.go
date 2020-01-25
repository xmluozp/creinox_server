package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/models"
)

func GetFunc_RowsWithHTTPReturn(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{}) func() {

	// 完整query:  page=1&rowCount=5&perPage=15&totalCount=10&totalPage=2&order=desc&orderBy=id&q=%7B%22fullName%22%3A%22%E7%8E%8B%E6%80%9D%E8%81%AA%22%7D

	// 权限判断
	// auth := r.Header.Get("Authorization")

	// 声明变量
	var returnValue models.JsonRowsReturn
	items := reflect.Zero(reflect.SliceOf(modelType)).Interface()
	item := reflect.New(modelType).Elem().Interface()

	pagination := GetPagination(r)
	searchTerms := GetSearchTerms(r)

	gerRows := reflect.ValueOf(repo).MethodByName("GetRows")
	args := []reflect.Value{
		reflect.ValueOf(db),
		reflect.ValueOf(item),
		reflect.ValueOf(items),
		reflect.ValueOf(pagination),
		reflect.ValueOf(searchTerms)}

	// 运行数据库语句: db, model, array of model, pagination, query
	out := gerRows.Call(args)
	rows := out[0].Interface()
	paginationOut := out[1].Interface()
	err := ParseError(out[2])

	if err != nil {

		fmt.Println("ge rows error: ", err)
		returnValue.Info = "Server error" + err.Error()
		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
	}

	// 准备返回值
	returnValue.Pagination = paginationOut
	returnValue.SearchTerms = searchTerms
	returnValue.Rows = rows

	w.Header().Set("Content-Type", "application/json")

	return func() { SendSuccess(w, returnValue) }
}

func GetFunc_RowWithHTTPReturn(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{}) func() {

	// 声明变量
	var returnValue models.JsonRowsReturn

	// 接参数
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	getRow := reflect.ValueOf(repo).MethodByName("GetRow")
	args := []reflect.Value{
		reflect.ValueOf(db),
		reflect.ValueOf(id)}
	out := getRow.Call(args)

	row := out[0].Interface()
	err := ParseError(out[1])

	if err != nil {

		fmt.Println("取单数据出错", err.Error())
		if err == sql.ErrNoRows {
			returnValue.Info = "Not Found" + err.Error()
			return func() { SendError(w, http.StatusNotFound, returnValue) }

		} else {
			returnValue.Info = "Server error" + err.Error()
			return func() { SendError(w, http.StatusInternalServerError, returnValue) }
		}
	}

	returnValue.Row = row

	w.Header().Set("Content-Type", "application/json")
	return func() { SendSuccess(w, returnValue) }
}

func GetFunc_AddWithHTTPReturn(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) func() {

	// 声明变量
	var returnValue models.JsonRowsReturn

	itemPtr := reflect.New(modelType).Interface()

	// 不用指针取了再转的话，item会被强行变成map类型
	err := json.NewDecoder(r.Body).Decode(itemPtr)
	item := reflect.ValueOf(itemPtr).Elem().Interface()

	if err != nil {
		fmt.Println("Insert error on controller: ", err)
		returnValue.Info = "Server error" + err.Error()
		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
	}

	// ------------------------------validation : 过后移出去
	isPassedValidation, returnValue := ValidateInputs(item)

	if !isPassedValidation {
		return func() { SendJsonError(w, http.StatusBadRequest, returnValue) } // 400
	}

	addRow := reflect.ValueOf(repo).MethodByName("AddRow")
	args := []reflect.Value{
		reflect.ValueOf(db),
		reflect.ValueOf(item),
		reflect.ValueOf(userId)}
	out := addRow.Call(args)

	row := out[0].Interface()
	errAdd := ParseError(out[1])

	if errAdd != nil {
		returnValue.Info = "Server error" + errAdd.Error()
		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
	}

	returnValue.Row = row

	w.Header().Set("Content-Type", "application/json")
	return func() { SendSuccess(w, returnValue) }
}

func GetFunc_UpdateWithHTTPReturn(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) func() {

	// 声明变量
	var returnValue models.JsonRowsReturn

	itemPtr := reflect.New(modelType).Interface()

	// 不用指针取了再转的话，item会被强行变成map类型
	err := json.NewDecoder(r.Body).Decode(itemPtr)
	item := reflect.ValueOf(itemPtr).Elem().Interface()

	if err != nil {
		returnValue.Info = "Server error" + err.Error()
		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
	}

	isPassedValidation, returnValue := ValidateInputs(item)

	if !isPassedValidation {
		return func() { SendJsonError(w, http.StatusBadRequest, returnValue) } // 400
	}

	addRow := reflect.ValueOf(repo).MethodByName("UpdateRow")
	args := []reflect.Value{
		reflect.ValueOf(db),
		reflect.ValueOf(item),
		reflect.ValueOf(userId)}
	out := addRow.Call(args)

	rowsUpdated := out[0].Interface()
	errAdd := ParseError(out[1])

	if errAdd != nil {
		returnValue.Info = "Server error" + errAdd.Error()
		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
	}

	returnValue.Info = fmt.Sprintf("修改了了%d条记录", rowsUpdated)

	w.Header().Set("Content-Type", "application/json")
	return func() { SendSuccess(w, returnValue) }
}

func GetFunc_DeleteWithHTTPReturn(
	db *sql.DB,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) func() {

	var returnValue models.JsonRowsReturn

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	getRow := reflect.ValueOf(repo).MethodByName("DeleteRow")
	args := []reflect.Value{
		reflect.ValueOf(db),
		reflect.ValueOf(id),
		reflect.ValueOf(userId)}
	out := getRow.Call(args)

	// rowsDeleted := out[0].Interface()
	err := ParseError(out[1])

	if err != nil {

		returnValue.Info = "Server error" + err.Error()

		return func() { SendError(w, http.StatusInternalServerError, returnValue) }
		// 千万不要忘了return。否则下面的数据也会加在返回的json后
	}

	if out[0].IsZero() {

		returnValue.Info = "Not Found"
		return func() { SendError(w, http.StatusNotFound, returnValue) }
	}

	w.Header().Set("Content-Type", "application/json")
	return func() { SendSuccess(w, returnValue) }
}

// deleteRole================================================================================
// var returnValue models.JsonRowsReturn

// params := mux.Vars(r)

// id, _ := strconv.Atoi(params["id"])
// roleRepo := roleRepository.RoleRepository{}

// rowsDeleted, err := roleRepo.DeleteRow(db, id)

// if err != nil {
// 	returnValue.Info = "Server error"
// 	utils.SendError(w, http.StatusInternalServerError, returnValue) //500
// 	return
// 	// 千万不要忘了return。否则下面的数据也会加在返回的json后
// }

// if rowsDeleted == 0 {
// 	returnValue.Info = "Not Found"
// 	utils.SendError(w, http.StatusNotFound, returnValue) //404
// 	return
// }

// returnValue.Info = fmt.Sprintf("删除了%d条记录", rowsDeleted)

// w.Header().Set("Content-Type", "application/json")
// utils.SendSuccess(w, returnValue)

// updateRole================================================================================

// var role models.Role
// var returnValue models.JsonRowsReturn
// json.NewDecoder(r.Body).Decode(&role)

// // validation
// isPassedValidation, returnValue := utils.ValidateInputs(role)

// if !isPassedValidation {
// 	utils.SendJsonError(w, http.StatusBadRequest, returnValue) // 400
// 	return
// }

// roleRepo := roleRepository.RoleRepository{}

// rowsUpdated, err := roleRepo.UpdateRow(db, role)

// if err != nil {
// 	returnValue.Info = "Server error"
// 	utils.SendError(w, http.StatusInternalServerError, returnValue)
// 	return
// }

// returnValue.Info = fmt.Sprintf("新增了%d条记录", rowsUpdated)

// w.Header().Set("Content-Type", "application/json")
// utils.SendSuccess(w, returnValue)

// addRole================================================================================

// var role models.Role
// var returnValue models.JsonRowsReturn
// err := json.NewDecoder(r.Body).Decode(&role)

// fmt.Println("add role", role)

// // ------------------------------validation : 过后移出去
// isPassedValidation, returnValue := utils.ValidateInputs(role)

// if !isPassedValidation {
// 	utils.SendJsonError(w, http.StatusBadRequest, returnValue) // 400
// 	return
// }

// // validation
// // if role.Rank < 0 || role.Name == "" || role.Auth == "" {
// // 	error.Message = "Enter missing fields."
// // 	utils.SendError(w, http.StatusBadRequest, error) // 400
// // 	return
// // }

// roleRepo := roleRepository.RoleRepository{}

// // TODO: 改成传指针试试，减少一个变量
// role, err = roleRepo.AddRow(db, role)

// if err != nil {
// 	returnValue.Info = "Server error"

// 	// TODO: 错误信息返回map[string]string
// 	utils.SendError(w, http.StatusInternalServerError, returnValue)
// 	return
// }

// returnValue.Row = role

// w.Header().Set("Content-Type", "application/json")
// utils.SendSuccess(w, returnValue)

// getRoles================================================================================
// 完整query:  page=1&rowCount=5&perPage=15&totalCount=10&totalPage=2&order=desc&orderBy=id&q=%7B%22fullName%22%3A%22%E7%8E%8B%E6%80%9D%E8%81%AA%22%7D

// 权限判断
// auth := r.Header.Get("Authorization")

// fmt.Println("权限：", auth)

// // 声明变量
// var roles []models.Role
// var role models.Role
// var returnValue models.JsonRowsReturn
// roles = []models.Role{}
// roleRepo := roleRepository.RoleRepository{} // 数据库语句

// // pagination & search
// pagination := utils.GetPagination(r)
// searchTerms := utils.GetSearchTerms(r)

// // 运行数据库语句: db, model, array of model, pagination, query
// roles, err := roleRepo.GetRoles(db, role, roles, &pagination, searchTerms)

// if err != nil {
// 	returnValue.Info = "Server error"
// 	utils.SendError(w, http.StatusInternalServerError, returnValue)
// 	return
// }

// // 准备返回值
// returnValue.Pagination = pagination
// returnValue.SearchTerms = searchTerms
// returnValue.Rows = roles

// w.Header().Set("Content-Type", "application/json")
// utils.SendSuccess(w, returnValue)

// getRole================================================================================
// var role models.Role
// var returnValue models.JsonRowsReturn

// // 接参数
// params := mux.Vars(r)

// id, _ := strconv.Atoi(params["id"])
// roleRepo := roleRepository.RoleRepository{}
// role, err := roleRepo.GetRow(db, id)

// if err != nil {
// 	if err == sql.ErrNoRows {
// 		returnValue.Info = "Not Found"
// 		utils.SendError(w, http.StatusNotFound, returnValue)
// 		return
// 	} else {
// 		returnValue.Info = "Server error"
// 		utils.SendError(w, http.StatusInternalServerError, returnValue)
// 		return
// 	}
// }

// returnValue.Row = role

// w.Header().Set("Content-Type", "application/json")
// utils.SendSuccess(w, returnValue)
