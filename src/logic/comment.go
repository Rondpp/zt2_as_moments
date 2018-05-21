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

func IncMomentReadNum(moment_id string) {
        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"_id":  bson.ObjectIdHex(moment_id)}

        data := bson.M{"$inc": bson.M{"read_num":1}}
        _,err := sMoments.Upsert(selector, data)
        if err != nil {
                log.Error(err)
        }
}


func GetCommentByID(comment_id string) *CommentMgo {
        var comment_mgo CommentMgo

        sComment := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        selector  := bson.M{"_id" : bson.ObjectIdHex(comment_id)}
        err       := sComment.Find(selector).One(&comment_mgo)

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        return &comment_mgo
}

func GetMomentComment(my_accid int64, moment_id string, start_id string, limit_num int) *[]proto.CommentRet {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        log.Debug("获取动态的评论,moment_id:%s", moment_id)
        var comment_mgo_list []CommentMgo
        selector := bson.M{"moment_id":bson.ObjectIdHex(moment_id),"comment_id":bson.M{"$exists":false}}

        if start_id != "" {
                selector = bson.M{"moment_id":bson.ObjectIdHex(moment_id),"comment_id":bson.M{"$exists":false}, "_id": bson.M{"$gt": bson.ObjectIdHex(start_id)}}
        }

        sComment := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        err := sComment.Find(selector).Sort("time").Limit(limit_num).All(&comment_mgo_list)
        if err != nil {
                log.Error(err)
                return nil
        }

        var rsp []proto.CommentRet
        for _, v := range comment_mgo_list {
                var comment_ret proto.CommentRet
                CommentMgoToRet(my_accid, &v, &comment_ret)
                rsp = append(rsp, comment_ret)
        }

        return &rsp
}

func GetCommentComment(my_accid int64, comment_id string, start_id string, limit_num int) *[]proto.CommentCommentRet {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var comment_mgo_list []CommentMgo

        sComment := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        selector  := bson.M{"comment_id" : bson.ObjectIdHex(comment_id)}
        if start_id != "" {
                selector = bson.M{"comment_id":bson.ObjectIdHex(comment_id), "_id": bson.M{"$gt": bson.ObjectIdHex(start_id)}}
        }

        err       := sComment.Find(selector).Sort("time").Limit(limit_num).All(&comment_mgo_list)

        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }
        var rsp []proto.CommentCommentRet
        for _, v := range comment_mgo_list {
                var commentcomment_ret proto.CommentCommentRet
                CommentCommentMgoToRet(my_accid, &v, &commentcomment_ret)
                rsp = append(rsp, commentcomment_ret)
        }

        return &rsp
}

func UploadMomentComment(my_accid int64, req proto.CommentReq) int {
        log.Debug("玩家评论动态")

        sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
        selector := bson.M{"_id":  bson.ObjectIdHex(req.MomentID)}

        type CommentInfo struct {
                AccID   int64  `bson:"accid"`
                Num     uint32 `bson:"comment_num"`
        }
        var commentinfo CommentInfo
        err_comment_info := sMoments.Find(selector).Select(bson.M{"accid":1,"comment_num":1,"_id":0}).One(&commentinfo)
        if err_comment_info != nil && err_comment_info != mgo.ErrNotFound {
                log.Error(err_comment_info)
                return proto.ReturnCodeServerError
        }

        data := bson.M{"$inc": bson.M{"comment_num":1}}

        _,err_upsert := sMoments.Upsert(selector, data)
        if err_upsert != nil {
                log.Error(err_upsert)
                return proto.ReturnCodeServerError
        }

        var comment_mgo CommentMgo
        comment_mgo.MomentID        = bson.ObjectIdHex(req.MomentID)
        comment_mgo.Time            = util.GetTimestamp()
        comment_mgo.Content         = req.Content
        comment_mgo.AccID           = my_accid
        comment_mgo.ID              = bson.NewObjectId()
        comment_mgo.CommentedAccID  = commentinfo.AccID

        sComment := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        err_insert := sComment.Insert(&comment_mgo)
        if err_insert != nil {
                log.Error(err_insert)
                return proto.ReturnCodeServerError
        }

        return proto.ReturnCodeOK
}

