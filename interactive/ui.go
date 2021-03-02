package interactive

import (
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strconv"
	"time"
	"twimg/services"
	"twimg/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func numCheck(s string) int {
	if s != "" {
		rule := regexp.MustCompile(`^[\d]+$`)
		result := rule.Match([]byte(s))
		if result {
			i, _ := strconv.Atoi(s)
			if i < 0 || i > 200 {
				return -1
			}
			return i
		}
		return -1
	}
	return 0
}

// UIRun UI Mode
func UIRun() {
	myApp := app.New()
	myApp.SetIcon(theme.MyLogo())
	myWin := myApp.NewWindow("Twimg")
	myWin.Resize(fyne.NewSize(300, 300))
	myWin.SetFixedSize(true)
	myWin.CenterOnScreen()

	twitter := services.Twitter

	var tExclude bool
	uiUsername := widget.NewEntry()
	uiLimit := widget.NewEntry()
	uiStatusID := widget.NewEntry()
	proxyEntry := widget.NewEntry()
	uiHelp := widget.NewLabel("")
	uiStatus := widget.NewLabel("")

	uiUsername.SetPlaceHolder("@username")
	uiLimit.SetPlaceHolder("Limit")
	uiStatusID.SetPlaceHolder("Latest ID")
	proxyEntry.SetPlaceHolder("Socks5 Proxy")

	uiExclude := widget.NewCheck("Exclude", func(e bool) {
		tExclude = e
	})
	uiExclude.SetChecked(true)

	uiDownload := widget.NewButton("Download", func() {
		start := time.Now()

		tUsername := uiUsername.Text
		tStatusID := uiStatusID.Text
		if tUsername != "" {
			uiHelp.SetText("Setting Token...")
			twitter.SetToken()
			if twitter.Token != "" {
				uiHelp.SetText("Set Target...")
				tLimit := numCheck(uiLimit.Text)
				if tLimit != -1 {
					var proxy string
					if tLimit != 0 {
						twitter.SetLimit(tLimit)
					}
					if tStatusID != "" {
						twitter.SetLastID(tStatusID)
					}
					if proxyEntry.Text != "" {
						proxy = proxyEntry.Text
					}
					twitter.SetExclude(tExclude)
					twitter.SetUser(tUsername)
					uiHelp.SetText("Checking Media...")
					urls, total := twitter.MediaURLs()
					if len(urls) != 0 {
						urlGroups, groupNum := services.DataGroups(urls, 20)
						uiHelp.SetText(fmt.Sprintf("Media: %d | Groups: %d", total, groupNum))
						uiStatus.SetText("Downloading media in groups...")
						for index, urlGroup := range urlGroups {
							uiStatus.SetText(fmt.Sprintf(" - Group %d", index+1))
							urlG := urlGroup.([]interface{})
							twitter.MediaDownload(urlG, runtime.NumCPU(), proxy)
							time.Sleep(time.Duration(2) * time.Second)
						}
						uiHelp.SetText("Finished.")
					} else {
						uiHelp.SetText("No Media")
					}
				} else {
					errMsg := "Limit can only be a number (<200)"
					uiHelp.SetText(errMsg)
					err := errors.New(errMsg)
					dialog.ShowError(err, myWin)
				}
			} else {
				errMsg := "Token Error\nPlease Check your configs/apikeys.json"
				uiHelp.SetText(errMsg)
				err := errors.New(errMsg)
				dialog.ShowError(err, myWin)
			}
		} else {
			errMsg := "Username is empty"
			uiHelp.SetText(errMsg)
			err := errors.New(errMsg)
			dialog.ShowError(err, myWin)
		}
		tc := time.Since(start)
		uiStatus.SetText(fmt.Sprintf("Time: %v", tc))
	})

	content := container.NewVBox(
		uiUsername,
		uiLimit,
		uiStatusID,
		proxyEntry,
		uiExclude,
		uiDownload,
		uiHelp,
		uiStatus,
	)

	myWin.SetContent(content)
	myWin.ShowAndRun()
}
