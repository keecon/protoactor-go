package etcd

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/keecon/protoactor-go/actor"
	"github.com/keecon/protoactor-go/cluster"
	"github.com/keecon/protoactor-go/remote"
	"github.com/stretchr/testify/assert"
)

func newClusterForTest(name string, addr string, cp cluster.ClusterProvider) *cluster.Cluster {
	host, _port, err := net.SplitHostPort(addr)
	if err != nil {
		panic(err)
	}
	port, _ := strconv.Atoi(_port)
	remoteConfig := remote.Configure(host, port)
	config := cluster.Configure(name, cp, nil, remoteConfig)

	system := actor.NewActorSystem()
	c := cluster.New(system, config)
	// use for test without start remote
	c.ActorSystem.ProcessRegistry.Address = addr
	c.MemberList = cluster.NewMemberList(c)
	c.Remote = remote.NewRemote(c.ActorSystem, c.Config.RemoteConfig)

	return c
}

type roleChangedFunc func(RoleType)

func (f roleChangedFunc) OnRoleChanged(role RoleType) {
	f(role)
}

func TestStartMember(t *testing.T) {
	if testing.Short() {
		return
	}

	a := assert.New(t)

	p, err := New()
	a.NoError(err)
	defer func() { _ = p.Shutdown(true) }()

	c := newClusterForTest("test_etcd_provider", "127.0.0.1:8000", p)
	eventstream := c.ActorSystem.EventStream
	ch := make(chan interface{}, 16)

	eventstream.Subscribe(func(m interface{}) {
		if _, ok := m.(*cluster.ClusterTopology); ok {
			ch <- m
		}
	})

	err = p.StartMember(c)
	a.NoError(err)

	select {
	case <-time.After(5 * time.Second):
		a.FailNow("no member joined yet")

	case m := <-ch:
		// member joined
		msg, _ := m.(*cluster.ClusterTopology)

		members := []*cluster.Member{
			{
				// Id:    "test_etcd_provider@127.0.0.1:8000",
				Id:    fmt.Sprintf("test_etcd_provider@%s", c.ActorSystem.ID),
				Host:  "127.0.0.1",
				Port:  8000,
				Kinds: []string{},
			},
		}

		expected := &cluster.ClusterTopology{
			Members:      members,
			Joined:       members,
			Left:         []*cluster.Member{},
			TopologyHash: msg.TopologyHash,
		}
		a.Equal(expected, msg)

	}
}

func TestConcurrentShutdown(t *testing.T) {
	if testing.Short() {
		return
	}

	p, err := New()
	assert.NoError(t, err)
	c := newClusterForTest(t.Name(), "127.0.0.1:8010", p)
	assert.NoError(t, p.StartMember(c))

	var wg sync.WaitGroup
	errors := make(chan error, 10)
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors <- p.Shutdown(true)
		}()
	}
	wg.Wait()
	close(errors)

	for err := range errors {
		assert.NoError(t, err)
	}
}

func TestRoleChangedListenerCanShutdown(t *testing.T) {
	if testing.Short() {
		return
	}

	var provider *Provider
	shutdownResult := make(chan error, 1)
	listener := roleChangedFunc(func(role RoleType) {
		if role == Leader {
			shutdownResult <- provider.Shutdown(true)
		}
	})
	var err error
	provider, err = New()
	assert.NoError(t, err)
	provider.roleChangedListener = listener
	t.Cleanup(func() { _ = provider.Shutdown(true) })
	c := newClusterForTest(t.Name(), "127.0.0.1:8011", provider)
	assert.NoError(t, provider.StartMember(c))

	select {
	case err := <-shutdownResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("role listener deadlocked while shutting down the provider")
	}
}

func TestTopologyListenerCanShutdownDuringStart(t *testing.T) {
	if testing.Short() {
		return
	}

	provider, err := New()
	assert.NoError(t, err)
	c := newClusterForTest(t.Name(), "127.0.0.1:8013", provider)
	shutdownResult := make(chan error, 1)
	var shutdownOnce sync.Once
	subscription := c.ActorSystem.EventStream.Subscribe(func(event interface{}) {
		if _, ok := event.(*cluster.ClusterTopology); !ok {
			return
		}
		shutdownOnce.Do(func() {
			shutdownResult <- provider.Shutdown(true)
		})
	})
	t.Cleanup(func() {
		c.ActorSystem.EventStream.Unsubscribe(subscription)
		_ = provider.Shutdown(true)
	})
	startResult := make(chan error, 1)
	go func() { startResult <- provider.StartMember(c) }()

	select {
	case err := <-shutdownResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("topology listener deadlocked while shutting down the provider")
	}
	select {
	case err := <-startResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("StartMember did not finish after topology listener shutdown")
	}
}