func UploadCommentComment(my_accid int64, req proto.CommentReq) int  {
        log.Debug("玩家评论评论,accid:%d,comment_id:%s", my_accid,req.CommentID)

        sComments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        selector := bson.M{"_id":  bson.ObjectIdHex(req.CommentID)}

        type CommentInfo struct {
                MomentID    bson.ObjectId       `bson:"moment_id"`
                AccID       int64               `bson:"accid"`
                Num         uint32              `bson:"comment_num"`
        }
        var commentinfo CommentInfo
        err_comment_info := sComments.Find(selector).Select(bson.M{"moment_id":1,"accid":1,"comment_num":1,"_id":0}).One(&commentinfo)
        if err_comment_info != nil && err_comment_info != mgo.ErrNotFound {
                log.Error(err_comment_info)
                return proto.ReturnCodeServerError
        }

        data := bson.M{"$inc": bson.M{"comment_num":1}}

        _,err_upsert := sComments.Upsert(selector, data)
        if err_upsert != nil {
                log.Error(err_upsert)
                return proto.ReturnCodeServerError
        }


        var comment_mgo CommentMgo
        comment_mgo.MomentID    = commentinfo.MomentID
        comment_mgo.CommentID   = bson.ObjectIdHex(req.CommentID)
        comment_mgo.Time        = util.GetTimestamp()
        comment_mgo.Content     = req.Content
        comment_mgo.AccID       = my_accid
        comment_mgo.ID          = bson.NewObjectId()
        comment_mgo.CommentedAccID  = commentinfo.AccID


        sComment := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
        err_insert := sComment.Insert(&comment_mgo)
        if err_insert != nil {
                log.Error(err_insert)
                return proto.ReturnCodeServerError
        }

        return proto.ReturnCodeOK
}

func UploadCommentRsp(r *http.Request) (interface {}, int)  {
        check_token := CheckToken(r)
        if check_token  != proto.ReturnCodeOK {
                return nil, check_token
        }

        body, body_err := ioutil.ReadAll(r.Body)
        if body_err != nil {
                log.Debug("body err:%s",body_err)
                return nil, proto.ReturnCodeMissParm
        }

        var req proto.CommentReq
        json_err := json.Unmarshal(body, &req)
        if json_err != nil {
                log.Debug("json err:%s", json_err)
                return nil, proto.ReturnCodeMissParm
        }
        my_accid    := GetMyAccID(r)
        if req.MomentID != "" {
                retcode := UploadMomentComment(my_accid, req) 
                if retcode == proto.ReturnCodeOK {
                        rsp := GetMomentComment(my_accid, req.MomentID, "", 0)
                        if rsp != nil {
                                return rsp, proto.ReturnCodeOK
                        } else {
                                return nil, proto.ReturnCodeServerError
                        }
                } else {
                        return nil, retcode
                }
        } else if req.CommentID != "" {
                retcode := UploadCommentComment(my_accid, req)
                if retcode == proto.ReturnCodeOK {
                        rsp := GetCommentComment(my_accid, req.CommentID, "", 0)
                        if rsp != nil {
                                return rsp, proto.ReturnCodeOK
                        } else {
                                return nil, proto.ReturnCodeServerError
                        }

                } else {
                        return nil, retcode
                }
        } else {
                return nil, proto.ReturnCodeMissParm
        }
}

func GetCommentRsp(r *http.Request) (interface {}, int)  {
        vars := r.URL.Query();
        my_accid    := GetMyAccID(r)
        limit_num := GetIntUrlParmByName(r, "num")
        start_id  := GetObjectIDByName(r, "start_id")
 
        if len(vars["moment_id"]) > 0 {
                moment_id := vars["moment_id"][0]
                rsp := GetMomentComment(my_accid, moment_id, start_id, limit_num)
                if rsp != nil {
                        IncMomentReadNum(moment_id)
                        return rsp, proto.ReturnCodeOK
                } else {
                        return nil, proto.ReturnCodeServerError
                }

        } else if len(vars["comment_id"]) > 0 {
                comment_id := vars["comment_id"][0]
                rsp := GetCommentComment(my_accid, comment_id, start_id, limit_num)
                if rsp != nil {
                        return rsp, proto.ReturnCodeOK
                } else {
                        return nil, proto.ReturnCodeServerError
                }
        } else {
                return nil, proto.ReturnCodeMissParm
        }
}

func  DeleteCommentRsp(r *http.Request) (int) {
        if CheckUrlParm(r , "moment_id") != proto.ReturnCodeOK {
                return proto.ReturnCodeMissParm
        }

        return  proto.ReturnCodeOK
}
