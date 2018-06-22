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
	redishelper "redis"
	"strconv"
	"util"
)

func GetAdminUserListRsp(r *http.Request) (*[]proto.UserInfoRet, int) {

	session := mgohelper.GetSession()
	defer session.Close()

	coll := mgohelper.GetCollection(session, "users")
	selector := bson.M{"permission": bson.M{"$gt": 0}}

	var userinfo_mgo_list []UserInfoMgo
	err := coll.Find(selector).All(&userinfo_mgo_list)
	if err != nil && err != mgo.ErrNotFound {
		log.Error(err)
		return nil, proto.ReturnCodeServerError
	}

	var rsp []proto.UserInfoRet
	for _, v := range userinfo_mgo_list {
		var userinfo_ret proto.UserInfoRet
		UserInfoMgoToRet(0, &v, &userinfo_ret)
		rsp = append(rsp, userinfo_ret)
	}
	log.Debug("管理员:%d", len(userinfo_mgo_list))
	return &rsp, proto.ReturnCodeOK
}

func UploadAdminUserRsp(r *http.Request) int {
	log.Debug("设置玩家权限")
	body, body_err := ioutil.ReadAll(r.Body)
	if body_err != nil {
		log.Debug("body err:%s", body_err)
		return proto.ReturnCodeMissParm
	}

	var req proto.AdminUserPermissionSetReq
	json_err := json.Unmarshal(body, &req)
	if json_err != nil {
		log.Debug("json err:%s", json_err)
		return proto.ReturnCodeMissParm
	}
	if req.AccID == conf.GetCfg().AdminUser.AccID || (req.AccID == 0 && len(req.Account) == 0) {
		return proto.ReturnCodeParmWrong
	}
	session := mgohelper.GetSession()
	defer session.Close()

	accid := req.AccID
	if accid == 0 {
		value, err := redishelper.HGet("account:"+req.Account, "accid")
		if err != proto.ReturnCodeOK {
			return proto.ReturnCodeServerError
		}
		accid, _ = strconv.ParseInt(value, 10, 64)
	}
	if accid == 0 {
		log.Error("accid = 0")
		return proto.ReturnCodeParmWrong
	}

	selector := bson.M{"accid": accid}
	data := bson.M{"$set": bson.M{"permission": req.Permission}}
	sUsers := mgohelper.GetCollection(session, "users")
	_, err := sUsers.Upsert(selector, data)

	if err != nil {
		log.Error(err)
		return proto.ReturnCodeServerError
	}

	log.Debug("设置玩家权限成功,account:%s,accid:%d,permission:%d", req.Account, accid, req.Permission)
	return proto.ReturnCodeOK
}

func UploadAdminForbiddenRsp(r *http.Request) int {

	body, body_err := ioutil.ReadAll(r.Body)
	if body_err != nil {
		log.Debug("body err:%s", body_err)
		return proto.ReturnCodeMissParm
	}

	var req proto.AdminForbiddenReq
	json_err := json.Unmarshal(body, &req)
	if json_err != nil {
		log.Debug("json err:%s", json_err)
		return proto.ReturnCodeMissParm
	}
	session := mgohelper.GetSession()
	defer session.Close()

	selector := bson.M{"accid": req.AccID}
	data := bson.M{"$set": bson.M{"forbidden_last_time": req.Time, "forbidden_start_time": util.GetTimestamp()}}
	sUsers := mgohelper.GetCollection(session, "users")
	_, err := sUsers.Upsert(selector, data)

	if err != nil {
		log.Error(err)
		return proto.ReturnCodeServerError
	}

	log.Debug("禁言,accid:%d,time:%d", req.AccID, req.Time)
	return proto.ReturnCodeOK
}

func GetOldestToTopMomentID() *bson.ObjectId {
	// 如果有两条置顶的,顶替掉最老的

	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "moments")

	type MomentInfo struct {
		ID bson.ObjectId `bson:"_id"`
	}
	var moment_info_list []MomentInfo

	selector := bson.M{"to_top_time": bson.M{"$exists": true}}
	err := sMoments.Find(selector).Sort("to_top_time").Select(bson.M{"_id": 1}).All(&moment_info_list)

	if err != nil && err != mgo.ErrNotFound {
		log.Error(err)
		return nil
	}
	if len(moment_info_list) >= 2 {
		return &moment_info_list[0].ID
	}
	return nil
}

