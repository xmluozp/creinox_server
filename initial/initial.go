package initial

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Unknwon/goconfig"
)

func GetMySql() (address string, username string, password string, database string) {

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	address, _ = cfg.GetValue("mysql", "address")
	username, _ = cfg.GetValue("mysql", "username")
	password, _ = cfg.GetValue("mysql", "password")
	database, _ = cfg.GetValue("mysql", "database")

	return address, username, password, database

}

func GetConfig() (rootUrl string, portStr string, port int, uploads string) {

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	cfgIp, err := goconfig.LoadConfigFile("ip.ini")

	if err != nil {
		panic("错误，找不到 ip.ini配置文件")
	}

	rootUrl, _ = cfgIp.GetValue("site", "root")
	portStr, _ = cfg.GetValue("site", "port")
	port, _ = cfg.Int("site", "port")
	uploads, _ = cfg.GetValue("site", "uploads")

	return rootUrl, portStr, port, uploads
}

func getPublicIP() string {

	url := "https://api.ipify.org?format=text" // we are using a pulib IP API, we're using ipify here, below are some others
	// https://www.ipify.org
	// http://myexternalip.com
	// http://api.ident.me
	// http://whatismyipaddress.com/api
	fmt.Printf("Getting IP address from  ipify ...\n")
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(ip)
}

func WriteFile() (clientURL string) {

	var rootUrl string
	var clientUrl string

	f, err := os.Create("ip.ini")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 如果是启动公共版
	var b = flag.Bool("public", false, "bool类型参数")
	flag.Parse()
	fmt.Println("使用公网IP：", *b)

	cfg, err := goconfig.LoadConfigFile("conf.ini")

	if err != nil {
		panic("错误，找不到conf.ini配置文件")
	}

	if *b {
		publicIP := getPublicIP()
		rootUrl = publicIP
		clientUrl = rootUrl
	} else {
		rootUrl, err = cfg.GetValue("site", "root")
		clientUrl, err = cfg.GetValue("site", "client")
	}

	d := []string{
		"[site]",
		"root=" + rootUrl,
		"client=" + clientUrl}

	for _, v := range d {
		fmt.Fprintln(f, v)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return clientUrl

}
