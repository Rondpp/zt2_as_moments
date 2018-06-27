package main

import (
	"conf"
	log "github.com/jeanphorn/log4go"
	"net/http"
	"router"
)

func main() {
	log.LoadConfiguration(conf.GetCfg().LogCfgName, "xml")

	log.Info("初始化log成功")

	router.InitRouter()

	defer log.Close()

	http.ListenAndServe(conf.GetCfg().ServerCfg.Host+":"+conf.GetCfg().ServerCfg.Port, nil)
}
