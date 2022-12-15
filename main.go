package main

import (
	"encoding/json"
	"fmt"
	"main/chatgpt"
	"time"

	pubnub "github.com/pubnub/go/v7"
)

type Msg struct {
	Content    string    `json:"content"`
	SenderUUID string    `json:"sender_uuid"`
	Date       time.Time `json:"date"`
	HashKey    string    `json:"$$hashKey"`
}

var pn *pubnub.PubNub
var serverID pubnub.UserId = "AxonifAI"

func init() {
	config := pubnub.NewConfigWithUserId("")
	config.SubscribeKey = ""
	config.PublishKey = ""
	config.SetUserId(serverID)

	chatgpt.CreateChatGPTClient()

	pn = pubnub.NewPubNub(config)
}

func main() {
	listener := pubnub.NewListener()
	doneSubscribe := make(chan bool)

	go messageListener(listener)

	pn.Subscribe().
		Channels([]string{"AdamJoshWarren"}).
		Execute()

	<-doneSubscribe

	pn.Unsubscribe().
		Channels([]string{"AdamJoshWarren"}).
		Execute()
}

func messageListener(listener *pubnub.Listener) {
	pn.AddListener(listener)

	for {
		select {
		case status := <-listener.Status:
			switch status.Category {
			case pubnub.PNConnectedCategory:
				fmt.Println("Server has started listening.")
			}
		case message := <-listener.Message:
			if msg, ok := message.Message.(map[string]interface{}); ok {
				msgObj := formatMessage(msg)
				if msgObj.SenderUUID == "" || msgObj.SenderUUID == string(serverID) {
					continue //ignore my own messages
				}
				fmt.Printf("Incoming message from %s\nMessage:%s \n", msgObj.SenderUUID, msgObj.Content)
				sendResponse(askChatGPT(msgObj.Content))
			}
		case <-listener.Presence:
		}
	}
}

func formatMessage(incomingMsg interface{}) Msg {
	jsonStr, err := json.Marshal(incomingMsg)
	if err != nil {
		fmt.Println(err)
	}

	var msg Msg
	if err := json.Unmarshal(jsonStr, &msg); err != nil {
		fmt.Println(err)
	}
	return msg
}

func askChatGPT(question string) string {
	response, err := chatgpt.AskChatGPT(question)

	if err != nil {
		fmt.Println(err)
		return "Sorry. Could not process your question."
	}

	return response
}

func sendResponse(message string) {
	if message == "" {
		return
	}
	msgResp := Msg{
		Content:    message,
		SenderUUID: string(serverID),
		Date:       time.Now(),
	}
	jsonStr, err := json.Marshal(msgResp)
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[string]interface{})

	if err := json.Unmarshal(jsonStr, &data); err != nil {
		fmt.Println(err)
	}

	data["content"] = message
	_, _, _ = pn.Publish().
		Channel("AdamJoshWarren").
		Message(data).
		Execute()
}
