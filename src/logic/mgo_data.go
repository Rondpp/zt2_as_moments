package logic 

import (
        "gopkg.in/mgo.v2/bson"
)

type UserInfoMgo struct {
        AccID       int64               `bson:"accid"`
        Name        string              `bson:"name"`
        Follows     []int64             `bson:"follows"`
        Fans        []int64             `bson:"fans"`
        Moments     []bson.ObjectId     `bson:"moments"`
        Avatar      string              `bson:"avatar"`
        Sex         int32               `bson:"sex"`
        Birthday    int64               `bson:"birthday"`
        Type        uint32              `bons:"type"`
}

type MomentMgo struct {
        ID              bson.ObjectId   `bson:"_id"`
        AccID           int64           `bson:"accid"`
        Content         string          `bson:"content"`
        Time            int64           `bson:"time"`
        Pic             []string        `bson:"pic,omitempty"`
        Video           string          `bson:"video,omitempty"`
        ReadNum         uint32          `bson:"read_num"`
        CommentNum       uint32         `bson:"comment_num"`
        Like            []uint32        `bson:"likes"`
}

type FansInfoMgo struct {
        Fans    []int64     `bson:"fans"`
}

type FollowsInfoMgo struct {
        Follows []uint64    `bson:"follows"`
}

type CommentMgo struct {
        MomentID            bson.ObjectId        `bson:"moment_id"`
        CommentID           bson.ObjectId        `bson:"comment_id,omitempty"`
        ID                  bson.ObjectId        `bson:"_id"`
        Time                int64                `bson:"time"`
        Content             string               `bson:"content"`
        AccID               int64                `bson:"accid"`
        CommentedAccID      int64                `bson:"commented_accid"`
        Like                []int64              `bson:"likes,omitempty"`
        CommentNum          uint32               `bson:"comment_num,omitempty"`
}
/*
type CommentCommentMgo struct {
        CommentID      uint32    `bson:"comment_id"`
        Time           int64     `bson:"time"`
        Content        string    `bson:"content"`
        AccID          uint32    `bson:"accid"`
}

type CommentCommentListMgo struct {
        CommentComments []CommentCommentMgo `bson:"comments"`
}
*/
