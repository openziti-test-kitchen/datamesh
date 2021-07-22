package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti-incubator/datamesh/channel"
	"github.com/openziti/foundation/identity/identity"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	cfg      *ListenerConfig
	id       *identity.TokenId
	listener channel.UnderlayListener
}

func NewListener(cfg *ListenerConfig, id *identity.TokenId) *Listener {
	return &Listener{cfg: cfg, id: id}
}

func (self *Listener) Listen(incoming chan<- *link) {
	self.listener = channel.NewClassicListener(self.id, self.cfg.BindAddress, channel.DefaultConnectOptions(), nil)
	if err := self.listener.Listen(); err != nil {
		logrus.Errorf("error starting listener [%s] (%v)", self.cfg.BindAddress, err)
		return
	}
	pfxlog.ContextLogger(self.id.Token).Infof("started")

	for {
		l := newLink(self.cfg.LinkConfig, InboundLink)

		options := channel.DefaultOptions()
		options.BindHandlers = []channel.BindHandler{newLinkBindHandler(l)}

		ch, err := channel.NewChannel("link", self.listener, options)
		if err != nil {
			logrus.Errorf("error accepting new link for [%s] (%v)", self.cfg.BindAddress, err)
		}
		l.setChannel(ch)

		incoming <- l
	}
}
