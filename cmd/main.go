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
	defer println("Goodbye!")

	origin := "http://localhost/"
	url := "wss://localhost:6868"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Println(err)
		return
	}

	_, err = apiCall[cortex.Response](ws, cortex.GetDefaultInfoRequest())
	if err != nil {
		log.Println(err)
		return
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
		log.Println(err)
		return
	}

	headsetResp, err := apiCall[cortex.ResponseSlice](ws, cortex.Request{
		ID:      2,
		JsonRPC: "2.0",
		Method:  "queryHeadsets",
		Params:  nil,
	})
	if err != nil {
		log.Println(err)
		return
	}
	headsetID := headsetResp.Result[0].(map[string]interface{})["id"].(string)

	resp, err := apiCall[cortex.Response](ws, cortex.Request{
		ID:      3,
		JsonRPC: "2.0",
		Method:  "authorize",
		Params: cortex.AuthParams{
			ClientID:     os.Getenv("CLIENT_ID"),
			ClientSecret: os.Getenv("CLIENT_SECRET"),
			Debit:        1,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	cortexToken := resp.Result.(map[string]interface{})["cortexToken"].(string)
	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      4,
		JsonRPC: "2.0",
		Method:  "setupProfile",
		Params: cortex.SetupProfileParams{
			CortexToken: cortexToken,
			Headset:     headsetID,
			Profile:     os.Getenv("PROFILE_NAME"),
			Status:      "load",
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	_, err = apiCall[cortex.Response](ws, cortex.Request{
		ID:      5,
		JsonRPC: "2.0",
		Method:  "getCurrentProfile",
		Params: cortex.GetProfileParams{
			CortexToken: cortexToken,
			Headset:     headsetID,
		},
	})
	if err != nil {
		log.Println(err)
		return
	}

	resp, err = apiCall[cortex.Response](ws, cortex.Request{
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

	_, err = apiCall[cortex.Response](ws, cortex.Request{
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
	receive(ws)
	// training ended
	receive(ws)

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
	receive(ws)

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

	// _, err = apiCall[cortex.Response](ws, cortex.Request{
	// 	ID:      8,
	// 	JsonRPC: "2.0",
	// 	Method:  "subscribe",
	// 	Params: cortex.SubscribeParams{
	// 		CortexToken: cortexToken,
	// 		Session:     sessionID,
	// 		Streams:     []string{"met", "com"},
	// 	},
	// })
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	// defer func() {
	// 	_, err = apiCall[cortex.Response](ws, cortex.Request{
	// 		ID:      9,
	// 		JsonRPC: "2.0",
	// 		Method:  "unsubscribe",
	// 		Params: cortex.SubscribeParams{
	// 			CortexToken: cortexToken,
	// 			Session:     sessionID,
	// 			Streams:     []string{"met", "com"},
	// 		},
	// 	})
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// }()

	// for i := 0; i < 10; i++ {
	// 	var msg = make([]byte, 2048)
	// 	var n int
	// 	if n, err = ws.Read(msg); err != nil {
	// 		log.Println("Error reading:", err)
	// 		return
	// 	}

	// 	fmt.Printf("Received: %s.\n", msg[:n])
	// }
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

func receive(ws *websocket.Conn) error {
	var msg = make([]byte, 2048)
	var n int
	var err error
	if n, err = ws.Read(msg); err != nil {
		return err
	}

	fmt.Printf("Received: %s.\n", msg[:n])
	return nil
}
