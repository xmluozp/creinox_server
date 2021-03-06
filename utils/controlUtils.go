package utils

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/models"
)

func GetFunc_FetchListHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	methodName string, // repo方法名
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, err error) {

	// items := reflect.Zero(reflect.SliceOf(modelType)).Interface()
	// item := reflect.New(modelType).Elem().Interface()

	// query sample:  page=1&rowCount=5&perPage=15&totalCount=10&totalPage=2&order=desc&orderBy=id&q=%7B%22fullName%22%3A%22%E7%8E%8B%E6%80%9D%E8%81%AA%22%7D

	pagination := GetPagination(r)
	searchTerms := GetSearchTerms(r)

	// 因为map是指针，所以searchTerms会被内页的方法篡改，这里复制一份，作为干净的searchTerm返回用
	searchTerms_returnValue := copyMap(searchTerms)

	fmt.Println("搜索条件", searchTerms)

	gerRows := reflect.ValueOf(repo).MethodByName(methodName)
	args := []reflect.Value{
		reflect.ValueOf(mydb),
		// reflect.ValueOf(item),
		// reflect.ValueOf(items),
		reflect.ValueOf(pagination),
		reflect.ValueOf(searchTerms),
		reflect.ValueOf(userId)}

	fmt.Println("call values:", args)
	// 运行数据库语句: db, model, array of model, pagination, query
	out := gerRows.Call(args)
	rows := out[0].Interface()
	paginationOut := out[1].Interface()
	err = ParseError(out[2])

	if err != nil {

		fmt.Println("ge rows error: ", err)
		returnValue.Info = "Server error" + err.Error()
		return http.StatusInternalServerError, returnValue, err
	}

	// 如果数据是空的，返回一个空数组（避免反复取数据）

	// 准备返回值
	returnValue.Pagination = paginationOut
	returnValue.SearchTerms = searchTerms_returnValue
	returnValue.Rows = rows
	return http.StatusOK, returnValue, nil

}

func GetFunc_RowsWithHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, err error) {

	return GetFunc_FetchListHTTPReturn(mydb, w, r, modelType, "GetRows", repo, userId)
}

func GetFunc_RowWithHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, err error) {

	return GetFunc_RowWithHTTPReturn_MethodName(mydb, w, r, modelType, "GetRow", repo, userId)
}

func GetFunc_RowWithHTTPReturn_MethodName(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	methodName string, // repo方法名
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, err error) {
	// 接参数
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	getRow := reflect.ValueOf(repo).MethodByName(methodName)
	args := []reflect.Value{
		reflect.ValueOf(mydb),
		reflect.ValueOf(id),
		reflect.ValueOf(userId)}
	out := getRow.Call(args)

	row := out[0].Interface()
	err = ParseError(out[1])

	if err != nil {

		fmt.Println("取单数据出错", err.Error())
		if err == sql.ErrNoRows {
			returnValue.Info = "找不到记录. " + err.Error()
			return http.StatusNotFound, returnValue, err

		} else {
			returnValue.Info = "服务器错误. " + err.Error()
			return http.StatusNotFound, returnValue, err
		}
	}

	returnValue.Row = row

	return http.StatusOK, returnValue, nil

}

// Here returns function to shows message to the front. In case we want to pop a message afterall instead of pop it on halfway
func GetFunc_AddWithHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型

	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, returnItem interface{}, err error) {

	itemPtr := reflect.New(modelType).Interface()

	// 不用指针取了再转的话，item会被强行变成map类型
	err = json.NewDecoder(r.Body).Decode(itemPtr)
	item := reflect.ValueOf(itemPtr).Elem().Interface()

	if err != nil {
		fmt.Println("Insert error on controller: ", err)
		returnValue.Info = "服务器错误. " + err.Error()
		return http.StatusInternalServerError, returnValue, item, err
	}

	// ====================================== 保存数据库
	status, returnValue, row, err := addDateBase(mydb, item, repo, userId)

	return status, returnValue, row, err
}

func GetFunc_UpdateWithHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, returnItem interface{}, err error) {

	return GetFunc_UpdateWithHTTPReturn_MethodName(mydb, w, r, modelType, "UpdateRow", repo, userId)
}

func GetFunc_UpdateWithHTTPReturn_MethodName(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	methodName string,
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, returnItem interface{}, err error) {

	itemPtr := reflect.New(modelType).Interface()

	// 不用指针取了再转的话，item会被强行变成map类型
	err = json.NewDecoder(r.Body).Decode(itemPtr)

	item := reflect.ValueOf(itemPtr).Elem().Interface()

	if err != nil {
		returnValue.Info = "服务器错误. " + err.Error()
		return http.StatusInternalServerError, returnValue, item, err
	}

	//  ---------------------------------------  保存数据库
	status, returnValue, err = updateDateBase(mydb, item, repo, methodName, userId)
	//  ---------------------------------------
	return status, returnValue, item, err
}

