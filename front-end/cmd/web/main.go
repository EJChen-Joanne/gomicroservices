package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Front-end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, html_root string) {
	html_partials := []string{
		"base.layout.gohtml",
		"footer.partial.gohtml",
		"header.partial.gohtml",
	}
	templatesSli := make([]string, 0)
	templatesSli = append(templatesSli, fmt.Sprintf("templates/%s", html_root)) //format templates

	for _, v := range html_partials {
		templatesSli = append(templatesSli, fmt.Sprintf("templates/%s", v))
	}

	tpl, err := template.ParseFS(templateFS, templatesSli...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
