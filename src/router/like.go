package router

import (
        "net/http"
        "logic"
        "proto"
)

func UploadLikeHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r)
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        } else {
                retcode := logic.UploadLikeRsp(r)
                logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
        }
}

func DeleteLikeHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r)
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        } else {

                retcode := logic.DeleteLikeRsp(r)
                logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
        }
}

func LikeHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "POST" {
                UploadLikeHandler(w, r)
        } else if  r.Method == "DELETE" {
                DeleteLikeHandler(w, r)
        }
}
