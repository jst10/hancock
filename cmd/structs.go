package main

type QueryOptions struct {
	countryCode string
	platform    string
	osVersion   string
	appName     string
	appVersion  string
}

func newQueryOptions() *QueryOptions {
	qO := QueryOptions{}
	return &qO
}

type AdPerformance struct {
	CountryCode string	`json:"country_code"`
	Network     string	`json:"network"`
	AddType     string	`json:"add_type"`
	Score       int	`json:"score"`
}

func newAdPerformance() *AdPerformance {
	ap := AdPerformance{}
	return &ap
}


type AdPerformanceResponse struct{
	banner []AdPerformance
	interstitial []AdPerformance
	reward []AdPerformance
}

type HttpError struct {
	Details string `json:"details"`
}

func newHttpError(details string) *HttpError {
	he := HttpError{Details: details}
	return &he
}
