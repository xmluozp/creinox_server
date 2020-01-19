package controllers

import (
	"database/sql"
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/models"
	roleRepository "github.com/xmluozp/creinox_server/repository/role"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}

func (c Controller) GetRoles(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// auth := r.Header.Get("Authorization")

		var role models.Role
		roleRepo := roleRepository.RoleRepository{}
		f := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(role), roleRepo)
		f()
	}
}

func (c Controller) GetRole(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		roleRepo := roleRepository.RoleRepository{}
		f := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(role), roleRepo)
		f()
	}
}
func (c Controller) AddRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		roleRepo := roleRepository.RoleRepository{}
		f := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(role), roleRepo)
		f()
	}
}

func (c Controller) UpdateRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		roleRepo := roleRepository.RoleRepository{}
		f := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(role), roleRepo)
		f()
	}
}

func (c Controller) DeleteRole(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var role models.Role
		roleRepo := roleRepository.RoleRepository{}
		f := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(role), roleRepo)
		f()
	}
}
