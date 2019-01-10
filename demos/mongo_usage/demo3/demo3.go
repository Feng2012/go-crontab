package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

//
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

// 日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`
	Command   string    `bson:"command"`
	Err       string    `bson:"err"`
	Content   string    `bson:"content"`
	TimePoint TimePoint `bson:"timePoint"`
}

func main() {
	var (
		client     *mongo.Client
		db         *mongo.Database
		collection *mongo.Collection
		res        *mongo.InsertManyResult
		err        error
	)

	// 1.建立连接
	if client, err = mongo.Connect(context.TODO(), "mongodb://127.0.0.1:27017", options.Client().SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println("connect error:", err)
		return
	}

	// 2.选择(或创建)数据库
	db = client.Database("cron")

	// 3.选择(或创建)collection
	collection = db.Collection("log")

	// 4.记录(bson)
	record := &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}

	// 5.批量插入多条document
	logArr := []interface{}{record, record, record}
	if res, err = collection.InsertMany(context.TODO(), logArr); err != nil {
		fmt.Println("insert many error:", err)
		return
	}

	for _, insertId := range res.InsertedIDs {
		docId := insertId.(primitive.ObjectID)
		fmt.Println("自增ID:", docId.Hex())
	}
}
