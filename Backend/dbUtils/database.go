package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Constant for reserved deleted User
var DeletedUserInfo = &USER{
	UID:           "000",
	BIO:           "Deleted User",
	DISPLAY_NAME:  "[DELETED]",
	PROFILE_PIC:   "none",
	CREATION_DATE: "0/0/0",
	SETTINGS:      "",
}

// Create a reserved-system user for assigning deleted users' files to
func createDeletedUser(ctx context.Context, db *sql.DB) error {
	//Change to INSERT IGNORE when porting to MySQL
	createDeletedUserCommand := `INSERT OR IGNORE INTO USERS (UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, createDeletedUserCommand, DeletedUserInfo.UID, DeletedUserInfo.DISPLAY_NAME, DeletedUserInfo.BIO, DeletedUserInfo.PROFILE_PIC, DeletedUserInfo.CREATION_DATE, DeletedUserInfo.SETTINGS)
	if err != nil {
		return err
	}
	return nil
}

func InitializeDB(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./uploads/DATABASE.db?_foreign_keys=on")
	if err != nil {
		return nil, err
	}

	//TODO: switch Time and Settings to a more efficient type
	userCreateCommand := `CREATE TABLE IF NOT EXISTS USERS (
		UID TEXT PRIMARY KEY, 
		DISPLAY_NAME TEXT,
		BIO TEXT,
		PROFILE_PIC TEXT,
		CREATION_DATE TEXT NOT NULL,
		SETTINGS TEXT
	);`

	_, err = db.ExecContext(ctx, userCreateCommand)
	if err != nil {
		return nil, err
	}

	//TODO: switch UID to a more efficient type
	dataCreateCommand := `CREATE TABLE IF NOT EXISTS DATA (
	UUID TEXT PRIMARY KEY, 
	NAME TEXT NOT NULL,
	DESCRIPTION TEXT,
	FILEPATH TEXT NOT NULL,
	CREATOR_ID TEXT NOT NULL,
	PREVIEW_IMG_PATH TEXT,
	FOREIGN KEY (CREATOR_ID) REFERENCES USERS(UID)
	);`
	_, err = db.ExecContext(ctx, dataCreateCommand)
	if err != nil {
		return nil, err
	}

	err = createDeletedUser(ctx, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateUser(ctx context.Context, user *USER, db *sql.DB) error {
	userInsertCommand := `INSERT INTO USERS (UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, userInsertCommand, user.UID, user.DISPLAY_NAME, user.BIO, user.PROFILE_PIC, user.CREATION_DATE, user.SETTINGS)
	if err != nil {
		return err
	}

	return nil
}

