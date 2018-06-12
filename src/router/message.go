package router

import (
        "net/http"
        "logic"
        "proto"
        "strings"
)

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetMessageRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func GetHasUnReadHandler(w http.ResponseWriter, r *http.Request) {
        rsp := logic.HasUnReadMesssage(r)
        logic.SendResponse(w, logic.GetResponseWithCode(proto.ReturnCodeOK, rsp))
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.DeleteMessageRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                if strings.Index(r.RequestURI, "/message/unread") == 0 {
                        GetHasUnReadHandler(w, r)
                } else {
                        GetMessageHandler(w, r)
                }
        } else if r.Method == "DELETE" {
                DeleteMessageHandler(w, r)
        }
}
