package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/util/sequence"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
)

func (self *link) pinger() {
	log := pfxlog.ContextLogger(self.Id().Token)
	log.Info("started")
	defer log.Warn("exited")

	seq := sequence.NewSequence()
	pending := make(map[string]*channel2.Message)

	for {
		select {
		case msg := <-self.pingResponses:
			self.pingerReceive(msg, pending)

		case <-time.After(self.cfg.PingPeriod):
			if err := self.pingerSend(seq, pending); err != nil {
				log.Errorf("error sending ping request (%v)", err)
			}
		}
	}
}

func (self *link) pingerReceive(msg *channel2.Message, pending map[string]*channel2.Message) {
	var found bool
	var stamp int64
	stamp, found = msg.GetInt64Header(PingTimestampHeaderKey)
	if !found {
		logrus.Errorf("ping response missing timestamp")
		return
	}

	var pingId string
	pingId, found = msg.GetStringHeader(PingIdHeaderKey)
	if !found {
		logrus.Errorf("ping response missing ping identifier")
		return
	}

	delta := time.Since(time.Unix(0, stamp))
	logrus.Infof("ping response [ping/%s] in [%s ms]", pingId, delta/time.Millisecond)
	delete(pending, pingId)
}

func (self *link) pingerSend(seq *sequence.Sequence, pending map[string]*channel2.Message) error {
	pingId, err := seq.NextHash()
	if err != nil {
		return errors.Wrap(err, "generating ping sequence")
	}
	headers := newHeaders()
	headers.PutBytes(PingIdHeaderKey, []byte(pingId))
	headers.PutInt64(PingTimestampHeaderKey, time.Now().UnixNano())
	ctrl := NewControl(uint32(PingRequestControlFlag), headers)
	err = self.SendControl(ctrl)
	if err != nil {
		return errors.Wrapf(err, "error sending [ping/%s]", pingId)
	}
	logrus.Infof("sent [ping/%s]", pingId)
	return nil
}
