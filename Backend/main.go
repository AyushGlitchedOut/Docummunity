package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AyushGlitchedOut/Docummunity/dbUtils"
	"github.com/AyushGlitchedOut/Docummunity/routes"
)

func main() {

	//CREATE UPLOADS DIRECTORY
	err := os.MkdirAll("./uploads", 0o755)
	if err != nil {
		log.Println(err)
	}

	//Init DB
	db, err := dbUtils.InitDatabase()
	if err != nil {
		log.Println(err)
	}

	fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")
	port := ":8080"

	server := routes.InitServer(port, db)

	if err := server.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}

	// db, err := sql.Open("sqlite3", "./uploads/data.db")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	// sqlCreate := `CREATE TABLE DATA (
	// ID INTEGER NOT NULL PRIMARY KEY,
	// NAME TEXT
	// )`

	// _, err = db.Exec(sqlCreate)
	// if err != nil {
	// 	log.Println(err)
	// }

	// sqlInsertData := `INSERT INTO DATA VALUES (` + strconv.FormatInt(time.Now().Unix(), 10) + `,'AYUSH GUPTA');`

	// _, err = db.Exec(sqlInsertData)
	// if err != nil {
	// 	log.Println(err)
	// }

	// retrieveCommand := `SELECT ID, NAME FROM DATA;`

	// results, err := db.Query(retrieveCommand)

	// var rows []models.TestData

	// for results.Next() {
	// 	var ID int
	// 	var NAME string

	// 	err := results.Scan(&ID, &NAME)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	rows = append(rows, models.TestData{
	// 		ID:   ID,
	// 		NAME: NAME,
	// 	})
	// }

	// for _, val := range rows {
	// 	fmt.Println(val)
	// }

	// fmt.Println("-----------------------------------------------AFTERWARDS-----------------------")
	// _, err = db.Exec("DELETE FROM DATA")

	// if err != nil {
	// 	log.Println()
	// }
}
