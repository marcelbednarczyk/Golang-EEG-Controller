package main

import (
	"log"
	"os"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
	"golang.org/x/net/websocket"
)

func trainMental(ws *websocket.Conn, action, cortexToken, sessionID, headsetID string) {
	_, err := apiCall[cortex.Response](ws, cortex.Request{
		ID:      8,
		JsonRPC: "2.0",
		Method:  "subscribe",
		Params: cortex.SubscribeParams{
			CortexToken: cortexToken,
			Session:     sessionID,
			Streams:     []string{"sys"},
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		_, err = apiCall[cortex.Response](ws, cortex.Request{
			ID:      9,
			JsonRPC: "2.0",
			Method:  "unsubscribe",
			Params: cortex.SubscribeParams{
				CortexToken: cortexToken,
				Session:     sessionID,
				Streams:     []string{"sys"},
			},
		})
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      10,
		JsonRPC: "2.0",
		Method:  "training",
		Params: cortex.TrainingParams{
			CortexToken: cortexToken,
			Session:     sessionID,
			Detection:   "mentalCommand",
			Status:      "start",
			Action:      "push",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// training started
	receivePrint(ws)
	// training ended
	receivePrint(ws)

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      11,
		JsonRPC: "2.0",
		Method:  "training",
		Params: cortex.TrainingParams{
			CortexToken: cortexToken,
			Session:     sessionID,
			Detection:   "mentalCommand",
			Status:      "accept",
			Action:      "push",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
	receivePrint(ws)

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      12,
		JsonRPC: "2.0",
		Method:  "setupProfile",
		Params: cortex.SetupProfileParams{
			CortexToken: cortexToken,
			Headset:     headsetID,
			Profile:     os.Getenv("PROFILE_NAME"),
			Status:      "save",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}
}
