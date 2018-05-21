package router

import (
        "net/http"
        "logic"
)

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetMessageRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func DeleteMessageHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.DeleteMessageRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetMessageHandler(w, r)
        } else if r.Method == "DELETE" {
                DeleteMessageHandler(w, r)
        }
}
