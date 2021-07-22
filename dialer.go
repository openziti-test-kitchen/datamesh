package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti-incubator/datamesh/channel"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
)

type Dialer struct {
	cfg *DialerConfig
	id  *identity.TokenId
}

func NewDialer(cfg *DialerConfig, id *identity.TokenId) *Dialer {
	return &Dialer{cfg, id}
}

func (self *Dialer) Dial(endpoint transport.Address) (*link, error) {
	pfxlog.ContextLogger(endpoint.String()).Infof("dialing")

	l := newLink(self.cfg.LinkConfig, OutboundLink)

	options := channel.DefaultOptions()
	options.BindHandlers = []channel.BindHandler{newLinkBindHandler(l)}
	dialer := channel.NewClassicDialer(self.id, endpoint, nil)
	ch, err := channel.NewChannel("link", dialer, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating channel")
	}
	l.setChannel(ch)

	return l, nil
}
