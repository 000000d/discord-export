package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type MessageBlock struct {
	ID        string `json:"id,omitempty"`
	Content   string `json:"content,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("pass channel id as cli arg.")
	}

	var authFileName string = "auth.txt"

	readAuthFile, err := os.ReadFile(authFileName)
	if err != nil || string(readAuthFile) == "" {
		os.Create(authFileName)
		log.Printf("paste discord token inside '%s'.", authFileName)
		os.Exit(-1)
	}

	var auth string = string(readAuthFile)
	var apiVersion string = "v9"
	var channelID string = os.Args[1]
	var apiLimit uint8 = 100
	var baseAPIUrl string = fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d", apiVersion, channelID, apiLimit)

	Call(baseAPIUrl, auth, apiVersion, channelID, apiLimit)
}

func Call(baseAPIUrl, auth, apiVersion, channelID string, apiLimit uint8) {
	var prevMessageID string = "0"
	var client *http.Client = &http.Client{
		Timeout: 5 * time.Second,
	}
	var msgJSONList []MessageBlock

	req, err := http.NewRequest("GET", baseAPIUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", auth)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(resBody, &msgJSONList)
	if err != nil {
		log.Printf("incorrect auth / token. received: '%s'", auth)
		log.Fatal(err)
	}

	for _, message := range msgJSONList {
		prevMessageID = message.ID
		fmt.Println(message.Content)
	}

	var beforeParamAPIUrl string = fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d&before=%s", apiVersion, channelID, apiLimit, prevMessageID)

	if len(msgJSONList) < int(apiLimit) {
		os.Exit(0)
	}

	time.Sleep(500 * time.Millisecond)
	Call(beforeParamAPIUrl, auth, apiVersion, channelID, apiLimit)
}
