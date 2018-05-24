package proto

const (
        ReturnCodeServerError   = -1
        ReturnCodeOK            = 0
        ReturnCodeMissParm      = 2
        ReturnCodeMissHeader    = 3
        ReturnCodeTokenWrong    = 4
        ReturnCodeTokenExpired  = 5
        ReturnCodeParmWrong     = 6
        ReturnCodeNoPermission  = 7
        ReturnCodeForbidden     = 8
)

var (
        statusMessages = map[int]string{
                ReturnCodeServerError   :   "服务器错误",
                ReturnCodeOK            :   "OK",
                ReturnCodeMissParm      :   "缺少必填参数",
                ReturnCodeMissHeader    :   "缺少header",
                ReturnCodeTokenWrong    :   "Token不对",
                ReturnCodeTokenExpired  :   "Token过期",
                ReturnCodeParmWrong     :   "参数错误",
                ReturnCodeNoPermission  :   "没有权限",
                ReturnCodeForbidden     :   "禁言中",
        }
)

func StatusMessage(statusCode int) string {
        s := statusMessages[statusCode]
        if s == "" {
                s = "未知错误"
        }
        return s
}

