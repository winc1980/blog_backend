package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Qiita struct {
	Title      string    `json:"title"`
	Created_at time.Time `json:"created_at"`
	Link       string    `json:"url"`
	Image      string
}

func (s *Server) FeedCollector() {
	db := s.client.Database("winc")
	collection := db.Collection("members")
	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println(err)
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var doc bson.Raw
		if err := cursor.Decode(&doc); err != nil {
			log.Fatal(err)
		}
		var mapData map[string]interface{}
		json.Unmarshal([]byte(doc.String()), &mapData)
		var member Member
		json.Unmarshal([]byte(doc.String()), &member)
		if member.Zenn != "" {
			s.ZennLinkCollector(member.Zenn)
		}
		if member.Qiita != "" {
			s.QiitaLinkCollector(member.Qiita, member.GithubID)
		}
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ZennLinkCollector(id string) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://zenn.dev/" + id + "/feed?all=1")
	log.Println(feed)
	ctx := context.TODO()
	db := s.client.Database("winc")
	collection := db.Collection("articles")

	for _, item := range feed.Items {
		isExist, err := s.checkArticleExists(item.Link)
		if err != nil && err != mongo.ErrNoDocuments {
			return
		} else if isExist {
			continue
		}
		imageurl, _ := getOGImage(item.Link)
		_, err = collection.InsertOne(ctx, bson.D{
			{Key: "type", Value: "zenn"},
			{Key: "githubid", Value: id},
			{Key: "name", Value: item.Authors[0].Name},
			{Key: "link", Value: item.Link},
			{Key: "title", Value: item.Title},
			{Key: "image", Value: imageurl},
			{Key: "published", Value: *item.PublishedParsed},
		})
		if err != nil {
			return
		}
	}
}

func (s *Server) QiitaLinkCollector(qiitaID string, githubID string) {
	req, err := http.NewRequest(
		"GET",
		"https://qiita.com/api/v2/users/"+qiitaID+"/items?per_page=100",
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var response []Qiita
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println(err)
		return
	}
	ctx := context.TODO()
	db := s.client.Database("winc")
	collection := db.Collection("articles")
	for _, item := range response {
		isExist, err := s.checkArticleExists(item.Link)
		if err != nil && err != mongo.ErrNoDocuments {
			return
		} else if isExist {
			continue
		}
		imageurl, _ := getOGImage(item.Link)
		_, err = collection.InsertOne(ctx, bson.D{
			{Key: "type", Value: "qiita"},
			{Key: "githubid", Value: githubID},
			{Key: "name", Value: qiitaID},
			{Key: "link", Value: item.Link},
			{Key: "title", Value: item.Title},
			{Key: "image", Value: imageurl},
			{Key: "published", Value: item.Created_at},
		})
		if err != nil {
			return
		}
	}
}

func getOGImage(link string) (string, error) {
	req, err := http.NewRequest(
		"GET",
		link,
		nil,
	)
	if err != nil {
		return "", err
	}
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)
	og := opengraph.NewOpenGraph()
	err = og.ProcessHTML(strings.NewReader(html))

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return og.Images[0].URL, nil
}
