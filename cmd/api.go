package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Id      string `json:"Id"`
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}


func advanceFilter(ro *QueryOptions) {

//1. ad mob doesn't work on any android os version 9
//2. no FB in CN
//3. AdMob-OptOut should be present in list only if there is no AdMob in list
}

func handleGetPerformances(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	fmt.Println(r.URL.Query())

	Articles := []Article{
		Article{Id: "1", Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Id: "2", Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}
	for index, article := range Articles {
		if index == 0 {
			json.NewEncoder(w).Encode(article)
		}
	}
}
func handleCreatePerformances(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(HttpError{"No kul1 "})
}

func performancesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {

		handleGetPerformances(w, r)
	} else if r.Method == "POST" {

		handleCreatePerformances(w, r)
	} else {
		w.WriteHeader(http.StatusFound)
		json.NewEncoder(w).Encode(HttpError{"No kul"})
	}

}
func createNewArticle(w http.ResponseWriter, r *http.Request) {

	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "%+v", string(reqBody))
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/performances", performancesHandler)
	log.Println("Starting http server on port 10000")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
