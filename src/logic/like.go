package logic

import (
        "proto"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
)

func UploadLikeRsp(r *http.Request) (int)  {
        vars := r.URL.Query();
        my_accid        := GetMyAccID(r)

        session := mgohelper.GetSession()
        defer session.Close()

        if len(vars["moment_id"]) > 0 {

                sMoments := mgohelper.GetCollection(session, "moments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["moment_id"][0])}
                var tmp interface {}
                exists := sMoments.Find(selector).One(tmp)
                if exists != nil {
                        log.Debug(exists)
                        return proto.ReturnCodeServerError
                }

                data     := bson.M{"$addToSet":bson.M{"likes":my_accid}}
                _, err := sMoments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家点赞动态,moment_id:%s,accid:%d", vars["moment_id"][0], my_accid)

        } else if len(vars["comment_id"]) > 0 {
                sComments := mgohelper.GetCollection(session, "comments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["comment_id"][0])}

                var tmp interface {}
                exists := sComments.Find(selector).One(tmp)
                if exists != nil {
                        log.Debug(exists)
                        return proto.ReturnCodeServerError
                }


                data     := bson.M{"$addToSet":bson.M{"likes":my_accid}}
                _, err := sComments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家点赞评论,comment_id:%s,accid:%d", vars["comment_id"][0], my_accid)

        }

        return  proto.ReturnCodeOK
}

func  DeleteLikeRsp(r *http.Request) (int) {
        vars := r.URL.Query();
        my_accid        := GetMyAccID(r)

        session := mgohelper.GetSession()
        defer session.Close()

        if len(vars["moment_id"]) > 0 {

                sMoments := mgohelper.GetCollection(session, "moments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["moment_id"][0])}

                var tmp interface {}
                exists := sMoments.Find(selector).One(tmp)
                if exists != nil {
                        log.Debug(exists)
                        return proto.ReturnCodeServerError
                }


                data     := bson.M{"$pull":bson.M{"likes":my_accid}}
                _, err := sMoments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家取消点赞动态,moment_id:%s,accid:%d", vars["moment_id"][0], my_accid)

        } else if len(vars["comment_id"]) > 0 {
                sComments := mgohelper.GetCollection(session, "comments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["comment_id"][0])}

                var tmp interface {}
                exists := sComments.Find(selector).One(tmp)
                if exists != nil {
                        log.Debug(exists)
                        return proto.ReturnCodeServerError
                }


                data     := bson.M{"$pull":bson.M{"likes":my_accid}}
                _, err := sComments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家取消点赞评论,comment_id:%s,accid:%d", vars["comment_id"][0], my_accid)

        }

        return  proto.ReturnCodeOK
}
