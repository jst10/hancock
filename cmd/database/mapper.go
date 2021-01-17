package database

import (
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

func insertCodeListInMapperIfNot(codeListMapper *CodeListMapper, name string) {
	_, prs := codeListMapper.itemNameToId[name]
	if !prs {
		id := len(codeListMapper.items)
		item := CodeList{ID: id, Name: name}
		codeListMapper.items = append(codeListMapper.items, &item)
		codeListMapper.itemNameToId[item.Name] = item.ID
		codeListMapper.itemIdToName[item.ID] = item.Name
	}
}

func buildMappersFromRawData(performances []*structs.Performance) *Mappers {
	countryMapper := NewCodeListMapper()
	appMapper := NewCodeListMapper()
	sdkMapper := NewCodeListMapper()

	for _, performance := range performances {
		insertCodeListInMapperIfNot(countryMapper, performance.Country)
		insertCodeListInMapperIfNot(appMapper, performance.App)
		insertCodeListInMapperIfNot(sdkMapper, performance.Sdk)
	}
	return &Mappers{
		countryMapper: countryMapper,
		appMapper:     appMapper,
		sdkMapper:     sdkMapper,
	}
}

func buildCodeListMapperFromDbData(items []*CodeList) *CodeListMapper {
	itemNameToId := make(map[string]int)
	itemIdToName := make(map[int]string)
	for _, item := range items {
		itemNameToId[item.Name] = item.ID
		itemIdToName[item.ID] = item.Name
	}
	return &CodeListMapper{
		items:        items,
		itemNameToId: itemNameToId,
		itemIdToName: itemIdToName,
	}
}

func buildMappersFromDBData(
	countries []*CodeList,
	apps []*CodeList,
	sdks []*CodeList,
) *Mappers {

	return &Mappers{
		countryMapper: buildCodeListMapperFromDbData(countries),
		appMapper:     buildCodeListMapperFromDbData(apps),
		sdkMapper:     buildCodeListMapperFromDbData(sdks),
	}
}
