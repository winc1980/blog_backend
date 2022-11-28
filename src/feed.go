package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"
)

func (s *Server) HandleFeeds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		return
	}
}

func (s *Server) HandleFeedsGet(w http.ResponseWriter, r *http.Request) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://zenn.dev/mattn/feed?all=1")
	fmt.Println(feed)

	for _, line := range feed.FeedLink {
		log.Println(line)
	}
}
