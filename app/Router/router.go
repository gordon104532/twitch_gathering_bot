package Router

import (
	"main/app/TwitchBot"
	"main/app/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Router() {
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

	r.Static("/assetPath", "./app/asset")
	r.LoadHTMLGlob("./app/view/*")

	// init 做一次
	TwitchBot.GatheringCheckLevelUp()
	r.GET("/87", func(c *gin.Context) {
		_, _, _, level, startPoint, endPoint := TwitchBot.GatheringCheckLevelUp()
		model.GatheringLevel = level

		if endPoint < 1 {
			endPoint = 1
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":        "八七集氣挑戰",
			"startPoint":   startPoint,
			"endPoint":     endPoint,
			"nowPoint":     model.BotSetting.GatheringEvent.InitPoint,
			"level":        level,
			"percent":      (model.BotSetting.GatheringEvent.InitPoint * 100) / endPoint,
			"titleColor":   model.DetailSetting.ProgressBar.TitleColor,
			"barCollor":    model.DetailSetting.ProgressBar.BarCollor,
			"barTxtCollor": model.DetailSetting.ProgressBar.BarTxtCollor,
		})
	})

	r.Run(":8787")
}
