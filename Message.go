package main

import (
	"errors"
	"math/rand"
)

// SynoError syno api error code to string
// login and create
func SynoError(code int) error {
	m := map[int]string{
		// login
		100: "Unknown error",
		101: "Invalid parameter",
		102: "The requested API does not exist",
		103: "The requested method does not exist",
		104: "The requested version does not support the functionality",
		105: "The logged in session does not have permission",
		106: "Session timeout",
		107: "Session interrupted by duplicate login",
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

	return errors.New(m[code])
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
