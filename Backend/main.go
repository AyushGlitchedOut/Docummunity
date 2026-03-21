package main

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	err := os.MkdirAll("./uploads", 0o755)

	// fmt.Println("-----------------DOCUMMUNITY BACKEND----------------------------")
	// port := ":8080"

	// router := gin.Default()

	// //TEMPORARY CORS POLICY! REMOVE IN PRODUCTION
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins: []string{"http://localhost:3000"},
	// 	AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
	// 	AllowHeaders: []string{"Origin",
	// 		"Content-Length",
	// 		"Content-Type",
	// 		"Authorization"},
	// }))

	// router.GET("/", routes.HandleGET)
	// //CREATE
	// router.POST("/upload", routes.HandleCREATE)
	// //READ
	// router.GET("/download/:ID", routes.HandleRead)
	// // UPDATE
	// router.PUT("/update/:ID", routes.HandleUpdate)
	// // DELETE
	// router.DELETE("/delete/:ID", routes.HandleDelete)

	// if err := router.Run(port); err != nil {
	// 	log.Fatal("Failed to run server: ", err)
	// }

	db, err := sql.Open("sqlite3", "./uploads/data.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlCreate := `CREATE TABLE DATA (
	ID INTEGER NOT NULL PRIMARY KEY,
	NAME TEXT
	)`

	_, err = db.Exec(sqlCreate)
	if err != nil {
		log.Println(err)
	}

	sqlInsertData := `INSERT INTO DATA VALUES (` + strconv.FormatInt(time.Now().Unix(), 10) + `,'AYUSH GUPTA');`

	_, err = db.Exec(sqlInsertData)
	if err != nil {
		log.Println(err)
	}

}
