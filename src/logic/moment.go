package logic

import (
        "proto"
        "net/http"
        "encoding/json"
        "io/ioutil"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "conf"
        "util"
)
func GetCommentNumByID(moment_id string) int {
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"_id":  bson.ObjectIdHex(moment_id)}

        type CommentInfo struct {
                Num     int `bson:"comment_num"`
        }
        var commentinfo CommentInfo
        err_comment_info := sMoments.Find(selector).Select(bson.M{"accid":1,"comment_num":1,"_id":0}).One(&commentinfo)
        if err_comment_info != nil && err_comment_info != mgo.ErrNotFound {
                log.Error(err_comment_info)
        }
        return commentinfo.Num
}

func GetMomentByID(moment_id string) *MomentMgo {
        var moment_mgo MomentMgo
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")

        err := sMoments.Find(bson.M{"_id":bson.ObjectIdHex(moment_id)}).One(&moment_mgo)
        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        return &moment_mgo
}

func GetUserMoments(accid int64, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var moment_mgo_list []MomentMgo

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"accid" : accid}

        if start_id != "" {
                selector = bson.M{"accid" : accid, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
        }

        err      := sMoments.Find(selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        return &moment_mgo_list
}

func GetNotVideoMoments(sort_type int, start_id string, limit_num int) *[]MomentMgo {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var moment_mgo_list []MomentMgo

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"video":bson.M{"$exists":false}}

        if start_id != "" {
                comment_num := GetCommentNumByID(start_id)
                selector = bson.M{"video":bson.M{"$exists":false}, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "comment_num":bson.M{"lte" :comment_num}}
        }
        log.Debug(selector)
        var err error
        if sort_type == 1 {
                err = sMoments.Find(selector).Sort("-comment_num", "-time").Limit(limit_num).All(&moment_mgo_list)
        } else {

                err = sMoments.Find(selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)
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
        selector := bson.M{"video":bson.M{"$exists":true}}

        if start_id != "" {
                comment_num := GetCommentNumByID(start_id)
                selector = bson.M{"video":bson.M{"$exists":true}, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "comment_num":bson.M{"lte" :comment_num}}
        }

        var err error
        if sort_type == 1 {
                err = sMoments.Find(selector).Sort("-comment_num", "-time").Limit(limit_num).All(&moment_mgo_list)
        } else {

                err = sMoments.Find(selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)
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
        moment_selector     := bson.M{"accid":bson.M{"$in":follow_mgo.Follows}}
        if start_id != "" {
                moment_selector = bson.M{"accid":bson.M{"$in":follow_mgo.Follows}, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
        }

        err_moment          := sMoments.Find(moment_selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)

        if err_moment != nil && err_moment != mgo.ErrNotFound {
                log.Error(err_moment)
                return nil
        }

        return &moment_mgo_list
}

func UploadMomentRsp(r *http.Request) (*[]proto.MomentRet, int)  {

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
        my_accid    := GetMyAccID(r)

        var moments MomentMgo
        moments.AccID = my_accid
        moments.Content = req.Content
        moments.Pic = req.Pic
        moments.Video = req.Video
        moments.ID = bson.NewObjectId()
        moments.Time = util.GetTimestamp()
        moments.CommentNum = 0 
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
        moment_mgo := GetMomentByID(moments.ID.Hex())
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
        start_id  := GetStartID(r)
 
        log.Debug("查询动态:my_accid:%d,limit_num:%d,start_id:%s", my_accid, limit_num, start_id)
        var moment_mgo_list *[]MomentMgo

        if vars["type"][0] == "0" {
                //查询某个人的
                log.Debug("查询某个人的动态")

                if CheckUrlParm(r , "accid") != proto.ReturnCodeOK {
                        return nil, proto.ReturnCodeMissParm
                }

                query_accid := GetInt64UrlParmByName(r, "accid")
                moment_mgo_list = GetUserMoments(query_accid, start_id, limit_num)
        } else if vars["type"][0] == "1" {
                //查询最新动态
                log.Debug("查询最新动态111")
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

                moment_mgo := GetMomentByID(vars["moment_id"][0])
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

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")

        remove_err := sMoments.Remove(bson.M{"_id":bson.ObjectIdHex(vars["moment_id"][0]),"accid":my_accid})
        if remove_err != nil {
                log.Error(remove_err)
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
