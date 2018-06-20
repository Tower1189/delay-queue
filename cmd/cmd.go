package cmd

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ouqiang/delay-queue/config"
	"github.com/ouqiang/delay-queue/delayqueue"
	"github.com/ouqiang/delay-queue/servers/http_server"
	"github.com/ouqiang/delay-queue/servers/grpc"
	"os/signal"
	"syscall"
)

// Cmd 应用入口Command
type Cmd struct{}

var (
	version    bool
	configFile string
)

const (
	// AppVersion 应用版本号
	AppVersion = "0.4"
)

// Run 运行应用
func (cmd *Cmd) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// 解析命令行参数
	cmd.parseCommandArgs()
	if version {
		fmt.Println(AppVersion)
		os.Exit(0)
	}
	// 初始化配置
	config.Init(configFile)
	// 初始化队列
	delayqueue.Init()

	// 运行web server
	go cmd.runWeb()
	// 运行 grpc
	go cmd.runGrpc()
	<-c
}

// 解析命令行参数
func (cmd *Cmd) parseCommandArgs() {
	// 配置文件
	flag.StringVar(&configFile, "c", "", "./delay-delay-queue -c /path/to/delay-delay-queue.conf")
	// 版本
	flag.BoolVar(&version, "v", false, "./delay-delay-queue -v")
	flag.Parse()
}

// 运行Web Server
func (cmd *Cmd) runWeb() {
	http.HandleFunc("/push", http_server.Push)
	http.HandleFunc("/pop", http_server.Pop)
	http.HandleFunc("/finish", http_server.Delete)
	http.HandleFunc("/delete", http_server.Delete)
	http.HandleFunc("/get", http_server.Get)

	log.Printf("http server listen %s\n", config.Setting.BindAddress)
	err := http.ListenAndServe(config.Setting.BindAddress, nil)
	log.Fatalln(err)
}

func (cmd *Cmd) runGrpc() {
	log.Printf("grpc server listen %s\n", config.Setting.GrpcBindAddress)
	grpc.Run(config.Setting.GrpcBindAddress)
}