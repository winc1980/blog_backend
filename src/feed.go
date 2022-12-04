package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mmcdole/gofeed"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		s.ZennFeedCollector(member.Zenn)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) ZennFeedCollector(id string) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("https://zenn.dev/" + id + "/feed?all=1")
	log.Println(feed)
	ctx := context.TODO()
	db := s.client.Database("winc")
	collection := db.Collection("articles")

	for _, item := range feed.Items {
		_, err := s.findArticleByLink(item.Link)
		if err != mongo.ErrNoDocuments && err != nil {
			return
		}
		_, err = collection.InsertOne(ctx, bson.D{
			{Key: "Name", Value: item.Authors[0].Name},
			{Key: "Link", Value: item.Link},
			{Key: "Title", Value: item.Title},
			{Key: "Published", Value: *item.PublishedParsed},
		})
		if err != nil {
			return
		}
	}
}