// TODO: 批量删除
// func GetFunc_DeleteWithHTTPReturn_Multiple(
// 	mydb models.MyDb,
// 	w http.ResponseWriter,
// 	r *http.Request,
// 	modelType reflect.Type, // 数据模型
// 	repo interface{},
// 	userId int) (status int, returnValue models.JsonRowsReturn, returnItem interface{}, err error) {

// 	var idList []int

// 	err = json.NewDecoder(r.Body).Decode(&idList)

// 	fmt.Println("idList", idList)

// 	deleteRows := reflect.ValueOf(repo).MethodByName("DeleteRows")
// 	args := []reflect.Value{
// 		reflect.ValueOf(mydb),
// 		reflect.ValueOf(idList),
// 		reflect.ValueOf(userId)}
// 	out := deleteRows.Call(args)

// 	rowsDeleted := out[0].Interface()
// 	err = ParseError(out[1])

// 	// if err != nil {

// 	// 	returnValue.Info = "Server error" + err.Error()

// 	// 	return http.StatusInternalServerError, returnValue, err
// 	// 	// 千万不要忘了return。否则下面的数据也会加在返回的json后
// 	// }

// 	// if out[0].IsZero() {

// 	// 	returnValue.Info = "Not Found"
// 	// 	return http.StatusNotFound, returnValue, err
// 	// }

// 	return http.StatusOK, returnValue, rowsDeleted, err
// }

func GetFunc_DeleteWithHTTPReturn(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type, // 数据模型
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, deletedItem interface{}, err error) {

	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])

	getRow := reflect.ValueOf(repo).MethodByName("DeleteRow")
	args := []reflect.Value{
		reflect.ValueOf(mydb),
		reflect.ValueOf(id),
		reflect.ValueOf(userId)}
	out := getRow.Call(args)

	// rowsDeleted := out[0].Interface()
	err = ParseError(out[1])

	if err != nil {

		returnValue.Info = "Server error" + err.Error()

		return http.StatusInternalServerError, returnValue, nil, err
		// 千万不要忘了return。否则下面的数据也会加在返回的json后
	}

	if out[0].IsNil() {
		returnValue.Info = "Not Found"
		return http.StatusNotFound, returnValue, nil, err
	}

	row := out[0].Interface()

	return http.StatusOK, returnValue, row, nil
}

// Here returns function to shows message to the front. In case we want to pop a message afterall instead of pop it on halfway
func GetFunc_AddWithHTTPReturn_FormData(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type,
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, itemReturn interface{}, files map[string][]byte, err error) {

	status, returnValue, item, files, err := DecodeFormData(r, modelType)

	if err != nil {
		fmt.Println("decode err:", err)
		return status, returnValue, nil, nil, err
	}

	//  ---------------------------------------  保存数据库
	status, returnValue, row, err := addDateBase(mydb, item, repo, userId)

	if err != nil {
		fmt.Println("update err:", err)
		return status, returnValue, nil, nil, err
	}

	//  ---------------------------------------
	return http.StatusOK, returnValue, row, files, nil
}

// Here returns function to shows message to the front. In case we want to pop a message afterall instead of pop it on halfway
func GetFunc_UpdateWithHTTPReturn_FormData(
	mydb models.MyDb,
	w http.ResponseWriter,
	r *http.Request,
	modelType reflect.Type,
	repo interface{},
	userId int) (status int, returnValue models.JsonRowsReturn, itemReturn interface{}, files map[string][]byte, err error) {

	status, returnValue, item, files, err := DecodeFormData(r, modelType)

	if err != nil {
		fmt.Println("decode err:", err)
		return status, returnValue, nil, nil, err
	}

	// --------------------------------------- 保存数据库
	status, returnValue, err = updateDateBase(mydb, item, repo, "UpdateRow", userId)

	if err != nil {
		fmt.Println("update err:", err)
		return status, returnValue, nil, nil, err
	}

	// ---------------------------------------
	return http.StatusOK, returnValue, item, files, nil
}

// ====================== private ======================

func updateDateBase(
	mydb models.MyDb,
	item interface{}, // 数据模型
	repo interface{},
	methodName string,
	userId int) (int, models.JsonRowsReturn, error) {

	isPassedValidation, returnValue := ValidateInputs(item)

	if !isPassedValidation {
		err := errors.New("表单验证失败")
		return http.StatusBadRequest, returnValue, err // 400
	}

	addRow := reflect.ValueOf(repo).MethodByName(methodName)
	args := []reflect.Value{
		reflect.ValueOf(mydb),
		reflect.ValueOf(item),
		reflect.ValueOf(userId)}
	out := addRow.Call(args)

	rowsUpdated := out[0].Interface()
	errAdd := ParseError(out[1])

	if errAdd != nil {
		returnValue.Info = "Server error" + errAdd.Error()
		return http.StatusInternalServerError, returnValue, errAdd
	}
	returnValue.Info = fmt.Sprintf("更新了%d条记录", rowsUpdated)
	returnValue.Row = item

	return http.StatusAccepted, returnValue, nil
}

