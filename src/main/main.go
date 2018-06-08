package main 

import (
        log "github.com/jeanphorn/log4go"
        "net/http"
        "fmt"
        "flag"
        "conf"
        "router"
)

func  main() {
        var cfg_name string
        flag.StringVar(&cfg_name, "c", "conf/conf.xml", "set confuration `file`")

        if conf.Init(cfg_name) {
                fmt.Printf("读取配置%s成功\n", cfg_name)
        } else {
                fmt.Printf("读取配置%s失败\n", cfg_name)
                return
        }

        log.LoadConfiguration(conf.GetCfg().LogCfgName, "xml")

        log.Info("初始化log成功")

        router.InitRouter()

        defer log.Close()

        http.ListenAndServe(conf.GetCfg().ServerCfg.Host + ":" + conf.GetCfg().ServerCfg.Port, nil)
}

