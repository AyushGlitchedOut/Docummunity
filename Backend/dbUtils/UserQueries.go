package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
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

func CreateUser(ctx context.Context, user *USER, db *sql.DB) error {
	userInsertCommand := `INSERT INTO USERS (UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, userInsertCommand, user.UID, user.DISPLAY_NAME, user.BIO, user.PROFILE_PIC, user.CREATION_DATE, user.SETTINGS)
	if err != nil {
		return err
	}

	return nil
}

func GetUserAccount(ctx context.Context, UID string, db DbTxCombiner) (*USER, error) {
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

// To get the records created by a specific user
func GetUserRecords(ctx context.Context, UID string, db DbTxCombiner) ([]*DATA, error) {
	var records []*DATA

	getUserRecordsCommand := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE CREATOR_ID = ?`

	results, err := db.QueryContext(ctx, getUserRecordsCommand, UID)
	if err != nil {
		return nil, err
	}
	defer results.Close()

	for results.Next() {
		data := &DATA{}

		err := results.Scan(&data.UUID, &data.NAME, &data.DESCRIPTION, &data.FILEPATH, &data.CREATOR_ID, &data.PREVIEW_IMG_PATH)
		if err != nil {
			return nil, err
		}
		records = append(records, data)

	}
	if results.Err() != nil {
		return nil, results.Err()
	}

	return records, nil
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

	var filesToDelete []string
	var previewsToDelete []string

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
		var records []*DATA
		records, err = GetUserRecords(ctx, userID, transaction)
		if err != nil {
			return err
		}
		for _, value := range records {
			filesToDelete = append(filesToDelete, value.FILEPATH)
			previewsToDelete = append(previewsToDelete, value.PREVIEW_IMG_PATH)
		}

		deleteCommand := `DELETE FROM DATA WHERE CREATOR_ID = ?`
		_, err = transaction.ExecContext(ctx, deleteCommand, userID)
		if err != nil {
			return err
		}

	}

	//Get preview Image Path
	profilePicturePath := ""
	result, err := GetUserAccount(ctx, userID, transaction)
	if err != nil {
		return err
	}
	profilePicturePath = result.PROFILE_PIC

	deleteUserCommand := `DELETE FROM USERS WHERE UID = ?`
	results, err := transaction.ExecContext(ctx, deleteUserCommand, userID)
	if err != nil {
		log.Println("2")
		return err
	}
	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No User found with UID: %s", userID)
	}

	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("Failed to Commit Transaction %w", err)
	}

	//Remove profile picture image
	if profilePicturePath != "" {
		err = os.Remove(profilePicturePath)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Error Deleting Preview Image")
			return err
		}
	}

	//Remove files
	for _, val := range filesToDelete {
		if val == "" {
			continue
		}
		err = os.Remove(val)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Error deleting File:", val)
		}

	}

	//Remove Previews
	for _, val := range previewsToDelete {
		if val == "" {
			continue
		}
		err = os.Remove(val)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Error deleting Preview:", val)
		}
	}

	return nil
}
