package cortex

type AuthResult struct {
	CortexToken string      `json:"cortexToken"`
	Warning     interface{} `json:"warning"`
}
