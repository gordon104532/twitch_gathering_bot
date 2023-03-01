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

var AutoHelloList map[string]bool
var SendMsgQueue []string
var TwitchClient *twitch.Client
var retryCount int = 0

func Init() {
	defer func() {
		if x := recover(); x != nil {
			// recovering from a panic; x contains whatever was passed to panic()
			ErrorHandle.Panic.Printf("機器人遇到預期外的錯誤。\n請截圖送到DC，並先重啟機器人。\nerr: %v", x)
		}
	}()

	//初始化活動檔案
	InitGatheringFile()  // 總分
	InitExpSettingFile() // 詳細設定
	InitIndexFile()      // 進度條檔案
	InitControlFile()    // 控制頁檔案

	AutoHelloList = make(map[string]bool)
	AutoHelloList = map[string]bool{
		"nightbot":                              true,
		"streamelements":                        true,
		model.BotSetting.General.TargetTwitchID: true,
		model.BotSetting.Twitch.ChatTwitchID:    true,
	}

	SendMsgQueue = make([]string, 0)

	// 連線物件
	TwitchClient = twitch.NewClient(model.BotSetting.Twitch.ChatTwitchID, model.BotSetting.Twitch.TwitchOAth)

	TwitchClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		twitchMessageHandle(TwitchClient, message)
	})

	TwitchClient.OnUserNoticeMessage(func(message twitch.UserNoticeMessage) {
		twitchUserNoticeHandle(TwitchClient, message)
	})

	// 成功連線
	TwitchClient.OnConnect(func() {
		ErrorHandle.Info.Println("TwitchBot OnConnect")
		retryCount = 0
	})
	// 重連事件
	TwitchClient.OnReconnectMessage(func(message twitch.ReconnectMessage) {
		ErrorHandle.Info.Println("TwitchBot Reconnect Success")
	})

	ErrorHandle.Info.Printf("加入Twitch頻道: %s \n", model.BotSetting.General.TargetTwitchID)
	// 加入頻道
	TwitchClient.Join(model.BotSetting.General.TargetTwitchID)

	err := TwitchClient.Connect()
	if err != nil {
		ErrorHandle.Error.Println("TwitchBot Init error", err)

		// 重新init TwitchBot
		if retryCount < 10 {
			time.Sleep(30 * time.Second)
			retryCount++
			ErrorHandle.Warning.Println("TwitchBot Reconnect Count : ", retryCount)
			go Init()
		} else {
			// 重連失敗送警報
			ErrorHandle.Error.Println("TwitchBot重連失敗。 檢查設定檔、重啟機器人", err)
		}
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

// 訊息處理
func twitchMessageHandle(client *twitch.Client, message twitch.PrivateMessage) {
	var context string

	// 集氣挑戰功能
	if model.BotSetting.GatheringEvent.GatheringSwitch {
		context = cheerEventPoint(client, message)
	}

	// 自動打招呼
	if model.BotSetting.Twitch.AutoHello {
		context = autoHello(message)
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

// 自動打招呼模組
func autoHello(message twitch.PrivateMessage) (context string) {
	if _, ok := AutoHelloList[message.User.Name]; !ok && len(message.User.DisplayName) > 0 {
		AutoHelloList[message.User.Name] = true
		context = fmt.Sprintf("%s %s %s", message.User.DisplayName, model.BotSetting.Twitch.AutoHelloMsg, model.BotSetting.Twitch.AutoHelloEmoji)
	}
	return
}
