package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/sirupsen/logrus"
	"time"
)

func (self *link) pinger() {
	log := pfxlog.ContextLogger(self.Id().Token)
	log.Info("started")
	defer log.Warn("exited")

	for {
		select {
		case <-self.pingResponses:
			logrus.Infof("ping response received")

		case <-time.After(self.cfg.PingPeriod):
			if err := self.SendControl(NewControl(uint32(PingRequestControlFlag), nil)); err == nil {
				log.Info("sent control ping request")
			} else {
				log.Errorf("error sending control ping request (%v)", err)
			}
		}
	}
}
