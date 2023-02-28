package TwitchBot

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"main/app/ErrorHandle"
	"main/app/model"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	twitch "github.com/gempir/go-twitch-irc/v3"
)

func UpdateBotSetting() {

	ErrorHandle.Info.Printf("寫入設定檔，總分: %d\n", model.BotSetting.GatheringEvent.InitPoint)
	msgJSON, _ := json.Marshal(model.BotSetting)
	err := os.WriteFile("botSetting.txt", msgJSON, 0644)
	if err != nil {
		ErrorHandle.Error.Printf("寫回總分失敗 請備份botSetting.txt內容如下:\n %v", model.BotSetting)
	}
}

var pointLogger *log.Logger

// 初始化活動紀錄檔案
func InitGatheringFile() {
	filename := "gatTotalPoint.txt"

	if _, err := os.Stat(filename); err == nil {
		// path/to/whatever exists // 有檔案則離開
		return
	} else if errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}
		// path/to/whatever does *not* exist
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Error.Printf("InitGatheringFile else err: %v\n", err)
	}

CREATE:
	// 建立檔案
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			//建立資料夾
			err = os.MkdirAll(filepath.Dir(filename), 0777)
			if err != nil {
				ErrorHandle.Error.Println("ERROR", "CreateFile: 建立相關資料夾錯誤, "+err.Error())
			}
			goto CREATE
		}
		ErrorHandle.Error.Printf("Failed to opening file: %v", err)
	}

	ErrorHandle.Info.Println("建立檔案:" + filename)
	// 修改檔案權限
	err = f.Chmod(0666)
	if err != nil {
		ErrorHandle.Error.Printf("Failed to chmod file: %v \n", err)
	}

	multiInfo := io.MultiWriter(f)
	pointLogger = log.New(multiInfo, "", 0)
	pointLogger.SetOutput(multiInfo)
}

// 初始化基本設定檔案
func InitBotSettingFile() {
	filename := "botSetting.txt"

	if _, err := os.Stat(filename); err == nil {
		// 有檔案則離開
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// 無檔則建立
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}

		//msgJSON, _ := json.Marshal(model.DetailSetting)

		var detailRaw string = `
		{
			"general":{
				"targetTwitchID":""
			},
			"opay":{
					"checkDonate": true,
					"opayID":"7BF5D2184771810862F9070719909401",
					"opayToken":"4jR8bmQj%2FyIxCbuczdpPhRFCTTOCGOStfu9laNR9RT1L3ZUgBvJFe9iJtkB%2FIIhCPpNxDwSSaOAqoxxvNOXm7RgGG1200uwIoZPib%2BNiE5%2FQwtaFkYC2wLLIFmMrCqbpMYQFjr6BMLYPJMDdm%2BIvrLBKuKo%3D",
					"opayCookie":"YlSbHQpkKPWeyFc6CVnOZ5skpidCYIxvjK4aCaGs40CCgs9pU/hRDbF3aWzf5QHT/k+p1BFd634KTum6IDkvYsIBsyubKneBuQKHkmo4mu9Vl0LxDYO/8SEFYGo/kHenXUXYbXmsvn9yrE6u5y39uZzak54=",
					"opayMsg":"/me 感謝 %s 贊助了 %d 元, %s"
			},
			"ecpay":{
					"checkDonate": true,
					"ecpayID":"EF382462D39404EADA212AF924B879C0",
					"ecpayMsg":"/me 感謝 %s 贊助了 %d 元, %s"
			},
			"twitch":{
				"chatTwitchID":"",
				"twitchOAth":""
			},
			"gatheringEvent":{
				"gatheringSwitch": true,
				"gatheringTitle": "活動名稱",
				"queryCommand": "!進度",
				"stampPoint": 3,
				"subPoint": 150,
				"cheerPoint": 1,
				"opayPoint": 3,
				"ecpayPoint": 3,
				"levelSetting":{
					"lv1":87,
					"lv2":587,
					"lv3":1487,
					"lv4":3487,
					"lv5":5487,
					"lv6":9487,
					"lv7":13487,
					"lv8":15487,
					"lv9":17487,
					"lv10":19487
				},
				"initPoint":0
			}
		}`

		werr := os.WriteFile(filename, []byte(detailRaw), 0644)
		if werr != nil {
			ErrorHandle.Error.Println("寫分數設定檔失敗")
		}
		ErrorHandle.Info.Println("建立檔案:" + filename)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Info.Printf("InitGatheringFile else err: %v\n", err)
	}

	// 有檔讀設定
	var tempStr string
	// open the file
	file, err := os.Open(filename)

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		tempStr = tempStr + fileScanner.Text()
	}

	err = json.Unmarshal([]byte(tempStr), &model.DetailSetting)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s", err)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}

