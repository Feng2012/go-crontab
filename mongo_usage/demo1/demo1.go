package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

func main() {
	var (
		client *mongo.Client
		db     *mongo.Database
		err    error
	)
	// 1.建立连接
	if client, err = mongo.Connect(context.TODO(), "mongodb://127.0.0.1:27017", options.Client().SetConnectTimeout(5*time.Second)); err != nil {
		fmt.Println("connect error:", err)
		return
	}
	// 2.选择数据库
	db = client.Database("my_db")

	// 3.选择my_collection
	db.Collection("my_collection")
}
