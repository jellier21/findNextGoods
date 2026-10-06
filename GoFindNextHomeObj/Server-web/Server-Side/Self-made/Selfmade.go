package Self_made

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// ========================<结构体>==========================
type MadeDB struct {
	Id           int     `json:"id"`
	Account      string  `json:"account"`
	WorkName     string  `json:"name"`
	Introduction string  `json:"introduction"`
	Price        float32 `json:"price"`
	Type         int     `json:"type"`
	Image        string  `json:"image"`
}

type MadeData struct {
	Account      string  `json:"account"`
	WorkName     string  `json:"name"`
	Introduction string  `json:"introduction"`
	Price        float32 `json:"price"`
	Type         int     `json:"type"`
	Image        string  `json:"image"`
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

func FindSelfIssue(account_a string) (backDB []MadeDB, err error) {
	sqlStr := "select id, account, work_name, introduce, price, type, image from self_made where account = ?"
	rows, err := db.Query(sqlStr, account_a)
	if err != nil {
		return backDB, err
	}
	defer rows.Close()
	for rows.Next() {
		var made MadeDB
		err = rows.Scan(&made.Id, &made.Account, &made.WorkName, &made.Introduction, &made.Price, &made.Type, &made.Image)
		if err != nil {
			return backDB, err
		}
		backDB = append(backDB, MadeDB{
			Id:           made.Id,
			Account:      made.Account,
			WorkName:     made.WorkName,
			Introduction: made.Introduction,
			Price:        made.Price,
			Type:         made.Type,
			Image:        made.Image,
		})
	}
	return backDB, nil

}

func FindAllSelfIssue() (backDB []MadeDB, err error) {
	sqlStr := "select id, account, work_name, introduce, price, type, image from self_made where type = ?"
	rows, err := db.Query(sqlStr, 1)
	if err != nil {
		return backDB, err
	}
	defer rows.Close()
	for rows.Next() {
		var made MadeDB
		err = rows.Scan(&made.Id, &made.Account, &made.WorkName, &made.Introduction, &made.Price, &made.Type, &made.Image)
		if err != nil {
			return backDB, err
		}
		backDB = append(backDB, made)
	}
	return backDB, nil
}

func AddSelfIssue(account_a string, work_name_a string, introduce_a string, price_a float32, type_a int, image_a string) (success bool, err error) {
	ret, err := db.Exec("insert into self_made(account, work_name, introduce, price, type, image) values (?, ?, ?, ?, ?, ?)",
		account_a, work_name_a, introduce_a, price_a, type_a, image_a)
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

func DeleteSelfIssue(Id_s int) (success bool, err error) {
	ret, err := db.Exec("delete from self_made where id = ? ", Id_s)
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

// ========================<实现api服务器功能>=========================
func FindAndBackSelfInformationApi(w http.ResponseWriter, r *http.Request) {
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

	WorkData, err := FindSelfIssue(req.Account)
	//for _, WorkData := range WorkData {
	//	fmt.Println(WorkData.Account, WorkData.WorkName, WorkData.Introduction)
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

func FindAndBackAllInformationApi(w http.ResponseWriter, r *http.Request) {
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

	WorkData, err := FindAllSelfIssue()
	if err != nil {
		http.Error(w, "该单元项目不存在", http.StatusNotFound)
		return
	}
	//for _, WorkData := range WorkData {
	//	fmt.Println(WorkData.Account, WorkData.WorkName, WorkData.Introduction)
	//}

	jsonData, err := json.Marshal(WorkData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func AddSelfIssueApi(w http.ResponseWriter, r *http.Request) {
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

	var req MadeData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "JSON 解析失败："+err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("收到发布自制品请求：", req.Account, req.WorkName, req.Introduction)

	WorkData, err := AddSelfIssue(req.Account, req.WorkName, req.Introduction, req.Price, req.Type, req.Image)

	if err != nil {
		http.Error(w, "发布失败", http.StatusNotFound)
		return
	}

	type backBool struct {
		Account string `json:"b_account"`
		Success bool   `json:"success"`
	}

	var resBool backBool
	resBool.Account = req.Account
	resBool.Success = WorkData

	//log.Println("成功发布作品", resBool.Success)

	jsonData, err := json.Marshal(resBool)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	//log.Println("成功发布作品")

}

func DeleteSelfIssueApi(w http.ResponseWriter, r *http.Request) {
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

	log.Println("收到删除发布请求：", req.Id)

	WorkData, err := DeleteSelfIssue(req.Id)

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

// =========================<启动api服务器>=========================
func UseSelfMadeApiServer() {
	err := InitDB()
	if err != nil {
		log.Fatal("数据库初始化失败：", err)
	}
	fmt.Println("自制作品数据库连接成功")
	http.HandleFunc("/api/selfmade/send", AddSelfIssueApi)
	http.HandleFunc("/api/selfmade/showmyself", FindAndBackSelfInformationApi)
	http.HandleFunc("/api/selfmade/showallself", FindAndBackAllInformationApi)
	http.HandleFunc("/api/selfmade/delete", DeleteSelfIssueApi)
}
