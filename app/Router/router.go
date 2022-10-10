package Router

import (
	"main/app/ErrorHandle"
	"main/app/TwitchBot"
	"main/app/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() {
	defer func() {
		if x := recover(); x != nil {
			// recovering from a panic; x contains whatever was passed to panic()
			ErrorHandle.Panic.Printf("機器人遇到預期外的錯誤。\n請截圖送到DC，並先重啟機器人。\nerr: %v", x)
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	// r.Use(gin.Logger())

	r.Use(gin.Recovery())

	//根目錄
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"error_code": "",
			"error_text": "",
			"result":     "hi",
		})
	})

	// r.Static("/assetPath", "./app/asset")
	// r.LoadHTMLGlob("./app/view/*")
	r.LoadHTMLGlob("./*.tmpl")

	// init 做一次
	TwitchBot.GatheringCheckLevelUp()
	r.GET("/87", func(c *gin.Context) {
		_, _, _, level, startPoint, endPoint := TwitchBot.GatheringCheckLevelUp()
		model.GatheringLevel = level

		if endPoint < 1 {
			endPoint = 1
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":       model.BotSetting.GatheringEvent.GatheringTitle,
			"startPoint":  startPoint,
			"endPoint":    endPoint,
			"nowPoint":    model.BotSetting.GatheringEvent.InitPoint,
			"level":       level,
			"percent":     ((model.BotSetting.GatheringEvent.InitPoint - startPoint) * 100) / (endPoint - startPoint),
			"titleColor":  model.DetailSetting.ProgressBar.TitleColor,
			"barColor":    model.DetailSetting.ProgressBar.BarColor,
			"barTxtColor": model.DetailSetting.ProgressBar.BarTxtColor,
		})
	})

	r.Run(":8787")
}