func GetUserAccount(ctx context.Context, UID string, db *sql.DB) (*USER, error) {
	user := &USER{}
	getUserCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS FROM USERS WHERE UID = ?`
	result := db.QueryRowContext(ctx, getUserCommand, UID)
	err := result.Scan(&user.UID, &user.DISPLAY_NAME, &user.BIO, &user.PROFILE_PIC, &user.CREATION_DATE, &user.SETTINGS)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No Rows Found")
		}
		return nil, err
	}

	return user, nil
}

func GetUserInfo(ctx context.Context, UID string, db *sql.DB) (*USER_PUBLIC, error) {
	user := &USER_PUBLIC{}
	getUserCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE FROM USERS WHERE UID = ?`
	result := db.QueryRowContext(ctx, getUserCommand, UID)
	err := result.Scan(&user.UID, &user.DISPLAY_NAME, &user.BIO, &user.PROFILE_PIC, &user.CREATION_DATE)
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
	updateUserInfoCommand := `UPDATE USERS SET DISPLAY_NAME = ?, BIO = ?,PROFILE_PIC = ?, SETTINGS = ? WHERE UID = ?`
	results, err := db.ExecContext(ctx, updateUserInfoCommand, data.DISPLAY_NAME, data.BIO, data.PROFILE_PIC, data.SETTINGS, UID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil

}

func SearchUser(ctx context.Context, query []string, db *sql.DB) ([]*USER_PUBLIC, error) {
	var users []*USER_PUBLIC

	if len(query) == 0 {
		return nil, fmt.Errorf("No Queries Provided")
	}
	searchCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE FROM USERS WHERE `

	var keywords []string
	var args []any
	for i := range query {
		keywords = append(keywords, "DISPLAY_NAME LIKE ?")
		args = append(args, "%"+query[i]+"%")
	}

	results, err := db.QueryContext(ctx, searchCommand+strings.Join(keywords, " OR "), args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		row := &USER_PUBLIC{}
		err = results.Scan(&row.UID, &row.DISPLAY_NAME, &row.BIO, &row.PROFILE_PIC, &row.CREATION_DATE)
		if err != nil {
			return nil, err
		}
		users = append(users, row)
	}
	if results.Err() != nil {
		return nil, results.Err()
	}

	return users, nil
}

func DeleteUser(ctx context.Context, userID string, db *sql.DB, keepRecords bool) error {
	//Using Transaction to Avoid Inconsistency
	if userID == "" {
		return fmt.Errorf("No User ID Given")
	}
	transaction, err := db.BeginTx(ctx, nil)

	if err != nil {
		return fmt.Errorf("Error Starting Transaction: %w", err)
	}
	defer transaction.Rollback()

	if keepRecords {
		detachRecordsCommand := `UPDATE DATA SET CREATOR_ID = ? WHERE CREATOR_ID = ?`
		_, err = transaction.ExecContext(ctx, detachRecordsCommand, DeletedUserInfo.UID, userID)
		if err != nil {
			return err
		}
	} else {
		deleteRecordsCommand := `DELETE FROM DATA WHERE CREATOR_ID = ?`

		_, err = transaction.ExecContext(ctx, deleteRecordsCommand, userID)
		if err != nil {
			return err
		}
	}

	deleteUserCommand := `DELETE FROM USERS WHERE UID = ?`

	results, err := transaction.ExecContext(ctx, deleteUserCommand, userID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No User found with UID %s", userID)
	}
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("Failed to Commit Transaction %w", err)
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
	dataInsertCommand := `INSERT INTO DATA (UUID, NAME, DESCRIPTION ,FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH) VALUES (?, ?, ?, ?, ?, ?);`

	if data.UUID == "" {
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

	_, err := db.ExecContext(ctx, dataInsertCommand, data.UUID, data.NAME, data.DESCRIPTION, data.FILEPATH, data.CREATOR_ID, data.PREVIEW_IMG_PATH)

	return err
}
func DeleteRecord(ctx context.Context, UID string, creatorID string, db *sql.DB) error {
	dataDeleteCommand := `DELETE FROM DATA WHERE UUID = ? AND CREATOR_ID = ?`
	results, err := db.ExecContext(ctx, dataDeleteCommand, UID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil
}
func GetRecord(ctx context.Context, UID string, db *sql.DB) (*DATA, error) {
	data := &DATA{}
	getRecordQuery := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE UUID = ?`

	result := db.QueryRowContext(ctx, getRecordQuery, UID)
	err := result.Scan(&data.UUID, &data.NAME, &data.DESCRIPTION, &data.FILEPATH, &data.CREATOR_ID, &data.PREVIEW_IMG_PATH)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No Rows Found")
		}
		return nil, err
	}

	return data, nil
}
func UpdateRecord(ctx context.Context, UID string, data *DataInfoUpdate, creatorID string, db *sql.DB) error {
	updateQuery := `UPDATE DATA SET NAME = ?,DESCRIPTION = ?,PREVIEW_IMG_PATH = ? WHERE UUID = ? AND CREATOR_ID = ?`

	if creatorID == "" {
		return fmt.Errorf("No Creator UID provided")
	}

	if data.NAME == "" {
		return fmt.Errorf("Name Not provided for Updating")
	}

	results, err := db.ExecContext(ctx, updateQuery, data.NAME, data.DESCRIPTION, data.PREVIEW_IMG_PATH, UID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Record found with UUID %s", UID)
	}

	return nil
}

func SearchRecord(ctx context.Context, query []string, db *sql.DB, useDescription bool) ([]*DATA, error) {
	var data []*DATA

	searchCommand := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE `

	if len(query) == 0 {
		return nil, fmt.Errorf("No Queries Provided")
	}

	var keyWords []string
	var args []any
	if useDescription {
		for i := range query {
			keyWords = append(keyWords, "(NAME LIKE ? OR DESCRIPTION LIKE ?)")
			args = append(args, "%"+query[i]+"%", "%"+query[i]+"%")
		}
	} else {
		for i := range query {
			keyWords = append(keyWords, "NAME LIKE ?")
			args = append(args, "%"+query[i]+"%")
		}
	}

	results, err := db.QueryContext(ctx, searchCommand+strings.Join(keyWords, " OR "), args...)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		row := &DATA{}
		err := results.Scan(&row.UUID, &row.NAME, &row.DESCRIPTION, &row.FILEPATH, &row.CREATOR_ID, &row.PREVIEW_IMG_PATH)
		if err != nil {
			return nil, err
		}
		data = append(data, row)
	}
	if results.Err() != nil {
		return nil, results.Err()
	}

	return data, nil
}
