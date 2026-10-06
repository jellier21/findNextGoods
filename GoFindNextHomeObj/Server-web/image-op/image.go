package image_op

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

//========================《接收前端图片》==========================

// 确保上传目录存在
func ensureUploadDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// ... existing code ...

func uploadHeaderimageHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("uploadHeaderimageHandler")
	// 添加 CORS 响应头，允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		return
	}
	// 1. 限制请求方法
	if r.Method != http.MethodPost {
		http.Error(w, "只支持 POST 请求", http.StatusMethodNotAllowed)
		return
	}

	// 2. 设置最大上传大小（例如 10MB）
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10 MB
	defer r.Body.Close()

	// 3. 解析 multipart 表单
	err := r.ParseMultipartForm(10 << 20) // 内存中最大 10MB
	if err != nil {
		if strings.Contains(err.Error(), "request body too large") {
			http.Error(w, "文件大小不能超过 10MB", http.StatusRequestEntityTooLarge)
		} else {
			http.Error(w, "解析表单失败："+err.Error(), http.StatusBadRequest)
		}
		return
	}

	// 4. 获取要保存的文件名（前端指定）
	targetFilename := strings.TrimSpace(r.FormValue("filename"))

	// 如果没有提供文件名，则返回错误
	if targetFilename == "" {
		http.Error(w, "必须提供文件名参数", http.StatusBadRequest)
		return
	}

	// 5. 安全过滤文件名，防止路径遍历攻击
	targetFilename = filepath.Base(targetFilename) // 去掉路径部分，只保留文件名
	if targetFilename == "." || targetFilename == ".." {
		http.Error(w, "文件名不合法", http.StatusBadRequest)
		return
	}

	// 6. 只允许图片文件扩展名
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
	}
	ext := strings.ToLower(filepath.Ext(targetFilename))
	if !allowedExts[ext] {
		http.Error(w, "只允许上传图片文件 (jpg,jpeg,png,gif,bmp,webp)", http.StatusBadRequest)
		return
	}

	// 7. 获取上传的文件
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "无法读取上传的文件："+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 8. 可选：验证文件确实是图片（读取文件头）
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil && err != io.EOF {
		http.Error(w, "读取文件失败："+err.Error(), http.StatusInternalServerError)
		return
	}

	contentType := http.DetectContentType(buffer)
	if !strings.HasPrefix(contentType, "image/") {
		http.Error(w, "上传的文件不是有效的图片", http.StatusBadRequest)
		return
	}

	// 9. 重置文件读取位置
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, "文件读取失败："+err.Error(), http.StatusInternalServerError)
		return
	}

	// 10. 使用绝对路径创建上传目录（关键修改）
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		http.Error(w, "服务器内部错误：无法获取当前路径", http.StatusInternalServerError)
		return
	}
	currentDir := filepath.Dir(currentFile)
	uploadDir := filepath.Join(currentDir, "..", "images", "header")

	log.Printf("当前文件路径：%s\n", currentFile)
	log.Printf("上传目录路径：%s\n", uploadDir)

	if err := ensureUploadDir(uploadDir); err != nil {
		log.Printf("创建目录失败：%s\n", err.Error())
		http.Error(w, "无法创建上传目录："+err.Error(), http.StatusInternalServerError)
		return
	}

	// 11. 创建/覆盖目标文件
	filePath := filepath.Join(uploadDir, targetFilename)

	dst, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		http.Error(w, "无法创建文件："+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 12. 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "保存文件失败："+err.Error(), http.StatusInternalServerError)
		return
	}

	// 13. 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{
		"success":  true,
		"message":  "文件上传/覆盖成功",
		"filename": targetFilename,
		"size":     header.Size,
		"url":      "/uploads/" + targetFilename,
	}

	json.NewEncoder(w).Encode(response)
}

// ... existing code ...

