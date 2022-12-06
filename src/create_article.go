package main

import (
	"context"
	"net/http"
	"time"
)

type Article struct {
	UUID      string
	ID        string `json:"id"`
	Title     string `json:"title"`
	Published time.Time
}

func (s *Server) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.HandleCreateArticlePost(w, r)
		return
	}
}

func (s *Server) HandleCreateArticlePost(w http.ResponseWriter, r *http.Request) {
	var article Article
	err := decodeBody(r, &article)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	githubid, err := s.GetCurrentUser(w, r)
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	if githubid != article.ID {
		respondErr(w, r, http.StatusBadRequest, "", err)
		return
	}
	db := s.client.Database("winc")
	collection := db.Collection("articles_original")
	_, err = collection.InsertOne(context.TODO(), article)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, "article created")
}
