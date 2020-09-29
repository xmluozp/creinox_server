package mytestController

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	// repository "github.com/xmluozp/creinox_server/repository/bankAccount"
	"github.com/gobuffalo/nulls"
	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}

// type modelName = models.BankAccount

var authName = ""

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

}

type TestModel struct {
	Name  nulls.String `json:"name"`
	Money nulls.Int    `json:"money"`
}

type MyDbModel struct {
	Db *sql.DB
	Tx *sql.Tx
}

type MakeErr struct {
	Test *sql.Tx
}

// =============================================== HTTP REQUESTS

// ===========================================================事务回滚
func (c Controller) MytestTx(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		v, _ := params["v"]

		mydb := MyDbModel{Db: db}

		// 初始值
		result, err := mydb.Db.Exec(`
			UPDATE product SET name = ?
			WHERE id = 16538
		`, "begin")

		tx, err := db.Begin()
		defer tx.Rollback()

		mydb.Tx = tx

		result, err = mydb.Tx.Exec(`
			UPDATE product SET name = ?
			WHERE id = 16538
		`, v)

		result, err = mydb.Tx.Exec(`
			UPDATE product SET name = ?
			WHERE id = 16538
		`, v+"jjj")

		// =============================== 制造一个错误
		// me := MakeErr{}
		// result, err = db.Exec(`
		// 	lkjlkjlkj
		// `, v+"jjj")

		// if err != nil {
		// 	fmt.Println("出错了")
		// 	return
		// }
		// ===============================

		// result, err = mydb.Tx.Exec(`
		// 	UPDATE product SET name = ?
		// 	WHERE id = 16538
		// `, v+"kkk")

		result, err = PassTx(mydb, v)

		tx.Commit()

		fmt.Println("mydb:", mydb)

		fmt.Println("Test output:", result, err)
	}
}

func PassTx(mydb MyDbModel, v string) (result sql.Result, err error) {

	return mydb.Tx.Exec(`
	UPDATE product SET name = ?
	WHERE id = 16538
	`, v+"kjk")
}

// ===========================================================事务回滚

// =============================================== HTTP REQUESTS
func (c Controller) Mytest(mydb models.MyDb) http.HandlerFunc {
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
func (c Controller) MytestApp(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ====================================================== 申请
		fmt.Println("==================第一层: 接到了原始的request")
		fmt.Println(r.URL.String())
		fmt.Println(r.Body)
		fmt.Println(r.Header)

		// ========== 把Application提取出来
		bodyBytes, err := ioutil.ReadAll(r.Body)

		// 整个表单存到数据库里
		content := hex.EncodeToString(bodyBytes)

		fmt.Println(content)

		// ====================================================== 审批
		newBodyBytes, _ := hex.DecodeString("7b0d0a226e616d65223a2022e5bca0e4b889222c0d0a226d6f6e6579223a203130300d0a7d0d0a")

		// 完整链接在前端生成. 存到 attemptUrl
		req, err := http.NewRequest(http.MethodPost, "http://192.168.0.10:8000/api/testAppReceive/asdf123", bytes.NewBuffer(newBodyBytes))

		// req.Header.Set("Content-type", "application/json")

		req.Header = r.Header.Clone()

		// 当前权限
		pass, userId := auth.CheckAuth(mydb, w, r, "")
		fmt.Println("========当前权限")
		fmt.Println(pass, "userId:", userId)

		// 加权限
		// TODO: 取出审批人的token，作为权限加进去
		//  req.Header.Set("Authorization", grant)
		// req.Header.Set("test", "test")
		req.Header.Set("Authorization", "okok")

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
func (c Controller) MytestAppReceive(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("==================第二层")
		fmt.Println(r.URL)
		fmt.Println(r.Body)
		fmt.Println(r.Header)

		fmt.Println("==================解开")

		pass, userId := auth.CheckAuth(mydb, w, r, "")
		fmt.Println("========确认以后审批人权限")
		fmt.Println(pass, "userId:", userId)

		var item TestModel
		json.NewDecoder(r.Body).Decode(&item)
		fmt.Println("item:", item)

		// 看权限

	}
}

func Hello() string {
	return "Hello, world"
}
