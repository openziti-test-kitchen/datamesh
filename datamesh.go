package datamesh

import (
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
	incoming  chan *link
	links     map[string]*link
	lock      sync.Mutex
}

func NewDatamesh(cf *Config) *Datamesh {
	d := &Datamesh{
		cf:        cf,
		listeners: make(map[string]*Listener),
		dialers:   make(map[string]*Dialer),
		incoming:  make(chan *link, 128),
		links:     make(map[string]*link),
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
		go v.Listen(self.incoming)
	}
	if len(self.listeners) > 0 {
		go self.linkAccepter()
	} else {
		logrus.Warn("no listeners, not starting accepter")
	}
}

func (self *Datamesh) Dial(id string, endpoint transport.Address) (Link, error) {
	if dialer, found := self.dialers[id]; found {
		l, err := dialer.Dial(endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "error dialing [%s]", endpoint)
		}
		self.addLink(l)
		return l, nil
	} else {
		return nil, errors.Errorf("no dialer [%s]", id)
	}
}

func (self *Datamesh) addLink(l *link) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if err := l.Start(); err == nil {
		self.links[l.Id().Token] = l
		logrus.Infof("added link [link/%s]", l.Id().Token)
	} else {
		logrus.Errorf("error starting [link/%s] (%v)", l.Id().Token, err)
	}
}

func (self *Datamesh) linkAccepter() {
	logrus.Info("started")
	defer logrus.Warn("exited")

	for {
		select {
		case l := <-self.incoming:
			self.addLink(l)
		}
	}
}
