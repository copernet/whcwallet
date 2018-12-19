package main

import (
	"github.com/copernet/whccommon/log"
	"github.com/copernet/whcwallet/config"
	"github.com/copernet/whcwallet/logic"
	"github.com/copernet/whcwallet/routers"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
)

func main() {
	// Setup server according to routers
	// now handle gin framework log
	log.InitLog(config.GetConf().Log)
	currentPath, err := filepath.Abs("./")
	if err != nil {
		panic("can not get current file abs path")
	}
	ginAccessLog := currentPath + log.DefaultLogDir + "access.log"
	ginLog, err := os.OpenFile(ginAccessLog, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0644))
	if err != nil {
		panic("gin log can not access: " + err.Error())
	}
	gin.DefaultWriter = ginLog

	//Dump headers
	go logic.DumpHeaders(log.NewContext(), nil)

	engine := routers.InitRouter()
	engine.Run(":9999")
}
