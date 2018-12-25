package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		config         clientv3.Config
		client         *clientv3.Client
		err            error
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId        clientv3.LeaseID
		putResp        *clientv3.PutResponse
		getResp        *clientv3.GetResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		kv             clientv3.KV
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

	// 申请一个lease（租约）
	lease = clientv3.NewLease(client)
	// 申请一个10s的租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println("lease grant err:", err)
		return
	}
	// 租约ID
	leaseId = leaseGrantResp.ID
	// 租约续租 续租5s后停止续租，总共15s的生命期
	// ctx, _ := context.WithTimeout(context.TODO(), 5 * time.Second)
	if keepRespChan, err = lease.KeepAlive(context.TODO(), leaseId); err != nil {
		fmt.Println("lease keep alive err:", err)
		return
	}
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					return
				} else {
					fmt.Println("收到自动续租应答")
				}
			}
		}
	}()

	// 获得kv对象
	kv = clientv3.NewKV(client)
	// Put一个kv，与租约关联，从而实现10s后自动过期
	if putResp, err = kv.Put(context.TODO(), "/cron/lock/job1", "", clientv3.WithLease(leaseId)); err != nil {
		fmt.Println("put err:", err)
		return
	}
	fmt.Println("写入成功:", putResp.Header.Revision)

	// 检查key是否过期
	for {
		if getResp, err = kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println("get err:", err)
			return
		}
		// 没有过期
		if getResp.Count == 0 {
			fmt.Println("kv过期")
			break
		}
		// 过期了
		fmt.Println("没有过期:", getResp.Kvs)
		time.Sleep(2 * time.Second)
	}
}
