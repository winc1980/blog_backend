package main

import (
	"context"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *Server) HandleFeeds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.HandleFeedsGet(w, r)
		return
	}
	respondErr(w, r, http.StatusNotFound)
}

func (s *Server) HandleFeedsGet(w http.ResponseWriter, r *http.Request) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://zenn.dev/mattn/feed?all=1")
	log.Println(feed)
	ctx := context.TODO()
	db := s.client.Database("winc")
	collection := db.Collection("articles")

	for _, item := range feed.Items {
		_, err := collection.InsertOne(ctx, bson.D{
			{Key: "Name", Value: item.Authors[0].Name},
			{Key: "Link", Value: item.Link},
			{Key: "Title", Value: item.Title},
			{Key: "Published", Value: *item.PublishedParsed},
		})
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
	}
}
