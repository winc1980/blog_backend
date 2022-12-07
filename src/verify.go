package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type Token struct {
	Token string `json:"token"`
}

type GithubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

var (
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "https://blog.winc.ne.jp/oauth/callback",
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
)

func NeedToken(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var OAuthToken Token
		decodeBody(r, &OAuthToken)
		token := &oauth2.Token{AccessToken: OAuthToken.Token}
		client := oauthConfig.Client(context.Background(), token)
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
		fn(w, r)
	}
}

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) (string, error) {
	var OAuthToken Token
	decodeBody(r, &OAuthToken)
	token := &oauth2.Token{AccessToken: OAuthToken.Token}
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve user info: %s", err.Error()))
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return "", err
	}
	var user GithubUser
	err = json.Unmarshal(body, &user)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return "", err
	}
	return user.Login, nil
}
