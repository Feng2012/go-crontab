package master

import (
	"net"
	"net/http"
	"time"
)

// 任务http接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	ApiSvr *ApiServer
)

//
func handleJobSave(w http.ResponseWriter, r *http.Request) {

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
	if listener, err = net.Listen("tcp", ":8070"); err != nil {
		return err
	}

	// 创建http服务
	server = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}

	ApiSvr = &ApiServer{
		httpServer: server,
	}

	// 启动服务
	go server.Serve(listener)

	return nil
}