// 初始化分數設定檔案
func InitExpSettingFile() {
	filename := "ExpSetting.txt"

	if _, err := os.Stat(filename); err == nil {
		// 有檔案則離開
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// 無檔則建立
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}

		//msgJSON, _ := json.Marshal(model.DetailSetting)

		var detailRaw string = `
		{
			"checkEmoji": "",
			"progressBar":{
				"titleColor":"#000000",
				"barColor":"",
				"barTxtColor": ""
			},
			"subgift":{
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"resub":{
				"zero": 0,
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"sub":{
				"one": 1,
				"three": 3,
				"six": 6,
				"twelve":12
			},
			"tier":{
				"one": 1,
				"two": 2,
				"three": 5
			}
		}`

		werr := os.WriteFile(filename, []byte(detailRaw), 0644)
		if werr != nil {
			ErrorHandle.Error.Println("寫分數設定檔失敗")
		}
		ErrorHandle.Info.Println("建立檔案:" + filename)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Info.Printf("InitGatheringFile else err: %v\n", err)
	}

	// 有檔讀設定
	var tempStr string
	// open the file
	file, err := os.Open(filename)

	//handle errors while opening
	if err != nil {
		log.Fatalf("Error when opening file: %s", err)
	}

	fileScanner := bufio.NewScanner(file)

	// read line by line
	for fileScanner.Scan() {
		tempStr = tempStr + fileScanner.Text()
	}

	err = json.Unmarshal([]byte(tempStr), &model.DetailSetting)
	if err != nil {
		log.Fatalf("json.Unmarshal: %s", err)
	}

	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}

	file.Close()
}

// 初始化index檔 進度條
func InitIndexFile() {
	filename := "index.tmpl"

	if _, err := os.Stat(filename); err == nil {
		// 有檔案則離開
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// 無檔則建立
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}

		//msgJSON, _ := json.Marshal(model.DetailSetting)

		var detailRaw string = `
		<!doctype html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport"
				content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
			<meta http-equiv="X-UA-Compatible" content="ie=edge">
			<meta http-equiv="refresh" content="60;url=http://127.0.0.1:8787/87">

			<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/css/bootstrap.min.css" integrity="sha384-ggOyR0iXCbMQv3Xipma34MD+dH/1fQ784/j6cY/iJTQUOhcWr7x9JvoRxT2MZw1T" crossorigin="anonymous">
			<script src="https://code.jquery.com/jquery-3.3.1.slim.min.js" integrity="sha384-q8i/X+965DzO0rT7abK41JStQIAqVgRVzpbzo5smXKp4YfRvH+8abtTE1Pi6jizo" crossorigin="anonymous"></script>
			<script src="https://cdn.jsdelivr.net/npm/popper.js@1.14.7/dist/umd/popper.min.js" integrity="sha384-UO2eT0CpHqdSJQ6hJty5KVphtPhzWj9WO1clHTMGa3JDZwrnQq4sF86dIHNDz0W1" crossorigin="anonymous"></script>
			<script src="https://cdn.jsdelivr.net/npm/bootstrap@4.3.1/dist/js/bootstrap.min.js" integrity="sha384-JjSmVgyd0p3pXB1rRibZUAYoIIy6OrQ6VrjIEaFf/nJGzIxFDsf4x0xIM+B07jRM" crossorigin="anonymous"></script>

			<title>{{.title}}</title>
		</head>
		<body style="color:{{.titleColor}};">
		<div style="text-align:center;">{{.title}}(Lv.{{.level}})</div>

		<div class="progress" style="height: 25px;">
		<div class="progress-bar progress-bar-striped progress-bar-animated" role="progressbar" style="width: {{.percent}}%; color:{{.barTxtColor}} ;background-color: {{.barColor}};font-size: larger;" aria-valuenow="{{.percent}}" aria-valuemin="0" aria-valuemax="100">目前: {{.nowPoint}}</div>
		</div>

		<div style="float:left;">{{.startPoint}}</div>
		<div style="float:right;">{{.endPoint}}</div>
		</body>
		</html>
		`

		werr := os.WriteFile(filename, []byte(detailRaw), 0644)
		if werr != nil {
			ErrorHandle.Error.Println("寫入index檔")
		}
		ErrorHandle.Info.Println("建立檔案:" + filename)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Info.Printf("InitIndexFile else err: %v\n", err)
	}
}

