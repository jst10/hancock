package api

import (
	"encoding/json"
	"made.by.jst10/outfit7/hancock/cmd/auth"
	"made.by.jst10/outfit7/hancock/cmd/structs"
	"net/http"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "DELETE" {
		handleDeAuthenticate(w, r)
	} else if r.Method == "POST" {
		handleAuthenticate(w, r)
	} else if r.Method == "PUT" {
		handleReAuthenticate(w, r)
	}
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
	// returning just for easier debugging from console
	json.NewEncoder(w).Encode(structs.TokensResponse{Token: tokenWrapper.Token, RefreshToken: refreshTokenWrapper.Token})
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
	tokenData, err := authenticateRequest(r)
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
