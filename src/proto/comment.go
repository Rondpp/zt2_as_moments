package proto
import (
        "gopkg.in/mgo.v2/bson"
)

type CommentReq struct {
        MomentID        string      `json:"moment_id"`
        CommentID       string      `json:"comment_id"`
        Content         string      `json:"content"`
}


type CommentRet struct {
        ID              bson.ObjectId   `json:"id"`
        Time            int64           `json:"time"`
        Content         string          `json:"content"`
        CommentNum      uint32          `json:"comment_num"`
        LikeNum         uint32          `json:"like_num"`
        Liked           bool            `json:"liked"`
        User            UserInfoRet     `json:"user"`
}

type CommentCommentRet struct {
        CommentID       bson.ObjectId   `json:"comment_comment_id"`
        Time            int64           `json:"time"`
        Content         string          `json:"content"`
        User            UserInfoRet     `json:"user"`
}

