package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xmluozp/creinox_server/models"
	roleRepository "github.com/xmluozp/creinox_server/repository/role"
	"github.com/xmluozp/creinox_server/utils"

	"github.com/gorilla/mux"
)

type Controller struct{}

func (c Controller) GetRoles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 为了在这个页面的其他地方可以复用
		getRoles(db, w, r)
	}
}

func (c Controller) GetRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		getRole(db, w, r)

	}
}
func (c Controller) AddRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		addRole(db, w, r)
	}
}

func (c Controller) UpdateRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		updateRole(db, w, r)
	}
}

func (c Controller) DeleteRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		deleteRole(db, w, r)
	}
}

// =========
func getRoles(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// 完整query:  page=1&rowCount=5&perPage=15&totalCount=10&totalPage=2&order=desc&orderBy=id&q=%7B%22fullName%22%3A%22%E7%8E%8B%E6%80%9D%E8%81%AA%22%7D

	// 权限判断
	auth := r.Header.Get("Authorization")

	fmt.Println("权限：", auth)

	// 基本变量
	var roles []models.Role
	var role models.Role
	var error models.Error
	roles = []models.Role{}
	roleRepo := roleRepository.RoleRepository{} // 数据库语句

	// pagination & search
	pagination := utils.GetPagination(r)
	searchTerms := utils.GetSearchTerms(r)

	// 运行数据库语句: db, model, array of model, pagination, query
	roles, err := roleRepo.GetRoles(db, role, roles, &pagination, searchTerms)

	if err != nil {
		error.Message = "Server error"
		utils.SendError(w, http.StatusInternalServerError, error)
		return
	}

	var returnValue models.JsonRowsReturn
	returnValue.Pagination = pagination
	returnValue.SearchTerms = searchTerms
	returnValue.Rows = roles

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, returnValue)
}

func getRole(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	// 接参数
	params := mux.Vars(r)

	// 需要用到的变量
	var role models.Role
	var error models.Error

	id, _ := strconv.Atoi(params["id"])
	roleRepo := roleRepository.RoleRepository{}
	role, err := roleRepo.GetRole(db, id)

	if err != nil {
		if err == sql.ErrNoRows {
			error.Message = "Not Found"
			utils.SendError(w, http.StatusNotFound, error)
			return
		} else {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

	}

	var returnValue models.JsonRowsReturn

	returnValue.Row = role

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, returnValue)
}

func addRole(db *sql.DB, w http.ResponseWriter, r *http.Request) {

	// 用来记录add以后自增长的id: 返回全部
	// var roleID int

	var role models.Role
	var error models.Error
	var returnValue models.JsonRowsReturn
	err := json.NewDecoder(r.Body).Decode(&role)

	fmt.Println("add role", role)

	// ------------------------------validation : 过后移出去
	isPassedValidation, returnValue := utils.ValidateInputs(role)

	if !isPassedValidation {
		utils.SendJsonError(w, http.StatusBadRequest, returnValue) // 400
		return
	}

	// validation
	// if role.Rank < 0 || role.Name == "" || role.Auth == "" {
	// 	error.Message = "Enter missing fields."
	// 	utils.SendError(w, http.StatusBadRequest, error) // 400
	// 	return
	// }

	roleRepo := roleRepository.RoleRepository{}

	// TODO: 改成传指针试试，减少一个变量
	role, err = roleRepo.AddRole(db, role)

	if err != nil {
		error.Message = "Server error"

		// TODO: 错误信息返回map[string]string
		utils.SendError(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, role)
}

func updateRole(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var role models.Role
	var error models.Error
	json.NewDecoder(r.Body).Decode(&role)

	// validation
	if role.Rank == 0 || role.Name == "" || role.Auth == "" {
		error.Message = "Enter missing fields."
		utils.SendError(w, http.StatusBadRequest, error) // 400
		return
	}

	roleRepo := roleRepository.RoleRepository{}

	// TODO: 和add一样错误处理
	rowsUpdated, err := roleRepo.UpdateRole(db, role)

	if err != nil {
		fmt.Println(err)
		error.Message = "Server error"
		utils.SendError(w, http.StatusInternalServerError, error)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, rowsUpdated)
}

func deleteRole(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, _ := strconv.Atoi(params["id"])
	roleRepo := roleRepository.RoleRepository{}
	var error models.Error

	rowsDeleted, err := roleRepo.DeleteRole(db, id)

	if err != nil {
		error.Message = "Server error"
		utils.SendError(w, http.StatusInternalServerError, error) //500
		return
		// 千万不要忘了return。否则下面的数据也会加在返回的json后
	}

	if rowsDeleted == 0 {
		error.Message = "Not Found"
		utils.SendError(w, http.StatusNotFound, error) //404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.SendSuccess(w, rowsDeleted)
}
