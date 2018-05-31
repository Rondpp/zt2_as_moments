package logic

import (
        "proto"
        "net/http"
        "fmt"
        "encoding/json"
        "io/ioutil"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "conf"
        "util"
        "errors"
        "time"
)

func CheckForbidden(r *http.Request) (int, error) {
        my_accid    := GetMyAccID(r)
        sUsers   := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        selector   := bson.M{"accid": my_accid}

        type ForbiddenInfo struct {
                ForbiddenStartTime    int64 `bson:"forbidden_start_time"`
                ForbiddenLastTime     int64 `bson:"forbidden_last_time"`
        }
        var forbiddeninfo ForbiddenInfo
        err := sUsers.Find(selector).Select(bson.M{"forbidden_start_time":1, "forbidden_last_time":1,"_id":0}).One(&forbiddeninfo)
        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
        }
        now := util.GetTimestamp()
        stop_forbidden_time := forbiddeninfo.ForbiddenLastTime + forbiddeninfo.ForbiddenStartTime
        if  (forbiddeninfo.ForbiddenLastTime > 0) && (now < stop_forbidden_time) {
                err_msg := fmt.Sprintf("您被禁言%s,将于%s后解除", util.FormatTimeCH(forbiddeninfo.ForbiddenLastTime),time.Unix(stop_forbidden_time/1000, 0).Format("2006-01-02 15:04:05"))
                return proto.ReturnCodeForbidden, errors.New(err_msg)
        }
        return proto.ReturnCodeOK, nil
}

func GetCommentNumByID(moment_id string) int {
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"_id":  bson.ObjectIdHex(moment_id)}

        type MomentInfo struct {
                Num     int `bson:"comment_num"`
        }
        var momentinfo MomentInfo
        err := sMoments.Find(selector).Select(bson.M{"accid":1,"comment_num":1,"_id":0}).One(&momentinfo)
        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
        }
        return momentinfo.Num
}

func GetMomentOwnerByID(moment_id string) int64 {
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"_id":  bson.ObjectIdHex(moment_id)}

        type MomentInfo struct {
                AccID     int64 `bson:"accid"`
        }
        var momentinfo MomentInfo
        err := sMoments.Find(selector).Select(bson.M{"accid":1,"_id":0}).One(&momentinfo)
        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
        }
        return momentinfo.AccID
}

func GetMomentByID(moment_id string, self bool) *MomentMgo {
        var moment_mgo MomentMgo
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        var err error
        if self {
                err = sMoments.Find(bson.M{"_id":bson.ObjectIdHex(moment_id), "$or":[]bson.M{bson.M{"valid":proto.ValidOK},bson.M{"valid":proto.ValidWaitForCheck},bson.M{"valid":proto.ValidDeleteByAdmin}}}).One(&moment_mgo)
        } else {
                err = sMoments.Find(bson.M{"_id":bson.ObjectIdHex(moment_id), "valid":proto.ValidOK}).One(&moment_mgo)
        }
        if err != nil {
                log.Error(err)
                return nil
        }
        return &moment_mgo
}

func GetUserMoments(my_accid int64, query_accid int64, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var moment_mgo_list []MomentMgo

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")

        var selector interface {}
        if start_id != "" {
                if my_accid == query_accid  {
                        selector = bson.M{"accid" : query_accid, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "$or":[]bson.M{bson.M{"valid":proto.ValidOK},bson.M{"valid":proto.ValidWaitForCheck}}}
                } else {
                        selector = bson.M{"accid" : query_accid, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "valid":proto.ValidOK}
                }

        } else {
                if my_accid == query_accid {
                        selector = bson.M{"accid" : query_accid, "$or":[]bson.M{bson.M{"valid":proto.ValidOK},bson.M{"valid":proto.ValidWaitForCheck}}}
                } else {
                        selector = bson.M{"accid" : query_accid, "valid":proto.ValidOK}
                }
        }

        err      := sMoments.Find(selector).Sort("-to_top_time", "-time").Limit(limit_num).All(&moment_mgo_list)

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        log.Debug(selector)
        return &moment_mgo_list
}

