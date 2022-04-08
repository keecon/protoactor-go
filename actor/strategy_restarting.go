package actor

func NewRestartingStrategy() SupervisorStrategy {
	return &restartingStrategy{}
}

type restartingStrategy struct{}

var _ SupervisorStrategy = &restartingStrategy{}

func (strategy *restartingStrategy) HandleFailure(actorSystem *ActorSystem, supervisor Supervisor, child *PID, rs *RestartStatistics, reason interface{}, message interface{}) {
	// always restart
	logFailure(actorSystem, child, reason, RestartDirective)
	supervisor.RestartChildren(child)
}
