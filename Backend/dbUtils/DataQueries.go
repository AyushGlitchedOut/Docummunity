package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//Honestly speaking, This interface is the only part in the entire project that I didn't come up with myself or design with myself
//I was stuck, and took the help of ChatGPT. I undestand it now of course, but it is the only thing I couldn;t
// think of myself and had to take the help of AI for.

// Interface implementing methods so that both Transaction and DB can be used as an argument in the same function
type DbTxCombiner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// Function to Initialize the database (Create DB file, establish connection, Create Tables, Create Deleted User)
func InitializeDB(ctx context.Context) (*sql.DB, error) {

	//Open the DB file
	DB, err := sql.Open("sqlite3", "./uploads/DATABASE.db?_foreign_keys=on")
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

	//Create User Table
	_, err = DB.ExecContext(ctx, userCreateCommand)
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
	FOREIGN KEY (CREATOR_ID) REFERENCES USERS(UID) ON DELETE CASCADE
	);`

	//Create Data Table
	_, err = DB.ExecContext(ctx, dataCreateCommand)
	if err != nil {
		return nil, err
	}

	//Create deleted User
	err = createDeletedUser(ctx, DB)
	if err != nil {
		return nil, err
	}

	return DB, nil
}

// To create a Record in the Table
func CreateRecord(ctx context.Context, data *DATA, db *sql.DB) error {
	dataInsertCommand := `INSERT INTO DATA (UUID, NAME, DESCRIPTION ,FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, dataInsertCommand, data.UUID, data.NAME, data.DESCRIPTION, data.FILEPATH, data.CREATOR_ID, data.PREVIEW_IMG_PATH)

	return err
}

// To obtain a record from the Table.
func GetRecord(ctx context.Context, UUID string, db DbTxCombiner) (*DATA, error) {
	data := &DATA{}
	getRecordQuery := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE UUID = ?`

	//QueryRowContext used since it returns only one row, which is what we want
	result := db.QueryRowContext(ctx, getRecordQuery, UUID)

	//scan the result into data struct
	err := result.Scan(&data.UUID, &data.NAME, &data.DESCRIPTION, &data.FILEPATH, &data.CREATOR_ID, &data.PREVIEW_IMG_PATH)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("No Records Found")
		}
		return nil, err
	}

	return data, nil
}

// To delete a record from the table, along with its files
func DeleteRecord(ctx context.Context, UUID string, creatorID string, db *sql.DB) error {

	//Using transaction instead of DB directly for atomicity
	transaction, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error Starting Transaction: %w", err)
	}

	//make sure the transaction rollsback even if the function fails and errors
	defer transaction.Rollback()

	//Use the previously created DB function GetRecord to get additional info (paths of files to delete) about the record to delete
	record, err := GetRecord(ctx, UUID, transaction)
	if err != nil {
		return err
	}

	//File and Preview Images to delete
	filePath := record.FILEPATH
	previewPath := record.PREVIEW_IMG_PATH

	//Delete the record from the DB
	dataDeleteCommand := `DELETE FROM DATA WHERE UUID = ? AND CREATOR_ID = ?`
	results, err := transaction.ExecContext(ctx, dataDeleteCommand, UUID, creatorID)
	if err != nil {
		return err
	}

	//Error out If no rows are deleted i.e. Record Not Found
	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Records Found")
	}

	//commit the transaction
	err = transaction.Commit()
	if err != nil {
		return fmt.Errorf("Failed to Commit Transaction %w", err)
	}

	//NOTE: Preview Image and Document File are deleted once the transaction is commited because FS operations are not and cannot be atomic.
	// If they are put inside the transaction, it will be of no use since if the transaction is rolled back upon failure of file deletion, there will be an incosistent state where the record exists but the actual files Dont.
	//Instead, its better to completely delete the Record first and then delete the file, since orphaned files(which are possible to be cleaned up by a future tool) are better than incosistent state

	//delete preview File
	if previewPath != "" {
		err = os.Remove(previewPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	//delete Document file
	if filePath != "" {
		err = os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// To update a Record in the Table
func UpdateRecord(ctx context.Context, UID string, data *DataInfoUpdate, creatorID string, db *sql.DB) error {

	var updateQuery string
	var results sql.Result
	var err error

	//Update the entry
	updateQuery = `UPDATE DATA SET NAME = ?,DESCRIPTION = ?,PREVIEW_IMG_PATH = ? WHERE UUID = ? AND CREATOR_ID = ?`
	results, err = db.ExecContext(ctx, updateQuery, data.NAME, data.DESCRIPTION, data.PREVIEW_IMG_PATH, UID, creatorID)
	if err != nil {
		return err
	}

	//Error out If no rows are Updated i.e. Record Not Found
	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Records Found")
	}

	return nil
}

// To search a record by Name and/or Description query. The last boolean is useDescription, stating whether to check Description while querying or not
func SearchRecord(ctx context.Context, query []string, db *sql.DB, useDescription bool) ([]*DATA, error) {
	var data []*DATA

	searchCommand := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE `

	//Combine the query words given into one string that can be appended to the SQL command
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

	//search the results
	results, err := db.QueryContext(ctx, searchCommand+strings.Join(keyWords, " OR "), args...) //searchCommand is the initial part of SQL command, strings.Join(keywords, " OR ") constructs the repeating string of NAME/DESCRIPTION LIKE ?, and args... finally pass the queries
	if err != nil {
		return nil, err
	}
	//For results containing 1+ rows, its important to close them afterwards.
	defer results.Close()

	//scan the results into data struct
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
	if len(data) < 1 {
		return nil, fmt.Errorf("No Records Found")
	}

	return data, nil
}
