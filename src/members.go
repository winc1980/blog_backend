package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return
	}
	respondErr(w, r, http.StatusNotFound)
}

func (s *Server) HandleMembersPost(w http.ResponseWriter, r *http.Request) {
	var member Member
	err := decodeBody(r, &member)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	githubid, err := s.GetCurrentUser(w, r)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	if githubid != member.ID {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	db := s.client.Database("winc")
	collection := db.Collection("members")
	_, err = s.findMemberByID(member.ID)
	if err != mongo.ErrNoDocuments && err != nil {
		respondErr(w, r, http.StatusBadRequest, "member already exists")
		return
	}
	_, err = collection.InsertOne(context.TODO(), member)
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
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	if githubid != member.ID {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	db := s.client.Database("winc")
	collection := db.Collection("members")
	_, err = collection.UpdateOne(context.TODO(), bson.D{{Key: "id", Value: member.ID}}, member)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
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
