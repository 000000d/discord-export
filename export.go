package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "log"
    "net/http"
    "os"
    "strings"
    "time"
)

type MessageBlock struct {
    ID        string `json:"id,omitempty"`
    Content   string `json:"content,omitempty"`
    ChannelID string `json:"channel_id,omitempty"`
    Author    struct {
        ID       string `json:"id"`
        Username string `json:"username"`
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

const AUTH_FILE_NAME string = "auth.txt"

const API_VERSION string = "v10"
const API_LIMIT uint8 = 100
const SLEEP_TIME time.Duration = 500 // Stated in milliseconds

var channelID string = "0"

func exportDirSetup(channelID string) *os.File {
    exportDirectory := "message-exports"
    if _, err := os.Stat(exportDirectory); errors.Is(err, os.ErrNotExist) {
        err := os.Mkdir(exportDirectory, 0755)
        if err != nil {
            log.Fatalf("Error creating directory '%s' to store exported messages. ERROR: %v", exportDirectory, err)
        }
    }
    exportFileName := fmt.Sprintf("%s_%v.json", channelID, time.Now().Format("2006-01-02_15.04.05"))
    exportPath := fmt.Sprintf("./%s/%s", exportDirectory, exportFileName)
    exportFile, err := os.OpenFile(exportPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("Error creating export file '%s'. ERROR: %v\n", exportFile,  err)
    }
    defer exportFile.Close()

    return exportFile
}

func apiCall(baseAPI, auth string, exportFile *os.File, exported ExportedContent) {
    prevMessageID := "0"
    client := &http.Client{
        Timeout: 5 * time.Second,
    }
    var msgJSONList []MessageBlock

    req, err := http.NewRequest("GET", baseAPI, nil)
    if err != nil {
        log.Fatalln("Error creating new HTTP request. ERROR:", err)
    }
    req.Header.Set("Authorization", auth)

    res, err := client.Do(req)
    if err != nil {
        log.Fatalln("Error sending HTTP request. ERROR:", err)
    }
    defer res.Body.Close()

    resBody, err := io.ReadAll(res.Body)
    if err != nil {
        log.Fatalln("Error reading HTTP response body. ERROR:", err)
    }

    err = json.Unmarshal(resBody, &msgJSONList)
    if err != nil {
        log.Fatalf("Incorrect auth token. Received: '%s'\n", auth)
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
            User:    message.Author.Username,
        })
    }

    if len(msgJSONList) < int(API_LIMIT) {
        export, err := json.Marshal(exported)
        if err != nil {
            log.Fatalln("Error with JSON marshal. ERROR:", err)
        }
        err = os.WriteFile(exportFile.Name(), export, os.ModePerm)
        if err != nil {
            log.Fatalf("Error writing exported messages to file '%s'. ERROR: %v\n", exportFile.Name(),  err)
        }
        log.Println("Finished")
        os.Exit(0)
    }

    beforeParamAPI := fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d&before=%s", API_VERSION, channelID, API_LIMIT, prevMessageID)

    time.Sleep(SLEEP_TIME * time.Millisecond)
    apiCall(beforeParamAPI, auth, exportFile, exported)
}

func main() {
    if len(os.Args) < 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s <channel_id> ...\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "ERROR: Include channel ID as CLI arg\n")
        os.Exit(1);
    }

    readAuthFile, err := os.ReadFile(AUTH_FILE_NAME)
    if err != nil || len(readAuthFile) == 0 {
        os.Create(AUTH_FILE_NAME)
        fmt.Fprintf(os.Stderr, "File '%s' must contain Discord token to export messages\n", AUTH_FILE_NAME)
        os.Exit(1)
    }

    auth := strings.Trim(string(readAuthFile), "\n")
    channelID = os.Args[1]
    baseAPI := fmt.Sprintf("https://discord.com/api/%s/channels/%s/messages?limit=%d", API_VERSION, channelID, API_LIMIT)

    var exported ExportedContent
    exported.ChannelID = channelID
    exportFile := exportDirSetup(channelID)

    log.Printf("Started on channel: '%s', API version: '%s', limit: '%d'", channelID, API_VERSION, API_LIMIT)
    apiCall(baseAPI, auth, exportFile, exported)
}
