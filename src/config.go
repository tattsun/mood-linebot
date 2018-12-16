package src

type Config struct {
	APIUserID   string
	APIPassword string

	Port string
	DB   DBConfig
	LINE LINEConfig
}

type LINEConfig struct {
	ChannelSecret string
	ChannelToken  string
}

type DBConfig struct {
	MongoAddr     string
	MongoDatabase string
}
