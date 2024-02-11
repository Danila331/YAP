package controllers

import (
	"fmt"
	"github/Danila331/YAP/internal/models"
	"html/template"
	"net/http"
	"path/filepath"
)

func FormTimeHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Ошибка при разборе формы", http.StatusInternalServerError)
		return
	}

	timePulse := r.Form.Get("field1")
	timeMinus := r.Form.Get("field2")
	timeProz := r.Form.Get("field3")
	timeDel := r.Form.Get("field4")

	time := models.TimeOperations{
		TimePulse: timePulse,
		TimeMinus: timeMinus,
		TimeProz:  timeProz,
		TimeDel:   timeDel,
	}

	err = time.Update()

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Ошибка при работе с базой данных", http.StatusInternalServerError)
		return
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
