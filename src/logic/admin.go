package logic

import (
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "proto"
        "io/ioutil"
        "encoding/json"
        "conf"
        "util"
)

func  GetAdminUserListRsp(r *http.Request) (*[]proto.UserInfoRet, int) {
        coll := mgohelper.GetCollection("users")
        selector := bson.M{"permission":bson.M{"$gt": 0}}

        var userinfo_mgo_list []UserInfoMgo
        err := coll.Find(selector).All(&userinfo_mgo_list)
        if err != nil  && err != mgo.ErrNotFound {
                log.Error(err)
                return nil, proto.ReturnCodeServerError
        }

        var  rsp []proto.UserInfoRet
        for _,v := range userinfo_mgo_list {
                var userinfo_ret proto.UserInfoRet
                UserInfoMgoToRet(0, &v, &userinfo_ret)
                rsp = append(rsp, userinfo_ret)
        }
        log.Debug("管理员:%d", len(userinfo_mgo_list))
        return &rsp, proto.ReturnCodeOK
}

func UploadAdminUserRsp(r *http.Request) (int)  {

        body, body_err := ioutil.ReadAll(r.Body)
        if body_err != nil {
                log.Debug("body err:%s",body_err)
                return proto.ReturnCodeMissParm
        }

        var req proto.AdminUserPermissionSetReq
        json_err := json.Unmarshal(body, &req) 
        if json_err != nil {
                log.Debug("json err:%s", json_err)
                return proto.ReturnCodeMissParm
        }

        selector    := bson.M{"accid" : req.AccID}
        data        := bson.M{"$set":bson.M{"permission":req.Permission}}
        sUsers      := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        _, err      := sUsers.Upsert(selector, data);

        if err != nil {
                log.Error(err)
                return proto.ReturnCodeServerError
        }

        log.Debug("设置玩家权限,accid:%d,permission:%d", req.AccID, req.Permission)
        return  proto.ReturnCodeOK
}

func UploadAdminForbiddenRsp(r *http.Request) (int)  {

        body, body_err := ioutil.ReadAll(r.Body)
        if body_err != nil {
                log.Debug("body err:%s",body_err)
                return proto.ReturnCodeMissParm
        }

        var req proto.AdminForbiddenReq
        json_err := json.Unmarshal(body, &req) 
        if json_err != nil {
                log.Debug("json err:%s", json_err)
                return proto.ReturnCodeMissParm
        }

        selector    := bson.M{"accid" : req.AccID}
        data        := bson.M{"$set":bson.M{"forbidden":req.Time}}
        sUsers      := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        _, err      := sUsers.Upsert(selector, data);

        if err != nil {
                log.Error(err)
                return proto.ReturnCodeServerError
        }

        log.Debug("禁言,accid:%d,time:%d", req.AccID, req.Time)
        return  proto.ReturnCodeOK
}

func UploadAdminToTopRsp(r *http.Request) (int)  {
        moment_id   := GetObjectIDByName(r, "moment_id")
        if moment_id == "" {
                return proto.ReturnCodeMissParm
        }

        selector    := bson.M{"_id" :bson.ObjectIdHex(moment_id)}
        data        := bson.M{"$set":bson.M{"to_top_time":util.GetTimestamp()}}
        sMoments    := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        _, err      := sMoments.Upsert(selector, data);

        if err != nil {
                log.Error(err)
                return proto.ReturnCodeServerError
        }

        log.Debug("置顶,moment_id:%d", moment_id)
        return  proto.ReturnCodeOK
}

func UploadAdminDeleteRsp(r *http.Request) (int)  {
        moment_id   := GetObjectIDByName(r, "moment_id")
        comment_id  := GetObjectIDByName(r, "comment_id")
        data        := bson.M{"$set":bson.M{"valid":0}}

        var coll *mgo.Collection
        var selector interface{}
        if moment_id != "" {
                selector    = bson.M{"_id" :bson.ObjectIdHex(moment_id)}
                coll = mgohelper.GetCollection("moments")
                log.Debug("管理员删除动态,moment_id:%d", moment_id)

        } else if comment_id != "" {
                selector    = bson.M{"_id" :bson.ObjectIdHex(comment_id)}
                coll = mgohelper.GetCollection("comments")
                log.Debug("管理员删除评论,comment_id:%d", comment_id)
        } else {
                return proto.ReturnCodeMissParm
        }
        _, err := coll.Upsert(selector, data);

        if err != nil {
                log.Error(err)
                return proto.ReturnCodeServerError
        }
        return  proto.ReturnCodeOK
}
