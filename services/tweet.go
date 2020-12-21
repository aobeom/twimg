package services

import (
	"fmt"
	"path/filepath"
	"time"
	"twimg/configs"
	"twimg/utils"
)

// 循环终止标记
var (
	statusFlag     int    = 0     // 推文获取终止标记
	statusCount    int    = 1     // 推文终止参照值
	excludeReplies bool   = false // 排除回复
	includeRTS     bool   = true  // 包含转发
	saveFolder     string = ""    // 保存文件夹
)

// TwitterBasic 基类
type TwitterBasic struct {
	User     string
	Token    string
	StatusID string
	Lastid   string
	Limit    int
}

// Twitter 初始化
var Twitter *TwitterBasic

func init() {
	Twitter = new(TwitterBasic)
	Twitter.Limit = 200
}

// tokenReq 获取Token
func (twi *TwitterBasic) tokenReq(key, secret string) string {
	tokenURL := "https://api.twitter.com/oauth2/token"

	formData := utils.MiniFormData{
		"grant_type": "client_credentials",
	}

	auth := utils.MiniAuth{
		key, secret,
	}

	res := utils.Minireq.Post(tokenURL, formData, auth)
	if res.RawRes.StatusCode == 200 {
		data := res.RawJSON().(map[string]interface{})
		token := data["access_token"].(string)
		return token
	}
	return ""
}

// mediaFilter 推文过滤
func (twi *TwitterBasic) mediaFilter(tweets []interface{}) (imgURLs []string) {
	imgURLs = make([]string, 0)
	for _, tweetR := range tweets {
		tweet := tweetR.(map[string]interface{})
		tweetUser := tweet["user"].(map[string]interface{})
		tweetStatusCount := tweetUser["statuses_count"].(float64)
		if tweetStatusCount < 3300 {
			statusCount = int(tweetStatusCount/200) + 2
		} else {
			statusCount = 20
		}
		tweetStatusID := Param2str(tweet["id"].(float64))
		tweetCreateAt := DateFormat("20060102150405", tweet["created_at"].(string))
		if _, ok := tweet["extended_entities"]; ok {
			tweetEntities := tweet["extended_entities"].(map[string]interface{})
			tweetMedia := tweetEntities["media"].([]interface{})

			tURLs := make([]string, 0)
			for _, tMedia := range tweetMedia {
				tMedium := tMedia.(map[string]interface{})
				if _, ok := tMedium["video_info"]; ok {
					tVideos := tMedium["video_info"].(map[string]interface{})
					tVariants := tVideos["variants"].([]interface{})

					zeroVal := float64(0)
					var tVURL string
					for _, tVar := range tVariants {
						tBitrates := tVar.(map[string]interface{})
						if _, ok := tBitrates["bitrate"]; ok {
							tBitrate := tBitrates["bitrate"].(float64)
							if tBitrate > zeroVal {
								zeroVal = tBitrate
								tVURL = tweetCreateAt + "_" + tweetStatusID + "#" + tBitrates["url"].(string)
							}
						}
					}
					if tVURL != "" {
						tURLs = append(tURLs, tVURL)
					}
				} else {
					imgURL := tMedium["media_url_https"].(string)
					if imgURL != "" {
						tURL := tweetCreateAt + "_" + tweetStatusID + "#" + imgURL + "?format=jpg&name=orig"
						tURLs = append(tURLs, tURL)
					}
				}
			}
			imgURLs = append(imgURLs, tURLs...)
			if tweetStatusID != twi.Lastid {
				twi.Lastid = tweetStatusID
			}
		} else {
			twi.Lastid = tweetStatusID
		}
	}
	return
}

// dlcore 下载函数
func (twi *TwitterBasic) dlcore(u string) interface{} {
	url, mediaName := SaveInfo(u)
	savepath := filepath.Join(saveFolder, mediaName)

	res := utils.Minireq.Get(url)
	utils.FileSuite.Write(savepath, res.RawData())
	return nil
}

// SetLimit 设置每次获取的条数
func (twi *TwitterBasic) SetLimit(i int) {
	twi.Limit = i
}

// SetLastID 设置起点（降序）
func (twi *TwitterBasic) SetLastID(s string) {
	twi.Lastid = s
}

// SetUser 指定用户
func (twi *TwitterBasic) SetUser(u string) {
	twi.User = u
}

// SetExclude 只获取推主原创推文的媒体内容
func (twi *TwitterBasic) SetExclude(b bool) {
	if b {
		excludeReplies = true
		includeRTS = false
	}
}

// SetToken 设置认证 获取令牌
func (twi *TwitterBasic) SetToken() {
	apikeys := configs.APIKeys()
	if len(apikeys) != 0 {
		consumerKey := apikeys["consumer_key"].(string)
		consumerSecret := apikeys["consumer_secret"].(string)
		if consumerKey != "" && consumerSecret != "" {
			token := twi.tokenReq(consumerKey, consumerSecret)
			twi.Token = token
		}
	}
}

// GetTweets 获取所有Tweets
func (twi *TwitterBasic) GetTweets(user string, excludeReplies, includeRTS bool) (tweets []interface{}) {
	if twi.Token == "" {
		// utils.TokenError()
	} else {
		timelineURL := "https://api.twitter.com/1.1/statuses/user_timeline.json"
		headers := utils.MiniHeaders{
			"Authorization": "Bearer " + twi.Token,
		}
		params := utils.MiniParams{
			"screen_name":     user,
			"count":           Param2str(twi.Limit),
			"exclude_replies": Param2str(excludeReplies),
			"include_rts":     Param2str(includeRTS),
			"tweet_mode":      "extended",
		}
		if twi.Lastid != "" {
			params["max_id"] = twi.Lastid
		}

		res := utils.Minireq.Get(timelineURL, headers, params)
		if res.RawRes.StatusCode == 200 {
			statusFlag++
			if statusFlag > statusCount {
				return nil
			}
			tweets = res.RawJSON().([]interface{})
			fmt.Printf("Fetched %d responses\n", len(tweets))
		}
	}
	return
}

// MediaURLs 获取媒体地址
func (twi *TwitterBasic) MediaURLs() (media []string) {
	media = make([]string, 0)
	if twi.User == "" {
		// utils.UserError()
	} else {
		tweets := twi.GetTweets(twi.User, excludeReplies, includeRTS)
		for len(tweets) > 0 {
			mURLs := twi.mediaFilter(tweets)
			media = append(media, mURLs...)
			if twi.Limit < 200 {
				return
			}
			tweets = twi.GetTweets(twi.User, excludeReplies, includeRTS)
		}
	}
	return
}

// MediaDownload 下载
func (twi *TwitterBasic) MediaDownload(urls []string, thread int) {
	utils.Minireq.Header.Set("User-Agent", utils.UserAgent)

	folderName := twi.User + "_" + time.Now().Format("20060102_150405")
	saveFolder = filepath.Join(utils.FileSuite.LocalPath(configs.Deployment()), folderName)
	utils.FileSuite.Create(saveFolder)
	if len(urls) != 0 {
		utils.MultiRun(twi.dlcore, urls, thread)
	}
}
