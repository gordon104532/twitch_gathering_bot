package model

var BotSetting botSettingModel

type botSettingModel struct {
	TargetTwitchID string `json:"targetTwitchID"`
	ChatTwitchID   string `json:"chatTwitchID"`
	CheckDonate    bool   `json:"checkDonate"`
	OpayID         string `json:"opayID"`
	OpayToken      string `json:"opayToken"`
	OpayCookie     string `json:"opayCookie"`
	TwitchOAth     string `json:"twitchOAth"`
	AutoHello      bool   `json:"autoHello"`
	AutoHelloMsg   string `json:"autoHelloMsg"`
	AutoHelloEmoji string `json:"autoHelloEmoji"`
}
