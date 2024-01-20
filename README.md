# YASP
YASP (Yet Another Status Page) is a status page system using [Theater](https://github.com/lauevrar77/go-theater) (an actor system implemented in Go).

It is made as a module that is easilly pluggable in your existing actor system.

Example usage : 
```golang
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
```
