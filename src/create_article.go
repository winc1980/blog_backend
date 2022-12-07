package main

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	UUID      string
	GithubID  string `json:"githubid"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	Published time.Time
}

func (s *Server) HandleCreateArticle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		s.HandleCreateArticlePost(w, r)
		return
	}
	respondErr(w, r, http.StatusNotFound)
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
	db := s.client.Database("winc")
	collection := db.Collection("articles_original")
	article.Published = time.Now()
	article.Published.Format("2006-01-02")
	article.GithubID = githubid
	u, err := uuid.NewRandom()
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	article.UUID = u.String()
	_, err = collection.InsertOne(context.TODO(), article)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError, err)
		return
	}
	respond(w, r, http.StatusOK, "article created")
}
