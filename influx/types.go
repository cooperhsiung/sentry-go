package influx

type Point struct {
	Type      string      `json:"type"`
	Project   string      `json:"project"`
	Field     string      `json:"field"`
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
	Method    string      `json:"method"`
	Tag       string      `json:"tag"`
}

type Qs struct {
	Project string
	Field   string
	Period  string
	Value   string
	Offset  string
}
