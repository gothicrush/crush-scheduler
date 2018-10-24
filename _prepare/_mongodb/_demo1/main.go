package main

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"log"
	"time"
)

func main() {

	// 建立连接
	client, err := mongo.Connect(context.TODO(), "mongodb://127.0.0.1:27017", clientopt.ConnectTimeout(5*time.Second))

	if err != nil {
		log.Println("mongo.Connect error : ", err)
		return
	}

	// 选择数据库
	database := client.Database("my_db")

	// 选择数据表
	collection := database.Collection("my_collection")

	collection = collection
}
