package datamesh

import (
	"github.com/pkg/errors"
	"sync"
)

type Forwarder struct {
	lock   sync.RWMutex
	dests  map[Address]Destination
	routes map[Circuit]map[Address]Address
}

func newForwarder() *Forwarder {
	return &Forwarder{
		dests:  make(map[Address]Destination),
		routes: make(map[Circuit]map[Address]Address),
	}
}

func (fw *Forwarder) AddDestination(d Destination) {
	fw.lock.Lock()
	fw.dests[d.Address()] = d
	fw.lock.Unlock()
}

func (fw *Forwarder) RemoveDestination(d Destination) {
	fw.lock.Lock()
	delete(fw.dests, d.Address())
	fw.lock.Unlock()
}

func (fw *Forwarder) AddRoute(circuitId Circuit, srcAddr, destAddr Address) {
	fw.lock.Lock()
	routeMap, found := fw.routes[circuitId]
	if !found {
		routeMap = make(map[Address]Address)
	}
	routeMap[srcAddr] = destAddr
	fw.routes[circuitId] = routeMap
	fw.lock.Unlock()
}

func (fw *Forwarder) RemoveRoute(circuitId Circuit, srcAddr Address) {
	fw.lock.Lock()
	routeMap, found := fw.routes[circuitId]
	if found {
		delete(routeMap, srcAddr)
		if len(routeMap) > 0 {
			fw.routes[circuitId] = routeMap
		} else {
			delete(fw.routes, circuitId)
		}
	}
	fw.lock.Unlock()
}

func (fw *Forwarder) Forward(srcAddr Address, payload *Payload) error {
	fw.lock.RLock()
	defer fw.lock.RUnlock()

	if destination := fw.destination(payload.CircuitId, srcAddr); destination != nil {
		if err := destination.FromNetwork(payload); err != nil {
			return errors.Wrapf(err, "unable to forward [circuit/%s][src/%s]", payload.CircuitId, srcAddr)
		}
		return nil
	} else {
		return errors.Errorf("no destination for [circuit/%s][src/%s]", payload.CircuitId, srcAddr)
	}
}

func (fw *Forwarder) destination(circuitId Circuit, srcAddr Address) Destination {
	routeMap, found := fw.routes[circuitId]
	if found {
		destAddr, found := routeMap[srcAddr]
		if found {
			destination, found := fw.dests[destAddr]
			if found {
				return destination
			}
		}
	}
	return nil
}
