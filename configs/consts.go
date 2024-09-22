package configs

const (
	DB_NAME string = "videotube"
)

var ServiceCodes = map[int]string{
	1001: "unable to read user form data.",
	1002: "username cannot be blank.",
	1003: "username is required.",
	1004: "email cannot be blank.",
	1005: "invalid email.",
	1006: "username is required.",
	1007: "failed to verify user. Please check again later OR try with different username and email.",
	1008: "username is taken.",
	1009: "email already registed.",
	1010: "allowed extensions for avatar - .jpg, .jpeg, .png",
	1011: "avatar filetype is not image",
	1012: "please check avatar filesize. Limit 0-800KB",
	1013: "unable to read avatar file. Please try again.",
	1014: "allowed extensions for coverImage - .jpg, .jpeg, .png",
	1015: "coverImage filetype is not image",
	1016: "please check coverImage filesize. Limit 0-500KB",
	1017: "unable to read coverImage file. Please try again.",
	1018: "coverImage file is required.",
	1019: "fullName cannot exceed 64 characters.",
	5000: "unexpected error occured. please contact admin.",
}

var AllowedImagesExt = map[string]bool{
	"jpg":  true,
	"png":  true,
	"jpeg": true,
}
