package services

import (
	"fmt"
	"net/http"
	"path/filepath"
	"sync"
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
	errImgs        []interface{}
	errLock        sync.Mutex
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
	utils.Minireq.Header.Set("User-Agent", utils.UserAgent)
	utils.Minireq.TimeOut(60)
}

func (twi *TwitterBasic) setDefault() {
	statusFlag = 0
	statusCount = 1
	twi.Lastid = ""
	twi.Limit = 200
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
func (twi *TwitterBasic) mediaFilter(tweets []interface{}) (imgURLs []interface{}, tweetStatusCount float64) {
	imgURLs = make([]interface{}, 0)
	for _, tweetR := range tweets {
		tweet := tweetR.(map[string]interface{})
		tweetUser := tweet["user"].(map[string]interface{})
		tweetStatusCount = tweetUser["statuses_count"].(float64)

		tweetStatusID := Param2str(tweet["id"].(float64))
		tweetCreateAt := DateFormat("20060102150405", tweet["created_at"].(string))
		if _, ok := tweet["extended_entities"]; ok {
			tweetEntities := tweet["extended_entities"].(map[string]interface{})
			tweetMedia := tweetEntities["media"].([]interface{})

			tTweetDetails := make(map[string]interface{})
			tTweetMediaURLs := make([]string, 0)
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
								tVURL = tBitrates["url"].(string)
							}
						}
					}
					if tVURL != "" {
						tTweetMediaURLs = append(tTweetMediaURLs, tVURL)
					}
				} else {
					imgURL := tMedium["media_url_https"].(string)
					if imgURL != "" {
						tURL := imgURL + "?format=jpg&name=orig"
						tTweetMediaURLs = append(tTweetMediaURLs, tURL)
					}
				}
			}
			tTweetDetails["date"] = tweetCreateAt
			tTweetDetails["urls"] = tTweetMediaURLs
			tTweetDetails["total"] = len(tTweetMediaURLs)
			imgURLs = append(imgURLs, tTweetDetails)
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
func (twi *TwitterBasic) dlcore(u interface{}) interface{} {
	data := u.(map[string]interface{})
	uDate := data["date"].(string)
	uURLs := data["urls"].([]string)

	for _, url := range uURLs {
		defer func() {
			if err := recover(); err != nil {
				errLock.Lock()
				errImgs = append(errImgs, u)
				errLock.Unlock()
			}
		}()
		savepath := SaveInfo(uDate, url, saveFolder)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer resp.Body.Close()
		Save2File(resp.Body, savepath)
		time.Sleep(time.Duration(2) * time.Second)
	}
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
func (twi *TwitterBasic) GetTweets() (tweets []interface{}) {
	if twi.Token != "" {
		timelineURL := "https://api.twitter.com/1.1/statuses/user_timeline.json"
		headers := utils.MiniHeaders{
			"Authorization": fmt.Sprintf("Bearer %s", twi.Token),
		}
		params := utils.MiniParams{
			"screen_name":     twi.User,
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
			fmt.Printf(" - Fetched %d responses\n", len(tweets))
		}
	}
	return
}

// MediaURLs 获取媒体地址
func (twi *TwitterBasic) MediaURLs() (media []interface{}, total int) {
	media = make([]interface{}, 0)
	if twi.User != "" {
		// 创建文件夹
		now := time.Now().Format("20060102150405")
		folderName := twi.User + "_" + now
		saveFolder = filepath.Join(utils.FileSuite.LocalPath(configs.Deployment()), folderName)
		utils.FileSuite.Create(saveFolder)
		// 获取地址
		tweets := twi.GetTweets()
		for len(tweets) > 0 {
			mURLs, tCounts := twi.mediaFilter(tweets)
			media = append(media, mURLs...)
			if twi.Limit < 200 {
				mediaUnique := RemoveDuplicate(media)
				for _, i := range mediaUnique {
					data := i.(map[string]interface{})
					total = total + data["total"].(int)
				}
				return
			}
			if tCounts < 3300 {
				statusCount = int(tCounts/200) + 2
			} else {
				statusCount = 20
			}
			tweets = twi.GetTweets()
		}
		mediaUnique := RemoveDuplicate(media)
		for _, i := range mediaUnique {
			data := i.(map[string]interface{})
			total = total + data["total"].(int)
		}
		return
	}
	return
}

// MediaDownload 下载
func (twi *TwitterBasic) MediaDownload(urls []interface{}, thread int) {
	if len(urls) != 0 {
		utils.TaskBoard(twi.dlcore, urls, thread)
		if len(errImgs) != 0 {
			fmt.Printf("-----\nPlease wait 10 seconds\nBefore retrying the failed task...\nTotal: (%d)\n-----\n", len(errImgs))
			time.Sleep(time.Duration(10) * time.Second)
			for _, errImg := range errImgs {
				twi.dlcore(errImg)
			}
		}
		twi.setDefault()
	}
}
