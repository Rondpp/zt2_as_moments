package logic

import (
	"conf"
	"encoding/json"
	log "github.com/jeanphorn/log4go"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	mgohelper "mongo"
	"net/http"
	"proto"
	"util"
)

func GetMomentIDByCommentID(comment_id string) bson.ObjectId {
	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "comments")
	selector := bson.M{"_id": bson.ObjectIdHex(comment_id)}

	type MomentInfo struct {
		MomentID bson.ObjectId `bson:"moment_id"`
	}
	var momentinfo MomentInfo
	err := sMoments.Find(selector).Select(bson.M{"moment_id": 1, "_id": 0}).One(&momentinfo)
	if err != nil && err != mgo.ErrNotFound {
		log.Error(err)
	}
	return momentinfo.MomentID
}

func GetCommentOwnerByID(comment_id string) int64 {
	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "comments")
	selector := bson.M{"_id": bson.ObjectIdHex(comment_id)}

	type MomentInfo struct {
		AccID int64 `bson:"accid"`
	}
	var momentinfo MomentInfo
	err := sMoments.Find(selector).Select(bson.M{"accid": 1, "_id": 0}).One(&momentinfo)
	if err != nil && err != mgo.ErrNotFound {
		log.Error(err)
	}
	return momentinfo.AccID
}

func IncMomentReadNum(moment_id string) {
	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "moments")
	selector := bson.M{"_id": bson.ObjectIdHex(moment_id)}

	data := bson.M{"$inc": bson.M{"read_num": 1}}
	_, err := sMoments.Upsert(selector, data)
	if err != nil {
		log.Error(err)
	}
}

func GetCommentByID(comment_id string) *CommentMgo {
	session := mgohelper.GetSession()
	defer session.Close()

	var comment_mgo CommentMgo

	sComment := mgohelper.GetCollection(session, "comments")
	selector := bson.M{"_id": bson.ObjectIdHex(comment_id), "valid": proto.ValidOK}
	err := sComment.Find(selector).One(&comment_mgo)

	if err != nil {
		log.Error(err)
		return nil
	}
	return &comment_mgo
}

func GetMomentComment(my_accid int64, moment_id string, start_id string, limit_num int) *[]proto.CommentRet {
	session := mgohelper.GetSession()
	defer session.Close()

	if limit_num == 0 {
		limit_num = conf.GetCfg().MgoCfg.PageLimit
	}

	log.Debug("获取动态的评论,moment_id:%s,start_id:%s", moment_id, start_id)
	var comment_mgo_list []CommentMgo
	selector := bson.M{"moment_id": bson.ObjectIdHex(moment_id), "comment_id": bson.M{"$exists": false}, "valid": proto.ValidOK}

	if start_id != "" {
		selector = bson.M{"moment_id": bson.ObjectIdHex(moment_id), "comment_id": bson.M{"$exists": false}, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "valid": proto.ValidOK}
	}

	sComment := mgohelper.GetCollection(session, "comments")
	err := sComment.Find(selector).Sort("-time").Limit(limit_num).All(&comment_mgo_list)
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
	session := mgohelper.GetSession()
	defer session.Close()

	if limit_num == 0 {
		limit_num = conf.GetCfg().MgoCfg.PageLimit
	}

	var comment_mgo_list []CommentMgo

	sComment := mgohelper.GetCollection(session, "comments")
	selector := bson.M{"comment_id": bson.ObjectIdHex(comment_id), "valid": proto.ValidOK}
	if start_id != "" {
		selector = bson.M{"comment_id": bson.ObjectIdHex(comment_id), "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "valid": proto.ValidOK}
	}

	err := sComment.Find(selector).Sort("-time").Limit(limit_num).All(&comment_mgo_list)

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

	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "moments")
	selector := bson.M{"_id": bson.ObjectIdHex(req.MomentID)}

	type CommentInfo struct {
		AccID int64  `bson:"accid"`
		Num   uint32 `bson:"comment_num"`
	}
	var commentinfo CommentInfo
	err_comment_info := sMoments.Find(selector).Select(bson.M{"accid": 1, "comment_num": 1, "_id": 0}).One(&commentinfo)
	if err_comment_info != nil && err_comment_info != mgo.ErrNotFound {
		log.Error(err_comment_info)
		return proto.ReturnCodeServerError
	}

	data := bson.M{"$inc": bson.M{"comment_num": 1}}

	_, err_upsert := sMoments.Upsert(selector, data)
	if err_upsert != nil {
		log.Error(err_upsert)
		return proto.ReturnCodeServerError
	}

	var comment_mgo CommentMgo
	comment_mgo.MomentID = bson.ObjectIdHex(req.MomentID)
	comment_mgo.Time = util.GetTimestamp()
	comment_mgo.Content = req.Content
	comment_mgo.AccID = my_accid
	comment_mgo.ID = bson.NewObjectId()
	comment_mgo.CommentedAccID = commentinfo.AccID
	comment_mgo.Valid = proto.ValidOK
	comment_mgo.Type = proto.MessageTypeUser

	sComment := mgohelper.GetCollection(session, "comments")
	err_insert := sComment.Insert(&comment_mgo)
	if err_insert != nil {
		log.Error(err_insert)
		return proto.ReturnCodeServerError
	}
	NotifyNewMessageToMe(commentinfo.AccID)

	return proto.ReturnCodeOK
}

