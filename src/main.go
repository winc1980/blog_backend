package main

import (
	"context"
	"log"
	"net/http"
	"os"

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
	log.Println("start server")
	mux := http.NewServeMux()
	mux.HandleFunc("/links/", withCORS(s.HandleLinks))
	mux.HandleFunc("/feeds/", withCORS(s.HandleFeeds))
	http.ListenAndServe(":8888", mux)
}

type Server struct {
	client *mongo.Client
}

func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
