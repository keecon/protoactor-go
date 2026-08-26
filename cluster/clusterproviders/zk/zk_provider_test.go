package zk

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/asynkron/protoactor-go/cluster"
	"github.com/asynkron/protoactor-go/cluster/identitylookup/disthash"
	"github.com/asynkron/protoactor-go/remote"
	"github.com/stretchr/testify/suite"
)

type ZookeeperTestSuite struct {
	suite.Suite
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
