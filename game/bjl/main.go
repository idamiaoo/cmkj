package main

import (
	"go/cmkj_server_go/conf"
	"go/cmkj_server_go/models"
	//"go/cmkj_server_go/network"
	"go/cmkj_server_go/util"

	//"strconv"

	"github.com/gin-gonic/gin"
)

func init() {
	conf.Loadconf("bjl.conf")
}

func main() {
	models.InitDb()
	util.InitLog("bjl.log", "DEBUG")
	Bjl = NewBjl()
	name := conf.Conf.String("name")
	util.Log.Debug(name)

	r := gin.Default()
	r.GET("/bjl/game", handleClient)
	r.GET("/", handleClient)
	r.Run(":3001")
}
