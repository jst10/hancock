package api

import (
	"encoding/json"
	"made.by.jst10/outfit7/hancock/cmd/constants"
	"made.by.jst10/outfit7/hancock/cmd/custom_errors"
	"made.by.jst10/outfit7/hancock/cmd/database"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"net/http"
	"strconv"
)

func handleGetPerformances(w http.ResponseWriter, r *http.Request) {
	_, err := authenticateRequest( r)
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
	performances = postProcessPerformancesResults(&qo, performances)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(performances)
}
func handleCreatePerformances(w http.ResponseWriter, r *http.Request) {
	tokenData, err := authenticateRequest( r)
	if err != nil {
		respondError(w, http.StatusUnauthorized, err)
		return
	}
	if tokenData.Role < constants.UserRoleAdmin {
		respondError(w, http.StatusForbidden, custom_errors.GetNotAllowed("Minimum role: "+strconv.Itoa(constants.UserRoleAdmin)))
		return
	}

	performances := make([]*structs.Performance, 0)
	err = decodeJSONBody(w, r, &performances)
	if err != nil {
		respondError(w, http.StatusBadRequest, err)
		return
	}
	if !arePerformancesValid(performances) {
		respondError(w, http.StatusForbidden, custom_errors.GetNotValidDataError("Make sure that are all 3 ad types in you data, and complete matrix for type,country and app_name"))
		return
	}
	_, err = database.StorePerformances(performances)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func performancesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGetPerformances(w, r)
	} else if r.Method == "POST" {
		handleCreatePerformances(w, r)
	}
}