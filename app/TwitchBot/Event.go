package TwitchBot

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/app/ErrorHandle"
	"main/app/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

func UpdateBotSetting() {
	ErrorHandle.Info.Printf("寫入設定檔，總分: %d\n", model.BotSetting.GatheringEvent.InitPoint)
	msgJSON, _ := json.Marshal(model.BotSetting)
	err := os.WriteFile("botSetting.txt", msgJSON, 0644)
	if err != nil {
		ErrorHandle.Error.Printf("寫回總分失敗 請備份botSetting.txt內容如下:\n %v", model.BotSetting)
	}
}

var pointLogger *log.Logger

// 初始化活動紀錄檔案
func InitGatheringFile() {
	filename := "gatTotalPoint.txt"

	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists
	} else if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}
		// path/to/whatever does *not* exist
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Error.Printf("InitGatheringFile else err: %v\n", err)
	}

CREATE:
	// 建立檔案
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			//建立資料夾
			err = os.MkdirAll(filepath.Dir(filename), 0777)
			if err != nil {
				ErrorHandle.Error.Println("ERROR", "CreateFile: 建立相關資料夾錯誤, "+err.Error())
			}
			goto CREATE
		}
		ErrorHandle.Error.Printf("Failed to opening file: %v", err)
	}

	ErrorHandle.Info.Println("建立檔案:" + filename)
	// 修改檔案權限
	err = f.Chmod(0666)
	if err != nil {
		ErrorHandle.Error.Printf("Failed to chmod file: %v \n", err)
	}

	multiInfo := io.MultiWriter(f)
	pointLogger = log.New(multiInfo, "", 0)
	pointLogger.SetOutput(multiInfo)
}

// 活動紀錄
func GatLog(event, alias, memo string, point int) {
	now := time.Now().Format("01/02 15:04:05")
	message := fmt.Sprintf("[%s] 事件: %s, 暱稱: %s, 分數: %d, 總分: %d, 備註: %s", now, event, alias, point, model.BotSetting.GatheringEvent.InitPoint, memo)
	pointLogger.Println(message)
	// 寫回總分
	UpdateBotSetting()
}

// 檢查升級
func GatheringCheckLevelUp() (isLevelUp bool, levelUpMsg, checkMsg string, newLevel, basePoint, nextPoint int) {
	newLevel, basePoint, nextPoint = 0, 0, 0
	nowPoint := model.BotSetting.GatheringEvent.InitPoint
	levelPoint := map[int]int{
		0:  0,
		1:  model.BotSetting.GatheringEvent.LevelSetting.Lv1,
		2:  model.BotSetting.GatheringEvent.LevelSetting.Lv2,
		3:  model.BotSetting.GatheringEvent.LevelSetting.Lv3,
		4:  model.BotSetting.GatheringEvent.LevelSetting.Lv4,
		5:  model.BotSetting.GatheringEvent.LevelSetting.Lv5,
		6:  model.BotSetting.GatheringEvent.LevelSetting.Lv6,
		7:  model.BotSetting.GatheringEvent.LevelSetting.Lv7,
		8:  model.BotSetting.GatheringEvent.LevelSetting.Lv8,
		9:  model.BotSetting.GatheringEvent.LevelSetting.Lv9,
		10: model.BotSetting.GatheringEvent.LevelSetting.Lv10,
		11: 999999999,
	}

	for i := 0; i < 12; i++ {
		if nowPoint < levelPoint[i] {
			newLevel = i - 1
			basePoint = levelPoint[i-1]
			nextPoint = levelPoint[i]
			break
		}
	}

	checkMsg = fmt.Sprintf("目前進度Lv.%d  %d/%d %s", newLevel, model.BotSetting.GatheringEvent.InitPoint, nextPoint, model.DetailSetting.CheckEmoji)

	if newLevel > model.GatheringLevel {
		levelUpMsg = fmt.Sprintf("目前等級Lv.%d %d/%d", newLevel, model.BotSetting.GatheringEvent.InitPoint, nextPoint)
		model.GatheringLevel = newLevel
		isLevelUp = true
	}
	return
}

