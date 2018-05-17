package mgohelper

import (
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "conf"
)
var (
        Session *mgo.Session
)

func GetSession() *mgo.Session {
        if Session == nil {
                var err error
                Session, err = mgo.Dial(conf.GetCfg().MgoCfg.Server + ":" + conf.GetCfg().MgoCfg.Port)
                if err != nil {
                        log.Error(err)
                        panic(err) //直接终止程序运行
                }
        }
        //最大连接池默认为4096
        return Session.Clone()
}

func GetCollection(col string) *mgo.Collection {
        return GetSession().DB(conf.GetCfg().MgoCfg.DB).C(col)
}
