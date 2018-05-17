package proto
type UpdateUserInfoReq struct {
        Name        string      `json:"name"`
        Avatar      string      `json:"avatar"`
        Sex         int32       `json:"sex"`
        Birthday    int64       `json:"birthday"`
}
