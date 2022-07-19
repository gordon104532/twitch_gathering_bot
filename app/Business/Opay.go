package Business

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/app/TwitchBot"
	"main/app/model"
	"net/http"
	"strings"
	"time"
)

var donateCache map[string]bool

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
	u := fmt.Sprintf("https://payment.opay.tw/Broadcaster/CheckDonate/%s", "7BF5D2184771810862F9070719909401")
	req, err := http.NewRequest("POST", u, strings.NewReader("__RequestVerificationToken="+model.BotSetting.OpayToken))
	if err != nil {
		return
	}

	cookie := &http.Cookie{
		Name:   "__RequestVerificationToken_Lw__",
		Value:  model.BotSetting.OpayCookie,
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
		fmt.Println("Opay http response code : ", resp.StatusCode)
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		fmt.Println("Opay resp not json")
		return
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll err:", err)
		return
	}
	oResp := opayResp{}
	err = json.Unmarshal(bodyBytes, &oResp)
	if err != nil {
		log.Println("Unmarshal err:", err)
		return
	}
	if len(oResp.LstDonate) == 0 {
		return
	}

	for _, v := range oResp.LstDonate {
		// 不在快取中 則加入快取
		if _, ok := donateCache[v.DonateID]; !ok {
			donateCache[v.DonateID] = true
			msg := fmt.Sprintf("/me 感謝 %s 贊助了 %d 元, %s", v.Name, v.Amount, v.MSG)

			fmt.Printf("[%s]  %s 贊助了 %d 元: %s\n", time.Now().In(time.FixedZone("", +8*3600)).Format("2006-01-02 15:04:05"), v.Name, v.Amount, v.MSG)

			TwitchBot.SendMessage(msg)
		}
	}
}
