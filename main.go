package main

import (
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/willGuimont/todo/app"
	"github.com/willGuimont/todo/db"
	"log"
	"net/http"
)

func main() {
	database, _ := db.CreateDatabase()
	defer database.Close()

	application := &app.App{
		Router:   mux.NewRouter(),
		Database: database,
	}
	application.SetupRouter()

	log.Fatal(http.ListenAndServe(":8081", application.Router))
}