func GetNotVideoMoments(sort_type int, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var moment_mgo_list []MomentMgo

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"video":bson.M{"$exists":false}, "valid":proto.ValidOK}

        if start_id != "" {
                if sort_type == 1 {
                        comment_num := GetCommentNumByID(start_id)
                        selector = bson.M{"video":bson.M{"$exists":false}, "valid":proto.ValidOK, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "comment_num":bson.M{"$lte" :comment_num}}
                } else {
                        selector = bson.M{"video":bson.M{"$exists":false}, "valid":proto.ValidOK, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
                }
        }
        log.Debug(selector)
        var err error
        if sort_type == 1 {
                err = sMoments.Find(selector).Sort("-to_top_time","-comment_num", "-time").Limit(limit_num).All(&moment_mgo_list)
        } else {

                err = sMoments.Find(selector).Sort("-to_top_time","-time").Limit(limit_num).All(&moment_mgo_list)
        }

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        log.Debug("查询非视频的动态:start_id:%s,limit_num:%d,count:%d", start_id, limit_num, len(moment_mgo_list))
        return &moment_mgo_list
}

func GetVideoMoments(sort_type int, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var moment_mgo_list []MomentMgo

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"video":bson.M{"$exists":true}, "valid":proto.ValidOK}

        if start_id != "" {
                if sort_type == 1 {
                        comment_num := GetCommentNumByID(start_id)
                        selector = bson.M{"video":bson.M{"$exists":true},  "valid":proto.ValidOK, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "comment_num":bson.M{"$lte" :comment_num}}
                } else {
                        selector = bson.M{"video":bson.M{"$exists":true},  "valid":proto.ValidOK, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
                }
        }

        var err error
        if sort_type == 1 {
                err = sMoments.Find(selector).Sort("-to_top_time", "-comment_num", "-time").Limit(limit_num).All(&moment_mgo_list)
        } else {

                err = sMoments.Find(selector).Sort("-to_top_time","-time").Limit(limit_num).All(&moment_mgo_list)
        }

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        return &moment_mgo_list
}

