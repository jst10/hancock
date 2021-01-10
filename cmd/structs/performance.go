package structs

type Performance struct {
	AdType  string `json:"ad_type"`
	Country string `json:"country"`
	App     string `json:"app"`
	Sdk     string `json:"sdk"`
	Score   int    `json:"score"`
}

func NewPerformance() *Performance {
	ap := Performance{}
	return &ap
}
