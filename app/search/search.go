package search

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	//"time"
)

const (
	url = "https://api.t-a-a-s.ru/client"
)

// RequestInline ...
type RequestInline struct {
	apiKey        string
	typeReq       string
	chatID        string
	limit         string
	offset        string
	fromMessageID string
}

type respGetChatHistory struct {
	TotalCount int       `json:"total_count"`
	Messages   []Message `json:"messages"`
}

// Message ...
type Message struct {
	ID      int     `json:"id"`
	Content Content `json:"content"`
}

// Content ...
type Content struct {
	VoiceNote VoiceNote `json:"voice_note"`
	Caption   Caption   `json:"caption"`
}

// VoiceNote ...
type VoiceNote struct {
	Voice File `json:"voice"`
}

// File ...
type File struct {
	FileID int        `json:"id"`
	Size   int        `json:"size"`
	Remote RemoteFile `json:"remote"`
}

// RemoteFile ...
type RemoteFile struct {
	ID       string `json:"id"`
	UniqueID string `json:"unique_id"`
}

// Caption ...
type Caption struct {
	Text string `json:"text"`
}

// GetRequestInline ...
func GetRequestInline(q RequestInline) ([]Message, int) {
	reqBody, err := json.Marshal(map[string]string{
		"api_key":         q.apiKey,
		"@type":           q.typeReq,
		"chat_id":         q.chatID,
		"limit":           q.limit,
		"offset":          q.offset,
		"from_message_id": q.fromMessageID,
	})
	if err != nil {
		log.Print("Marshal", err)
	}

	client := &http.Client{}

	// _, err = client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	// if err != nil {
	// 	log.Print("Post", err)
	// }
	// time.Sleep(time.Second*1)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log.Print("Post", err)
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ReadAll", err)
	}

	responce := new(respGetChatHistory)
	if err := json.Unmarshal(b, responce); err != nil {
		log.Print("Unmarshal", err)
	}

	return responce.Messages, responce.TotalCount

}

// SetRequestInline ...
func (r *RequestInline) SetRequestInline(apikey, typereq, cgatid, limit, offset, frommessageid string) {
	r.apiKey = apikey
	r.typeReq = typereq
	r.chatID = cgatid
	r.limit = limit
	r.offset = offset
	r.fromMessageID = frommessageid
}

// GetMap ...
func GetMap(messages []Message) map[string]string {
	m := make(map[string]string)
	for _, message := range messages {
		if message.Content.VoiceNote.Voice.Size != 0 {
			m[message.Content.VoiceNote.Voice.Remote.ID] = message.Content.Caption.Text
		}
	}

	return m
}

// 	} else {
// 		return []tgbotapi.InlineQueryResultAudio{
// 			{
// 				Type:    "article",
// 				ID:      "1",
// 				Title:   "error",
// 				Caption: "error",
// 				//InputMessageContent: "error",
// 			},
// 		}
// 	}
// }
