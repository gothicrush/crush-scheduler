package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"log"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Content   string    `bson:"content"`
	Err       string    `bson:"err"`
	TimePoint TimePoint `bson:"timePoint"`
}

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

	// 插入
	record := &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "content",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}
	result, err := collection.InsertOne(context.TODO(), record)

	if err != nil {
		log.Println("InsertOne error : ", err)
		return
	} else {
		fmt.Println(result.InsertedID.(objectid.ObjectID))
	}
}
