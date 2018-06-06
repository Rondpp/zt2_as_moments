package redis
import (
        redis "github.com/garyburd/redigo/redis"
        log "github.com/jeanphorn/log4go"
        "conf"
        "proto"
)
func init() {
        pool = newRdsPool()
}

var (
        pool *redis.Pool
)

func newRdsPool() *redis.Pool {
        return &redis.Pool{
                MaxIdle: 80,
                MaxActive: 12000, // max number of connections
                Dial: func() (redis.Conn, error) {
                        c, err := redis.Dial("tcp", conf.GetCfg().RdsCfg.Server + ":" + conf.GetCfg().RdsCfg.Port)
                        if err != nil {
                                log.Debug(err)
                                panic(err.Error())
                        } else {
                                log.Debug("newRdsPool 成功")
                        }

                        return c, err
                },
        }

}

func HGetAll(query string) (map[string]string, int32) {
        session := pool.Get()
        defer session.Close()
 
        values, err := redis.StringMap(session.Do("HGETALL", query))
        if  err != nil {
                log.Debug(err)
                return nil, proto.ReturnCodeServerError
        }

        return values, proto.ReturnCodeOK
}
