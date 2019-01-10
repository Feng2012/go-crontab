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
		delResp *clientv3.DeleteResponse
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
	// 删除
	if delResp, err = kv.Delete(context.TODO(), "/cron/jobs/job2", clientv3.WithPrevKV()); err != nil {
		fmt.Println("delete err:", err)
		return
	}
	// 删除前的值
	if len(delResp.PrevKvs) != 0 {
		for _, kvPair := range delResp.PrevKvs {
			fmt.Println("删除了:", string(kvPair.Key), string(kvPair.Value))
		}
	}
}
