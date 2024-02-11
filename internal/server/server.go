package server

import (
	"github/Danila331/YAP/internal/controllers"
	"net/http"
)

// 1. Написать форму для отправки времени для каждой операции
// 2. Сделать отдельную базу для этого куда будем это сохранять
// 3. Так же расширть бд tasks добавить туда id agent где выполянются вырадения
// 4. Сделать бд куда регестируются сервера и к ним записываются их выражения
// придумать потом как сделать и чекать отвалился он или нет

// Функция для запуска сервера
func StartServer() error {

	http.HandleFunc("/", controllers.MainPageHandler)
	http.HandleFunc("/submit", controllers.FormArifmHandler)
	http.HandleFunc("/expression", controllers.ExpressionHandler)
	http.HandleFunc("/expressions", controllers.ExpressionsHandler)
	http.HandleFunc("/submit-time", controllers.FormTimeHandler)
	http.HandleFunc("/form-time", controllers.FormTimePageHanlder)
	// http.Handle("../static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		return err
	}
	return nil
}
