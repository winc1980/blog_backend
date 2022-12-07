package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleLink struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Title     string
	Image     string
	Published time.Time
}

type ArticlesResponse struct {
	Items       []primitive.M
	Pages_count int64
}

type Items struct {
	Links    []ArticleLink
	Articles []Article
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
	var limit int64 = 18
	var opts *options.FindOptions
	var filter bson.D = bson.D{}
	query := r.URL.Query()
	pageQuery := query.Get("page")
	page, err := strconv.ParseInt(pageQuery, 10, 64)
	page -= 1
	if err != nil {
		opts = options.Find().SetSort(bson.D{{Key: "published", Value: -1}}).SetLimit(limit)
	} else {
		opts = options.Find().SetSort(bson.D{{Key: "published", Value: -1}}).SetLimit(limit).SetSkip(limit * page)
	}
	typeQuery := query.Get("type")
	if typeQuery != "" {
		filter = bson.D{{Key: "type", Value: typeQuery}}
	}
	cursor, err := collection.Find(context.TODO(), filter, opts)

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

	count, err := collection.CountDocuments(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}
	respond(w, r, http.StatusOK, ArticlesResponse{Items: results, Pages_count: (count / limit) + 1})
}

func (s *Server) findArticleByLink(link string) (ArticleLink, error) {
	db := s.client.Database("winc")
	collection := db.Collection("articles")
	var result bson.Raw
	err := collection.FindOne(context.TODO(), bson.D{{Key: "link", Value: link}}).Decode(&result)
	if err != nil {
		return ArticleLink{}, err
	}
	var mapData map[string]interface{}
	json.Unmarshal([]byte(result.String()), &mapData)
	var article ArticleLink
	json.Unmarshal([]byte(result.String()), &article)
	return article, nil
}

func (s *Server) checkArticleExists(link string) (bool, error) {
	db := s.client.Database("winc")
	collection := db.Collection("articles")
	count, err := collection.CountDocuments(context.TODO(), bson.D{{Key: "link", Value: link}})
	if err != nil {
		log.Println(err)
		return false, err
	}
	return count != 0, nil
}