// 歐富寶加分
func GatheringDonatePoint(platform, donorName string, donorValue int) {
	var addPoint int
	switch platform {
	case "opay":
		{
			addPoint = donorValue * model.BotSetting.GatheringEvent.OpayPoint
		}
	case "ecpay":
		{
			addPoint = donorValue * model.BotSetting.GatheringEvent.EcpayPoint
		}
	}

	model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

	if addPoint > 0 {
		//活動紀錄
		GatLog(platform, donorName, fmt.Sprintf("金額:%d", donorValue), addPoint)
	}
	// 檢查升級
	isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
	if isLevelUp && model.BotSetting.Twitch.LevelUpNotice {
		SendMessage(levelUpMsg)
	}
}

// 活動訂閱分數處理
func subEventPoint(client *twitch.Client, message twitch.UserNoticeMessage) {
	if message.MsgID == "subgift" || message.MsgID == "resub" || message.MsgID == "sub" { // || message.MsgID == "submysterygift"
		var event string
		var month, tier int = 1, 1
		switch message.MsgID {
		case "subgift":
			event = "贈訂閱"
			switch message.MsgParams["msg-param-gift-months"] {
			case "1":
				month = model.DetailSetting.Subgift.One
			case "3":
				month = model.DetailSetting.Subgift.Three
			case "6":
				month = model.DetailSetting.Subgift.Six
			case "12":
				month = model.DetailSetting.Subgift.Twelve
			}
		// 社群贈訂submysterygift後 仍會跳subgift 事件造成重複計分
		// case "submysterygift":
		// 	event = "贈訂閱"
		// 	giftCount, err := strconv.Atoi(message.MsgParams["msg-param-mass-gift-count"])
		// 	if err != nil {
		// 		ErrorHandle.Error.Printf("贈訂加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s", message.SystemMsg)
		// 	} else {
		// 		// 其實是贈訂份數
		// 		month = giftCount
		// 		event = event + fmt.Sprintf("*%d", giftCount)
		// 	}
		case "resub":
			event = "續訂閱"
			switch message.MsgParams["msg-param-multimonth-duration"] {
			case "0":
				// 舊的多月份訂閱續訂 不加分
				month = model.DetailSetting.Resub.Zero
			case "1":
				month = model.DetailSetting.Resub.One
			case "3":
				month = model.DetailSetting.Resub.Three
			case "6":
				month = model.DetailSetting.Resub.Six
			case "12":
				month = model.DetailSetting.Resub.Twelve
			}
		case "sub":
			event = "新訂閱"
			switch message.MsgParams["msg-param-multimonth-duration"] {
			case "1":
				month = model.DetailSetting.Sub.One
			case "3":
				month = model.DetailSetting.Sub.Three
			case "6":
				month = model.DetailSetting.Sub.Six
			case "12":
				month = model.DetailSetting.Sub.Twelve
			}
		case "midnightsquid": // 超級貼圖
			// 超級貼圖處理
			superStampProgress(client, message)

		case "giftpaidupgrade": // 繼續使用贈禮訂閱
			event = "繼續訂"
			// 原生無相關資訊，預設層級一、一個月
			month = model.DetailSetting.Sub.One
		case "standardpayforward", "communitypayforward": // 接力贈訂
			event = "接訂閱"
			// 原生無相關資訊，預設層級一、一個月
			month = model.DetailSetting.Sub.One
		}

		switch message.MsgParams["msg-param-sub-plan"] {
		case "1000", "Prime":
			tier = model.DetailSetting.Tier.One
		case "2000":
			tier = model.DetailSetting.Tier.Two
		case "3000":
			tier = model.DetailSetting.Tier.Three
		default:
			tier = model.DetailSetting.Tier.One
		}

		addPoint := model.BotSetting.GatheringEvent.SubPoint * month * tier
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

		GatLog(event, message.User.DisplayName, fmt.Sprintf("月數:%d, 層級:%d", month, tier), addPoint)

		// 檢查升級
		isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
		if isLevelUp && model.BotSetting.Twitch.LevelUpNotice {
			client.Say(message.Channel, levelUpMsg)
		}
	}
}

