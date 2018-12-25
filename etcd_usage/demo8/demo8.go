package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		opResp clientv3.OpResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println("clientv3 new err: ", err)
		return
	}

	kv = clientv3.NewKV(client)

	// 创建op
	putOp := clientv3.OpPut("/cron/jobs/job8", "job8")
	// 执行op
	if opResp, err = kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println("kv do putOp err:", err)
		return
	}
	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	getOp := clientv3.OpGet("/cron/jobs/job8")
	if opResp, err = kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println("kv do getOp err:", err)
		return
	}
	fmt.Println("数据Revision:", opResp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value:", string(opResp.Get().Kvs[0].Value))
}
