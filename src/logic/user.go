package logic

import (
        "strconv"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "proto"
        "io/ioutil"
        "encoding/json"
)

const DEFAULT_BIRTHDAY_TIME = int64(333907200000) // 1980年8月1号0点0分0秒0毫秒
const DEFAULT_AVATAR_URL = "https://zt2as.oss-cn-hangzhou.aliyuncs.com/common/ic_avatar_default.png"

func UploadDefaultUserInfo(accid int64) error{
        session := mgohelper.GetSession()
        defer session.Close()

        coll := mgohelper.GetCollection(session, "users")
        selector := bson.M{"accid" : accid}

        data := bson.M{"$set":bson.M{"name":"英雄" + strconv.FormatUint(uint64(accid), 10), "birthday": DEFAULT_BIRTHDAY_TIME,"avatar":DEFAULT_AVATAR_URL}}
        _,err := coll.Upsert(selector, data)
        if err != nil {
                log.Error(err)
        }
        log.Debug("自动生成玩家信息:%d",accid)
        return err
}

func IsUserExists(accid int64) bool{

        session := mgohelper.GetSession()
        defer session.Close()

        coll := mgohelper.GetCollection(session, "users")
        selector := bson.M{"accid" : accid}

        var  data interface{}
        err := coll.Find(selector).One(data)
        if err != nil  {
                log.Error(err)
                return false
        }

        return true
}

func GetUserInfoRet(my_accid int64, accid int64) *proto.UserInfoRet{

        session := mgohelper.GetSession()
        defer session.Close()

        coll := mgohelper.GetCollection(session, "users")
        selector := bson.M{"accid" : accid}

        var req UserInfoMgo
        err := coll.Find(selector).One(&req)
        if err != nil  {
                log.Error(err)
                return nil
        }

        var  rsp proto.UserInfoRet
        UserInfoMgoToRet(my_accid, &req, &rsp)

        return &rsp

}

func  GetUserInfoRsp(r *http.Request) (*proto.UserInfoRet, int) {
        if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                return nil, proto.ReturnCodeMissParm
        } 

        my_accid    := GetMyAccID(r)
        query_accid := GetInt64UrlParmByName(r, "accid")

        log.Debug("查询个人信息,my_accid:%d,query_accid:%d", my_accid, query_accid)

        var rsp *proto.UserInfoRet
        rsp = GetUserInfoRet(my_accid, query_accid)
        if rsp == nil  {
                if my_accid == query_accid {
                        // 没有自己的信息,上传默认的
                        check_token := CheckToken(r)
                        if check_token  != proto.ReturnCodeOK {
                                return nil, check_token
                        }
                        err := UploadDefaultUserInfo(my_accid)
                        if err != nil {
                                return nil, proto.ReturnCodeServerError
                        }

                        rsp = GetUserInfoRet(my_accid, query_accid)
                        if rsp == nil {
                                return nil, proto.ReturnCodeServerError
                        }
                } else {
                        return nil, proto.ReturnCodeOK
                }
        }
        return rsp, proto.ReturnCodeOK
}

func UpdateUserInfoRsp(r *http.Request) (*proto.UserInfoRet, int) {
        body, err := ioutil.ReadAll(r.Body)
        if err != nil {
                log.Debug("read body err, %v\n", err)
                return nil, proto.ReturnCodeMissParm
        }

        var req proto.UpdateUserInfoReq
        json_err := json.Unmarshal(body, &req)
        if json_err != nil {
                log.Error(json_err)
                return nil,proto.ReturnCodeMissParm
        }

        my_accid    := GetMyAccID(r)
        session := mgohelper.GetSession()
        defer session.Close()

        coll := mgohelper.GetCollection(session, "users")

        selector := bson.M{"accid":my_accid}
        data     := bson.M{"$set":bson.M{"avatar":req.Avatar, "sex":req.Sex, "birthday":req.Birthday, "name":req.Name}}
        _,mgo_err := coll.Upsert(selector, data)
        if mgo_err != nil {
                log.Error(mgo_err)
                return nil,proto.ReturnCodeServerError
        }

        rsp := GetUserInfoRet(0, my_accid)
        if rsp == nil {
                return nil, proto.ReturnCodeServerError
        }
        log.Debug("修改个人信息:accid:%d,avatar:%s,sex:%d,birthday:%d,name:%s", my_accid, rsp.Avatar, rsp.Sex, rsp.Birthday, rsp.Name)
        return rsp, proto.ReturnCodeOK
}
