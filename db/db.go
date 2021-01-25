package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

func CreateDatabase() (db *sql.DB, err error) {
	_ = os.Remove("./todos.sqlite")
	db, err = sql.Open("sqlite3", "./todos.sqlite")
	if err != nil {
		log.Fatal(err)
	}

	err = populateDatabase(err, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func populateDatabase(err error, db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE todos (id INTEGER NOT NULL PRIMARY KEY, message TEXT, done BOOLEAN);
	DELETE FROM todos;
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO todos(id, message, done) VALUES (?, ?, FALSE)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, fmt.Sprintf("Do thing %03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}
