package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"os"
)

var userAuthData = &structs.AuthData{Username: "user", Password: "user"}
var adminAuthData = &structs.AuthData{Username: "admin", Password: "admin"}
var client http.Client

func makeARequest(req *http.Request, expectedStatusCode int) *http.Response {
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error at making request", err)
	}
	if resp.StatusCode != expectedStatusCode {
		log.Fatal("Invalid status ", resp.Status)
	}
	return resp
}

func extractBody(body io.ReadCloser, dst interface{}) {
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {
		log.Fatal("Error at decoding data", err)
	}
	body.Close()
}
func authenticate(user *structs.AuthData) *structs.TokensResponse {
	var jsonStr, _ = json.Marshal(user)
	req, err := http.NewRequest("POST", "http://localhost:10000/api/auth", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal("Error authenticating user", err)
	}
	resp := makeARequest(req, http.StatusOK)
	dst := &structs.TokensResponse{}
	extractBody(resp.Body, dst)
	return dst
}

func getPerformances(qo *structs.QueryOptions) *structs.PerformanceResponse {
	url := fmt.Sprintf("http://localhost:10000/api/performances?country=%s&platform=%s&os_version=%s&app_name=%s&app_version=%s",
		qo.Country, qo.Platform, qo.OsVersion, qo.AppName, qo.AppVersion)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Error authenticating user", err)
	}
	resp := makeARequest(req, http.StatusOK)
	dst := &structs.PerformanceResponse{}
	extractBody(resp.Body, dst)
	return dst
}

func createPerformances(performances []structs.Performance) {
	var jsonStr, _ = json.Marshal(performances)
	req, err := http.NewRequest("POST", "http://localhost:10000/api/performances", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		log.Fatal("Error making request for creating performances", err)
	}
	makeARequest(req, http.StatusCreated)
}

func initClient() {
	options := cookiejar.Options{}
	jar, err := cookiejar.New(&options)
	if err != nil {
		log.Fatal("error creating client", err)
	}
	client = http.Client{Jar: jar}
}
func main() {
	initClient()
	fmt.Println("Client created")
	authenticate(adminAuthData)
	fmt.Println("Start tests:")
	test1()
	test2()
	test3()
	test4()
	fmt.Println("All ok :)")
}

func test1() {
	nPerformances := make([]structs.Performance, 0)
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s1", 2))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s1", 2))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s1", 2))
	createPerformances(nPerformances)

	qo := structs.NewQueryOptions("c1", "p1", "os", "a1", "av")
	rPerformances := getPerformances(qo)

	if len(rPerformances.Banner) != 1 || len(rPerformances.Interstitial) != 1 || len(rPerformances.Reward) != 1 {
		log.Fatal("test1 length not ok")
	}
}

func ensureOrder(scores []structs.SdkScore, order []string) {
	if len(scores) != len(order) {
		log.Fatal("Length do not match")
	}
	for i, _ := range scores {
		if scores[i].Sdk != order[i] {
			log.Fatal("Order is not ok")
		}
	}
}

func test2() {
	nPerformances := make([]structs.Performance, 0)
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s1", 20))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s2", 30))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s3", 40))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s1", 4))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s2", 3))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s3", 3))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s1", 4))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s2", 8))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s3", 2))
	createPerformances(nPerformances)

	qo := structs.NewQueryOptions("c1", "p1", "os", "a1", "av")
	rPerformances := getPerformances(qo)

	ensureOrder(rPerformances.Banner, []string{"s3", "s2", "s1"})
	ensureOrder(rPerformances.Interstitial, []string{"s1", "s2", "s3"})
	ensureOrder(rPerformances.Reward, []string{"s2", "s1", "s3"})
}

func test3() {
	nPerformances := make([]structs.Performance, 0)
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s1", 20))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s2", 30))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a1", "s3", 40))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, "c1", "a2", "s3", 40))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s1", 4))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s2", 3))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a1", "s3", 3))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, "c1", "a2", "s3", 3))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s1", 4))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s2", 8))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a1", "s3", 2))
	nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, "c1", "a2", "s3", 2))
	createPerformances(nPerformances)

	qo := structs.NewQueryOptions("c1", "p1", "os", "a1", "av")
	rPerformances := getPerformances(qo)

	ensureOrder(rPerformances.Banner, []string{"s3", "s2", "s1"})
	ensureOrder(rPerformances.Interstitial, []string{"s1", "s2", "s3"})
	ensureOrder(rPerformances.Reward, []string{"s2", "s1", "s3"})

	qo = structs.NewQueryOptions("c1", "p1", "os", "a2", "av")
	rPerformances = getPerformances(qo)
	ensureOrder(rPerformances.Banner, []string{"s3"})
	ensureOrder(rPerformances.Interstitial, []string{"s3"})
	ensureOrder(rPerformances.Reward, []string{"s3"})
}

func loadJsonFileOfStrings(filePath string) []string {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Error loading json file", err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal("Error loading json file", err)
	}
	var results []string
	json.Unmarshal(byteValue, &results)
	return results
}
func test4() {
	randSource := rand.NewSource(1)
	random := rand.New(randSource)
	countries := loadJsonFileOfStrings("data/countries.json")
	sdks := loadJsonFileOfStrings("data/sdks.json")
	nPerformances := make([]structs.Performance, 0)
	// add 10 means that we have 156870 entries to insert into db (was tested also with 100)
	for _, country := range countries {
		for i := 0; i < 10; i++ {
			app := fmt.Sprintf("app_%d", i)
			for _, sdk := range sdks {
				nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeBanner, country, app, sdk, random.Intn(100)))
				nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeInterstitial, country, app, sdk,random.Intn(100)))
				nPerformances = append(nPerformances, *structs.NewPerformance(constants.AdTypeRewarded, country, app, sdk, random.Intn(100)))
			}
		}
	}
	fmt.Printf("Creating %d performances", len(nPerformances))
	createPerformances(nPerformances)
	// TODO make request validate
}
