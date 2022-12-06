package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load("/app/.env")
	if err != nil {
		log.Panicln(err)
	}
	ctx := context.TODO()

	uri := os.Getenv("MONGODB_URI")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	s := &Server{client}
	go func() {
		s.WINCMembers()
		s.FeedCollector()
		t := time.NewTicker(15 * time.Minute)
		for {
			<-t.C
			log.Println("start collect")
			s.WINCMembers()
			s.FeedCollector()
		}
	}()
	log.Println("start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/create_article/", withCORS(NeedToken(s.HandleCreateArticle)))
	mux.HandleFunc("/articles/", withCORS(s.HandleArticles))
	mux.HandleFunc("/members/", withCORS(NeedToken(s.HandleMembers)))
	mux.HandleFunc("/members_list/", withCORS(s.HandleMembersList))
	mux.HandleFunc("/github_team/", withCORS(s.HandleGithubTeam))
	http.ListenAndServe(":8888", mux)
}

type Server struct {
	client *mongo.Client
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		switch r.Method {
		case "OPTIONS":
			log.Println("preflight")
			return
		}
		fn(w, r)
	}
}
