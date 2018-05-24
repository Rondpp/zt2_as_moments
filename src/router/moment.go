package router

import (
        "net/http"
        "logic"
        "proto"
)

func GetMomentsHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetMomentRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func UploadMomentsHandler(w http.ResponseWriter, r *http.Request) {
         ret := logic.CheckToken(r) 
         if ret != proto.ReturnCodeOK {
                 logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
         } else {
                 ret, err := logic.CheckForbidden(r)
                 if ret != proto.ReturnCodeOK {
                         logic.SendResponse(w, logic.GetErrResponseWithCodeMsg(ret,err.Error()))
                         return
                 }

                 rsp, retcode := logic.UploadMomentRsp(r)
                 logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
         }
}

func DeleteMomentsHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r)
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        } else {
                retcode := logic.DeleteMomentRsp(r)
                logic.SendResponse(w, logic.GetErrResponseWithCode(retcode))
        }
}

func MomentsHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetMomentsHandler(w, r)
        } else if  r.Method == "POST" {
                UploadMomentsHandler(w, r)
        } else if  r.Method == "DELETE" {
                DeleteMomentsHandler(w, r)
        }
}
