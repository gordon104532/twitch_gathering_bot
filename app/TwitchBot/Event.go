package TwitchBot

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/app/ErrorHandle"
	"main/app/model"
	"os"
	"path/filepath"
	"time"
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
	_, err := os.Create(filename)
	if err != nil {
		ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
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
func GatLog(event string, point int) {
	now := time.Now().Format("01/02 15:04:05")
	messeg := fmt.Sprintf("[%s] 事件:%s, 分數:%d, 總分: %d", now, event, point, model.BotSetting.GatheringEvent.InitPoint)
	pointLogger.Println(messeg)
}
