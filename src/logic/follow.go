package logic

import (
        "proto"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
)

func UploadFollowRsp(r *http.Request) (int)  {
        if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                return proto.ReturnCodeMissParm
        }

        my_accid        := GetMyAccID(r)
        accid           := GetInt64UrlParmByName(r, "accid")

        session := mgohelper.GetSession()
        defer session.Close()

        sUsers          := mgohelper.GetCollection(session, "users")
        if my_accid == accid {
                return proto.ReturnCodeParmWrong
        }

        fans_selector    := bson.M{"accid":accid}
        fans_data        := bson.M{"$addToSet":bson.M{"fans":my_accid}}
        _,fans_err       := sUsers.Upsert(fans_selector, fans_data);
        if fans_err != nil {
                log.Debug(fans_err)
                return proto.ReturnCodeServerError
        }

        follow_selector   := bson.M{"accid":my_accid}
        follow_data       := bson.M{"$addToSet":bson.M{"follows":accid}}
        _,follow_err      := sUsers.Upsert(follow_selector, follow_data);
        if follow_err != nil {
                log.Debug(follow_err)
                return proto.ReturnCodeServerError
        }

        log.Debug("玩家关注,accid:%d,follow:%d", my_accid, accid)
        return  proto.ReturnCodeOK
}

func  GetFollowRsp(r *http.Request) (*[]proto.UserInfoRet, int) {
        if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                return nil, proto.ReturnCodeMissParm
        }

        my_accid    := GetMyAccID(r)
        accid       := GetInt64UrlParmByName(r, "accid")

        session := mgohelper.GetSession()
        defer session.Close()

        sUsers      := mgohelper.GetCollection(session, "users")


        var follows FollowsInfoMgo
        selector := bson.M{"accid":accid}
        err := sUsers.Find(selector).Select(bson.M{"follows":1,"_id":0}).One(&follows)
        if err != nil && err != mgo.ErrNotFound {
                log.Debug(err)
                return nil, proto.ReturnCodeServerError
        }

        log.Debug("获取玩家关注的人,accid:%d,follows:%d", accid, len(follows.Follows))

        var rsp []proto.UserInfoRet 
        for _, v := range follows.Follows {

                var userinfo_mgo UserInfoMgo
                sUsers.Find(bson.M{"accid": v}).One(&userinfo_mgo)
                if err != nil {
                        log.Debug(err)
                }

                var userinfo_ret  proto.UserInfoRet
                UserInfoMgoToRet(my_accid, &userinfo_mgo, &userinfo_ret)
                rsp = append(rsp, userinfo_ret)
        }
        return &rsp, proto.ReturnCodeOK
}

func  DeleteFollowRsp(r *http.Request) (int) {
        if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                return proto.ReturnCodeMissParm
        }

        my_accid    := GetMyAccID(r)
        query_accid := GetInt64UrlParmByName(r, "accid")

        session := mgohelper.GetSession()
        defer session.Close()

        sUsers := mgohelper.GetCollection(session, "users")

        fans_selector    := bson.M{"accid":query_accid}
        fans_data        := bson.M{"$pull":bson.M{"fans":my_accid}}
        _,fans_err       := sUsers.Upsert(fans_selector, fans_data);
        if fans_err != nil {
                return proto.ReturnCodeServerError
        }

        follow_selector    := bson.M{"accid":my_accid}
        follow_data        := bson.M{"$pull":bson.M{"follows":query_accid}}
        _,follow_err       := sUsers.Upsert(follow_selector, follow_data);
        if follow_err != nil {
                return proto.ReturnCodeServerError
        }
        log.Debug("玩家取消关注,accid:%d,follow:%d", my_accid, query_accid)
        return  proto.ReturnCodeOK
}
