package database

type CodeList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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
	DbIndex   int    `json:"db_index"`
}

type CodeListMapper struct {
	items        []*CodeList
	itemNameToId map[string]int
	itemIdToName map[int]string
}

func NewCodeListMapper() *CodeListMapper {
	codeListMapper := CodeListMapper{
		items:        make([]*CodeList, 0),
		itemNameToId: make(map[string]int),
		itemIdToName: make(map[int]string),
	}
	return &codeListMapper
}

type Mappers struct {
	countryMapper *CodeListMapper
	appMapper     *CodeListMapper
	sdkMapper     *CodeListMapper
}
