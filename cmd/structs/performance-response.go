package structs

type PerformanceResponse struct {
	Banner       []Performance `json:"banner"`
	Interstitial []Performance `json:"interstitial"`
	Reward       []Performance `json:"reward"`
}
