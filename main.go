package main

import (
	"github.com/lauevrar77/go-theater"
	"github.com/lauevrar77/yasp/yasp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("yasp.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	configuration := yasp.Configuration{
		Websites: []yasp.WebsiteConfiguration{
			{Name: "google", URL: "https://google.com", CheckInterval: 30},
		},
		NotifierURL: "http://100.82.213.51:8011/yasp",
	}

	actorSystem := theater.NewActorSystem()
	yasp.Configure(configuration, &actorSystem, db)

	actorSystem.Run()
}
