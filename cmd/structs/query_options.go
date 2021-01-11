package structs

type QueryOptions struct {
	Country    string `json:"country"`
	Platform   string `json:"platform"`
	OsVersion  string `json:"os_version"`
	AppName    string `json:"app_name"`
	AppVersion string `json:"app_version"`
}

func NewQueryOptions(country, platform, osVersion, appName, appVersion string) *QueryOptions {
	qO := QueryOptions{
		Country:    country,
		Platform:   platform,
		OsVersion:  osVersion,
		AppName:    appName,
		AppVersion: appVersion,
	}
	return &qO
}
