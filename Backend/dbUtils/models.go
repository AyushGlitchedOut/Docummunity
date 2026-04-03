package dbUtils

type TestData struct {
	TIMEID   int
	NAME     string
	FILEPATH string
}

type DATA struct {
	TIME_UUID        string
	NAME             string
	DESCRIPTION      string
	FILEPATH         string
	CREATOR_ID       string
	PREVIEW_IMG_PATH string
}
type DataInfoUpdate struct {
	NAME             string
	PREVIEW_IMG_PATH string
	DESCRIPTION      string
}

type USER struct {
	UID           string
	EMAIL         string
	DISPLAY_NAME  string
	BIO           string
	PROFILE_PIC   string
	CREATION_DATE string
	SETTINGS      string
}
type UserInfoUpdate struct {
	DISPLAY_NAME string
	PROFILE_PIC  string
	SETTINGS     string
	BIO          string
}
