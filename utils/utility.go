package utils

func IsFormKeyPresent(key string, mapCheck map[string][]string) bool {
	if _, ex := mapCheck[key]; ex{
		return true
	}
	return false
}


func IsKeyPresent(key any, mapCheck map[any][]any) bool {
	if _, ex := mapCheck[key]; ex{
		return true
	}
	return false
}