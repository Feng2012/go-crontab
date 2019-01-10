package common

import "encoding/json"

// 定时任务
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//
type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func BuildResponse(errno int, msg string, data interface{}) ([]byte, error) {
	//
	var (
		resp Response
		rb   []byte
		err  error
	)
	resp.Errno = errno
	resp.Msg = msg
	resp.Data = data
	if rb, err = json.Marshal(resp); err != nil {
		return nil, err
	}
	return rb, nil
}
