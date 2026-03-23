package dbUtils

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitDatabase() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./uploads/data.db")
	if err != nil {
		log.Println(err)
	}

	createCommand := `CREATE TABLE DATA (TIMEID INTEGER PRIMARY KEY, NAME TEXT NOT NULL, FILEPATH TEXT NOT NULL);`

	_, err = db.Exec(createCommand)

	return db, err
}

func InsertData(data TestData, db *sql.DB) error {

	createCommand := `INSERT INTO DATA (TIMEID, NAME, FILEPATH) VALUES (?, ?, ?)`

	_, err := db.Exec(createCommand, data.TIMEID, data.NAME, data.FILEPATH)

	return err
}

func ReadData(NAME string, db *sql.DB) ([]TestData, error) {

	var result []TestData

	readCommand := `SELECT TIMEID,NAME,FILEPATH FROM DATA WHERE NAME = ?`

	rows, err := db.Query(readCommand, NAME)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var row TestData

		err := rows.Scan(&row.TIMEID, &row.NAME, &row.FILEPATH)

		if err != nil {
			return result, err
		}
		result = append(result, row)
	}

	return result, nil
}

func DeleteData(NAME string, db *sql.DB) error {
	//Delete File

	var result []string

	readCommand := `SELECT FILEPATH FROM DATA WHERE NAME = ?`

	rows, err := db.Query(readCommand, NAME)
	if err != nil {
		return err
	}

	for rows.Next() {
		var row string

		err := rows.Scan(&row)

		if err != nil {
			return err
		}
		result = append(result, row)

	}

	defer rows.Close()

	if len(result) == 0 {
		return fmt.Errorf("No Records Found")
	}

	for _, val := range result {
		err := os.Remove(val)

		if err != nil {
			return err
		}
	}

	//Delete Record
	deleteCommand := `DELETE FROM DATA WHERE NAME = ?;`

	_, err = db.Exec(deleteCommand, NAME)

	if err != nil {
		return err
	}

	return nil
}
