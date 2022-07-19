package core

import (
	"fmt"
	"main/app/Business"
	"main/app/TwitchBot"
	"time"

	cron "github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()
	fmt.Printf("[%s] 背景啟動\n", time.Now().In(time.FixedZone("", +8*3600)).Format("2006-01-02 15:04:05"))
	c.AddFunc("@every 3s", func() {
		Business.GetOpayData()
		TwitchBot.TwitchCron()
	})

	c.Start()

}
