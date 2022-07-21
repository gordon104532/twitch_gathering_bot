package TwitchBot

import (
	"main/app/ErrorHandle"
	"main/app/model"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

// 申請一個 OAUTH 的密鑰 https://twitchapps.com/tmi/

// 套件
// https://pkg.go.dev/github.com/gempir/go-twitch-irc/v3#section-readme

var SendMsgQueue []string
var TwitchClient *twitch.Client
var ofaAutoHiList map[string]bool

func Init() {

	// 初始化對話紀錄
	ofaAutoHiList = make(map[string]bool)
	ofaAutoHiList = map[string]bool{
		"nightbot":       true,
		"streamelements": true,
	}
	ofaAutoHiList[model.BotSetting.ChatTwitchID] = true

	SendMsgQueue = make([]string, 0)

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	TwitchClient = twitch.NewClient(model.BotSetting.ChatTwitchID, model.BotSetting.TwitchOAth)

	TwitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		twitchMessageHandle(TwitchClient, message)
	})

	ErrorHandle.Info.Printf("加入Twitch頻道: %s \n", model.BotSetting.TargetTwitchID)
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

func twitchMessageHandle(client *twitch.Client, message twitch.PrivateMessage) {

	// 自動打招呼
	var context string
	if model.BotSetting.AutoHello {
		if _, ok := ofaAutoHiList[message.User.Name]; !ok {
			ofaAutoHiList[message.User.Name] = true
			context = message.User.DisplayName + " " + model.BotSetting.AutoHelloMsg + " " + model.BotSetting.AutoHelloEmoji
		}
	}

	if len(context) > 1 {
		client.Say(message.Channel, context)
	}
}
