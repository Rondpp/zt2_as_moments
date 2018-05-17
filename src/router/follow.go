package router

import (
        "net/http"
        "logic"
        "proto"
)

func GetFollowHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetFollowRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func UploadFollowHandler(w http.ResponseWriter, r *http.Request) {
         ret := logic.CheckToken(r) 
         if ret != proto.ReturnCodeOK {
                 logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
         } else {

                 retcode := logic.UploadFollowRsp(r)
                 logic.SendResponse(w, logic.GetErrResponseWithCode(retcode))
         }
}

func DeleteFollowHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r) 
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        } else {
                retcode := logic.DeleteFollowRsp(r)
                logic.SendResponse(w, logic.GetErrResponseWithCode(retcode))
        }
}

func FollowHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetFollowHandler(w, r)
        } else if  r.Method == "POST" {
                UploadFollowHandler(w, r)
        } else if  r.Method == "DELETE" {
                DeleteFollowHandler(w, r)
        }
}
