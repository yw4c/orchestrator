package orchestrator

import (
	"sync"
)

// 協調者
type IOrchestrator interface {
	// 註冊事務流程
	SetSyncFlows(facade Facade, flow ISyncFlow)
	SetAsyncFlows(facade Facade, flow IAsyncFlow)
	// 用 facade 領取事務流程
	GetSyncFlow(facade Facade) ISyncFlow
	GetAsyncFlow(facade Facade) IAsyncFlow
}

// 用 facade 領取事務流程
type Facade string

type IFlow interface {
	// 開始準備接收 MQ 訊息
	ConsumeRollback(rollback *TopicHandlerPair)
}

// Singleton
var instance IOrchestrator
var once sync.Once
func GetInstance() IOrchestrator {
	once.Do(func() {
		c:= &Orchestrator{
			registeredSyncFlows: make(map[Facade]ISyncFlow),
			registeredAsyncFlows: make(map[Facade]IAsyncFlow),
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