func uploadNeedImageHandler(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 响应头，允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求（OPTIONS）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. 限制请求方法：只允许 POST
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析 FormData（包含文件和字段）
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		println("表单解析失败:", err.Error())
		http.Error(w, "文件过大或表单解析失败", http.StatusBadRequest)
		return
	}

	// 3. 获取图片名称（从同一个请求的表单中）
	tempName := r.FormValue("data")
	if tempName == "" {
		tempName = fmt.Sprintf("%d_%d", time.Now().UnixNano(), os.Getpid())
	}

	// 4. 从车厢里找到标签为 'image' 的货物（文件）
	file, header, err := r.FormFile("image")
	if err != nil {
		println("读取文件失败:", err.Error())
		http.Error(w, "无法读取上传的文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 5. 获取原文件后缀
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%s%s", tempName, ext)

	// 6. 创建一个存放图片的目录，使用绝对路径
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	uploadDir := filepath.Join(currentDir, "..", "images", "find_need")

	//println("上传目录：" + uploadDir)

	if err := ensureUploadDir(uploadDir); err != nil {
		println("创建目录失败:", err.Error())
		http.Error(w, "服务器无法创建上传目录", http.StatusInternalServerError)
		return
	}

	// 7. 在服务器上创建目标文件
	filePath := filepath.Join(uploadDir, newFilename)
	//println("文件路径：" + filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		println("创建文件失败:", err.Error())
		http.Error(w, "服务器无法创建文件", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 8. 将上传的文件内容拷贝到新文件中（卸货）
	_, err = io.Copy(dst, file)
	if err != nil {
		println("文件保存失败:", err.Error())
		http.Error(w, "文件保存失败", http.StatusInternalServerError)
		return
	}

	// 9. 告诉前端"取件码"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"filePath": "%s"}`, newFilename)
}

func uploadSelfImageHandler(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 响应头，允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求（OPTIONS）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. 限制请求方法：只允许 POST
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析 FormData（包含文件和字段）
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		println("表单解析失败:", err.Error())
		http.Error(w, "文件过大或表单解析失败", http.StatusBadRequest)
		return
	}

	// 3. 获取图片名称（从同一个请求的表单中）
	tempName := r.FormValue("data")
	if tempName == "" {
		tempName = fmt.Sprintf("%d_%d", time.Now().UnixNano(), os.Getpid())
	}

	// 4. 从车厢里找到标签为 'image' 的货物（文件）
	file, header, err := r.FormFile("image")
	if err != nil {
		println("读取文件失败:", err.Error())
		http.Error(w, "无法读取上传的文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 5. 获取原文件后缀
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%s%s", tempName, ext)

	// 6. 创建一个存放图片的目录，使用绝对路径
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	uploadDir := filepath.Join(currentDir, "..", "images", "selfmade")

	//println("上传目录：" + uploadDir)

	if err := ensureUploadDir(uploadDir); err != nil {
		println("创建目录失败:", err.Error())
		http.Error(w, "服务器无法创建上传目录", http.StatusInternalServerError)
		return
	}

	// 7. 在服务器上创建目标文件
	filePath := filepath.Join(uploadDir, newFilename)
	//println("文件路径：" + filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		println("创建文件失败:", err.Error())
		http.Error(w, "服务器无法创建文件", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 8. 将上传的文件内容拷贝到新文件中（卸货）
	_, err = io.Copy(dst, file)
	if err != nil {
		println("文件保存失败:", err.Error())
		http.Error(w, "文件保存失败", http.StatusInternalServerError)
		return
	}

	// 9. 告诉前端"取件码"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"filePath": "%s"}`, newFilename)
}

func uploadGoodsImageHandler(w http.ResponseWriter, r *http.Request) {
	// 添加 CORS 响应头，允许跨域
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// 处理预检请求（OPTIONS）
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// 1. 限制请求方法：只允许 POST
	if r.Method != "POST" {
		http.Error(w, "只支持 POST 方法", http.StatusMethodNotAllowed)
		return
	}

	// 2. 解析 FormData（包含文件和字段）
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		println("表单解析失败:", err.Error())
		http.Error(w, "文件过大或表单解析失败", http.StatusBadRequest)
		return
	}

	// 3. 获取图片名称（从同一个请求的表单中）
	tempName := r.FormValue("data")
	if tempName == "" {
		tempName = fmt.Sprintf("%d_%d", time.Now().UnixNano(), os.Getpid())
	}

	// 4. 从车厢里找到标签为 'image' 的货物（文件）
	file, header, err := r.FormFile("image")
	if err != nil {
		println("读取文件失败:", err.Error())
		http.Error(w, "无法读取上传的文件", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 5. 获取原文件后缀
	ext := filepath.Ext(header.Filename)
	newFilename := fmt.Sprintf("%s%s", tempName, ext)

	// 6. 创建一个存放图片的目录，使用绝对路径
	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	uploadDir := filepath.Join(currentDir, "..", "images", "goods")

	//println("上传目录：" + uploadDir)

	if err := ensureUploadDir(uploadDir); err != nil {
		println("创建目录失败:", err.Error())
		http.Error(w, "服务器无法创建上传目录", http.StatusInternalServerError)
		return
	}

	// 7. 在服务器上创建目标文件
	filePath := filepath.Join(uploadDir, newFilename)
	//println("文件路径：" + filePath)

	dst, err := os.Create(filePath)
	if err != nil {
		println("创建文件失败:", err.Error())
		http.Error(w, "服务器无法创建文件", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// 8. 将上传的文件内容拷贝到新文件中（卸货）
	_, err = io.Copy(dst, file)
	if err != nil {
		println("文件保存失败:", err.Error())
		http.Error(w, "文件保存失败", http.StatusInternalServerError)
		return
	}

	// 9. 告诉前端"取件码"
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"filePath": "%s"}`, newFilename)
}

// ===========================《返回相应图片》===========================================
func headerimageHandler(w http.ResponseWriter, r *http.Request) {
	//检查本函数api是否被调用
	//println("头像图片")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 1. 从 URL 查询参数中获取"取件码"（文件名）
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "缺少文件名参数", http.StatusBadRequest)
		return
	}

	cleanFileName := filepath.Base(fileName) // Base() 会去掉路径部分，只留文件名
	if cleanFileName != fileName {
		http.Error(w, "无效的文件名", http.StatusBadRequest)
		return
	}

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	filePath := filepath.Join(currentDir, "..", "images", "header", cleanFileName)

	// 4. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		println("文件不存在：" + filePath)
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	//println("文件路径：" + filePath)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 6. 禁用缓存（开发时方便，生产环境可调整）
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// 7. 将文件内容直接写入响应体
	http.ServeFile(w, r, filePath)
}

func NeedimageHandler(w http.ResponseWriter, r *http.Request) {
	//检查本函数api是否被调用
	//println("悬赏图片")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 1. 从 URL 查询参数中获取"取件码"（文件名）
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "缺少文件名参数", http.StatusBadRequest)
		return
	}

	cleanFileName := filepath.Base(fileName) // Base() 会去掉路径部分，只留文件名
	if cleanFileName != fileName {
		http.Error(w, "无效的文件名", http.StatusBadRequest)
		return
	}

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	filePath := filepath.Join(currentDir, "..", "images", "find_need", cleanFileName)

	// 4. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		println("文件不存在：" + filePath)
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	//println("文件路径：" + filePath)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 6. 禁用缓存（开发时方便，生产环境可调整）
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// 7. 将文件内容直接写入响应体
	http.ServeFile(w, r, filePath)
}

func SelfimageHandler(w http.ResponseWriter, r *http.Request) {
	//检查本函数api是否被调用
	//println("悬赏图片")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 1. 从 URL 查询参数中获取"取件码"（文件名）
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "缺少文件名参数", http.StatusBadRequest)
		return
	}

	cleanFileName := filepath.Base(fileName) // Base() 会去掉路径部分，只留文件名
	if cleanFileName != fileName {
		http.Error(w, "无效的文件名", http.StatusBadRequest)
		return
	}

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	filePath := filepath.Join(currentDir, "..", "images", "selfmade", cleanFileName)

	// 4. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		println("文件不存在：" + filePath)
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	//println("文件路径：" + filePath)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 6. 禁用缓存（开发时方便，生产环境可调整）
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// 7. 将文件内容直接写入响应体
	http.ServeFile(w, r, filePath)
}

func GoodsimageHandler(w http.ResponseWriter, r *http.Request) {
	//检查本函数api是否被调用
	//println("悬赏图片")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// 1. 从 URL 查询参数中获取"取件码"（文件名）
	fileName := r.URL.Query().Get("filename")
	if fileName == "" {
		http.Error(w, "缺少文件名参数", http.StatusBadRequest)
		return
	}

	cleanFileName := filepath.Base(fileName) // Base() 会去掉路径部分，只留文件名
	if cleanFileName != fileName {
		http.Error(w, "无效的文件名", http.StatusBadRequest)
		return
	}

	_, currentFile, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(currentFile)
	filePath := filepath.Join(currentDir, "..", "images", "goods", cleanFileName)

	// 4. 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		println("文件不存在：" + filePath)
		http.Error(w, "图片不存在", http.StatusNotFound)
		return
	}

	//println("文件路径：" + filePath)

	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".gif":
		w.Header().Set("Content-Type", "image/gif")
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
	}

	// 6. 禁用缓存（开发时方便，生产环境可调整）
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	// 7. 将文件内容直接写入响应体
	http.ServeFile(w, r, filePath)
}

// =======================Use The Image File Api Server
func UseTheImageFileApiServer() {
	http.HandleFunc("/api/image/header/upload", uploadHeaderimageHandler)
	http.HandleFunc("/api/image/need/upload", uploadNeedImageHandler)
	http.HandleFunc("/api/image/self/upload", uploadSelfImageHandler)
	http.HandleFunc("/api/image/goods/upload", uploadGoodsImageHandler)
	http.HandleFunc("/api/image/header/show", headerimageHandler)
	http.HandleFunc("/api/image/need/show", NeedimageHandler)
	http.HandleFunc("/api/image/self/show", SelfimageHandler)
	http.HandleFunc("/api/image/goods/show", GoodsimageHandler)
}
