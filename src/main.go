package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
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
	mux.HandleFunc("create_article", withCORS(NeedToken(s.HandleCreateArticle)))
	mux.HandleFunc("/articles/", withCORS(s.HandleArticles))
	mux.HandleFunc("/members/", withCORS(NeedToken(s.HandleMembers)))
	mux.HandleFunc("/settoken/", withCORS(s.HandleSetToken))
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
		switch r.Method {
		case "OPTIONS":
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			return
		}
		fn(w, r)
	}
}

type GithubUser struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
}

var (
	oauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  "https://blog.winc.ne.jp/oauth/callback",
		Scopes:       []string{"user"},
		Endpoint:     github.Endpoint,
	}
)

func NeedToken(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("OAuthToken")
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, "OAuthToken cookie not found")
			return
		}

		token := &oauth2.Token{AccessToken: cookie.Value}
		client := oauthConfig.Client(context.Background(), token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			respondErr(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve user info: %s", err.Error()))
			return
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
		var user GithubUser
		err = json.Unmarshal(body, &user)
		if err != nil {
			respondErr(w, r, http.StatusInternalServerError)
			return
		}
		fn(w, r)
	}
}

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("OAuthToken")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, "OAuthToken cookie not found")
		return "", err
	}

	token := &oauth2.Token{AccessToken: cookie.Value}
	client := oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		respondErr(w, r, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve user info: %s", err.Error()))
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return "", err
	}
	var user GithubUser
	err = json.Unmarshal(body, &user)
	log.Println("user:", user)
	if err != nil {
		respondErr(w, r, http.StatusInternalServerError)
		return "", err
	}
	return user.Login, nil
}
