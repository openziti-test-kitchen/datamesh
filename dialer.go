package datamesh

import (
	"github.com/openziti/foundation/channel2"
	"github.com/openziti/foundation/identity/identity"
	"github.com/openziti/foundation/transport"
	"github.com/pkg/errors"
)

type Dialer struct {
	id       *identity.TokenId
	endpoint transport.Address
	bind     transport.Address
}

func NewDialer(id *identity.TokenId, endpoint, bind transport.Address) *Dialer {
	return &Dialer{id, endpoint, bind}
}

func (self *Dialer) Dial() (channel2.Channel, error) {
	options := channel2.DefaultOptions()
	options.BindHandlers = []channel2.BindHandler{&bindHandler{}}
	dialer := channel2.NewClassicDialer(self.id, self.endpoint, nil)
	ch, err := channel2.NewChannel("link", dialer, options)
	if err != nil {
		return nil, errors.Wrap(err, "error creating channel")
	}
	return ch, nil
}
