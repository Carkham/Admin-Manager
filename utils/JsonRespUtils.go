package utils

import "admin/model"

func SetBadRequestResp(data interface{}, extra interface{}) (resp model.JSONResp) {
	resp.Code = -1
	resp.Message = "参数格式错误"
	resp.Data = data
	resp.Extra = extra
	return
}

func SetServerErrorResp(data interface{}, extra interface{}) (resp model.JSONResp) {
	resp.Code = -1
	resp.Message = "服务器内部错误"
	resp.Data = data
	resp.Extra = extra
	return
}

func SetUnAuthResp(data interface{}, extra interface{}) (resp model.JSONResp) {
	resp.Code = -1
	resp.Message = "无权限访问"
	resp.Data = data
	resp.Extra = extra
	return
}

func SetOKResp(data interface{}, extra interface{}) (resp model.JSONResp) {
	resp.Code = 0
	resp.Message = "success"
	resp.Data = data
	resp.Extra = extra
	return
}
