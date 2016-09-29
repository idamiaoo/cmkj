package main

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/models"
	_ "go/cmkj_server_go/timer"
	"go/cmkj_server_go/util"

	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof" // 引用pprof package
	"syscall"

	"github.com/gin-gonic/gin"
)

var f *os.File

func wait() {
	ch := make(chan os.Signal, 1)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGHUP,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	for {
		sig := <-ch
		util.Log.Debug("signal:", sig.String())
		switch sig {
		case syscall.SIGTERM, syscall.SIGINT:
			pprof.StopCPUProfile()
			return
		default:
			return
		}
	}
}

func main() {

	f, _ = os.Create("profile_file")
	pprof.StartCPUProfile(f) // 开始cpu profile，结果写到文件f中
	/*
		defer func() {
			util.Log.Debug("center stop")
			pprof.StopCPUProfile() // 结束profile
		}()
	*/
	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		http.ListenAndServe(":6060", nil)

	}()
	util.InitLog("login.log", "DEBUG")
	conf.Loadconf("login.conf")
	name := conf.Conf.String("name")
	util.Log.Debug(name)
	models.InitDb()
	DefaultCenter = NewCenter()
	DefaultCenter.Start()

	gin.SetMode("debug") //test debug release
	r := gin.Default()

	r.GET("/login", loginServer)
	r.GET("/", loginServer)

	go r.Run(":3000")
	wait()
}
