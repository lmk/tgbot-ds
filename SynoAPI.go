package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

type sError struct {
	code int
}

type sLoginData struct {
	Sid string
}

type sLoginResult struct {
	Success bool
	Data    sLoginData
	Error   sError
}

type sCreateResult struct {
	Success bool
	Error   sError
}

// CreateMagnet create magent uri to destnation with SYNO.API
func CreateMagnet(dst string, uri string) error {

	host := os.Getenv("SYNOHOST")
	synoID := os.Getenv("SYNOID")
	synoPwd := os.Getenv("SYNOPWD")

	// login
	resp, err := http.PostForm(fmt.Sprintf("http://%s/webapi/auth.cgi", host),
		url.Values{
			"api":     {"SYNO.API.Auth"},
			"version": {"2"},
			"method":  {"login"},
			"account": {synoID},
			"passwd":  {synoPwd},
			"session": {"DownloadStation"},
			"format":  {"sid"}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	//log.Printf("[LOGIN] %s", string(byteBody))

	var loginResult sLoginResult

	json.Unmarshal(byteBody, &loginResult)
	//log.Printf("[LOGIN] success: %t", loginResult.Success)
	//log.Printf("[LOGIN] sid: %s", loginResult.Data.Sid)
	if loginResult.Success == false {
		return SynoError(loginResult.Error.code)
	}

	sid := loginResult.Data.Sid

	// create magent
	resp, err = http.PostForm(fmt.Sprintf("http://%s/webapi/DownloadStation/task.cgi", host),
		url.Values{
			"api":         {"SYNO.DownloadStation.Task"},
			"version":     {"1"},
			"method":      {"create"},
			"destination": {dst},
			"uri":         {uri},
			"_sid":        {sid}})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	//log.Printf("[CREATE] %s", string(byteBody))

	var createResult sCreateResult
	json.Unmarshal(byteBody, &createResult)
	//log.Printf("[CREATE] success: %t", createResult.Success)
	if createResult.Success == false {
		//log.Printf("[CREATE] ERROR %d", createResult.Error.code)
		return SynoError(createResult.Error.code)
	}

	return nil
}
