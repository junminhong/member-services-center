package handler

import "time"

const (
	// OK 請求成功，適用Read、Update
	// Created 請求成功，適用Create
	// Accepted 此請求已被接受但未做任何處理
	// NoContent server已經處理請求，但未返回任何內容，適用Delete
	// BadRequest server無法理解請求，適用錯誤api參數
	// Unauthorized 未經過token認證
	// Forbidden 無權限訪問，這邊跟Unauthorized不同的是，這邊有token但無權限
	// NotFound	找不到資源
	OK                  = 200
	Created             = 201
	Accepted            = 202
	NoContent           = 204
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	RequestFormatError1 = 1000
	AuthError1          = 1200
	AuthError2          = 1201
	AuthError3          = 1202
	AuthOK1             = 1300
	RegisterError1      = 1400
	RegisterError2      = 1401
	RegisterOK1         = 1500
)

var ResponseFlag = map[int]string{
	OK:                  "請求成功",
	Created:             "請求成功",
	Accepted:            "請求成功",
	NoContent:           "請求成功",
	BadRequest:          "請依照API文件重新發起請求",
	Unauthorized:        "該請求未經過認證",
	Forbidden:           "你的權限不足以發起該請求",
	NotFound:            "",
	RequestFormatError1: "請依照API文件重新發起請求",
	AuthError1:          "token驗證失敗",
	AuthError2:          "email token已過期，請重新要求寄發新的驗證信件",
	AuthError3:          "信箱驗證失敗，請洽網站管理員",
	AuthOK1:             "信箱驗證成功",
	RegisterError1:      "帳號註冊失敗",
	RegisterError2:      "已存在該信箱",
	RegisterOK1:         "帳號註冊成功",
}

type Response struct {
	ResultCode int         `json:"result_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	TimeStamp  time.Time   `json:"time_stamp"`
}
