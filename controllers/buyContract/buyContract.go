package buyContractController

import (
	"net/http"
	"reflect"

	"github.com/xmluozp/creinox_server/auth"
	"github.com/xmluozp/creinox_server/models"
	repository "github.com/xmluozp/creinox_server/repository/buyContract"
	"github.com/xmluozp/creinox_server/utils"
)

type Controller struct{}
type modelName = models.BuyContract

var authName = "buycontract"

// =============================================== HTTP REQUESTS

func (c Controller) GetItems(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItems(w, r, mydb)
	}
}

func (c Controller) GetItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_GetItem(w, r, mydb)
	}
}
func (c Controller) AddItem(mydb models.MyDb) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		c.C_AddItem(w, r, mydb)
	}
}

func (c Controller) UpdateItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_UpdateItem(w, r, mydb)
	}
}

func (c Controller) DeleteItem(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_DeleteItem(w, r, mydb)
	}
}

func (c Controller) Print(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c.C_Print(w, r, mydb)
	}
}

// =============================================== basic CRUD
func (c Controller) C_GetItems(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	status, returnValue, err := utils.GetFunc_RowsWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_GetItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}
	status, returnValue, err := utils.GetFunc_RowWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)
	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_AddItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_AddWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		tx.Rollback()
	} else {
		tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_UpdateItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_UpdateWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		tx.Rollback()
	} else {
		tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_DeleteItem(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	var item modelName
	repo := repository.Repository{}

	// ====================================== tx begin
	tx, err := mydb.Db.Begin()
	if err != nil {
		utils.Log(err)
		return
	}
	mydb.Tx = tx
	// ====================================== tx begin

	status, returnValue, _, err := utils.GetFunc_DeleteWithHTTPReturn(mydb, w, r, reflect.TypeOf(item), repo, userId)

	if err != nil {
		utils.Log(err, "事务运行失败")
		tx.Rollback()
	} else {
		tx.Commit()
	}

	utils.SendJson(w, status, returnValue, err)
}

func (c Controller) C_Print(w http.ResponseWriter, r *http.Request, mydb models.MyDb) {

	pass, userId := auth.CheckAuth(mydb, w, r, authName)
	if !pass {
		return
	}

	// 从param里取出id，模板所在目录，打印格式（做在utils里面是为了方便日后修改）
	id, _path, printFormat := utils.FetchPrintPathAndId(r)

	// 生成打印数据(取map出来而不是item，是为了方便篡改)
	repo := repository.Repository{}
	dataSource, err := repo.GetPrintSource(mydb, id, userId)

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

// =============================================== customized

func (c Controller) GetLast(mydb models.MyDb) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass, userId := auth.CheckAuth(mydb, w, r, authName)
		if !pass {
			return
		}

		var item modelName
		repo := repository.Repository{}

		status, returnValue, err := utils.GetFunc_RowWithHTTPReturn_MethodName(mydb, w, r, reflect.TypeOf(item), "GetRow_GetLast", repo, userId)
		utils.SendJson(w, status, returnValue, err)
	}
}

// func (c Controller) Print_Bak(mydb models.MyDb) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		pass, userId := auth.CheckAuth(mydb, w, r, authName)
// 		if !pass {
// 			return
// 		}
// 		params := mux.Vars(r)
// 		id, _ := strconv.Atoi(params["id"])
// 		templateFolder, _ := params["templateFolder"]
// 		template, _ := params["template"]
// 		printFormat, _ := params["printFormat"]

// 		repo := repository.Repository{}

// 		_cache := "cache" + utils.RandomString(8)

// 		os.Mkdir(_cache, os.ModePerm)
// 		defer os.RemoveAll(_cache)
// 		// 打印用数据 =========================
// 		// 这里得自己处理. getrow得手写。尤其是有子列表的/
// 		// ** 取出来以后过滤，nil丢掉
// 		// ** 取数据的部分不能在model，要移到repo。因为取下属产品只能从repo.
// 		// ** : 所有repo弄个 repo.GetPrintSource(). 默认原地返回GetRow。
// 		// 有特殊处理再处理（比如buyContract需要取产品列表）

// 		m, err := repo.GetPrintSource(mydb, id, userId)

// 		if err != nil {
// 			fmt.Println("生成模板出错, ", err)
// 		}
// 		// item转json
// 		// fmt.Println("反序列化：=====================")
// 		// fmt.Println(m)

// 		// 把打印用数据弄到xsl ===============/
// 		_path := "." + "/templates/" + templateFolder + "/" + template
// 		ext := filepath.Ext(_path)

// 		// 判断扩展名如果不对就终止
// 		if ext != ".xlsx" {
// 			w.Write([]byte("模板必须是: .xlsx文件"))
// 			return
// 		}

// 		// 生成打印文件
// 		xls := utils.XlsxTemplate{}

// 		resultPath := _cache + "/temporary.xlsx"
// 		xls.PrintOut(_path, resultPath, m)

// 		// --------------------------------------------------------- 转成pdf TODO: 加个判断，看是要pdf还是excel
// 		if printFormat == "pdf" {
// 			fullpath, _ := filepath.Abs(resultPath)

// 			var fireparams []string

// 			fireparams = []string{"/c", "soffice", "--headless", "--invisible", "--convert-to", "pdf", "--outdir", _cache + "/pdf/", fullpath}
// 			interactiveToexec("cmd", fireparams)

// 			// resultPath 针对pdf改成是pdf文件
// 			resultPath = _cache + "/pdf/temporary.pdf"

// 			// http.ServeFile(w, r, "./cache/pdf/test.pdf")
// 			// err := os.RemoveAll("./cache")
// 		}
// 		// --------------------------------------------------------- 转成pdf

// 		file, _ := os.Open(resultPath)
// 		defer file.Close()
// 		buff, _ := ioutil.ReadAll(file)

// 		w.Write(buff)

// 		file.Close()
// 		err = os.RemoveAll(_cache)
// 		if err != nil {
// 			fmt.Println("error from remove cache", err)
// 		}
// 	}
// }
