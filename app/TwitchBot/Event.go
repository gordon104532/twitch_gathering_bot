package TwitchBot

import (
	"bufio"
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

	ErrorHandle.Info.Printf("總分異動: %d\n", model.BotSetting.GatheringEvent.InitPoint)
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

// 初始化分數設定檔案
func InitExpSettingFile() {
	filename := "ExpSetting.txt"

	if _, err := os.Stat(filename); err == nil {

	} else if errors.Is(err, os.ErrNotExist) {
		// 無檔則建立
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}

		//msgJSON, _ := json.Marshal(model.DetailSetting)

		var detailRaw string = `
		{
			"checkEmoji": "",
			"progressBar":{
				"titleColor":"#000000",
				"barCollor":"",
				"barTxtCollor": ""
			},
			"logSetting":{
				"cheerPerRoll":100,
				"opayPerRoll":50
			},
			"subgift":{
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"resub":{
				"zero": 0,
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"sub":{
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"tier":{
				"one": 1,
				"two": 2,
				"three": 5
			}
		}`

		werr := os.WriteFile(filename, []byte(detailRaw), 0644)
		if werr != nil {
			ErrorHandle.Error.Println("寫分數設定檔失敗")
		}
		ErrorHandle.Info.Println("建立檔案:" + filename)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Info.Printf("InitGatheringFile else err: %v\n", err)
	}

	// 有檔讀設定
	var tempStr string
	// open the file
	file, err := os.Open(filename)

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		tempStr = tempStr + fileScanner.Text()
	}

	err = json.Unmarshal([]byte(tempStr), &model.DetailSetting)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s", err)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}

// 活動紀錄
func GatLog(event, alias, memo string, point, logTimes int) {
	now := time.Now().Format("01/02 15:04:05")
	messeg := fmt.Sprintf("[%s] 事件:%s, 暱稱:%s, 分數:%d, 總分: %d, 可抽: %d, 備註: %s", now, event, alias, point, model.BotSetting.GatheringEvent.InitPoint, logTimes, memo)
	pointLogger.Println(messeg)
}

// 檢查升級
func GatheringCheckLevelUp() (isLevelup bool, levelUpMsg, checkMsg string, newLevel, basePoint, nextPoint int) {
	newLevel, basePoint, nextPoint = 0, 0, 0
	nowpoint := model.BotSetting.GatheringEvent.InitPoint
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
		if nowpoint < levelPoint[i] {
			newLevel = i - 1
			basePoint = levelPoint[i-1]
			nextPoint = levelPoint[i]
			break
		}
	}

	checkMsg = fmt.Sprintf("目前進度Lv.%d  %d/%d %s", newLevel, model.BotSetting.GatheringEvent.InitPoint, nextPoint, model.DetailSetting.CheckEmoji)

	if newLevel > model.GatheringLevel {
		levelUpMsg = fmt.Sprintf("八七升級！目前集氣等級Lv.%d", newLevel)
		model.GatheringLevel = newLevel
		isLevelup = true
	}
	return
}

// 歐富寶加分
func GatheringOpayPoint(opayName string, opayValue int) {
	addPoint := opayValue * model.BotSetting.GatheringEvent.OpayPoint
	model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

	//活動紀錄
	var logTimes int = opayValue / model.DetailSetting.LogSetting.OpayPerRoll
	GatLog("歐富寶", opayName, fmt.Sprintf("金額:%d", opayValue), addPoint, logTimes)

	// 檢查升級
	isLevelup, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
	if isLevelup {
		SendMessage(levelUpMsg)
	}
}

// 活動訂閱分數處理
func subEventPoint(client *twitch.Client, message twitch.UserNoticeMessage) {
	if message.MsgID == "subgift" || message.MsgID == "resub" || message.MsgID == "sub" || message.MsgID == "submysterygift" {
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
		case "submysterygift":
			event = "贈訂閱"
			giftCount, err := strconv.Atoi(message.MsgParams["msg-param-mass-gift-count"])
			if err != nil {
				ErrorHandle.Error.Printf("贈訂加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s", message.SystemMsg)
			} else {
				// 其實是贈訂份數
				month = giftCount
				event = event + fmt.Sprintf("*%d", giftCount)
			}

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
		}

		switch message.MsgParams["msg-param-sub-plan"] {
		case "1000":
			tier = model.DetailSetting.Tier.One
		case "2000":
			tier = model.DetailSetting.Tier.Two
		case "3000":
			tier = model.DetailSetting.Tier.Three
		}

		addPoint := model.BotSetting.GatheringEvent.SubPoint * month * tier
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

		if addPoint > 0 {
			logTime := month * tier
			GatLog(event, message.User.DisplayName, fmt.Sprintf("月數:%d 層級:%d", month, tier), addPoint, logTime)
		}

		// 檢查升級
		isLevelup, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
		if isLevelup {
			client.Say(message.Channel, levelUpMsg)
		}
	}
}

// 小奇點加分與手動加分
func cheerEventPoint(client *twitch.Client, message twitch.PrivateMessage) (context string) {
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
				var logTimes int = cheerPoint / model.DetailSetting.LogSetting.CheerPerRoll
				GatLog("小奇點", message.User.DisplayName, fmt.Sprintf("%d點", cheerPoint), addPoint, logTimes)
			}
		}

		// 檢查升級
		isLevelup, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
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
			GatLog("手動加", message.User.DisplayName, "", manualPoint, 1)
			// 檢查升級
			isLevelup, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
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
			GatLog("手動減", message.User.DisplayName, "", manualPoint, 1)
		}
	}

	if strings.Contains(message.Message, "!87LV") || strings.Contains(message.Message, "!87lv") {
		_, _, context, _, _, _ = GatheringCheckLevelUp()
	} else {
		// 寫回總分
		UpdateBotSetting()
	}

	return
}
