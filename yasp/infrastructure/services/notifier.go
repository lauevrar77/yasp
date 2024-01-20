package services

import (
	"net/http"
	"strings"

	"github.com/lauevrar77/go-theater"
)

type Notifier struct {
	me      *theater.ActorRef
	mailbox *theater.Mailbox
	system  *theater.ActorSystem

	NotificationURL string
}

func NewNotifier(notificationURL string) Notifier {
	return Notifier{NotificationURL: notificationURL}
}

func (n *Notifier) Initialize(me *theater.ActorRef, mailbox *theater.Mailbox, system *theater.ActorSystem) {
	n.me = me
	n.mailbox = mailbox
	n.system = system
}

func (n *Notifier) Run() {
	for msg := range *n.mailbox {
		if n.NotificationURL == "" {
			continue
		}

		msg := msg.Content.(string)
		n.notify(msg)
	}
}

func (n *Notifier) notify(message string) {
	http.Post(n.NotificationURL, "text/plain", strings.NewReader(message))
}
