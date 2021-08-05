package datamesh

// LinkUpHandler is invoked whenever a new link is brought up
type LinkUpHandler func(string)

// LinkDownHandler is invoked whenever a link goes down
type LinkDownHandler func(string)

type Handlers struct {
	linkUpHandlers   []LinkUpHandler
	linkDownHandlers []LinkDownHandler
}

func (handlers *Handlers) AddLinkUpHandler(h LinkUpHandler) {
	handlers.linkUpHandlers = append(handlers.linkUpHandlers, h)
}

func (handlers *Handlers) AddLinkDownHandler(h LinkDownHandler) {
	handlers.linkDownHandlers = append(handlers.linkDownHandlers, h)
}
