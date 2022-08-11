package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Update is a Tg object that the handler receives every time a user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Message is a Tg object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Chat is s Tg object that indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// parseTgRequest handles incoming update from the Tg web hook
func parseTgRequest(r *http.Request) (*Update, error) {
	var upd Update
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		log.Printf("could not decode incoming update %s", err)
		return nil, err
	}
	return &upd, nil
}

// HandleTgWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func handleTgWebHook(_ http.ResponseWriter, r *http.Request, updHandler func(upd string) string) {
	// Parse incoming request
	upd, err := parseTgRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err)
		return
	}
	responseText := updHandler(upd.Message.Text)
	responseBody, err := sendTextToChat(upd.Message.Chat.Id, responseText)
	if err != nil {
		log.Printf("got error %s from telegram, response body is %s", err, responseBody)
	} else {
		log.Printf("responseText %s succesfully distributed to chat id %d", responseText, upd.Message.Chat.Id)
	}

}

func CreateWebHookHandler(updHandler func(upd string) string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handleTgWebHook(w, r, updHandler)
	}
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat id
func sendTextToChat(id int, text string) (string, error) {
	log.Printf("Sending %s to chat_id: %d", text, id)
	response, err := http.PostForm(
		"https://api.telegram.org/bot"+os.Getenv("TELEGRAM_BOT_TOKEN")+"/sendMessage",
		url.Values{
			"chat_id": {strconv.Itoa(id)},
			"text":    {text},
		})
	if err != nil {
		log.Printf("error when posting text to thee chat: %s", err)
		return "", err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error in parsing telegram answer %s", err)
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)
	return bodyString, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func LoadEnv(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		keyAndVal := strings.Split(sc.Text(), "=")
		os.Setenv(keyAndVal[0], keyAndVal[1])
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	LoadEnv(".env")
	http.HandleFunc("/", handler)
	http.HandleFunc("/echo", CreateWebHookHandler(func(s string) string {
		return s
	}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
