package cortex

type CreateSessionParams struct {
	CortexToken string `json:"cortexToken"`
	Status      string `json:"status"`
	Headset     string `json:"headset"`
}

type UpdateSessionParams struct {
	CortexToken string `json:"cortexToken"`
	Session     string `json:"session"`
	Status      string `json:"status"`
}
