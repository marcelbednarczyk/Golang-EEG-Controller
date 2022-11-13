package main

import (
	"log"
	"os"

	"golang.org/x/net/websocket"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
)

func main() {
	println("Welcome to EEG with Golang!")
	defer println("Goodbye!")

	origin := "http://localhost/"
	url := "wss://localhost:6868"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println(err)
		return
	}

	cortexToken, headsetID, err := connectHeadset(ws)
	if err != nil {
		log.Println(err)
		return
	}

	resp, err := apiCall[cortex.Response](ws, cortex.Request{
		ID:      6,
		JsonRPC: "2.0",
		Method:  "createSession",
		Params: cortex.CreateSessionParams{
			CortexToken: cortexToken,
			Status:      "open",
			Headset:     headsetID,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	sessionID := resp.Result.(map[string]interface{})["id"].(string)

	defer func() {
		_, err = apiCall[cortex.Response](ws, cortex.Request{
			ID:      7,
			JsonRPC: "2.0",
			Method:  "updateSession",
			Params: cortex.UpdateSessionParams{
				CortexToken: cortexToken,
				Session:     sessionID,
				Status:      "close",
			},
		})
		if err != nil {
			log.Println(err)
			return
		}
	}()

	switch os.Getenv("AIM") {
	case "TRAIN":
		trainMental(ws, "push", cortexToken, sessionID, headsetID)
	case "CONTROL":
		controlMental(ws, cortexToken, sessionID)
	default:
		println("No aim specified. Exiting...")
	}
}
