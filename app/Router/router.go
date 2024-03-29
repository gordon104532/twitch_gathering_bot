package Router

import (
	"fmt"
	"main/app/ErrorHandle"
	"main/app/TwitchBot"
	"main/app/model"
	"net/http"
	"strconv"
	"strings"

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

		// 控制進度條版型
		var displayFirst, displaySecond string = "none", "none"
		switch model.BotSetting.General.TemplateSwitch {
		case 1:
			displayFirst = ""
		default: // 預設使用新模板
			displaySecond = ""
		}
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":           model.BotSetting.GatheringEvent.GatheringTitle,
			"startPoint":      startPoint,
			"endPoint":        endPoint,
			"nowPoint":        model.BotSetting.GatheringEvent.InitPoint,
			"level":           level,
			"percent":         ((model.BotSetting.GatheringEvent.InitPoint - startPoint) * 100) / (endPoint - startPoint),
			"titleColor":      model.DetailSetting.ProgressBar.TitleColor,
			"barColor":        model.DetailSetting.ProgressBar.BarColor,
			"barTxtColor":     model.DetailSetting.ProgressBar.BarTxtColor,
			"backgroundColor": model.DetailSetting.ProgressBar.BackgroundColor,
			"displayFirst":    displayFirst,
			"displaySecond":   displaySecond,
			"secondBarColor":  model.DetailSetting.ProgressBar.SecondBarColor,
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

	// 撈出分數紀錄做排行
	r.GET("/rank", func(c *gin.Context) {
		target, err := strconv.Atoi(strings.TrimSpace(c.Query("target")))
		if err != nil && len(c.Query("target")) > 0 {
			c.String(400, fmt.Sprintf("參數錯誤, err: %+v", err))
			return
		}

		// 預設最多10名
		c.String(http.StatusOK, TwitchBot.RankByPoint(target))
	})

	r.Run(":8787")
}
