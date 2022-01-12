package types

import (
	"fmt"
)

// The router is a map from module name to the XIBCModule
// which contains all the module-defined callbacks
type Router struct {
	routes map[string]XIBCModule
	sealed bool
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]XIBCModule),
	}
}

// Seal prevents the Router from any subsequent route handlers to be registered.
// Seal will panic if called more than once.
func (rtr *Router) Seal() {
	if rtr.sealed {
		panic("router already sealed")
	}
	rtr.sealed = true
}

// Sealed returns a boolean signifying if the Router is sealed or not.
func (rtr Router) Sealed() bool {
	return rtr.sealed
}

// AddRoute adds XIBCModule for a given port. It returns the Router
// so AddRoute calls can be linked. It will panic if the Router is sealed.
func (rtr *Router) AddRoute(port string, cbs XIBCModule) *Router {
	if rtr.sealed {
		panic(fmt.Sprintf("router sealed; cannot register %s route callbacks", port))
	}
	if rtr.HasRoute(port) {
		panic(fmt.Sprintf("route %s has already been registered", port))
	}

	rtr.routes[port] = cbs
	return rtr
}

// HasRoute returns true if the Router has a module registered or false otherwise.
func (rtr *Router) HasRoute(port string) bool {
	_, ok := rtr.routes[port]
	return ok
}

// GetRoute returns a XIBCModule for a given port.
func (rtr *Router) GetRoute(port string) (XIBCModule, bool) {
	if !rtr.HasRoute(port) {
		return nil, false
	}
	return rtr.routes[port], true
}