func UploadAdminToTopRsp(r *http.Request) int {
	moment_id := GetObjectIDByName(r, "moment_id")
	if moment_id == "" {
		return proto.ReturnCodeMissParm
	}

	if !IsMomentExists(bson.ObjectIdHex(moment_id)) {
		return proto.ReturnCodeParmWrong
	}

	session := mgohelper.GetSession()
	defer session.Close()

	// 如果有两条置顶的,顶替掉最老的
	old_id := GetOldestToTopMomentID()
	if old_id != nil {
		log.Debug(old_id)
		selector := bson.M{"_id": old_id}
		data := bson.M{"$unset": bson.M{"to_top_time": 1}}
		sMoments := mgohelper.GetCollection(session, "moments")
		_, err := sMoments.Upsert(selector, data)
		if err != nil {
			log.Error(err)
			return proto.ReturnCodeServerError
		}
	}

	selector := bson.M{"_id": bson.ObjectIdHex(moment_id)}
	data := bson.M{"$set": bson.M{"to_top_time": util.GetTimestamp()}}
	sMoments := mgohelper.GetCollection(session, "moments")
	_, err := sMoments.Upsert(selector, data)

	if err != nil {
		log.Error(err)
		return proto.ReturnCodeServerError
	}

	log.Debug("置顶,moment_id:%d", moment_id)
	return proto.ReturnCodeOK
}

func DeleteAdminToTopRsp(r *http.Request) int {
	moment_id := GetObjectIDByName(r, "moment_id")
	if moment_id == "" {
		return proto.ReturnCodeMissParm
	}
	session := mgohelper.GetSession()
	defer session.Close()

	selector := bson.M{"_id": bson.ObjectIdHex(moment_id)}
	data := bson.M{"$unset": bson.M{"to_top_time": 1}}
	sMoments := mgohelper.GetCollection(session, "moments")
	_, err := sMoments.Upsert(selector, data)

	if err != nil {
		log.Error(err)
		return proto.ReturnCodeServerError
	}

	log.Debug("取消置顶,moment_id:%d", moment_id)
	return proto.ReturnCodeOK
}

func UploadAdminDeleteRsp(r *http.Request) int {
	moment_id := GetObjectIDByName(r, "moment_id")
	comment_id := GetObjectIDByName(r, "comment_id")
	my_accid := GetMyAccID(r)
	data := bson.M{"$set": bson.M{"valid": proto.ValidDeleteByAdmin}}

	session := mgohelper.GetSession()
	defer session.Close()

	var comment_mgo CommentMgo

	var coll *mgo.Collection
	var selector interface{}
	if moment_id != "" {
		selector = bson.M{"_id": bson.ObjectIdHex(moment_id)}
		coll = mgohelper.GetCollection(session, "moments")
		log.Debug("管理员删除动态,moment_id:%d", moment_id)
		commented_accid := GetMomentOwnerByID(moment_id)

		comment_mgo.MomentID = bson.ObjectIdHex(moment_id)
		comment_mgo.CommentedAccID = commented_accid
		if !SetVideoFlag(moment_id, VideoFlagDelByAdmin) {
			return proto.ReturnCodeServerError
		}
		NotifyNewMessageToMe(commented_accid)

	} else if comment_id != "" {
		selector = bson.M{"_id": bson.ObjectIdHex(comment_id)}
		coll = mgohelper.GetCollection(session, "comments")
		log.Debug("管理员删除评论,comment_id:%d", comment_id)
		commented_accid := GetCommentOwnerByID(comment_id)

		comment_mgo.MomentID = GetMomentIDByCommentID(comment_id)
		comment_mgo.CommentID = bson.ObjectIdHex(comment_id)
		comment_mgo.CommentedAccID = commented_accid
		NotifyNewMessageToMe(commented_accid)

	} else {
		return proto.ReturnCodeMissParm
	}
	_, err := coll.Upsert(selector, data)

	if err != nil {
		log.Error(err)
		return proto.ReturnCodeServerError
	}

	comment_mgo.Time = util.GetTimestamp()
	comment_mgo.Content = "管理员删除"
	comment_mgo.AccID = my_accid
	comment_mgo.ID = bson.NewObjectId()
	comment_mgo.Valid = proto.ValidOK
	comment_mgo.Type = proto.MessageTypeAdmin

	sComment := mgohelper.GetCollection(session, "comments")
	err_insert := sComment.Insert(&comment_mgo)
	if err_insert != nil {
		log.Error(err_insert)
		return proto.ReturnCodeServerError
	}

	return proto.ReturnCodeOK
}

