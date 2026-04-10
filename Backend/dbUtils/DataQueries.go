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

// Interface implementing methods so that both Transaction and DB can be used as an argument
type DbTxCombiner interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
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
	FOREIGN KEY (CREATOR_ID) REFERENCES USERS(UID) ON DELETE CASCADE
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

func CreateRecord(ctx context.Context, data *DATA, db *sql.DB) error {
	dataInsertCommand := `INSERT INTO DATA (UUID, NAME, DESCRIPTION ,FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH) VALUES (?, ?, ?, ?, ?, ?);`

	_, err := db.ExecContext(ctx, dataInsertCommand, data.UUID, data.NAME, data.DESCRIPTION, data.FILEPATH, data.CREATOR_ID, data.PREVIEW_IMG_PATH)

	return err
}
func DeleteRecord(ctx context.Context, UUID string, creatorID string, db *sql.DB) error {

	record, err := GetRecord(ctx, UUID, db)
	if err != nil {
		return err
	}
	filePath := record.FILEPATH
	previewPath := record.PREVIEW_IMG_PATH

	dataDeleteCommand := `DELETE FROM DATA WHERE UUID = ? AND CREATOR_ID = ?`
	results, err := db.ExecContext(ctx, dataDeleteCommand, UUID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, _ := results.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("No Rows Found")
	}

	//delete preview
	if previewPath != "" {
		err = os.Remove(previewPath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	//delete file
	if filePath != "" {
		err = os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
func GetRecord(ctx context.Context, UUID string, db *sql.DB) (*DATA, error) {
	data := &DATA{}
	getRecordQuery := `SELECT UUID, NAME, DESCRIPTION, FILEPATH, CREATOR_ID, PREVIEW_IMG_PATH FROM DATA WHERE UUID = ?`

	result := db.QueryRowContext(ctx, getRecordQuery, UUID)
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
