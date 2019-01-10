package master

import (
	"context"
	"encoding/json"
	"go-crontab/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"time"
)

// 任务管理器
type JobManager struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	JM *JobManager
)

// 初始化管理器
func InitJobManager() error {
	var (
		client *clientv3.Client
		config clientv3.Config
		kv     clientv3.KV
		lease  clientv3.Lease
		err    error
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   Cfg.EtcdEndpoints,
		DialTimeout: time.Duration(Cfg.EtcdDialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return err
	}

	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	//
	JM = &JobManager{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return nil
}

// 保存任务
func (jm *JobManager) SaveJob(newJob *common.Job) (*common.Job, error) {
	var (
		jobKey   string
		jobValue []byte
		err      error
		putResp  *clientv3.PutResponse
		oldJob   *common.Job
	)
	// 任务的key值
	jobKey = "/cron/jobs/" + newJob.Name

	// 序列化newJob为json
	if jobValue, err = json.Marshal(newJob); err != nil {
		return nil, err
	}

	// 保存到etcd
	if putResp, err = JM.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return nil, err
	}
	// 如果是是更新，返回旧值
	if putResp.PrevKv != nil {
		// 反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJob); err != nil {
			return nil, nil
		}

		return oldJob, nil
	}

	return nil, nil
}
