package main

import (
	"context"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Server) HandleMembersList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleMembersListGet(w, r)
		return
	}
	respondErr(w, r, http.StatusNotFound)
}

func (s *Server) handleMembersListGet(w http.ResponseWriter, r *http.Request) {
	db := s.client.Database("winc")
	collection := db.Collection("members")
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
	if results == nil {
		respond(w, r, http.StatusOK, []Member{})
		return
	}

	respond(w, r, http.StatusOK, results)
}
