package pkgerror

import (
	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"net/http"
)

var Format = eris.NewDefaultStringFormat(eris.FormatOptions{
	InvertOutput: true, // flag that inverts the error output (wrap errors shown first)
	WithTrace: true,    // flag that enables stack trace output
	InvertTrace: true,  // flag that inverts the stack trace output (top of call stack shown first)
})


func SetGRPCErrorResp(requestID string, err error) error {

	// log tracked errors
	formattedStr := eris.ToCustomString(err, Format)

	log.Error().
		Str("track", formattedStr).
		Str("x-request-id", requestID).
		Str("error", err.Error()).
		Msg("gRPC Error Response")


	causeErr := eris.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		_err = &_error{
			Status:   http.StatusInternalServerError,
			Code:     "Unknown",
			GRPCCode: codes.Internal,
			Message:  err.Error(),
			Details:  nil,
		}
	}
	detail := map[string]interface{}{
		"x-request-id": requestID,
	}
	_err.SetDetails(detail)

	return convertProtoErr(_err)

}
