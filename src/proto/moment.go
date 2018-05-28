package proto
import (
        "gopkg.in/mgo.v2/bson"
)

const (
            ValidOK             = 1
            ValidDeleteByAdmin  = 2
            ValidDeleteByMe     = 3
            ValidWaitForCheck   = 4
            ValidCheckNoPass    = 5
)

type PublishMomentReq struct {
        Content string   `json:"content"`
        Pic     []string `json:"pic,omitempty"`
        Video   string   `json:"video,omitempty"`
}

type MomentRet struct {
        ID          bson.ObjectId   `json:"id"`
        Content     string          `json:"content"`
        Time        int64           `json:"time"` 
        Pic         []string        `json:"pic,omitempty"`
        Video       string          `json:"video,omitempty"`
        ReadNum     uint32          `json:"read_num"`
        CommentNum   uint32         `json:"comment_num"`
        LikeNum     uint32          `json:"like_num"`
        Liked       bool            `json:"liked"`
        User        UserInfoRet     `json:"user"`
        ToTopTime   int64           `json:"to_top_time,omitempty"`
        Valid       int             `json:"valid",omitempty"`
}

