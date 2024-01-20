package services

import (
	"log"
	"net/http"
	"time"

	"github.com/lauevrar77/go-theater"
	"github.com/lauevrar77/yasp/yasp/domain"
)

type Requester struct {
	me      *theater.ActorRef
	mailbox *theater.Mailbox
	system  *theater.ActorSystem
}

func NewRequester() Requester {
	return Requester{}
}

func (r *Requester) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
	r.me = me
	r.mailbox = mailbox
	r.system = system
}

func (r *Requester) Run() {
	for msg := range *r.mailbox {
		msg := msg.Content.(domain.CheckWebsite)
		resp, err := http.Get(msg.URL)
		check := domain.UpdateWebsiteState{
			Error:      err,
			StatusCode: resp.StatusCode,
			CheckedAt:  time.Now(),
		}

		target, err := r.system.ByRef(*msg.RespondTo)
		if err != nil {
			log.Println("Failed to find target")
		}
		*target <- theater.Message{Type: "UpdateWebsiteState", Content: check}
	}
}
