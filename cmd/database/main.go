package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	"made.by.jst10/outfit7/hancock/cmd/structs"
)

var db *sql.DB
var mappersVersionId = -1
var mappers *Mappers

func createTables() {
	success := true
	err := dbUserCreateTableIfNot()
	if err != nil {
		success = false
	}
	err = dbVersionCreateTableIfNot()
	if err != nil {
		success = false
	}
	err = dbSessionCreateTableIfNot()
	if err != nil {
		success = false
	}
	for i := 0; i < 2; i++ {
		err = dbSdkCreateTablesIfNot(i)
		if err != nil {
			success = false
		}
		err = dbAppCreateTablesIfNot(i)
		if err != nil {
			success = false
		}
		err = dbCountryCreateTablesIfNot(i)
		if err != nil {
			success = false
		}
		err = dbPerformanceCreateTablesIfNot(i)
		if err != nil {
			success = false
		}
	}

	if success {
		fmt.Println("All tables were successfully created.")
	}

}

func InitDatabase() {
	var err error
	db, err = sql.Open("mysql",
		"root:root@tcp(127.0.0.1:3306)/hancock")
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {

		log.Fatal(err)
	}
	fmt.Println("DB connection has successfully initialized")
	createTables()
	res, e := dbVersionGetLatest()
	fmt.Println(e)
	fmt.Println(res)
	fmt.Print("Read")
	res1, err1 := dbVersionAll()
	fmt.Println("res1")
	fmt.Println(res1)
	fmt.Println(err1)
	e = dbVersionCreate(Version{DbIndex: 1})
	fmt.Println(e)
	res, _ = dbVersionGetLatest()
	fmt.Println(res)
	//err = dbUserCreate(structs.User{Role: constants.UserRoleGuest, Username: "test", Password: "test"})
	r, e := dbUserAll()
	//e:= dbUserDeleteAll()
	//e:= dbUserDeleteById(3)
	//r,e:= dbUserGetUserByUsername("test")
	fmt.Println(r)
	fmt.Println(e)
	//defer db.Close()
}

func reloadMappersFromDb(version *Version) error {
	countries, countryErr := dbCountryAll(version.DbIndex)
	if countryErr != nil {
		return countryErr
	}
	apps, appErr := dbAppAll(version.DbIndex)
	if appErr != nil {
		return appErr
	}
	sdks, sdkErr := dbSdkAll(version.DbIndex)
	if sdkErr != nil {
		return sdkErr
	}
	mappers = buildMappersFromDBData(countries, apps, sdks)
	mappersVersionId = version.ID
	return nil

}

func reloadCacheFromDb(version *Version) error {
	performances, err := dbPerformanceAll(version.DbIndex)
	if err != nil {
		return err
	}
	savePerformancesInCache(version.ID, mappers.countries, mappers.apps, performances)
	return nil
}

func GetPerformances(options *structs.QueryOptions) (*structs.PerformanceResponse, error) {
	latestVersion, err := dbVersionGetLatest()
	if err != nil {
		return nil, err
	}

	if latestVersion.ID != mappersVersionId {
		err := reloadMappersFromDb(latestVersion)
		if err != nil {
			return nil, err
		}
	}

	if latestVersion.ID != cacheVersionId {
		err := reloadCacheFromDb(latestVersion)
		if err != nil {
			return nil, err
		}
	}

	countryID, existCountryName := mappers.countryNameToId[options.Country]
	appID, existAppName := mappers.appNameToId[options.AppName]
	if !existAppName {
		return nil, errors.New("App name is not valid")
	}
	if !existCountryName {
		return nil, errors.New("Country name is not valid")
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
func StorePerformances(performances []structs.Performance) error {
	latestVersion, err := dbVersionGetLatest()
	if err != nil {
		return err
	}
	newMappers := buildMappersFromRawData(performances)
	newTableIndex := (latestVersion.DbIndex + 1) % 2

	err = dbCountryDeleteAll(newTableIndex)
	if err != nil {
		return err
	}
	err = dbAppDeleteAll(newTableIndex)
	if err != nil {
		return err
	}
	err = dbSdkDeleteAll(newTableIndex)
	if err != nil {
		return err
	}
	err = dbPerformanceDeleteAll(newTableIndex)
	if err != nil {
		return err
	}

	for _, item := range newMappers.countries {
		err = dbCountryCreate(newTableIndex, &item)
		if err != nil {
			return err
		}
	}

	for _, item := range newMappers.apps {
		err = dbAppCreate(newTableIndex, &item)
		if err != nil {
			return err
		}
	}

	for _, item := range newMappers.sdks {
		err = dbSdkCreate(newTableIndex, &item)
		if err != nil {
			return err
		}
	}
	// TODO consider bulk insert for performance...
	for index, item := range performances {
		adType, addErr := constants.AdTypeNameToId(item.AdType)
		if err != nil {
			return addErr
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
			return err
		}
	}

	err = dbVersionCreate(Version{DbIndex: newTableIndex})
	if err != nil {
		return err
	}
	// TODO place notify all services through some push pull service
	return nil
}

func CreateUser(user *structs.User) error {
	return dbUserCreate(user)
}
func GetUserByUsername(username string)  (*structs.User, error) {
	return dbUserGetUserByUsername(username)
}
func GetUserById(userId int)  (*structs.User, error) {
	return dbUserGetUserById(userId)
}
func GetSessionById(sessionId int)  (*structs.Session, error) {
	return dbSessionGetSessionById(sessionId)
}
func DeleteUserSessions(userId int) error {
	return dbSessionDeleteByUserId(userId)
}
