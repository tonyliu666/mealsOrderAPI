package main

import (
	"log"
	"weather/config"
	router "weather/router"

	"github.com/beego/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	err := orm.RunSyncdb("default", true, true)
	if err != nil {
		log.Fatalf("Error running syncdb: %v", err)
	}
	err = config.Init()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {
	// build aa simple hello world gin server
	router.Init().Run(":8080")
}
