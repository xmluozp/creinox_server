package printControllr

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"

	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}

var authName = ""

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, db)
	}
}

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, _ := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	// status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	params := mux.Vars(r)
	templateFolder, _ := params["templateFolder"]

	// 根据文件夹名称取模板 (首字母要大写，不然读不到)
	type template struct {
		Name     string
		Path     string
		FileName string
	}

	returnRows := []template{}
	returnValue := models.JsonRowsReturn{}

	root := "templates/" + templateFolder
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if info != nil && !info.IsDir() && filepath.Ext(path) == ".xlsx" {
			tmp := template{
				Name:     strings.TrimSuffix(info.Name(), filepath.Ext(path)),
				Path:     templateFolder,
				FileName: info.Name()}
			returnRows = append(returnRows, tmp)

			fmt.Println("找到文件", info)
		}

		return err
	})

	returnValue.Rows = returnRows

	utils.SendJson(w, http.StatusAccepted, returnValue, err)
}
