package zk

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	zookeeper "github.com/go-zookeeper/zk"
	"github.com/stretchr/testify/suite"
)

type ZookeeperTestSuite struct {
	suite.Suite
}

func newClusterForProviderTest(t *testing.T, name, address string, provider cluster.ClusterProvider) *cluster.Cluster {
	t.Helper()
	host, portText, err := net.SplitHostPort(address)
	if err != nil {
		t.Fatal(err)
	}
	port, err := strconv.Atoi(portText)
	if err != nil {
		t.Fatal(err)
	}
	remoteConfig := remote.Configure(host, port)
	config := cluster.Configure(name, provider, nil, remoteConfig)
	system := actor.NewActorSystem()
	c := cluster.New(system, config)
	c.ActorSystem.ProcessRegistry.Address = address
	c.MemberList = cluster.NewMemberList(c)
	c.Remote = remote.NewRemote(c.ActorSystem, c.Config.RemoteConfig)
	return c
}

func (suite *ZookeeperTestSuite) SetupTest() {

}

func (suite *ZookeeperTestSuite) TearDownTest() {
}

func TestZookeeperTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Zookeeper integration test in short mode")
	}
	suite.Run(t, new(ZookeeperTestSuite))
}

func TestProviderConcurrentShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	provider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	remoteConfig := remote.Configure("127.0.0.1", 8020)
	config := cluster.Configure(t.Name(), provider, nil, remoteConfig)
	system := actor.NewActorSystem()
	c := cluster.New(system, config)
	c.ActorSystem.ProcessRegistry.Address = "127.0.0.1:8020"
	c.MemberList = cluster.NewMemberList(c)
	c.Remote = remote.NewRemote(c.ActorSystem, c.Config.RemoteConfig)
	if err := provider.StartMember(c); err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	errors := make(chan error, 10)
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errors <- provider.Shutdown(true)
		}()
	}
	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

func TestRoleChangedListenerCanShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	var provider *Provider
	shutdownResult := make(chan error, 1)
	var err error
	provider, err = New([]string{`localhost:8000`}, WithRoleChangedFunc(func(role RoleType) {
		if role == Leader {
			shutdownResult <- provider.Shutdown(true)
		}
	}))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = provider.Shutdown(true) })
	c := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8021", provider)
	if err := provider.StartMember(c); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-shutdownResult:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("role listener deadlocked while shutting down the provider")
	}
}

func TestTopologyListenerCanShutdownDuringStart(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	provider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	c := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8024", provider)
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
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("topology listener deadlocked while shutting down the provider")
	}
	select {
	case err := <-startResult:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("StartMember did not finish after topology listener shutdown")
	}
}

func TestWatcherTopologyListenerCanShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	provider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	first := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8025", provider)
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
	if err := provider.StartMember(first); err != nil {
		t.Fatal(err)
	}

	secondProvider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = secondProvider.Shutdown(true) })
	second := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8026", secondProvider)
	if err := secondProvider.StartMember(second); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-shutdownResult:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("watcher topology listener deadlocked while shutting down the provider")
	}
}

func TestDisconnectedProviderCannotRegainLeadershipFromStaleNodes(t *testing.T) {
	provider := &Provider{
		done:            make(chan struct{}),
		self:            &Node{ID: "self"},
		roleChangedChan: make(chan RoleType, 2),
		connected:       true,
	}
	provider.role.Store(int32(Leader))
	provider.onEvent(zookeeper.Event{Type: zookeeper.EventSession, State: zookeeper.StateDisconnected})
	provider.updateLeadership([]*Node{provider.self})

	if provider.IsLeader() {
		t.Fatal("disconnected provider regained leadership from a stale member snapshot")
	}
}

func TestStartMemberConcurrentWithShutdown(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	provider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	c := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8022", provider)
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
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Shutdown did not finish")
	}
}

func TestShutdownIgnoresMissingMemberNode(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping ZooKeeper integration test in short mode")
	}

	provider, err := New([]string{`localhost:8000`})
	if err != nil {
		t.Fatal(err)
	}
	c := newClusterForProviderTest(t, t.Name(), "127.0.0.1:8023", provider)
	if err := provider.StartMember(c); err != nil {
		t.Fatal(err)
	}
	if err := provider.conn.Delete(provider.fullpath, -1); err != nil {
		t.Fatal(err)
	}
	if err := provider.Shutdown(true); err != nil {
		t.Fatal(err)
	}
	if provider.fullpath != "" {
		t.Fatalf("member path was not cleared: %s", provider.fullpath)
	}
}

type ClusterAndSystem struct {
	Cluster *cluster.Cluster
	System  *actor.ActorSystem
}

func (cs *ClusterAndSystem) Shutdown() {
	cs.Cluster.Shutdown(true)
}

func (suite *ZookeeperTestSuite) start(name string, opts ...cluster.ConfigOption) *ClusterAndSystem {
	cp, _ := New([]string{`localhost:8000`})
	remoteConfig := remote.Configure("localhost", 0)
	config := cluster.Configure(name, cp, disthash.New(), remoteConfig, opts...)
	system := actor.NewActorSystem()
	c := cluster.New(system, config)
	c.StartMember()
	return &ClusterAndSystem{Cluster: c, System: system}
}

func (suite *ZookeeperTestSuite) TestEmptyExecute() {
	name := `cluster0`
	suite.start(name).Shutdown()
}

func (suite *ZookeeperTestSuite) TestMultiNodes() {
	var actorCount int32
	props := actor.PropsFromFunc(func(ctx actor.Context) {
		switch ctx.Message().(type) {
		case *actor.Started:
			atomic.AddInt32(&actorCount, 1)
		}
	})
	helloKind := cluster.NewKind("hello", props)

	name := `cluster1`
	c1 := suite.start(name, cluster.WithKinds(helloKind))
	defer c1.Shutdown()
	c2 := suite.start(name, cluster.WithKinds(helloKind))
	defer c2.Shutdown()
	c1.Cluster.Get(`a1`, `hello`)
	c2.Cluster.Get(`a2`, `hello`)
	suite.Require().Eventually(func() bool {
		return atomic.LoadInt32(&actorCount) == 2
	}, 10*time.Second, 5*time.Millisecond)
	suite.Assert().Equal(2, c1.Cluster.MemberList.Members().Len(), "Expected 2 members in the cluster")
	suite.Assert().Equal(2, c2.Cluster.MemberList.Members().Len(), "Expected 2 members in the cluster")
}