// 初始化index檔 設定頁
func InitControlFile() {
	filename := "index.html"

	if _, err := os.Stat(filename); err == nil {
		// 有檔案則離開
		return
	} else if errors.Is(err, os.ErrNotExist) {
		// 無檔則建立
		_, err := os.Create(filename)
		if err != nil {
			ErrorHandle.Error.Println("ERROR", "CreateFile: 建立檔案錯誤, "+err.Error())
		}

		var detailRaw string = `
		<!-- html -->
<head>
   <meta charset="UTF-8">
   <meta name="viewport"
      content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
   <meta http-equiv="X-UA-Compatible" content="ie=edge">
   <script src="https://cdn.jsdelivr.net/npm/vue@2.5.17/dist/vue.js"></script>
   <!-- axios -->
   <script src="https://unpkg.com/axios/dist/axios.min.js"></script>
   <!-- bootstrap -->
   <link href="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0/dist/css/bootstrap.min.css" rel="stylesheet" integrity="sha384-wEmeIV1mKuiNpC+IOBjI7aAzPcEZeedi5yW5f2yOq55WWLwNGmvvx4Um1vskeMj0" crossorigin="anonymous">
   <script src="https://cdn.jsdelivr.net/npm/bootstrap@5.0.0/dist/js/bootstrap.bundle.min.js" integrity="sha384-p34f1UUtsS3wqzfto5wAAmdvj+osOnFyQFpp4Ua3gs/ZVWx6oOypYoCJhGGScy+8" crossorigin="anonymous"></script>
</head>
<div>
   <p>
</div>
<div id="setting-input">
   <h3>主要設定</h3>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">目標Twitch</span>
      <input v-model.trim="setting.general.targetTwitchID" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="圖奇帳號">
   </div>
   <br>
   <h3>歐富寶相關</h3>
   <div class="form-check form-switch">
      <input v-model="setting.opay.checkDonate" class="form-check-input" type="checkbox" id="flexSwitchCheckChecked" checked>
      <label class="form-check-label" for="flexSwitchCheckChecked">是否啟用檢查</label>
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">opayID</span>
      <input v-model.trim="setting.opay.opayID" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: 7BF5D2184771810862F9070719909401">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">opayMsg</span>
      <input v-model.trim="setting.opay.opayMsg" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: /me 感謝 %s 贊助了 %d 元, %s">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">opayToken</span>
      <input v-model.trim="setting.opay.opayToken" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" >
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">opayCookie</span>
      <input v-model.trim="setting.opay.opayCookie" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <br>
   <h3>綠界相關</h3>
   <div class="form-check form-switch">
      <input v-model="setting.ecpay.checkDonate" class="form-check-input" type="checkbox" id="flexSwitchCheckChecked" checked>
      <label class="form-check-label" for="flexSwitchCheckChecked">是否啟用檢查</label>
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">ecpayID</span>
      <input v-model.trim="setting.ecpay.ecpayID" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: EF382462D39404EADA212AF924B879C0">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 110;">ecpayMsg</span>
      <input v-model.trim="setting.ecpay.ecpayMsg" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: /me 感謝 %s 贊助了 %d 元, %s">
   </div>
   <br>
   <h3>圖奇帳號相關</h3>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">發話者圖奇ID</span>
      <input v-model.trim="setting.twitch.chatTwitchID" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="圖奇帳號">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">發話者OAth金鑰</span>
      <input v-model.trim="setting.twitch.twitchOAth" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: oauth:ijf94mqvg2x0u7mv8n7keidwoowk">
   </div>
   <br>
   <h3>分數設定相關</h3>
   <div class="form-check form-switch">
      <input v-model="setting.gatheringEvent.gatheringSwitch" class="form-check-input" type="checkbox" id="flexSwitchCheckChecked" checked>
      <label class="form-check-label" for="flexSwitchCheckChecked" style="width: 150;">是否啟用活動</label>
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">活動名稱</span>
      <input v-model.trim="setting.gatheringEvent.gatheringTitle" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" >
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">分數查詢指令</span>
      <input v-model.trim="setting.gatheringEvent.queryCommand" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" placeholder="ex: !lv">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">訂閱基本分數</span>
      <input v-model.number="setting.gatheringEvent.subPoint" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">小奇點基本分數</span>
      <input v-model.number="setting.gatheringEvent.cheerPoint" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">貼圖基本分數</span>
      <input v-model.number="setting.gatheringEvent.stampPoint" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">歐富寶基本分數</span>
      <input v-model.number="setting.gatheringEvent.opayPoint" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 150;">綠界基本分數</span>
      <input v-model.number="setting.gatheringEvent.ecpayPoint" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <br>
   <h4>各等級分數</h4>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv1</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv1" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv2</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv2" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv3</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv3" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv4</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv4" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv5</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv5" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv6</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv6" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv7</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv7" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv8</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv8" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv9</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv9" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default" style="width: 55;">lv10</span>
      <input v-model.number="setting.gatheringEvent.levelSetting.lv10" type="text" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default">
   </div>
   <br>
   <h4>服務起始分數</h4>
   <div class="input-group mb-3">
      <span class="input-group-text" id="inputGroup-sizing-default"  style="width: 55;">總分</span>
      <input v-model="setting.gatheringEvent.initPoint" type="number" min="0" step="1" class="form-control" aria-label="Sizing example input" aria-describedby="inputGroup-sizing-default" required>
   </div>
   <br>
   <svg xmlns="http://www.w3.org/2000/svg" style="display: none;">
      <symbol id="check-circle-fill" fill="currentColor" viewBox="0 0 16 16">
         <path d="M16 8A8 8 0 1 1 0 8a8 8 0 0 1 16 0zm-3.97-3.03a.75.75 0 0 0-1.08.022L7.477 9.417 5.384 7.323a.75.75 0 0 0-1.06 1.06L6.97 11.03a.75.75 0 0 0 1.079-.02l3.992-4.99a.75.75 0 0 0-.01-1.05z"/>
      </symbol>
      <symbol id="info-fill" fill="currentColor" viewBox="0 0 16 16">
         <path d="M8 16A8 8 0 1 0 8 0a8 8 0 0 0 0 16zm.93-9.412-1 4.705c-.07.34.029.533.304.533.194 0 .487-.07.686-.246l-.088.416c-.287.346-.92.598-1.465.598-.703 0-1.002-.422-.808-1.319l.738-3.468c.064-.293.006-.399-.287-.47l-.451-.081.082-.381 2.29-.287zM8 5.5a1 1 0 1 1 0-2 1 1 0 0 1 0 2z"/>
      </symbol>
      <symbol id="exclamation-triangle-fill" fill="currentColor" viewBox="0 0 16 16">
         <path d="M8.982 1.566a1.13 1.13 0 0 0-1.96 0L.165 13.233c-.457.778.091 1.767.98 1.767h13.713c.889 0 1.438-.99.98-1.767L8.982 1.566zM8 5c.535 0 .954.462.9.995l-.35 3.507a.552.552 0 0 1-1.1 0L7.1 5.995A.905.905 0 0 1 8 5zm.002 6a1 1 0 1 1 0 2 1 1 0 0 1 0-2z"/>
      </symbol>
   </svg>
   <transition name="fade">
      <div v-if="visibility.saveSuccess" class="alert alert-success d-flex align-items-center alert-dismissible fade show" role="alert">
         <svg class="bi flex-shrink-0 me-2" width="24" height="24" role="img" aria-label="Success:">
            <use xlink:href="#check-circle-fill"/>
         </svg>
         <div>
            儲存成功
         </div>
         <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      </div>
   </transition>
   <transition name="fade">
      <div v-if="visibility.saveFailed" class="alert alert-danger d-flex align-items-center alert-dismissible fade show" role="alert">
         <svg class="bi flex-shrink-0 me-2" width="24" height="24" role="img" aria-label="Danger:">
            <use xlink:href="#exclamation-triangle-fill"/>
         </svg>
         <div>
            儲存失敗
         </div>
         <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
      </div>
   </transition>
   <div class="d-grid gap-2 d-md-flex justify-content-md-center">
      <button class="btn btn-primary btn-lg me-md-2" type="button" v-on:click="apiUpdateSetting">儲存</button>
      <button class="btn btn-secondary btn-lg" type="button" v-on:click="apiGetSetting">還原</button>
   </div>
   <br>
</div>
<!-- js -->
<script>
   var settingInput =new Vue({
     el:'#setting-input',
     data:{
      setting:{
         general:{
           targetTwitchID:'讀不到設定檔'
         },
         opay:{
           checkDonate:true,
           opayID:"7BF5D2184771810862F9070719909401",
           opayMsg:"/me 感謝 %s 贊助了 %d 元, %s",
           opayToken: "4jR8bmQj%2FyIxCbuczdpPhRFCTTOCGOStfu9laNR9RT1L3ZUgBvJFe9iJtkB%2FIIhCPpNxDwSSaOAqoxxvNOXm7RgGG1200uwIoZPib%2BNiE5%2FQwtaFkYC2wLLIFmMrCqbpMYQFjr6BMLYPJMDdm%2BIvrLBKuKo%3D",
           opayCookie:"YlSbHQpkKPWeyFc6CVnOZ5skpidCYIxvjK4aCaGs40CCgs9pU/hRDbF3aWzf5QHT/k+p1BFd634KTum6IDkvYsIBsyubKneBuQKHkmo4mu9Vl0LxDYO/8SEFYGo/kHenXUXYbXmsvn9yrE6u5y39uZzak54=",
         },
         ecpay:{
           checkDonate:true,
           ecpayID:"EF382462D39404EADA212AF924B879C0",
           ecpayMsg:"/me 感謝 %s 贊助了 %d 元, %s"
         },
         twitch:{
           chatTwitchID:"",
           twitchOAth:"oauth:ijf94mqvg2x0u7mv8n7keidwoowk"
         },
         gatheringEvent:{
           gatheringSwitch:true,
           gatheringTitle:"預設名稱",
           queryCommand:"!87",
           subPoint:0,
           cheerPoint:0,
           stampPoint:0,
           opayPoint:0,
           ecpayPoint:0,
           levelSetting:{
             lv1:10,
             lv2:20,
             lv3:30,
             lv4:40,
             lv5:50,
             lv6:60,
             lv7:70,
             lv8:80,
             lv9:90,
             lv10:100,
           },
           initPoint:100 
         },
      },
      visibility:{
         saveSuccess: false,
         saveFailed: false,
      }
     },
     const: {},
     methods: {
       // 取得設定檔
       apiGetSetting() {
         axios
         .get('/setting')
         .then( response => {
             this.setting  = response.data;
           })
         .catch(function (error) { // 請求失敗處理
           console.log("get setting err: "+error);
         });
       },
   
       // 更新設定檔
       apiUpdateSetting() {
         // 字串轉回數字
         this.transStringToInt();

         axios
         .post('/setting', this.setting)
         .then((response) => {
            // console.log(response.status);
            settingInput.visibility.saveSuccess = true;
            setTimeout(() => {
               settingInput.visibility.saveSuccess = false;
            }, 3000)
         })
         .catch((error) => {
            console.log(error);
            settingInput.visibility.saveFailed = true
            setTimeout(() => {
               settingInput.visibility.saveFailed = false;
            }, 3000)
         });
       },

       // 網頁輸入的數字存資料時從字串改回數字
       transStringToInt(){

         for (let key of Object.keys(this.setting.gatheringEvent)) {
            if (key !== "gatheringSwitch" && key !== "gatheringTitle" && key !== "queryCommand" && key !== "levelSetting"){
               this.setting.gatheringEvent[key] = parseInt(this.setting.gatheringEvent[key]);
            }

            if (key === "levelSetting"){
               for (let nextKey of Object.keys(this.setting.gatheringEvent.levelSetting)) {
                  this.setting.gatheringEvent.levelSetting[nextKey] = parseInt(this.setting.gatheringEvent.levelSetting[nextKey]);
               }
            }
         }
       }
     },

     created(){
       this.apiGetSetting()
     }
     });
</script>
		`

		werr := os.WriteFile(filename, []byte(detailRaw), 0644)
		if werr != nil {
			ErrorHandle.Error.Println("寫入index檔")
		}
		ErrorHandle.Info.Println("建立檔案:" + filename)

	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		ErrorHandle.Info.Printf("InitIndexFile else err: %v\n", err)
	}
}

