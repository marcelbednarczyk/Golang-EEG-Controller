package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/websocket"
)

func receivePrint(ws *websocket.Conn) error {
	var msg = make([]byte, 2048)
	var n int
	var err error
	if n, err = ws.Read(msg); err != nil {
		return err
	}

	fmt.Printf("Received: %s.\n", msg[:n])
	return nil
}

func receive[T any](ws *websocket.Conn) (*T, error) {
	var msg = make([]byte, 2048)
	var n int
	var err error
	if n, err = ws.Read(msg); err != nil {
		return nil, err
	}

	// fmt.Printf("Received: %s.\n", msg[:n])
	var result *T
	if err = json.Unmarshal(msg[:n], &result); err != nil {
		return nil, err
	}

	return result, nil
}
