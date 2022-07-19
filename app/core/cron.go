package core

import (
	"main/app/Business"
	"main/app/ErrorHandle"
	"main/app/TwitchBot"

	cron "github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()
	ErrorHandle.Info.Printf("背景啟動\n")
	c.AddFunc("@every 3s", func() {
		Business.GetOpayData()
		TwitchBot.TwitchCron()
	})

	c.Start()

}
