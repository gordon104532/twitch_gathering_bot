# 歐法機器人

## 功能項目
1. 歐富寶檢查機器人
2. 自動打招呼

## 打包
- go build -mod=mod -o ofaBot_v2.exe

## 使用說明
1. 解壓縮 (主程式\設定檔\source Code)
2. 把 主程式\設定檔 放在相同資料夾下(或是都丟在桌面)
3. 打開設定檔 編輯下opayID 與 twitchOAth 的內容(如下)
4. 點兩下主程式，看到"背景啟動" 與 "加入Twitch頻道"  且沒有其他錯誤訊息就可以了
5.  source Code 沒有用 跟我一樣

＊如果主程式被防毒軟體殺掉，請就救救他

## 設定檔說明 點開來編輯內容就可以了
botSetting.txt
- targetTwitchID
> 發話的圖奇頻道
- chatTwitchID
> 發話的圖奇的帳號
- opayID 
> 歐富寶提供的動畫頁網址 https://payment.opay.tw/Broadcaster/CheckDonate/這裡的一大串
- twitchOAth
> 從圖奇申請一個 OAUTH 的密鑰 https://twitchapps.com/tmi/ (請妥善保管)
- opayToken、opayCookie
> 這兩個比較複雜暫且保持原樣試試看
- autoHello
> 自動打招呼的開關 true:開 false:關
- autoHelloMsg
> 自動打招呼後接的文字
- autoHelloEmoji
> 自動打招呼後接的貼圖

```
{
    "targetTwitchID":"ofadoraifu",
    "chatTwitchID":"ofadoraifu",
    "checkDonate": true,
    "opayID":"7BF5D2184771810862F90707199",
    "opayToken":"4jR8bmQj%2FyIxCbuczdpPhRFCTTOCGOStfu9laNR9RT1L3ZUgBvJFe9iJtkB%2FIIhCPpNxDwSSaOAqoxxvNOXm7RgGG1200uwIoZPib%2BNiE5%2FQwtaFkYC2wLLIFmMrCqbpMYQFjr6BMLYPJMDdm%2BIvrLBKuKo%3D",
    "opayCookie":"YlSbHQpkKPWeyFc6CVnOZ5skpidCYIxvjK4aCaGs40CCgs9pU/hRDbF3aWzf5QHT/k+p1BFd634KTum6IDkvYsIBsyubKneBuQKHkmo4mu9Vl0LxDYO/8SEFYGo/kHenXUXYbXmsvn9yrE6u5y39uZzak54=",
    "twitchOAth":"oauth:ijf94mqvg2x0u7mv8n7keidwo",
    "autoHello":true,
    "autoHelloMsg": "安安",
    "autoHelloEmoji": "ofadorYeah"
}

```