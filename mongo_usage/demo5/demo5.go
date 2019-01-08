package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

// startTime小于某时间
// {"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

// {"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	var (
		client     *mongo.Client
		db         *mongo.Database
		collection *mongo.Collection
		delCond    *DeleteCond
		res        *mongo.DeleteResult
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

	// 4. 删除起始时间早于当前时间的所有日志
	delCond = &DeleteCond{
		BeforeCond: TimeBeforeCond{
			Before: time.Now().Unix(),
		},
	}
	// 删除
	if res, err = collection.DeleteMany(context.TODO(), delCond); err != nil {
		fmt.Println("delete error:", err)
		return
	}

	fmt.Println("删除的行数:", res.DeletedCount)
}
