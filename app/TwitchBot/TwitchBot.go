package TwitchBot

import (
	"fmt"
	"main/app/ErrorHandle"
	"main/app/model"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

// 申請一個 OAUTH 的密鑰 https://twitchapps.com/tmi/

// 套件
// https://pkg.go.dev/github.com/gempir/go-twitch-irc/v3#section-readme

var SendMsgQueue []string
var TwitchClient *twitch.Client

func Init() {
	SendMsgQueue = make([]string, 0)

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	TwitchClient = twitch.NewClient(model.BotSetting.ChatTwitchID, model.BotSetting.TwitchOAth)

	fmt.Printf("[%s] 加入Twitch頻道: %s \n", time.Now().In(time.FixedZone("", +8*3600)).Format("2006-01-02 15:04:05"), model.BotSetting.TargetTwitchID)
	// 加入頻道
	TwitchClient.Join(model.BotSetting.TargetTwitchID)

	err := TwitchClient.Connect()
	if err != nil {
		ErrorHandle.Error.Println("TwitchBot Init error", err)
	} else {
		ErrorHandle.Info.Println("TwitchBot Start", err)
	}

}

// 把要傳的訊息queue起來
func SendMessage(msg string) {
	if len(msg) == 0 {
		return
	}
	SendMsgQueue = append(SendMsgQueue, msg)
}

// 定時執行功能
func TwitchCron() {

	if len(SendMsgQueue) > 0 {
		for _, msg := range SendMsgQueue {
			TwitchClient.Say(model.BotSetting.TargetTwitchID, msg)
		}

		// 清空queue
		SendMsgQueue = make([]string, 0)
	}
}
