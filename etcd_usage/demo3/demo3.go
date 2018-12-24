package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config  clientv3.Config
		client  *clientv3.Client
		err     error
		kv      clientv3.KV
		getResp *clientv3.GetResponse
	)
	// 客户端配置
	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	// 建立客户端
	if client, err = clientv3.New(config); err != nil {
		fmt.Println("clientv3 new err: ", err)
		return
	}
	// 用于读写etcd键值对
	kv = clientv3.NewKV(client)

	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job1"); err != nil {
		fmt.Println("get err:", err)
		return
	}
	// fmt.Println(getResp.Kvs)
	fmt.Println(getResp.Kvs, getResp.Count)
}
