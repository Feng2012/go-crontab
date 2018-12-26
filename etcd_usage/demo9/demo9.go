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
		kv             clientv3.KV
		lease          clientv3.Lease
		leaseGrantResp *clientv3.LeaseGrantResponse
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		keepResp       *clientv3.LeaseKeepAliveResponse
		txnResp        *clientv3.TxnResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	if client, err = clientv3.New(config); err != nil {
		fmt.Println("clientv3 new err: ", err)
		return
	}

	// 1. 上锁(创建租约，自动续租，拿着租约抢占一个key)
	lease = clientv3.NewLease(client)

	// 创建一个5s租约
	if leaseGrantResp, err = lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println("lease grant err:", err)
		return
	}

	// 租约ID
	leaseId := leaseGrantResp.ID

	// 准备一个用于取消自动续租的context
	ctx, cancelFunc := context.WithCancel(context.TODO())

	// 程序退出后，停止续租，取消租约
	defer cancelFunc()
	defer lease.Revoke(context.TODO(), leaseId)

	// 5s后自动取消续租
	if keepRespChan, err = lease.KeepAlive(ctx, leaseId); err != nil {
		fmt.Println("lease keep alive err:", err)
		return
	}

	// 处理合约应答的协程
	go func() {
		for {
			select {
			case keepResp = <-keepRespChan:
				if keepRespChan == nil {
					fmt.Println("租约已经失效")
					return
				} else { // 每秒会续租一次，就会收到续租应答
					fmt.Println("收到自动续租应答:", keepResp.ID)
				}
			}
		}
	}()

	// if 不存在key，then设置它，else抢锁失败
	kv = clientv3.NewKV(client)
	// 创建事务
	txn := kv.Txn(context.TODO())
	//
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/jobs/job9", "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/jobs/job9"))
	// 提交事务
	if txnResp, err = txn.Commit(); err != nil {
		fmt.Println("txn commit err:", err)
		return
	}
	// 没抢到锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}
	fmt.Println("处理任务")
	time.Sleep(10 * time.Second)

	// 2. 处理业务

	// 3. 释放锁(取消自动续租，释放租约)
	// defer

}
