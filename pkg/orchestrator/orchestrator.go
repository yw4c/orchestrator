package orchestrator

import (
	"sync"
)

// 協調者
type IOrchestrator interface {
	// 註冊事務流程
	SetFlows(facade Facade,flow IFlow)
	// 用 facade 領取事務流程
	GetFlow(facade Facade) IFlow

}

// 用 facade 領取事務流程
type Facade string

type IFlow interface {
	// 執行事務流程
	Run(requestID string, requestParam interface{}) (response interface{}, err error)
	// 開始準備接收 MQ 訊息
	ConsumeRollback(rollback *TopicHandlerPair)
}

// Singleton
var instance IOrchestrator
var once sync.Once
func GetInstance() IOrchestrator {
	once.Do(func() {
		c:= &Orchestrator{
			registeredFlows: make(map[Facade]IFlow),
		}
		instance = c
	})
	return instance
}


type Orchestrator struct {
	registeredFlows  map[Facade]IFlow
	syncFlowMu           sync.RWMutex
}

func (o *Orchestrator) SetFlows(facade Facade,flow IFlow) {
	o.syncFlowMu.Lock()
	if sf, ok := flow.(IFlow);ok {
		o.registeredFlows[facade] = sf
	}

	o.syncFlowMu.Unlock()
}

func (o *Orchestrator) GetFlow(facade Facade) IFlow{
	if flow, ok := o.registeredFlows[facade] ; ok {
		return flow
	}
	panic("no such flow : "+facade)
}

