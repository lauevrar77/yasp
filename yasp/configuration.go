package yasp

import (
	"fmt"
	"github.com/lauevrar77/go-theater"
	"github.com/lauevrar77/yasp/yasp/domain"
	"github.com/lauevrar77/yasp/yasp/infrastructure/services"
	"gorm.io/gorm"
)

type WebsiteConfiguration struct {
	Name          string
	URL           string
	CheckInterval int
}

type Configuration struct {
	DefaultMailboxSize uint
	Websites           []WebsiteConfiguration
	NotifierURL        string
}

func Configure(configuration Configuration, actorSystem *theater.ActorSystem, db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&domain.Website{}, &domain.WebsiteState{})

	requester := services.NewRequester()
	if configuration.DefaultMailboxSize == 0 {
		configuration.DefaultMailboxSize = 100
	}
	actorSystem.Spawn("yasp.requester", &requester, int(configuration.DefaultMailboxSize))

	notifier := services.NewNotifier(configuration.NotifierURL)
	actorSystem.Spawn("yasp.notifier", &notifier, int(configuration.DefaultMailboxSize))

	for _, website := range configuration.Websites {
		yago := domain.NewWebsite(website.URL, uint(website.CheckInterval), db)
		actorSystem.Spawn(theater.ActorRef(
			fmt.Sprintf("yasp.%s", website.Name),
		), &yago, int(configuration.DefaultMailboxSize))
	}
}
