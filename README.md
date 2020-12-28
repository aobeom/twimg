# twimg
Apply for your own [Twitter](https://developer.twitter.com/en/docs) API Key.  

Modify configs/apikeys.json
```go
{
    "consumer_key": "your-key",
    "consumer_secret": "you-secret",
    "Access_token": "",
    "Access_token_secret": ""
}
```

## Build cmd
1. Rename ui.go to ui.go1 and fyne.syso to fyne.syso1
2. Modify main.go
```
package main

import (
	"twimg/interactive"
)

func main() {
	interactive.CmdRun()
}
```
3. Build
```
go build
```
4. Rename ui.go1 to ui.go and fyne.syso1 to fyne.syso
## Build GUI
1. Modify main.go
```
package main

import (
	"twimg/interactive"
)

func main() {
	interactive.UIRun()
}
```
2. Build
```
fyne package -os windows -icon twimg.png
```