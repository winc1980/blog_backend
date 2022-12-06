package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
)

type Member struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Zenn  string `json:"zenn"`
	Qiita string `json:"qiita"`
}

func (s *Server) HandleMembers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.HandleMembersPost(w, r)
		return
	case "PUT":
		s.HandleMembersPut(w, r)
		return
	}
	respondErr(w, r, http.StatusNotFound)
}

func (s *Server) HandleMembersPost(w http.ResponseWriter, r *http.Request) {
	var token Token
	decodeBody(r, &token)
	OAuthToken := &oauth2.Token{AccessToken: token.Token}
	client := oauthConfig.Client(context.Background(), OAuthToken)
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

func (s *Server) HandleMembersPut(w http.ResponseWriter, r *http.Request) {
	var member Member
	err := decodeBody(r, &member)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	githubid, err := s.GetCurrentUser(w, r)
	if err != nil {
		return
	}
	db := s.client.Database("winc")
	collection := db.Collection("members")
	filter := bson.D{{Key: "id", Value: githubid}}
	if member.Name != "" {
		update := bson.D{{"$set", bson.D{{Key: "name", Value: member.Name}}}}
		_, err = collection.UpdateOne(
			context.TODO(),
			filter,
			update,
		)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
	}
	if member.Zenn != "" {
		update := bson.D{{"$set", bson.D{{Key: "zenn", Value: member.Zenn}}}}
		_, err = collection.UpdateOne(
			context.TODO(),
			filter,
			update,
		)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
	}
	if member.Qiita != "" {
		update := bson.D{{"$set", bson.D{{Key: "qiita", Value: member.Qiita}}}}
		_, err = collection.UpdateOne(
			context.TODO(),
			filter,
			update,
		)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
	}
	respond(w, r, http.StatusOK, "")
}

func (s *Server) findMemberByID(id string) (Member, error) {
	ctx := context.TODO()
	db := s.client.Database("winc")
	collection := db.Collection("users")
	var result bson.Raw
	err := collection.FindOne(ctx, bson.D{{Key: "id", Value: id}}, options.FindOne()).Decode(&result)
	if err != nil {
		return Member{}, err
	}
	var mapData map[string]interface{}
	json.Unmarshal([]byte(result.String()), &mapData)
	var member Member
	json.Unmarshal([]byte(result.String()), &member)
	return member, nil
}
