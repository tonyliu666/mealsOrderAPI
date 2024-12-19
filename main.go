package main

import (
	"log"
	router "weather/router"

	"github.com/beego/beego/orm"
	_ "github.com/lib/pq"
)

func init() {
	err := orm.RunSyncdb("default", false, true)
	if err != nil {
		log.Fatalf("Error running syncdb: %v", err)
	}
}

func main() {
	// build aa simple hello world gin server
	router.Init().Run(":8080")
}
