package userController

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/nulls"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	roleRepository "github.com/xmluozp/creinox_server/repository/role"
	repository "github.com/xmluozp/creinox_server/repository/user"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.User

var authName = "user"

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	// 浏览不设权限，因为要下拉
	pass, userId := auth.CheckAuth(db, w, r, "")
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItems_DropDown(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// ...
}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName

	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItemsForLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	userId := 0

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_FetchListHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRowsForLogin", repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	pass, userId := auth.CheckAuth(db, w, r, authName)
	if !pass {
		return
	}

	// 从param里取出id，模板所在目录，打印格式（做在utils里面是为了方便日后修改）
	id, _path, printFormat := utils.FetchPrintPathAndId(r)

	// 生成打印数据(取map出来而不是item，是为了方便篡改)
	repo := repository.Repository{}
	dataSource, err := repo.GetPrintSource(db, id, userId)

	if err != nil {
		w.Write([]byte("error on generating source data," + err.Error()))
	}

	// 直接打印到writer(因为打印完毕需要删除cache，所以要在删除之前使用writer)
	err = utils.PrintFromTemplate(w, dataSource, _path, printFormat, userId)

	if err != nil {
		w.Write([]byte("error on printing," + err.Error()))
		return
	}
}

// =============================================== HTTP REQUESTS
func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, db)
	}
}

func (c Controller) GetItems_DropDown(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems_DropDown(w, r, db)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItem(w, r, db)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_AddItem(w, r, db)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem(w, r, db)
	}
}

func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_DeleteItem(w, r, db)
	}
}

func (c Controller) Print(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_Print(w, r, db)
	}
}

func (c Controller) GetItemsForLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItemsForLogin(w, r, db)
	}
}

// =============================================== customized: login
func (c Controller) Login(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var creds auth.Credentials
		var user models.User
		var returnValue models.JsonRowsReturn
		repo := repository.Repository{}
		roleRepo := roleRepository.Repository{}

		// Get the JSON body and decode into credentials
		err := json.NewDecoder(r.Body).Decode(&creds)

		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			returnValue.Info = err.Error()
			utils.SendError(w, http.StatusUnauthorized, returnValue)
			return
		}

		// Get the expected password from our in memory map
		// expectedPassword, ok := users[creds.Username]

		fmt.Println(creds.UserName)
		// 传统登录验证，取出userid
		user, err = repo.GetLoginRow(db, creds.UserName)
		expectedPassword := user.Password.String
		fmt.Println(expectedPassword)

		if err != nil || auth.CheckPasswordHash(expectedPassword, creds.Password) {
			fmt.Println("登录错误", err)
			// 找不到用户或者密码对不上
			returnValue.Info = "用户名或者密码错误，或用户被禁用"
			utils.SendError(w, http.StatusUnauthorized, returnValue)
			return
		}

		// 如果登录成功 取角色：
		role, err := roleRepo.GetRow(db, user.Role_id.Int, 0)

		if err != nil {
			// 找不到用户或者密码对不上
			returnValue.Info = "角色出错"
			utils.SendError(w, http.StatusUnauthorized, returnValue)
			return
		}

		// Declare the expiration time of the token
		// here, we have kept it as 5 minutes
		fmt.Println("从角色取到的auth", role)

		expirationTime := time.Now().Add(10000 * time.Hour)

		// Create the JWT claims, which includes the username and expiry time
		claims := &auth.Claims{
			UserId:   user.ID.Int,
			UserName: creds.UserName,
			Auth:     role.Auth.String,
			StandardClaims: jwt.StandardClaims{
				// In JWT, the expiry time is expressed as unix milliseconds
				ExpiresAt: expirationTime.Unix(),
			},
		}

		// Declare the token with the algorithm used for signing, and the claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// Create the JWT string
		tokenString, err := token.SignedString(auth.JwtKey)

		if err != nil {
			returnValue.Info = "Server error" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		// 设置新的token，ip，登录时间
		user.Token = nulls.NewString(tokenString)
		user.IP = nulls.NewString(GetIP(r))

		// 存入数据库
		_, err = repo.UpdateLoginRow(db, user)

		// 例子中返回前端的还有过期时间
		// http.SetCookie(w, &http.Cookie{
		// 	Name:    "token",
		// 	Value:   tokenString,
		// 	Expires: expirationTime,
		// })

		// 登录后的user信息返回给前端，包括token

		if err != nil {
			returnValue.Info = "Server error" + err.Error()
			utils.SendError(w, http.StatusInternalServerError, returnValue)
			return
		}

		user.RoleItem = role
		user.Password = nulls.NewString("")

		// returnValue.Row = rowsUpdated

		// w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, user) // 成功了直接返回row
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
