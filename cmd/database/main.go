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
var mappers *Mappers

func createTables() {
	success := true
	err := dbUserCreateTableIfNot()
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
	success = success && dbExec(db, "CREATE TABLE IF NOT EXISTS versions("+
		"id int primary key auto_increment,"+
		"db_index int NOT NULL,"+
		"created_at TIMESTAMP default CURRENT_TIMESTAMP);")

	success = success && dbExec(db, "CREATE TABLE IF NOT EXISTS sessions ("+
		"id int primary key auto_increment,"+
		"created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,"+
		"user_id int NOT NULL);")

	if success {
		fmt.Println("All tables were successfully created.")
	}

}

func getLatestVersion() (*Version, error) {
	fmt.Println("getla")
	var version Version
	results, err := db.Query("SELECT id, created_at FROM versions ORDER BY id DESC LIMIT 1;")
	if err != nil {
		return nil, err
	}
	fmt.Println("Gavesomething")
	for results.Next() {
		err = results.Scan(&version.ID, &version.CreatedAt)
		if err != nil {
			return nil, err
		}
	}
	return &version, nil
}

func createNewVersion() (*Version, error) {
	success := dbExec(db, "INSERT INTO versions (created_at) VALUES (CURRENT_TIMESTAMP);")
	if success {
		return getLatestVersion()
	} else {
		return nil, errors.New("Error at inseritng into versions")
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
	getLatestVersion()
	res, _ := createNewVersion()
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

func loadMappersFromDb() error {
	return errors.New("bleh")
}

func reloadEverythingFromDBIntoCache() error {
	return nil
}

func GetSdks(options *structs.QueryOptions) (*structs.PerformanceResponse, error) {
	latestVersion, err := getLatestVersion()
	if err != nil {
		return nil, err
	}
	if mappers == nil {
		err := loadMappersFromDb()
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
	cacheVersionId := getCacheVersionId()
	if latestVersion.ID != cacheVersionId {
		err := reloadEverythingFromDBIntoCache()
		if err != nil {
			return nil, err
		}
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
func StorePerformances(performances []structs.Performance) {

}
