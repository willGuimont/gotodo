package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router   *mux.Router
	Database *sql.DB
}

func (app *App) SetupRouter() {
	app.Router.HandleFunc("/todos", app.GetTodosHandler).Methods("GET")
	app.Router.HandleFunc("/todo/{id}", app.GetTodoHandler).Methods("GET")
	app.Router.HandleFunc("/todo", app.CreateTodoHandler).Methods("POST")
	app.Router.HandleFunc("/todo/{id}", app.MarkTodoDoneHandler).Methods("PUT")
}

type Todo struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
	Done    bool   `json:"done"`
}

func (app *App) GetTodoHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	todo := app.getTodo(id)
	_ = json.NewEncoder(w).Encode(todo)
}

func (app *App) getTodo(id string) (*Todo) {
	result, err := app.Database.Query("SELECT id, message, done FROM todos")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var todo Todo
		err := result.Scan(&todo.ID, &todo.Message, &todo.Done)
		if err != nil {
			panic(err.Error())
		}
		if todo.ID == id {
			return &todo
		}
	}
	return nil
}

func (app *App) GetTodosHandler(w http.ResponseWriter, req *http.Request) {
	var todos []Todo
	result, err := app.Database.Query("SELECT id, message, done FROM todos")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var todo Todo
		err = result.Scan(&todo.ID, &todo.Message, &todo.Done)
		todos = append(todos, todo)
	}

	_ = json.NewEncoder(w).Encode(todos)
}

func (app *App) CreateTodoHandler(w http.ResponseWriter, req *http.Request) {
	var todo Todo
	_ = json.NewDecoder(req.Body).Decode(&todo)
	res, err := app.Database.Exec("INSERT INTO todos(message, done) VALUES (?, FALSE)", todo.Message)

	if err != nil {
		log.Fatal(err)
	} else {
		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err)
		}
		todo.ID = fmt.Sprintf("%d", id)
	}

	_ = json.NewEncoder(w).Encode(todo)
}

func (app *App) MarkTodoDoneHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	_, _ = app.Database.Exec("UPDATE todos SET done = true WHERE id = ?", id)
	todo :=	app.getTodo(id)
	_ = json.NewEncoder(w).Encode(todo)
}
