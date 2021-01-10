package structs

type PerformanceResponse struct {
	Country string `json:"country"`
	App     string `json:"app"`
	Banner       []SdkScore `json:"banner"`
	Interstitial []SdkScore `json:"interstitial"`
	Reward       []SdkScore `json:"reward"`
}
