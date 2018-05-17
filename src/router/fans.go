package router

import (
        "net/http"
        "logic"
) 

func FansHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetUserFansRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}
