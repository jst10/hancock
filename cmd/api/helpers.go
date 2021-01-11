package api

import (
	"encoding/json"
	"fmt"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	"made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"net/http"
	"strings"
)

func removeSdksFromPerformanceResults(sdksToRemove map[string]bool, listOfPerformances []structs.SdkScore) []structs.SdkScore {
	results := make([]structs.SdkScore, 0)
	for _, performance := range listOfPerformances {
		if !sdksToRemove[performance.Sdk] {
			results = append(results, performance)
		}
	}
	return results
}

func postProcessPerformancesResults(qo *structs.QueryOptions, performances *structs.PerformanceResponse) *structs.PerformanceResponse {
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

func respondError(w http.ResponseWriter, statusCode int, err *custom_errors.CustomError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(err)
	fmt.Println("Error")
	fmt.Println(err)
	fmt.Println(err.StackTrace)
	fmt.Println(err.OriginalError)
}
func getCookieFromRequest(r *http.Request, cookieName string) (*http.Cookie, *custom_errors.CustomError) {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return nil, custom_errors.GetErrorGettingCookie(err, cookieName)
	} else {
		return cookie, nil
	}
}

func authenticateRequest( r *http.Request) (*structs.TokenData, *custom_errors.CustomError) {
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