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
	r.LoadHTMLGlob("./index.*")

	// init 做一次確認等級
	TwitchBot.GatheringCheckLevelUp()

	// 進度條 // 60秒刷新一次
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

	// 設定控制頁
	// r.LoadHTMLGlob("./view/*.html")
	// r.GET("/control", func(c *gin.Context) {
	// 	c.HTML(http.StatusOK, "index.html", nil)
	// })
	r.Static("/control", "./")

	// 取得基本設定
	r.GET("/setting", func(c *gin.Context) {
		c.JSON(http.StatusOK, model.BotSetting)
	})

	// 更新基本設定
	r.POST("/setting", func(c *gin.Context) {
		// copyCtx := c
		// data, _ := copyCtx.GetRawData()
		// fmt.Println(string(data))
		var getSetting model.BotSettingModel
		if err := c.ShouldBindJSON(&getSetting); err != nil {
			ErrorHandle.Error.Println("更新基本設定 Err:", err)
			c.JSON(500, gin.H{
				"Code": 500,
				"Msg":  err.Error(),
			})
			return
		}

		// 設定為全域參數
		model.BotSetting = getSetting
		// 寫回檔案
		TwitchBot.UpdateBotSetting()
		c.String(http.StatusOK, "ok")
	})

	r.Run(":8787")
}
