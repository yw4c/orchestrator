package orchestrator

import (
	"orchestrator/pkg/pkgerror"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
)

// ISyncFacade
type ISyncFacade interface {
	// Register Nodes in serialize
	Use(syncNode SyncNode) ISyncFacade
	// Execute facade in request
	Run(requestID string, requestParam interface{}) (response interface{}, err error)
	IFacade
}

// SyncNode
type SyncNode func(requestID string, ctx *Context) error

// FacadeReqContextKey: Request is stored in context, give it a key to get request payload
type FacadeReqContextKey string

// FacadeRespContextKey: Response is stored in context, give it a key to get response data
type FacadeRespContextKey string

type SyncFacade struct {
	nodes         []SyncNode
	rollbackTopic Topic
}

func (s *SyncFacade) ConsumeRollback(rollback *TopicRollbackNodePair) {
	GetMQInstance().ConsumeRollback(rollback.Topic, rollback.Node)
	s.rollbackTopic = rollback.Topic
}

func (s *SyncFacade) Use(syncNode SyncNode) ISyncFacade {
	s.nodes = append(s.nodes, syncNode)
	return s
}

func (s *SyncFacade) Run(requestID string, requestParam interface{}) (response interface{}, err error) {

	if s.rollbackTopic == "" {
		return nil, eris.Wrap(pkgerror.ErrInternalError, "rollback topic is not set")
	}

	context := &Context{}
	context.setPbReq(requestParam)

	for _, v := range s.nodes {
		if err := v(requestID, context); err != nil {

			// if error occur，send Rollback Topic to MQ
			log.Error().
				Str("Rollback Topic", string(s.rollbackTopic)).
				Str("Request", requestID).
				Msg("Flow went wrong, Publishing Rollback topic")
			rollback(s.rollbackTopic, requestID)

			return nil, err
		}
	}

	// Flow finished successfully
	resp, _ := context.Get(pbRespCtxKey)

	return resp, nil
}

func NewSyncFacade() ISyncFacade {
	return &SyncFacade{
		nodes: nil,
	}
}
