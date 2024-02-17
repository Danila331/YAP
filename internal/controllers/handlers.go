package controllers

import (
	"fmt"
	"github/Danila331/YAP/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
)

var Task models.Task

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	htmlFiles := []string{
		filepath.Join("../", "templates", "main.html"),
	}

	tmpl, err := template.ParseFiles(htmlFiles...)
	tmpl.ExecuteTemplate(w, "main", nil)
	if err != nil {
		fmt.Println("Internal server error")
		return
	}
}

// Handler для отображения страницы с выражениями
func ExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	htmlFiles := []string{
		filepath.Join("../", "templates", "expressions.html"),
	}

	data, err := Task.ReadAll()

	if err != nil {
		fmt.Println(err)
		return
	}

	tmpl, err := template.ParseFiles(htmlFiles...)
	tmpl.ExecuteTemplate(w, "expressions", data)
	if err != nil {
		fmt.Println("Internal server error")
		return
	}
}

// Handler для отображения страницы с определнным выражением
func ExpressionHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	Task.Id = id
	data, err := Task.ReadById()

	if err != nil {
		fmt.Println(err)
		return
	}

	htmlFiles := []string{
		filepath.Join("../", "templates", "expression.html"),
	}

	tmpl, err := template.ParseFiles(htmlFiles...)
	tmpl.ExecuteTemplate(w, "expression", data)

	if err != nil {
		fmt.Println("Internal server error")
		return
	}
}

func FormTimePageHanlder(w http.ResponseWriter, r *http.Request) {
	htmlFiles := []string{
		filepath.Join("../", "templates", "form-time.html"),
	}

	tmpl, err := template.ParseFiles(htmlFiles...)
	tmpl.ExecuteTemplate(w, "form-time", nil)

	if err != nil {
		fmt.Println("Internal server error")
		return
	}
}
