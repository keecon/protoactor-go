package router

import (
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

type broadcastGroupRouter struct {
	GroupRouter
}

type broadcastPoolRouter struct {
	PoolRouter
}

type broadcastRouterState struct {
	routees *actor.PIDSet
	sender  actor.SenderContext
	mu      sync.RWMutex
}

func (state *broadcastRouterState) SetSender(sender actor.SenderContext) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.sender = sender
}

func (state *broadcastRouterState) SetRoutees(routees *actor.PIDSet) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.routees = routees.Clone()
}

func (state *broadcastRouterState) GetRoutees() *actor.PIDSet {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return state.routees.Clone()
}

func (state *broadcastRouterState) RouteMessage(message interface{}) {
	state.mu.RLock()
	routees := state.routees
	sender := state.sender
	state.mu.RUnlock()

	routees.ForEach(func(_ int, pid *actor.PID) {
		sender.Send(pid, message)
	})
}

func NewBroadcastPool(size int, opts ...actor.PropsOption) *actor.Props {
	return (&actor.Props{}).
		Configure(actor.WithSpawnFunc(spawner(&broadcastPoolRouter{PoolRouter{PoolSize: size}}))).
		Configure(opts...)
}

func NewBroadcastGroup(routees ...*actor.PID) *actor.Props {
	return (&actor.Props{}).Configure(actor.WithSpawnFunc(spawner(&broadcastGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}})))
}

func (config *broadcastPoolRouter) CreateRouterState() State {
	return &broadcastRouterState{}
}

func (config *broadcastGroupRouter) CreateRouterState() State {
	return &broadcastRouterState{}
}
