package cortex

type DataSample struct {
	Com  []interface{} `json:"com"`
	SID  string        `json:"sid"`
	Time float32       `json:"time"`
}
