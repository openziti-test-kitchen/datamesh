package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/sirupsen/logrus"
)

type Datamesh struct {
	cf        *Config
	listeners map[string]*Listener
	dialers   map[string]*Dialer
	incoming  chan channel2.Channel
}

func NewDatamesh(cf *Config) *Datamesh {
	d := &Datamesh{
		cf:        cf,
		listeners: make(map[string]*Listener),
		dialers:   make(map[string]*Dialer),
		incoming:  make(chan channel2.Channel, 128),
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
}
