package Images

import (
	"net/http"
	"os"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// 打开图片文件
	file, err := os.Open("../../images/261001.jpg")
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	// 设置响应头（关键步骤）
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Cache-Control", "max-age=3600") // 缓存优化

	// 将文件流直接写入响应
	//http.ServeContent(w, r, "image.jpg", file.ModTime(), file)
}

func UseTheImagesServer() {
	http.HandleFunc("/api/image/hander", imageHandler)
}
