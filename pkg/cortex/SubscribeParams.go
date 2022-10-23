package cortex

type SubscribeParams struct {
	CortexToken string   `json:"cortexToken"`
	Session     string   `json:"session"`
	Streams     []string `json:"streams"`
}
