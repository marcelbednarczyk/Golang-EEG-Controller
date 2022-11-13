package main

import (
	"log"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
	"golang.org/x/net/websocket"
)

func controlMental(ws *websocket.Conn, cortexToken, sessionID string) error {
	_, err := apiCall[cortex.Response](ws, cortex.Request{
		ID:      8,
		JsonRPC: "2.0",
		Method:  "subscribe",
		Params: cortex.SubscribeParams{
			CortexToken: cortexToken,
			Session:     sessionID,
			Streams:     []string{"com"},
		},
	})
	if err != nil {
		return err
	}

	for i := 0; i < 100; i++ {
		data, err := receive[cortex.DataSample](ws)
		if err != nil {
			return err
		}

		for i := range data.Com {
			if i+1 < len(data.Com) {
				log.Printf("Action: %s\t Power: %f\n", data.Com[i], data.Com[i+1])
			}
		}
	}

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      9,
		JsonRPC: "2.0",
		Method:  "unsubscribe",
		Params: cortex.SubscribeParams{
			CortexToken: cortexToken,
			Session:     sessionID,
			Streams:     []string{"com"},
		},
	})
	if err != nil {
		log.Println(err)
	}

	return nil
}
