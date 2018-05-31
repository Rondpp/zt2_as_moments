package proto
import (
        "gopkg.in/mgo.v2/bson"
)
const (
        MessageTypeUser   = 1
        MessageTypeAdmin  = 2
)

type MessageCommentMeRet struct {
        ID                  bson.ObjectId       `json:"id"`
        MomentRet           MomentRet           `json:"moment"`
        CommentRet          CommentRet          `json:"comment,omitempty"`
        CommentCommentRet   CommentCommentRet   `json:"comment_comment,omitempty"`
        Valid               int                 `json:"valid,omitempty"`
        Type                int                 `json:"type"`
}
