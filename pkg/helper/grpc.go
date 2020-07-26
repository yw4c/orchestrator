package helper

import (
	"context"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc/metadata"
	"orchestrator/pkg/pkgerror"
)

func GetRequestID(ctx context.Context) (string, error){
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", eris.Wrap(pkgerror.ErrInternalServerError, "Read Metadata Fail")
	}
	requestID := md.Get("x-request-id")
	if len(requestID) == 0 {
		return "", eris.Wrap(pkgerror.ErrInvalidHeaderValue, "Missing X-Request-ID on header")
	}
	return requestID[0], nil
}
