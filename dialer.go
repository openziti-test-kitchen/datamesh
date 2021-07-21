package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
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

	options := channel2.DefaultOptions()
	options.BindHandlers = []channel2.BindHandler{&linkBindHandler{}}
	dialer := channel2.NewClassicDialer(self.id, endpoint, nil)
	ch, err := channel2.NewChannel("link", dialer, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating channel")
	}

	l := newLink(self.cfg.LinkConfig, &identity.TokenId{Token: ch.ConnectionId()}, nil, ch, OutboundLink)
	return l, nil
}
