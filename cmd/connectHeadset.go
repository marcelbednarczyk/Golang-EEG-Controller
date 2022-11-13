package main

import (
	"os"

	"github.com/marcelbednarczyk/Golang-EEG-Controller/pkg/cortex"
	"golang.org/x/net/websocket"
)

func connectHeadset(ws *websocket.Conn) (string, string, error) {
	_, err := apiCall[cortex.Response](ws, cortex.GetDefaultInfoRequest())
	if err != nil {
		return "", "", err
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
		return "", "", err
	}

	headsetResp, err := apiCall[cortex.ResponseSlice](ws, cortex.Request{
		ID:      2,
		JsonRPC: "2.0",
		Method:  "queryHeadsets",
		Params:  nil,
	})
	if err != nil {
		return "", "", err
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
		return "", "", err
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
		return "", "", err
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
		return "", "", err
	}

	return cortexToken, headsetID, nil
}
