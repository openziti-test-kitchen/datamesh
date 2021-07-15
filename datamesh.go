package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type Datamesh struct {
	cf        *Config
	self      *identity.TokenId
	listeners map[string]*Listener
	dialers   map[string]*Dialer
	incoming  chan channel2.Channel
	links     map[string]*link
	lock      sync.Mutex
}

func NewDatamesh(cf *Config) *Datamesh {
	d := &Datamesh{
		cf:        cf,
		listeners: make(map[string]*Listener),
		dialers:   make(map[string]*Dialer),
		incoming:  make(chan channel2.Channel, 128),
		links:     make(map[string]*link),
	}
	for _, listenerCf := range cf.Listeners {
		d.listeners[listenerCf.Id] = NewListener(&identity.TokenId{Token: listenerCf.Id}, listenerCf.BindAddress)
		logrus.Infof("added listener at [%s]", listenerCf.BindAddress)
	}
	for _, dialerCf := range cf.Dialers {
		d.dialers[dialerCf.Id] = NewDialer(&identity.TokenId{Token: dialerCf.Id}, dialerCf.BindAddress)
		logrus.Infof("added dialer at [%s]", dialerCf.BindAddress)
	}
	return d
}

func (self *Datamesh) Start() {
	for _, v := range self.listeners {
		go v.Listen(self.incoming)
	}
	go self.accepter()
}

func (self *Datamesh) Dial(id string, endpoint transport.Address) (Link, error) {
	if dialer, found := self.dialers[id]; found {
		ch, err := dialer.Dial(endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "error dialing [%s]", endpoint)
		}
		l := &link{ch: ch, id: &identity.TokenId{Token: ch.ConnectionId()}}
		self.lock.Lock()
		self.links[ch.Id().Token] = l
		self.lock.Unlock()

		return l, nil

	} else {
		return nil, errors.Errorf("no dialer [%s]", id)
	}
}

func (self *Datamesh) accepter() {
	for {
		select {
		case ch := <-self.incoming:
			l := &link{ch: ch, id: &identity.TokenId{Token: ch.ConnectionId()}}
			self.lock.Lock()
			self.links[ch.Id().Token] = l
			self.lock.Unlock()

			logrus.Infof("accepted link [%s]", l.Id().Token)
		}
	}
}
