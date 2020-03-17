package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type sResult struct {
	Success bool
	Sid     string
	Error   int
}

func (s *sResult) UnmarshalJSON(b []byte) error {
	var f interface{}
	json.Unmarshal(b, &f)

	m := f.(map[string]interface{})

	s.Success = m["success"].(bool)

	sub := m["error"]
	if sub != nil {
		v := sub.(map[string]interface{})
		s.Error = int(v["code"].(float64))
	}

	sub = m["data"]
	if sub != nil {
		v := sub.(map[string]interface{})
		s.Sid = v["sid"].(string)
	}

	return nil
}

// CreateMagnet create magent uri to destnation with SYNO.API
func CreateMagnet(dst string, uri string) error {

	host := os.Getenv("SYNOHOST")
	synoID := os.Getenv("SYNOID")
	synoPwd := os.Getenv("SYNOPWD")

	if debug {
		log.Printf("[LOGIN] SYNOHOST:[%s]", string(host))
		log.Printf("[LOGIN] SYNOID:[%s]", string(synoID))
		log.Printf("[LOGIN] SYNOPWD:[%s]", string(synoPwd))
	}

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
		log.Printf("[LOGIN] Fail read body")
		return err
	}

	if debug {
		log.Printf("[LOGIN] %s", string(byteBody))
	}

	var loginResult sResult
	err = loginResult.UnmarshalJSON(byteBody)
	if err != nil {
		return err
	}

	if loginResult.Success == false {
		log.Printf("[LOGIN] error code: %d", loginResult.Error)
		return SynoLoginError(loginResult.Error)
	}

	// create magent
	resp, err = http.PostForm(fmt.Sprintf("http://%s/webapi/DownloadStation/task.cgi", host),
		url.Values{
			"api":         {"SYNO.DownloadStation.Task"},
			"version":     {"2"},
			"method":      {"create"},
			"destination": {dst},
			"uri":         {uri},
			"_sid":        {loginResult.Sid}})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if debug {
		log.Printf("[CREATE] %s", string(byteBody))
	}

	var createResult sResult
	err = createResult.UnmarshalJSON(byteBody)
	if err != nil {
		return err
	}

	if createResult.Success == false {
		log.Printf("[CREATE] ERROR %d", createResult.Error)
		return SynoCreateError(createResult.Error)
	}

	return nil
}

// QuerySynoAPI query
func QuerySynoAPI(bot *tgbotapi.BotAPI, query string) error {

	host := os.Getenv("SYNOHOST")
	synoID := os.Getenv("SYNOID")
	synoPwd := os.Getenv("SYNOPWD")

	if debug {
		log.Printf("[LOGIN] SYNOHOST:[%s]", string(host))
		log.Printf("[LOGIN] SYNOID:[%s]", string(synoID))
		log.Printf("[LOGIN] SYNOPWD:[%s]", string(synoPwd))
	}

	if query == "" {
		query = "SYNO.API.Auth,SYNO.DownloadStation.Task"
	}

	// query
	resp, err := http.PostForm(fmt.Sprintf("http://%s/webapi/query.cgi", host),
		url.Values{
			"api":     {"SYNO.API.Info"},
			"version": {"1"},
			"method":  {"query"},
			"query":   {query}})
	//			"query":   {"all"}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	byteBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[QUERY] Fail read body")
		return err
	}

	if debug {
		log.Printf("[QUERY] %s", string(byteBody))
	}

	f, err := os.Create("query.json")
	if err != nil {
		log.Printf("[QUERY] Fail Create query.json")
		return err
	}

	f.Write(byteBody)
	if err != nil {
		log.Printf("[QUERY] Fail Write query.json")
		return err
	}

	return nil
}
