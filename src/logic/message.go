package logic

import (
        "proto"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "conf"
)

func GetCommentMeRsp(my_accid int64, start_id string, limit_num int) *[]proto.MessageCommentMeRet {
        if limit_num == 0 {
                limit_num = conf.GetCfg().MgoCfg.PageLimit
        }

        var comment_mgo_list []CommentMgo

        session := mgohelper.GetSession()
        defer session.Close()

        sComment := mgohelper.GetCollection(session, "comments")
        selector := bson.M{"commented_accid":my_accid, "read":bson.M{"$ne": 1}}
        if start_id != "" {
                selector = bson.M{"commented_accid" : my_accid, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}}
        }

        err := sComment.Find(selector).Sort("-time").Limit(limit_num).All(&comment_mgo_list)
        if err != nil && err != mgo.ErrNotFound {
                log.Error(err)
                return nil
        }

        var rsp []proto.MessageCommentMeRet
        for _, v := range comment_mgo_list {
                var message_comment_me_ret proto.MessageCommentMeRet
                message_comment_me_ret.ID = v.ID
                moment_mgo := GetMomentByID(v.MomentID.Hex(), true)
                if moment_mgo == nil {
                        continue;
                }

                var comment_ret proto.CommentRet
                if v.CommentID == "" {
                        CommentMgoToRet(my_accid, &v, &comment_ret)
                } else {
                        comment_mgo := GetCommentByID(v.CommentID.Hex())
                        if comment_mgo == nil {
                                continue;
                        }

                        CommentMgoToRet(my_accid, comment_mgo, &comment_ret)


                        var comment_comment_ret proto.CommentCommentRet
                        CommentCommentMgoToRet(my_accid, &v, &comment_comment_ret)
                        message_comment_me_ret.CommentCommentRet = comment_comment_ret
                }
                message_comment_me_ret.CommentRet = comment_ret

                var moment_ret proto.MomentRet
                MomentMgoToRet(my_accid, moment_mgo, &moment_ret)
                message_comment_me_ret.MomentRet = moment_ret
                message_comment_me_ret.Valid = v.Valid
                if v.Type != 0 {
                        message_comment_me_ret.Type = v.Type
                } else {
                        message_comment_me_ret.Type = proto.MessageTypeUser
                }

                rsp = append(rsp, message_comment_me_ret)
        }

        return &rsp
}

func GetMessageRsp(r *http.Request) (interface {}, int) {
        log.Debug("查询信息")
        if CheckUrlParm(r, "type") != proto.ReturnCodeOK {
                return nil, proto.ReturnCodeMissParm
        }
        vars     := r.URL.Query();
        my_accid := GetMyAccID(r)
        limit_num := GetIntUrlParmByName(r, "num")

        start_id := GetObjectIDByName(r, "start_id")

        if vars["type"][0] == "1" {
                rsp:= GetCommentMeRsp(my_accid, start_id, limit_num)
                return rsp, proto.ReturnCodeOK
        }
        return nil, proto.ReturnCodeMissParm
}

func DeleteMessageRsp(r *http.Request) int {
        id  := GetObjectIDByName(r, "id")

        session := mgohelper.GetSession()
        defer session.Close()

        sComment := mgohelper.GetCollection(session, "comments")
        selector := bson.M{"_id":bson.ObjectIdHex(id)}
        data     := bson.M{"$set":bson.M{"read": 1}}

        _, err := sComment.Upsert(selector, data)
        if err != nil {
                log.Error(err)
                return proto.ReturnCodeServerError
        }
        return proto.ReturnCodeOK
}

func HasUnReadMesssage(r *http.Request) (interface {}) {
        var rsp proto.MessageHasUnRead
        my_accid := GetMyAccID(r)

        session := mgohelper.GetSession()
        defer session.Close()

        sComment := mgohelper.GetCollection(session, "comments")
        selector := bson.M{"commented_accid":my_accid, "read":bson.M{"$ne": 1}}
        var data interface{}
        err      := sComment.Find(selector).One(data)
        if err != nil {
                rsp.UnRead = false
        }
        rsp.UnRead = true
        return &rsp
}
