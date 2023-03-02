package main

import (
	"flag"
	"fmt"
	"github.com/Serendipity-sw/gutil"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/swgloomy/gutil/glog"
	"os"
	"os/signal"
	"syscall"
	"wechat-chatting-service/wechat"
)

var (
	pidStrPath = "./wechat-chatting-service.pid"
	debugFlag  = flag.Bool("d", false, "debug mode")                        //是否为调试模式
	configFn   = flag.String("config", "./config.json", "config file path") //配置文件地址
	logsDir    string
)

func main() {
	flag.Parse()

	err := config.ReadCfg(*configFn)
	if err != nil {
		fmt.Printf("main ReadCfg read err! filePath: %s err: %+v \n", *configFn, err.Error())
		return
	}
	readConfig()
	serverRun(*debugFlag)

	wechat.WechatLogin()

	c := make(chan os.Signal, 1)
	gutil.WritePid(pidStrPath)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	//信号等待
	<-c
	fmt.Println("main exit application!")
	serverExit()
}

func serverRun(debug bool) {

	gutil.LogInit(debug, logsDir)

	gutil.SetCPUUseNumber(0)
	fmt.Println("set many cpu successfully!")

	deferinit.InitAll()
	fmt.Println("init all module successfully!")

	deferinit.RunRoutines()
	fmt.Println("init all run successfully!")
}

func serverExit() {
	deferinit.StopRoutines()
	fmt.Println("stop routine successfully!")

	deferinit.FiniAll()
	fmt.Println("stop all modules successfully!")

	glog.Close()

	os.Exit(0)
}
