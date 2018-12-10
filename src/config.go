package src

type Config struct {
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
