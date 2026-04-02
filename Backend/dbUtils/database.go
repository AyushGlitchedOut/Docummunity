package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func InitializeDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./uploads/DATABASE.db")
	if err != nil {
		return nil, err
	}

	//TODO: switch UID to a more efficient type
	dataCreateCommand := `CREATE TABLE IF NOT EXISTS DATA (
	TIME_UUID TEXT PRIMARY KEY, 
	NAME TEXT NOT NULL,
	FILEPATH TEXT NOT NULL,
	CREATOR_ID TEXT NOT NULL,
	PREVIEW_IMG_PATH TEXT
	);`
	_, err = db.ExecContext(ctx, dataCreateCommand)
	if err != nil {
		return nil, err
	}

	//TODO: switch Time and Settings to a more efficient type
	userCreateCommand := `CREATE TABLE IF NOT EXISTS USERS (
		UID TEXT PRIMARY KEY, 
		EMAIL TEXT,
		DISPLAY_NAME TEXT,
		PROFILE_PIC TEXT,
		CREATION_DATE TEXT NOT NULL,
		SETTINGS TEXT
	);`

	_, err = db.ExecContext(ctx, userCreateCommand)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(ctx context.Context, user *USER, db *sql.DB) error {
	userInsertCommand := `INSERT INTO USERS (UID, EMAIL, DISPLAY_NAME, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`
	if user.UID == "" {
		return fmt.Errorf("NO UID Provided")
	}
	if user.CREATION_DATE == "" {
		return fmt.Errorf("NO Creation Date Provided")
	}
	if user.EMAIL == "" && user.DISPLAY_NAME == "" {
		return fmt.Errorf("Bad User")
	}

	_, err := db.ExecContext(ctx, userInsertCommand, user.UID, user.EMAIL, user.DISPLAY_NAME, user.PROFILE_PIC, user.CREATION_DATE, user.SETTINGS)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUser(ctx context.Context, UID string, db *sql.DB) error {

	deleteUserCommand := `DELETE FROM USERS WHERE UID = ?`

	results, err := db.ExecContext(ctx, deleteUserCommand, UID)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No User found with UID %s", UID)
	}

	return nil
}

func GetUserInfo(ctx context.Context, UID string, db *sql.DB) (*USER, error) {
	user := &USER{}
	getUserCommand := `SELECT UID, EMAIL,DISPLAY_NAME, PROFILE_PIC, CREATION_DATE, SETTINGS FROM USERS WHERE UID = ?`
	result := db.QueryRowContext(ctx, getUserCommand, UID)
	err := result.Scan(&user.UID, &user.EMAIL, &user.DISPLAY_NAME, &user.PROFILE_PIC, &user.CREATION_DATE, &user.SETTINGS)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No Rows Found")
		}
		return nil, err
	}

	return user, nil
}

func UpdateUserInfo(ctx context.Context, UID string, data *UserInfoUpdate, db *sql.DB) error {
	if data.DISPLAY_NAME == "" {
		return fmt.Errorf("No Name Provided")
	}
	updateUserInfoCommand := `UPDATE USERS SET DISPLAY_NAME = ?,PROFILE_PIC = ?, SETTINGS = ? WHERE UID = ?`
	results, err := db.ExecContext(ctx, updateUserInfoCommand, data.DISPLAY_NAME, data.PROFILE_PIC, data.SETTINGS, UID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil
}

// .
// .
// .
// .
// .
// .
func CreateRecord(ctx context.Context, data *DATA, db *sql.DB) error {
	dataInsertCommand := `INSERT INTO DATA (TIME_UUID, NAME, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH) VALUES (?, ?, ?, ?, ?);`

	if data.TIME_UUID == "" {
		return fmt.Errorf("No UUID Provided")
	}
	if data.NAME == "" {
		return fmt.Errorf("No Name Provided")
	}
	if data.FILEPATH == "" {
		return fmt.Errorf("No File Provided")
	}
	if data.CREATOR_ID == "" {
		return fmt.Errorf("No Creator Provided")
	}

	_, err := db.ExecContext(ctx, dataInsertCommand, data.TIME_UUID, data.NAME, data.FILEPATH, data.CREATOR_ID, data.PREVIEW_IMG_PATH)

	return err
}
func DeleteRecord(ctx context.Context, UID string, db *sql.DB) error {
	dataDeleteCommand := `DELETE FROM DATA WHERE TIME_UUID = ?`
	results, err := db.ExecContext(ctx, dataDeleteCommand, UID)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil
}
func GetRecord(ctx context.Context, UID string, db *sql.DB) (*DATA, error) {
	data := &DATA{}
	getRecordQuery := `SELECT TIME_UUID,NAME,FILEPATH,CREATOR_ID,PREVIEW_IMG_PATH FROM DATA WHERE TIME_UUID = ?`

	result := db.QueryRowContext(ctx, getRecordQuery, UID)
	err := result.Scan(&data.TIME_UUID, &data.NAME, &data.FILEPATH, &data.CREATOR_ID, &data.PREVIEW_IMG_PATH)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No Rows Found")
		}
		return nil, err
	}

	return data, nil
}
func UpdateRecord(ctx context.Context, UID string, data *DataInfoUpdate, db *sql.DB) error {
	updateQuery := `UPDATE DATA SET NAME = ?,PREVIEW_IMG_PATH = ? WHERE TIME_UUID = ?`

	if data.NAME == "" {
		return fmt.Errorf("Name Not provided for Updating")
	}

	results, err := db.ExecContext(ctx, updateQuery, data.NAME, data.PREVIEW_IMG_PATH, UID)
	if err != nil {
		return err
	}

	rowsAffected, err := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil
}

// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .
// .

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

	if len(result) == 0 {
		return result, fmt.Errorf("No Files Found!")
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
