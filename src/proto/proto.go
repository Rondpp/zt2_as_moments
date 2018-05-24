package proto 

type Response struct {
        Code int         `json:"code"`
        Msg  string      `json:"msg,omitempty"`
        Data interface{} `json:"data,omitempty"`
}

const (
        PermissionOfficial  = 1 << 0 // 官方
        PermissionDelete    = 1 << 1 // 删除
        PermissionForbidden = 1 << 2 // 禁言
        PermissionToTop     = 1 << 3 // 置顶
)

type UserInfoRet struct {
        AccID       int64       `json:"accid"`
        Account     string      `json:"account"`
        Name        string      `json:"name"`
        FollowNum   uint32      `json:"follow_num"`
        FanNum      uint32      `json:"fan_num"`
        MomentNum   uint32      `json:"moment_num"`
        Avatar      string      `json:"avatar"`
        Sex         int32       `json:"sex"`
        Birthday    int64       `json:"birthday"`
        Followed    bool        `json:"followd"`
        Type        uint32      `json:"type"`       // 1超级管理员
        Permission  int64       `json:"permission"` // 1官方发布 2删除动态评论 3禁言 4置顶
}

