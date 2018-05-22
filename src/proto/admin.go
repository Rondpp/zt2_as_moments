package proto

type AdminUserPermissionSetReq struct {
        Account     string  `json:"account"`
        AccID       int64   `json:"accid"`
        Permission  int64   `json:"permission"`
}

type AdminForbiddenReq struct {
        AccID       int64   `json:"accid"`
        Time        int64   `json:"time"`
}
