package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	cfg      *ListenerConfig
	id       *identity.TokenId
	listener channel2.UnderlayListener
}

func NewListener(cfg *ListenerConfig, id *identity.TokenId) *Listener {
	return &Listener{cfg: cfg, id: id}
}

func (self *Listener) Listen(incoming chan<- *link) {
	self.listener = channel2.NewClassicListener(self.id, self.cfg.BindAddress, channel2.DefaultConnectOptions(), nil)
	if err := self.listener.Listen(); err != nil {
		logrus.Errorf("error starting listener [%s] (%v)", self.cfg.BindAddress, err)
		return
	}
	pfxlog.ContextLogger(self.id.Token).Infof("started")

	options := channel2.DefaultOptions()
	options.BindHandlers = []channel2.BindHandler{&linkBindHandler{}}
	for {
		ch, err := channel2.NewChannel("link", self.listener, options)
		if err != nil {
			logrus.Errorf("error accepting new link for [%s] (%v)", self.cfg.BindAddress, err)
		}

		l := newLink(self.cfg.LinkConfig, &identity.TokenId{Token: ch.ConnectionId()}, nil, ch, InboundLink)
		incoming <- l
	}
}
