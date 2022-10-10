package model

var GatheringLevel int
var BotSetting BotSettingModel
var DetailSetting DetailSettingModel

type BotSettingModel struct {
	General struct {
		TargetTwitchID string `json:"targetTwitchID"`
	} `json:"general"`
	Opay struct {
		CheckDonate bool   `json:"checkDonate"`
		OpayID      string `json:"opayID"`
		OpayToken   string `json:"opayToken"`
		OpayCookie  string `json:"opayCookie"`
		OpayMsg     string `json:"opayMsg"`
	} `json:"opay"`
	Ecpay struct {
		CheckDonate bool   `json:"checkDonate"`
		EcpayID     string `json:"ecpayID"`
		EcpayMsg    string `json:"ecpayMsg"`
	} `json:"ecpay"`
	Twitch struct {
		ChatTwitchID string `json:"chatTwitchID"`
		TwitchOAth   string `json:"twitchOAth"`
		// AutoHello      bool   `json:"autoHello"`
		// AutoHelloMsg   string `json:"autoHelloMsg"`
		// AutoHelloEmoji string `json:"autoHelloEmoji"`
	} `json:"twitch"`
	GatheringEvent struct {
		GatheringSwitch bool   `json:"gatheringSwitch"`
		GatheringTitle  string `json:"gatheringTitle"`
		QueryCommand    string `json:"queryCommand"`
		SubPoint        int    `json:"subPoint"`
		CheerPoint      int    `json:"cheerPoint"`
		StampPoint      int    `json:"stampPoint"`
		OpayPoint       int    `json:"opayPoint"`
		EcpayPoint      int    `json:"ecpayPoint"`
		LevelSetting    struct {
			Lv1  int `json:"lv1"`
			Lv2  int `json:"lv2"`
			Lv3  int `json:"lv3"`
			Lv4  int `json:"lv4"`
			Lv5  int `json:"lv5"`
			Lv6  int `json:"lv6"`
			Lv7  int `json:"lv7"`
			Lv8  int `json:"lv8"`
			Lv9  int `json:"lv9"`
			Lv10 int `json:"lv10"`
		} `json:"levelSetting"`
		InitPoint int `json:"initPoint"`
	} `json:"gatheringEvent"`
}

type DetailSettingModel struct {
	CheckEmoji  string `json:"checkEmoji"`
	ProgressBar struct {
		TitleColor   string `json:"titleColor"`
		BarCollor    string `json:"barCollor"`
		BarTxtCollor string `json:"barTxtCollor"`
	} `json:"progressBar"`
	Subgift struct {
		One    int `json:"one"`
		Three  int `json:"three"`
		Six    int `json:"six"`
		Twelve int `json:"twelve"`
	} `json:"subgift"`
	Resub struct {
		Zero   int `json:"zero"`
		One    int `json:"one"`
		Three  int `json:"three"`
		Six    int `json:"six"`
		Twelve int `json:"twelve"`
	} `json:"resub"`
	Sub struct {
		One    int `json:"one"`
		Three  int `json:"three"`
		Six    int `json:"six"`
		Twelve int `json:"twelve"`
	} `json:"sub"`
	Tier struct {
		One   int `json:"one"`
		Two   int `json:"two"`
		Three int `json:"three"`
	} `json:"tier"`
}
