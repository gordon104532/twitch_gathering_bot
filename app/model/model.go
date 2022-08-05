package model

var GatheringLevel int
var BotSetting BotSettingModel

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
	Twitch struct {
		ChatTwitchID   string `json:"chatTwitchID"`
		TwitchOAth     string `json:"twitchOAth"`
		AutoHello      bool   `json:"autoHello"`
		AutoHelloMsg   string `json:"autoHelloMsg"`
		AutoHelloEmoji string `json:"autoHelloEmoji"`
	} `json:"twitch"`
	GatheringEvent struct {
		GatheringSwitch bool `json:"gatheringSwitch"`
		SubPoint        int  `json:"subPoint"`
		CheerPoint      int  `json:"cheerPoint"`
		OpayPoint       int  `json:"opayPoint"`
		LevelOne        int  `json:"levelOne"`
		LevelTwo        int  `json:"levelTwo"`
		LevelThree      int  `json:"levelThree"`
		LevelFour       int  `json:"levelFour"`
		LevelFive       int  `json:"levelFive"`
		InitPoint       int  `json:"initPoint"`
	} `json:"gatheringEvent"`
}
