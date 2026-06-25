package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
)

// Constant User data for reserved deleted User
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
	//TODO: Change to INSERT IGNORE when porting to MySQL
	createDeletedUserCommand := `INSERT OR IGNORE INTO USERS (UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`

	//Create the User
	_, err := db.ExecContext(ctx, createDeletedUserCommand, DeletedUserInfo.UID, DeletedUserInfo.DISPLAY_NAME, DeletedUserInfo.BIO, DeletedUserInfo.PROFILE_PIC, DeletedUserInfo.CREATION_DATE, DeletedUserInfo.SETTINGS)
	if err != nil {
		return err
	}

	return nil
}

// To Create a User in the table
func CreateUser(ctx context.Context, user *USER, db *sql.DB) error {
	userInsertCommand := `INSERT INTO USERS (UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, userInsertCommand, user.UID, user.DISPLAY_NAME, user.BIO, user.PROFILE_PIC, user.CREATION_DATE, user.SETTINGS)
	if err != nil {
		return err
	}

	return nil
}

// To get User's Account Info (All info about the user, meant for User's own use upon authentication)
func GetUserAccount(ctx context.Context, UID string, db DbTxCombiner) (*USER, error) {
	user := &USER{}

	getUserCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE, SETTINGS FROM USERS WHERE UID = ?`
	//QueryRowContext used since it returns only one entry which is what we want
	result := db.QueryRowContext(ctx, getUserCommand, UID)

	//Scan the results into user struct
	err := result.Scan(&user.UID, &user.DISPLAY_NAME, &user.BIO, &user.PROFILE_PIC, &user.CREATION_DATE, &user.SETTINGS)
	if err != nil {
		//Error out If no rows are Found i.e. User Not Found
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No User Found")
		}
		return nil, err
	}

	return user, nil
}

// Get User's Public Info (Public and safe info about user, meant for other users to see)
func GetUserInfo(ctx context.Context, UID string, db *sql.DB) (*USER_PUBLIC, error) {
	user := &USER_PUBLIC{}

	getUserCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE FROM USERS WHERE UID = ?`
	//QueryRowContext used since it returns only one entry which is what we want
	result := db.QueryRowContext(ctx, getUserCommand, UID)

	//Scan the results into user struct
	err := result.Scan(&user.UID, &user.DISPLAY_NAME, &user.BIO, &user.PROFILE_PIC, &user.CREATION_DATE)
	if err != nil {
		if err == sql.ErrNoRows {
			//Error out If no rows are Found i.e. User Not Found
			return nil, fmt.Errorf("No User Found")
		}
		return nil, err
	}

	return user, nil
}

// To get the records created by a specific user
func GetUserRecords(ctx context.Context, UID string, db DbTxCombiner) ([]*DATA, error) {
	var records []*DATA

	getUserRecordsCommand := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE CREATOR_ID = ?`

	//get the User's Records
	results, err := db.QueryContext(ctx, getUserRecordsCommand, UID)
	if err != nil {
		return nil, err
	}
	//For results containing 1+ rows, its important to close them afterwards.
	defer results.Close()

	//Scan the results into the records struct
	for results.Next() {
		data := &DATA{}

		err := results.Scan(&data.UUID, &data.NAME, &data.DESCRIPTION, &data.FILEPATH, &data.CREATOR_ID, &data.PREVIEW_IMG_PATH)
		if err != nil {
			return nil, err
		}
		records = append(records, data)

	}
	//Check for errors during scanning
	if results.Err() != nil {
		return nil, results.Err()
	}

	//Check if the query was empty
	if len(records) < 1 {
		return nil, fmt.Errorf("No Records Found")
	}

	return records, nil
}

// To update the Information for a specific user
func UpdateUserInfo(ctx context.Context, UID string, data *UserInfoUpdate, db *sql.DB) error {
	updateUserInfoCommand := `UPDATE USERS SET DISPLAY_NAME = ?, BIO = ?,PROFILE_PIC = ?, SETTINGS = ? WHERE UID = ?`

	//Update the user's Info
	results, err := db.ExecContext(ctx, updateUserInfoCommand, data.DISPLAY_NAME, data.BIO, data.PROFILE_PIC, data.SETTINGS, UID)
	if err != nil {
		return err
	}

	//Error out If no rows are Found i.e. User Not Found
	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No User Found")
	}

	return nil

}

