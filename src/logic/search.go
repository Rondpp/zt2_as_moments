package logic
import (
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "proto"
        "conf"
)

func GetSearchMoments(my_accid int64, keyword string, start_id string, limit_num int) *[]proto.MomentRet {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }
        sMoments := mgohelper.GetCollection("moments")

        selector := bson.M{"content" : bson.M{"$regex":keyword}}
        if start_id != "" {
                selector = bson.M{"content" : bson.M{"$regex":keyword}, "_id": bson.M{"$gt": bson.ObjectIdHex(start_id)}}
        }

        var moment_mgo_list []MomentMgo
        err := sMoments.Find(selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)
        if err != nil  && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        var rsp []proto.MomentRet
        for _, v := range moment_mgo_list {
                var moment_ret proto.MomentRet
                MomentMgoToRet(my_accid, &v, &moment_ret)
                rsp = append(rsp, moment_ret)
        }

        log.Debug("搜索动态:keyword:%s,start_id:%s,limit_num:%d,num:%d", keyword, start_id, limit_num, len(moment_mgo_list))
        return &rsp
}

func GetSearchUsers(my_accid int64, keyword string, start_id string, limit_num int) *[]proto.UserInfoRet {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        sUsers := mgohelper.GetCollection("users")
        selector := bson.M{"name" : bson.M{"$regex":keyword}}
        if start_id != "" {
                selector = bson.M{"name" : bson.M{"$regex":keyword}, "_id": bson.M{"$gt": bson.ObjectIdHex(start_id)}}
        }

        var user_list []UserInfoMgo
        err := sUsers.Find(selector).Sort("-name").Limit(limit_num).All(&user_list)
        if err != nil  && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }

        var rsp []proto.UserInfoRet 
        for _, v := range user_list {
                var userinforet  proto.UserInfoRet
                UserInfoMgoToRet(my_accid, &v, &userinforet)
                rsp = append(rsp, userinforet)
        }

        log.Debug("搜索玩家,keyword:%s,start_id:%s,limit_num:%d,num:%d", keyword, start_id, limit_num, len(user_list))
        return &rsp
}

func GetSearchRsp(r *http.Request) (interface {}, int) {
        if CheckUrlParm(r, "type", "keyword") != proto.ReturnCodeOK {
                return nil, proto.ReturnCodeMissParm
        }

        vars        := r.URL.Query();
        my_accid    := GetMyAccID(r)
        limit_num   := GetIntUrlParmByName(r, "num")
        keyword     := vars["keyword"][0]
        start_id    := GetStartID(r)

        if vars["type"][0] == "1" {
                rsp := GetSearchMoments(my_accid, keyword, start_id, limit_num)
                if rsp == nil {
                        return nil, proto.ReturnCodeServerError
                } else {
                        return rsp, proto.ReturnCodeOK
                }
        } else if vars["type"][0] == "2" {
                rsp := GetSearchUsers(my_accid, keyword, start_id, limit_num)
                if rsp == nil {
                        return nil, proto.ReturnCodeServerError
                } else {
                        return rsp, proto.ReturnCodeOK
                }
        } else {
                return nil, proto.ReturnCodeServerError
        }
}
