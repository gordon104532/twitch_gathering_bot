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
	//初始化活動檔案
	InitGatheringFile()
	InitExpSettingFile()

	// 初始化對話紀錄
	ofaAutoHiList = make(map[string]bool)
	ofaAutoHiList = map[string]bool{
		"nightbot":       true,
		"streamelements": true,
	}
	ofaAutoHiList[model.BotSetting.Twitch.ChatTwitchID] = true

	SendMsgQueue = make([]string, 0)

	// or client := twitch.NewAnonymousClient() for an anonymous user (no write capabilities)
	TwitchClient = twitch.NewClient(model.BotSetting.Twitch.ChatTwitchID, model.BotSetting.Twitch.TwitchOAth)

	TwitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		twitchMessageHandle(TwitchClient, message)
	})

	TwitchClient.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		twitchUserNoticeHandle(TwitchClient, message)
	})

	ErrorHandle.Info.Printf("加入Twitch頻道: %s \n", model.BotSetting.General.TargetTwitchID)
	// 加入頻道
	TwitchClient.Join(model.BotSetting.General.TargetTwitchID)

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
			TwitchClient.Say(model.BotSetting.General.TargetTwitchID, msg)
		}

		// 清空queue
		SendMsgQueue = make([]string, 0)
	}
}

func twitchMessageHandle(client *twitch.Client, message twitch.PrivateMessage) {
	// 自動打招呼
	var context string
	if model.BotSetting.Twitch.AutoHello {
		if _, ok := ofaAutoHiList[message.User.Name]; !ok {
			ofaAutoHiList[message.User.Name] = true
			context = message.User.DisplayName + " " + model.BotSetting.Twitch.AutoHelloMsg + " " + model.BotSetting.Twitch.AutoHelloEmoji
		}
	}

	// 集氣挑戰功能
	if model.BotSetting.GatheringEvent.GatheringSwitch {
		context = cheerEventPoint(client, message)
	}

	// 有內容才發話
	if len(context) > 1 {
		client.Say(message.Channel, context)
	}

}

// 使用者通知處理
func twitchUserNoticeHandle(client *twitch.Client, message twitch.UserNoticeMessage) {
	subEventPoint(client, message)
}
