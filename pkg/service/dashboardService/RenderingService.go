package dashboardService

import (
	"log"
	"net/http"
	"text/template"
)

type RenderingService struct {
}

// NewRenderingService - constructor of RenderingService structure.
func NewRenderingService() *RenderingService {
	return &RenderingService{}
}

// RenderFromFiles - render templates from files and pass the `date` into it.
// The order of the files in `templates` is important, daughters first, then parents.
func (service *RenderingService) RenderFromFiles(w http.ResponseWriter, templates []string, data interface{}) {
	// Используем функцию template.ParseFiles() для чтения файлов шаблона.
	// Если возникла ошибка, мы запишем детальное сообщение ошибки и
	// используя функцию http.Error() мы отправим пользователю
	// ответ: 500 Internal Server Error (Внутренняя ошибка на сервере)
	ts, err := template.ParseFiles(templates...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Затем мы используем метод Execute() для записи содержимого
	// шаблона в тело HTTP ответа. Последний параметр в Execute() предоставляет
	// возможность отправки динамических данных в шаблон.
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
