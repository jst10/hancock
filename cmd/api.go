package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"net/http"
	"strings"
)

// Since there are just few static rules I am ok with that.
const SdkAdMob = "AdMob"
const SdkAdMobOptOut = "SdkAdMob-OptOut"
const SdkFacebook = "Facebook"
const CountryCN = "CN"

const tokenCookie = "token"
const refreshCookie = "refresh_token"

func removeSdksFromPerformanceResults(sdksToRemove map[string]bool, listOfPerformances []structs.Performance) []structs.Performance {
	results := make([]structs.Performance, 0)
	for _, performance := range listOfPerformances {
		if !sdksToRemove[performance.Sdk] {
			results = append(results, performance)
		}
	}
	return results
}

func postProcess(qo *structs.QueryOptions, performances *structs.PerformanceResponse) *structs.PerformanceResponse {
	sdksToRemove := make(map[string]bool)
	if qo.Country == CountryCN {
		sdksToRemove[SdkFacebook] = true
	}
	if qo.OsVersion == "9" || strings.Index(qo.OsVersion, "9.") == 0 {
		sdksToRemove[SdkAdMob] = true
	} else {
		// since we know that are always all sdk available in performance data
		sdksToRemove[SdkAdMobOptOut] = true
	}

	performances.Banner = removeSdksFromPerformanceResults(sdksToRemove, performances.Banner)
	performances.Interstitial = removeSdksFromPerformanceResults(sdksToRemove, performances.Interstitial)
	performances.Reward = removeSdksFromPerformanceResults(sdksToRemove, performances.Reward)
	return performances
}

func handleAuthenticate(w http.ResponseWriter, r *http.Request) {
	var authData structs.AuthData
	err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenWrapper, refreshTokenWrapper, err := auth.AuthenticateUser(&authData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    tokenCookie,
		Value:   tokenWrapper.Token,
		Expires: *tokenWrapper.Expiration,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    refreshCookie,
		Value:   refreshTokenWrapper.Token,
		Expires: *refreshTokenWrapper.Expiration,
	})
	w.WriteHeader(http.StatusOK)
}

func handleReAuthenticate(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(refreshCookie)
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	refreshToken := cookie.Value
	tokenData, err := auth.VerifyJWRT(refreshToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tokenWrapper, err := auth.ReAuthenticateUser(tokenData.SessionId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    tokenCookie,
		Value:   tokenWrapper.Token,
		Expires: *tokenWrapper.Expiration,
	})
	w.WriteHeader(http.StatusOK)
}

func handleDeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authenticateRequest(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = auth.DeAuthenticateUser(tokenData.UserId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) (*structs.TokenData, error) {
	cookie, err := r.Cookie(tokenCookie)
	if err != nil {
		return nil, err
	}
	token := cookie.Value
	tokenData, err := auth.VerifyJWT(token)
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

func handleGetPerformances(w http.ResponseWriter, r *http.Request) {
	//_, err := authenticateRequest(w, r)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	query := r.URL.Query()
	qo := structs.QueryOptions{}
	qo.Country = query.Get("country")
	qo.Platform = query.Get("platform")
	qo.OsVersion = query.Get("os_version")
	qo.AppName = query.Get("app_name")
	qo.AppVersion = query.Get("app_version")
	if len(qo.Country) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetMissingQueryParamError("country"))
		return
	}
	if len(qo.Platform) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetMissingQueryParamError("platform"))
		return
	}
	if len(qo.OsVersion) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetMissingQueryParamError("os_version"))
		return
	}
	if len(qo.AppName) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetMissingQueryParamError("app_name"))
		return
	}

	performances, err := database.GetPerformances(&qo)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	performances = postProcess(&qo, performances)
	json.NewEncoder(w).Encode(performances)
}
func handleCreatePerformances(w http.ResponseWriter, r *http.Request) {
	var err error
	//tokenData, err := authenticateRequest(w, r)
	//if err != nil {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//if tokenData.Role < constants.UserRoleAdmin {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}

	performances := make([]structs.Performance, 0)
	err = json.NewDecoder(r.Body).Decode(&performances)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = database.StorePerformances(performances)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func performancesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGetPerformances(w, r)
	} else if r.Method == "POST" {
		handleCreatePerformances(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetNotFoundError())
	}
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		handleDeAuthenticate(w, r)
	} else if r.Method == "POST" {
		handleAuthenticate(w, r)
	} else if r.Method == "PUT" {
		handleReAuthenticate(w, r)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(custom_errors.GetNotFoundError())
	}
}

func startApi() {

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/auth", authHandler)
	myRouter.HandleFunc("/api/performances", performancesHandler)
	log.Println("Starting http server on port 10000")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
