package orchestrator

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/ctx"
	"orchestrator/pkg/pkgerror"
)

// 同步的事務流程
type ISyncFlow interface {
	// 將䩞點依序加入
	Use(syncHandler SyncHandler) ISyncFlow
	// 執行事務流程
	Run(requestID string, requestParam interface{}, reqKey FlowContextKeyReq, respKey FlowContextKeyResp)(response interface{}, err error)
	IFlow
}

// 同步的事務節點
type SyncHandler func(requestID string, ctx *ctx.Context) error
// 從 Context 取得 Request Value 的 Key
type FlowContextKeyReq string
// 從 Context 取得 Response Value 的 Key
type FlowContextKeyResp string

type SyncFlow struct {
	handlers      []SyncHandler
	rollbackTopic Topic
}

func (s *SyncFlow) ConsumeRollback(rollback *TopicRollbackHandlerPair) {
	GetMQInstance().ConsumeRollback(rollback.Topic, rollback.Handler)
	s.rollbackTopic = rollback.Topic
}

func (s *SyncFlow) Use(syncHandler SyncHandler) ISyncFlow {
	s.handlers = append(s.handlers, syncHandler)
	return s
}

func (s *SyncFlow) Run(requestID string, requestParam interface{}, reqKey FlowContextKeyReq, respKey FlowContextKeyResp) (response interface{}, err error) {

	if s.rollbackTopic == "" {
		return nil, eris.Wrap(pkgerror.ErrInternalError, "rollback topic is not set")
	}

	context := &ctx.Context{}
	context.Set(string(reqKey), requestParam)

	for _,v := range s.handlers {
		if err := v(requestID, context); err != nil {

			// 發生錯誤，發送 Rollback Topic 給 MQ
			log.Error().
				Str("Rollback Topic", string(s.rollbackTopic)).
				Str("Request", requestID).
				Msg("Flow went wrong, Publishing Rollback topic")
			rollback(s.rollbackTopic, requestID)

			return nil, err
		}
	}

	// Flow finished successfully
	resp, exist := context.Get(string(respKey))
	if !exist {
		return nil, eris.Wrap(pkgerror.ErrInternalError, "Can not get DTO from Context")
	}
	return  resp, nil
}

func NewSyncFlow() ISyncFlow {
	return &SyncFlow{
		handlers: nil,
	}
}