// 活動紀錄
func GatLog(event, alias, memo string, point int) {
	now := time.Now().Format("01/02 15:04:05")
	message := fmt.Sprintf("[%s] 事件:%s, 暱稱:%s, 分數:%d, 總分: %d, 備註: %s", now, event, alias, point, model.BotSetting.GatheringEvent.InitPoint, memo)
	pointLogger.Println(message)
	// 寫回總分
	UpdateBotSetting()
}

// 檢查升級
func GatheringCheckLevelUp() (isLevelUp bool, levelUpMsg, checkMsg string, newLevel, basePoint, nextPoint int) {
	newLevel, basePoint, nextPoint = 0, 0, 0
	nowPoint := model.BotSetting.GatheringEvent.InitPoint
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
		11: 999999999,
	}

	for i := 0; i < 12; i++ {
		if nowPoint < levelPoint[i] {
			newLevel = i - 1
			basePoint = levelPoint[i-1]
			nextPoint = levelPoint[i]
			break
		}
	}

	checkMsg = fmt.Sprintf("目前進度Lv.%d  %d/%d %s", newLevel, model.BotSetting.GatheringEvent.InitPoint, nextPoint, model.DetailSetting.CheckEmoji)

	if newLevel > model.GatheringLevel {
		levelUpMsg = fmt.Sprintf("目前等級Lv.%d", newLevel)
		model.GatheringLevel = newLevel
		isLevelUp = true
	}
	return
}

