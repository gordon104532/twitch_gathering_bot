package Business

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"main/app/ErrorHandle"
	"main/app/TwitchBot"
	"main/app/model"
	"net/http"
	"strings"
)

var donateCache map[string]bool
var opayApiCount int = 0

type opayResp struct {
	LstDonate []donateList `json:"lstDonate"`
	Settings  opaySetting  `json:"settings"`
}
type donateList struct {
	DonateID string `json:"donateid"`
	Name     string `json:"name"`
	Amount   int    `json:"amount"`
	MSG      string `json:"msg"`
}
type opaySetting struct {
	BgColor          string `json:"BgColor"`
	FontAnimate      string `json:"FontAnimate"`
	MsgTemplate      string `json:"MsgTemplate"`
	AlertSound       string `json:"AlertSound"`
	AlertSec         int    `json:"AlertSec"`
	AlertStyle       int    `json:"AlertStyle"`
	TTSStatus        int    `json:"TTSStatus"`
	TTSVolume        int    `json:"TTSVolume"`
	AlertSoundVolume int    `json:"AlertSoundVolume"`
	FontSize         int    `json:"FontSize"`
}

func OpayInit() {
	donateCache = make(map[string]bool)
}
func GetOpayData() {
	u := fmt.Sprintf("https://payment.opay.tw/Broadcaster/CheckDonate/%s", model.BotSetting.Opay.OpayID)
	req, err := http.NewRequest("POST", u, strings.NewReader("__RequestVerificationToken="+model.BotSetting.Opay.OpayToken))
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:   "__RequestVerificationToken_Lw__",
		Value:  model.BotSetting.Opay.OpayCookie,
		Domain: "payment.opay.tw",
		Path:   "/",
		MaxAge: 0,
	}
	req.AddCookie(cookie)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:63.0) Gecko/20100101 Firefox/80.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		ErrorHandle.Error.Println("Opay http response code : ", resp.StatusCode)
		return
	}

	opayApiCount++
	if resp.StatusCode == 200 {
		if opayApiCount == 1 {
			ErrorHandle.Info.Println("歐付寶Donate檢查已成功 1 次")
		}
		if opayApiCount%600 == 0 {
			ErrorHandle.Info.Printf("歐付寶Donate檢查已成功 %d 次(%d分鐘)\n", opayApiCount, opayApiCount/20)
		}
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		ErrorHandle.Error.Println("Opay resp not json")
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ErrorHandle.Error.Println("ioutil.ReadAll err:", err)
		return
	}
	oResp := opayResp{}
	err = json.Unmarshal(bodyBytes, &oResp)
	if err != nil {
		ErrorHandle.Error.Println("Unmarshal err:", err)
		return
	}
	if len(oResp.LstDonate) == 0 {
		return
	}

	for _, v := range oResp.LstDonate {
		// 不在快取中 則加入快取
		if _, ok := donateCache[v.DonateID]; !ok {
			donateCache[v.DonateID] = true
			msg := fmt.Sprintf(model.BotSetting.Opay.OpayMsg, v.Name, v.Amount, v.MSG)

			ErrorHandle.Info.Printf("%s 贊助了 %d 元: %s\n", v.Name, v.Amount, v.MSG)

			TwitchBot.SendMessage(msg)
			TwitchBot.GatheringOpayPoint(v.Name, v.Amount)
		}
	}
}
