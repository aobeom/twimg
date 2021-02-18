package services

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
	cnLocal, _ := time.LoadLocation("PRC")
	// time.RubyDate
	rawTime, _ := time.Parse("Mon Jan 02 15:04:05 -0700 2006", t)
	uTime := rawTime.In(cnLocal).Format(layout)
	return uTime
}

// SaveInfo 设置文件名
func SaveInfo(date, sid, url, saveFolder string) (savepath string) {
	urlParts := strings.Split(url, "/")
	urlFilename := urlParts[len(urlParts)-1]
	urlFnParts := strings.Split(urlFilename, "?")
	uName := urlFnParts[0]
	fn := fmt.Sprintf("%s_%s_%s", date, sid, uName)
	savepath = filepath.Join(saveFolder, fn)
	return
}

// RemoveDuplicate 去重
func RemoveDuplicate(data []interface{}) []interface{} {
	result := make([]interface{}, 0, len(data))
	tmp := map[interface{}]struct{}{}
	for _, i := range data {
		statusID := i.(map[string]interface{})["id"]
		if _, ok := tmp[statusID]; !ok {
			tmp[statusID] = struct{}{}
			result = append(result, i)
		}
	}
	return result
}

// Save2File 保存图片
func Save2File(raw io.Reader, savepath string) {
	dst, _ := os.Create(savepath)
	io.Copy(dst, raw)
	defer dst.Close()
}

// DataGroups 数据分组
func DataGroups(data []interface{}, piece int) ([]interface{}, int) {
	newData := make([]interface{}, 0)
	dataCounts := len(data)

	groupCounts := dataCounts / piece
	groupExtra := dataCounts % piece
	groupNums := groupCounts

	startIndex := 0
	endIndex := 0
	for i := 0; i < groupCounts; i++ {
		startIndex = i * piece
		endIndex = startIndex + piece
		newData = append(newData, data[startIndex:endIndex])
	}
	if groupExtra != 0 {
		lastData := data[endIndex : endIndex+groupExtra]
		groupNums++
		newData = append(newData, lastData)
	}
	return newData, groupNums
}
