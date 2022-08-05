package Router

import (
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
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + 133
		c.JSON(200, gin.H{
			"error_code": "",
			"error_text": "",
			"result":     "hi",
		})
	})

	r.Static("/assetPath", "./app/asset")
	r.LoadHTMLGlob("./app/view/*")
	r.GET("/87", func(c *gin.Context) {
		var startPoint, endPoint int
		if model.BotSetting.GatheringEvent.InitPoint < model.BotSetting.GatheringEvent.LevelOne {
			model.GatheringLevel = 0
			startPoint = 0
			endPoint = model.BotSetting.GatheringEvent.LevelOne
		}
		if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelOne {
			model.GatheringLevel = 1
			startPoint = model.BotSetting.GatheringEvent.LevelOne
			endPoint = model.BotSetting.GatheringEvent.LevelTwo
		}
		if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelTwo {
			model.GatheringLevel = 2
			startPoint = model.BotSetting.GatheringEvent.LevelTwo
			endPoint = model.BotSetting.GatheringEvent.LevelThree
		}
		if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelThree {
			model.GatheringLevel = 3
			startPoint = model.BotSetting.GatheringEvent.LevelThree
			endPoint = model.BotSetting.GatheringEvent.LevelFour
		}
		if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelFour {
			model.GatheringLevel = 4
			startPoint = model.BotSetting.GatheringEvent.LevelFour
			endPoint = model.BotSetting.GatheringEvent.LevelFive
		}
		if model.BotSetting.GatheringEvent.InitPoint >= model.BotSetting.GatheringEvent.LevelFive {
			model.GatheringLevel = 5
			startPoint = model.BotSetting.GatheringEvent.LevelFive
			endPoint = 878787
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":      "八七集氣挑戰",
			"startPoint": startPoint,
			"endPoint":   endPoint,
			"nowPoint":   model.BotSetting.GatheringEvent.InitPoint,
			"level":      model.GatheringLevel,
			"percent":    (model.BotSetting.GatheringEvent.InitPoint * 100) / endPoint,
		})
	})

	r.Run(":8787")
}