// 查詢指令CD時間
var commandCD bool = true

func ContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

// 小奇點加分與手動加分
func cheerEventPoint(client *twitch.Client, message twitch.PrivateMessage) (context string) {
	// 小奇點
	if message.Bits != 0 {
		var addPoint int = 0
		addPoint = message.Bits * model.BotSetting.GatheringEvent.CheerPoint
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

		//活動紀錄
		GatLog("小奇點", message.User.DisplayName, fmt.Sprintf("%d點", message.Bits), addPoint)

		// 檢查升級
		isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
		if isLevelUp && model.BotSetting.Twitch.LevelUpNotice {
			context = levelUpMsg
		}
	}

	if strings.Contains(message.User.Name, model.BotSetting.Twitch.ChatTwitchID) || strings.Contains(message.User.Name, model.BotSetting.General.TargetTwitchID) {
		if strings.Contains(message.Message, "+") {
			t := strings.Replace(message.Message, "+", "", -1)
			manualPoint, err := strconv.Atoi(t)

			if err != nil {
				ErrorHandle.Error.Printf("手動加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s\nerr:%v", message.Message, err)
				return
			} else {
				if manualPoint == 0 {
					return
				}
				model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + manualPoint
			}

			//活動紀錄
			GatLog("手動加", message.User.DisplayName, "", manualPoint)
			// 檢查升級
			isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
			if isLevelUp && model.BotSetting.Twitch.LevelUpNotice {
				context = levelUpMsg
			}
		}

		if strings.Contains(message.Message, "-") {
			t := strings.Replace(message.Message, "-", "", -1)
			manualPoint, err := strconv.Atoi(t)

			if err != nil {
				ErrorHandle.Error.Printf("手動減分失敗，請後續開botSetting直接在initPoint減分 並重啟bot: %s\nerr:%v", message.Message, err)
				return
			} else {
				if manualPoint == 0 {
					return
				}
				model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint - manualPoint
			}

			//活動紀錄
			GatLog("手動減", message.User.DisplayName, "", manualPoint)
		}
	}

	// 依設定檔指令做查詢
	if strings.Contains(message.Message, model.BotSetting.GatheringEvent.QueryCommand) && model.BotSetting.GatheringEvent.QueryCommand != "" {
		if commandCD {
			_, _, context, _, _, _ = GatheringCheckLevelUp()
			commandCD = false
			go cdCoolDown()
		}

	}

	return
}

// CD 時間倒數
func cdCoolDown() {
	time.Sleep(time.Second * 30)
	commandCD = true
}

// 超級貼圖處理
func superStampProgress(client *twitch.Client, message twitch.UserNoticeMessage) {
	event := "贈貼圖"
	rawAmount := message.MsgParams["msg-param-amount"]

	amount, err := strconv.Atoi(rawAmount)
	if err != nil {
		ErrorHandle.Error.Printf("超級貼圖加分失敗: 原始金額:%s \nerr:%v", rawAmount, err)
		return
	}

	showAmount := amount / 100
	var addPoint int = 0
	currency := message.MsgParams["msg-param-currency"]
	if currency == "TWD" {
		addPoint = showAmount * model.BotSetting.GatheringEvent.StampPoint
	}

	// 活動紀錄
	GatLog(event, message.User.DisplayName, fmt.Sprintf("金額:%d 幣別:%s", showAmount, currency), addPoint)

	// 檢查升級
	isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
	if isLevelUp && model.BotSetting.Twitch.LevelUpNotice {
		client.Say(message.Channel, levelUpMsg)
	}
}
