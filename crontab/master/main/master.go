package main

import (
	"flag"
	"fmt"
	"go-crontab/crontab/master"
	"runtime"
)

var (
	confFile string // 配置文件路径
)

// 初始化线程
func initEnv() {
	// 设置线程数量与CPU核数相等
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// 解析命令行参数
func initArgs() {
	// master -config ./master.json
	// mater -h
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

func main() {
	var (
		err error
	)
	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 启动api http服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	return

ERR:
	fmt.Println(err)
}
