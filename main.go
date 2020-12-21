package main

import (
	"fmt"
	"twimg/services"
)

var (
	userName     string
	tweetLimit   int
	tweetID      string
	tweetExclude string
	excludeRTS   bool
)

func help() {
	fmt.Println("-----")
	fmt.Println("You need to provide the following information for download.")
	fmt.Println("Example:")
	fmt.Println("  https://twitter.com/Twitter/status/1334542969530183683")
	fmt.Println()
	fmt.Println("@username: Twitter")
	fmt.Println("ExcludeRTS: no")
	fmt.Println("Starting ID: 1334542969530183683")
	fmt.Println("Limit: 20")
	fmt.Println()
	fmt.Println("[ExcludeRTS] means to exclude forwarding and reply, Default yes.")
	fmt.Println("[Starting ID] means to start fetched from this one.")
	fmt.Println("[Limit] is the number of items fetched each time.")
	fmt.Println("Limit and Starting ID are usually used together and can be empty.")
	fmt.Println("-----")
}

func userInterface() {
	fmt.Printf("@username: ")
	fmt.Scanln(&userName)

	fmt.Printf("ExcludeRTS (yes/no): ")
	fmt.Scanln(&tweetExclude)

	fmt.Printf("Limit: ")
	fmt.Scanln(&tweetLimit)

	fmt.Printf("Starting ID: ")
	fmt.Scanln(&tweetID)

	fmt.Println()
}

func main() {
	help()
	for {
		userInterface()
		if tweetExclude == "no" || tweetExclude == "n" {
			excludeRTS = false
		} else {
			excludeRTS = true
		}

		if userName == "" {
			fmt.Println("Username is empty.")
			continue
		}
		
		if userName != "" {
			break
		}
	}
	twitter := services.Twitter
	fmt.Println("1.Setting Token...")
	twitter.SetToken()
	if twitter.Token != "" {
		fmt.Println("2.Set Data...")
		twitter.SetUser(userName)
		if tweetLimit != 0 {
			twitter.SetLimit(tweetLimit)
		}
		if tweetID != "" {
			twitter.SetLastID(tweetID)
		}
		twitter.SetExclude(excludeRTS)
		fmt.Println("3.Checking Media...")
		urls := twitter.MediaURLs()
		if len(urls) != 0 {
			fmt.Printf("4.Total: %d\n", len(urls))
			fmt.Println("5.Downloading Media...")
			twitter.MediaDownload(urls, 16)
			fmt.Println("6.Finished.")
		} else {
			fmt.Println("4.No Media")
		}
	} else {
		fmt.Println("Token Error, Please Check your configs/apikeys.json")
	}
}
