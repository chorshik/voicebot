package v1

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	url = "https://api.convertio.co/convert"
)

// postBody ...

// GetID ...
func GetID(urlFile string, apikey string) (string, int64) {

	reqBody, err := json.Marshal(map[string]string{
		"apikey":       apikey,
		"input":        "url",
		"file":         urlFile,
		"outputformat": "opus",
	})
	if err != nil {
		log.Print(err)
	}

	client := &http.Client{}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	responseStr := string(read)

	id := gjson.Get(responseStr, "data.id").String()
	log.Print("tries left:" + "" + gjson.Get(responseStr, "data.minutes").String())
	tries := gjson.Get(responseStr, "data.minutes").Int()

	return id, tries
}

// GetStatus ...
func GetStatus(id string) string {
	status := isFinished(id)
	if status == false {
		log.Print("non id")
	}

	return id
}

// GetFile ...
func GetFile(id string) []byte {
	// Get(https://api.convertio.co/convert/:ID/dl/)
	path := url + "/" + id + "/dl/"

	resp, err := http.Get(path)
	if err != nil {
		log.Print(err)
	}
	defer resp.Body.Close()

	read, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print(err)
	}
	responseStr := string(read)

	sDec, err := base64.StdEncoding.DecodeString(gjson.Get(responseStr, "data.content").String())
	if err != nil {
		log.Print(err)
	}

	return sDec
}

func isFinished(id string) bool {
	var Isfinished bool = false
	duration := time.Duration(5 * time.Second)

	for Isfinished != true {
		resp, err := http.Get(url + "/" + id + "/status")
		if err != nil {
			log.Print(err)
		}
		defer resp.Body.Close()

		read, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
		}
		responseStr := string(read)

		step := gjson.Get(responseStr, "data.step").String()
		if step == "finish" {
			Isfinished = true
		}

		time.Sleep(duration)

	}

	return Isfinished
}