func addDateBase(
	mydb models.MyDb,
	item interface{}, // 数据模型
	repo interface{},
	userId int) (int, models.JsonRowsReturn, interface{}, error) {

	isPassedValidation, returnValue := ValidateInputs(item)

	fmt.Println("表单验证", isPassedValidation, "returnValue", returnValue)

	if !isPassedValidation {
		err := errors.New("表单验证失败")
		return http.StatusBadRequest, returnValue, nil, err // 400
	}

	addRow := reflect.ValueOf(repo).MethodByName("AddRow")
	args := []reflect.Value{
		reflect.ValueOf(mydb),
		reflect.ValueOf(item),
		reflect.ValueOf(userId)}
	out := addRow.Call(args)

	row := out[0].Interface()
	errAdd := ParseError(out[1])

	if errAdd != nil {
		returnValue.Info = "Server error" + errAdd.Error()
		return http.StatusInternalServerError, returnValue, nil, errAdd
	}

	returnValue.Row = row

	return http.StatusOK, returnValue, row, nil
}

// 把前端传来的multipart/formData变成 string: []byte的map
func DecodeFormData(
	r *http.Request,
	modelType reflect.Type,
) (status int, returnValue models.JsonRowsReturn, itemReturn interface{}, files map[string][]byte, err error) {

	// prepare to store the form values
	var jsonString []byte
	files = make(map[string][]byte)

	mr, err := r.MultipartReader()

	if err != nil {
		Log(err)
		returnValue.Info = "Server error" + err.Error()
		return http.StatusInternalServerError, returnValue, nil, nil, err
	}
	itemPtr := reflect.New(modelType).Interface()

	// scan
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			returnValue.Info = "Server error" + err.Error()
			return http.StatusInternalServerError, returnValue, nil, nil, err
		}

		data, _ := ioutil.ReadAll(part)

		if part.FileName() != "" {

			files[part.FormName()] = data

		} else if part.FormName() == "doc" { // only a convention: 前端传过来的doc下面放的是整个json，直接用它unmrshal

			jsonString = data
			err = json.Unmarshal(jsonString, itemPtr)

			if err != nil {
				returnValue.Info = "Server error" + err.Error()
				return http.StatusInternalServerError, returnValue, nil, nil, err
			}
		}
	}

	item := reflect.ValueOf(itemPtr).Elem().Interface()

	return http.StatusOK, returnValue, item, files, err
}

func FetchPrintPathAndId(r *http.Request) (int, string, string) {

	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])
	templateFolder, _ := params["templateFolder"]
	template, _ := params["template"]
	printFormat, _ := params["printFormat"]
	_path := "." + "/templates/" + templateFolder + "/" + template
	return id, _path, printFormat
}

// 从模板和数据源生成打印文件。插入回 writer里
func PrintFromTemplate(
	w http.ResponseWriter,
	m map[string]interface{},
	_path string,
	printFormat string,
	userId int) (err error) {

	_cache := "cache" + RandomString(8)

	os.Mkdir(_cache, os.ModePerm)
	defer os.RemoveAll(_cache)

	ext := filepath.Ext(_path)

	// 判断扩展名如果不对就终止
	if ext != ".xlsx" {
		return ErrorFromString("template has to be: .xlsx file")
	}

	// 生成打印文件
	xls := XlsxTemplate{}

	resultPath := _cache + "/temporary.xlsx"
	xls.PrintOut(_path, resultPath, m)

	// --------------------------------------------------------- 看是要pdf还是excel，如果pdf就调用程序转
	if printFormat == "pdf" {
		fullpath, _ := filepath.Abs(resultPath)

		var fireparams []string

		fireparams = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", _cache + "/pdf/", fullpath}
		interactiveToexec("cmd", fireparams)

		// resultPath 针对pdf改成是pdf文件
		resultPath = _cache + "/pdf/temporary.pdf"

		// http.ServeFile(w, r, "./cache/pdf/test.pdf")
		// err := os.RemoveAll("./cache")
	}
	// --------------------------------------------------------- 转成pdf

	file, _ := os.Open(resultPath)
	defer file.Close()
	buff, _ := ioutil.ReadAll(file)

	w.Write(buff)

	return err
}

// deleteRole================================================================================
// var returnValue models.JsonRowsReturn

// params := mux.Vars(r)

// id, _ := strconv.Atoi(params["id"])
// roleRepo := roleRepository.RoleRepository{}

// rowsDeleted, err := roleRepo.DeleteRow(mydb, id)

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

// rowsUpdated, err := roleRepo.UpdateRow(mydb, role)

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
// role, err = roleRepo.AddRow(mydb, role)

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
// roles, err := roleRepo.GetRoles(mydb, role, roles, &pagination, searchTerms)

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
// role, err := roleRepo.GetRow(mydb, id)

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
