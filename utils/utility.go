package utils

import (
	"mime"
	"mime/multipart"
	"strings"
)

func IsFormKeyPresent(key string, mapCheck map[string][]string) bool {
	if _, ex := mapCheck[key]; ex {
		return true
	}
	return false
}

func IsFormFileKeyPresent(key string, mapCheck map[string][]*multipart.FileHeader) bool {
	if _, ex := mapCheck[key]; ex {
		return true
	}
	return false
}

func IsKeyPresent(key any, mapCheck map[any][]any) bool {
	if _, ex := mapCheck[key]; ex {
		return true
	}
	return false
}

func FileExtension(filename string) (ext string) {
	if strings.Contains(filename, ".") {
		names := strings.Split(filename, ".")
		ext = strings.ToLower(names[len(names)-1])
	}

	return
}

func IsImageFile(filename string) bool {
	ext := FileExtension(filename)
	if ext == "" {
		return false
	}

	// Guess the MIME type based on the file extension
	mimeType := mime.TypeByExtension("." + ext)
	if mimeType == "" {
		return false
	}

	// Check if the MIME type starts with "image/"
	return strings.HasPrefix(mimeType, "image/")
}
