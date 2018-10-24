package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
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
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {

	// 建立连接
	client, err := mongo.Connect(context.TODO(), "mongodb://127.0.0.1:27017", clientopt.ConnectTimeout(1*time.Second))

	if err != nil {
		log.Println("mongo.Connect error : ", err)
		return
	}

	// 选择数据库
	database := client.Database("my_db")

	// 选择数据表
	collection := database.Collection("my_collection")

	// 查询

	// 创建查询结构体
	condition := &FindByJobName{
		JobName: "job10",
	}

	cursor, err := collection.Find(context.TODO(), condition, findopt.Skip(0), findopt.Limit(2))

	if err != nil {
		log.Println("collection.Find error : ", err)
		return
	}

	defer cursor.Close(context.TODO())

	// 遍历查询结果
	var logRecord LogRecord
	for cursor.Next(context.TODO()) {
		if err = cursor.Decode(&logRecord); err != nil {
			log.Println("cursor.Decode error : ", err)
			return
		}

		fmt.Println(logRecord)
	}
}
