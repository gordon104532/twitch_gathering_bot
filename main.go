package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"main/app/Business"
	"main/app/ErrorHandle"
	"main/app/Router"
	"main/app/TwitchBot"
	"main/app/core"
	"main/app/model"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

func main() {
	// 初始化log
	ErrorHandle.Init(os.Stdout, os.Stdout, os.Stderr)

	// 監聽外部輸入已關閉
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := &sync.WaitGroup{}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		sig := <-c
		_ = sig
		wg.Done()
	}()
	wg.Add(1)
	ErrorHandle.Info.Printf("Ctrl+C兩次 或 按叉叉 以結束。\n\n")

	//讀取設定檔
	readBotSetting()
	ErrorHandle.Info.Printf("%s Bot Start\n", model.BotSetting.General.TargetTwitchID)

	// 啟用背景
	Business.OpayInit()
	core.StartCron()

	// api服務
	go Router.Router()

	// TwitchBot 啟動
	go TwitchBot.Init()

	wg.Wait()
	ErrorHandle.Info.Printf("Bot End\n")
}

// 讀取txt作為設定檔
func readBotSetting() {
	var tempStr string
	// open the file
	file, err := os.Open("botSetting.txt")

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		tempStr = tempStr + fileScanner.Text()
	}

	err = json.Unmarshal([]byte(tempStr), &model.BotSetting)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s", err)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()

	// 印出開關
	ErrorHandle.Info.Printf("開關-斗內檢查: %v\n", model.BotSetting.Opay.CheckDonate)
	ErrorHandle.Info.Printf("開關-自動安安: %v\n", model.BotSetting.Twitch.AutoHello)
	ErrorHandle.Info.Printf("開關-八七集氣: %v\n", model.BotSetting.GatheringEvent.GatheringSwitch)

	if model.BotSetting.GatheringEvent.GatheringSwitch {
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
		}

		for i := 0; i < 10; i++ {
			if levelPoint[i] > levelPoint[i+1] {
				ErrorHandle.Error.Printf("八七集氣 等級設定有誤: lv.%d 比 lv.%d 小，請關閉bot並檢查設定檔\n", i+1, i)
				fmt.Scanln()
			}
		}
	}

}
