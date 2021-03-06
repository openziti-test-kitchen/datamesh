package datamesh

import (
	"fmt"
	"github.com/openziti-incubator/datamesh/channel"
	"github.com/openziti/dilithium"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/openziti/foundation/util/sequence"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type Circuit string

type Datamesh struct {
	cf         *Config
	self       *identity.TokenId
	listeners  map[string]*Listener
	dialers    map[string]*Dialer
	nics       map[string]NIC
	overlay    *Overlay
	Fwd        *Forwarder
	Handlers   *Handlers
	sequence   *sequence.Sequence
	pool       *dilithium.Pool
	lock       sync.Mutex
	instrument dilithium.Instrument
}

func NewDatamesh(cf *Config) *Datamesh {
	d := &Datamesh{
		cf:        cf,
		listeners: make(map[string]*Listener),
		dialers:   make(map[string]*Dialer),
		nics:      make(map[string]NIC),
		overlay:   newGraph(),
		Fwd:       newForwarder(),
		Handlers:  &Handlers{},
		sequence:  sequence.NewSequence(),
		pool:      dilithium.NewPool("link", 128*1024, dilithium.NilInstrumentInstance{}),
	}
	channel.SetUnderlayRegistrySequence(d.sequence)
	d.overlay.addLinkCb = d.addLinkCb
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
		go v.Listen(self, self.overlay.incoming)
	}
	self.overlay.start()
}

func (self *Datamesh) DialLink(id string, endpoint transport.Address) (Link, error) {
	if dialer, found := self.dialers[id]; found {
		l, err := dialer.Dial(self, endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "error dialing link at [%s]", endpoint)
		}
		self.overlay.addLink(l)
		return l, nil
	} else {
		return nil, errors.Errorf("no dialer [%s]", id)
	}
}

func (self *Datamesh) InsertNIC(circuitId Circuit, endpoint Endpoint) (NIC, error) {
	self.lock.Lock()
	defer self.lock.Unlock()

	addr, err := self.sequence.NextHash()
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/home/plorenz/tmp/dilithium/%v", os.Getpid())
	if err := os.MkdirAll(path, 0700); err != nil {
		return nil, err
	}
	instrument, err := dilithium.NewMetricsInstrumentNoCf(path, 250, true)
	if err != nil {
		return nil, err
	}
	ii := instrument.NewInstance("test")
	nic := newNIC(self, circuitId, Address(addr), endpoint, ii)
	txAlg, err := self.cf.Profile.(dilithium.TxAlgorithmProfile).Create(ii)
	if err != nil {
		return nil, err
	}
	if err := nic.(*nicImpl).SetTxAlgorithm(txAlg); err != nil {
		return nil, err
	}
	if err := nic.(*nicImpl).Start(); err != nil {
		return nil, err
	}
	self.nics[addr] = nic

	self.Fwd.AddDestination(nic)

	return nic, nil
}

func (self *Datamesh) addLinkCb(l *link) {
	self.Fwd.AddDestination(l)
	for _, handler := range self.Handlers.linkUpHandlers {
		handler(l)
	}
}
