package User

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

// -----------------------------------------------------连接数据库
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

// ========================<结构体>==========================
// 从数据库打捞数据至服务端
type UserDB struct {
	Id       int
	Nickname string
	Account  string
	Password string
}

// 自服务端向客户端传递
type UserData struct {
	Nickname string
	Account  string
	Passkey  string
}

// =========================<数据库操作工具>==========================
func QueryOneUser(account string) (userDB UserDB, err error) {
	sqlStr := "select id, nickname, account, password from accounts where account = ?"
	err = db.QueryRow(sqlStr, account).Scan(&userDB.Id, &userDB.Nickname, &userDB.Account, &userDB.Password)
	return userDB, err
}

func AddTheUserDBTools(nick_a string, account_a string, password_a string) error {
	ret, err := db.Exec("insert into accounts(nickname,account,password) values(?,?,?)",
		nick_a, account_a, password_a)
	if err != nil {
		fmt.Println("Exec failed, err:", err)
		return err
	}
	Uid, err := ret.LastInsertId()
	if err != nil {
		fmt.Println("GetLastInsertId failed, err:", err)
		return err
	}
	fmt.Println("Insert success, id:", Uid)
	return nil
}

func WhetherUserDBRepetition(account_a string) bool {
	// 判断账号是否重复
	var err error
	_, err = QueryOneUser(account_a)
	if err != nil {
		return false
	} else {
		return true
	}
}

func FindAndReturnNick(account_a string) string {
	var err error
	var userDB UserDB
	userDB, err = QueryOneUser(account_a)
	if err != nil {
		return ""
	}
	return userDB.Nickname
}

// ========================<实现api服务器功能>=========================
func LoginTheUserAPI(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "数据库连接失败", http.StatusInternalServerError)
		return
	}
	// 设置允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求
	if r.Method == "OPTIONS" {
		return
	}

	// 只处理 POST 请求
	if r.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 检查 Content-Type 是否为 JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type 必须为 application/json", http.StatusBadRequest)
		return
	}

	// 解析 JSON 数据
	type Login struct {
		AccountDB string `json:"accDB"`
		Passkey   string `json:"passkey"`
	}
	var receiveLogin Login
	err := json.NewDecoder(r.Body).Decode(&receiveLogin)
	if err != nil {
		http.Error(w, "JSON 解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("\n接收到登录请求：\n\t账号=%s,\n\t密码=%s\n", receiveLogin.AccountDB, receiveLogin.Passkey)

	// 查询用户
	var userRowData UserDB
	userRowData, err = QueryOneUser(receiveLogin.AccountDB)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	// 验证密码
	if userRowData.Password != receiveLogin.Passkey {
		http.Error(w, "账号或密码错误", http.StatusBadRequest)
		return
	}

	var reData UserData
	reData.Account = userRowData.Account
	reData.Nickname = userRowData.Nickname
	reData.Passkey = userRowData.Password

	// 返回成功响应
	jsonData, err := json.Marshal(reData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	log.Printf("登录成功")

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func AddTheUserDBAPI(w http.ResponseWriter, r *http.Request) {

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

	// 检查 Content-Type 是否为 JSON
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Content-Type 必须为 application/json", http.StatusBadRequest)
		return
	}

	// 解析 JSON 数据
	var u UserData
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "JSON 解析失败: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("接收到数据：%+v", u)
	if WhetherUserDBRepetition(u.Account) {
		http.Error(w, "用户已存在", http.StatusBadRequest)
		return
	}
	var errDB error
	errDB = AddTheUserDBTools(u.Nickname, u.Account, u.Passkey)
	if errDB != nil {
		http.Error(w, "数据库操作失败: "+errDB.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "注册成功")
}

func FindAndReturnNickAPI(w http.ResponseWriter, r *http.Request) {
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

	userDB, err := QueryOneUser(req.Account)
	if err != nil {
		http.Error(w, "用户不存在", http.StatusNotFound)
		return
	}

	responseData := map[string]string{
		"account":  userDB.Account,
		"nick":     userDB.Nickname,
		"nickname": userDB.Nickname,
	}

	jsonData, err := json.Marshal(responseData)
	if err != nil {
		http.Error(w, "Failed to marshal data", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	//log.Printf("返回用户信息成功：Account=%s, Nickname=%s", userDB.Account, userDB.Nickname)
}

// ===================<汇总>==============================
func UseTheUserServer() {
	if err := InitDB(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	fmt.Println("个人信息数据库连接成功")
	http.HandleFunc("/api/user/login", LoginTheUserAPI)
	http.HandleFunc("/api/user/add", AddTheUserDBAPI)
	http.HandleFunc("/api/user/account-nick", FindAndReturnNickAPI)

}
