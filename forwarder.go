package datamesh

import (
	"github.com/pkg/errors"
	"sync"
)

type Forwarder struct {
	lock   sync.RWMutex
	dests  map[Address]Destination
	routes map[CircuitId]map[Address]Address
}

func newForwarder() *Forwarder {
	return &Forwarder{
		dests:  make(map[Address]Destination),
		routes: make(map[CircuitId]map[Address]Address),
	}
}

func (fw *Forwarder) addDestination(d Destination) {
	fw.lock.Lock()
	fw.dests[d.Address()] = d
	fw.lock.Unlock()
}

func (fw *Forwarder) removeDestination(d Destination) {
	fw.lock.Lock()
	delete(fw.dests, d.Address())
	fw.lock.Unlock()
}

func (fw *Forwarder) addRoute(circuitId CircuitId, srcAddr, destAddr Address) {
	fw.lock.Lock()
	routeMap, found := fw.routes[circuitId]
	if !found {
		routeMap = make(map[Address]Address)
	}
	routeMap[srcAddr] = destAddr
	fw.routes[circuitId] = routeMap
	fw.lock.Unlock()
}

func (fw *Forwarder) removeRoute(circuitId CircuitId, srcAddr, destAddr Address) {
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

func (fw *Forwarder) ForwardPayload(circuitId CircuitId, srcAddr Address, p *Payload) error {
	fw.lock.RLock()
	defer fw.lock.RUnlock()

	if destination := fw.destination(circuitId, srcAddr); destination != nil {
		if err := destination.SendPayload(p); err != nil {
			return errors.Wrapf(err, "unable to forward [circuit/%s][src/%s]", circuitId, srcAddr)
		}
		return nil
	} else {
		return errors.Errorf("no destination for [circuit/%s][src/%s]", circuitId, srcAddr)
	}
}

func (fw *Forwarder) ForwardAcknowledgement(circuitId CircuitId, srcAddr Address, a *Acknowledgement) error {
	fw.lock.RLock()
	defer fw.lock.RUnlock()

	if destination := fw.destination(circuitId, srcAddr); destination != nil {
		if err := destination.SendAcknowledgement(a); err != nil {
			return errors.Wrapf(err, "unable to forward [circuit/%s][src/%s]", circuitId, srcAddr)
		}
		return nil
	} else {
		return errors.Errorf("no destination for [circuit/%s][src/%s]", circuitId, srcAddr)
	}
}

func (fw *Forwarder) destination(circuitId CircuitId, srcAddr Address) Destination {
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
