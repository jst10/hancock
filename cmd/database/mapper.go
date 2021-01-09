package database

import "made.by.jst10/outfit7/hancock/cmd/structs"

func buildMappersFromRawData(performances []structs.Performance) *Mappers {
	countries := make([]Country, 1)
	countryNameToId := make(map[string]int)
	countryIdToName := make(map[int]string)

	apps := make([]App, 1)
	appNameToId := make(map[string]int)
	appIdToName := make(map[int]string)

	sdks := make([]Sdk, 1)
	sdkNameToId := make(map[string]int)
	sdkIdToName := make(map[int]string)

	for _, performance := range performances {
		countryName := performance.Country
		appName := performance.App
		sdkName := performance.Sdk

		_, prs := countryNameToId[countryName]
		if !prs {
			id := len(countries)
			country := Country{ID: id, Name: countryName}
			countries = append(countries, country)
			countryNameToId[country.Name] = country.ID
			countryIdToName[country.ID] = country.Name
		}

		_, prs = appNameToId[appName]
		if !prs {
			id := len(countries)
			app := App{ID: id, Name: appName}
			apps = append(apps, app)
			appNameToId[app.Name] = app.ID
			appIdToName[app.ID] = app.Name
		}

		_, prs = sdkNameToId[sdkName]
		if !prs {
			id := len(countries)
			sdk := Sdk{ID: id, Name: sdkName}
			sdks = append(sdks, sdk)
			sdkNameToId[sdk.Name] = sdk.ID
			sdkIdToName[sdk.ID] = sdk.Name
		}
	}
	return &Mappers{
		countries:       countries,
		countryIdToName: countryIdToName,
		apps:            apps,
		appNameToId:     appNameToId,
		appIdToName:     appIdToName,
		sdks:            sdks,
		sdkNameToId:     sdkNameToId,
		sdkIdToName:     sdkIdToName,
	}
}

func buildMappersFromDBData(
	countries []Country,
	apps []App,
	sdks []Sdk,
) *Mappers {

	countryNameToId := make(map[string]int)
	countryIdToName := make(map[int]string)

	appNameToId := make(map[string]int)
	appIdToName := make(map[int]string)

	sdkNameToId := make(map[string]int)
	sdkIdToName := make(map[int]string)

	for _, country := range countries {
		countryNameToId[country.Name] = country.ID
		countryIdToName[country.ID] = country.Name
	}
	for _, app := range apps {
		appNameToId[app.Name] = app.ID
		appIdToName[app.ID] = app.Name
	}
	for _, sdk := range sdks {
		sdkNameToId[sdk.Name] = sdk.ID
		sdkIdToName[sdk.ID] = sdk.Name
	}

	return &Mappers{
		countries:       countries,
		countryIdToName: countryIdToName,
		apps:            apps,
		appNameToId:     appNameToId,
		appIdToName:     appIdToName,
		sdks:            sdks,
		sdkNameToId:     sdkNameToId,
		sdkIdToName:     sdkIdToName,
	}
}
