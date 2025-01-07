package main

import (
	"log"
	"weather/config"
	"weather/models/cache"
	"weather/models/gateway"
	router "weather/router"

	"github.com/beego/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		log.Fatalf("Error running syncdb: %v", err)
	}
	err = config.Init()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	err = gateway.Init()
	if err != nil {
		log.Fatalf("Error initializing gateway: %v", err)
	}
	err = cache.Init()
	if err != nil {
		log.Fatalf("Error initializing redis cache: %v", err)
	}
}

func main() {
	// build aa simple hello world gin server
	router.Init().Run(":8080")
}
