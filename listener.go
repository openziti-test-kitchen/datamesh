package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti-incubator/datamesh/channel"
	"github.com/openziti/foundation/identity/identity"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	datamesh *Datamesh
	cfg      *ListenerConfig
	id       *identity.TokenId
	listener channel.UnderlayListener
}

func NewListener(cfg *ListenerConfig, id *identity.TokenId) *Listener {
	return &Listener{cfg: cfg, id: id}
}

func (self *Listener) Listen(datamesh *Datamesh, incoming chan<- *link) {
	self.datamesh = datamesh
	self.listener = channel.NewClassicListener(self.id, self.cfg.BindAddress, channel.DefaultConnectOptions(), nil)
	if err := self.listener.Listen(); err != nil {
		logrus.Errorf("error starting listener [%s] (%v)", self.cfg.BindAddress, err)
		return
	}
	pfxlog.ContextLogger(self.id.Token).Infof("started")

	for {
		l := newLink(self.cfg.LinkConfig, InboundLink, self.datamesh)

		options := channel.DefaultOptions()
		options.BindHandlers = []channel.BindHandler{newLinkBindHandler(self.datamesh, l)}

		ch, err := channel.NewChannel("link", self.listener, options)
		if err != nil {
			logrus.Errorf("error accepting new link for [%s] (%v)", self.cfg.BindAddress, err)
		}
		l.setChannel(ch)

		incoming <- l
	}
}
