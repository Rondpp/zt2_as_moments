package proto 

type Response struct {
        Code int         `json:"code"`
        Msg  string      `json:"msg,omitempty"`
        Data interface{} `json:"data,omitempty"`
}


type UserInfoRet struct {
        AccID       int64       `json:"accid"`
        Name        string      `json:"name"`
        FollowNum   uint32      `json:"follow_num"`
        FanNum      uint32      `json:"fan_num"`
        MomentNum   uint32      `json:"moment_num"`
        Avatar      string      `json:"avatar"`
        Sex         int32       `json:"sex"`
        Birthday    int64       `json:"birthday"`
        Followed    bool        `json:"followd"`
        Type        uint32      `json:"type"`
}

