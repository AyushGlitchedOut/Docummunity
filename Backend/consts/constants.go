package consts

// Constants used in the App
const (

	//Rate-Limits
	//Requests the User can send per second to not reach their Limit
	RequestsPerSecond = 1
	//Maximum Requests the user can send together in a spike
	Burst = 5

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

// Outside const since only primary types can be used in const declaration

// Image Types Allowed
var AllowedImageExtensions = []string{
	"image/jpeg",
	"image/png",
	"image/webp",
}

// Document Types allowed
var AllowedDocumentExtensions = []string{
	"application/pdf",
}
