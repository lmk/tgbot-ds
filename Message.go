package main

import (
	"math/rand"
)

type synoError struct {
	Code    int
	Message string
}

func (e *synoError) Error() string {
	return e.Message
}

// SynoLoginError syno api error code to string
func SynoLoginError(code int) error {
	m := map[int]string{
		// common
		100: "Unknown error",
		101: "Invalid parameter",
		102: "The requested API does not exist",
		103: "The requested method does not exist",
		104: "The requested version does not support the functionality",
		105: "The logged in session does not have permission",
		106: "Session timeout",
		107: "Session interrupted by duplicate login",
		// login
		400: "No such account or incorrect password",
		401: "Account disabled",
		402: "Permission denied",
		403: "2-step verification code required",
		404: "Failed to authenticate 2-step verification code",
	}

	return &synoError{Code: code, Message: m[code]}
}

// SynoCreateError syno api error code to string
func SynoCreateError(code int) error {
	m := map[int]string{
		// common
		100: "Unknown error",
		101: "Invalid parameter",
		102: "The requested API does not exist",
		103: "The requested method does not exist",
		104: "The requested version does not support the functionality",
		105: "The logged in session does not have permission",
		106: "Session timeout",
		// create
		400: "File upload failed",
		401: "Max number of tasks reached",
		402: "Destination denied",
		403: "Destination does not exist",
		404: "Invalid task id",
		405: "Invalid task action",
		406: "No default destination",
		407: "Set destination failed",
		408: "File does not exist",
	}

	return &synoError{Code: code, Message: m[code]}
}

// OkMesssge random ok message
func OkMesssge() string {
	msgs := []string{
		"Ok",
		"I got it",
		"Check",
		"I love it",
		"I like it",
		"interesting",
		"your taste",
	}

	return msgs[rand.Intn(len(msgs))]
}