// 歐富寶加分
func GatheringDonatePoint(platform, donorName string, donorValue int) {
	var addPoint int
	switch platform {
	case "opay":
		{
			addPoint = donorValue * model.BotSetting.GatheringEvent.OpayPoint
		}
	case "ecpay":
		{
			addPoint = donorValue * model.BotSetting.GatheringEvent.EcpayPoint
		}
	}

	model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

	if addPoint > 0 {
		//活動紀錄
		GatLog(platform, donorName, fmt.Sprintf("金額:%d", donorValue), addPoint)
	}
	// 檢查升級
	isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
	if isLevelUp {
		SendMessage(levelUpMsg)
	}
}

// 活動訂閱分數處理
func subEventPoint(client *twitch.Client, message twitch.UserNoticeMessage) {
	if message.MsgID == "subgift" || message.MsgID == "resub" || message.MsgID == "sub" { // || message.MsgID == "submysterygift"
		var event string
		var month, tier int = 1, 1
		switch message.MsgID {
		case "subgift":
			event = "贈訂閱"
			switch message.MsgParams["msg-param-gift-months"] {
			case "1":
				month = model.DetailSetting.Subgift.One
			case "3":
				month = model.DetailSetting.Subgift.Three
			case "6":
				month = model.DetailSetting.Subgift.Six
			case "12":
				month = model.DetailSetting.Subgift.Twelve
			}
		// 社群贈訂submysterygift後 仍會跳subgift 事件造成重複計分
		// case "submysterygift":
		// 	event = "贈訂閱"
		// 	giftCount, err := strconv.Atoi(message.MsgParams["msg-param-mass-gift-count"])
		// 	if err != nil {
		// 		ErrorHandle.Error.Printf("贈訂加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s", message.SystemMsg)
		// 	} else {
		// 		// 其實是贈訂份數
		// 		month = giftCount
		// 		event = event + fmt.Sprintf("*%d", giftCount)
		// 	}
		case "resub":
			event = "續訂閱"
			switch message.MsgParams["msg-param-multimonth-duration"] {
			case "0":
				// 舊的多月份訂閱續訂 不加分
				month = model.DetailSetting.Resub.Zero
			case "1":
				month = model.DetailSetting.Resub.One
			case "3":
				month = model.DetailSetting.Resub.Three
			case "6":
				month = model.DetailSetting.Resub.Six
			case "12":
				month = model.DetailSetting.Resub.Twelve
			}
		case "sub":
			event = "新訂閱"
			switch message.MsgParams["msg-param-multimonth-duration"] {
			case "1":
				month = model.DetailSetting.Sub.One
			case "3":
				month = model.DetailSetting.Sub.Three
			case "6":
				month = model.DetailSetting.Sub.Six
			case "12":
				month = model.DetailSetting.Sub.Twelve
			}
		case "midnightsquid": // 超級貼圖
			// 超級貼圖處理
			superStampProgress(client, message)

		case "giftpaidupgrade": // 繼續使用贈禮訂閱
			event = "繼續訂"
			// 原生無相關資訊，預設層級一、一個月
			month = model.DetailSetting.Sub.One
		case "standardpayforward", "communitypayforward": // 接力贈訂
			event = "接訂閱"
			// 原生無相關資訊，預設層級一、一個月
			month = model.DetailSetting.Sub.One
		}

		switch message.MsgParams["msg-param-sub-plan"] {
		case "1000", "Prime":
			tier = model.DetailSetting.Tier.One
		case "2000":
			tier = model.DetailSetting.Tier.Two
		case "3000":
			tier = model.DetailSetting.Tier.Three
		default:
			tier = model.DetailSetting.Tier.One
		}

		addPoint := model.BotSetting.GatheringEvent.SubPoint * month * tier
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

		GatLog(event, message.User.DisplayName, fmt.Sprintf("月數:%d 層級:%d", month, tier), addPoint)

		// 檢查升級
		isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
		if isLevelUp {
			client.Say(message.Channel, levelUpMsg)
		}
	}
}

