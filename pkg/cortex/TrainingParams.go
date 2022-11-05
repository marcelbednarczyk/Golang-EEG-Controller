package cortex

type TrainingParams struct {
	CortexToken string `json:"cortexToken"`
	Session     string `json:"session"`
	Detection   string `json:"detection"`
	Status      string `json:"status"`
	Action      string `json:"action"`
}
