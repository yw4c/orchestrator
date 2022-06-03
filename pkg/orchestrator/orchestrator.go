package orchestrator

import (
	"sync"
)

type IOrchestrator interface {
	RegisterSyncFacade(name FacadeName, facade ISyncFacade)
	RegisterAsyncFacade(name FacadeName, facade IAsyncFacade)
	RegisterThrottlingFacade(name FacadeName, facade IThrottlingFacade)
	GetSyncFacade(name FacadeName) ISyncFacade
	GetAsyncFacade(name FacadeName) IAsyncFacade
	GetThrottlingFacade(name FacadeName) IThrottlingFacade
}

type FacadeName string

type IFacade interface {
	// Start rollback consumer
	ConsumeRollback(rollback *TopicRollbackNodePair)
}

// Singleton
var instance IOrchestrator
var once sync.Once

func GetInstance() IOrchestrator {
	once.Do(func() {
		c := &Orchestrator{
			registeredSyncFlows:       make(map[FacadeName]ISyncFacade),
			registeredAsyncFlows:      make(map[FacadeName]IAsyncFacade),
			registeredThrottlingFlows: make(map[FacadeName]IThrottlingFacade),
		}
		instance = c
	})
	return instance
}

type Orchestrator struct {
	registeredSyncFlows       map[FacadeName]ISyncFacade
	syncFlowMu                sync.RWMutex
	registeredAsyncFlows      map[FacadeName]IAsyncFacade
	asyncFlowMu               sync.RWMutex
	registeredThrottlingFlows map[FacadeName]IThrottlingFacade
	throttlingFlowMu          sync.RWMutex
}

func (o *Orchestrator) RegisterAsyncFacade(name FacadeName, facade IAsyncFacade) {
	o.asyncFlowMu.Lock()
	o.registeredAsyncFlows[name] = facade
	o.asyncFlowMu.Unlock()
}

func (o *Orchestrator) GetAsyncFacade(name FacadeName) IAsyncFacade {
	if flow, ok := o.registeredAsyncFlows[name]; ok {
		return flow
	}
	return nil
}

func (o *Orchestrator) RegisterSyncFacade(name FacadeName, facade ISyncFacade) {
	o.syncFlowMu.Lock()
	o.registeredSyncFlows[name] = facade
	o.syncFlowMu.Unlock()
}

func (o *Orchestrator) GetSyncFacade(name FacadeName) ISyncFacade {
	if flow, ok := o.registeredSyncFlows[name]; ok {
		return flow
	}
	return nil
}

func (o *Orchestrator) RegisterThrottlingFacade(name FacadeName, facade IThrottlingFacade) {
	o.throttlingFlowMu.Lock()
	o.registeredThrottlingFlows[name] = facade
	o.throttlingFlowMu.Unlock()
}

func (o *Orchestrator) GetThrottlingFacade(name FacadeName) IThrottlingFacade {
	if flow, ok := o.registeredThrottlingFlows[name]; ok {
		return flow
	}
	return nil
}
