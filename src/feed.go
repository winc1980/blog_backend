package main

import (
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
)

type Article struct {
	Name      string
	Link      string
	Title     string
	Published time.Time
}

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
	var articles []Article

	for _, item := range feed.Items {
		var article Article
		article.Name = item.Authors[0].Name
		article.Link = item.Link
		article.Title = item.Title
		article.Published = *item.PublishedParsed
		articles = append(articles, article)
	}
	respond(w, r, http.StatusOK, articles)
}
