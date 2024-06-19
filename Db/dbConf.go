package Db

type DbConfig struct {
	DatabaseURL string `toml:"database_url" `
}

func NewDbConf() *DbConfig {
	return &DbConfig{}
}
