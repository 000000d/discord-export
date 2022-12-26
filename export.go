package main

import (
	"encoding/json"
	"errors"
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
	Author    struct {
		ID            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
	} `json:"author"`
}

type ExportedContent struct {
	ChannelID string `json:"channel_id"`
	Messages  []struct {
		Content string `json:"message"`
		UserID  string `json:"user_id"`
		User    string `json:"user"`
	} `json:"messages"`
}

const g_AUTH_FILE_NAME string = "auth.txt"
const g_API_VERSION string = "v9"
const g_API_LIMIT uint8 = 100
const g_SLEEP_TIME time.Duration = 500

var g_channelID string = "0"

func main() {
	LogSetup()

	if len(os.Args) < 2 {
		log.Fatal("pass channel id as cli arg.")
	}

	readAuthFile, err := os.ReadFile(g_AUTH_FILE_NAME)
	if err != nil || string(readAuthFile) == "" {
		os.Create(g_AUTH_FILE_NAME)
		log.Printf("paste discord token inside '%s'.", g_AUTH_FILE_NAME)
		os.Exit(-1)
	}

	// var auth string = strings.Trim(string(readAuthFile), "\n")
	var auth string = string(readAuthFile)
	g_channelID = os.Args[1]
	var baseAPIUrl string = fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d", g_API_VERSION, g_channelID, g_API_LIMIT)

	var exported ExportedContent
	exported.ChannelID = g_channelID
	var exportFile *os.File = ExportDirSetup(g_channelID)

	log.Printf("started on channel: '%s', api version: '%s', limit: '%d'", g_channelID, g_API_VERSION, g_API_LIMIT)
	Call(baseAPIUrl, auth, exportFile, exported)
}

func LogSetup() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	if _, err := os.Stat("./logs"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("./logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	logFile, err := os.Create(fmt.Sprintf("./logs/%s.log", time.Now().Format("2006-01-02_15.04.05")))
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func ExportDirSetup(g_channelID string) *os.File {
	if _, err := os.Stat("./message-exports"); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir("./message-exports", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	var exportFileName string = fmt.Sprintf("%s-%v.json", g_channelID, time.Now().Unix())
	var exportPath string = fmt.Sprintf("./message-exports/%s", exportFileName)
	exportFile, err := os.OpenFile(exportPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer exportFile.Close()

	return exportFile
}

func Call(APIUrl, auth string, exportFile *os.File, exported ExportedContent) {
	var prevMessageID string = "0"
	var client *http.Client = &http.Client{
		Timeout: 5 * time.Second,
	}
	var msgJSONList []MessageBlock

	req, err := http.NewRequest("GET", APIUrl, nil)
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
		exported.Messages = append(exported.Messages, struct {
			Content string `json:"message"`
			UserID  string `json:"user_id"`
			User    string `json:"user"`
		}{
			Content: message.Content,
			UserID:  message.Author.ID,
			User:    message.Author.Username + "#" + message.Author.Discriminator,
		})
	}

	if len(msgJSONList) < int(g_API_LIMIT) {
		export, err := json.Marshal(exported)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(exportFile.Name(), export, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("finished")
		os.Exit(0)
	}

	var beforeParamAPIUrl string = fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d&before=%s", g_API_VERSION, g_channelID, g_API_LIMIT, prevMessageID)

	time.Sleep(g_SLEEP_TIME * time.Millisecond)
	Call(beforeParamAPIUrl, auth, exportFile, exported)
}
