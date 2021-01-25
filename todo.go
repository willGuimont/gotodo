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
	db, _ := db.CreateDatabase()
	defer db.Close()

	app := &app.App{
		Router:   mux.NewRouter(),
		Database: db,
	}
	app.SetupRouter()

	log.Fatal(http.ListenAndServe(":8081", app.Router))
}
