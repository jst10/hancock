package structs

type Performance struct {
	AdType  string `json:"ad_type"`
	Country string `json:"country"`
	App     string `json:"app"`
	Sdk     string `json:"sdk"`
	Score   int    `json:"score"`
}

func NewPerformance(adType, country, app, sdk string, score int) *Performance {
	ap := Performance{
		AdType:  adType,
		Country: country,
		App:     app,
		Sdk:     sdk,
		Score:   score,
	}
	return &ap
}
