package database

type App struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Country struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Sdk struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Performance struct {
	ID      int `json:"id"`
	AdType  int `json:"ad_type"`
	Country int `json:"country"`
	App     int `json:"app"`
	Sdk     int `json:"sdk"`
	Score   int `json:"score"`
}

func NewPerformance() *Performance {
	ap := Performance{}
	return &ap
}

type SdkPerformance struct {
	Sdk   uint8
	Score uint8
}

type Version struct {
	ID        int    `json:"id"`
	CreatedAt string `json:"created_at"`
	DbIndex   string `json:"db_index"`
}

type Mappers struct {
	countries       []Country
	countryNameToId map[string]int
	countryIdToName map[int]string

	apps        []App
	appNameToId map[string]int
	appIdToName map[int]string

	sdks        []Sdk
	sdkNameToId map[string]int
	sdkIdToName map[int]string
}
