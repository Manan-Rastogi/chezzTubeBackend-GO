package configs

const (
	DB_NAME string = "videotube"
)

var ServiceCodes = map[int]string{
	1001: "unable to read user form data.",
	1002: "username cannot be blank.",
	1003: "username is required.",
}