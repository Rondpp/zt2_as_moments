package router

import (
        "net/http"
        "logic"
        "proto"
)

func GetCommentHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetCommentRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func UploadCommentHandler(w http.ResponseWriter, r *http.Request) {
         ret := logic.CheckToken(r) 
         if ret != proto.ReturnCodeOK {
                 logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
         } else {

                 rsp, retcode := logic.UploadCommentRsp(r)
                 logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
         }
}

func DeleteCommentHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r) 
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
        } else {
                retcode := logic.DeleteCommentRsp(r)
                logic.SendResponse(w, logic.GetErrResponseWithCode(retcode))
        }
}

func CommentHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetCommentHandler(w, r)
        } else if  r.Method == "POST" {
                UploadCommentHandler(w, r)
        } else if  r.Method == "DELETE" {
                DeleteCommentHandler(w, r)
        }
}
