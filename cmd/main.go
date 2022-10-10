package main

import (
	"encoding/json"
	"fmt"
	"log"

	"golang.org/x/net/websocket"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
)

func main() {
	println("Hello, World!")

	origin := "http://localhost/"
	url := "wss://localhost:6868"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	req, err := json.Marshal(cortex.GetDefaultInfoRequest())
	if err != nil {
		log.Fatal(err)
	}

	if _, err := ws.Write(req); err != nil {
		log.Fatal(err)
	}
	var msg = make([]byte, 512)
	var n int
	if n, err = ws.Read(msg); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Received: %s.\n", msg[:n])
}