func GetAdminMomentsRsp(r *http.Request) (interface{}, int) {
	my_accid := GetMyAccID(r)
	limit_num := GetIntUrlParmByName(r, "num")
	start_id := GetObjectIDByName(r, "start_id")

	if limit_num == 0 {
		limit_num = conf.GetCfg().MgoCfg.PageLimit
	}

	var moment_mgo_list []MomentMgo

	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "moments")
	selector := bson.M{"video": bson.M{"$exists": true}, "valid": proto.ValidWaitForCheck}

	// 置顶只有2个,翻页没有置顶的
	var err error
	if start_id != "" {
		selector = bson.M{"video": bson.M{"$exists": true}, "to_top_time": bson.M{"$exists": false}, "_id": bson.M{"$lt": bson.ObjectIdHex(start_id)}, "valid": proto.ValidWaitForCheck}
		err = sMoments.Find(selector).Sort("-time").Limit(limit_num).All(&moment_mgo_list)
	} else {
		err = sMoments.Find(selector).Sort("-to_top_time", "-time").Limit(limit_num).All(&moment_mgo_list)
	}

	if err != nil && err != mgo.ErrNotFound {
		log.Error(err)
		return nil, proto.ReturnCodeServerError
	}

	var rsp []proto.MomentRet
	for _, v := range moment_mgo_list {
		var moment_ret proto.MomentRet
		MomentMgoToRet(my_accid, &v, &moment_ret)
		rsp = append(rsp, moment_ret)
	}

	return &rsp, proto.ReturnCodeOK
}

func UploadAdminCheckMomentsRsp(r *http.Request) int {
	if CheckUrlParm(r, "moment_id") != proto.ReturnCodeOK {
		return proto.ReturnCodeMissParm
	}
	my_accid := GetMyAccID(r)
	moment_id := GetObjectIDByName(r, "moment_id")
	pass := GetIntUrlParmByName(r, "pass")

	log.Debug("管理员:%d审核视频动态,moment_id:%d,pass:%d", my_accid, moment_id, pass)

	session := mgohelper.GetSession()
	defer session.Close()

	sMoments := mgohelper.GetCollection(session, "moments")
	moment_selector := bson.M{"_id": bson.ObjectIdHex(moment_id)}
	var moment_data interface{}
	if pass == 0 {
		moment_data = bson.M{"$set": bson.M{"valid": proto.ValidDeleteByAdmin}}
		if !SetVideoFlag(moment_id, VideoFlagCheckNoPass) {
			return proto.ReturnCodeServerError
		}
	} else {
		if !SetVideoFlag(moment_id, VideoFlagCheckPass) {
			return proto.ReturnCodeServerError
		}
		moment_data = bson.M{"$set": bson.M{"valid": proto.ValidOK}}
	}
	_, moment_err := sMoments.Upsert(moment_selector, moment_data)
	if moment_err != nil {
		log.Error(moment_err)
		return proto.ReturnCodeServerError
	}
	user_accid := GetMomentOwnerByID(moment_id)

	selector := bson.M{"accid": user_accid}
	data := bson.M{"$pull": bson.M{"moments": bson.ObjectIdHex(moment_id)}}

	sUsers := mgohelper.GetCollection(session, "users")

	_, upsert_err := sUsers.Upsert(selector, data)
	if upsert_err != nil {
		log.Error(upsert_err)
		return proto.ReturnCodeServerError
	}
	return proto.ReturnCodeOK
}
