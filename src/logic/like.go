package logic

import (
        "proto"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2/bson"
        mgohelper "mongo"
        "conf"
)

func UploadLikeRsp(r *http.Request) (int)  {
        vars := r.URL.Query();
        my_accid        := GetMyAccID(r)
        if len(vars["moment_id"]) > 0 {

                sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["moment_id"][0])}
                data     := bson.M{"$addToSet":bson.M{"likes":my_accid}}
                _, err := sMoments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家点赞动态,moment_id:%s,accid:%d", vars["moment_id"][0], my_accid)

        } else if len(vars["comment_id"]) > 0 {
                sComments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["comment_id"][0])}
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
        if len(vars["moment_id"]) > 0 {

                sMoments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("moments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["moment_id"][0])}
                data     := bson.M{"$pull":bson.M{"likes":my_accid}}
                _, err := sMoments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家取消点赞动态,moment_id:%s,accid:%d", vars["moment_id"][0], my_accid)

        } else if len(vars["comment_id"]) > 0 {
                sComments := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("comments")
                selector := bson.M{"_id":  bson.ObjectIdHex(vars["comment_id"][0])}
                data     := bson.M{"$pull":bson.M{"likes":my_accid}}
                _, err := sComments.Upsert(selector, data)
                if err != nil {
                        return proto.ReturnCodeServerError
                }

                log.Debug("玩家取消点赞评论,comment_id:%s,accid:%d", vars["comment_id"][0], my_accid)

        }

        return  proto.ReturnCodeOK
}
