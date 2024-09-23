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
	1006: "email is required.",
	1007: "failed to verify user. Please check again later OR try with different username and email.",
	1008: "username is taken.",
	1009: "email already registed.",
	1010: "allowed extensions for avatar - .jpg, .jpeg, .png",
	1011: "avatar filetype is not image",
	1012: "please check avatar filesize. Limit 0-500KB",
	1013: "unable to read avatar file. Please try again.",
	1014: "allowed extensions for coverImage - .jpg, .jpeg, .png",
	1015: "coverImage filetype is not image",
	1016: "please check coverImage filesize. Limit 0-800KB",
	1017: "unable to read coverImage file. Please try again.",
	1018: "coverImage file is required.",
	1019: "fullName cannot exceed 64 characters.",
	1020: "failed to upload avatar file. Please try again.",
	1021: "failed to upload coverImage file. Please try again.",
	1022: "failed to generate tokens.",
	1023: "password is required.",
	1024: "password must have more than 8 and less than 25 characters, 1 uppercase, 1 lowercase and 1 special character.",
	1025: "unexpected error occured. failed to check password, please try with a different one.",
	1026: "unable to read input.",
	1027: "credentials are required to login.",
	1028: "username is not registered.",
	1029: "password didn't match.",

	2001:"unexpected error occured. Failed to register user",


	5000: "unexpected error occured. please contact admin.",
	5001: "request timed out. Please try again later.",
}

var AllowedImagesExt = map[string]bool{
	"jpg":  true,
	"png":  true,
	"jpeg": true,
}
