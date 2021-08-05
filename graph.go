package datamesh

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Graph struct {
	incoming  chan *link
	links     map[string]*link
	addLinkCb func(*link)
	lock      sync.Mutex
}

func newGraph() *Graph {
	return &Graph{
		incoming: make(chan *link, 128),
		links:    make(map[string]*link),
	}
}

func (self *Graph) start() {
	go self.linkAccepter()
}

func (self *Graph) addLink(l *link) {
	self.lock.Lock()
	defer self.lock.Unlock()

	if err := l.Start(); err == nil {
		self.links[l.Id().Token] = l
		logrus.Infof("added link [link/%s]", l.Id().Token)
	} else {
		logrus.Errorf("error starting [link/%s] (%v)", l.Id().Token, err)
	}
}

func (self *Graph) linkAccepter() {
	logrus.Info("started")
	defer logrus.Warn("exited")

	for {
		select {
		case l := <-self.incoming:
			self.addLink(l)
		}
	}
}