package UserNeed

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ========================<结构体>==========================
type NeedDB struct {
	Id          int    `json:"id"`
	Account     string `json:"account"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type NeedData struct {
	Account     string `json:"account"`
	Description string `json:"description"`
	Image       string `json:"image"`
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
func AddmyNeedSQL(account_s string, Description_s string, Image_s string) (success bool, err error) {
	sqlStr := "insert into need(account,description,image) values(?,?,?)"
	ret, err := db.Exec(sqlStr, account_s, Description_s, Image_s)
	if err != nil {
		log.Println("添加失败")
		return false, err
	}
	_, err = ret.LastInsertId()

	if err != nil {
		return false, err
	}
	return true, err
}

func BackAllNeedList() (backDB []NeedDB, err error) {
	sqlStr := "select id, account, description,image from need "
	rows, err := db.Query(sqlStr)
	if err != nil {
		return backDB, err
	}
	defer rows.Close()
	for rows.Next() {
		var n NeedDB
		err = rows.Scan(&n.Id, &n.Account, &n.Description, &n.Image)
		if err != nil {
			return backDB, err
		}
		backDB = append(backDB, n)
	}
	return backDB, nil
}

func FindNeedByAccount(account_a string) (backDB []NeedDB, err error) {
	sqlStr := "select id, account, description,image from need where account = ? "
	rows, err := db.Query(sqlStr, account_a)
	if err != nil {
		return backDB, err
	}
	defer rows.Close()
	for rows.Next() {
		var n NeedDB
		err = rows.Scan(&n.Id, &n.Account, &n.Description, &n.Image)
		if err != nil {
			return backDB, err
		}
		backDB = append(backDB, n)
	}
	return backDB, nil

}

func deleteSQLById(ID int) (success bool, err error) {
	ret, err := db.Exec("delete from need where id = ? ", ID)
	if err != nil {
		log.Println("删除失败")
		return false, err
	}
	_, err = ret.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, err

}

// =======================< api接口 >==========================
func AddmyNeedApi(w http.ResponseWriter, r *http.Request) {
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

	var req NeedData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("收到发布悬赏请求：", req.Account, req.Description)

	WorkData, err := AddmyNeedSQL(req.Account, req.Description, req.Image)

	if err != nil {
		http.Error(w, "发布失败", http.StatusNotFound)
		return
	}

	type backBool struct {
		Account string `json:"account"`
		Success bool   `json:"success"`
	}

	var resBool backBool
	resBool.Account = req.Account
	resBool.Success = WorkData

	jsonData, err := json.Marshal(resBool)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func BackAllNeedListApi(w http.ResponseWriter, r *http.Request) {
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

	WorkData, err := BackAllNeedList()
	if err != nil {
		http.Error(w, "该单元项目不存在", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(WorkData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func FindMyNeedListApi(w http.ResponseWriter, r *http.Request) {
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

	type requestDB struct {
		Account string `json:"Account"`
	}
	var req requestDB
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	WorkData, err := FindNeedByAccount(req.Account)
	//for _, WorkData := range WorkData {
	//	fmt.Println(WorkData.Account, WorkData.Description)
	//}
	if err != nil {
		http.Error(w, "该单元项目不存在", http.StatusNotFound)
		return
	}

	jsonData, err := json.Marshal(WorkData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func DeleteNEEDIssueApi(w http.ResponseWriter, r *http.Request) {
	//log.Println("删除发布")
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

	type receiveDB struct {
		Id int `json:"id"`
	}
	var req receiveDB
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	//log.Println("收到删除悬赏请求：", req.Id)

	WorkData, err := deleteSQLById(req.Id)

	if err != nil {
		http.Error(w, "删除失败", http.StatusNotFound)
		return
	}

	type backBool struct {
		Success bool `json:"success"`
	}

	var resBool backBool
	resBool.Success = WorkData

	jsonData, err := json.Marshal(resBool)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// =======================< 启用服务器 >==========================
func UseNeedApiServer() {
	err := InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败：", err)
	}
	fmt.Println("悬赏数据库连接成功")
	http.HandleFunc("/api/need/add", AddmyNeedApi)
	http.HandleFunc("/api/need/showall", BackAllNeedListApi)
	http.HandleFunc("/api/need/showmy", FindMyNeedListApi)
	http.HandleFunc("/api/need/delete", DeleteNEEDIssueApi)
}
