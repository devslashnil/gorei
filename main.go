package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Chat is s Telegram object that indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var upd Update
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	return &upd, nil
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {

}

func main() {

}
