package testController

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// repository "github.com/xmluozp/creinox_server/repository/bankAccount"
	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}

// type modelName = models.BankAccount

var authName = ""

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, db *sql.DB) {

}

type TestModel struct {
	Name  nulls.String `json:"name"`
	Money nulls.Int    `json:"money"`
}

// =============================================== HTTP REQUESTS
func (c Controller) Test(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)

		v, _ := params["v"]

		a1 := utils.ParseFlightSlice(v)

		fmt.Println("Test output:", a1)
		// returnValue.Message =
		// utils.SendJson(w, http.StatusOK, "hi", nil)

	}
}

// 测试提交http请求
func (c Controller) TestApp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("==================第一层")
		fmt.Println(r.URL.String())
		fmt.Println(r.Body)
		fmt.Println(r.Header)

		// ========== body转json和还原
		// var content interface{}
		// body到json
		bodyBytes, err := ioutil.ReadAll(r.Body)

		// bodyJson, err := json.Marshal(bodyBytes)

		req, err := http.NewRequest(http.MethodPost, "http://192.168.0.10:8000/api/testAppReceive/asdf123", bytes.NewBuffer(bodyBytes))

		// req.Header.Set("Content-type", "application/json")

		req.Header = r.Header.Clone()

		// 当前权限
		pass, userId := auth.CheckAuth(db, w, r, "")
		fmt.Println("========当前权限")
		fmt.Println(pass, "userId:", userId)

		// 加权限
		// TODO: 取出审批人的token，作为权限加进去
		//  req.Header.Set("Authorization", grant)
		req.Header.Set("test", "test")

		// if err != nil {
		// 	fmt.Println("err", err)
		// }
		// fmt.Println("==================准备发送")
		// fmt.Println(req.URL)
		// fmt.Println(req.Body)
		// fmt.Println(req.Header)

		client := &http.Client{}

		_, err = client.Do(req)

		if err != nil {
			fmt.Println("出错", err)
		}

		// json到body

		// body到具体的model
		// err = json.NewDecoder(r.Body).Decode(itemPtr)

		// params := mux.Vars(r)

		// v, _ := params["v"]

		// fmt.Println("Test output:", v)
		// returnValue.Message =
		// utils.SendJson(w, http.StatusOK, "hi", nil)

	}
}

// 测试提交http请求
func (c Controller) TestAppReceive(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("==================第二层")
		fmt.Println(r.URL)
		fmt.Println(r.Body)
		fmt.Println(r.Header)

		fmt.Println("==================解开")

		pass, userId := auth.CheckAuth(db, w, r, "")
		fmt.Println("========确认以后审批人权限")
		fmt.Println(pass, "userId:", userId)

		var item TestModel
		json.NewDecoder(r.Body).Decode(&item)
		fmt.Println("item:", item)

		// 看权限

	}
}
