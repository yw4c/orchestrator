package orchestrator

import (
	"sync"
)

// 協調者
type IOrchestrator interface {
	// 註冊事務流程
	SetSyncFlows(facade Facade, flow ISyncFlow)
	SetAsyncFlows(facade Facade, flow IAsyncFlow)
	SetThrottlingFlows(facade Facade, flow IThrottlingFlow)
	// 用 facade 領取事務流程
	GetSyncFlow(facade Facade) ISyncFlow
	GetAsyncFlow(facade Facade) IAsyncFlow
	GetThrottlingFlows(facade Facade) IThrottlingFlow
}

// 用 facade 領取事務流程
type Facade string

type IFlow interface {
	// 開始準備接收 MQ 訊息
	ConsumeRollback(rollback *TopicRollbackHandlerPair)
}

// Singleton
var instance IOrchestrator
var once sync.Once
func GetInstance() IOrchestrator {
	once.Do(func() {
		c:= &Orchestrator{
			registeredSyncFlows: make(map[Facade]ISyncFlow),
			registeredAsyncFlows: make(map[Facade]IAsyncFlow),
			registeredThrottlingFlows: make(map[Facade]IThrottlingFlow),
		}
		instance = c
	})
	return instance
}


type Orchestrator struct {
	registeredSyncFlows map[Facade]ISyncFlow
	syncFlowMu          sync.RWMutex
	registeredAsyncFlows map[Facade]IAsyncFlow
	asyncFlowMu          sync.RWMutex
	registeredThrottlingFlows map[Facade]IThrottlingFlow
	throttlingFlowMu sync.RWMutex
}

func (o *Orchestrator) SetAsyncFlows(facade Facade, flow IAsyncFlow) {
	o.asyncFlowMu.Lock()
	o.registeredAsyncFlows[facade] = flow
	o.asyncFlowMu.Unlock()
}

func (o *Orchestrator) GetAsyncFlow(facade Facade) IAsyncFlow {
	if flow, ok := o.registeredAsyncFlows[facade] ; ok {
		return flow
	}
	return nil
}

func (o *Orchestrator) SetSyncFlows(facade Facade, flow ISyncFlow) {
	o.syncFlowMu.Lock()
	o.registeredSyncFlows[facade] = flow
	o.syncFlowMu.Unlock()
}

func (o *Orchestrator) GetSyncFlow(facade Facade) ISyncFlow {
	if flow, ok := o.registeredSyncFlows[facade] ; ok {
		return flow
	}
	return nil
}

func (o *Orchestrator) SetThrottlingFlows(facade Facade, flow IThrottlingFlow) {
	o.throttlingFlowMu.Lock()
	o.registeredThrottlingFlows[facade] = flow
	o.throttlingFlowMu.Unlock()
}

func (o *Orchestrator) GetThrottlingFlows(facade Facade) IThrottlingFlow {
	if flow, ok := o.registeredThrottlingFlows[facade] ; ok {
		return flow
	}
	return nil
}