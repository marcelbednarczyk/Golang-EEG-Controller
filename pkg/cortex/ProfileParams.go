package cortex

type SetupProfileParams struct {
	CortexToken string `json:"cortexToken"`
	Headset     string `json:"headset"`
	Profile     string `json:"profile"`
	Status      string `json:"status"`
}

type GetProfileParams struct {
	CortexToken string `json:"cortexToken"`
	Headset     string `json:"headset"`
}
