package datamesh

// LinkUpHandler is invoked whenever a new link is brought up
type LinkUpHandler func(Link)

// LinkDownHandler is invoked whenever a link goes down
type LinkDownHandler func(Link)

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
