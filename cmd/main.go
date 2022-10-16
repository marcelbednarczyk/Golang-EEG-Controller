package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/websocket"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
)

func main() {
	println("Welcome to EEG with Golang!")

	origin := "http://localhost/"
	url := "wss://localhost:6868"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}

	_, err = apiCall[cortex.Response](ws, cortex.GetDefaultInfoRequest())
	if err != nil {
		log.Fatal(err)
	}

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      1,
		JsonRPC: "2.0",
		Method:  "requestAccess",
		Params: cortex.AuthParams{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      2,
		JsonRPC: "2.0",
		Method:  "queryHeadsets",
		Params:  nil,
	})
	if err != nil {
		log.Fatal(err)
	}

	resp, err := apiCall[cortex.Response](ws, cortex.Request{
		ID:      3,
		JsonRPC: "2.0",
		Method:  "authorize",
		Params: cortex.AuthParams{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	cortexToken := resp.Result.(map[string]interface{})["cortexToken"].(string)
	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      4,
		JsonRPC: "2.0",
		Method:  "createSession",
		Params: cortex.SessionParams{
			CortexToken: cortexToken,
			Headset:     "EPOCX-E2020839",
			Status:      "open",
		},
	})
	if err != nil {
		log.Fatal(err)
	}

}

func apiCall[T any](ws *websocket.Conn, request cortex.Request) (*T, error) {
	req, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	if _, err = ws.Write(req); err != nil {
		return nil, err
	}

	var msg = make([]byte, 2048)
	var n int
	if n, err = ws.Read(msg); err != nil {
		return nil, err
	}

	fmt.Printf("Received: %s.\n", msg[:n])
	var result *T
	if err = json.Unmarshal(msg[:n], &result); err != nil {
		return nil, err
	}

	return result, nil
}
