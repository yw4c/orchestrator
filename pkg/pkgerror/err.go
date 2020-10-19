package pkgerror

import (
	"encoding/json"
	"net/http"
	"strings"

	errors "github.com/rotisserie/eris"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// ErrBadRequest                       =  &_error{Code: "400000", Message: http.StatusText(http.StatusBadRequest), Status: http.StatusBadRequest}
	ErrInvalidInput          = &_error{Code: "400001", Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidAmount         = &_error{Code: "400002", Message: "This amount is not allow ", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrSuccessRate           = &_error{Code: "400003", Message: "This SuccessRate is not allow ", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidHeaderValue    = &_error{Code: "400004", Message: "The value provided for one of the HTTP headers was invalid .", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidTimestamp      = &_error{Code: "400005", Message: "Invaild timestamp", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrMissingRequiredHeader = &_error{Code: "400017", Message: "A required HTTP header was not specified.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInternalDataNotSync   = &_error{Code: "400041", Message: "The specified data not sync in system, please wait a moment.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrNotMatchSetting       = &_error{Code: "400087", Message: "The specified data not match setting, please adjust your inputs.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}

	ErrUnauthorized                = &_error{Code: "401001", Message: http.StatusText(http.StatusUnauthorized), Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrInvalidAuthenticationInfo   = &_error{Code: "401001", Message: "The authentication information was not provided in the correct Format. Verify the value of Authorization header.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrInvalidSign                 = &_error{Code: "401001", Message: "Invaild sign", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}
	ErrUsernameOrPasswordIncorrect = &_error{Code: "401002", Message: "Username or Password is incorrect.", Status: http.StatusUnauthorized, GRPCCode: codes.Unauthenticated}

	// ErrForbidden                                   =  &_error{Code: "403000", Message: http.StatusText(http.StatusForbidden), Status: http.StatusForbidden}
	ErrAccountIsDisabled                           = &_error{Code: "403001", Message: "The specified account is disabled.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrAuthenticationFailed                        = &_error{Code: "403002", Message: "Server failed to authenticate the request. Make sure the value of the Authorization header is formed correctly including the signature.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrNotAllowed                                  = &_error{Code: "403003", Message: "The request is understood, but it has been refused or access is not allowed.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpExpired                                  = &_error{Code: "403004", Message: "OTP is expired.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrInsufficientAccountPermissionsWithOperation = &_error{Code: "403005", Message: "The account being accessed does not have sufficient permissions to execute this operation.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOptRequired                                 = &_error{Code: "403007", Message: "OTP Binding is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpAuthorizationRequired                    = &_error{Code: "403008", Message: "Two-factor authorization is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrOtpIncorrect                                = &_error{Code: "403009", Message: "OTP is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrResetPasswordRequired                       = &_error{Code: "403010", Message: "Reset password is required.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrResetOTPAccountNameIncorrect                = &_error{Code: "403011", Message: "Reset otp failed,requested otp name is incorrect.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrIPNotAllowed                                = &_error{Code: "403012", Message: "IP is invalid", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}

	// ErrNotFound         =  &_error{Code: "404000", Message: http.StatusText(http.StatusNotFound), Status: http.StatusNotFound}
	ErrResourceNotFound = &_error{Code: "404001", Message: "The specified resource does not exist.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrAccountNotFound  = &_error{Code: "404002", Message: "cant find any account.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrPageNotFound     = &_error{Code: "404003", Message: "Page Not Fount.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrTimeout            = &_error{Code: "408000", Message: "Request timeout", Status: http.StatusRequestTimeout, GRPCCode: codes.Canceled}
	ErrConflict              = &_error{Code: "409000", Message: http.StatusText(http.StatusConflict), Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrAccountAlreadyExists  = &_error{Code: "409001", Message: "The specified account already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrAccountBeingCreated   = &_error{Code: "409002", Message: "The specified account is in the process of being created.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrResourceAlreadyExists = &_error{Code: "409004", Message: "The specified resource already exists.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}


	ErrInternalServerError = &_error{Code: "500000", Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
	ErrInternalError       = &_error{Code: "500001", Message: "The server encountered an internal error. Please retry the request.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}

	// Internal usage
	ErrOrderNoChange = &_error{Code: "500002", Message: "Order status No change", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)

type _error struct {
	Status   int                    `json:"-"`
	Code     string                 `json:"code"`
	GRPCCode codes.Code             `json:"grpccode"`
	Message  string                 `json:"message"`
	Details  map[string]interface{} `json:"details"`
}

func (e *_error) Error() string {
	var b strings.Builder
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(e.Code)
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// SetDetails set details as you wish =)
func (e *_error) SetDetails(details map[string]interface{}) {
	e.Details = details
	return
}

// CompareErrorCode 比較兩個錯誤代碼是否一致
func CompareErrorCode(errA error, errB error) bool {
	var aErr, bErr *_error
	if err, exists := errors.Cause(errA).(*_error); exists {
		aErr = err
	}
	if err, exists := errors.Cause(errB).(*_error); exists {
		bErr = err
	}
	if aErr.Code == bErr.Code {
		return true
	}
	return false
}

//convertProtoErr Convert _error to grpc error
func convertProtoErr(err error) error {
	if err == nil {
		return nil
	}
	causeErr := errors.Cause(err)
	_err, ok := causeErr.(*_error)
	if !ok {
		return status.Error(ErrInternalError.GRPCCode, ErrInternalError.Message)
	}
	return status.Error(_err.GRPCCode, _err.Message)
}

func switchCode(s *status.Status) *_error {
	httperr := ErrInternalError
	switch s.Code() {
	case codes.Unknown:
		httperr = ErrInternalError
	case codes.InvalidArgument:
		httperr = ErrInvalidInput
	case codes.NotFound:
		httperr = ErrResourceNotFound
	case codes.AlreadyExists:
		httperr = ErrResourceAlreadyExists
	case codes.PermissionDenied:
		httperr = ErrNotAllowed
	case codes.Unauthenticated:
		httperr = ErrUnauthorized
	case codes.OutOfRange:
		httperr = ErrInvalidInput
	case codes.Internal:
		httperr = ErrInternalError
	case codes.DataLoss:
		httperr = ErrInternalError
	}
	httperr.Message = s.Message()
	return httperr
}

//ConvertHttpErr Convert  grpc error to _error
func ConvertFromGrpc(err error) *_error {
	if err == nil {
		return nil
	}
	s := status.Convert(err)
	if s == nil {
		return ErrInternalError
	}
	interErr := _error{}
	jerr := json.Unmarshal([]byte(s.Message()), &interErr)
	if jerr != nil {
		return switchCode(s)
	}
	return &interErr
}
