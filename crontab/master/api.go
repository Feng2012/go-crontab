package master

import (
	"encoding/json"
	"go-crontab/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
)

// 任务http接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	ApiSvr *ApiServer
)

// 保存任务
// POST job={"name":"job1", "command":"echo hello", "cronExpr":"*****"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	// 保存任务到etcd中
	var (
		postedJob string
		err       error
		newJob    *common.Job
		oldJob    *common.Job
		bytes     []byte
	)
	// 解析POST表单
	if err = r.ParseForm(); err != nil {
		goto ERR
	}
	// 取表单中的job字段
	postedJob = r.PostForm.Get("job")
	// 反序列化
	if err = json.Unmarshal([]byte(postedJob), newJob); err != nil {
		goto ERR
	}
	// 保存任务到etcd
	if oldJob, err = JM.SaveJob(newJob); err != nil {
		goto ERR
	}
	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}

	return

ERR:
	// 异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

// 初始化服务
func InitApiServer() error {
	var (
		mux      *http.ServeMux
		listener net.Listener
		server   *http.Server
		err      error
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	// 启动TCP监听
	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(Cfg.ApiPort)); err != nil {
		return err
	}

	// 创建http服务
	server = &http.Server{
		ReadTimeout:  time.Duration(Cfg.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(Cfg.ApiWriteTimeout) * time.Millisecond,
		Handler:      mux,
	}

	ApiSvr = &ApiServer{
		httpServer: server,
	}

	// 启动服务
	go server.Serve(listener)

	return nil
}
