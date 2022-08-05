package TwitchBot

import (
	"fmt"
	"main/app/ErrorHandle"
	"main/app/model"
	"strconv"
	"strings"

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

	if model.BotSetting.GatheringEvent.GatheringSwitch {

		if strings.Contains(message.Message, "Cheer") {
			strSlice := strings.Split(message.Message, " ")
			for i := range strSlice {
				if strings.Contains(strSlice[i], "Cheer") {
					var addPoint = 0
					cheerStr := strings.Replace(strSlice[i], "Cheer", "", -1)
					cheerPoint, err := strconv.Atoi(cheerStr)

					if err != nil {
						ErrorHandle.Error.Printf("小奇點加分失敗，請手動換算與加分: %s", message.Message)
					} else {
						addPoint = cheerPoint * model.BotSetting.GatheringEvent.CheerPoint
						model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint
					}
					//活動紀錄
					GatLog("小奇點", addPoint)
				}
			}

			// 檢查升級
			isLevelup, levelUpMsg, _ := gatheringCheckLevelUp()
			if isLevelup {
				context = levelUpMsg
			}
		}

		if strings.Contains(message.User.Name, model.BotSetting.Twitch.ChatTwitchID) {
			if strings.Contains(message.Message, "+") {
				t := strings.Replace(message.Message, "+", "", -1)
				manualPoint, err := strconv.Atoi(t)
				if err != nil {
					ErrorHandle.Error.Printf("手動加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s", message.Message)
				} else {
					model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + manualPoint
				}
				//活動紀錄
				GatLog("手動加", manualPoint)
				// 檢查升級
				isLevelup, levelUpMsg, _ := gatheringCheckLevelUp()
				if isLevelup {
					context = levelUpMsg
				}
			}

			if strings.Contains(message.Message, "-") {
				t := strings.Replace(message.Message, "-", "", -1)
				manualPoint, err := strconv.Atoi(t)
				if err != nil {
					ErrorHandle.Error.Printf("手動減分失敗，請後續開botSetting直接在initPoint減分 並重啟bot: %s", message.Message)
				} else {
					model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint - manualPoint
				}
				//活動紀錄
				GatLog("手動減", manualPoint)
			}
		}

		if strings.Contains(message.Message, "!87") {
			_, _, context = gatheringCheckLevelUp()
		}
		// 寫回總分
		UpdateBotSetting()
	}

	if len(context) > 1 {
		client.Say(message.Channel, context)
	}

}

// 檢查升級
func gatheringCheckLevelUp() (isLevelup bool, levelUpMsg string, checkMsg string) {
	var statment map[int]string = map[int]string{
		0: "不太87",
		1: "有點87",
		2: "稍微87",
		3: "真的87",
		4: "非常87",
		5: "瘋狂87",
	}
	newLevel := 0
	if model.BotSetting.GatheringEvent.InitPoint < model.BotSetting.GatheringEvent.LevelOne {
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelOne, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}
	if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelOne {
		newLevel = 1
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelTwo, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}
	if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelTwo {
		newLevel = 2
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelThree, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}
	if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelThree {
		newLevel = 3
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelFour, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}
	if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelFour {
		newLevel = 4
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelFive, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}
	if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelFive {
		newLevel = 5
		checkMsg = fmt.Sprintf("%d/%d 八七程度%d: %s %s", model.BotSetting.GatheringEvent.InitPoint, model.BotSetting.GatheringEvent.LevelFive, newLevel, statment[newLevel], model.BotSetting.Twitch.AutoHelloEmoji)
	}

	if newLevel > model.GatheringLevel {
		levelUpMsg = fmt.Sprintf("八七升級! 台主現在: lv.%d %s", newLevel, statment[newLevel])
		model.GatheringLevel = newLevel
		isLevelup = true
	}
	return
}

// 歐富寶加分
func GatheringOpayPoint(opayValue int) {
	addPoint := opayValue * model.BotSetting.GatheringEvent.OpayPoint
	model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

	//活動紀錄
	GatLog("歐富寶", addPoint)

	// 檢查升級
	isLevelup, levelUpMsg, _ := gatheringCheckLevelUp()
	if isLevelup {
		SendMessage(levelUpMsg)
	}
}

// 使用者通知處理
func twitchUserNoticeHandle(client *twitch.Client, message twitch.UserNoticeMessage) {
	if message.MsgID == "subgift" || message.MsgID == "resub" || message.MsgID == "sub" {
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + model.BotSetting.GatheringEvent.SubPoint

		GatLog("新訂閱", model.BotSetting.GatheringEvent.SubPoint)
		// 檢查升級
		isLevelup, levelUpMsg, _ := gatheringCheckLevelUp()
		if isLevelup {
			client.Say(message.Channel, levelUpMsg)
		}
	}

}
