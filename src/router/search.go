package router

import (
        "net/http"
        "logic"
) 

func SearchHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetSearchRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}