func GetFollowUserMoments(my_accid int64, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var follow_mgo FollowsInfoMgo

        sFollows   := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        selector   := bson.M{"accid": my_accid}

        err_follow := sFollows.Find(selector).Select(bson.M{"follows":1,"_id":0}).One(&follow_mgo)

        if err_follow != nil && err_follow != mgo.ErrNotFound {
                log.Error(err_follow)
                return nil
        }
        log.Debug(follow_mgo.Follows)

        var moment_mgo_list []MomentMgo
        sMoments            := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        moment_selector     := bson.M{"accid":bson.M{"$in":follow_mgo.Follows}, "valid":proto.ValidOK}
        if start_id != "" {
                moment_selector = bson.M{"accid":bson.M{"$in":follow_mgo.Follows}, "valid":proto.ValidOK, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
        }

        err_moment          := sMoments.Find(moment_selector).Sort("-to_top_time", "-time").Limit(limit_num).All(&moment_mgo_list)

        if err_moment != nil && err_moment != mgo.ErrNotFound {
                log.Error(err_moment)
                return nil
        }

        return &moment_mgo_list
}

func UploadMomentRsp(r *http.Request) (*[]proto.MomentRet, int)  {
        log.Debug("玩家发布动态,method:%s,Host:%s,url:%s", r.Method, r.Host, r.URL)

        check_token := CheckToken(r)
        if check_token  != proto.ReturnCodeOK {
                return nil, check_token
        }

        body, body_err := ioutil.ReadAll(r.Body)
        if body_err != nil {
                log.Debug("body err:%s",body_err)
                return nil, proto.ReturnCodeMissParm
        }

        var req proto.PublishMomentReq
        json_err := json.Unmarshal(body, &req) 
        if json_err != nil {
                log.Debug("json err:%s", json_err)
                return nil, proto.ReturnCodeMissParm
        }

        if len(req.Content) == 0 {
                return nil, proto.ReturnCodeMissParm
        }

        my_accid    := GetMyAccID(r)
        if !IsUserExists(my_accid)  {
                UploadDefaultUserInfo(my_accid)
        }

        var moments MomentMgo
        moments.AccID = my_accid
        moments.Content = req.Content
        moments.Pic = req.Pic
        moments.Video = req.Video
        moments.ID = bson.NewObjectId()
        moments.Time = util.GetTimestamp()
        moments.CommentNum = 0 
        if len(req.Video) > 0 {
                moments.Valid = proto.ValidWaitForCheck
        } else {
                moments.Valid = proto.ValidOK
        }
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        err := sMoments.Insert(&moments)
        if err != nil {
                log.Error(err)
                return nil, proto.ReturnCodeServerError
        }

        selector := bson.M{"accid" : my_accid}
        data := bson.M{"$addToSet":bson.M{"moments":moments.ID}}
        sUsers := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        sUsers.Upsert(selector, data);

        if err != nil {
                log.Error(err)
                return nil,proto.ReturnCodeServerError
        }
        moment_mgo := GetMomentByID(moments.ID.Hex(), true)
        if moment_mgo == nil {
                return nil, proto.ReturnCodeServerError
        }

        var rsp []proto.MomentRet
        var moment_ret proto.MomentRet
        MomentMgoToRet(my_accid, moment_mgo, &moment_ret)
        rsp = append(rsp, moment_ret)

        log.Debug("玩家发布动态,method:%s,Host:%s,url:%s,body:%s:content:%s,len(pic):%d", r.Method, r.Host, r.URL, string(body),req.Content, req.Pic)
        return &rsp, proto.ReturnCodeOK
}

func  GetMomentRsp(r *http.Request) (interface {}, int) {
        vars := r.URL.Query();

        if CheckUrlParm(r , "type") != proto.ReturnCodeOK {
                return nil, proto.ReturnCodeMissParm
        }

        my_accid  := GetMyAccID(r)
        limit_num := GetIntUrlParmByName(r, "num")
        start_id  := GetObjectIDByName(r, "start_id")
 
        log.Debug("查询动态:my_accid:%d,limit_num:%d,start_id:%s", my_accid, limit_num, start_id)
        var moment_mgo_list *[]MomentMgo

        if vars["type"][0] == "0" {
                //查询某个人的
                log.Debug("查询某个人的动态")

                if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                        return nil, proto.ReturnCodeMissParm
                }

                query_accid := GetInt64UrlParmByName(r, "accid")
                moment_mgo_list = GetUserMoments(my_accid, query_accid, start_id, limit_num)
        } else if vars["type"][0] == "1" {
                //查询最新动态
                log.Debug("查询最新动态")
                sort_type := 0
                if len(vars["sort_type"]) > 0 {
                        if vars["sort_type"][0] == "1" {
                                sort_type = 1
                        }
                }
                moment_mgo_list = GetNotVideoMoments(sort_type, start_id, limit_num)

        } else if vars["type"][0] == "2" {
                // 查询最新视频
                log.Debug("查询最新视频")
                sort_type := 0
                if len(vars["sort_type"]) > 0 {
                        if vars["sort_type"][0] == "1" {
                                sort_type = 1
                        }
                }
                moment_mgo_list = GetVideoMoments(sort_type, start_id, limit_num)

        } else if vars["type"][0] == "3" {
                // 查询关注的人的动态
                log.Debug("查询关注的人的动态")

                moment_mgo_list = GetFollowUserMoments(my_accid, start_id, limit_num)
        } else if vars["type"][0] == "4"{
                if CheckUrlParm(r , "moment_id") != proto.ReturnCodeOK  || len(vars["moment_id"][0]) != 24{
                        return nil, proto.ReturnCodeMissParm
                }

                moment_mgo := GetMomentByID(vars["moment_id"][0], false)
                if moment_mgo == nil {
                        return nil, proto.ReturnCodeOK
                }

                var rsp proto.MomentRet
                MomentMgoToRet(my_accid, moment_mgo, &rsp)
                return &rsp, proto.ReturnCodeOK

        } else {
                return nil, proto.ReturnCodeMissParm
        }

        if moment_mgo_list == nil {
                return nil, proto.ReturnCodeServerError
        }

        var rsp []proto.MomentRet
        for _, v := range *moment_mgo_list {
                var moment_ret proto.MomentRet
                MomentMgoToRet(my_accid, &v, &moment_ret)
                rsp = append(rsp, moment_ret)
        }

        return &rsp, proto.ReturnCodeOK
}

func  DeleteMomentRsp(r *http.Request) (int) {
        if CheckUrlParm(r , "moment_id") != proto.ReturnCodeOK {
                return proto.ReturnCodeMissParm
        }

        my_accid    := GetMyAccID(r)
        vars        := r.URL.Query();

        sMoments        := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        moment_selector := bson.M{"_id":bson.ObjectIdHex(vars["moment_id"][0]),"accid":my_accid}
        moment_data     := bson.M{"$set":bson.M{"valid":proto.ValidDeleteByMe}}
        _, moment_err   := sMoments.Upsert(moment_selector, moment_data)
        if moment_err != nil {
                log.Error(moment_err)
                return proto.ReturnCodeServerError
        }

        selector    := bson.M{"accid" : my_accid}
        data        := bson.M{"$pull":bson.M{"moments":bson.ObjectIdHex(vars["moment_id"][0])}}
        sUsers      := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")

        _,upsert_err  := sUsers.Upsert(selector, data);
        if upsert_err != nil {
                log.Error(upsert_err)
                return proto.ReturnCodeServerError
        }
        return  proto.ReturnCodeOK
}


