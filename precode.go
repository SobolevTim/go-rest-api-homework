package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// обработка метода GET для эндпоинта `/tasks`
// обработчик должен вернуть все задачи, которые хранятся в мапе.
// При успешном запросе сервер должен вернуть статус 200 OK.
// При ошибке сервер должен вернуть статус 500 Internal Server Error.
func getTasks(w http.ResponseWriter, r *http.Request) {
	resp, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		log.Print(err)
	}
}

// обработка метода POST для эндпоинта `/tasks`
// обработчик должен принимать задачу в теле запроса и сохранять ее в мапе.
// При успешном запросе сервер должен вернуть статус 201 Created.
// При ошибке сервер должен вернуть статус 400 Bad Request.
func postTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tasks[task.ID] = task

	w.Header().Set("Content-Type", "applicaion/json")
	w.WriteHeader(http.StatusCreated)
}

// обработка метода GET для эндпоинта `/tasks/{id}`
// обработчик должен вернуть задачу с указанным в запросе пути ID, если такая есть в мапе
// при успешном выполнении запроса сервер должен вернуть статус 200 OK.
// в случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request
func getTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	task, ok := tasks[id]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// обработка метода DELETE для эндпоинта `/tasks/{id}`
// обработчик должен удалить задачу из мапы по её ID
// при успешном выполнении запроса сервер должен вернуть статус 200 OK
// в случае ошибки или отсутствия задачи в мапе сервер должен вернуть статус 400 Bad Request
func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, ok := tasks[id]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	delete(tasks, id)
	_, ok = tasks[id]
	if ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/tasks` с методом GET
	r.Get("/tasks", getTasks)
	// регистрируем в роутере эндпоинт `/tasks` с методом POST
	r.Post("/tasks", postTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом GET
	r.Get("/tasks/{id}", getTask)
	// регистрируем в роутере эндпоинт `/tasks/{id}` с методом DELETE
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
