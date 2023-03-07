# 集氣機器人
[87集器活動-宣傳頁](https://sites.google.com/view/87cheerup/%E9%A6%96%E9%A0%81?authuser=0)
## 功能項目
1. 歐富寶檢查
2. 綠界檢查
3. 小奇點計分
4. 圖奇訂閱事件計分
5. 計分紀錄
6. 詳細分數倍率設定檔
7. 顯示集氣條進度
8. 自動打招呼功能
## 打包指令
- go build -mod=mod -o 87Bot_v1.7.exe
- 64-bit
> Mac打包方式 $ GOOS=darwin GOARCH=amd64 go build -o bin/app-amd64-darwin 87Bot_v1

## 使用說明
- 前置作業
1. 解壓縮後 只會有一個檔案，點兩下會再變出五個檔案
(註:機器人很容易被防毒軟體抓走，如果有發生這樣的情況再請跟我說)
2. 機器人啟動後會有個黑色的框框，暫時不用管他
3. 用任何瀏覽器打開設定頁 (http://127.0.0.1:8787/control/)
4. 參考歐法的說明填上個欄位後，滾到最下面案儲存，成功會顯示儲存成功
5. 按叉叉關閉機器人 再開啟一次<如果資料都填正確的話 那機器人就開始運作了

- 檢查方式
0. 先把"升級通知" 關閉，並儲存
1. 黑色框框會顯示 歐富寶檢查成功1次
2. 有人訂閱/小奇點 可以在 gatTotalPoint.txt 中察看到
3. (沒開台時 打開升級通知) 可以測試 手動加減分、升級通知、測試斗內功能
4. 分數有變化後進度條約是60秒才會更新，勿緊張
5. 點兩下主程式，看到"背景啟動" 與 "加入Twitch頻道"  且沒有其他錯誤訊息就可以

＊詳細使用說明請看 /doc/bot使用教學.md
＊如果主程式被防毒軟體殺掉，請就救救他(設定排除路徑)

## 設定檔說明 
v1.6 版開始可從網頁進行設定 http://127.0.0.1:8787/control/
> 舊有直接更改檔案方式仍可使用，更改後請重啟程式，botSetting.txt
- templateSwitch
    > 進度條版型 1:舊 2:新
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
- ecpayPoint
    > 綠界加分倍率(aka一塊加幾分)
- initPoint
    > bot啟動時初始分數，後續加分也會持續寫入這一欄
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
- 集氣條插件 顏色設定titleColor(標題與目標分數)、barColor(進度條顏色)、barTxtColor(進度條上文字顏色)，可寫orange 或 色碼#ed42e5

* botSetting 有改結構 還請用新的
* 集氣調改顏色 就不會有條紋動畫了 不是壞掉
* 匿名贈訂 匿名小奇點 還沒看過事件名稱，目前需要手動加分 

for 新使用 botSetting.txt 必改欄位
targetTwitchID
opayID
chatTwitchID
twitchOAth

### 正式版 87Bot_v1(v5) 
- 移除自動打招呼功能
- 移除perRoll功能 在紀錄中寫入"可抽"訊息
- index.tmpl 改由程式建出，引入bootstrap從網路載(不需要app資料夾了)
- 八七集氣條中的文字大一咪咪
    
### 87Bot_v1.1 
- 新增綠界斗內檢查(設定檔新增Ecpay類別、ecpayPoint綠界倍率)
- 斗內檢查由三秒改為五秒一次
- 修正:歐富寶斗內、綠界斗內、訂閱後總分沒寫為botSetting.text的問題
- 註1:主播若不需綠界檢查，可以更新87Bot.exe就好 (botSetting可照舊)
- 註2:若有遇到bot被防毒軟體移除的問題，請將放機器人的資料夾暫時排除

### 87Bot_v1.2
- 斗內檢查由五秒改為四秒一次 (五秒會被windows defender 判定為有害軟體)
- 修正:寫回總分會跑出自動安安設定的問題

### 87Bot_v1.3
- [Fixed] 社群贈訂事件後，仍會跳贈訂事件，導致重複計分
- [Changed] 執行bot後，不會跳看不懂的集氣條相關log(比較整齊)
- ! 結束bot要等約五秒後才關閉是正常現象
- [Fixed] 手動+分-分只認可設定檔中chatTwitchID的帳號操作。改成也認可設定檔中的targetTwitchID
- [Changed] 設定檔格式有誤與設定檔無法開啟，顯示錯誤訊息而非閃退
- [Fixed] 修正社群贈訂沒刪乾淨，導致多算一個月分數的問題 (問題1、2是同問題)
- [Fixed] 小奇點加分找到錯誤後，應要把整句處理完(Ex: Cheer10 Cheer20A Cheer30，結果要是10+30=40)，原本只會記錄到10 與errror便跳出
- [Fixed] 修正貼圖字串包含Cheer被小奇點加分判斷到的問題 (問題5)
- 1.3.1
- [Fixed] 修正小寫cheer 沒有被記入的問題
- [Fixed] 修正亂打200Cheer 卻可以加分的問題
- [Added] 除Cheer外 新增25種小奇點表符判斷，並支援全小寫樣式 (問題3的延伸)
- [Changed] 手動加減0分 不紀錄
- [Changed] 小奇點加分錯誤，只記錄錯誤的字段而非原始整句
- [Changed] 連續手動加減分，應要噴錯
- 1.3.2
- [Fixed] 哭哭饅頭貼圖被誤判為小奇點動作，造成數字檢查時Panic問題
- [Changed] 於goRoutine中加入panic recover()，印出panic訊息並避免閃退
- [Changed] 26種小奇點表符判斷改用迴圈跑，且全用小寫判斷而不是每次比較才轉小寫
### 87Bot_v1.4
- [Fixed] 修正集氣調顯示百分比錯誤
- [無解] resub 單月、多月續訂，若以分享按鈕顯示會被視為 "舊的多月續訂"，暫時解決方法: ExpSetting.txt 中"resub":"zero": 0 改為1。讓舊的多月續訂也計分。
- 目前觀察的確只有初次訂閱會自動跳訊息在聊天室，續訂會以按鈕方式讓觀眾決定何時要分享(但不是非常肯定)。

### 87Bot_v1.5 for 雨鼠
- [Changed] 小奇點不用判斷貼圖了 有原生數字可用 
- [Added] twitch 斷線重連機制30秒一次 重連10次，重連log
- [Added] 接力贈訂事件 
- [Added] 繼續使用贈禮訂閱事件
- [Added] 活動名稱、進度指令可自訂 "gatheringTitle":"活動名稱","queryCommand":"!進度"
- [Added] 新的超級貼圖事件 與計分設定檔 gatheringEvent - stampPoint，一新台幣幾分。若是其他幣別則暫時記0分還請手動加分
- [Changed] 訂閱0分也紀錄(舊月續訂 或 多月)
- [Changed] 小奇點*0分也紀錄
- [Changed] Prime訂閱算層級一
- [Changed] 修正color拼錯字，影響靜態檔與設定檔

### v1.6 for 
- [x] 網頁版設定頁面
- [x] 讀取服務設定檔
- [x] 寫入服務設定檔
- [x] 寫入成功/失敗提示
- [x] 修正寫入成功/失敗提示
- [x] 修正int被axios傳出去變string的問題 => 送出前全部轉換一次
- [x] 輸入框標題底色方塊統一寬度
- [x] index.html index.tmpl 由bot產生
- [x] botSetting設定檔由bot產生
- [x] 升級通知開關
- [x] 自動安安機器人復活
- [x] 進度條底色可自訂，預設加深
> 之後改控制頁 先改/bak/裡的index檔，再貼過去TwitchBot/Event.go裡
> 之後改設定檔 先改Readme說明，再貼過去TwitchBot/Event.go裡

### v1.7
- [x] 進度條優化
- [x] 設定檔可切換版型
- [x] 列出rank
>  預設最多10名
>  因為是讀取gatTotalPoint.txt 檔案，所以如果有用手動加減分的話，還請記得調整該筆"暱稱" 內容為正確的人，不然排名出來最高分可能會是自己

待優化項目
- [ ] 控制頁參數驗證(寫一個功能監聽變化來改變dom?)
- [ ] 讀取細節設定檔
- [ ] 寫入細節設定檔 
- [ ] 新的超級留言事件 與計分設定檔 (尚未找到怎麼抓) => 檢查是不是圖奇放棄了
- [ ] 控制頁區塊收納
- [ ] html js css 檔案分開
- [ ] 感謝追隨
- [ ] 感謝訂閱
- [ ] 感謝揪團
- [ ] paypal斗內監控做得到嗎?
- [ ] 排名顯示在控制頁
- [ ] 排名在控制頁可選列出幾個名次
