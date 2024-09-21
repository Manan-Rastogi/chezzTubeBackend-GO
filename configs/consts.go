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

	5000: "unexpected error occured. please contact admin.",
}
