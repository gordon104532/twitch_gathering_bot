# 歐法機器人

## 功能項目
1. 歐富寶檢查機器人
2. 自動打招呼

## 打包
- go build -mod=mod -o ofaBot_v4.exe

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
- checkDonate
    > 斗內檢查開關
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
- opayMsg
    > 自訂歐富寶感謝訊息，%s %d %s 請依序保留,分別代表贊助者姓名, 贊助金額, 附註的話
- gatheringEvent 
    > 集氣挑戰相關設定
- subPoint
    > 一份訂閱加幾分
- cheerPoint
    > 小奇點加分倍率(aka一點加幾分)
- opayPoint
    > 歐富寶加分倍率(aka一塊加幾分)
- initPoint
    > bot啟動時初始分數，後續加分也會持續寫入這一欄
```
{
    "general":{
        "targetTwitchID":"ofadoraifu"
    },
    "opay":{
            "checkDonate": true,
            "opayID":"7BF5D2184771810862F90707199",
            "opayToken":"4jR8bmQj%2FyIxCbuczdpPhRFCTTOCGOStfu9laNR9RT1L3ZUgBvJFe9iJtkB%2FIIhCPpNxDwSSaOAqoxxvNOXm7RgGG1200uwIoZPib%2BNiE5%2FQwtaFkYC2wLLIFmMrCqbpMYQFjr6BMLYPJMDdm%2BIvrLBKuKo%3D",
            "opayCookie":"YlSbHQpkKPWeyFc6CVnOZ5skpidCYIxvjK4aCaGs40CCgs9pU/hRDbF3aWzf5QHT/k+p1BFd634KTum6IDkvYsIBsyubKneBuQKHkmo4mu9Vl0LxDYO/8SEFYGo/kHenXUXYbXmsvn9yrE6u5y39uZzak54=",
            "opayMsg":"/me 感謝 %s 贊助了 %d 元, %s"
    },
    "twitch":{
        "chatTwitchID":"ofadoraifu",
        "twitchOAth":"oauth:ijf94mqvg2x0u7mv8n7keidwo",
        "autoHello":true,
        "autoHelloMsg": "安安",
        "autoHelloEmoji": "ofadorYeah"
    },
    "gatheringEvent":{
        "gatheringSwitch": true,
        "subPoint": 100,
        "cheerPoint": 1,
        "opayPoint": 1,
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
}
```

## 版本紀錄
### ver1
- 自動打招呼
- opay檢查
- twitch聊天室發話
### ver2
- 修正opayID沒有被呼叫
- checkDonate(斗內檢查)開關
- autoHelloMsg(自訂打招呼文字)
- 跳一次api成功訊息
### ver3
- 自訂感謝斗內訊息
- 調整設定檔結構
- 八七集氣活動
### ver3.1
- 修正集氣挑戰追蹤條件，改吃圖奇訊息 MsgID
- 新增歐富寶集氣設定
- template+ bootstrap 做進度條/87
- 升級時會在聊天室發話
### ver3.1
- getTotalPoint.txt 重啟bot後不會清空，會繼續記錄
- 訂閱分數補上訂閱月數與層級

### ver4
- windows 防火牆請改按取消
- getTotalPoint.txt 已經不會因bot重啟而重置了
- 查詢指令 !87 改為 !87LV 與 !87lv
- botSetting中，levelSetting 改成10級結構
- 等級設定檢查(檢查lv大 分數卻較少問題)
- 移除稍微87等說明設定

新的細節設定檔，有需要才設定
- !87 後的emoji 有獨立的設定檔了 預設空
- 各等級的語助詞設定
- 多少小奇點 多少歐富寶 算一次抽獎機會的設定
- sub, subgift 可對多月份訂閱做倍數設定
- resub zero 代表舊的多月續訂不算入集氣
- tier 訂閱層級倍數設定
- 集氣條插件 顏色設定titleColor(標題與目標分數)、barCollor(進度條顏色)、barTxtCollor(進度條上文字顏色)，可寫orange 或 色碼#ed42e5

* botSetting 有改結構 還請用新的
* 集氣調改顏色 就不會有條紋動畫了 不是壞掉
* 匿名贈訂 匿名小奇點 還沒看過事件名稱，目前需要手動加分 

for 新使用 botSetting.txt 必改欄位
targetTwitchID
opayID
chatTwitchID
twitchOAth
