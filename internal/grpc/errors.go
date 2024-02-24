package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrMissingRequiredField returns an error message format for missing required fields.
func ErrMissingRequiredField(field string) error {
	return status.Errorf(codes.InvalidArgument, "error: missing required field '%s'", field)
}
