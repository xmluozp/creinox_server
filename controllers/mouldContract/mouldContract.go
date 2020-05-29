package mouldContractController

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"

	xlst "github.com/ivahaev/go-xlsx-templater"

	"github.com/gorilla/mux"
	"github.com/xmluozp/creinox_server/auth"
	folderController "github.com/xmluozp/creinox_server/controllers/folder"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/mouldContract"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.MouldContract

var authName = "mouldcontract"

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

		status, returnValue, itemReturn, err := utils.GetFunc_DeleteWithHTTPReturn(db, w, r, reflect.TypeOf(item), repo, userId)

		contract := itemReturn.(modelName)

		if err == nil {
			// 删除folder (folder下面images的删除在folder处理)
			folderCtrl := folderController.Controller{}
			err = folderCtrl.Delete(db, contract.Gallary_folder_id.Int, userId)

			if err != nil {
				returnValue.Info = "删除产品开发合同对应图库失败" + err.Error()
				utils.SendJson(w, http.StatusFailedDependency, returnValue, err)
			}
		}

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

		fmt.Println("打印开始")

		_cache := "cache" + utils.RandomString(8)
		os.Mkdir(_cache, os.ModePerm)
		// 打印用数据=========================/
		item, _ := repo.GetRow(db, id, 0)
		fmt.Println(item.Code)

		// 把打印用数据弄到xsl ===============/
		_path := "." + "/templates/" + templateFolder + "/" + template
		ext := filepath.Ext(_path)

		// 判断扩展名如果不对就终止
		if ext != ".xlsx" {
			w.Write([]byte("template has to be: xlsx or xls"))
			return
		}

		// TODO: 现在是直接传file回去，改成传回处理过的文件。hardcoding处理
		// 生成打印文件
		ctx := map[string]interface{}{
			"name": "Github User"}
		ctx["code"] = "okokokok"

		// doc := xlst.New()

		xlsfile, _ := os.Open(_path)
		defer xlsfile.Close()
		xlsbuff, _ := ioutil.ReadAll(xlsfile)

		doc, err := xlst.NewFromBinary(xlsbuff)

		// 这个和上面那个readbinary应该一样的
		doc.ReadTemplate(_path)

		if err != nil {
			fmt.Println("read xls from template error,", err)
		}

		err = doc.Render(ctx)

		if err != nil {
			fmt.Println("error when render", err)
		}

		doc.Save(_cache + "/temporary.xlsx")

		fullpath, _ := filepath.Abs(_cache + "/temporary.xlsx")

		fmt.Println("fullpath:", fullpath)

		var fireparams []string

		fireparams = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", _cache + "/pdf/", fullpath}
		interactiveToexec("cmd", fireparams)

		resultPath := _cache + "/pdf/temporary.pdf"

		// http.ServeFile(w, r, "./cache/pdf/test.pdf")
		// err := os.RemoveAll("./cache")

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

// dataByte, _ := ioutil.ReadFile("html/pdf.html")
// dataStr := string(dataByte)
// pdfUrl := "pdf_asset/" + path.Base(filePath)
// dataStr = strings.Replace(dataStr, "{{url}}", pdfUrl, -1)
// dataByte = []byte(dataStr)
// return dataByte

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
