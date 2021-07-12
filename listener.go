package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/sirupsen/logrus"
)

type Listener struct {
	id       *identity.TokenId
	bind     transport.Address
	listener channel2.UnderlayListener
}

func NewListener(id *identity.TokenId, bind transport.Address) *Listener {
	return &Listener{id: id, bind: bind}
}

func (self *Listener) Listen(incoming chan<- channel2.Channel) {
	self.listener = channel2.NewClassicListener(self.id, self.bind, channel2.DefaultConnectOptions(), nil)
	if err := self.listener.Listen(); err != nil {
		logrus.Errorf("error starting listener [%s] (%v)", self.bind, err)
		return
	}
	pfxlog.ContextLogger(self.id.Token).Infof("started")

	options := channel2.DefaultOptions()
	options.BindHandlers = []channel2.BindHandler{&bindHandler{}}
	for {
		ch, err := channel2.NewChannel("link", self.listener, options)
		if err != nil {
			logrus.Errorf("error accepting new link for [%s] (%v)", self.bind, err)
		}
		incoming <- ch
	}
}
