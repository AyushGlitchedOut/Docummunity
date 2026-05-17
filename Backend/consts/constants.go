package consts

// Constants used in the App
const (
	//Maximum Memory Server is allowed to use
	MaxPerRequestServerMemorySize = 64 << 20

	MaxPictureSize  = 2 << 20  //2mb
	MaxDocumentSize = 40 << 20 //40mb

	//File storage Locations (Hardcoded disk directories for now, will change in the future when architecture changes)
	UploadsDirectory    = "./uploads"
	DBLocation          = UploadsDirectory + "/DATABASE.db"
	ProfilePicDirectory = UploadsDirectory + "/PROFILE_PIC/"
	FileDirectory       = UploadsDirectory + "/FILES/"
	PreviewImgDirectory = UploadsDirectory + "/PREVIEW/"
)
