package api

import (
	"github.com/gorilla/mux"
	"log"
	"made.by.jst10/outfit7/hancock/cmd/config"
	"net/http"
)

// Since there are just few static rules I am ok with that.
const SdkAdMob = "AdMob"
const SdkAdMobOptOut = "SdkAdMob-OptOut"
const SdkFacebook = "Facebook"
const CountryCN = "CN"

const tokenCookieName = "token"
const refreshCookieName = "refresh_token"


func StartApi(configs *config.ApiConfigs) {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/api/auth", authHandler).Methods("POST", "PUT", "DELETE")
	myRouter.HandleFunc("/api/performances", performancesHandler).Methods("POST", "GET")
	log.Println("Starting http server on port:", configs.Port)
	log.Fatal(http.ListenAndServe(":"+configs.Port, myRouter))
}
