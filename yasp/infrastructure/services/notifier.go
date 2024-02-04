package services

import (
	"net/http"
	"strings"

	"github.com/lauevrar77/go-theater"
)

type Notifier struct {
	me         theater.ActorRef
	dispatcher theater.MessageDispatcher
	system     *theater.ActorSystem

	NotificationURL string
}

func NewNotifier(notificationURL string) Notifier {
	return Notifier{NotificationURL: notificationURL}
}

func (n *Notifier) Initialize(me theater.ActorRef, dispatcher theater.MessageDispatcher, system *theater.ActorSystem) {
	n.dispatcher.RegisterDefaultHandler(n.notify)

	n.me = me
	n.dispatcher = dispatcher
	n.system = system
}

func (n *Notifier) Run() {
	for {
		n.dispatcher.Receive()
	}
}

func (n *Notifier) notify(msg theater.Message) {
	message := msg.Content.(string)
	http.Post(n.NotificationURL, "text/plain", strings.NewReader(message))
}
