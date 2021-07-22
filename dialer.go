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

	options := channel.DefaultOptions()
	options.BindHandlers = []channel.BindHandler{&linkBindHandler{}}
	dialer := channel.NewClassicDialer(self.id, endpoint, nil)
	ch, err := channel.NewChannel("link", dialer, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating channel")
	}

	l := newLink(self.cfg.LinkConfig, &identity.TokenId{Token: ch.ConnectionId()}, nil, ch, OutboundLink)
	return l, nil
}
