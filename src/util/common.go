package util

import (
        "reflect"
        "time"
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

