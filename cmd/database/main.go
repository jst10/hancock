package database

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

var db *sql.DB
var mappersVersionId = -1
var mappers *Mappers

func createTablesIfNot() *custom_errors.CustomError {
	err := dbUserCreateTableIfNot()
	if err != nil {
		return err
	}
	err = dbVersionCreateTableIfNot()
	if err != nil {
		return err
	}
	err = dbSessionCreateTableIfNot()
	if err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		err = dbSdkCreateTablesIfNot(i)
		if err != nil {
			return err
		}
		err = dbAppCreateTablesIfNot(i)
		if err != nil {
			return err
		}
		err = dbCountryCreateTablesIfNot(i)
		if err != nil {
			return err
		}
		err = dbPerformanceCreateTablesIfNot(i)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitDatabase() *custom_errors.CustomError {
	var err *custom_errors.CustomError
	db, err = sqlOpen("mysql",
		"root:root@tcp(127.0.0.1:3306)/hancock")
	if err != nil {
		return err
	}
	err = dbPing()
	if err != nil {
		return err
	}
	fmt.Println("DB connection has successfully established")
	err = createTablesIfNot()
	if err != nil {
		return err
	}
	if err == nil {
		fmt.Println("DB  initialized")
	}
	return nil
	// TODO figure it out when is smart to open close db connection (definitely if we implement push/pull system)
	//defer db.Close()
}

func reloadMappersFromDb(version *Version) *custom_errors.CustomError {
	countries, countryErr := dbCountryAll(version.DbIndex)
	if countryErr != nil {
		return countryErr.AST("reload mappers from db")
	}
	apps, appErr := dbAppAll(version.DbIndex)
	if appErr != nil {
		return appErr.AST("reload mappers from db")
	}
	sdks, sdkErr := dbSdkAll(version.DbIndex)
	if sdkErr != nil {
		return sdkErr.AST("reload mappers from db")
	}
	mappers = buildMappersFromDBData(countries, apps, sdks)
	mappersVersionId = version.ID
	return nil

}

func reloadCacheFromDb(version *Version) *custom_errors.CustomError {
	performances, err := dbPerformanceAll(version.DbIndex)
	if err != nil {
		return err.AST("reload cache from db")
	}
	savePerformancesInCache(version.ID, mappers.countries, mappers.apps, performances)
	return nil
}

func GetPerformances(options *structs.QueryOptions) (*structs.PerformanceResponse, *custom_errors.CustomError) {
	latestVersion, err := dbVersionGetLatest()
	if err != nil {
		return nil, err.AST("get performances")
	}

	if latestVersion.ID != mappersVersionId {
		err := reloadMappersFromDb(latestVersion)
		if err != nil {
			return nil, err.AST("get performances")
		}
	}

	if latestVersion.ID != cacheVersionId {
		err := reloadCacheFromDb(latestVersion)
		if err != nil {
			return nil, err.AST("get performances")
		}
	}

	countryID, existCountryName := mappers.countryNameToId[options.Country]
	appID, existAppName := mappers.appNameToId[options.AppName]
	if !existAppName {
		return nil, custom_errors.GetNotValidDataError("App name is not valid")
	}
	if !existCountryName {
		return nil, custom_errors.GetNotValidDataError("Country name is not valid")
	}

	bannerPerformances := getSdksFromCache(constants.AdTypeBannerId, countryID, appID)
	interstitialPerformances := getSdksFromCache(constants.AdTypeInterstitialId, countryID, appID)
	rewardedPerformances := getSdksFromCache(constants.AdTypeRewardedId, countryID, appID)

	bannerR := make([]structs.Performance, len(bannerPerformances))
	interstitialR := make([]structs.Performance, len(interstitialPerformances))
	rewardedR := make([]structs.Performance, len(rewardedPerformances))

	for index, performance := range bannerPerformances {
		bannerR[index] = structs.Performance{
			AdType:  constants.AdTypeBanner,
			Country: options.Country,
			App:     options.AppName,
			Sdk:     mappers.sdkIdToName[int(performance.Sdk)],
			Score:   int(performance.Score),
		}
	}

	for index, performance := range interstitialPerformances {
		interstitialR[index] = structs.Performance{
			AdType:  constants.AdTypeInterstitial,
			Country: options.Country,
			App:     options.AppName,
			Sdk:     mappers.sdkIdToName[int(performance.Sdk)],
			Score:   int(performance.Score),
		}
	}

	for index, performance := range rewardedPerformances {
		rewardedR[index] = structs.Performance{
			AdType:  constants.AdTypeInterstitial,
			Country: options.Country,
			App:     options.AppName,
			Sdk:     mappers.sdkIdToName[int(performance.Sdk)],
			Score:   int(performance.Score),
		}
	}

	return &structs.PerformanceResponse{
		Banner:       bannerR,
		Interstitial: interstitialR,
		Reward:       rewardedR,
	}, nil

}
func StorePerformances(performances []structs.Performance) (*Version, *custom_errors.CustomError) {
	latestVersion, err := dbVersionGetLatest()
	if err != nil {
		return nil, err.AST("store performances")
	}
	newMappers := buildMappersFromRawData(performances)
	newTableIndex := (latestVersion.DbIndex + 1) % 2

	err = dbCountryDeleteAll(newTableIndex)
	if err != nil {
		return nil, err.AST("store performances")
	}
	err = dbAppDeleteAll(newTableIndex)
	if err != nil {
		return nil, err.AST("store performances")
	}
	err = dbSdkDeleteAll(newTableIndex)
	if err != nil {
		return nil, err.AST("store performances")
	}
	err = dbPerformanceDeleteAll(newTableIndex)
	if err != nil {
		return nil, err.AST("store performances")
	}

	for _, item := range newMappers.countries {
		err = dbCountryCreate(newTableIndex, &item)
		if err != nil {
			return nil, err.AST("store performances")
		}
	}

	for _, item := range newMappers.apps {
		err = dbAppCreate(newTableIndex, &item)
		if err != nil {
			return nil, err.AST("store performances")
		}
	}

	for _, item := range newMappers.sdks {
		err = dbSdkCreate(newTableIndex, &item)
		if err != nil {
			return nil, err.AST("store performances")
		}
	}
	// TODO consider bulk insert for performance...
	for index, item := range performances {
		adType, err := constants.AdTypeNameToId(item.AdType)
		if err != nil {
			return nil, err.AST("store performances")
		}
		err = dbPerformanceCreate(newTableIndex, &Performance{
			ID:      index,
			AdType:  adType,
			Country: newMappers.countryNameToId[item.Country],
			App:     newMappers.appNameToId[item.App],
			Sdk:     newMappers.sdkNameToId[item.Sdk],
			Score:   item.Score,
		})
		if err != nil {
			return nil, err.AST("store performances")
		}
	}

	version, err := dbVersionCreate(&Version{DbIndex: newTableIndex})
	if err != nil {
		return nil, err.AST("store performances")
	}
	// TODO place notify all services through some push pull service
	return version, nil
}

func CreateUser(user *structs.User) (*structs.User, *custom_errors.CustomError) {
	return dbUserCreate(user)
}
func GetUserByUsername(username string) (*structs.User, *custom_errors.CustomError) {
	return dbUserGetUserByUsername(username)
}
func GetUserById(userId int) (*structs.User, *custom_errors.CustomError) {
	return dbUserGetUserById(userId)
}
func CreateSession(session *structs.Session) (*structs.Session, *custom_errors.CustomError) {
	return dbSessionCreate(session)
}
func GetSessionById(sessionId int) (*structs.Session, *custom_errors.CustomError) {
	return dbSessionGetSessionById(sessionId)
}
func DeleteUserSessions(userId int) *custom_errors.CustomError {
	return dbSessionDeleteByUserId(userId)
}
