package router

import (
        "net/http"
        "logic"
        "proto"
)

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetUserInfoRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func UploadUserHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r) 
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        }  else {
                rsp, ret := logic.UpdateUserInfoRsp(r)
                logic.SendResponse(w, logic.GetResponseWithCode(ret, rsp))
        }

}


func UserHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetUserHandler(w, r)
        } else if  r.Method == "POST" {
                UploadUserHandler(w, r)
        }
}
