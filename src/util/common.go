package util

import (
        "reflect"
        "time"
        "fmt"
)

func IsEmpty(a interface{}) bool {
        defer func() {
                if err := recover(); err != nil {

                }
        }()

        if a == nil {
                return true
        }
        v := reflect.ValueOf(a)
        if v.Kind() == reflect.Ptr {
                v = v.Elem()
        } 
        return v.Interface() == reflect.Zero(v.Type()).Interface()
}

func GetTimestamp() int64 {
        return time.Now().UnixNano() / 1000000
}

func FormatTimeCH(timestamp_msec int64) string {
        timestamp := timestamp_msec * 1000 * 1000

        var str string
        day := int64(time.Duration(timestamp).Hours() / 24)
        if day > 0 {
                str += fmt.Sprintf("%d天",day)  
        }

        hour := int64(time.Duration(timestamp).Hours()) % 24
        if hour > 0 {
                str += fmt.Sprintf("%d小时",hour)  
        }

        min := int64(time.Duration(timestamp).Minutes()) % 60
        if min > 0 {
                str += fmt.Sprintf("%d分钟", min)  
        }
        return str
}
