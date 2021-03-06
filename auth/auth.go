package auth

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/xmluozp/creinox_server/models"
	"github.com/xmluozp/creinox_server/utils"
	"golang.org/x/crypto/bcrypt"
)

var JwtKey = make([]byte, 64)

var allAuthList = []string{
	"all",
	"setting",
	"user",
	"test",
	"region",
	"commonitem",
	"image",
	"companyinternal",
	"companyfactory",
	"companyoverseas",
	"companydomestic",
	"companyshipping",
	"product",
	"category",
	"productpurchase",
	"commodity",
	"port",
	"sellcontract",
	"buycontract",
	"mouldcontract",
	"financial",
	"paymentRequest",
	"expressOrder",
	"userLog",
	"confirm-payment",
	"confirm-contract"}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId   int    `json:"userId"`
	UserName string `json:"userName"`
	Auth     string `json:"auth"`
	jwt.StandardClaims
}

type Credentials struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// Check if current user has the authentication
func CheckAuth(mydb models.MyDb, w http.ResponseWriter, r *http.Request, authTag string) (bool, int) {

	//----------------------------------------/ only for Postman testing
	tknstr := r.Header.Get("test")
	if tknstr == "test" {
		fmt.Println("Postman测试专用权限")
		return true, 1
	}

	//----------------------------------------\ only for Postman testing

	// get token from header
	var returnValue models.JsonRowsReturn
	tknstr = r.Header.Get("Authorization")

	// get userName and auth from token
	userId, _, auth, err := GetUserNameFromToken(w, r, tknstr)

	if err != nil {
		returnValue.Info = "权限验证出错," + err.Error()
		utils.SendError(w, http.StatusInternalServerError, returnValue)
		return false, 0
	}

	authTagIdx := -1

	// get int from authTag string
	for i, v := range allAuthList {
		if v == authTag {
			authTagIdx = i
			break
		}
	}

	// turn user's auth into Array
	userAuthList := strings.Split(auth, ",")

	if authTag == "" {
		return true, userId
	} else {
		// try to match integerfied authTag with User's authList
		for _, v := range userAuthList {

			// because it's int, need to convert
			userAuthIdx, errInt := strconv.Atoi(v)

			if errInt == nil && (userAuthIdx == 0 || userAuthIdx == authTagIdx) { // 0 是全权限
				return true, userId
			}
		}
	}

	return false, userId
}

func GetUserNameFromToken(w http.ResponseWriter, r *http.Request, tknStr string) (int, string, string, error) {

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		fmt.Println("身份验证错误", err)
		return 0, "", "", err
	}

	if !tkn.Valid {
		return 0, "", "", err
	}

	// 取到用户名
	return claims.UserId, claims.UserName, claims.Auth, nil
}

func GetRankFromUser(mydb models.MyDb, userId int) int {

	rank := -1

	result := mydb.Db.QueryRow("SELECT a.rank FROM role a LEFT JOIN user b ON a.id = b.role_id WHERE b.id = ?", userId)

	result.Scan((&rank))

	return rank
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(hash, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
