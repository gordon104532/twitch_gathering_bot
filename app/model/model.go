package model

var BotSetting botSettingModel

type botSettingModel struct {
	TargetTwitchID string `json:"targetTwitchID"`
	OpayID         string `json:"opayID"`
	OpayToken      string `json:"opayToken"`
	OpayCookie     string `json:"opayCookie"`
	ChatTwitchID   string `json:"chatTwitchID"`
	TwitchOAth     string `json:"twitchOAth"`
}
