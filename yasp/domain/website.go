package domain

import (
	"fmt"
	"log"
	"time"

	"github.com/lauevrar77/go-theater"
	"gorm.io/gorm"
)

type WebsiteState struct {
	gorm.Model

	WebsiteID  uint
	CheckError *string
	StatusCode int
	CheckAt    time.Time
}

func (WebsiteState) TableName() string {
	return "yasp_website_states"
}

type Website struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	URL           string `gorm:"unique"`
	States        []WebsiteState
	CheckInterval uint

	lastCheckAt *time.Time    `gorm:"-:all"`
	lastState   *WebsiteState `gorm:"-:all"`

	me       *theater.ActorRef    `gorm:"-:all"`
	mailbox  *theater.Mailbox     `gorm:"-:all"`
	system   *theater.ActorSystem `gorm:"-:all"`
	notifier *theater.Mailbox     `gorm:"-:all"`

	db *gorm.DB `gorm:"-:all"`
}

func (Website) TableName() string {
	return "yasp_websites"
}

func NewWebsite(url string, checkInterval uint, db *gorm.DB) Website {
	website := Website{}
	db.FirstOrCreate(&website, Website{URL: url})
	website.db = db
	website.CheckInterval = checkInterval
	db.Save(&website)
	return website
}

func (w *Website) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
	w.me = me
	w.mailbox = mailbox
	w.system = system

	notifier, err := system.ByRef("notifier")
	if err != nil {
		log.Println("Failed to find notifier")
	}
	w.notifier = notifier
}

func (w *Website) Run() {
	for {
		select {
		case msg := <-*w.mailbox:
			switch msg.Type {
			case "UpdateWebsiteState":
				msg := msg.Content.(UpdateWebsiteState)
				w.updateState(msg)
			default:
				log.Println("Unknown message type")
			}
		default:
			if w.lastCheckAt == nil || time.Since(*w.lastCheckAt) > time.Duration(w.CheckInterval)*time.Second {
				requester, err := w.system.ByRef("requester")
				if err != nil {
					log.Println("Failed to find requester")
				}
				log.Printf("Requesting check for website %s\n", w.URL)
				*requester <- theater.Message{Type: "CheckWebsite", Content: CheckWebsite{URL: w.URL, RespondTo: w.me}}
				now := time.Now()
				w.lastCheckAt = &now
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}
}

func (w *Website) updateState(msg UpdateWebsiteState) {
	log.Printf("Updating state for %s\n", w.URL)
	var checkError *string
	if msg.Error != nil {
		errorString := msg.Error.Error()
		checkError = &errorString
	}

	state := WebsiteState{
		WebsiteID:  w.ID,
		CheckError: checkError,
		StatusCode: msg.StatusCode,
		CheckAt:    msg.CheckedAt,
	}
	w.db.Create(&state)
	w.lastState = &state
	if w.shouldNotify(state) {
		w.notify(state)
	}
}

func (w *Website) shouldNotify(state WebsiteState) bool {
	if w.lastState == nil {
		return false
	}
	if state.CheckError != w.lastState.CheckError {
		return true
	}
	if state.StatusCode != w.lastState.StatusCode {
		return true
	}
	return false
}

func (w *Website) notify(state WebsiteState) {
	message := fmt.Sprintf("Website %s is %d", w.URL, state.StatusCode)
	if state.CheckError != nil {
		message = fmt.Sprintf("%s, error: %s", message, *state.CheckError)
	}
	*w.notifier <- theater.Message{Type: "Notify", Content: message}
}

type UpdateWebsiteState struct {
	Error      error
	StatusCode int
	CheckedAt  time.Time
}

type CheckWebsite struct {
	URL       string
	RespondTo *theater.ActorRef
}
