package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/oauth2"
)

type Token struct {
	Token string `json:"token"`
}

func (s *Server) HandleSetToken(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.HandleSetTokenPost(w, r)
		return
	}
}
func (s *Server) HandleSetTokenPost(w http.ResponseWriter, r *http.Request) {
	var token Token
	decodeBody(r, &token)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth",
		Value:    token.Token,
		Secure:   true,
		HttpOnly: true,
		Path:     "/",
	})
	authtoken := &oauth2.Token{AccessToken: token.Token}
	client := oauthConfig.Client(context.Background(), authtoken)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve user info: %s", err.Error()))
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return
	}
	var user GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return
	}

	respond(w, r, http.StatusOK, "")
}
