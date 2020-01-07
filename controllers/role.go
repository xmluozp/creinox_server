package controllers

import (
	"github.com/xmluozp/creinox_server/models"
	roleRepository "github.com/xmluozp/creinox_server/repository/role"
	"github.com/xmluozp/creinox_server/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Controller struct{}

func (c Controller) GetRoles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var roles []models.Role
		var role models.Role
		var error models.Error
		roles = []models.Role{}

		roleRepo := roleRepository.RoleRepository{}

		roles, err := roleRepo.GetRoles(db, role, roles)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, roles)
	}
}

func (c Controller) GetRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, role)
	}
}
func (c Controller) AddRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// 用来记录add以后自增长的id
		var roleID int

		var role models.Role
		var error models.Error
		err := json.NewDecoder(r.Body).Decode(&role)

		fmt.Println("add")

		// validation
		if role.Author == "" || role.Title == "" || role.Year == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error) // 400
			return
		}

		roleRepo := roleRepository.RoleRepository{}
		roleID, err = roleRepo.AddRole(db, role)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, roleID)

	}
}

func (c Controller) UpdateRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		var error models.Error
		json.NewDecoder(r.Body).Decode(&role)

		// validation
		if role.Author == "" || role.Title == "" || role.Year == "" {
			error.Message = "Enter missing fields."
			utils.SendError(w, http.StatusBadRequest, error) // 400
			return
		}

		roleRepo := roleRepository.RoleRepository{}
		rowsUpdated, err := roleRepo.UpdateRole(db, role)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controller) DeleteRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
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
}
