package api

import (
	"made.by.jst10/outfit7/hancock/cmd/constants"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

// TODO just for debugging, so it may be one or two for loops more than needed
// in production I would trust internal services that would provide me data

func arePerformancesValid(performances []structs.Performance) bool {
	// 1. we need to have all 3 types
	// 2. add each type we need to have data for all countries
	// 3. add each country we need to have data for each app name
	adTypes := make(map[string]bool)
	countries := make(map[string]bool)
	apps := make(map[string]bool)
	adTypeToCountries := make(map[string]map[string]bool)
	adTypeCountryToApps := make(map[string]map[string]bool)

	for _, performance := range performances {
		adTypes[performance.AdType] = true
		countries[performance.Country] = true
		apps[performance.App] = true
		adTypeToCountryDict, exist := adTypeToCountries[performance.AdType]
		if !exist {
			adTypeToCountries[performance.AdType] = make(map[string]bool)
			adTypeToCountryDict = adTypeToCountries[performance.AdType]
		}
		adTypeToCountryDict[performance.Country] = true

		key := performance.AdType + performance.Country
		adTypeCountryToAppDict, exist := adTypeCountryToApps[key]
		if !exist {
			adTypeCountryToApps[key] = make(map[string]bool)
			adTypeCountryToAppDict = adTypeCountryToApps[key]
		}
		adTypeCountryToAppDict[performance.App] = true
	}
	// since we have collected all different items, checking the length is ok
	if len(constants.AddTypes) != len(adTypes) {
		return false
	}
	for _, adType := range constants.AddTypes {
		_, exist := adTypes[adType]
		if !exist {
			return false
		}
	}

	numberOfCountries := len(countries)
	for _, v := range adTypeToCountries {
		tmp := len(v)
		if tmp != numberOfCountries {
			return false
		}
	}

	numberOfApps := len(apps)
	for _, v := range adTypeCountryToApps {
		tmp := len(v)
		if tmp != numberOfApps {
			return false
		}
	}
	return true
}
