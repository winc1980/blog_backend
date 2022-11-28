package main

import (
	"fmt"
	"net/http"

	"github.com/mmcdole/gofeed"
)

func (s *Server) HandleFeeds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		return
	}
}

func feed() {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://zenn.dev/mattn/feed?all=1")
	fmt.Println(feed)
}
