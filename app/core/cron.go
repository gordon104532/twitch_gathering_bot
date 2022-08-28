package core

import (
	"main/app/Business"
	"main/app/ErrorHandle"
	"main/app/TwitchBot"
	"main/app/model"

	cron "github.com/robfig/cron/v3"
)

func StartCron() {
	c := cron.New()
	ErrorHandle.Info.Printf("背景啟動\n")
	c.AddFunc("@every 4s", func() {
		if model.BotSetting.Opay.CheckDonate {
			Business.GetOpayData()
		}
		if model.BotSetting.Ecpay.CheckDonate {
			Business.GetECpayData()
		}
		TwitchBot.TwitchCron()
	})

	c.Start()
}
