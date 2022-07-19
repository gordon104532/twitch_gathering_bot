package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"main/app/Business"
	"main/app/ErrorHandle"
	"main/app/TwitchBot"
	"main/app/core"
	"main/app/model"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"
)

func main() {

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
	fmt.Printf("Ctrl+C兩次 或 按叉叉 以結束。\n\n")
	fmt.Printf("[%s] ofadoraifu Bot Start\n", time.Now().In(time.FixedZone("", +8*3600)).Format("2006-01-02 15:04:05"))

	//讀取設定檔
	readBotSetting()

	// 初始化log
	ErrorHandle.Init(os.Stdout, os.Stdout, os.Stderr)

	// 啟用背景
	Business.OpayInit()
	core.StartCron()

	// TwitchBot 啟動
	go TwitchBot.Init()

	wg.Wait()
	fmt.Printf("[%s] Bot End\n", time.Now().In(time.FixedZone("", +8*3600)).Format("2006-01-02 15:04:05"))
}

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
}