func UploadCommentComment(my_accid int64, req proto.CommentReq) int {
	log.Debug("玩家评论评论,accid:%d,comment_id:%s", my_accid, req.CommentID)

	session := mgohelper.GetSession()
	defer session.Close()

	sComments := mgohelper.GetCollection(session, "comments")
	selector := bson.M{"_id": bson.ObjectIdHex(req.CommentID)}

	type CommentInfo struct {
		MomentID bson.ObjectId `bson:"moment_id"`
		AccID    int64         `bson:"accid"`
		Num      uint32        `bson:"comment_num"`
	}
	var commentinfo CommentInfo
	err_comment_info := sComments.Find(selector).Select(bson.M{"moment_id": 1, "accid": 1, "comment_num": 1, "_id": 0}).One(&commentinfo)
	if err_comment_info != nil && err_comment_info != mgo.ErrNotFound {
		log.Error(err_comment_info)
		return proto.ReturnCodeServerError
	}

	data := bson.M{"$inc": bson.M{"comment_num": 1}}

	_, err_upsert := sComments.Upsert(selector, data)
	if err_upsert != nil {
		log.Error(err_upsert)
		return proto.ReturnCodeServerError
	}

	var comment_mgo CommentMgo
	comment_mgo.MomentID = commentinfo.MomentID
	comment_mgo.CommentID = bson.ObjectIdHex(req.CommentID)
	comment_mgo.Time = util.GetTimestamp()
	comment_mgo.Content = req.Content
	comment_mgo.AccID = my_accid
	comment_mgo.ID = bson.NewObjectId()
	comment_mgo.CommentedAccID = commentinfo.AccID
	comment_mgo.Valid = proto.ValidOK
	comment_mgo.Type = proto.MessageTypeUser

	sComment := mgohelper.GetCollection(session, "comments")
	err_insert := sComment.Insert(&comment_mgo)
	if err_insert != nil {
		log.Error(err_insert)
		return proto.ReturnCodeServerError
	}
	NotifyNewMessageToMe(commentinfo.AccID)

	return proto.ReturnCodeOK
}

func UploadCommentRsp(r *http.Request) (interface{}, int) {
	check_token := CheckToken(r)
	if check_token != proto.ReturnCodeOK {
		return nil, check_token
	}

	body, body_err := ioutil.ReadAll(r.Body)
	if body_err != nil {
		log.Debug("body err:%s", body_err)
		return nil, proto.ReturnCodeMissParm
	}

	var req proto.CommentReq
	json_err := json.Unmarshal(body, &req)
	if json_err != nil {
		log.Debug("json err:%s", json_err)
		return nil, proto.ReturnCodeMissParm
	}
	my_accid := GetMyAccID(r)
	if req.MomentID != "" {
		retcode := UploadMomentComment(my_accid, req)
		if retcode == proto.ReturnCodeOK {
			rsp := GetMomentComment(my_accid, req.MomentID, "", 0)
			return rsp, proto.ReturnCodeOK
		} else {
			return nil, retcode
		}
	} else if req.CommentID != "" {
		retcode := UploadCommentComment(my_accid, req)
		if retcode == proto.ReturnCodeOK {
			rsp := GetCommentComment(my_accid, req.CommentID, "", 0)
			return rsp, proto.ReturnCodeOK

		} else {
			return nil, retcode
		}
	} else {
		return nil, proto.ReturnCodeMissParm
	}
}

func GetCommentRsp(r *http.Request) (interface{}, int) {
	vars := r.URL.Query()
	my_accid := GetMyAccID(r)
	limit_num := GetIntUrlParmByName(r, "num")
	start_id := GetObjectIDByName(r, "start_id")

	if len(vars["moment_id"]) > 0 {
		moment_id := GetObjectIDByName(r, "moment_id")
		if moment_id == "" {
			return nil, proto.ReturnCodeParmWrong
		}

		rsp := GetMomentComment(my_accid, moment_id, start_id, limit_num)

		if rsp != nil {
			IncMomentReadNum(moment_id)
		}
		return rsp, proto.ReturnCodeOK

	} else if len(vars["comment_id"]) > 0 {
		comment_id := GetObjectIDByName(r, "comment_id")
		if comment_id == "" {
			return nil, proto.ReturnCodeParmWrong
		}
		_, ok := vars["start_id"]
		var rsp interface{}
		if ok {
			rsp = GetCommentComment(my_accid, comment_id, start_id, limit_num)
		} else {
			rsp = GetComment(my_accid, comment_id)
		}
		return rsp, proto.ReturnCodeOK

	} else {
		return nil, proto.ReturnCodeMissParm
	}
}

func GetComment(my_accid int64, comment_id string) interface{} {
	comment_mgo := GetCommentByID(comment_id)

	var comment_ret proto.CommentRet
	CommentMgoToRet(my_accid, comment_mgo, &comment_ret)
	return &comment_ret
}

func DeleteCommentRsp(r *http.Request) int {
	if CheckUrlParm(r, "moment_id") != proto.ReturnCodeOK {
		return proto.ReturnCodeMissParm
	}

	return proto.ReturnCodeOK
}
