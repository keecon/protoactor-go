package router

import (
	"math/rand"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
)

type randomGroupRouter struct {
	GroupRouter
}

type randomPoolRouter struct {
	PoolRouter
}

type randomRouterState struct {
	routees *actor.PIDSet
	sender  actor.SenderContext
	mu      sync.RWMutex
}

func (state *randomRouterState) SetSender(sender actor.SenderContext) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.sender = sender
}

func (state *randomRouterState) SetRoutees(routees *actor.PIDSet) {
	state.mu.Lock()
	defer state.mu.Unlock()

	state.routees = routees.Clone()
}

func (state *randomRouterState) GetRoutees() *actor.PIDSet {
	state.mu.RLock()
	defer state.mu.RUnlock()

	return state.routees.Clone()
}

func (state *randomRouterState) RouteMessage(message interface{}) {
	state.mu.RLock()
	routees := state.routees
	sender := state.sender
	state.mu.RUnlock()

	pid := randomRoutee(routees)
	sender.Send(pid, message)
}

func NewRandomPool(size int, opts ...actor.PropsOption) *actor.Props {
	return (&actor.Props{}).
		Configure(actor.WithSpawnFunc(spawner(&randomPoolRouter{PoolRouter{PoolSize: size}}))).
		Configure(opts...)
}

func NewRandomGroup(routees ...*actor.PID) *actor.Props {
	return (&actor.Props{}).Configure(actor.WithSpawnFunc(spawner(&randomGroupRouter{GroupRouter{Routees: actor.NewPIDSet(routees...)}})))
}

func (config *randomPoolRouter) CreateRouterState() State {
	return &randomRouterState{}
}

func (config *randomGroupRouter) CreateRouterState() State {
	return &randomRouterState{}
}

func randomRoutee(routees *actor.PIDSet) *actor.PID {
	l := routees.Len()
	r := rand.Intn(l)
	pid := routees.Get(r)
	return pid
}
