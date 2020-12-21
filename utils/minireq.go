package utils

import (
	"github.com/aobeom/minireq"
)

// UserAgent 全局 UA
var UserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36 Edg/86.0.622.56"

// Minireq 初始化
var Minireq *minireq.MiniRequest

// MiniHeaders Headers
type MiniHeaders = minireq.Headers

// MiniParams Params
type MiniParams = minireq.Params

// MiniJSONData JSONData
type MiniJSONData = minireq.JSONData

// MiniFormData FormData
type MiniFormData = minireq.FormData

// MiniAuth AuthData
type MiniAuth = minireq.Auth

func init() {
	Minireq = minireq.Requests()
}
