package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	custom_errors "made.by.jst10/outfit7/hancock/cmd/custom-errors"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"net/http"
	"strconv"
	"strings"
)

// Since there are just few static rules I am ok with that.
const SdkAdMob = "AdMob"
const SdkAdMobOptOut = "SdkAdMob-OptOut"
const SdkFacebook = "Facebook"
const CountryCN = "CN"

const tokenCookieName = "token"
const refreshCookieName = "refresh_token"

func respondError(w http.ResponseWriter, statusCode int, err *custom_errors.CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
	fmt.Println("Error")
	fmt.Println(err)
	fmt.Println(err.OriginalError)
}

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
	err := decodeJSONBody(w, r, &authData)
	//err := json.NewDecoder(r.Body).Decode(&authData)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	tokenWrapper, refreshTokenWrapper, err := auth.AuthenticateUser(&authData)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    tokenCookieName,
		Value:   tokenWrapper.Token,
		Expires: *tokenWrapper.Expiration,
	})

	http.SetCookie(w, &http.Cookie{
		Name:    refreshCookieName,
		Value:   refreshTokenWrapper.Token,
		Expires: *refreshTokenWrapper.Expiration,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// TODO merge token wrapper and token response structs together
	json.NewEncoder(w).Encode(structs.TokensResponse{Token: tokenWrapper.Token, RefreshToken: refreshTokenWrapper.Token})
}

func getCookieFromRequest(r *http.Request, cookieName string) (*http.Cookie, *custom_errors.CustomError) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, custom_errors.GetErrorGettingCookie(err, cookieName)
	} else {
		return cookie, nil
	}
}

func handleReAuthenticate(w http.ResponseWriter, r *http.Request) {
	cookie, err := getCookieFromRequest(r, refreshCookieName)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	refreshToken := cookie.Value
	tokenData, err := auth.VerifyJWRT(refreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	tokenWrapper, err := auth.ReAuthenticateUser(tokenData.SessionId)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    tokenCookieName,
		Value:   tokenWrapper.Token,
		Expires: *tokenWrapper.Expiration,
	})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(structs.TokensResponse{Token: tokenWrapper.Token})
}

func handleDeAuthenticate(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authenticateRequest(w, r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	err = auth.DeAuthenticateUser(tokenData.UserId)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   tokenCookieName,
		Value:  "",
		MaxAge: -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:   refreshCookieName,
		Value:  "",
		MaxAge: -1,
	})
	w.WriteHeader(http.StatusNoContent)
}

func authenticateRequest(w http.ResponseWriter, r *http.Request) (*structs.TokenData, *custom_errors.CustomError) {
	cookie, err := getCookieFromRequest(r, tokenCookieName)
	if err != nil {
		return nil, custom_errors.GetCookieNotPresentError(tokenCookieName)
	}
	token := cookie.Value
	tokenData, err := auth.VerifyJWT(token)
	if err != nil {
		return nil, err
	}
	return tokenData, nil
}

func handleGetPerformances(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest(w, r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	query := r.URL.Query()
	qo := structs.QueryOptions{}
	qo.Country = query.Get("country")
	qo.Platform = query.Get("platform")
	qo.OsVersion = query.Get("os_version")
	qo.AppName = query.Get("app_name")
	qo.AppVersion = query.Get("app_version")
	if len(qo.Country) == 0 {
		respondError(w, http.StatusBadRequest, custom_errors.GetMissingQueryParamError("country"))
		return
	}
	if len(qo.Platform) == 0 {
		respondError(w, http.StatusBadRequest, custom_errors.GetMissingQueryParamError("platform"))
		return
	}
	if len(qo.OsVersion) == 0 {
		respondError(w, http.StatusBadRequest, custom_errors.GetMissingQueryParamError("os_version"))
		return
	}
	if len(qo.AppName) == 0 {
		respondError(w, http.StatusBadRequest, custom_errors.GetMissingQueryParamError("app_name"))
		return
	}

	performances, err := database.GetPerformances(&qo)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	performances = postProcess(&qo, performances)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(performances)
}
func handleCreatePerformances(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authenticateRequest(w, r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	if tokenData.Role < constants.UserRoleAdmin {
		respondError(w, http.StatusForbidden, custom_errors.GetNotAllowed("Minimum role: "+strconv.Itoa(constants.UserRoleAdmin)))
		return
	}

	performances := make([]structs.Performance, 0)
	err = decodeJSONBody(w, r, &performances)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}

	fmt.Println("GOT BODY")
	fmt.Println(performances)

	_, err = database.StorePerformances(performances)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func performancesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("performancesHandler", r.Method)
	if r.Method == "GET" {
		handleGetPerformances(w, r)
	} else if r.Method == "POST" {
		handleCreatePerformances(w, r)
	}
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("authHandler", r.Method)
	if r.Method == "DELETE" {
		handleDeAuthenticate(w, r)
	} else if r.Method == "POST" {
		handleAuthenticate(w, r)
	} else if r.Method == "PUT" {
		handleReAuthenticate(w, r)
	}
}

func startApi() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/auth", authHandler).Methods("POST", "PUT", "DELETE")
	myRouter.HandleFunc("/api/performances", performancesHandler).Methods("POST", "GET")
	log.Println("Starting http server on port 10000")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
