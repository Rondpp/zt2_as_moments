package router

import (
        "strings"
        "net/http"
        log "github.com/jeanphorn/log4go"
        "logic"
        "proto"
        "conf"
)

func GetAdminUserListHandler(w http.ResponseWriter, r *http.Request) {
        rsp, retcode := logic.GetAdminUserListRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, rsp))
}

func UploadAdminUserHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.UploadAdminUserRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func UploadAdminForbiddenHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.UploadAdminForbiddenRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func UploadAdminToTopHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.UploadAdminToTopRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func DeleteAdminToTopHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.DeleteAdminToTopRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func UploadAdminDeleteHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.UploadAdminDeleteRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func GetAdminMomentsHandler(w http.ResponseWriter, r *http.Request) {
        data, retcode := logic.GetAdminMomentsRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, data))
}

func UploadAdminCheckMomentsHandler(w http.ResponseWriter, r *http.Request) {
        retcode := logic.UploadAdminCheckMomentsRsp(r)
        logic.SendResponse(w, logic.GetResponseWithCode(retcode, nil))
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
        ret := logic.CheckToken(r)
        if ret != proto.ReturnCodeOK {
                logic.SendResponse(w, logic.GetErrResponseWithCode(ret))
                return
        }

        my_accid := logic.GetMyAccID(r)
        if strings.Index(r.RequestURI, "/admin/user/") == 0 {
                if my_accid != conf.GetCfg().AdminUser.AccID {
                        logic.SendResponse(w, logic.GetErrResponseWithCode(proto.ReturnCodeNoPermission))
                        return
                }

                if  r.Method == "GET" {
                        GetAdminUserListHandler(w, r)
                } else if  r.Method == "POST" {
                        UploadAdminUserHandler(w, r)
                }
        } else if strings.Index(r.RequestURI, "/admin/forbidden/") == 0 {
                permission := logic.GetPermissionByAccID(my_accid)
                if (permission & proto.PermissionForbidden) != proto.PermissionForbidden {
                        logic.SendResponse(w, logic.GetErrResponseWithCode(proto.ReturnCodeNoPermission))
                        return
                }
                UploadAdminForbiddenHandler(w, r)
        } else if strings.Index(r.RequestURI, "/admin/totop/") == 0 {
                permission := logic.GetPermissionByAccID(my_accid)
                if (permission & proto.PermissionToTop) != proto.PermissionToTop {
                        logic.SendResponse(w, logic.GetErrResponseWithCode(proto.ReturnCodeNoPermission))
                        return
                }
                if  r.Method == "POST" {
                        UploadAdminToTopHandler(w, r)
                } else if r.Method == "DELETE" {
                        DeleteAdminToTopHandler(w, r)
                }
        } else if strings.Index(r.RequestURI, "/admin/delete/") == 0 {
                permission := logic.GetPermissionByAccID(my_accid)
                if (permission & proto.PermissionDelete) != proto.PermissionDelete {
                        logic.SendResponse(w, logic.GetErrResponseWithCode(proto.ReturnCodeNoPermission))
                        return
                }
                UploadAdminDeleteHandler(w, r)
        } else if strings.Index(r.RequestURI, "/admin/moments/") == 0 {
                permission := logic.GetPermissionByAccID(my_accid)
                log.Debug("permission:%d,%d",permission,proto.PermissionCheck)
                if (permission & proto.PermissionCheck) != proto.PermissionCheck {
                        logic.SendResponse(w, logic.GetErrResponseWithCode(proto.ReturnCodeNoPermission))
                        return
                }
                if  r.Method == "GET" {
                        GetAdminMomentsHandler(w, r)
                } else if r.Method == "POST" {
                        UploadAdminCheckMomentsHandler(w, r)
                }
        }
}
