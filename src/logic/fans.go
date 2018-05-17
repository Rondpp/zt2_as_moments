package logic
import (
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "proto"
)

func GetUserFansRsp(r *http.Request) (*[]proto.UserInfoRet, int){
        var fans  FansInfoMgo
        query_accid := GetInt64UrlParmByName(r, "accid")
        my_accid    := GetMyAccID(r)

        sUsers   := mgohelper.GetCollection("users")
        selector := bson.M{"accid" : query_accid}
        err      := sUsers.Find(selector).Select(bson.M{"fans":1,"_id":0}).One(&fans)
        if err != nil  && err != mgo.ErrNotFound {
                log.Debug(err)
                return nil, proto.ReturnCodeServerError
        }

        log.Debug("获取玩家粉丝,accid:%d,follows:%d", query_accid, len(fans.Fans))

        var rsp []proto.UserInfoRet 
        for _, v := range fans.Fans {

                var userinfo_in_mgo UserInfoMgo
                sUsers.Find(bson.M{"accid": v}).One(&userinfo_in_mgo)
                if err != nil && err != mgo.ErrNotFound {
                        log.Debug(err)
                }

                var userinfo_ret  proto.UserInfoRet
                UserInfoMgoToRet(my_accid, &userinfo_in_mgo, &userinfo_ret)
                rsp = append(rsp, userinfo_ret)
        }
        return &rsp, proto.ReturnCodeOK
}
