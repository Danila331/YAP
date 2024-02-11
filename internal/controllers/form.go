package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"regexp"
	"time"
)

func isValidExpression(input string) bool {
	// Паттерн регулярного выражения для проверки
	pattern := `^[\d+\-*/]+$`

	// Компилируем регулярное выражение
	regex := regexp.MustCompile(pattern)

	// Проверяем соответствие строки паттерну
	return regex.MatchString(input)
}

func FormArifmHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	// Получаем текст из тела запроса
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusInternalServerError)
		return
	}

	expression := r.Form.Get("text")
	if !isValidExpression(expression) {
		fmt.Println("Expression is not correct")
		return
	}

	Task.Expression = expression
	Task.StartDate = time.Now().Format("2006-01-02 15:04:05")
	err = Task.Create()

	if err != nil {
		fmt.Println(err)
	}

	// Отправляем ответ клиенту
	htmlFiles := []string{
		filepath.Join("../", "templates", "submit.html"),
	}

	tmpl, err := template.ParseFiles(htmlFiles...)
	tmpl.ExecuteTemplate(w, "submit", nil)
	if err != nil {
		fmt.Println("Internal server error")
		return
	}
}
