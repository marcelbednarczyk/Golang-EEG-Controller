package cortex

type Response struct {
	ID      int         `json:"id"`
	JsonRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   interface{} `json:"error"`
}