func TestWatcherTopologyListenerCanShutdown(t *testing.T) {
	if testing.Short() {
		return
	}

	provider, err := New()
	assert.NoError(t, err)
	first := newClusterForTest(t.Name(), "127.0.0.1:8014", provider)
	shutdownResult := make(chan error, 1)
	var shutdownOnce sync.Once
	subscription := first.ActorSystem.EventStream.Subscribe(func(event interface{}) {
		topology, ok := event.(*cluster.ClusterTopology)
		if !ok || len(topology.Members) < 2 {
			return
		}
		shutdownOnce.Do(func() {
			shutdownResult <- provider.Shutdown(true)
		})
	})
	t.Cleanup(func() {
		first.ActorSystem.EventStream.Unsubscribe(subscription)
		_ = provider.Shutdown(true)
	})
	assert.NoError(t, provider.StartMember(first))

	secondProvider, err := New()
	assert.NoError(t, err)
	t.Cleanup(func() { _ = secondProvider.Shutdown(true) })
	second := newClusterForTest(t.Name(), "127.0.0.1:8015", secondProvider)
	assert.NoError(t, secondProvider.StartMember(second))

	select {
	case err := <-shutdownResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("watcher topology listener deadlocked while shutting down the provider")
	}
}

func TestStartMemberConcurrentWithShutdown(t *testing.T) {
	if testing.Short() {
		return
	}

	provider, err := New()
	assert.NoError(t, err)
	c := newClusterForTest(t.Name(), "127.0.0.1:8012", provider)
	startResult := make(chan error, 1)
	shutdownResult := make(chan error, 1)
	go func() { startResult <- provider.StartMember(c) }()
	go func() { shutdownResult <- provider.Shutdown(true) }()

	select {
	case <-startResult:
	case <-time.After(5 * time.Second):
		t.Fatal("StartMember did not finish")
	}
	select {
	case err := <-shutdownResult:
		assert.NoError(t, err)
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown did not finish")
	}
}

func TestStartMember_Multiple(t *testing.T) {
	if testing.Short() {
		return
	}

	a := assert.New(t)
	members := []struct {
		cluster string
		host    string
		port    int
	}{
		{"mycluster2", "127.0.0.1", 8001},
		{"mycluster2", "127.0.0.1", 8002},
		{"mycluster2", "127.0.0.1", 8003},
	}

	p := make([]*Provider, len(members))

	var err error

	t.Cleanup(func() {
		for i := range p {
			_ = p[i].Shutdown(true)
		}
	})

	for i, member := range members {
		addr := fmt.Sprintf("%s:%d", member.host, member.port)
		p[i], err = New()
		a.NoError(err)

		c := newClusterForTest(member.cluster, addr, p[i])
		err := p[i].StartMember(c)
		a.NoError(err)
	}

	isNodesEqual := func(nodes []*Node) bool {
		for _, node := range nodes {
			for _, member := range members {
				if node.Host == member.host && node.Port == member.port {
					return true
				}
			}
		}

		return false
	}

	for i := range p {
		nodes, err := p[i].fetchNodes()
		a.NoError(err)
		a.Equal(len(members), len(nodes))
		flag := isNodesEqual(nodes)
		a.Truef(flag, "Member not found - %+v", p[i].self)
	}
}

//func TestUpdateMemberState(t *testing.T) {
//	if testing.Short() {
//		return
//	}
//	assert := assert.New(t)
//
//	p, _ := New()
//	defer p.Shutdown(true)
//
//	c := newClusterForTest("mycluster3", "127.0.0.1:8000", p)
//	err := p.StartMember(c)
//	assert.NoError(err)
//
//	state := cluster.ClusterState{[]string{"yes"}}
//	err = p.UpdateClusterState(state)
//	assert.NoError(err)
//}
//
//func TestUpdateMemberState_DoesNotReregisterAfterShutdown(t *testing.T) {
//	if testing.Short() {
//		return
//	}
//	assert := assert.New(t)
//
//	p, _ := New()
//	c := newClusterForTest("mycluster4", "127.0.0.1:8001", p)
//	err := p.StartMember(c)
//	assert.NoError(err)
//	t.Cleanup(func() {
//		p.Shutdown(true)
//	})
//
//	state := cluster.ClusterState{[]string{"yes"}}
//	err = p.UpdateClusterState(state)
//	assert.NoError(err)
//
//	err = p.Shutdown(true)
//	assert.NoError(err)
//
//	err = p.UpdateClusterState(state)
//	assert.Error(err)
//}
