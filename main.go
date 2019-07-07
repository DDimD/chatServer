package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// // Message struct for sending message to clients
// type Message struct {
// 	userName string `json:"userName"`
// 	body     string `json:"messageBody"`
// }

// type Client struct {
// 	name      string
// 	webSocket *websocket.Conn
// }

func clientSendMsgHandler(respWriter http.ResponseWriter, request *http.Request) {
	message, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(fmt.Sprintf("sendClient request read error"))
	}
	messageString := string(message)
	messageString += "+1"
	message = []byte(messageString)
	respWriter.Write(message)
}

func main() {
	fileSys := http.FileServer(http.Dir("index/"))
	http.Handle("/", http.StripPrefix("/", fileSys))
	http.HandleFunc("/sendMessage", clientSendMsgHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
