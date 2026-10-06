package html_template

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"
)

// 获取 frontend 目录的绝对路径
func getFrontendPath() string {
	// 获取当前文件所在的目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	// 从 template.go 向上三级到 Server-web 目录，然后进入 frontend 目录
	// template.go → Server-Side → Server-web → GoFindNextHomeObj → frontend
	currentDir := filepath.Dir(filename)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(currentDir)))
	frontendPath := filepath.Join(projectRoot, "frontend")

	return frontendPath
}

var frontendPath = getFrontendPath()

// ... existing code ...

func renderTemplate(w http.ResponseWriter, templateName string) {
	templatePath := filepath.Join(frontendPath, templateName)
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		fmt.Fprintf(w, "template.ParseFiles err: %v\n", err)
		fmt.Fprintf(w, "Trying to load from: %s\n", templatePath)
		return
	}
	err = t.Execute(w, nil)
	if err != nil {
		fmt.Fprintf(w, "template.Execute err: %v", err)
		return
	}
}

func HomePageTemplateApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "homepage.html")
}

func FrontendRootApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html")
}

func ExhibitionTemplateApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "exhibition.html")
}

func FindNeedsTemplateApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "find_Needs.html")
}

func MyInterfaceTemplateApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "myInterface.html")
}

func SelfMadeTemplateApi(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "Self-made.html")
}

func UseTheTemplateApiServer() {
	http.HandleFunc("/frontend", FrontendRootApi)
	http.HandleFunc("/frontend/", FrontendRootApi)
	http.HandleFunc("/frontend/homepage", HomePageTemplateApi)
	http.HandleFunc("/frontend/exhibition", ExhibitionTemplateApi)
	http.HandleFunc("/frontend/find-needs", FindNeedsTemplateApi)
	http.HandleFunc("/frontend/my-interface", MyInterfaceTemplateApi)
	http.HandleFunc("/frontend/self-made", SelfMadeTemplateApi)

	// 静态文件服务
	http.Handle("/frontend/css/", http.StripPrefix("/frontend/css/", http.FileServer(http.Dir(filepath.Join(frontendPath, "css")))))
	http.Handle("/frontend/images/", http.StripPrefix("/frontend/images/", http.FileServer(http.Dir(filepath.Join(frontendPath, "images")))))
	http.Handle("/frontend/js/", http.StripPrefix("/frontend/js/", http.FileServer(http.Dir(filepath.Join(frontendPath, "js")))))
}
