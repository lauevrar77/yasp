package services

import (
	"net/http"
	"time"

	"github.com/lauevrar77/go-theater"
	"github.com/lauevrar77/yasp/yasp/domain"
)

type Requester struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem
}

func NewRequester() Requester {
	return Requester{}
}

func (r *Requester) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	dispatcher.RegisterDefaultHandler(r.performRequest)
	r.me = me
	r.dispatcher = dispatcher
	r.system = system
}

func (r *Requester) Run() {
	for {
		r.dispatcher.Receive()
	}
}

func (r *Requester) performRequest(message theater.Message) {
	msg := message.Content.(domain.CheckWebsite)
	resp, err := http.Get(msg.URL)
	check := domain.UpdateWebsiteState{
		Error:      err,
		StatusCode: resp.StatusCode,
		CheckedAt:  time.Now(),
	}

	r.dispatcher.Send(*msg.RespondTo, theater.Message{Type: "UpdateWebsiteState", Content: check})
}
