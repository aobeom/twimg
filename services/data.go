package services

import (
	"strconv"
	"strings"
	"time"
)

// Param2str 参数转换为字符串
func Param2str(p interface{}) string {
	switch t := p.(type) {
	case bool:
		return strconv.FormatBool(t)
	case int:
		return strconv.Itoa(t)
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	}
	return ""
}

// DateFormat UTC时间格式化
func DateFormat(layout string, t string) string {
	rawTime, _ := time.Parse("Mon Jan 02 15:04:05 Z0700 2006", t)
	uTime := rawTime.Format(layout)
	return uTime
}

// SaveInfo 设置文件名
func SaveInfo(url string) (u string, fn string) {
	rawURLParts := strings.Split(url, "#")
	uDate, uURL := rawURLParts[0], rawURLParts[1]

	urlParts := strings.Split(uURL, "/")
	urlFilename := urlParts[len(urlParts)-1]
	urlFnParts := strings.Split(urlFilename, "?")
	uName := urlFnParts[0]

	u = uURL
	fn = uDate + "_" + uName
	return
}
