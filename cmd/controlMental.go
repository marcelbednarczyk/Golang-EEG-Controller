package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	client := http.Client{}

	for i := 0; i < 500; i++ {
		data, err := receive[cortex.DataSample](ws)
		if err != nil {
			return err
		}

		for i := range data.Com {
			if i+1 < len(data.Com) {
				fmt.Printf("Action: %s\t Power: %f\n", data.Com[i], data.Com[i+1])
				s := 0.0
				if s, err = strconv.ParseFloat(fmt.Sprintf("%f", data.Com[i+1]), 32); err != nil {
					fmt.Println("Error: ", err)
					return err
				}
				if data.Com[i] == "lift" && s > 0.3 {
					fmt.Println("Fireplace!")
					httpGet(&client, os.Getenv("IOT_URL")+"/update", map[string]string{"state": "0"})
				} else if data.Com[i] == "push" && s > 0.3 {
					fmt.Println("Snow!")
					httpGet(&client, os.Getenv("IOT_URL")+"/update", map[string]string{"state": "1"})
				}
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
		fmt.Println(err)
	}

	return nil
}
