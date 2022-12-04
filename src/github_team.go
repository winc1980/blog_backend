package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/oauth2"
)

type Github struct {
	ID string `json:"id"`
}

func (s *Server) HandleGithubTeam(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.WINCMembers()
		s.handleGithubTeamGet(w, r)
		return
	}
}

func (s *Server) handleGithubTeamGet(w http.ResponseWriter, r *http.Request) {
	db := s.client.Database("winc")
	collection := db.Collection("github_team_members")
	cursor, err := collection.Find(context.TODO(), bson.M{})

	if err != nil {
		if err == mongo.ErrNoDocuments {
			respondErr(w, r, http.StatusInternalServerError, "mongo: no result")
			return
		}
		respondErr(w, r, http.StatusInternalServerError)
		return
	}
	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		return
	}

	respond(w, r, http.StatusOK, results)
}

func (s *Server) WINCMembers() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_ACCESS_TOKEN")},
	)
	tc := oauth2.NewClient(context.TODO(), ts)

	client := github.NewClient(tc)
	db := s.client.Database("winc")
	collection := db.Collection("github_team_members")
	pagecnt := 1
	for {
		listoption := github.ListOptions{Page: pagecnt}
		members, resp, err := client.Organizations.ListMembers(
			context.Background(),
			"winc1980",
			&github.ListMembersOptions{ListOptions: listoption},
		)
		if err != nil {
			return
		}

		for _, member := range members {

			count, err := collection.CountDocuments(context.TODO(), bson.D{{Key: "id", Value: *member.Login}})
			if err != nil {
				log.Println(err)
				return
			}
			if count == 0 {
				collection.InsertOne(context.TODO(), Github{ID: *member.Login})
			}
		}
		if resp.NextPage == 0 {
			break
		}
		pagecnt++
	}

	count, _ := collection.CountDocuments(context.TODO(), bson.D{})
	log.Println(count)
}
