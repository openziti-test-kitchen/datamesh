package datamesh

import (
	"github.com/michaelquigley/pfxlog"
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
)

type Dialer struct {
	id       *identity.TokenId
	bind     transport.Address
}

func NewDialer(id *identity.TokenId, bind transport.Address) *Dialer {
	return &Dialer{id, bind}
}

func (self *Dialer) Dial(endpoint transport.Address) (channel2.Channel, error) {
	pfxlog.ContextLogger(endpoint.String()).Infof("dialing")

	options := channel2.DefaultOptions()
	options.BindHandlers = []channel2.BindHandler{&linkBindHandler{}}
	dialer := channel2.NewClassicDialer(self.id, endpoint, nil)
	ch, err := channel2.NewChannel("link", dialer, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating channel")
	}
	return ch, nil
}
