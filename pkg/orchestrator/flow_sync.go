package orchestrator

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"orchestrator/pkg/ctx"
	"orchestrator/pkg/pkgerror"
)

// 同步的事務流程
type ISyncFlow interface {
	Use(syncHandler SyncHandler) ISyncFlow
	IFlow
}

// 同步的事務節點
type SyncHandler func(requestID string, ctx *ctx.Context) error


type SyncFlow struct {
	handlers      []SyncHandler
	rollbackTopic Topic
}

func (s *SyncFlow) ConsumeRollback(rollback *TopicHandlerPair) {
	GetMQInstance().ListenAndConsume(rollback.Topic, rollback.AsyncHandler)
	s.rollbackTopic = rollback.Topic
}

func (s *SyncFlow) Use(syncHandler SyncHandler) ISyncFlow {
	s.handlers = append(s.handlers, syncHandler)
	return s
}

func (s *SyncFlow) Run(requestID string, requestParam interface{}) (response interface{}, err error) {

	if s.rollbackTopic == "" {
		return nil, eris.Wrap(pkgerror.ErrInternalError, "rollback topic is not set")
	}

	context := &ctx.Context{}
	context.Set(BookingSyncPbReq, requestParam)

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
	resp, exist := context.Get(BookingSyncPbResp)
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

// 事務節點間，傳遞 Context 資料的 key
const (
	BookingSyncPbReq = "bookingSyncPbReq"
	BookingSyncPbResp = "bookingSyncPbResp"
)