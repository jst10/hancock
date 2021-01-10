package database

import (
	"made.by.jst10/outfit7/hancock/cmd/constants"
)

var cacheVersionId = -1
var cache [][][][]SdkPerformance

func savePerformancesInCache(
	versionId int,
	countries []Country,
	apps []App,
	performances []Performance) {

	numberOfTypes := len(constants.AddTypes)
	numberOfCountries := len(countries)
	numberOfApps := len(apps)

	tmpCache := make([][][][]SdkPerformance, numberOfTypes)

	for _, performance := range performances {
		adTypeId := performance.AdType
		countryId := performance.Country
		appId := performance.App
		sdkId := performance.Sdk
		score := performance.Score

		if tmpCache[adTypeId] == nil {
			tmpCache[adTypeId] = make([][][]SdkPerformance, numberOfCountries)
		}
		if tmpCache[adTypeId][countryId] == nil {
			tmpCache[adTypeId][countryId] = make([][]SdkPerformance, numberOfApps)
		}

		//It is not necessary that I have same number of SDK, because data is already preprocessed
		//eg. Removed FB from CN country
		if tmpCache[adTypeId][countryId][appId] == nil {
			tmpCache[adTypeId][countryId][appId] = make([]SdkPerformance, 0)
		}
		//	Data received into method are already sorted by score, so nothing more is needed to do here
		tmpCache[adTypeId][countryId][appId] = append(tmpCache[adTypeId][countryId][appId], SdkPerformance{Sdk: uint8(sdkId), Score: uint8(score)})
	}
	cacheVersionId = versionId
	cache = tmpCache
}

func getSdksFromCache(adTypeId int, countryId int, appId int) []SdkPerformance {
	if cacheVersionId != -1 {
		return cache[adTypeId][countryId][appId]
	} else {
		return nil
	}
}



func invalidateCache() {
	cacheVersionId = -1
	cache = nil
}
