package User_Location

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ========================<结构体>==========================
type SiteDB struct {
	id      int
	Account string
	RelWay  string
	Site    string
}

type backSite struct {
	Account  string `json:"Account"`
	SiteWay  string `json:"SiteWay"`
	Location string `json:"Location"`
}

// 连接数据库
var db *sql.DB

func InitDB() (err error) {
	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/find-nexthome")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

// =========================<数据库操作工具>==========================
func QueryOneSiteDB(account_s string) (site SiteDB, err error) {
	sqlStr := "select id, account, rel_way, location from address where account = ?"
	err = db.QueryRow(sqlStr, account_s).Scan(&site.id, &site.Account, &site.RelWay, &site.Site)
	return site, err
}

func UploadOneSiteDB(account_s string, rel_way_s string, site_s string) (err error) {
	_, err = QueryOneSiteDB(account_s)
	if err != nil {
		sqlStr := "insert into address (account, rel_way, location) values (?, ?, ?)"
		_, err = db.Exec(sqlStr, account_s, rel_way_s, site_s)
		return err
	} else {
		sqlStr := "update address set rel_way = ?, location = ? where account = ?"
		_, err = db.Exec(sqlStr, rel_way_s, site_s, account_s)
		return err
	}

}

// ========================<实现api服务器功能>=========================
func BackOneSiteApi(w http.ResponseWriter, r *http.Request) {
	//log.Printf("收到请求：%s %s", r.Method, r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type 必须为 application/json", http.StatusBadRequest)
		return
	}

	type AccountRequest struct {
		Account string `json:"account"`
	}

	var req AccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	//log.Println("收到查询联系方式请求：Account=%s", req.Account)

	SiteDB, err := QueryOneSiteDB(req.Account)
	if err != nil {
		http.Error(w, "该单元项目不存在", http.StatusNotFound)
		return
	}

	var responseData backSite

	responseData.Account = SiteDB.Account
	responseData.SiteWay = SiteDB.RelWay
	responseData.Location = SiteDB.Site

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	//log.Printf("返回用户联系成功：Account=%s, FindWay=%s, Location=%s", SiteDB.Account, SiteDB.RelWay, SiteDB.Site)
}

func UploadOneSiteApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type 必须为 application/json", http.StatusBadRequest)
		return
	}

	type AccountRequest struct {
		Account string `json:"account"`
		RelWay  string `json:"rel_way"`
		Site    string `json:"site"`
	}

	type backBool struct {
		Bool bool `json:"bool"`
	}
	var responseData backBool

	var req AccountRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("收到修改联系方式请求：Account=%s", req.Account)

	err = UploadOneSiteDB(req.Account, req.RelWay, req.Site)
	if err != nil {
		http.Error(w, "修改失败", http.StatusNotFound)
		return
	}

	responseData.Bool = true

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	//log.Printf("返回用户联系成功：Account=%s, FindWay=%s, Location", SiteDB.Account, SiteDB.RelWay, SiteDB.Site)
}

func UseTheSiteServer() {
	err := InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败：", err)
	}
	fmt.Println("联系方式数据库连接成功")
	http.HandleFunc("/api/site/mysite", BackOneSiteApi)
	http.HandleFunc("/api/site/upload", UploadOneSiteApi)
}
