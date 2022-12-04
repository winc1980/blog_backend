package main

import (
	"context"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Article struct {
	Name      string `json:"name"`
	Link      string `json:"link"`
	Title     string
	Published time.Time
}

type Articles struct {
	ID   string    `json:"id"`
	List []Article `json:"list"`
}

func (s *Server) HandleArticles(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.handleArticlesGet(w, r)
		return
	}
	respondErr(w, r, http.StatusNotFound)
}

func (s *Server) handleArticlesGet(w http.ResponseWriter, r *http.Request) {
	db := s.client.Database("winc")
	collection := db.Collection("articles")
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
