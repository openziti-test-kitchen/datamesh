package datamesh

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Overlay struct {
	incoming  chan *link
	links     map[LinkId]*link
	addLinkCb func(*link)
	lock      sync.Mutex
}

func newGraph() *Overlay {
	return &Overlay{
		incoming: make(chan *link, 128),
		links:    make(map[LinkId]*link),
	}
}

func (self *Overlay) start() {
	go self.linkAccepter()
}

func (self *Overlay) addLink(l *link) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if err := l.Start(); err == nil {
		self.links[l.LinkId()] = l
		logrus.Infof("added link [link/%s]", l.LinkId())
	} else {
		logrus.Errorf("error starting [link/%s] (%v)", l.LinkId(), err)
	}
}

func (self *Overlay) linkAccepter() {
	logrus.Info("started")
	defer logrus.Warn("exited")

	for {
		select {
		case l := <-self.incoming:
			self.addLink(l)
		}
	}
}