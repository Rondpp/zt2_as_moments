package router

import (
        "fmt"
        "net/http"
        log "github.com/jeanphorn/log4go"
)

func router(path string, handler func(http.ResponseWriter, *http.Request)) {
        http.HandleFunc(path, commonHandler(handler))
}

func commonHandler(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
        return func(writer http.ResponseWriter, request *http.Request) {
                if request.Method == "OPTIONS" {
                        handleGlobalOptions(writer)
                } else {
                        handler(writer, request)
                }
        }
}

func handleGlobalOptions(w http.ResponseWriter) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "accid,token,Content-Type,*")
        w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,*")
        w.WriteHeader(http.StatusOK)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
        fmt.Println(w, "welcome to zt2")
}

func  InitRouter() {
        router("/", IndexHandler)
        router("/moments/", MomentsHandler) // 发布动态，获取某人动态，获取最新动态
        router("/user/", UserHandler) // 个人信息
        router("/fans/", FansHandler) // 某人的粉丝
        router("/follow/", FollowHandler) // 关注取关,某人关注的人
        router("/like/", LikeHandler) // 点赞
        router("/search/", SearchHandler) // 搜索
        router("/comment/", CommentHandler) // 评论
        router("/message/", MessageHandler) // 评论
        log.Debug("InitRouter成功")
}
