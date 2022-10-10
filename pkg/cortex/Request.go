package cortex

type Request struct {
	ID      int         `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

func GetDefaultInfoRequest() Request {
	return Request{
		ID:      1,
		JsonRPC: "2.0",
		Method:  "getCortexInfo",
		Params:  nil,
	}
}
