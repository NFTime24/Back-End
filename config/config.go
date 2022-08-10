package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf("config/config.json", &configuration)
	return configuration
}

// {
// 	"DB_USERNAME": "root",
// 	"DB_PASSWORD": "dukecoon!",
// 	"DB_PORT": "3306",
// 	"DB_HOST": "database-1.cng84nshh1hx.us-west-2.rds.amazonaws.com",
// 	"DB_NAME": "test"
//   }
// {
//   "DB_USERNAME": "root",
//   "DB_PASSWORD": "111111",
//   "DB_PORT": "3306",
//   "DB_HOST": "localhost",
//   "DB_NAME": "test"
// }
