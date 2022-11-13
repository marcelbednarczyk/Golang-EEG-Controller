package main

import (
	"encoding/json"
	"fmt"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
	"golang.org/x/net/websocket"
)

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
