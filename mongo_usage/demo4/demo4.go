package main

import (
	"context"
	"fmt"
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

// 过滤条件
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

func main() {
	var (
		client     *mongo.Client
		db         *mongo.Database
		collection *mongo.Collection
		cond       *FindByJobName
		cursor     mongo.Cursor
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

	// 4.查找(过滤条件+翻页参数)
	// 按照jobName字段过滤，找出jobName=job10的记录
	cond = &FindByJobName{JobName: "job10"}
	findOpt := options.Find()
	if cursor, err = collection.Find(context.TODO(), cond, findOpt.SetSkip(0), findOpt.SetLimit(2)); err != nil {
		fmt.Println("find error:", err)
		return
	}

	// 延迟释放游标
	defer cursor.Close(context.TODO())

	// 遍历结果
	for cursor.Next(context.TODO()) {
		record := &LogRecord{}
		// 反序列化bson到对象
		if err = cursor.Decode(record); err != nil {
			fmt.Println("decode error:", err)
			return
		}
		fmt.Println(*record)
	}

}
