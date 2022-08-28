```
var gatherBotAMap = map[string]string{
	"bot被防毒軟體刪掉了(簡單的方法)":                                      "解壓縮前到直播結束期間，暫時防毒軟體中的掃描檔案功能",
	"bot被防毒軟體刪掉了(正確的方法)":                                      "於防毒軟體或windows的病毒與威脅防護中，將bot路徑下的87Bot.exe 排除項目 ",
	"點兩下bot，卻沒有跳出黑色框框":                                        "確認botSetting.txt是否跟87Bot放在相同的目錄下?\n 是:Bot設定檔格式有問題(多逗點、少雙引號、少括號)，請索取Bot範例重新填寫資料\n否:請放在相同目錄下",
	"點兩下bot，沒出現防火牆確認":                                         "[正常]不影響使用",
	"點兩下bot，沒出現botSetting.txt 機器人設定檔":                         "[正常]不會出現，請索取設定檔範例",
	"點兩下bot，沒出現gatTotalPoint.txt 集氣記錄檔":                       "從新執行87Bot.exe仍無，\n方法一:於87bot.exe相同目錄下，自己建立一個gatTotalPoint.txt文字檔。並重新執行bot\n方法二:方法一若仍無法寫入，則需再關閉bot前先行備份黑色框框內容",
	"點兩下bot，沒出現index.tmpl 集氣條檔案":                              "從新執行87Bot.exe仍無，\n請向主辦索取檔案並放於相同目錄下",
	"黑色框框，沒跳出歐付寶Donate檢查已成功 1 次":                              "請檢查設定檔checkDonate是否填 [ true, ] (不含雙引號)",
	"黑色框框，跳出歐付寶Donate檢查已成功 600次(30分鐘)":                        "[正常]代表背景仍正確執行中",
	"ERROR:八七集氣 等級設定有誤: lv.2 比 lv.1小":                         "請關閉bot並檢查設定檔中該等級的設定",
	"ERROR: Opay.go:83: Opay resp not json":                   "請確認設定檔opayID，是否與斗內頁網址上相同",
	"ERROR: TwitchBot Init error login authentication failed": "請確認設定檔twitchOAth，是否從圖奇換回的Token相同initPoint",
	"!87lv, !87LV沒反應":                                         "目前設定CD時間30秒，若30秒後仍無反應，請檢查設定檔chatTwitchID與twitchOAth是否正確",
	"+分數, -分數沒反應":                                             "若黑色框框有錯誤訊息請截圖傳到DC，並於botSetting中最末端對總分操作。 存檔後重新執行87Bot.exe",
	"分數到了卻沒有升級":                                               "重啟bot，會重算等級。若仍不行，\n1.請先記得原始分數\n2.botSetting設定檔initPoint歸零\n3.重啟機器人\n4.手動加分",
	"你的問題不在這裡":                                                "不急:請協助截圖畫面與黑色框框發至DC\n急:請使用 !help+{你想問的問題}",
	// 集氣條沒顯示
}
```