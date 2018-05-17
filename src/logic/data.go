package logic

import (
        "encoding/json"
        "net/http"
        "proto"
        "util"
)
func SendResponse(w http.ResponseWriter, data string) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Headers", "*")
        w.Header().Set("Access-Control-Allow-Methods", "*")
        w.Write([]byte(data))
}

func GetOKResponse(data interface{}) string {
        return GetResponseWithCode(proto.ReturnCodeOK, data)
}

func GetErrResponseWithCode(code int) string {
        return GetResponseWithCode(code, nil)
}

func GetResponseWithCode(code int, data interface{}) string {
        if util.IsEmpty(data) {
                data = nil
        }
        responseEntity := proto.Response{
                Code: code,
                Data: data,
        }
        if code != proto.ReturnCodeOK {
                responseEntity.Msg = proto.StatusMessage(code)
        }
        response, err := json.Marshal(responseEntity)
        if err != nil {
        } 
        return string(response)
}

