package dbUtils

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./uploads/data.db")
	if err != nil {
		log.Println(err)
	}

	createCommand := `CREATE TABLE DATA (ID INTEGER PRIMARY KEY, NAME TEXT NOT NULL, FILEPATH TEXT NOT NULL);`

	_, err = db.Exec(createCommand)

	return db, err
}

func InsertData(data TestData, db *sql.DB) error {

	createCommand := `INSERT INTO DATA (ID, NAME, FILEPATH) VALUES (?, ?, ?)`

	_, err := db.Exec(createCommand, data.ID, data.NAME, data.FILEPATH)

	return err
}
