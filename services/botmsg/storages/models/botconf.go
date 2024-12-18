package models

type BotType int

var (
	BotType_Default BotType = 0
	BotType_Custom  BotType = 1
	BotType_Dify    BotType = 2
	BotType_Minmax  BotType = 3
)

type BotConf struct {
	ID          int64
	AppKey      string
	BotId       string
	Nickname    string
	BotPortrait string
	BotType     BotType
	BotConf     string
}

type IBotConfStorage interface {
	Upsert(item BotConf) error
	FindById(appkey, botId string) (*BotConf, error)
}
