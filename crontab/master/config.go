package master

import (
	"encoding/json"
	"io/ioutil"
)

// 配置
type Config struct {
	ApiPort         int      `json:"apiPort"`
	ApiReadTimeout  int      `json:"apiReadTimeout"`
	ApiWriteTimeout int      `json:"apiWriteTimeout"`
	EtcdEndpoints   []string `json:"etcdEndpoints"`
	EtcdDialTimeout int      `json:"etcdDialTimeout"`
}

var (
	Cfg *Config
)

func InitConfig(filename string) error {
	var (
		content []byte
		cfg     Config
		err     error
	)

	// 读配置文件
	if content, err = ioutil.ReadFile(filename); err != nil {
		return err
	}

	// 反序列化
	if err = json.Unmarshal(content, cfg); err != nil {
		return err
	}

	Cfg = &cfg

	return nil
}