// 查詢指令CD時間
var commandCD bool = true

func ContainsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
}

// 小奇點加分與手動加分
func cheerEventPoint(client *twitch.Client, message twitch.PrivateMessage) (context string) {
	// 小奇點
	if message.Bits != 0 {
		var addPoint int = 0
		addPoint = message.Bits * model.BotSetting.GatheringEvent.CheerPoint
		model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + addPoint

		//活動紀錄
		GatLog("小奇點", message.User.DisplayName, fmt.Sprintf("%d點", message.Bits), addPoint)

		// 檢查升級
		isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
		if isLevelUp {
			context = levelUpMsg
		}
	}

	if strings.Contains(message.User.Name, model.BotSetting.Twitch.ChatTwitchID) || strings.Contains(message.User.Name, model.BotSetting.General.TargetTwitchID) {
		if strings.Contains(message.Message, "+") {
			t := strings.Replace(message.Message, "+", "", -1)
			manualPoint, err := strconv.Atoi(t)

			if err != nil {
				ErrorHandle.Error.Printf("手動加分失敗，請後續開botSetting直接在initPoint加分 並重啟bot: %s\nerr:%v", message.Message, err)
				return
			} else {
				if manualPoint == 0 {
					return
				}
				model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint + manualPoint
			}
			//活動紀錄
			GatLog("手動加", message.User.DisplayName, "", manualPoint)
			// 檢查升級
			isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
			if isLevelUp {
				context = levelUpMsg
			}
		}

		if strings.Contains(message.Message, "-") {
			t := strings.Replace(message.Message, "-", "", -1)
			manualPoint, err := strconv.Atoi(t)

			if err != nil {
				ErrorHandle.Error.Printf("手動減分失敗，請後續開botSetting直接在initPoint減分 並重啟bot: %s\nerr:%v", message.Message, err)
				return
			} else {
				if manualPoint == 0 {
					return
				}
				model.BotSetting.GatheringEvent.InitPoint = model.BotSetting.GatheringEvent.InitPoint - manualPoint
			}

			//活動紀錄
			GatLog("手動減", message.User.DisplayName, "", manualPoint)
		}
	}

	// 依設定檔指令做查詢
	if strings.Contains(message.Message, model.BotSetting.GatheringEvent.QueryCommand) && model.BotSetting.GatheringEvent.QueryCommand != "" {
		if commandCD {
			_, _, context, _, _, _ = GatheringCheckLevelUp()
			commandCD = false
			go cdCoolDown()
		}

	}

	return
}

// CD 時間倒數
func cdCoolDown() {
	time.Sleep(time.Second * 30)
	commandCD = true
}

// 超級貼圖處理
func superStampProgress(client *twitch.Client, message twitch.UserNoticeMessage) {
	event := "贈貼圖"
	rawAmount := message.MsgParams["msg-param-amount"]

	amount, err := strconv.Atoi(rawAmount)
	if err != nil {
		ErrorHandle.Error.Printf("超級貼圖加分失敗: 原始金額:%s \nerr:%v", rawAmount, err)
		return
	}

	showAmount := amount / 100
	var addPoint int = 0
	currency := message.MsgParams["msg-param-currency"]
	if currency == "TWD" {
		addPoint = showAmount * model.BotSetting.GatheringEvent.StampPoint
	}

	// 活動紀錄
	GatLog(event, message.User.DisplayName, fmt.Sprintf("金額:%d 幣別:%s", showAmount, currency), addPoint)

	// 檢查升級
	isLevelUp, levelUpMsg, _, _, _, _ := GatheringCheckLevelUp()
	if isLevelUp {
		client.Say(message.Channel, levelUpMsg)
	}
}
