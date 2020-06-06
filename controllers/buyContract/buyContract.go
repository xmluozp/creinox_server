package buyContractController

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/buyContract"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.BuyContract

var authName = "buycontract"

// =============================================== basic CRUD

func (c Controller) GetItems(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) GetItem(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}
func (c Controller) AddItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) UpdateItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}
		status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) DeleteItem(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

func (c Controller) Print(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// pass, userId := auth.CheckAuth(db, w, r, authName)
		// if !pass {
		// 	return
		// }
		params := mux.Vars(r)
		id, _ := strconv.Atoi(params["id"])
		templateFolder, _ := params["templateFolder"]
		template, _ := params["template"]

		repo := repository.Repository{}

		_cache := "cache" + utils.RandomString(8)
		os.Mkdir(_cache, os.ModePerm)

		// 打印用数据 =========================
		// 这里得自己处理. getrow得手写。尤其是有子列表的/
		// ** 取出来以后过滤，nil丢掉
		// 取数据的部分不能在model，要移到repo。因为取下属产品只能从repo.
		// : 所有repo弄个 ** repo.GetPrintSource(). 默认原地返回GetRow。
		// 有特殊处理再处理（比如这个buyContract需要取产品列表）

		item, _ := repo.GetPrintSource(db, id, 0)
		byteArray, err := json.Marshal(item)
		var m map[string]interface{}
		err = json.Unmarshal(byteArray, &m)
		m = utils.ClearNil(m)

		// item转json
		fmt.Println("反序列化：=====================")
		fmt.Println(m)

		// 把打印用数据弄到xsl ===============/
		_path := "." + "/templates/" + templateFolder + "/" + template
		ext := filepath.Ext(_path)

		// 判断扩展名如果不对就终止
		if ext != ".xlsx" {
			w.Write([]byte("template has to be: xlsx or xls"))
			return
		}

		// 生成打印文件
		xls := utils.XlsxTemplate{}

		resultPath := _cache + "/temporary.xlsx"
		xls.PrintOut(_path, resultPath, m)

		// --------------------------------------------------------- 转成pdf TODO: 加个判断，看是要pdf还是excel
		fullpath, _ := filepath.Abs(resultPath)

		var fireparams []string

		fireparams = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", _cache + "/pdf/", fullpath}
		interactiveToexec("cmd", fireparams)

		// resultPath 针对pdf改成是pdf文件
		resultPath = _cache + "/pdf/temporary.pdf"

		// http.ServeFile(w, r, "./cache/pdf/test.pdf")
		// err := os.RemoveAll("./cache")
		// --------------------------------------------------------- 转成pdf

		file, _ := os.Open(resultPath)
		defer file.Close()
		buff, _ := ioutil.ReadAll(file)

		w.Write(buff)

		file.Close()
		err = os.RemoveAll(_cache)
		if err != nil {
			fmt.Println("error from remove cache", err)
		}
	}
}

// 运行命令行来进行pdf转换
func interactiveToexec(commandName string, params []string) (string, bool) {
	cmd := exec.Command(commandName, params...)
	buf, err := cmd.Output()
	log.Println(cmd.Args)
	w := bytes.NewBuffer(nil)
	cmd.Stderr = w
	log.Printf("%s\n", w)
	if err != nil {
		log.Println("Error: <", err, "> when exec command read out buffer")
		return "", false
	} else {
		return string(buf), true
	}
}

// =============================================== customized

func (c Controller) GetLast(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(db, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_FetchRowHTTPReturn(db, w, r, reflect.TypeOf(item), "GetRow_GetLast", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}
