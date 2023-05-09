package main

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func InitTestConfig() (dbConf DBConfig) {
	dbConf = DBConfig{}
	dbConf.Host = "127.0.0.1"
	dbConf.Port = "3306"
	dbConf.Username = "root"
	dbConf.Password = "123456"
	dbConf.Database = "ksxt_dev"
	return
}
