package proto

const (
        ReturnCodeServerError   = -1
        ReturnCodeOK            = 0
        ReturnCodeMissParm      = 1
        ReturnCodeMissHeader    = 2
        ReturnCodeTokenWrong    = 3
        ReturnCodeTokenExpired  = 4
        ReturnCodeParmWrong     = 5
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
        }
)

func StatusMessage(statusCode int) string {
        s := statusMessages[statusCode]
        if s == "" {
                s = "未知错误"
        }
        return s
}

