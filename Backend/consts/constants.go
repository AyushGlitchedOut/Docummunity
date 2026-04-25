package consts

const (
	MaxPerRequestServerMemorySize = 64 << 20

	MaxPictureSize  = 2 << 20  //2mb
	MaxDocumentSize = 40 << 20 //40mb

	//File storage Locations
	UploadsDirectory    = "./uploads"
	DBLocation          = UploadsDirectory + "/DATABASE.db"
	ProfilePicDirectory = UploadsDirectory + "/PROFILE_PIC/"
	FileDirectory       = UploadsDirectory + "/FILES/"
	PreviewImgDirectory = UploadsDirectory + "/PREVIEW/"
)
