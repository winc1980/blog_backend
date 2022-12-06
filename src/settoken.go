package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	respondErr(w, r, http.StatusNotFound)
}
func (s *Server) HandleSetTokenPost(w http.ResponseWriter, r *http.Request) {
	var token Token
	decodeBody(r, &token)
	http.SetCookie(w, &http.Cookie{
		Name:     "OAuthToken",
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
	db := s.client.Database("winc")
	collection := db.Collection("github_team_members")
	count, err := collection.CountDocuments(context.TODO(), bson.D{{Key: "id", Value: user.Login}})
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return
	}
	if count == 0 {
		respondErr(w, r, http.StatusBadRequest, "")
		return
	}
	memberCollection := db.Collection("members")
	_, err = s.findMemberByID(user.Login)
	if err != mongo.ErrNoDocuments && err != nil {
		respondErr(w, r, http.StatusBadRequest, "member already exists")
		return
	}
	_, err = memberCollection.InsertOne(context.TODO(), bson.D{{Key: "id", Value: user.Login}})
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}

	respond(w, r, http.StatusOK, "")
}
