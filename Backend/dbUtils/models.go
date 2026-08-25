package dbUtils

// Base struct for model of Data (Record)
type DATA struct {
	UUID             string
	NAME             string
	DESCRIPTION      string
	FILEPATH         string
	CREATOR_ID       string
	PREVIEW_IMG_PATH string
	CREATION_DATE    int64
}

// Struct to handle Info regarding updating a record
type DataInfoUpdate struct {
	NAME             string
	PREVIEW_IMG_PATH string
	DESCRIPTION      string
}

// Base struct for model of User
type USER struct {
	UID           string
	DISPLAY_NAME  string
	BIO           string
	PROFILE_PIC   string
	CREATION_DATE int64
	SETTINGS      string
}

// Created as when more stuff will be stored per user, we need to omit personal info
type USER_PUBLIC struct {
	UID           string
	DISPLAY_NAME  string
	BIO           string
	PROFILE_PIC   string
	CREATION_DATE int64
}

// Struct to handle Info regarding updating a User
type UserInfoUpdate struct {
	DISPLAY_NAME string
	PROFILE_PIC  string
	SETTINGS     string
	BIO          string
}