// To search for a specific User for thw query in its Display Name
func SearchUser(ctx context.Context, query []string, db *sql.DB) ([]*USER_PUBLIC, error) {
	var users []*USER_PUBLIC

	searchCommand := `SELECT UID, DISPLAY_NAME, BIO, PROFILE_PIC, CREATION_DATE FROM USERS WHERE `

	//Merge the query words given into one string to append to the SQL query
	var keywords []string
	var args []any
	for i := range query {
		keywords = append(keywords, "DISPLAY_NAME LIKE ?")
		args = append(args, "%"+query[i]+"%")
	}

	results, err := db.QueryContext(ctx, searchCommand+strings.Join(keywords, " OR "), args...) //searchCommand is the initial SQL syntax, and strings.Join(keywords, " OR ") appends the repeated query words with "LIKE ? OR", and args... pass the query arguments
	if err != nil {
		return nil, err
	}
	//For results containing 1+ rows, its important to close them afterwards.
	defer results.Close()

	//Scan the rresults into the users struct
	for results.Next() {
		row := &USER_PUBLIC{}
		err = results.Scan(&row.UID, &row.DISPLAY_NAME, &row.BIO, &row.PROFILE_PIC, &row.CREATION_DATE)
		//To prevent the user from accessing deleted exclusive user
		if row.UID == DeletedUserInfo.UID {
			continue
		}
		if err != nil {
			return nil, err
		}
		users = append(users, row)
	}
	if results.Err() != nil {
		return nil, results.Err()
	}
	if len(users) < 1 {
		return nil, fmt.Errorf("No User Found")
	}

	return users, nil
}

// To delete a User from the table, along with its files, the last boolean is keepRecords which decides if upon deletion, the records belonging to the user will be preserved or not
func DeleteUser(ctx context.Context, userID string, db *sql.DB, keepRecords bool) error {

	var filesToDelete []string
	var previewsToDelete []string

	if userID == "" {
		return fmt.Errorf("No User ID Given")
	}

	//Using Transaction to Avoid Inconsistency
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error Starting Transaction: %w", err)
	}
	//make sure the transaction rollsback even if the function fails and errors
	defer transaction.Rollback()

	if keepRecords {
		//If records are to be kept, and only user is to be deleted, we modify all records to have the creatorID of our deleted system User so they are preserved
		detachRecordsCommand := `UPDATE DATA SET CREATOR_ID = ? WHERE CREATOR_ID = ?`
		_, err = transaction.ExecContext(ctx, detachRecordsCommand, DeletedUserInfo.UID, userID)
		if err != nil {
			return err
		}
	} else {
		//If records are to be deleted as well, we also delete the records from the table belonging to the user along with files belonging to those records.
		//Also, we are writing a custom delete function here and not using the pre-defined recordDelete function since the pre-defined one deletes files along with the record in one go for each record, whereas we want to delete all the records first amd then all the files later after commiting the transaction to preserve consistency.
		var records []*DATA

		//Get all records belonging to the user
		records, err = GetUserRecords(ctx, userID, transaction)
		if err != nil {
			if !strings.Contains(err.Error(), "No Records Found") {
				return err
			}
		}

		//append the filepaths for document files and preview Images to a seperate array
		for _, value := range records {
			filesToDelete = append(filesToDelete, value.FILEPATH)
			previewsToDelete = append(previewsToDelete, value.PREVIEW_IMG_PATH)
		}

		//Delete the Records from the Table
		deleteCommand := `DELETE FROM DATA WHERE CREATOR_ID = ?`
		_, err = transaction.ExecContext(ctx, deleteCommand, userID)
		if err != nil {
			return err
		}

	}

	//Get preview Image Path of the User to delete
	profilePicturePath := ""
	result, err := GetUserAccount(ctx, userID, transaction)
	if err != nil {
		return err
	}
	profilePicturePath = result.PROFILE_PIC

	//Delete the User finally
	deleteUserCommand := `DELETE FROM USERS WHERE UID = ?`
	results, err := transaction.ExecContext(ctx, deleteUserCommand, userID)
	if err != nil {
		return err
	}

	//Error out If no rows are Found i.e. User Not Found
	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No User Found")
	}

	//Commit the Transaction
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("Failed to Commit Transaction %w", err)
	}

	//Remove profile picture image File
	if profilePicturePath != "" {
		err = os.Remove(profilePicturePath)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Error Deleting Preview Image")
			return err
		}
	}

	//Remove Document Files from the paths stored earlier before deleting the records
	for _, val := range filesToDelete {
		if val == "" {
			continue
		}
		err = os.Remove(val)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Error deleting File:", val)
		}

	}

	//Remove Preview Image Files from the paths stored earlier before deleting the records
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
