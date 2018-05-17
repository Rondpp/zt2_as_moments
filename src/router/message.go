package router

import (
        "net/http"
        "logic"
)

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetMessageRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func MessageHandler(w http.ResponseWriter, r *http.Request) {
        if  r.Method == "GET" {
                GetMessageHandler(w, r)
        }
}
