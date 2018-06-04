package logic

import (
        "net/http"
        log "github.com/jeanphorn/log4go"
        "gopkg.in/mgo.v2/bson"
        "proto"
        "time"
        "strconv"
        redishelper "redis"
        mgohelper "mongo"
        "conf"
)

func CheckUrlParm(r *http.Request, parms ...string) int {
        vars := r.URL.Query();
        var missing_parm string
        for _, parm := range parms {

                if len(vars[parm]) != 0{
                        continue
                }

                missing_parm = parm
                log.Error("1缺少必填参数:%s", parm)
                break
        }
        if (missing_parm == "") {
                return proto.ReturnCodeOK
        }
        return proto.ReturnCodeMissParm
}

func CheckToken(r *http.Request) int {
        log.Debug("玩家header:accid:%s,token:%s", r.Header["Accid"], r.Header["Token"])
        return proto.ReturnCodeOK

        if  len(r.Header["Accid"]) == 0 {
                return proto.ReturnCodeMissHeader
        }


        kvs, err := redishelper.HGetAll("app_token:" + r.Header["Accid"][0])

        if err != proto.ReturnCodeOK {
                log.Error(err)
                return proto.ReturnCodeServerError
        }
        retcode := proto.ReturnCodeServerError
        now := time.Now().Unix()


        token_start_time,_ := strconv.ParseInt(kvs["time"], 10, 32)

        // token正确且在有效期内
        if   kvs["token"] != r.Header["Token"][0] {
                retcode = proto.ReturnCodeTokenWrong
        } else if   token_start_time + int64(conf.GetCfg().TokenLastTime) > now {
                retcode = proto.ReturnCodeTokenWrong
        } else {
                retcode = proto.ReturnCodeOK
        }
        log.Debug("玩家token验证:%d:accid:%s,redis_token:%s,my_token:%s,token_start_time:%s, now:%d", retcode, r.Header["Accid"][0], kvs["token"],  r.Header["Token"][0], kvs["time"], now)

        return retcode
}

func GetMyAccID(r *http.Request) int64 {
        if len(r.Header["Accid"]) == 0 {
                return 0
        }

        my_accid,err := strconv.ParseInt(r.Header["Accid"][0], 10, 64)
        if err != nil {
                log.Error(err)
                return 0
        }
        return my_accid
}

func GetInt64UrlParmByName(r *http.Request, name string) int64 {
        vars := r.URL.Query();
        value, err := strconv.ParseInt(vars.Get(name), 10, 64)
        if err != nil {
                return 0
        }
        return value
}

func GetObjectIDByName(r *http.Request, parm_name string) string {
        vars := r.URL.Query();
        if len(vars[parm_name]) > 0 && len(vars[parm_name][0]) == 24 {
                return vars[parm_name][0]
        }
        return ""
}

func GetIntUrlParmByName(r *http.Request, name string) int {
        return int(GetInt64UrlParmByName(r, name))
}


func UserInfoMgoToRet(my_accid int64, userinfo_mgo *UserInfoMgo, userinfo_ret *proto.UserInfoRet) {
        userinfo_ret.AccID       = userinfo_mgo.AccID
        userinfo_ret.Account     = userinfo_mgo.Account
        userinfo_ret.Name        = userinfo_mgo.Name
        userinfo_ret.FollowNum   = uint32(len(userinfo_mgo.Follows))
        userinfo_ret.FanNum      = uint32(len(userinfo_mgo.Fans))
        userinfo_ret.MomentNum   = uint32(len(userinfo_mgo.Moments))
        userinfo_ret.Avatar      = userinfo_mgo.Avatar
        userinfo_ret.Sex         = userinfo_mgo.Sex
        userinfo_ret.Birthday    = userinfo_mgo.Birthday
        userinfo_ret.Permission  = userinfo_mgo.Permission
        if userinfo_mgo.AccID == conf.GetCfg().AdminUser.AccID {
                userinfo_ret.Type = 1
        }
        for _, v := range userinfo_mgo.Fans {
                if v == my_accid {
                        userinfo_ret.Followed = true
                        break
                }
        }
}

func MomentMgoToRet(my_accid int64, moment_mgo *MomentMgo, moment_ret * proto.MomentRet) {
        if moment_mgo == nil {
                return
        }
        moment_ret.ID           = moment_mgo.ID
        moment_ret.Content      = moment_mgo.Content
        moment_ret.Time         = moment_mgo.Time
        moment_ret.Pic          = moment_mgo.Pic
        moment_ret.Video        = moment_mgo.Video
        moment_ret.ReadNum      = moment_mgo.ReadNum
        moment_ret.CommentNum   = moment_mgo.CommentNum
        moment_ret.LikeNum      = uint32(len(moment_mgo.Like))
        moment_ret.ToTopTime    = moment_mgo.ToTopTime
        moment_ret.Valid        = moment_mgo.Valid

        var userinfo_in_mgo UserInfoMgo
        sUsers := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        err := sUsers.Find(bson.M{"accid": moment_mgo.AccID}).One(&userinfo_in_mgo)
        if err != nil {
                log.Debug(err)
        }

        UserInfoMgoToRet(my_accid, &userinfo_in_mgo, &moment_ret.User)

        for _, like_accid := range moment_mgo.Like {
                if int64(like_accid) == my_accid {
                        moment_ret.Liked = true
                }
        }
}

func CommentMgoToRet(my_accid int64, mgo *CommentMgo, ret *proto.CommentRet) {
        ret.ID          = mgo.ID
        ret.Time        = mgo.Time
        ret.Content     = mgo.Content
        ret.CommentNum   = mgo.CommentNum
        ret.LikeNum     = uint32(len(mgo.Like))
        for _, v := range mgo.Like {
                if v == my_accid {
                        ret.Liked = true
                        break
                }
        }

        var userinfo_in_mgo UserInfoMgo
        sUsers := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        sUsers.Find(bson.M{"accid": mgo.AccID}).One(&userinfo_in_mgo)

        var userinforet  proto.UserInfoRet
        UserInfoMgoToRet(my_accid, &userinfo_in_mgo, &userinforet)
        ret.User = userinforet
}

func CommentCommentMgoToRet(my_accid int64, mgo *CommentMgo, ret *proto.CommentCommentRet) {
        ret.CommentID   = mgo.ID
        ret.Time        = mgo.Time
        ret.Content     = mgo.Content

        var userinfo_in_mgo UserInfoMgo

        sUsers := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        sUsers.Find(bson.M{"accid": mgo.AccID}).One(&userinfo_in_mgo)

        var userinforet  proto.UserInfoRet
        UserInfoMgoToRet(my_accid, &userinfo_in_mgo, &userinforet)
        ret.User = userinforet
}

func GetPermissionByAccID(accid int64) int64 {
        var userinfo_in_mgo UserInfoMgo
        sUsers := mgohelper.GetSession().DB(conf.GetCfg().MgoCfg.DB).C("users")
        err := sUsers.Find(bson.M{"accid": accid}).One(&userinfo_in_mgo)
        if err != nil {
                log.Debug(err)
        }
        return userinfo_in_mgo.Permission
}
