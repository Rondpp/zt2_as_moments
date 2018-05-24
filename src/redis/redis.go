package redis
import (
        redis "github.com/garyburd/redigo/redis"
        log "github.com/jeanphorn/log4go"
        "conf"
        "proto"
//        "fmt"
)
func init() {
        Pool = newRdsPool()//GetPool()
}

var (
        Pool *redis.Pool// = newRdsPool()//GetPool()
)

func newRdsPool() *redis.Pool {
        return &redis.Pool{
                MaxIdle: 80,
                MaxActive: 12000, // max number of connections
                Dial: func() (redis.Conn, error) {
                        c, err := redis.Dial("tcp", "localhost:27017")
                        if err != nil {
//                                fmt.Println(err.Error())
                                panic(err.Error())
                        }
                        return c, err
                },
        }

}

func GetPool() *redis.Pool {
        return &redis.Pool{
                MaxIdle     :   conf.GetCfg().RdsCfg.MaxIdle,
                MaxActive   :   conf.GetCfg().RdsCfg.MaxActive, // max number of connections
                Dial        :   func() (redis.Conn, error) {
                                    //c, err := redis.Dial("tcp", conf.GetCfg().RdsCfg.Server + ":"  conf.GetCfg().RdsCfg.Port)
                                    c, err := redis.Dial("tcp", "localhost:27017")
                                    if err != nil {
                                        panic(err.Error())
                                    }
                                    return c, err
                                },
        }
}

func HGetAll(query string) (map[string]string, int32) {
        session := Pool.Get()
        defer session.Close()
 
        values, err := redis.StringMap(session.Do("HGETALL", query))
        if  err != nil {
                log.Debug(err)
                return nil, proto.ReturnCodeServerError
        }

        return values, proto.ReturnCodeOK
}
