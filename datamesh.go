package datamesh

import (
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Datamesh struct {
	cf        *Config
	self      *identity.TokenId
	listeners map[string]*Listener
	dialers   map[string]*Dialer
	graph     *Graph
	Fwd       *Forwarder
	Handlers  *Handlers
}

func NewDatamesh(cf *Config) *Datamesh {
	d := &Datamesh{
		cf:        cf,
		listeners: make(map[string]*Listener),
		dialers:   make(map[string]*Dialer),
		graph:     newGraph(),
		Fwd:       newForwarder(),
		Handlers:  &Handlers{},
	}
	for _, listenerCf := range cf.Listeners {
		d.listeners[listenerCf.Id] = NewListener(listenerCf, &identity.TokenId{Token: listenerCf.Id})
		logrus.Infof("added listener at [%s]", listenerCf.BindAddress)
	}
	for _, dialerCf := range cf.Dialers {
		d.dialers[dialerCf.Id] = NewDialer(dialerCf, &identity.TokenId{Token: dialerCf.Id})
		logrus.Infof("added dialer at [%s]", dialerCf.BindAddress)
	}
	return d
}

func (self *Datamesh) Start() {
	for _, v := range self.listeners {
		go v.Listen(self, self.graph.incoming)
	}
	self.graph.start()
}

func (self *Datamesh) Dial(id string, endpoint transport.Address) (Link, error) {
	if dialer, found := self.dialers[id]; found {
		l, err := dialer.Dial(self, endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "error dialing [%s]", endpoint)
		}
		self.graph.addLink(l)
		return l, nil
	} else {
		return nil, errors.Errorf("no dialer [%s]", id)
	}
}
