# 歐法機器人

## 功能項目
1. 歐富寶檢查機器人

## 打包
- go build -mod=mod -o ofaBot.exe

## 起服務
1. 打開 botSetting.txt 填上對應的內容
2. 把ofaBot.exe 跟 botSetting.txt 

## 設定檔說明
botSetting.txt
- targetTwitchID
> 發話的圖奇頻道
- chatTwitchID
> 發話的圖奇的帳號
- opayID 
> 歐富寶提供的動畫頁網址 https://payment.opay.tw/Broadcaster/CheckDonate/{這裡的一大串}
- twitchOAth
> 從圖奇申請一個 OAUTH 的密鑰 https://twitchapps.com/tmi/ (請妥善保管)
- opayToken、opayCookie
> 這兩個比較複雜暫且保持原樣試試看
- autoHello
> 自動打招呼的開關 true:開 false:關

```
{
    "targetTwitchID":"ofadoraifu",
    "opayID":"7BF5D2184771810862F90707199",
    "opayToken":"4jR8bmQj%2FyIxCbuczdpPhRFCTTOCGOStfu9laNR9RT1L3ZUgBvJFe9iJtkB%2FIIhCPpNxDwSSaOAqoxxvNOXm7RgGG1200uwIoZPib%2BNiE5%2FQwtaFkYC2wLLIFmMrCqbpMYQFjr6BMLYPJMDdm%2BIvrLBKuKo%3D",
    "opayCookie":"YlSbHQpkKPWeyFc6CVnOZ5skpidCYIxvjK4aCaGs40CCgs9pU/hRDbF3aWzf5QHT/k+p1BFd634KTum6IDkvYsIBsyubKneBuQKHkmo4mu9Vl0LxDYO/8SEFYGo/kHenXUXYbXmsvn9yrE6u5y39uZzak54=",
    "chatTwitchID":"ofadoraifu",
    "twitchOAth":"oauth:ijf94mqvg2x0u7mv8n7keidwo",
    "autoHello":true
}

```